package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	lg := lg.With().
		Str(lkHandler, "MessageCreate").
		Str(lkMID, m.ID).
		Str(lkGuild, m.GuildID).
		Str(lkCh, m.ChannelID).
		Str(lkUsr, m.Author.ID).
		Str(lkName, m.Author.Username).
		Logger()

	// ignore bot messages
	if m.Author.Bot {
		return
	}

	// check DM
	isDM := false
	if m.GuildID == "" {
		isDM = true
	}

	// check mention
	isMention := false
	for _, u := range m.Mentions {
		if u.ID == s.State.Application.ID {
			isMention = true
		}
	}

	// check thread
	isThread := false
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get channel")
	}
	if ch.IsThread() {
		isThread = true
	}

	// check ref
	hasRef := false
	if m.MessageReference != nil {
		hasRef = true
	}

	// check function
	funcName := ""
	if !isDM && isMention && hasRef && containsReactionStats(m.Content) {
		funcName = FuncMessageCreateReactionStats
	}

	lg.Debug().Bool(lkDM, isDM).Msgf("isMention=%v isThread=%v hasRef=%v content=%s", isMention, isThread, hasRef, m.Content)

	if funcName == "" {
		return
	}

	lg.Info().Str(lkFunc, funcName).Bool(lkDM, isDM).Msg("OnMessageCreate")
	switch funcName {
	case FuncMessageCreateReactionStats:
		handleMessageCreateReactionStats(s, m)
	default:
		return
	}
}

func containsReactionStats(s string) bool {
	ss := strings.ToLower(s)
	return strings.Contains(ss, "集計") || strings.Contains(ss, "sum") || strings.Contains(ss, "stats")
}

func handleMessageCreateReactionStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	lg := lg.With().
		Str(lkFunc, FuncMessageCreateReactionStats).
		Str(lkMID, m.ID).
		Str(lkGuild, m.GuildID).
		Str(lkCh, m.ChannelID).
		Str(lkUsr, m.Author.ID).
		Str(lkName, m.Author.Username).
		Logger()

	lg.Info().Msgf("ReactionStats: called")

	if err := s.ChannelTyping(m.ChannelID); err != nil {
		lg.Error().Err(err).Msg("could not send typing")
	}

	parentMsg, err := s.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get parent message")
		return
	}

	userEmojis := map[string]map[string]bool{} // user -> emojis
	emojiUsers := map[string]map[string]bool{} // emoji -> users

	for _, u := range parentMsg.Mentions {
		userEmojis[u.Username] = map[string]bool{}
	}
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		lg.Error().Err(err).Msg("could not get GuildMembers")
		return
	}
	for _, role := range parentMsg.MentionRoles {
		for _, member := range members {
			if member.User.Bot {
				continue // skip bots (because bots don't skip what they need to do)
			}
			for _, memberRole := range member.Roles {
				if memberRole == role {
					userEmojis[member.User.Username] = map[string]bool{}
				}
			}
		}
	}

	for _, r := range parentMsg.Reactions {
		users, err := s.MessageReactions(m.ChannelID, parentMsg.ID, r.Emoji.APIName(), 100, "", "")
		if err != nil {
			lg.Error().Err(err).Msg("could not get messagereactions")
			continue
		}
		for _, u := range users {
			if _, ok := userEmojis[u.Username]; !ok {
				userEmojis[u.Username] = map[string]bool{}
			}
			userEmojis[u.Username][r.Emoji.APIName()] = true
			if _, ok := emojiUsers[r.Emoji.APIName()]; !ok {
				emojiUsers[r.Emoji.APIName()] = map[string]bool{}
			}
			emojiUsers[r.Emoji.APIName()][u.Username] = true
		}
	}

	// drop self
	delete(userEmojis, s.State.User.Username)
	for _, e := range emojiUsers {
		delete(e, s.State.User.Username)
	}
	// drop special emojis
	guildEmojis, err := s.GuildEmojis(m.GuildID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get GuildEmojis")
	}
	specialEmojis := []string{"🤖", getGuildEmojiAPINameByName(guildEmojis, ReactionMa), getGuildEmojiAPINameByName(guildEmojis, ReactionConbu01)}
	for _, u := range userEmojis {
		for _, se := range specialEmojis {
			delete(u, se)
		}
	}
	for _, se := range specialEmojis {
		delete(emojiUsers, se)
	}

	var emojiList []string
	for emoji := range emojiUsers {
		emojiList = append(emojiList, emoji)
	}

	var table strings.Builder
	table.WriteString("集計しました（2分間表示）\n")
	for _, emoji := range emojiList {
		table.WriteString(emoji2msg(emoji))
		table.WriteString(" | ")
	}
	table.WriteString("😀")
	table.WriteString("\n")
	for i := 0; i < len(emojiList)*5+10; i++ {
		table.WriteRune('-')
	}
	table.WriteString("\n")
	for user, emojis := range userEmojis {
		for _, emoji := range emojiList {
			if emojis[emoji] {
				table.WriteString("✅")
			} else {
				table.WriteString("➖")
			}
			table.WriteString(" | ")
		}
		table.WriteString(fmt.Sprintf("`%s`", user))
		table.WriteString("\n")
	}
	// TODO: send silent message
	reply, err := s.ChannelMessageSendReply(m.ChannelID, table.String(), m.Reference())
	if err != nil {
		lg.Error().Err(err).Msg("could not send reply")
	}
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		lg.Error().Err(err).Msg("could not delete summarization request message")
	}
	time.AfterFunc(time.Second*120, func() {
		lg.Info().Msgf("delete (AfterFunc)")
		if err := s.ChannelMessageDelete(m.ChannelID, reply.ID); err != nil {
			lg.Error().Err(err).Msg("could not delete reply (AfterFunc)")
		}
	})
}

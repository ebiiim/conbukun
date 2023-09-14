package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	EmojiConbu01 = "icon_conbu01"
	EmojiConbu02 = "icon_conbu02"
	EmojiMa      = "ma"
)

var (
	ReactionAddHandlers = map[string]func(s *discordgo.Session, r *discordgo.MessageReactionAdd){}

	emojisReactionAddReactionStats    = []string{"🤖", EmojiConbu01, EmojiMa}
	emojisReactionAddReactionRequired = []string{"👀", EmojiConbu02}
)

func init() {
	lg.Debug().Msgf("init: register ReactionAddHandlers")
	for _, emoji := range emojisReactionAddReactionStats {
		ReactionAddHandlers[emoji] = handleReactionAddReactionStats
	}
	for _, emoji := range emojisReactionAddReactionRequired {
		ReactionAddHandlers[emoji] = handleReactionAddReactionRequired
	}
}

func OnMessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	lg := lg.With().
		Str(lkHandler, "MessageReactionAdd").
		Str(lkGuild, r.GuildID).
		Str(lkCh, r.ChannelID).
		Str(lkMID, r.MessageID).
		Str(lkUsr, r.UserID).
		Str(lkEmoji, r.Emoji.Name).
		Str(lkEmojiA, r.Emoji.APIName()).
		Logger()

	lg.Debug().Msgf("OnMessageReactionAdd")
	if h, ok := ReactionAddHandlers[r.Emoji.Name]; ok {
		lg.Info().Msgf("OnMessageReactionAdd")
		h(s, r)
	}

}

func handleReactionAddReactionStats(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.GuildID == "" {
		lg.Debug().Msgf("ReactionAddReactionStats: return as no GuildID (this message is a DM)")
		return
	}

	lg := lg.With().
		Str(lkFunc, FuncReactionAddReactionStats).
		Str(lkGuild, r.GuildID).
		Str(lkCh, r.ChannelID).
		Str(lkMID, r.MessageID).
		Str(lkUsr, r.UserID).
		Str(lkEmoji, r.Emoji.Name).
		Str(lkEmojiA, r.Emoji.APIName()).
		Logger()

	lg.Info().Msgf("ReactionAddReactionStats: called")

	if err := s.ChannelTyping(r.ChannelID); err != nil {
		lg.Error().Err(err).Msg("could not send typing")
	}

	parentMsg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get message")
		return
	}

	userEmojis := map[string]map[string]bool{} // user -> emojis
	emojiUsers := map[string]map[string]bool{} // emoji -> users

	for _, u := range parentMsg.Mentions {
		userEmojis[u.Username] = map[string]bool{}
	}
	members, err := s.GuildMembers(r.GuildID, "", 1000)
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

	for _, rt := range parentMsg.Reactions {
		users, err := s.MessageReactions(r.ChannelID, parentMsg.ID, rt.Emoji.APIName(), 100, "", "")
		if err != nil {
			lg.Error().Err(err).Msg("could not get messagereactions")
			continue
		}
		for _, u := range users {
			if _, ok := userEmojis[u.Username]; !ok {
				userEmojis[u.Username] = map[string]bool{}
			}
			userEmojis[u.Username][rt.Emoji.APIName()] = true
			if _, ok := emojiUsers[rt.Emoji.APIName()]; !ok {
				emojiUsers[rt.Emoji.APIName()] = map[string]bool{}
			}
			emojiUsers[rt.Emoji.APIName()][u.Username] = true
		}
	}

	// drop self
	delete(userEmojis, s.State.User.Username)
	for _, e := range emojiUsers {
		delete(e, s.State.User.Username)
	}
	// drop special emojis
	guildEmojis, err := s.GuildEmojis(r.GuildID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get GuildEmojis")
	}
	specialEmojis := []string{}
	for _, emoji := range emojisReactionAddReactionStats {
		specialEmojis = append(specialEmojis, getGuildEmojiAPINameByName(guildEmojis, emoji))
	}
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
	msg, err := sendSilentMessage(s, r.ChannelID, &discordgo.MessageSend{
		Content: table.String(),
	})
	if err != nil {
		lg.Error().Err(err).Msg("could not send msg")
	}
	time.AfterFunc(time.Second*120, func() {
		lg.Info().Msgf("delete (AfterFunc)")
		if err := s.ChannelMessageDelete(r.ChannelID, msg.ID); err != nil {
			lg.Error().Err(err).Msg("could not delete msg (AfterFunc)")
		}
	})
}

func handleReactionAddReactionRequired(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.GuildID == "" {
		lg.Debug().Msgf("ReactionAddReactionRequired: return as no GuildID (this message is a DM)")
		return
	}

	lg := lg.With().
		Str(lkFunc, FuncReactionAddReactionRequired).
		Str(lkGuild, r.GuildID).
		Str(lkCh, r.ChannelID).
		Str(lkMID, r.MessageID).
		Str(lkUsr, r.UserID).
		Str(lkEmoji, r.Emoji.Name).
		Str(lkEmojiA, r.Emoji.APIName()).
		Logger()

	lg.Info().Msgf("ReactionAddReactionRequired: called")

	if err := s.ChannelTyping(r.ChannelID); err != nil {
		lg.Error().Err(err).Msg("could not send typing")
	}

	// Fetch data.
	parentMsg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get message")
		return
	}
	members, err := s.GuildMembers(r.GuildID, "", 1000)
	if err != nil {
		lg.Error().Err(err).Msg("could not get GuildMembers")
		return
	}

	// Parse mentions.
	// NOTE: channel mentions are not supported
	mentionedUserIDs := map[string]struct{}{}
	for _, u := range parentMsg.Mentions { // normal mentions
		if u.Bot {
			continue
		}
		mentionedUserIDs[u.ID] = struct{}{}
	}
	for _, role := range parentMsg.MentionRoles { // role mentions
		for _, member := range members {
			if member.User.Bot {
				continue // skip bots (because bots don't skip what they need to do)
			}
			for _, memberRole := range member.Roles {
				if memberRole == role {
					mentionedUserIDs[member.User.ID] = struct{}{}
				}
			}
		}
	}
	if parentMsg.MentionEveryone { // everyone mentions
		for _, member := range members {
			mentionedUserIDs[member.User.ID] = struct{}{}
		}
	}

	// Parse reactions.
	reactedUserIDs := map[string]struct{}{}
	for _, rt := range parentMsg.Reactions {
		skip := false
		for _, excludedEmoji := range emojisReactionAddReactionRequired {
			if rt.Emoji.Name == excludedEmoji {
				skip = true
			}
		}
		if skip {
			continue
		}
		users, err := s.MessageReactions(r.ChannelID, parentMsg.ID, rt.Emoji.APIName(), 100, "", "")
		if err != nil {
			lg.Error().Err(err).Msg("could not get messagereactions")
			continue
		}
		for _, u := range users {
			if u.Bot {
				continue
			}
			reactedUserIDs[u.ID] = struct{}{}
		}
	}

	// Generate the response.
	var sb strings.Builder
	sb.WriteString("集計しました（2分間表示）\n")
	sb.WriteString(fmt.Sprintf("### リアクションしたメンバー %d/%d\n", len(reactedUserIDs), len(mentionedUserIDs)))
	sb.WriteString(("### 🔔リマインダー🔔\n"))
	for id, _ := range mentionedUserIDs {
		if _, ok := reactedUserIDs[id]; ok {
			continue
		}
		sb.WriteString(fmt.Sprintf("`%s` ", id2nick(members, id)))
	}

	msg, err := sendSilentMessage(s, r.ChannelID, &discordgo.MessageSend{
		Content: sb.String(),
	})
	if err != nil {
		lg.Error().Err(err).Msg("could not send msg")
	}
	time.AfterFunc(time.Second*120, func() {
		lg.Info().Msgf("delete (AfterFunc)")
		if err := s.ChannelMessageDelete(r.ChannelID, msg.ID); err != nil {
			lg.Error().Err(err).Msg("could not delete msg (AfterFunc)")
		}
	})
}

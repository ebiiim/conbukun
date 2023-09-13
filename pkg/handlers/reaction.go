package handlers

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

const (
	ReactionConbu01 = "icon_conbu01"
	ReactionConbu02 = "icon_conbu02"
	ReactionMa      = "ma"
)

var (
	ReactionAddHandlers = map[string]func(s *discordgo.Session, r *discordgo.MessageReactionAdd){
		"🤖":             handleReactionAddReactionStats,
		ReactionConbu01: handleReactionAddReactionStats,
		ReactionMa:      handleReactionAddReactionStats,
	}
)

func OnMessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	lg = lg.With().
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

func emoji2msg(emojiAPIName string) string {
	if emojiAPIName == "" {
		return ""
	} else if utf8.RuneCountInString(emojiAPIName) == 1 {
		return emojiAPIName // normal emojis
	} else {
		return fmt.Sprintf("<:%s>", emojiAPIName) // e.g. <:ma:1151171171799269476>
	}
}

func handleReactionAddReactionStats(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.GuildID == "" {
		lg.Debug().Msgf("ReactionAddReactionStats: return as no GuildID (this message is a DM)")
		return
	}

	lg = lg.With().
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

	// drop special emojis
	guildEmojis, err := s.GuildEmojis(r.GuildID)
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
	msg, err := sendSilentMessage(s, r.ChannelID, &discordgo.MessageSend{
		Content: table.String(),
		TTS:     false,
	})
	// msg, err := s.ChannelMessageSend(r.ChannelID, table.String())
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

type messageSend struct {
	discordgo.MessageSend `json:",inline"`
	Flags                 discordgo.MessageFlags `json:"flags,omitempty"`
}

func sendSilentMessage(s *discordgo.Session, channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	// TODO: Remove this when compatibility is not required.
	if data.Embed != nil {
		if data.Embeds == nil {
			data.Embeds = []*discordgo.MessageEmbed{data.Embed}
		} else {
			err = fmt.Errorf("cannot specify both Embed and Embeds")
			return
		}
	}

	for _, embed := range data.Embeds {
		if embed.Type == "" {
			embed.Type = "rich"
		}
	}
	endpoint := discordgo.EndpointChannelMessages(channelID)

	// TODO: Remove this when compatibility is not required.
	files := data.Files
	if data.File != nil {
		if files == nil {
			files = []*discordgo.File{data.File}
		} else {
			err = fmt.Errorf("cannot specify both File and Files")
			return
		}
	}

	data2 := messageSend{
		MessageSend: *data,
		Flags:       MessageFlagsSilent,
	}

	var response []byte
	if len(files) > 0 {
		// NOTE: won't support
		return nil, fmt.Errorf("not supported: len(files) > 0")

		// contentType, body, encodeErr := discordgo.MultipartBodyWithJSON(data, files)
		// if encodeErr != nil {
		// 	return st, encodeErr
		// }

		// response, err = s.request("POST", endpoint, contentType, body, endpoint, 0, options...)
	} else {
		response, err = s.RequestWithBucketID("POST", endpoint, data2, endpoint, options...)
	}
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

func unmarshal(data []byte, v interface{}) error {
	err := discordgo.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%w: %s", discordgo.ErrJSONUnmarshal, err)
	}

	return nil
}

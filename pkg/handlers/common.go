package handlers

import (
	"fmt"
	"math/rand"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "dev" // main.go injects this value

var lg zerolog.Logger = log.With().Str("component", "conbukun/pkg/handlers").Logger()

// log keys
const (
	lkHandler = "handler"
	lkCmd     = "command"
	lkFunc    = "function"
	lkGuild   = "guild"
	lkCh      = "channel"
	lkUsr     = "user"
	lkName    = "username"
	lkDM      = "dm"
	lkIID     = "interaction_id"
	lkMID     = "message_id"
	lkEmoji   = "emoji"
	lkEmojiA  = "emoji_apiname"
)

const (
	CmdHelp      = "help"
	CmdMule      = "mule"
	CmdActReqMsg = "message"

	FuncMessageCreateSayHello       = "message-create/say-hello"
	FuncReactionAddReactionStats    = "reaction-add/reaction-stats"
	FuncReactionAddReactionRequired = "reaction-add/reaction-required"
)

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	lg = lg.With().Str(lkHandler, "Ready").Logger()

	lg.Info().Msgf("logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator)
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func emoji2msg(emojiAPIName string) string {
	if emojiAPIName == "" {
		return ""
	} else if !isASCII(emojiAPIName) {
		return emojiAPIName // normal emojis
	} else {
		return fmt.Sprintf("<:%s>", emojiAPIName) // e.g. <:ma:1151171171799269476>
	}
}

func getGuildEmojiAPINameByName(guildEmojis []*discordgo.Emoji, emojiName string) string {
	if emojiName == "" {
		return ""
	}
	if !isASCII(emojiName) {
		return emojiName // normal emojis
	}

	for _, e := range guildEmojis {
		if e.Name == emojiName {
			return e.APIName()
		}
	}

	return "" // GuildEmoji not found
}

// returns "" if member not found
func id2name(guildMembers []*discordgo.Member, id string) string {
	var member *discordgo.Member
	for _, m := range guildMembers {
		if m.User.ID == id {
			member = m
			break
		}
	}
	if member == nil {
		return ""
	}
	if member.Nick != "" {
		return member.Nick
	}
	if member.User.GlobalName != "" {
		return member.User.GlobalName
	}
	return member.User.Username
}

func sendSilentMessage(s *discordgo.Session, channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	data.Flags |= discordgo.MessageFlagsSuppressNotifications
	return s.ChannelMessageSendComplex(channelID, data, options...)
}

type choice struct {
	Weight uint8
	Data   string
}

func pickOne(choices []choice) string {
	sumWeight := 0
	for _, c := range choices {
		sumWeight += int(c.Weight)
	}
	if sumWeight == 0 {
		return ""
	}

	flattenChoices := make([]string, sumWeight)
	idx := 0
	for _, c := range choices {
		for i := 0; i < int(c.Weight); i++ {
			flattenChoices[idx] = c.Data
			idx++
		}
	}
	return flattenChoices[rand.Intn(sumWeight)]
}

package handlers

import (
	"fmt"
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

const MessageFlagsSilent = 1 << 12

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
func id2nick(guildMembers []*discordgo.Member, id string) string {
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
	if member.Nick == "" {
		return member.User.Username
	}
	return member.Nick
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

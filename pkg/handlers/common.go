package handlers

import (
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
)

const MessageFlagsSilent = 1 << 12

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	lg = lg.With().Str(lkHandler, "Ready").Logger()

	lg.Info().Msgf("logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator)
}

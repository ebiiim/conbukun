package handlers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

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
)

const (
	CmdHelp      = "help"
	CmdMule      = "mule"
	CmdActReq    = "action-required"
	CmdActReqMsg = "message"

	FuncReactionStats = "reaction-stats"
)

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	hn := "Ready"
	lg.Info().Str(lkHandler, hn).Msgf("logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator)
}

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	isDM := i.User != nil
	var usr *discordgo.User
	if isDM {
		usr = i.User
	} else {
		usr = i.Member.User
	}
	cmd := i.ApplicationCommandData().Name
	lg.Info().Str(lkIID, i.ID).Str(lkGuild, i.GuildID).Str(lkCh, i.ChannelID).Str(lkCmd, cmd).Bool(lkDM, isDM).Str(lkUsr, usr.ID).Str(lkName, usr.Username).Msg("OnInteractionCreate")
	if h, ok := CommandHandlers[cmd]; ok {
		h(s, i)
	}
}

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        CmdHelp,
			Description: "Show help message.",
		},
		{
			Name:        CmdMule,
			Description: "Show a random mule tips for 30 seconds.",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		CmdHelp: handleCmdHelp,
		CmdMule: handleCmdMule,
	}
)

func handleCmdHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var helpMsg = "" +
		"## コマンド\n" +
		"- `/help` このメッセージを表示します。\n" +
		"- `/mule` ラバに関するヒントをランダムに表示します。\n" +
		"## メンション\n" +
		"- **リアクション集計機能** 集計したいメッセージの返信に本botへのメンションとキーワード（`集計` `stats` `summary`）を入力すると表形式で出力します。\n" +
		"\n" +
		"> _[conbukun](https://github.com/ebiiim/conbukun) " + Version + " by ebiiim with ❤_" +
		""

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMsg,
			Flags:   discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsSuppressEmbeds | MessageFlagsSilent,
		},
	}); err != nil {
		lg.Error().Err(err).Str(lkIID, i.ID)
	}
}

var (
	muleMsgs = []string{
		"これはラバbotです（開発者より）",
		"【ラバ教豆知識】戦闘ラバの重さは110kg",
		"【ラバ教豆知識】ラバの重さは45kg",
		"あなたはラバを信じますか？ | Do you believe in Mule? | Ты веришь в мула? | 你相信骡子吗？",
		"ラバは世界を救います | Mule saves the world | мулы спасает мир | 骡子拯救世界",
		"ラバさえあればいい | No Mule, no life | все, что тебе нужно, это мул | 你只需要一头骡子",
		"ラバを讃えよ | Praise Mule | хвалите мула | 赞美骡子",
		"ラバは不滅です | Mule is immortal | мул бессмертен | 骡子是不朽的",
		"ラバ！ラバ！ラバ！ラバ！ラバ！ | Mule! Mule! Mule! Mule! Mule!",
		"[mule-n-img](https://render.albiononline.com/v1/item/Novice's%20Mule.png)",
		"[mule-b-img](https://render.albiononline.com/v1/item/Heretic%20Combat%20Mule.png)",
	}
)

const MessageFlagsSilent = 1 << 12

func handleCmdMule(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: muleMsgs[rand.Intn(len(muleMsgs))],
			Flags:   MessageFlagsSilent,
		},
	}); err != nil {
		lg.Error().Err(err).Str(lkIID, i.ID)
	}
	time.AfterFunc(time.Second*30, func() {
		if err := s.InteractionResponseDelete(i.Interaction); err != nil {
			lg.Error().Err(err).Str(lkIID, i.ID)
		}
	})
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not get channel")
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
		funcName = FuncReactionStats
	}

	lg.Debug().Str(lkMID, m.ID).Str(lkGuild, m.GuildID).Str(lkCh, m.ChannelID).Bool(lkDM, isDM).Str(lkUsr, m.Author.ID).Str(lkName, m.Author.Username).Msgf("isMention=%v isThread=%v hasRef=%v content=%s", isMention, isThread, hasRef, m.Content)

	if funcName == "" {
		return
	}

	lg.Info().Str(lkMID, m.ID).Str(lkGuild, m.GuildID).Str(lkCh, m.ChannelID).Str(lkFunc, funcName).Bool(lkDM, isDM).Str(lkUsr, m.Author.ID).Str(lkName, m.Author.Username).Msg("OnMessageCreate")
	switch funcName {
	case FuncReactionStats:
		handleReactionStats(s, m)
	default:
		return
	}
}

func containsReactionStats(s string) bool {
	ss := strings.ToLower(s)
	return strings.Contains(ss, "集計") || strings.Contains(ss, "sum") || strings.Contains(ss, "stats")
}

func handleReactionStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	lg.Info().Str(lkMID, m.ID).Str(lkFunc, FuncReactionStats).Msgf("ReactionStats: called")

	if err := s.ChannelTyping(m.ChannelID); err != nil {
		lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not send typing")
	}

	parentMsg, err := s.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
	if err != nil {
		lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not get parent message")
		return
	}

	userEmojis := map[string]map[string]bool{} // user -> emojis
	emojiUsers := map[string]map[string]bool{} // emoji -> users

	for _, u := range parentMsg.Mentions {
		userEmojis[u.Username] = map[string]bool{}
	}
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not get guildmembers")
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

	for _, r := range parentMsg.Reactions {
		users, err := s.MessageReactions(m.ChannelID, parentMsg.ID, r.Emoji.APIName(), 100, "", "")
		if err != nil {
			lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not get messagereactions")
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

	var emojiList []string
	for emoji := range emojiUsers {
		emojiList = append(emojiList, emoji)
	}

	var table strings.Builder
	table.WriteString("集計しました（2分間表示）\n")
	for _, emoji := range emojiList {
		if utf8.RuneCountInString(emoji) == 1 {
			table.WriteString(emoji)
		} else {
			table.WriteString(fmt.Sprintf("<:%s>", emoji)) // e.g. <:ma:1151171171799269476>
		}
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
		lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not send reply")
	}
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not delete summarization request message")
	}
	time.AfterFunc(time.Second*120, func() {
		lg.Info().Str(lkMID, m.ID).Str(lkFunc, FuncReactionStats).Msgf("ReactionStats: delete")
		if err := s.ChannelMessageDelete(m.ChannelID, reply.ID); err != nil {
			lg.Error().Err(err).Str(lkMID, m.ID).Msg("could not delete reply")
		}
	})
}
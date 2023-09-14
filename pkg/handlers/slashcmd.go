package handlers

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkHandler, "InteractionCreate").Str(lkIID, i.ID).Logger()

	isDM := i.User != nil
	var usr *discordgo.User
	if isDM {
		usr = i.User
	} else {
		usr = i.Member.User
	}
	cmd := i.ApplicationCommandData().Name
	lg.Info().Str(lkGuild, i.GuildID).Str(lkCh, i.ChannelID).Str(lkCmd, cmd).Bool(lkDM, isDM).Str(lkUsr, usr.ID).Str(lkName, usr.Username).Msg("OnInteractionCreate")
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

func getGuildEmojiAPINameByName(guildEmojis []*discordgo.Emoji, name string) string {
	for _, e := range guildEmojis {
		if e.Name == name {
			return e.APIName()
		}
	}
	return ""
}

func handleCmdHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdHelp).Str(lkIID, i.ID).Logger()

	guildEmojis, err := s.GuildEmojis(i.GuildID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get GuildEmojis")
	}

	geMa := emoji2msg(getGuildEmojiAPINameByName(guildEmojis, ReactionMa))
	geC1 := emoji2msg(getGuildEmojiAPINameByName(guildEmojis, ReactionConbu01))
	var helpMsg = "" +
		"## コマンド\n" +
		"- `/help` このメッセージを表示します。\n" +
		"- `/mule` ラバに関するヒントをランダムに表示します。\n" +
		"## リアクション\n" +
		"- [試験運用中] **リアクション集計機能** 集計したいメッセージにリアクション（🤖 " + geMa + " " + geC1 + " ）を行うと表形式で出力します。\n" +
		"## おまけ\n" +
		"- [試験運用中] 呼びかけに反応したりお昼寝したりします。\n" +
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
		lg.Error().Err(err).Msg("could not send InteractionResponse")
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

func handleCmdMule(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdMule).Str(lkIID, i.ID).Logger()

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: muleMsgs[rand.Intn(len(muleMsgs))],
			Flags:   MessageFlagsSilent,
		},
	}); err != nil {
		lg.Error().Err(err).Msg("could not send InteractionResponse")
	}
	time.AfterFunc(time.Second*30, func() {
		if err := s.InteractionResponseDelete(i.Interaction); err != nil {
			lg.Error().Err(err).Msg("could not delete InteractionResponse (AfterFunc)")
		}
	})
}

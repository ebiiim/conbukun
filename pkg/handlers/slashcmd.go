package handlers

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func NewOnInteractionCreateHandler(
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate),
) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	f := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		if h, ok := commandHandlers[cmd]; ok {
			h(s, i)
		}
	}

	return f
}

func InitializeApplicationCommands(cmds ...*discordgo.ApplicationCommand) []*discordgo.ApplicationCommand {
	return cmds
}

func InitializeCommandHandlers(handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return handlers
}

var (
	AppCmdHelp = &discordgo.ApplicationCommand{
		Name:        CmdHelp,
		Description: "こんぶくんについて知る",
	}

	AppCmdMule = &discordgo.ApplicationCommand{
		Name:        CmdMule,
		Description: "こんぶくんがラバ教の経典から引用してくれる（30秒後に自動削除）",
	}
)

func HandleCmdHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdHelp).Str(lkIID, i.ID).Logger()

	guildEmojis, err := s.GuildEmojis(i.GuildID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get GuildEmojis")
	}

	emojis2msg := func(guildEmojis []*discordgo.Emoji, emojis []string) string {
		var sb strings.Builder
		for _, emoji := range emojis {
			s := emoji2msg(getGuildEmojiAPINameByName(guildEmojis, emoji))
			if s == "" {
				continue
			}
			sb.WriteString(s)
			sb.WriteString(" ")
		}
		return strings.TrimRight(sb.String(), " ")
	}
	var helpMsg = "" +
		"## コマンド\n" +
		"- `/help` このメッセージを表示します。\n" +
		"- `/mule` ラバに関するヒントをランダムに投稿します（30秒後に自動削除）。\n" +
		"- `/route-add` アバロンのルートを追加します。\n" +
		"- `/route-print` アバロンのルートを画像で投稿します。\n" +
		"- `/route-clear` アバロンのルートをリセットします。\n" +
		"- `/route-mark` マップをマークします（色が変わります）。\n" +
		"## リアクション\n" +
		"- **リアクション集計** 集計したいメッセージにリアクション（" + emojis2msg(guildEmojis, EmojisReactionAddReactionRequired) + "）を行うとリマインダーを投稿します（2分後に自動削除）。\n" +
		// "- [試験運用中] **リアクション集計（表）** 集計したいメッセージにリアクション（" + emojis2msg(guildEmojis, emojisReactionAddReactionStats) + "）を行うと表形式で投稿します（2分後に削除）。\n" +
		"## おまけ\n" +
		"- 呼びかけに反応したりお昼寝したりします。\n" +
		"\n" +
		"> _[conbukun](https://github.com/ebiiim/conbukun) " + Version + " by ebiiim with ❤_" +
		""

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMsg,
			Flags:   discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsSuppressEmbeds | discordgo.MessageFlagsSuppressNotifications,
		},
	}); err != nil {
		lg.Error().Err(err).Msg("could not send InteractionResponse")
	}
}

var (
	muleMsgs = []choice{
		{10, "わん"},
		{10, "ファッキンラバ"},
		{4, "「これはラバbotです」（開発者より）"},
		{8, "「アルコールポーションうめぇ」（開発者より）"},
		{10, "【ラバ教豆知識】戦闘ラバの重さは110kg"},
		{10, "【ラバ教豆知識】ラバの重さは45kg"},
		// too messy not cool
		// {6, "> あなたはラバを信じますか？ | Do you believe in Mule? | Ты веришь в мула? | 你相信骡子吗？"},
		// {6, "> ラバは世界を救います | Mule saves the world | мулы спасает мир | 骡子拯救世界"},
		// {6, "> ラバさえあればいい | No Mule, no life | все, что тебе нужно, это мул | 你只需要一头骡子"},
		// {6, "> ラバを讃えよ | Praise Mule | хвалите мула | 赞美骡子"},
		// {6, "> ラバは不滅です | Mule is immortal | мул бессмертен | 骡子是不朽的"},
		// {6, "> ラバ！ラバ！ラバ！ラバ！ラバ！ | Mule! Mule! Mule! Mule! Mule!"},
		{6, "「あなたはラバを信じますか？」"},
		{6, "「ラバは世界を救います」"},
		{6, "「ラバさえあればいい」"},
		{6, "「ラバを讃えよ」"},
		{6, "「ラバは不滅です」"},
		{6, "「ラバ！ラバ！ラバ！ラバ！ラバ！」"},
		{6, "「シェイプシフターなぜラバに変身しない？」"},
		{6, "「ハッピーハロウィン！ラバ万歳！」"},
		{6, "[mule-t2-img](https://render.albiononline.com/v1/item/Novice's%20Mule.png)"},
		{6, "[mule-t6-img](https://render.albiononline.com/v1/item/Heretic%20Combat%20Mule.png)"},
		// too large not displayed
		// {6, "[mount-jod-img](https://wiki.albiononline.com/data/images/4/4e/JackODonkeyMountSkin.png)"},
		// {6, "[mount-mule-t2-img](https://wiki.albiononline.com/data/images/e/e7/NovicesMule.png)"},
	}
)

func HandleCmdMule(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdMule).Str(lkIID, i.ID).Logger()

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: pickOne(muleMsgs),
			Flags:   discordgo.MessageFlagsSuppressNotifications,
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

func respondEphemeralMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsSuppressEmbeds | discordgo.MessageFlagsSuppressNotifications,
		},
	})
}

func respondEphemeralMessageEdit(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) error {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	return err
}

func respondEphemeralMessageDelete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionResponseDelete(i.Interaction)
}

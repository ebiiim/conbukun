package handlers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/ebiiim/conbukun/pkg/ao/data"
	"github.com/ebiiim/conbukun/pkg/ao/roanav"
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
	AppCmdHelp = &discordgo.ApplicationCommand{
		Name:        CmdHelp,
		Description: "こんぶくんについて知る",
	}

	AppCmdMule = &discordgo.ApplicationCommand{
		Name:        CmdMule,
		Description: "こんぶくんがラバ教の経典から引用してくれる（30秒後に自動削除）",
	}

	AppCmdRouteAdd = &discordgo.ApplicationCommand{
		Name:        CmdRouteAdd,
		Description: "アバロンのルートを追加する",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "from",
				Description:  "出発地（何文字か入力して！）",
				Type:         discordgo.ApplicationCommandOptionString,
				Autocomplete: true,
				Required:     true,
			},
			{
				Name:         "to",
				Description:  "目的地（何文字か入力して！）",
				Type:         discordgo.ApplicationCommandOptionString,
				Autocomplete: true,
				Required:     true,
			},
			{
				Name:        "color",
				Description: "ポータルの色（青or黄 入れる人数が決まる）",
				Type:        discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Yellow",
						Value: roanav.PortalTypeYellow,
					},
					{
						Name:  "Blue",
						Value: roanav.PortalTypeBlue,
					},
				},
				Required: true,
			},
			{
				Name:        "time",
				Description: "残り時間（HHmm形式 3時間14分なら0314）",
				Type:        discordgo.ApplicationCommandOptionInteger,
				MinValue:    Ptr(0.0),
				MaxValue:    2359.0,
				Required:    true,
			},
		},
	}

	AppCmdRoutePrint = &discordgo.ApplicationCommand{
		Name:        CmdRoutePrint,
		Description: "アバロンのルートを表示する",
	}

	Commands = []*discordgo.ApplicationCommand{
		AppCmdHelp,
		AppCmdMule,
		AppCmdRouteAdd,
		AppCmdRoutePrint,
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		CmdHelp:       handleCmdHelp,
		CmdMule:       handleCmdMule,
		CmdRouteAdd:   handleCmdRouteAdd,
		CmdRoutePrint: handleCmdRoutePrint,
	}
)

func handleCmdHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		"- `/route-print` アバロンのルートを表示します。\n" +
		"## リアクション\n" +
		"- **リアクション集計** 集計したいメッセージにリアクション（" + emojis2msg(guildEmojis, emojisReactionAddReactionRequired) + "）を行うとリマインダーを投稿します（2分後に自動削除）。\n" +
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

func handleCmdMule(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

var (
	// TODO: make this private, persistent and thread-safe
	navigations                        = map[string]*roanav.Navigation{}
	mapNameCompleter *MapNameCompleter = nil
)

func init() {
	navigations = map[string]*roanav.Navigation{}
	mapNameCompleter = NewMapNameCompleter(5)
}

func navigationName(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	c, err := s.State.Channel(i.ChannelID)
	if err != nil {
		return "", err
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s#%s (conbukun@%s)", g.Name, c.Name, Version), nil
}

func handleCmdRouteAdd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteAdd+"/dispatcher").Str(lkIID, i.ID).Logger()
	switch i.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		handleCmdRouteAddAutocomplete(s, i)
	case discordgo.InteractionApplicationCommand:
		handleCmdRouteAddCommand(s, i)
	default:
		lg.Error().Msgf("unknown InteractionType: %s", i.Type.String())
	}
}

func handleCmdRouteAddAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteAdd+"/autocomplete").Str(lkIID, i.ID).Logger()

	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch {
	case data.Options[0].Focused: // from
		choices = mapNameCompleter.GetChoices(data.Options[0].StringValue())
	case data.Options[1].Focused: // to
		choices = mapNameCompleter.GetChoices(data.Options[1].StringValue())
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}); err != nil {
		lg.Error().Err(err).Msg("could not send InteractionResponse")
	}
}

func handleCmdRouteAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteAdd).Str(lkIID, i.ID).Logger()
	if i.Member == nil {
		if mErr := respondWithEphemeralMessage(s, i, "エラー: この機能はDMではなくサーバーで使ってほしいわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Fetch data.
	members, err := s.GuildMembers(i.GuildID, "", 1000)
	if err != nil {
		lg.Error().Err(err).Msg("could not get GuildMembers")
		return
	}
	userName := id2name(members, i.Member.User.ID)
	if userName == "" {
		userName = i.Member.User.Username
	}

	// Get the Navigation.
	navname, err := navigationName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name")
		if mErr := respondWithEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかチャンネルの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if _, ok := navigations[navname]; !ok {
		navigations[navname] = &roanav.Navigation{
			Name:    navname,
			Portals: []*roanav.Portal{},
		}
	}
	nav := navigations[navname]
	nav.DeleteExpiredPortals()

	// Get arguments.
	optFrom := i.ApplicationCommandData().Options[0]
	optTo := i.ApplicationCommandData().Options[1]
	optColor := i.ApplicationCommandData().Options[2]
	optTime := i.ApplicationCommandData().Options[3]
	from := optFrom.StringValue()
	to := optTo.StringValue()
	color := optColor.StringValue()
	timeVal := optTime.IntValue()
	timeMinute := timeVal % 100
	timeHour := (timeVal - timeMinute) / 100
	lg.Info().Str("from", from).Str("to", to).Str("color", color).Int("time", int(timeVal)).Msg("arguments")

	// Validate arguments.
	if from == to {
		lg.Error().Err(fmt.Errorf("from and to are the same")).Msg("invalid arguments")
		if mErr := respondWithEphemeralMessage(s, i, "エラー: `from` と `to` は異なるマップにしてほしいわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if _, ok := data.Maps[from]; !ok {
		lg.Error().Err(fmt.Errorf("invalid from")).Msg("invalid arguments")
		if mErr := respondWithEphemeralMessage(s, i, "エラー: `from` に知らないマップ名が入ってるわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if _, ok := data.Maps[to]; !ok {
		lg.Error().Err(fmt.Errorf("invalid to")).Msg("invalid arguments")
		if mErr := respondWithEphemeralMessage(s, i, "エラー: `to` に知らないマップ名が入ってるわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if timeHour < 0 || timeHour > 23 || timeMinute < 0 || timeMinute > 59 {
		lg.Error().Err(fmt.Errorf("invalid time")).Msg("invalid arguments")
		if mErr := respondWithEphemeralMessage(s, i, "エラー: `time` は `HHmm` のフォーマットで入力してほしいわん（3時間14分なら `0314` ）"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	portal := roanav.NewPortal(
		from, to,
		color,
		time.Now().Add(time.Hour*time.Duration(timeHour)+time.Minute*time.Duration(timeMinute)),
		map[string]string{
			roanav.PortalDataKeyUser: userName,
		},
	)
	nav.AddPortal(portal)

	if mErr := respondWithEphemeralMessage(s, i,
		fmt.Sprintf("追加したわん！いまこんな感じ！\n%s`/route-print` で画像を投稿できるわん！", roanav.BriefNavigation(nav, data.Maps)),
	); mErr != nil {
		lg.Error().Err(mErr).Msg("could not send InteractionResponse")
	}
}

func respondWithEphemeralMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsSuppressEmbeds | discordgo.MessageFlagsSuppressNotifications,
		},
	})
}

func handleCmdRoutePrint(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRoutePrint).Str(lkIID, i.ID).Logger()
	if i.Member == nil {
		if mErr := respondWithEphemeralMessage(s, i, "エラー: この機能はDMではなくサーバーで使ってほしいわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Get the Navigation.
	navname, err := navigationName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name")
		if mErr := respondWithEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかチャンネルの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if _, ok := navigations[navname]; !ok {
		navigations[navname] = &roanav.Navigation{
			Name:    navname,
			Portals: []*roanav.Portal{},
		}
	}
	nav := navigations[navname]
	nav.DeleteExpiredPortals()
	if nav.Portals == nil || len(nav.Portals) == 0 {
		lg.Error().Err(fmt.Errorf("no portals")).Msg("len(nav.Portals) == 0")
		if mErr := respondWithEphemeralMessage(s, i, "エラー: 有効なルートが1個もないわん。 `/route-add` で追加してからまた試してほしいわん。"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Generate PlantUML.
	p := roanav.NewKrokiPlantUMLPainter(roanav.DefaultKrokiEndpoint, roanav.DefaultKrokiTimeout, data.Maps)
	dist, err := p.Paint(nav)
	if err != nil {
		lg.Error().Err(err).Msg("could not generate PlantUML")
		if mErr := respondWithEphemeralMessage(s, i, fmt.Sprintf("エラー: 画像の生成に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Send PNG.
	pngFile, err := os.Open(dist)
	if err != nil {
		lg.Error().Err(err).Msg("could not open PNG file")
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s お待たせしましたわん！", i.Member.User.Mention()),
			Flags:   discordgo.MessageFlagsSuppressNotifications,
			Files: []*discordgo.File{
				{
					Name:        dist, // unsafe chars will be stripped
					ContentType: "image/png",
					Reader:      pngFile,
				},
			},
		},
	}); err != nil {
		lg.Error().Err(err).Msg("could not send InteractionResponse")
	}
}

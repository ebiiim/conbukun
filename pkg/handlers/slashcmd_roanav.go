package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/ebiiim/conbukun/pkg/ao/data"
	"github.com/ebiiim/conbukun/pkg/ao/roanav"
)

var (
	ROANavHandlerSuggestionsLimit = 5
)

type ROANavHandler struct {
	navigations      sync.Map
	MapNameCompleter *MapNameCompleter

	saveFile string
}

func NewROANavHandler(saveFile string) (*ROANavHandler, error) {
	lg := lg.With().Str("func", "NewROANavHandler").Logger()

	lg.Info().
		Str("saveFile", saveFile).
		Int("ROANavHandlerSuggestionsLimit", ROANavHandlerSuggestionsLimit).
		Msg("initializing ROANavHandler")

	rn := &ROANavHandler{
		navigations:      sync.Map{},
		MapNameCompleter: NewMapNameCompleter(ROANavHandlerSuggestionsLimit),
		saveFile:         saveFile,
	}

	if saveFile == "" {
		lg.Warn().Msg("saveFile is empty, ROANav will not be persistent")
		return rn, nil
	}
	if err := rn.Load(); err != nil {
		return nil, err
	}

	return rn, nil
}

func (h *ROANavHandler) Save() error {
	if h.saveFile == "" {
		return nil
	}

	// create save file if not exists
	if _, err := os.Stat(h.saveFile); os.IsNotExist(err) {
		f, err := os.Create(h.saveFile)
		if err != nil {
			return err
		}
		if _, err := f.Write([]byte("{}")); err != nil {
			return err
		}
		f.Close()
		lg.Info().Str("saveFile", h.saveFile).Msg("ROANavHandler created save file")
	}

	// open save file
	f, err := os.Open(h.saveFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// save
	if err := h.ExportNavigations(f); err != nil {
		return err
	}

	lg.Info().Str("saveFile", h.saveFile).Msg("ROANavHandler saved to file")
	return nil
}

func (h *ROANavHandler) Load() error {
	if h.saveFile == "" {
		return nil
	}

	// create save file if not exists
	if _, err := os.Stat(h.saveFile); os.IsNotExist(err) {
		f, err := os.Create(h.saveFile)
		if err != nil {
			return err
		}
		if _, err := f.Write([]byte("{}")); err != nil {
			return err
		}
		f.Close()
		lg.Info().Str("saveFile", h.saveFile).Msg("ROANavHandler created save file")
	}

	// open save file
	f, err := os.Open(h.saveFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// load
	if err := h.ImportNavigations(f); err != nil {
		return err
	}

	lg.Info().Str("saveFile", h.saveFile).Msg("ROANavHandler loaded from file")
	return nil
}

func (h *ROANavHandler) GetOrCreateNavigation(name string) *roanav.Navigation {
	if v, ok := h.navigations.Load(name); ok {
		return v.(*roanav.Navigation)
	}
	nav := roanav.NewNavigation(name)
	h.navigations.Store(name, nav)
	return nav
}

func (h *ROANavHandler) DeleteNavigation(name string) { h.navigations.Delete(name) }

func (h *ROANavHandler) ExportNavigations(w io.Writer) error {
	jm := make(map[string]*roanav.Navigation)
	h.navigations.Range(func(k, v interface{}) bool {
		jm[k.(string)] = v.(*roanav.Navigation)
		return true
	})
	return json.NewEncoder(w).Encode(jm)
}

func (h *ROANavHandler) ImportNavigations(r io.Reader) error {
	jm := make(map[string]*roanav.Navigation)
	if err := json.NewDecoder(r).Decode(&jm); err != nil {
		return err
	}
	for k, v := range jm {
		h.navigations.Store(k, v)
	}
	return nil
}

func (h *ROANavHandler) HandleCmdRouteAdd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteAdd+"/dispatcher").Str(lkIID, i.ID).Logger()
	switch i.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		h.HandleCmdRouteAddAutocomplete(s, i)
	case discordgo.InteractionApplicationCommand:
		h.HandleCmdRouteAddCommand(s, i)
	default:
		lg.Error().Msgf("unknown InteractionType: %s", i.Type.String())
	}
}

func (h *ROANavHandler) HandleCmdRouteAddAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteAdd+"/autocomplete").Str(lkIID, i.ID).Logger()

	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch {
	case data.Options[0].Focused: // from
		choices = h.MapNameCompleter.GetChoices(data.Options[0].StringValue())
	case data.Options[1].Focused: // to
		choices = h.MapNameCompleter.GetChoices(data.Options[1].StringValue())
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

func getNavNameAndUserName(s *discordgo.Session, i *discordgo.InteractionCreate) (navName, userName string, err error) {

	isDM := (i.Member == nil)

	if isDM {
		userName = i.User.Username
		navName = fmt.Sprintf("DM#%s (conbukun@%s)", i.User.Username, Version)
	} else {
		members, iErr := s.GuildMembers(i.GuildID, "", 1000)
		if iErr != nil {
			return "", "", iErr
		}

		userName = id2name(members, i.Member.User.ID)
		if userName == "" {
			userName = i.Member.User.Username
		}

		c, err := s.State.Channel(i.ChannelID)
		if err != nil {
			return "", "", err
		}
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return "", "", err
		}
		navName = fmt.Sprintf("%s#%s (conbukun@%s)", g.Name, c.Name, Version)
	}

	return
}

func (h *ROANavHandler) HandleCmdRouteAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteAdd).Str(lkIID, i.ID).Logger()

	// Get names.
	navName, userName, err := getNavNameAndUserName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name or user name")
		if mErr := respondWithEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Get the Navigation.
	nav := h.GetOrCreateNavigation(navName)
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

	// Add portal and save.
	portal := roanav.NewPortal(
		from, to,
		color,
		time.Now().Add(time.Hour*time.Duration(timeHour)+time.Minute*time.Duration(timeMinute)),
		map[string]string{
			roanav.PortalDataKeyUser: userName,
		},
	)
	nav.AddPortal(portal)
	if err := h.Save(); err != nil {
		lg.Error().Err(err).Msg("could not save navigations")
	}

	if mErr := respondWithEphemeralMessage(s, i,
		fmt.Sprintf("追加したわん！いまこんな感じ！\n%s`/route-print` で画像を投稿できるわん！", roanav.BriefNavigation(nav, data.Maps)),
	); mErr != nil {
		lg.Error().Err(mErr).Msg("could not send InteractionResponse")
	}
}

func (h *ROANavHandler) HandleCmdRoutePrint(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRoutePrint).Str(lkIID, i.ID).Logger()

	// Get names.
	navName, _, err := getNavNameAndUserName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name or user name")
		if mErr := respondWithEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Get the Navigation.
	nav := h.GetOrCreateNavigation(navName)
	nav.DeleteExpiredPortals()

	// Validate.
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
			Content: "お待たせしましたわん！",
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

func (h *ROANavHandler) HandleCmdRouteClear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRoutePrint).Str(lkIID, i.ID).Logger()

	// Get names.
	navName, _, err := getNavNameAndUserName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name or user name")
		if mErr := respondWithEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Get the Navigation.
	nav := h.GetOrCreateNavigation(navName)
	nav.DeleteExpiredPortals()

	// Validate.
	if nav.Portals == nil || len(nav.Portals) == 0 {
		if mErr := respondWithEphemeralMessage(s, i, "ルートをクリアしたわん！"); mErr != nil {
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
			Content: "ルートをクリアしたわん！念のためクリアする前の状態を投稿しておくわん！",
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

	// Clear.
	h.DeleteNavigation(navName)
}

func (*ROANavHandler) CommandRouteAdd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CmdRouteAdd,
		Description: "こんぶくんにアバロンのルートを覚えてもらう",
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
				Description: "残り時間（HHmm形式 3時間14分なら0314、0000でルートを削除）",
				Type:        discordgo.ApplicationCommandOptionInteger,
				MinValue:    Ptr(0.0),
				MaxValue:    2359.0,
				Required:    true,
			},
		},
	}
}

func (*ROANavHandler) CommandRoutePrint() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CmdRoutePrint,
		Description: "こんぶくんに覚えているアバロンのルートを画像で投稿してもらう",
	}
}

func (*ROANavHandler) CommandRouteClear() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CmdRouteClear,
		Description: "こんぶくんに覚えているアバロンのルートを全部忘れてもらう",
	}
}

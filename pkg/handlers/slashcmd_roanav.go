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
	DefaultROANavHandlerSuggestionsLimit = 5
)

type ROANavHandler struct {
	navigations      sync.Map
	MapNameCompleter *MapNameCompleter

	saveFile      string
	saveFileMutex sync.Mutex
}

func NewROANavHandler(saveFile string, suggestionsLimit int) (*ROANavHandler, error) {
	lg := lg.With().Str("func", "NewROANavHandler").Logger()

	lg.Info().
		Str("saveFile", saveFile).
		Int("suggestionsLimit", suggestionsLimit).
		Msg("initializing ROANavHandler")

	rn := &ROANavHandler{
		navigations:      sync.Map{},
		MapNameCompleter: NewMapNameCompleter(suggestionsLimit),
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
	h.saveFileMutex.Lock()
	defer h.saveFileMutex.Unlock()

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
	f, err := os.OpenFile(h.saveFile, os.O_WRONLY, 0644)
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
	h.saveFileMutex.Lock()
	defer h.saveFileMutex.Unlock()

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

// GetOrCreateNavigation always returns a non-nil value.
func (h *ROANavHandler) GetOrCreateNavigation(name string) *roanav.Navigation {
	if v, ok := h.navigations.Load(name); ok {
		return v.(*roanav.Navigation)
	}
	nav := roanav.NewNavigation(name)
	nav.Data[roanav.NavigationDataMarkedMaps] = "[]" // init
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

func getUser(i *discordgo.InteractionCreate) *discordgo.User {
	if i.Member == nil {
		return i.User
	}
	return i.Member.User
}

func getNavNameAndUserName(s *discordgo.Session, i *discordgo.InteractionCreate) (navName, userName string, err error) {

	isDM := (i.Member == nil)

	if isDM {
		userName = i.User.Username
		navName = fmt.Sprintf("@%s", i.User.Username)
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
		navName = fmt.Sprintf("%s#%s", g.Name, c.Name)
	}

	return
}

func (h *ROANavHandler) HandleCmdRouteList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteList).Str(lkIID, i.ID).Logger()

	// Get names.
	navName, _, err := getNavNameAndUserName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name or user name")
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	lg = lg.With().Str("navName", navName).Logger()

	// Get the Navigation.
	nav := h.GetOrCreateNavigation(navName)
	nav.DeleteExpiredPortals()

	// Get MarkedMaps.
	var markedMaps []roanav.MarkedMap
	if err := json.Unmarshal([]byte(nav.Data[roanav.NavigationDataMarkedMaps]), &markedMaps); err != nil {
		lg.Error().Err(err).Msg("could not unmarshal marked maps")
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: マークされているマップの取得に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Construct response.
	var resp string
	resp += "いまこんな感じだわん！ `/route-add` でルートを追加、 `/route-mark` でコメント、 `/route-print` で画像を投稿するわん！\n"
	resp += "\n"
	resp += "***有効なルート***\n"
	srt := roanav.BriefNavigation(nav, data.Maps)
	if srt == "" {
		resp += "（なし）\n"
	} else {
		resp += srt
	}
	resp += "\n"
	resp += "***マークされているマップ***\n"
	smm := briefMarkedMaps(markedMaps)
	if smm == "" {
		resp += "（なし）\n"
	} else {
		resp += smm
	}

	// Send response.
	if mErr := respondEphemeralMessage(s, i, resp); mErr != nil {
		lg.Error().Err(mErr).Msg("could not send InteractionResponse")
	}
}

func (h *ROANavHandler) HandleCmdRouteAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteAdd).Str(lkIID, i.ID).Logger()

	// Get names.
	navName, userName, err := getNavNameAndUserName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name or user name")
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	lg = lg.With().Str("navName", navName).Logger()

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
		if mErr := respondEphemeralMessage(s, i, "エラー: `from` と `to` は異なるマップにしてほしいわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if _, ok := data.Maps[from]; !ok {
		lg.Error().Err(fmt.Errorf("invalid from")).Msg("invalid arguments")
		if mErr := respondEphemeralMessage(s, i, "エラー: `from` に知らないマップ名が入ってるわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if _, ok := data.Maps[to]; !ok {
		lg.Error().Err(fmt.Errorf("invalid to")).Msg("invalid arguments")
		if mErr := respondEphemeralMessage(s, i, "エラー: `to` に知らないマップ名が入ってるわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	if timeHour < 0 || timeHour > 23 || timeMinute < 0 || timeMinute > 59 {
		lg.Error().Err(fmt.Errorf("invalid time")).Msg("invalid arguments")
		if mErr := respondEphemeralMessage(s, i, "エラー: `time` は `HHmm` のフォーマットで入力してほしいわん（3時間14分なら `0314` ）"); mErr != nil {
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

	if mErr := respondEphemeralMessage(s, i,
		"ありがとうわん！ `/route-list` で確認、 `/route-print` で画像を投稿してみんなと共有するわん！",
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
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	lg = lg.With().Str("navName", navName).Logger()

	// Get the Navigation.
	nav := h.GetOrCreateNavigation(navName)
	nav.DeleteExpiredPortals()

	// Validate.
	if nav.Portals == nil || len(nav.Portals) == 0 {
		lg.Error().Err(fmt.Errorf("no portals")).Msg("len(nav.Portals) == 0")
		if mErr := respondEphemeralMessage(s, i, "エラー: 有効なルートがないわん。 `/route-list` で確認したり `/route-add` で追加したりしてわん！"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// First reaction.
	if mErr := respondEphemeralMessage(s, i, "処理中だわん"); mErr != nil {
		lg.Error().Err(mErr).Msg("could not send InteractionResponse")
	}

	// Generate PlantUML.
	p := roanav.NewKrokiPlantUMLPNGPainter(roanav.DefaultKrokiEndpoint, roanav.DefaultKrokiTimeout, data.Maps, roanav.PlantUMLStyleAuto)
	dist, err := p.Paint(nav)
	if err != nil {
		lg.Error().Err(err).Msg("could not generate PlantUML")
		if mErr := respondEphemeralMessageEdit(s, i, fmt.Sprintf("エラー: 画像の生成に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Send PNG.
	pngFile, err := os.Open(dist)
	if err != nil {
		lg.Error().Err(err).Msg("could not open PNG file")
	}
	if _, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("%s お待たせしましたわん！", getUser(i).Mention()),
		Flags:   discordgo.MessageFlagsSuppressNotifications,
		Files: []*discordgo.File{
			{
				Name:        dist, // unsafe chars will be stripped
				ContentType: "image/png",
				Reader:      pngFile,
			},
		},
	}); err != nil {
		lg.Error().Err(err).Msg("could not send message")
		if mErr := respondEphemeralMessageEdit(s, i, fmt.Sprintf("エラー: 画像の投稿に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	respondEphemeralMessageDelete(s, i)
}

func (h *ROANavHandler) HandleCmdRouteClear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRoutePrint).Str(lkIID, i.ID).Logger()

	// Get names.
	navName, _, err := getNavNameAndUserName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name or user name")
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Get the Navigation.
	nav := h.GetOrCreateNavigation(navName)
	nav.DeleteExpiredPortals()

	// Validate.
	if nav.Portals == nil || len(nav.Portals) == 0 {
		if mErr := respondEphemeralMessage(s, i, "クリアするルートがないわん"); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// First reaction.
	if mErr := respondEphemeralMessage(s, i, "処理中だわん"); mErr != nil {
		lg.Error().Err(mErr).Msg("could not send InteractionResponse")
	}

	// Generate PlantUML.
	p := roanav.NewKrokiPlantUMLPNGPainter(roanav.DefaultKrokiEndpoint, roanav.DefaultKrokiTimeout, data.Maps, roanav.PlantUMLStyleAuto)
	dist, err := p.Paint(nav)
	if err != nil {
		lg.Error().Err(err).Msg("could not generate PlantUML")
		if mErr := respondEphemeralMessageEdit(s, i, fmt.Sprintf("エラー: 画像の生成に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Send PNG.
	pngFile, err := os.Open(dist)
	if err != nil {
		lg.Error().Err(err).Msg("could not open PNG file")
	}
	if _, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("%s ルートをクリアしたわん！念のため最後の状態を投稿しておくわん！（バックアップもあるので、間違えた場合は管理者に知らせるわん）", getUser(i).Mention()),
		Flags:   discordgo.MessageFlagsSuppressNotifications,
		Files: []*discordgo.File{
			{
				Name:        dist, // unsafe chars will be stripped
				ContentType: "image/png",
				Reader:      pngFile,
			},
		},
	}); err != nil {
		lg.Error().Err(err).Msg("could not send InteractionResponse")
		if mErr := respondEphemeralMessageEdit(s, i, fmt.Sprintf("エラー: 画像の投稿に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	// Backup before clear.
	navbakName := fmt.Sprintf("%s___bak_%d", navName, time.Now().Unix()) // hoge#fuga___bak_1234567890
	navbak := h.GetOrCreateNavigation(navbakName)
	nav.DeepCopyInto(navbak)
	navbak.Name = navbakName

	// Clear.
	nav.Portals = []*roanav.Portal{} // clear portals instead of deleting navigation to keep marked maps

	// Save.
	if err := h.Save(); err != nil {
		lg.Error().Err(err).Msg("could not save navigations")
	}
}

func (h *ROANavHandler) HandleCmdRouteMark(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteMark+"/dispatcher").Str(lkIID, i.ID).Logger()
	switch i.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		h.HandleCmdRouteMarkAutocomplete(s, i)
	case discordgo.InteractionApplicationCommand:
		h.HandleCmdRouteMarkCommand(s, i)
	default:
		lg.Error().Msgf("unknown InteractionType: %s", i.Type.String())
	}
}

func (h *ROANavHandler) HandleCmdRouteMarkAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteMark+"/autocomplete").Str(lkIID, i.ID).Logger()

	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch {
	case data.Options[0].Focused: // map
		choices = h.MapNameCompleter.GetChoices(data.Options[0].StringValue())
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

// upsertOrRemoveMarkedMap upserts or removes the marked map.
// If markedMap.Color is none and markedMap.Comment is empty, the element with same markedMap.ID will be removed.
func upsertOrRemoveMarkedMap(markedMaps []roanav.MarkedMap, markedMap roanav.MarkedMap) []roanav.MarkedMap {
	if markedMap.ID == "" {
		return markedMaps
	}
	var m2 []roanav.MarkedMap
	isNew := true
	for _, m := range markedMaps {
		if m.ID != markedMap.ID {
			m2 = append(m2, m)
		} else {
			isNew = false
			if markedMap.Color == roanav.MarkedMapColorNone && markedMap.Comment == "" {
				// remove
			} else {
				// update
				m2 = append(m2, markedMap)
			}
		}
	}
	if isNew && !(markedMap.Color == roanav.MarkedMapColorNone && markedMap.Comment == "") {
		m2 = append(m2, markedMap)
	}
	return m2
}

// briefMarkedMaps returns a string representation of marked maps for printing.
func briefMarkedMaps(markedMaps []roanav.MarkedMap) string {
	markedMapsStr := ""
	for _, m := range markedMaps {
		md, ok := data.Maps[m.ID]
		if !ok {
			continue
		}
		markedMapsStr += fmt.Sprintf("- %s", md.DisplayName)
		if m.Color != roanav.MarkedMapColorNone {
			markedMapsStr += fmt.Sprintf(" (%s)", m.Color)
		}
		if m.Comment != "" {
			markedMapsStr += fmt.Sprintf(" \"%s\"", m.Comment)
		}
		if m.User != "" {
			markedMapsStr += fmt.Sprintf(" by %s", m.User)
		} else {
			markedMapsStr += " by ユーザ情報なし" // for backward compatibility
		}
		markedMapsStr += "\n"
	}
	return markedMapsStr
}

func (h *ROANavHandler) HandleCmdRouteMarkCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lg := lg.With().Str(lkCmd, CmdRouteMark).Str(lkIID, i.ID).Logger()

	// Get names.
	navName, userName, err := getNavNameAndUserName(s, i)
	if err != nil {
		lg.Error().Err(err).Msg("could not get navigation name or user name")
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: サーバーかユーザーの名前が取得できなかったわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	lg = lg.With().Str("navName", navName).Logger()

	// Get the Navigation.
	nav := h.GetOrCreateNavigation(navName)
	nav.DeleteExpiredPortals()

	// Get arguments.
	optMap := i.ApplicationCommandData().Options[0]
	targetMap := optMap.StringValue()

	optColor := i.ApplicationCommandData().Options[1]
	targetColor := optColor.StringValue()

	targetComment := ""
	if len(i.ApplicationCommandData().Options) >= 3 {
		optComment := i.ApplicationCommandData().Options[2]
		targetComment = optComment.StringValue()
	}

	lg.Info().Str("map", targetMap).Str("color", targetColor).Str("comment", targetComment).Msg("arguments")

	// Validate arguments.
	// nothing to validate

	// Construct MarkedMap.
	markedMap := roanav.MarkedMap{
		ID:      targetMap,
		Color:   targetColor,
		Comment: targetComment,
		User:    userName,
	}

	// Init MarkedMaps if not exists.
	if _, ok := nav.Data[roanav.NavigationDataMarkedMaps]; !ok {
		nav.Data[roanav.NavigationDataMarkedMaps] = "[]"
	}

	// Get current marked maps.
	var markedMaps []roanav.MarkedMap
	if err := json.Unmarshal([]byte(nav.Data[roanav.NavigationDataMarkedMaps]), &markedMaps); err != nil {
		lg.Error().Err(err).Msg("could not unmarshal marked maps")
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: マークされているマップの取得に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}

	markedMaps = upsertOrRemoveMarkedMap(markedMaps, markedMap)

	// Set new marked maps.
	markedMapsJSON, err := json.Marshal(markedMaps)
	if err != nil {
		lg.Error().Err(err).Msg("could not marshal marked maps")
		if mErr := respondEphemeralMessage(s, i, fmt.Sprintf("エラー: マークされているマップの保存に失敗したわん。何回も発生する場合は管理者に知らせてほしいわん。 ```\n%v```", err)); mErr != nil {
			lg.Error().Err(mErr).Msg("could not send InteractionResponse")
		}
		return
	}
	nav.Data[roanav.NavigationDataMarkedMaps] = string(markedMapsJSON)

	// Save.
	if err := h.Save(); err != nil {
		lg.Error().Err(err).Msg("could not save navigations")
	}

	if mErr := respondEphemeralMessage(s, i,
		"わかったわん！ `/route-list` でいまの状態を確認できるわん！\n",
	); mErr != nil {
		lg.Error().Err(mErr).Msg("could not send InteractionResponse")
	}
}

func (*ROANavHandler) CommandRouteList() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CmdRouteList,
		Description: "こんぶくんに知ってる情報をこっそり教えてもらう（確認用）",
	}
}

func (*ROANavHandler) CommandRouteAdd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CmdRouteAdd,
		Description: "こんぶくんにアバロンのルートを覚えてもらう（チャンネルごとに分かれてるよ！）",
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
		Description: "こんぶくんに覚えているアバロンのルートを画像でチャンネルに投稿してもらう（みんなに共有！）",
	}
}

func (*ROANavHandler) CommandRouteClear() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CmdRouteClear,
		Description: "[取扱注意] こんぶくんに覚えているアバロンのルートを全部忘れてもらう（リセット機能）",
	}
}

func (*ROANavHandler) CommandRouteMark() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CmdRouteMark,
		Description: "こんぶくんに特別なマップに色をつけてもらう",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "map",
				Description:  "マップ名（何文字か入力して！）",
				Type:         discordgo.ApplicationCommandOptionString,
				Autocomplete: true,
				Required:     true,
			},
			{
				Name:        "color",
				Description: "なに色にする？（color=None かつ commentなし でマーク削除）",
				Type:        discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "None",
						Value: roanav.MarkedMapColorNone,
					},
					{
						Name:  "Green",
						Value: roanav.MarkedMapColorGreen,
					},
					{
						Name:  "Pink",
						Value: roanav.MarkedMapColorPink,
					},
					{
						Name:  "Purple",
						Value: roanav.MarkedMapColorPurple,
					},
					{
						Name:  "Orange",
						Value: roanav.MarkedMapColorOrange,
					},
					{
						Name:  "Brown",
						Value: roanav.MarkedMapColorBrown,
					},
				},
				Required: true,
			},
			{
				Name:        "comment",
				Description: "マップ名の下に表示するコメント（任意）",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
				MaxLength:   20,
			},
		},
	}
}

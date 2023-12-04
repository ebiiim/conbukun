package roanav

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	kroki "github.com/yuzutech/kroki-go"

	"github.com/ebiiim/conbukun/pkg/ao"
	"github.com/ebiiim/conbukun/pkg/ao/data"
)

var lg zerolog.Logger = log.With().Str("component", "conbukun/pkg/ao/roanav").Logger()

type KrokiPlantUMLPNGPainter struct {
	Client  kroki.Client
	MapData map[string]data.MapData
}

var _ Painter = (*KrokiPlantUMLPNGPainter)(nil)

const (
	DefaultKrokiEndpoint = "https://kroki.io"
	DefaultKrokiTimeout  = 10 * time.Second
)

func NewKrokiPlantUMLPNGPainter(endpoint string, timeout time.Duration, mapData map[string]data.MapData) *KrokiPlantUMLPNGPainter {
	lg.Info().Str("func", "NewKrokiPlantUMLPainter").Msgf("endpoint=%s, timeout=%s, len(mapData)=%d", endpoint, timeout, len(mapData))
	p := &KrokiPlantUMLPNGPainter{
		Client: kroki.New(kroki.Configuration{
			URL:     endpoint,
			Timeout: timeout,
		}),
		MapData: mapData,
	}
	return p
}

func (p *KrokiPlantUMLPNGPainter) Paint(n *Navigation) (path string, err error) {
	lg := lg.With().Str("func", "Paint").Str("Navigation.Name", n.Name).Logger()

	pu, err := p.ToPlantUML(n)
	if err != nil {
		return "", err
	}

	lg.Debug().Str("PlantUML", pu).Msg("sending to kroki...")
	result, err := p.Client.FromString(pu, kroki.PlantUML, kroki.PNG)
	if err != nil {
		return "", err
	}

	path = fmt.Sprintf("roanav-%s-%d.png", n.Name, time.Now().Unix())

	if err := p.Client.WriteToFile(path, result); err != nil {
		return "", err
	}

	return
}

func (p *KrokiPlantUMLPNGPainter) ToPlantUML(n *Navigation) (string, error) {
	lg := lg.With().Str("func", "ToPlantUML").Str("Navigation.Name", n.Name).Logger()

	tmplData, err := p.NavigationToTemplateData(n, time.Now())
	if err != nil {
		lg.Warn().Err(err).Msg("p.NavigationToTemplateData got more than one error (still continuing)")
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (p *KrokiPlantUMLPNGPainter) NavigationToTemplateData(n *Navigation, t time.Time) (templateData, error) {

	templateData := templateData{
		GeneratedAt:  "ERROR",
		Contributors: "ERROR",
		Credit:       "ERROR",
		Agents:       []templateDataAgent{},
		Links:        []templateDataLink{},
	}

	templateData.GeneratedAt = t.Format("2006-01-02 15:04:05 MST")

	contribM := map[string]struct{}{}
	for _, portal := range n.Portals {
		if v, ok := portal.Data[PortalDataKeyUser]; ok {
			contribM[v] = struct{}{}
		}
	}
	contribS := []string{}
	for k := range contribM {
		contribS = append(contribS, k)
	}
	templateData.Contributors = strings.Join(contribS, ", ")

	templateData.Credit = fmt.Sprintf("%s (conbukun@%s)", n.Name, ao.Version)

	var errs error

	agentsM := map[string]templateDataAgent{}
	for _, portal := range n.Portals {
		aliasFrom := toAlias(portal.From)
		aliasTo := toAlias(portal.To)

		agentFrom, err := p.toTemplateDataAgent(portal.From, n.Data)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("p.toTemplateDataAgent(portal.From): %w", err))
			continue
		}
		agentTo, err := p.toTemplateDataAgent(portal.To, n.Data)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("p.toTemplateDataAgent(portal.To): %w", err))
			continue
		}

		agentsM[aliasFrom] = agentFrom
		agentsM[aliasTo] = agentTo

		color := "red" // error
		switch portal.Type {
		case PortalTypeBlue:
			color = "darkblue"
		case PortalTypeYellow:
			color = "orange"
		}

		ts := strings.TrimSuffix(portal.ExpiredAt.Sub(t).Truncate(time.Minute).String(), "0s")
		if ts == "" {
			ts = "<1m"
		}

		templateData.Links = append(templateData.Links, templateDataLink{
			FromAlias: aliasFrom,
			ToAlias:   aliasTo,
			// 2h30m0s -> 2h30m
			Duration: ts,
			Color:    color,
		})
	}

	for _, v := range agentsM {
		templateData.Agents = append(templateData.Agents, v)
	}

	return templateData, errs
}

func (p *KrokiPlantUMLPNGPainter) toTemplateDataAgent(portalID string, navigationData map[string]string) (templateDataAgent, error) {
	d := templateDataAgent{}

	md, ok := p.MapData[portalID]
	if !ok {
		return d, fmt.Errorf("map data not found: %s", portalID)
	}

	mapName := md.DisplayName
	shortName := data.GetMapShortName(md)
	if shortName != "" {
		mapName = fmt.Sprintf("[%s] %s", shortName, mapName) // "[SO] Suyos-Onaytum"
	}

	// d.Name = fmt.Sprintf("%s\\n%s %s", mapName, data.GetMapType(md), data.GetMapTier(md))
	d.Name = fmt.Sprintf("%s (%s)", mapName, data.GetMapTier(md))
	d.Alias = toAlias(portalID)

	color := ""
	switch data.GetMapType(md) {
	case data.MapTypeBlackZone:
		color = "#dimgray;text:white"
	case data.MapTypeRedZone:
		color = "#firebrick;text:white"
	case data.MapTypeYellowZone:
		color = "#gold"
	case data.MapTypeCity, data.MapTypeBlueZone:
		color = "#lightskyblue"
	}

	// check hideout
	if v, ok := navigationData[NavigationDataHideouts]; ok {
		for _, hideout := range strings.Split(v, ",") {
			if hideout == portalID {
				color = "#green;text:white"
				break
			}
		}
	}

	d.Color = color

	return d, nil
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("plantuml").Parse(`
@startuml

skinparam handwritten true
skinparam backgroundColor #F8F5EB
caption "Contributors: {{ .Contributors }}\nTimestamp: {{ .GeneratedAt }}\nAffiliation: {{ .Credit }}"
{{ range $val := .Agents }}
agent "{{ $val.Name }}" as {{ $val.Alias }} {{ $val.Color }}
{{- end}}
{{ range $val := .Links }}
{{ $val.FromAlias }} <-[#{{ $val.Color }}]-> {{ $val.ToAlias }} : "{{ $val.Duration }}"
{{- end}}

@enduml
`))
}

func toAlias(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	s = "map" + s
	return s
}

type templateData struct {
	GeneratedAt  string
	Contributors string
	Credit       string

	Agents []templateDataAgent
	Links  []templateDataLink
}

type templateDataAgent struct {
	Name  string
	Alias string
	Color string
}

type templateDataLink struct {
	FromAlias string
	ToAlias   string
	Duration  string
	Color     string
}

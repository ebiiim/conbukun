package roanav

import (
	"bytes"
	"encoding/json"
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

const (
	PlantUMLStyleAuto     = "auto" // if len(edges)<32 use plantuml else use dot
	PlantUMLStylePlantUML = "plantuml"
	PlantUMLStyleDOT      = "dot"
)

type KrokiPlantUMLPNGPainter struct {
	Client  kroki.Client
	MapData map[string]data.MapData
	Style   string
}

var _ Painter = (*KrokiPlantUMLPNGPainter)(nil)

const (
	DefaultKrokiEndpoint = "https://kroki.io"
	DefaultKrokiTimeout  = 60 * time.Second // kroki.io is so slow :(
)

func NewKrokiPlantUMLPNGPainter(endpoint string, timeout time.Duration, mapData map[string]data.MapData, style string) *KrokiPlantUMLPNGPainter {
	lg.Info().Str("func", "NewKrokiPlantUMLPainter").Msgf("endpoint=%s, timeout=%s, style=%s, len(mapData)=%d", endpoint, timeout, style, len(mapData))
	p := &KrokiPlantUMLPNGPainter{
		Client: kroki.New(kroki.Configuration{
			URL:     endpoint,
			Timeout: timeout,
		}),
		MapData: mapData,
		Style:   style,
	}
	return p
}

func (p *KrokiPlantUMLPNGPainter) Paint(n *Navigation) (path string, err error) {
	lg := lg.With().Str("func", "Paint").Str("style", p.Style).Str("Navigation.Name", n.Name).Logger()

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
	lg := lg.With().Str("func", "ToPlantUML").Str("style", p.Style).Str("Navigation.Name", n.Name).Logger()

	tmplData, err := p.NavigationToTemplateData(n, time.Now())
	if err != nil {
		lg.Warn().Err(err).Msg("p.NavigationToTemplateData got more than one error (still continuing)")
	}

	var buf bytes.Buffer
	switch p.Style {
	case PlantUMLStyleAuto:
		if len(tmplData.Edges) < 32 {
			err = tmplStylePlantUML.Execute(&buf, tmplData)
		} else {
			err = tmplStyleDOT.Execute(&buf, tmplData)
		}
	case PlantUMLStylePlantUML:
		err = tmplStylePlantUML.Execute(&buf, tmplData)
	case PlantUMLStyleDOT:
		err = tmplStyleDOT.Execute(&buf, tmplData)
	default:
		err = fmt.Errorf("unknown style: %s", p.Style)
	}

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (p *KrokiPlantUMLPNGPainter) NavigationToTemplateData(n *Navigation, t time.Time) (templateData, error) {

	templateData := templateData{
		GeneratedAt:  "ERROR",
		Contributors: "ERROR",
		Affiliation:  "ERROR",
		Nodes:        []templateDataNode{},
		Edges:        []templateDataEdge{},
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

	templateData.Affiliation = fmt.Sprintf("%s (conbukun@%s)", n.Name, ao.Version)

	var errs error

	agentsM := map[string]templateDataNode{}
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

		templateData.Edges = append(templateData.Edges, templateDataEdge{
			FromAlias: aliasFrom,
			ToAlias:   aliasTo,
			// 2h30m0s -> 2h30m
			Duration: ts,
			Color:    color,
		})
	}

	for _, v := range agentsM {
		templateData.Nodes = append(templateData.Nodes, v)
	}

	return templateData, errs
}

func (p *KrokiPlantUMLPNGPainter) toTemplateDataAgent(portalID string, navigationData map[string]string) (templateDataNode, error) {
	d := templateDataNode{}

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

	fillColor := ""
	textColor := ""
	switch data.GetMapType(md) {
	case data.MapTypeBlackZone:
		fillColor = "dimgray"
		textColor = "white"
	case data.MapTypeRedZone:
		fillColor = "firebrick"
		textColor = "white"
	case data.MapTypeYellowZone:
		fillColor = "gold"
		textColor = "black"
	case data.MapTypeCity, data.MapTypeBlueZone:
		fillColor = "lightskyblue"
		textColor = "black"
	default:
		fillColor = "whitesmoke"
		textColor = "black"
	}

	// check marked maps
	if v, ok := navigationData[NavigationDataMarkedMaps]; ok {
		var markedMaps []MarkedMap
		if err := json.Unmarshal([]byte(v), &markedMaps); err != nil {
			return d, fmt.Errorf("json.Unmarshal: %w", err)
		}
		for _, m := range markedMaps {

			if m.ID == portalID { // found!

				// 1. add comment
				d.Name = fmt.Sprintf("%s\\n%s", d.Name, m.Comment)

				// 2. change color
				switch m.Color {
				case MarkedMapColorNone:
					// do nothing; herited from above
				case MarkedMapColorGreen:
					fillColor = "darkgreen"
					textColor = "white"
				case MarkedMapColorPink:
					fillColor = "deeppink"
					textColor = "white"
				case MarkedMapColorPurple:
					fillColor = "darkviolet"
					textColor = "white"
				case MarkedMapColorOrange:
					fillColor = "orange"
					textColor = "black"
				case MarkedMapColorBrown:
					fillColor = "saddlebrown"
					textColor = "white"
				default:
					// do nothing (normally unreachable)
				}

				// 3. stop searching
				break

			}

		}
	}

	d.FillColor = fillColor
	d.TextColor = textColor

	return d, nil
}

var (
	tmplStylePlantUML *template.Template
	tmplStyleDOT      *template.Template
)

func init() {
	tmplStylePlantUML = template.Must(template.New("plantuml").Parse(`
@startuml

skinparam handwritten true
skinparam backgroundColor #F8F5EB
caption "Contributors: {{ .Contributors }}\nTimestamp: {{ .GeneratedAt }}\nAffiliation: {{ .Affiliation }}"
{{ range $val := .Nodes }}
agent "{{ $val.Name }}" as {{ $val.Alias }} #{{ $val.FillColor }};text:{{ $val.TextColor }}
{{- end}}
{{ range $val := .Edges }}
{{ $val.FromAlias }} <-[#{{ $val.Color }}]-> {{ $val.ToAlias }} : "{{ $val.Duration }}"
{{- end}}

@enduml
`))

	tmplStyleDOT = template.Must(template.New("dot").Parse(`
@startdot

digraph g {
  
  charset = "UTF-8";
  bgcolor = "#F8F5EB"
  label = "Affiliation: {{ .Affiliation }}\nTimestamp: {{ .GeneratedAt }}\nContributors: {{ .Contributors }}"

  graph [
    overlap = scale;
    layout = fdp;
    labelloc = t;
    labeljust = c;
  ];

  node [
    shape = box;
	style = "rounded,filled";
  ];

  edge [
    dir = both;
    labelfloat = true; 
  ];

{{ range $val := .Nodes }}
  "{{ $val.Alias }}" [label="{{ $val.Name }}", color="darkgray", fillcolor="{{ $val.FillColor }}", fontcolor="{{ $val.TextColor }}"] 
{{- end}}

{{ range $val := .Edges }}
  {{ $val.FromAlias }} -> {{ $val.ToAlias }} [label="{{ $val.Duration }}", color="{{ $val.Color }}"]
{{- end}}

}

@enddot
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
	Affiliation  string

	Nodes []templateDataNode
	Edges []templateDataEdge
}

type templateDataNode struct {
	Name      string
	Alias     string
	FillColor string
	TextColor string
}

type templateDataEdge struct {
	FromAlias string
	ToAlias   string
	Duration  string
	Color     string
}

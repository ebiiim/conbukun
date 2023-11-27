package roanav

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	kroki "github.com/yuzutech/kroki-go"

	"github.com/ebiiim/conbukun/pkg/ao/data"
)

type KrokiPlantUMLPNGPainter struct {
	Client  kroki.Client
	MapData map[string]data.MapData
}

var _ Painter = (*KrokiPlantUMLPNGPainter)(nil)

const (
	DefaultKrokiEndpoint = "https://kroki.io"
	DefaultKrokiTimeout  = 10 * time.Second
)

func NewKrokiPlantUMLPainter(endpoint string, timeout time.Duration, mapData map[string]data.MapData) *KrokiPlantUMLPNGPainter {
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

	pu, err := p.ToPlantUML(n)
	if err != nil {
		return "", err
	}

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
	tmplData, err := p.NavigationToTemplateData(n, time.Now())
	if err != nil {
		// still continue
		// TODO: log
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

	templateData.Credit = n.Name

	var errs error

	agentsM := map[string]templateDataAgent{}
	for _, portal := range n.Portals {
		aliasFrom := toAlias(portal.From)
		aliasTo := toAlias(portal.To)

		agentFrom, err := p.toTemplateDataAgent(portal.From)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("p.toTemplateDataAgent(portal.From): %w", err))
			continue
		}
		agentTo, err := p.toTemplateDataAgent(portal.To)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("p.toTemplateDataAgent(portal.To): %w", err))
			continue
		}

		agentsM[aliasFrom] = agentFrom
		agentsM[aliasTo] = agentTo

		color := "red" // error
		switch portal.Type {
		case PortalTypeBlue:
			color = "blue"
		case PortalTypeYellow:
			color = "gold"
		}
		templateData.Links = append(templateData.Links, templateDataLink{
			FromAlias: aliasFrom,
			ToAlias:   aliasTo,
			Duration:  strings.TrimSuffix(portal.ExpiredAt.Sub(t).Truncate(time.Minute).String(), "0s"), // 2h30m0s -> 2h30m
			Color:     color,
		})
	}

	for _, v := range agentsM {
		templateData.Agents = append(templateData.Agents, v)
	}

	return templateData, errs
}

func (p *KrokiPlantUMLPNGPainter) toTemplateDataAgent(portalID string) (templateDataAgent, error) {
	d := templateDataAgent{}

	md, ok := p.MapData[portalID]
	if !ok {
		return d, fmt.Errorf("map data not found: %s", portalID)
	}

	d.Name = md.DisplayName
	d.Alias = toAlias(portalID)

	return d, nil
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("plantuml").Parse(`
@startuml

caption "Contributors: {{ .Contributors }}\nTimestamp: {{ .GeneratedAt }}\nAffiliation: {{ .Credit }}"
{{ range $val := .Agents }}
agent "{{ $val.Name }}" as {{ $val.Alias }}
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
}

type templateDataLink struct {
	FromAlias string
	ToAlias   string
	Duration  string
	Color     string
}

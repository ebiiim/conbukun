package roanav

import (
	"fmt"
	"strings"
	"time"

	"github.com/ebiiim/conbukun/pkg/ao/data"
)

// TODO: add brief interface?

func BriefPortal(p *Portal, mapData map[string]data.MapData) string {

	var sb strings.Builder

	from := p.From
	mdFrom, ok := mapData[p.From]
	if ok {
		from = mdFrom.DisplayName
	}
	to := p.To
	mdTo, ok := mapData[p.To]
	if ok {
		to = mdTo.DisplayName
	}

	typ := "E"
	switch p.Type {
	case PortalTypeBlue:
		typ = "B"
	case PortalTypeYellow:
		typ = "Y"
	}

	sb.WriteString(fmt.Sprintf("%s <-[%s|%s]-> %s", from, typ, strings.TrimSuffix(time.Until(p.ExpiredAt).Truncate(time.Minute).String(), "0s"), to))

	u, ok := p.Data[PortalDataKeyUser]
	if ok {
		sb.WriteString(fmt.Sprintf("  (%s)", u))
	}

	return sb.String()
}

func BriefNavigation(n *Navigation, mapData map[string]data.MapData) string {
	n.DeleteExpiredPortals()

	var sb strings.Builder

	// sb.WriteString(fmt.Sprintf("Name: %s\n", n.Name))
	// sb.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05 MST")))
	// sb.WriteString("Portals:\n")
	for _, p := range n.Portals {
		// sb.WriteString(fmt.Sprintf("  %s\n", BriefPortal(p, mapData)))
		sb.WriteString(fmt.Sprintf("- %s\n", BriefPortal(p, mapData)))
	}
	return sb.String()
}

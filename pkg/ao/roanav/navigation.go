package roanav

import (
	"strings"
	"time"
)

const (
	PortalTypeBlue   = "blue"
	PortalTypeYellow = "yellow"
)

const (
	PortalDataKeyUser = "user"
)

type Portal struct {
	// From is the map with a smaller name. Note that Route is an undirected path.
	From string
	// To is the map with a larger name. Note that Route is an undirected path.
	To string
	// Type is the portal type.
	Type      string
	ExpiredAt time.Time

	// Data contains additional data.
	// E.g. the user who added the portal.
	Data map[string]string
}

// NewPortal initializes a Portal.
// map1 and map2 are sorted in alphabetical order, then the smaller one is set to Portal.From and the other is set to Portal.To.
func NewPortal(map1, map2, typ string, expiredAt time.Time, data map[string]string) *Portal {
	from := map1
	to := map2
	if strings.Compare(map1, map2) > 0 {
		from = map2
		to = map1
	}
	if data == nil {
		data = map[string]string{}
	}
	return &Portal{
		From:      from,
		To:        to,
		Type:      typ,
		ExpiredAt: expiredAt,
		Data:      data,
	}
}

type Navigation struct {
	// Name is the name of the navigation, usually the name of the channel+guild.
	Name string
	// Portals is the list of portals.
	Portals []*Portal

	// Data contains additional data.
	Data map[string]string
}

// NewNavigation initializes a Navigation.
func NewNavigation(name string) *Navigation {
	return &Navigation{
		Name:    name,
		Portals: nil,
		Data:    map[string]string{},
	}
}

// AddPortal adds a portal to the navigation.
func (n *Navigation) AddPortal(p *Portal) {
	// Update ExpiredAt and Data if the portal already exists.
	for _, portal := range n.Portals {
		if portal.From == p.From && portal.To == p.To && portal.Type == p.Type {
			portal.ExpiredAt = p.ExpiredAt
			portal.Data = p.Data
			return
		}
	}
	n.Portals = append(n.Portals, p)
}

// DeleteExpiredPortals deletes expired portals.
func (n *Navigation) DeleteExpiredPortals() {
	now := time.Now()

	newPortals := make([]*Portal, 0, len(n.Portals))
	for i, p := range n.Portals {
		if p.ExpiredAt.After(now) {
			newPortals = append(newPortals, n.Portals[i])
		}
	}
	n.Portals = newPortals
}

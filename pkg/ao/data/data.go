package data

import (
	"bytes"
	"encoding/json"
	"strings"
	"unicode"

	_ "embed"
)

type world struct {
	World struct {
		Clusters struct {
			Cluster []struct {
				ID           string `json:"@id"`
				File         string `json:"@file"`
				DisplayName  string `json:"@displayName"`
				Type         string `json:"@type"`
				Distribution struct {
					Resource json.RawMessage `json:"resource"`
				} `json:"distribution"`
			} `json:"cluster"`
		} `json:"clusters"`
	} `json:"world"`
}

type MapData struct {
	ID           string `json:"@id"`
	File         string `json:"@file"`
	DisplayName  string `json:"@displayName"`
	Type         string `json:"@type"`
	Distribution struct {
		Resources []Resource `json:"resource"`
	} `json:"distribution"`
}

type Resource struct {
	Name  string `json:"@name"`
	Tier  string `json:"@tier"`
	Count string `json:"@count"`
}

//go:embed ao-bin-dumps/cluster/world.json
var worldjson []byte

var (
	worldData world
	Maps      = map[string]MapData{}
)

func init() {
	buf := bytes.NewBuffer(worldjson)
	worldData = world{}
	if err := json.NewDecoder(buf).Decode(&worldData); err != nil {
		panic(err)
	}
	for _, c := range worldData.World.Clusters.Cluster {

		mapdata := MapData{
			ID:          c.ID,
			File:        c.File,
			DisplayName: c.DisplayName,
			Type:        c.Type,
		}

		// try single resource
		if mapdata.Distribution.Resources == nil {
			resource := Resource{}
			if err := json.Unmarshal(c.Distribution.Resource, &resource); err == nil {
				mapdata.Distribution.Resources = []Resource{resource}
			}
		}

		// try multiple resources
		if mapdata.Distribution.Resources == nil {
			resources := []Resource{}
			if err := json.Unmarshal(c.Distribution.Resource, &resources); err == nil {
				mapdata.Distribution.Resources = resources
			}
		}

		Maps[c.ID] = mapdata
	}
}

func GetMapDataFromName(displayName string) (MapData, bool) {
	for _, m := range Maps {
		if m.DisplayName == displayName {
			return m, true
		}
	}
	return MapData{}, false
}

// MapTypes
// cat pkg/ao/data/ao-bin-dumps/cluster/world.json | jq .world.clusters | jq '.cluster[]."@type"' | sort | uniq
const (
	MapTypeRawOpenpvpBlack1                        = "OPENPVP_BLACK_1"
	MapTypeRawOpenpvpBlack2                        = "OPENPVP_BLACK_2"
	MapTypeRawOpenpvpBlack3                        = "OPENPVP_BLACK_3"
	MapTypeRawOpenpvpBlack4                        = "OPENPVP_BLACK_4"
	MapTypeRawOpenpvpBlack5                        = "OPENPVP_BLACK_5"
	MapTypeRawOpenpvpBlack6                        = "OPENPVP_BLACK_6"
	MapTypeRawOpenpvpRed                           = "OPENPVP_RED"
	MapTypeRawOpenpvpYellow                        = "OPENPVP_YELLOW"
	MapTypeRawSafearea                             = "SAFEAREA" // blue
	MapTypeRawTunnelBlackHigh                      = "TUNNEL_BLACK_HIGH"
	MapTypeRawTunnelBlackLow                       = "TUNNEL_BLACK_LOW"
	MapTypeRawTunnelBlackMedium                    = "TUNNEL_BLACK_MEDIUM"
	MapTypeRawTunnelDeep                           = "TUNNEL_DEEP"
	MapTypeRawTunnelDeepRaid                       = "TUNNEL_DEEP_RAID"
	MapTypeRawTunnelHideout                        = "TUNNEL_HIDEOUT"
	MapTypeRawTunnelHideoutDeep                    = "TUNNEL_HIDEOUT_DEEP"
	MapTypeRawTunnelHigh                           = "TUNNEL_HIGH"
	MapTypeRawTunnelLow                            = "TUNNEL_LOW"
	MapTypeRawTunnelMedium                         = "TUNNEL_MEDIUM"
	MapTypeRawTunnelRoyal                          = "TUNNEL_ROYAL"
	MapTypeRawPlayercityBlack                      = "PLAYERCITY_BLACK"                        // Brecilien, Rests
	MapTypeRawPlayercityBlackRoyal                 = "PLAYERCITY_BLACK_ROYAL"                  // Caerleon
	MapTypeRawPlayercityBlackPortalcityNofurniture = "PLAYERCITY_BLACK_PORTALCITY_NOFURNITURE" // Portals
	MapTypeRawPlayercitySafearea01                 = "PLAYERCITY_SAFEAREA_01"                  // Martlock, Thetford, Bridgewatch
	MapTypeRawPlayercitySafearea02                 = "PLAYERCITY_SAFEAREA_02"                  // Fort Sterling, Lymhurst
)

// User friendly map types
const (
	MapTypeUnknown    = "Unknown"
	MapTypeBlackZone  = "Black"
	MapTypeRedZone    = "Red"
	MapTypeYellowZone = "Yellow"
	MapTypeBlueZone   = "Blue"
	MapTypeAvalon     = "Avalon"
	MapTypeCity       = "City"
)

var (
	MapTypesBlackZone = []string{
		MapTypeRawOpenpvpBlack1,
		MapTypeRawOpenpvpBlack2,
		MapTypeRawOpenpvpBlack3,
		MapTypeRawOpenpvpBlack4,
		MapTypeRawOpenpvpBlack5,
		MapTypeRawOpenpvpBlack6,
	}
	MapTypesRedZone = []string{
		MapTypeRawOpenpvpRed,
	}
	MapTypesYellowZone = []string{
		MapTypeRawOpenpvpYellow,
	}
	MapTypesBlueZone = []string{
		MapTypeRawSafearea,
	}
	MapTypesAvalon = []string{
		MapTypeRawTunnelBlackHigh,
		MapTypeRawTunnelBlackLow,
		MapTypeRawTunnelBlackMedium,
		MapTypeRawTunnelDeep,
		MapTypeRawTunnelDeepRaid,
		MapTypeRawTunnelHideout,
		MapTypeRawTunnelHideoutDeep,
		MapTypeRawTunnelHigh,
		MapTypeRawTunnelLow,
		MapTypeRawTunnelMedium,
		MapTypeRawTunnelRoyal,
	}
	MapTypesCity = []string{
		MapTypeRawPlayercityBlack,
		MapTypeRawPlayercityBlackRoyal,
		MapTypeRawPlayercityBlackPortalcityNofurniture,
		MapTypeRawPlayercitySafearea01,
		MapTypeRawPlayercitySafearea02,
	}
)

func GetMapType(mapData MapData) string {
	for _, t := range MapTypesBlackZone {
		if t == mapData.Type {
			return MapTypeBlackZone
		}
	}
	for _, t := range MapTypesRedZone {
		if t == mapData.Type {
			return MapTypeRedZone
		}
	}
	for _, t := range MapTypesYellowZone {
		if t == mapData.Type {
			return MapTypeYellowZone
		}
	}
	for _, t := range MapTypesBlueZone {
		if t == mapData.Type {
			return MapTypeBlueZone
		}
	}
	for _, t := range MapTypesAvalon {
		if t == mapData.Type {
			return MapTypeAvalon
		}
	}
	for _, t := range MapTypesCity {
		if t == mapData.Type {
			return MapTypeCity
		}
	}
	return MapTypeUnknown
}

// Tiers
const (
	TierErr = "T0(Error)"
	Tier1   = "T1"
	Tier2   = "T2"
	Tier3   = "T3"
	Tier4   = "T4"
	Tier5   = "T5"
	Tier6   = "T6"
	Tier7   = "T7"
	Tier8   = "T8"
)

func GetMapTier(mapData MapData) string {
	switch {
	case strings.Contains(mapData.File, "_T1_"):
		return Tier1
	case strings.Contains(mapData.File, "_T2_"):
		return Tier2
	case strings.Contains(mapData.File, "_T3_"):
		return Tier3
	case strings.Contains(mapData.File, "_T4_"):
		return Tier4
	case strings.Contains(mapData.File, "_T5_"):
		return Tier5
	case strings.Contains(mapData.File, "_T6_"):
		return Tier6
	case strings.Contains(mapData.File, "_T7_"):
		return Tier7
	case strings.Contains(mapData.File, "_T8_"):
		return Tier8
	default:
		return TierErr
	}
}

func GetMapShortName(mapData MapData) string {
	if GetMapType(mapData) != MapTypeAvalon {
		return ""
	}

	displayName := mapData.DisplayName
	shortName := ""

	for _, c := range displayName {
		if unicode.IsUpper(c) {
			shortName += string(c)
		}
	}

	return shortName
}

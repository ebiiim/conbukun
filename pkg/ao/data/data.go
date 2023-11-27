package data

import (
	"bytes"
	"encoding/json"
	"strings"

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
	MapTypeOpenpvpBlack1     = "OPENPVP_BLACK_1"
	MapTypeOpenpvpBlack2     = "OPENPVP_BLACK_2"
	MapTypeOpenpvpBlack3     = "OPENPVP_BLACK_3"
	MapTypeOpenpvpBlack4     = "OPENPVP_BLACK_4"
	MapTypeOpenpvpBlack5     = "OPENPVP_BLACK_5"
	MapTypeOpenpvpBlack6     = "OPENPVP_BLACK_6"
	MapTypeOpenpvpRed        = "OPENPVP_RED"
	MapTypeOpenpvpYellow     = "OPENPVP_YELLOW"
	MapTypeSafearea          = "SAFEAREA"
	MapTypeTunnelBlackHigh   = "TUNNEL_BLACK_HIGH"
	MapTypeTunnelBlackLow    = "TUNNEL_BLACK_LOW"
	MapTypeTunnelBlackMedium = "TUNNEL_BLACK_MEDIUM"
	MapTypeTunnelDeep        = "TUNNEL_DEEP"
	MapTypeTunnelDeepRaid    = "TUNNEL_DEEP_RAID"
	MapTypeTunnelHideout     = "TUNNEL_HIDEOUT"
	MapTypeTunnelHideoutDeep = "TUNNEL_HIDEOUT_DEEP"
	MapTypeTunnelHigh        = "TUNNEL_HIGH"
	MapTypeTunnelLow         = "TUNNEL_LOW"
	MapTypeTunnelMedium      = "TUNNEL_MEDIUM"
	MapTypeTunnelRoyal       = "TUNNEL_ROYAL"
)

var (
	MapTypesBlack = []string{
		MapTypeOpenpvpBlack1,
		MapTypeOpenpvpBlack2,
		MapTypeOpenpvpBlack3,
		MapTypeOpenpvpBlack4,
		MapTypeOpenpvpBlack5,
		MapTypeOpenpvpBlack6,
	}
	MapTypesRed = []string{
		MapTypeOpenpvpRed,
	}
	MapTypesYellow = []string{
		MapTypeOpenpvpYellow,
	}
	MapTypesBlue = []string{
		MapTypeSafearea,
	}
	MapTypesTunnel = []string{
		MapTypeTunnelBlackHigh,
		MapTypeTunnelBlackLow,
		MapTypeTunnelBlackMedium,
		MapTypeTunnelDeep,
		MapTypeTunnelDeepRaid,
		MapTypeTunnelHideout,
		MapTypeTunnelHideoutDeep,
		MapTypeTunnelHigh,
		MapTypeTunnelLow,
		MapTypeTunnelMedium,
		MapTypeTunnelRoyal,
	}
)

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

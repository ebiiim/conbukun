package data

import (
	"bytes"
	_ "embed"
	"encoding/json"
)

type world struct {
	World struct {
		Clusters struct {
			Cluster []struct {
				ID           string `json:"@id"`
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

package data

import (
	"testing"

	_ "embed"
)

func TestGetMapTier(t *testing.T) {
	type args struct {
		mapData MapData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "f_t7", args: args{mapData: getMapData("Dryvein End")}, want: Tier7},
		{name: "f_t6", args: args{mapData: getMapData("Frostspring Passage")}, want: Tier6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMapTier(tt.args.mapData); got != tt.want {
				t.Errorf("GetMapTier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getMapData(name string) MapData {
	for _, m := range Maps {
		if m.DisplayName == name {
			return m
		}
	}
	panic("map not found")
}

func TestGetMapShortName(t *testing.T) {
	type args struct {
		mapData MapData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "f", args: args{mapData: getMapData("Dryvein End")}, want: ""},
		{name: "f", args: args{mapData: getMapData("Frostspring Passage")}, want: ""},
		{name: "t_2", args: args{mapData: getMapData("Suyos-Onaytum")}, want: "SO"},
		{name: "t_3", args: args{mapData: getMapData("Quaent-In-Odesum")}, want: "QIO"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMapShortName(tt.args.mapData); got != tt.want {
				t.Errorf("GetMapShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMapType(t *testing.T) {
	type args struct {
		mapData MapData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "bz", args: args{mapData: getMapData("Dryvein End")}, want: MapTypeBlackZone},
		{name: "rz", args: args{mapData: getMapData("Creag Morr")}, want: MapTypeRedZone},
		{name: "yz", args: args{mapData: getMapData("Sleetwater Basin")}, want: MapTypeYellowZone},
		{name: "safe", args: args{mapData: getMapData("Fog Fen")}, want: MapTypeBlueZone},
		{name: "city", args: args{mapData: getMapData("Bridgewatch")}, want: MapTypeCity},
		{name: "city", args: args{mapData: getMapData("Caerleon")}, want: MapTypeCity},
		{name: "city", args: args{mapData: getMapData("Morgana's Rest")}, want: MapTypeCity},
		{name: "city", args: args{mapData: getMapData("Brecilien")}, want: MapTypeCity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMapType(tt.args.mapData); got != tt.want {
				t.Errorf("GetMapType() = %v, want %v, md=%+v", got, tt.want, tt.args.mapData)
			}
		})
	}
}

package data

import (
	_ "embed"
	"testing"
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
		{
			name: "2345(Dryvein End)",
			args: args{mapData: Maps["2345"]},
			want: Tier7,
		},
		{
			name: "4336(Frostspring Passage)",
			args: args{mapData: Maps["4336"]},
			want: Tier6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMapTier(tt.args.mapData); got != tt.want {
				t.Errorf("GetMapTier() = %v, want %v", got, tt.want)
			}
		})
	}
}

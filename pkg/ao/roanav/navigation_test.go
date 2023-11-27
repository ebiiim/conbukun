package roanav

import (
	"reflect"
	"testing"
	"time"
)

func TestNewPortal(t *testing.T) {
	var (
		timeNow = time.Now()
	)
	type args struct {
		map1      string
		map2      string
		typ       string
		expiredAt time.Time
	}
	tests := []struct {
		name string
		args args
		want *Portal
	}{
		{
			name: "normal",
			args: args{map1: "a", map2: "b", typ: PortalTypeBlue, expiredAt: timeNow},
			want: &Portal{From: "a", To: "b", Type: PortalTypeBlue, ExpiredAt: timeNow},
		},
		{
			name: "sort",
			args: args{map1: "b", map2: "a", typ: PortalTypeYellow, expiredAt: timeNow},
			want: &Portal{From: "a", To: "b", Type: PortalTypeYellow, ExpiredAt: timeNow},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPortal(tt.args.map1, tt.args.map2, tt.args.typ, tt.args.expiredAt, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPortal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigation_AddPortal(t *testing.T) {
	var (
		time1 = time.Now().Add(-time.Hour)
		time2 = time.Now().Add(time.Hour)
	)
	type fields struct {
		Name    string
		Portals []*Portal
	}
	type args struct {
		p *Portal
	}
	tests := []struct {
		name   string
		args   args
		before fields
		after  fields
	}{
		{
			name: "0to1",
			args: args{
				p: NewPortal("a", "b", PortalTypeBlue, time1, nil),
			},
			before: fields{
				Name:    "test",
				Portals: []*Portal{},
			},
			after: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, time1, nil),
				},
			},
		},
		{
			name: "1to2",
			args: args{
				p: NewPortal("c", "d", PortalTypeYellow, time2, nil),
			},
			before: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, time1, nil),
				},
			},
			after: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, time1, nil),
					NewPortal("c", "d", PortalTypeYellow, time2, nil),
				},
			},
		},
		{
			name: "1to1",
			args: args{
				p: NewPortal("a", "b", PortalTypeBlue, time2, map[string]string{PortalDataKeyUser: "user1"}),
			},
			before: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, time1, nil),
				},
			},
			after: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, time2, map[string]string{PortalDataKeyUser: "user1"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Navigation{
				Name:    tt.before.Name,
				Portals: tt.before.Portals,
			}
			n.AddPortal(tt.args.p)
			if !reflect.DeepEqual(n.Portals, tt.after.Portals) {
				t.Errorf("Navigation.AddPortal() = %v, want %v", n.Portals, tt.after.Portals)
			}
		})
	}
}

func TestNavigation_DeleteExpiredPortals(t *testing.T) {
	var (
		timePast   = time.Now().Add(-time.Hour)
		timeFuture = time.Now().Add(time.Hour)
	)
	type fields struct {
		Name    string
		Portals []*Portal
	}
	tests := []struct {
		name   string
		before fields
		after  fields
	}{
		{
			name: "no-expired",
			before: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, timeFuture, nil),
					NewPortal("c", "d", PortalTypeBlue, timeFuture, nil),
				},
			},
			after: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, timeFuture, nil),
					NewPortal("c", "d", PortalTypeBlue, timeFuture, nil),
				},
			},
		},
		{
			name: "expired",
			before: fields{
				Name: "test",
				Portals: []*Portal{
					NewPortal("a", "b", PortalTypeBlue, timePast, nil),
					NewPortal("c", "d", PortalTypeBlue, timePast, nil),
				},
			},
			after: fields{
				Name:    "test",
				Portals: []*Portal{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Navigation{
				Name:    tt.before.Name,
				Portals: tt.before.Portals,
			}
			n.DeleteExpiredPortals()
			if !reflect.DeepEqual(n.Portals, tt.after.Portals) {
				t.Errorf("Navigation.DeleteExpiredPortals() = %v, want %v", n.Portals, tt.after.Portals)
			}
		})
	}
}

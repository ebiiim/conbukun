package handlers

import (
	"bytes"
	"io"
	"testing"

	"github.com/ebiiim/conbukun/pkg/ao/roanav"
)

func TestROANavHandler_ExportNavigations(t *testing.T) {
	type fields struct {
		navigations      map[string]*roanav.Navigation
		MapNameCompleter *MapNameCompleter
		saveFile         string
	}
	tests := []struct {
		name    string
		fields  fields
		wantW   string
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				navigations: map[string]*roanav.Navigation{},
			},
			wantW:   "{}\n",
			wantErr: false,
		},
		{
			name: "1",
			fields: fields{
				navigations: map[string]*roanav.Navigation{
					"foo": {Name: "foo", Portals: nil, Data: map[string]string{}},
				},
			},
			wantW:   "{\"foo\":{\"name\":\"foo\",\"portals\":null,\"data\":{}}}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ROANavHandler{
				MapNameCompleter: tt.fields.MapNameCompleter,
				saveFile:         tt.fields.saveFile,
			}
			for k, v := range tt.fields.navigations {
				h.navigations.Store(k, v)
			}

			w := &bytes.Buffer{}
			if err := h.ExportNavigations(w); (err != nil) != tt.wantErr {
				t.Errorf("ROANavHandler.ExportNavigations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ROANavHandler.ExportNavigations() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestROANavHandler_ImportNavigations(t *testing.T) {
	type fields struct {
		MapNameCompleter *MapNameCompleter
		saveFile         string
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    map[string]*roanav.Navigation
	}{
		{
			"empty", fields{},
			args{bytes.NewBufferString("{}")},
			false,
			map[string]*roanav.Navigation{},
		},
		{
			"1", fields{},
			args{bytes.NewBufferString("{\"foo\":{\"name\":\"foo\",\"portals\":null,\"data\":{}}}\n")},
			false,
			map[string]*roanav.Navigation{
				"foo": {Name: "foo", Portals: nil, Data: map[string]string{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ROANavHandler{
				MapNameCompleter: tt.fields.MapNameCompleter,
				saveFile:         tt.fields.saveFile,
			}
			if err := h.ImportNavigations(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("ROANavHandler.ImportNavigations() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := map[string]*roanav.Navigation{}
			h.navigations.Range(func(k, v interface{}) bool {
				got[k.(string)] = v.(*roanav.Navigation)
				return true
			})

			for k, v := range tt.want {
				if got[k] == nil {
					t.Errorf("ROANavHandler.ImportNavigations() = %v, want %v", got, tt.want)
				}
				if got[k].Name != v.Name {
					t.Errorf("ROANavHandler.ImportNavigations() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

package editorials

import (
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

func TestLoad(t *testing.T) {
	const figures = `[{"id":"steps","title":"Шаги","frames":[{"caption":"l = 0, r = 1.","rows":[{"label":"a","cells":[2,"x"],"pointers":{"r":1,"l":0},"window":[0,1]}]}]}]`
	tests := []struct {
		name    string
		files   fstest.MapFS
		want    map[string]Editorial
		wantErr string
	}{
		{
			// A number cell decodes as float64, and the pointers keep the order of the file.
			name:  "HTML with its figures",
			files: fstest.MapFS{"cf-1-a.html": {Data: []byte("<p>Разбор.</p>\n")}, "cf-1-a.json": {Data: []byte(figures)}},
			want: map[string]Editorial{"cf-1-a": {HTML: "<p>Разбор.</p>\n", Figures: []sets.Figure{{ID: "steps", Title: "Шаги", Frames: []sets.FigureFrame{{Caption: "l = 0, r = 1.", Rows: []sets.FigureRow{
				{Label: "a", Cells: []any{2.0, "x"}, Pointers: &sets.Pointers{{Name: "r", Cell: 1}, {Name: "l", Cell: 0}}, Window: &[2]int{0, 1}},
			}}}}}}},
		},
		{name: "HTML without figures", files: fstest.MapFS{"cf-1-a.html": {Data: []byte("<p>Разбор.</p>")}}, wantErr: "cf-1-a.json"},
		{name: "figures without HTML", files: fstest.MapFS{"cf-1-a.html": {Data: []byte("<p>Разбор.</p>")}, "cf-1-a.json": {Data: []byte("[]")}, "cf-2-a.json": {Data: []byte("[]")}}, wantErr: "cf-2-a.html"},
		{name: "unknown figure field", files: fstest.MapFS{"cf-1-a.html": {Data: []byte("<p>Разбор.</p>")}, "cf-1-a.json": {Data: []byte(strings.Replace(figures, "window", "mark", 1))}}, wantErr: `unknown field "mark"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := load(tt.files)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want an error with %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("load = %#v, want %#v", got, tt.want)
			}
		})
	}
}

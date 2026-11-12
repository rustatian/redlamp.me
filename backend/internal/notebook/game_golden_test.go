package notebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"
)

// goldenGame is testdata/game.json. The one-time generator in testdata/golden ran the Worker on the rows of each scenario
// and stored its answer of GET /leetcode/api/game. Now is the moment of that run: the dates of the rows are relative to it.
type goldenGame struct {
	Now       time.Time `json:"now"`
	Scenarios []struct {
		Name string `json:"name"`
		Role string `json:"role"`
		Rows []struct {
			Table string         `json:"table"`
			Raw   map[string]any `json:"values"`
		} `json:"rows"`
		Game   json.RawMessage `json:"game"`
		Visits []string        `json:"visits"`
		Badges []string        `json:"badges"`
	} `json:"scenarios"`
}

// The XP, the level, the awards, the late tasks and the missed days of the Go server equal the answers of the Worker for the same records on the same day.
func TestGameEqualsTheGoldenFile(t *testing.T) {
	data, err := os.ReadFile("testdata/game.json")
	if err != nil {
		t.Fatal(err)
	}
	var golden goldenGame
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&golden); err != nil {
		t.Fatal(err)
	}
	if len(golden.Scenarios) == 0 {
		t.Fatal("testdata/game.json has no scenario")
	}
	content := goldenSets(t)
	for _, tt := range golden.Scenarios {
		t.Run(tt.Name, func(t *testing.T) {
			pool := database(t)
			for _, row := range tt.Rows {
				columns := make([]string, 0, len(row.Raw))
				for column := range row.Raw {
					columns = append(columns, column)
				}
				sort.Strings(columns)
				placeholders, values := make([]string, len(columns)), make([]any, len(columns))
				for i, column := range columns {
					placeholders[i] = fmt.Sprintf("$%d", i+1)
					values[i] = row.Raw[column]
					// A JSON number is the value of a bigint column.
					if number, ok := values[i].(json.Number); ok {
						integer, err := number.Int64()
						if err != nil {
							t.Fatal(err)
						}
						values[i] = integer
					}
				}
				exec(t, pool, "INSERT INTO "+row.Table+" ("+strings.Join(columns, ",")+") VALUES ("+strings.Join(placeholders, ",")+")", values...)
			}
			s := newServerWithSets(t, pool, content)
			s.now = func() time.Time { return golden.Now }
			w := request(t, s, "/game", http.MethodGet, nil, loginAs(t, s, tt.Role))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, body %s", w.Code, w.Body.String())
			}
			// The golden file has the JSON of the Worker with indentation. Without it, the text must equal the body.
			var want bytes.Buffer
			if err := json.Compact(&want, tt.Game); err != nil {
				t.Fatal(err)
			}
			if got := w.Body.String(); got != want.String() {
				t.Errorf("game =\n%s\nwant\n%s", got, want.String())
			}
			var visits, badges []string
			if err := pool.QueryRow(t.Context(), "SELECT coalesce(array_agg(day ORDER BY day),'{}') FROM visits").Scan(&visits); err != nil {
				t.Fatal(err)
			}
			if err := pool.QueryRow(t.Context(), "SELECT coalesce(array_agg(id ORDER BY id),'{}') FROM earned_badges").Scan(&badges); err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(visits, tt.Visits) {
				t.Errorf("visits = %v, want %v", visits, tt.Visits)
			}
			if !slices.Equal(badges, tt.Badges) {
				t.Errorf("stored awards = %v, want %v", badges, tt.Badges)
			}
		})
	}
}

package content

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/sets/settest"
)

var update = flag.Bool("update", false, "write the file sets.json of package settest from the Go content")

// The file sets.json of package settest is the JSON copy of the learning sets. The route tests, the web tests and the browser tests read it.
// It equals the Go content: every text, problem, card, figure and lesson, in the same order.
// After a change of the content, run this test with -update. It writes the file, and the next run compares.
func TestSetsEqualTheFile(t *testing.T) {
	if *update {
		var text bytes.Buffer
		encoder := json.NewEncoder(&text)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "\t")
		if err := encoder.Encode(Sets); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile("../settest/sets.json", text.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		// The test binary has the earlier file inside it.
		t.Skip("wrote ../settest/sets.json; run the test again without -update")
	}
	file := settest.Golden(t)
	if len(file) != len(Sets) {
		t.Fatalf("the file has %d learning sets, and the content has %d", len(file), len(Sets))
	}
	for i, set := range Sets {
		t.Run(set.Slug, func(t *testing.T) {
			want := file[i]
			if len(set.Chapters) != len(want.Chapters) {
				t.Fatalf("chapters = %d, and the file has %d", len(set.Chapters), len(want.Chapters))
			}
			for j, chapter := range set.Chapters {
				if !reflect.DeepEqual(chapter, want.Chapters[j]) {
					t.Errorf("chapter %d differs from the chapter of the file", j+1)
				}
			}
			set.Chapters, want.Chapters = nil, nil
			if !reflect.DeepEqual(set, want) {
				t.Errorf("set = %+v, and the file has %+v", set, want)
			}
		})
	}
}

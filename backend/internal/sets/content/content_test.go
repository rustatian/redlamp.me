package content

import (
	"reflect"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
	"github.com/rustatian/redlamp.me/backend/internal/sets/settest"
)

// Each learning set in Go equals the learning set of the Worker with the same slug: every text, problem, card, figure and lesson.
func TestSetsEqualTheContentOfTheWorker(t *testing.T) {
	golden := map[string]sets.LearningSet{}
	for _, set := range settest.Golden(t) {
		golden[set.Slug] = set
	}
	if len(Sets) == 0 {
		t.Fatal("the content has no learning set")
	}
	for _, set := range Sets {
		t.Run(set.Slug, func(t *testing.T) {
			want, ok := golden[set.Slug]
			if !ok {
				t.Fatalf("the Worker has no learning set %q", set.Slug)
			}
			if len(set.Chapters) != len(want.Chapters) {
				t.Fatalf("chapters = %d, want %d", len(set.Chapters), len(want.Chapters))
			}
			for i, chapter := range set.Chapters {
				if !reflect.DeepEqual(chapter, want.Chapters[i]) {
					t.Errorf("chapter %d differs from the chapter of the Worker", i+1)
				}
			}
			set.Chapters, want.Chapters = nil, nil
			if !reflect.DeepEqual(set, want) {
				t.Errorf("set = %+v, want %+v", set, want)
			}
		})
	}
}

// The order of the learning sets is the order of the Worker.
func TestSetsKeepTheOrderOfTheWorker(t *testing.T) {
	var order []string
	for _, set := range settest.Golden(t) {
		for _, own := range Sets {
			if own.Slug == set.Slug {
				order = append(order, set.Slug)
			}
		}
	}
	for i, set := range Sets {
		if i >= len(order) || order[i] != set.Slug {
			t.Errorf("set %d is %q, and the order of the Worker is %v", i+1, set.Slug, order)
		}
	}
}

package notebook

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The Worker builds a recall pick as the whole set task with its slug, its label and its problem of the library
// (recallPicks in src/leetcode/recall.ts). So a pick of a learning set also has the statement and the tests of the task.
func TestRecallPickCarriesTheWholeTask(t *testing.T) {
	pool := database(t)
	s := newServerWithSets(t, pool, goldenSets(t))
	picks, err := s.recallPicks(context.Background(), pool, "2026-10-07")
	if err != nil {
		t.Fatal(err)
	}
	if len(picks) != 3 {
		t.Fatalf("picks = %d, want the 3 picks of the two-pointers pool for an empty library", len(picks))
	}
	tasks := map[string]sets.Task{}
	for _, task := range s.setTasks("two-pointers") {
		tasks[task.Slug()] = task
	}
	for _, pick := range picks {
		t.Run(pick.Slug, func(t *testing.T) {
			task, ok := tasks[pick.Slug]
			if !ok {
				t.Fatalf("the pick %q is not a task of the set", pick.Slug)
			}
			want := map[string]any{}
			data, err := json.Marshal(task)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(data, &want); err != nil {
				t.Fatal(err)
			}
			want["slug"], want["label"], want["problemId"] = task.Slug(), task.Label(), nil
			got := map[string]any{}
			if data, err = json.Marshal(pick); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("pick = %v, want %v", got, want)
			}
			if _, ok := got["statement"]; !ok {
				t.Errorf("the pick has no statement, and the task %q has one", pick.Slug)
			}
		})
	}
}

package notebook

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The cases of test/recall.test.ts.

// The prepared statements of test/recall.test.ts.
const (
	recallInsertProblem = "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,next_review_date,archived_at,created_at,updated_at) VALUES ($1,$2,'Problem','medium','[]','','',$3,$4,'now','now')"
	recallInsertReview  = "INSERT INTO reviews (id,problem_id,result,note,reviewed_on,created_at) VALUES ($1,$2,'with_help','',$3,'now')"
)

// The homework slugs of test/recall.test.ts.
const (
	recallXor        = "sum-of-all-subset-xor-totals"
	recallLetterCase = "letter-case-permutation"
	recallPhone      = "letter-combinations-of-a-phone-number"
)

// recallSetTopics are the topics of the first three tasks of the set two-pointers, as setTopics of test/recall.test.ts.
var recallSetTopics = map[string][]string{
	"two-sum-ii-input-array-is-sorted": {"Array", "Two Pointers"},
	"valid-palindrome":                 {"String", "Two Pointers"},
	"max-number-of-k-sum-pairs":        {"Array", "Two Pointers", "Sorting"},
}

// recallProblem is a library row of test/recall.test.ts. An empty due or reviewedOn is no value.
type recallProblem struct {
	slug, due, reviewedOn string
	archived              bool
}

// recallSlug returns the slug of a set task URL, as taskSlug of src/leetcode/sets.
func recallSlug(url string) string {
	if rest, ok := strings.CutPrefix(url, "https://codeforces.com/problemset/problem/"); ok {
		contest, index, _ := strings.Cut(rest, "/")
		return "cf-" + contest + "-" + strings.ToLower(index)
	}
	return strings.TrimSuffix(strings.TrimPrefix(url, "https://leetcode.com/problems/"), "/")
}

// recallLabel returns the label of a set task, as taskLabel of src/leetcode/sets.
func recallLabel(task sets.Task) string {
	if rest, ok := strings.CutPrefix(recallSlug(task.URL), "cf-"); ok {
		return "Codeforces " + strings.ToUpper(strings.Replace(rest, "-", "", 1))
	}
	if task.Number != nil && *task.Number != 0 {
		return fmt.Sprintf("LeetCode %d", *task.Number)
	}
	return "LeetCode"
}

// recallTwoPointers returns the slugs of the tasks of the set two-pointers in order, and each task by its slug.
func recallTwoPointers(t *testing.T) ([]string, map[string]sets.Task) {
	t.Helper()
	var slugs []string
	tasks := map[string]sets.Task{}
	for _, set := range goldenSets(t) {
		if set.Slug != "two-pointers" {
			continue
		}
		for _, chapter := range set.Chapters {
			for _, task := range chapter.Tasks {
				slug := recallSlug(task.URL)
				if _, ok := tasks[slug]; !ok {
					slugs = append(slugs, slug)
				}
				tasks[slug] = task
			}
		}
	}
	return slugs, tasks
}

// recallSetup returns an empty database and a server with the learning sets.
func recallSetup(t *testing.T) (*pgxpool.Pool, *Server) {
	t.Helper()
	pool := database(t)
	return pool, newServerWithSets(t, pool, goldenSets(t))
}

// recallAddLibrary adds the problems and their reviews, as addLibrary of test/recall.test.ts, and returns the problem IDs by slug.
func recallAddLibrary(t *testing.T, pool *pgxpool.Pool, library []recallProblem) map[string]string {
	t.Helper()
	ids := map[string]string{}
	for _, problem := range library {
		id := newID()
		ids[problem.slug] = id
		var due, archived any
		if problem.due != "" {
			due = problem.due
		}
		if problem.archived {
			archived = "now"
		}
		exec(t, pool, recallInsertProblem, id, problem.slug, due, archived)
		if problem.reviewedOn != "" {
			exec(t, pool, recallInsertReview, newID(), id, problem.reviewedOn)
		}
	}
	return ids
}

// recallPicksOf returns the recall picks of the day. The test stops on an error.
func recallPicksOf(t *testing.T, pool *pgxpool.Pool, s *Server, today string) []recallPick {
	t.Helper()
	picks, err := s.recallPicks(context.Background(), pool, today)
	if err != nil {
		t.Fatal(err)
	}
	return picks
}

func TestRecallPicks(t *testing.T) {
	order, tasks := recallTwoPointers(t)
	// All set tasks except the first three are due in the manual list, so the set pool has exactly three candidates.
	kept := order[:3]
	var due, keptRows []recallProblem
	for _, slug := range order[3:] {
		due = append(due, recallProblem{slug: slug, due: "2026-10-11"})
	}
	for _, slug := range kept {
		keptRows = append(keptRows, recallProblem{slug: slug})
	}
	tests := []struct {
		name    string
		library []recallProblem
		picks   []string
	}{
		{name: "set tasks without a library row", library: due, picks: kept},
		{name: "set tasks with a library row", library: slices.Concat(keptRows, due), picks: kept},
		{name: "homework slugs outside the library are skipped", library: slices.Concat(due, []recallProblem{{slug: recallXor}}), picks: slices.Concat(kept, []string{recallXor})},
		{name: "homework problem due today leaves the pool", library: slices.Concat(due, []recallProblem{{slug: recallXor}, {slug: recallLetterCase}, {slug: recallPhone, due: "2026-10-11"}}), picks: slices.Concat(kept, []string{recallXor, recallLetterCase})},
		{name: "overdue homework problem leaves the pool", library: slices.Concat(due, []recallProblem{{slug: recallXor}, {slug: recallLetterCase}, {slug: recallPhone, due: "2026-10-01"}}), picks: slices.Concat(kept, []string{recallXor, recallLetterCase})},
		{name: "homework problem due tomorrow stays in the pool", library: slices.Concat(due, []recallProblem{{slug: recallXor}, {slug: recallPhone, due: "2026-10-12"}}), picks: slices.Concat(kept, []string{recallXor, recallPhone})},
		{name: "homework problem reviewed today leaves the pool", library: slices.Concat(due, []recallProblem{{slug: recallXor}, {slug: recallLetterCase}, {slug: recallPhone, reviewedOn: "2026-10-11"}}), picks: slices.Concat(kept, []string{recallXor, recallLetterCase})},
		{name: "homework problem reviewed yesterday stays in the pool", library: slices.Concat(due, []recallProblem{{slug: recallXor}, {slug: recallPhone, reviewedOn: "2026-10-10"}}), picks: slices.Concat(kept, []string{recallXor, recallPhone})},
		{name: "archived homework problem leaves the pool", library: slices.Concat(due, []recallProblem{{slug: recallXor}, {slug: recallLetterCase}, {slug: recallPhone, archived: true}}), picks: slices.Concat(kept, []string{recallXor, recallLetterCase})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := recallSetup(t)
			ids := recallAddLibrary(t, pool, tt.library)
			var want []map[string]any
			for _, slug := range tt.picks {
				task, ok := tasks[slug]
				if !ok {
					want = append(want, map[string]any{"url": "https://leetcode.com/problems/" + slug + "/", "title": "Problem", "difficulty": "medium", "topics": []string{}, "note": "", "slug": slug, "label": "LeetCode", "problemId": ids[slug]})
					continue
				}
				// A pick of a learning set is the whole task with other topics, the slug, the label and the problem of the library.
				pick := map[string]any{}
				data, err := json.Marshal(task)
				if err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(data, &pick); err != nil {
					t.Fatal(err)
				}
				pick["topics"], pick["slug"], pick["label"], pick["problemId"] = recallSetTopics[slug], slug, recallLabel(task), nil
				if id, ok := ids[slug]; ok {
					pick["problemId"] = id
				}
				want = append(want, pick)
			}
			slices.SortFunc(want, func(a, b map[string]any) int { return strings.Compare(a["slug"].(string), b["slug"].(string)) })
			got := recallPicksOf(t, pool, s, "2026-10-11")
			slices.SortFunc(got, func(a, b recallPick) int { return strings.Compare(a.Slug, b.Slug) })
			reviewsEqualJSON(t, got, want)
		})
	}
}

func TestRecallNormalizesTopics(t *testing.T) {
	order, _ := recallTwoPointers(t)
	var due []recallProblem
	for _, slug := range order[3:] {
		due = append(due, recallProblem{slug: slug, due: "2026-10-11"})
	}
	tests := []struct {
		name     string
		slug     string
		topics   []string
		expected []string
	}{
		{name: "library aliases", slug: recallXor, topics: []string{"HashTable", "Хеш-таблица", "Hash Table"}, expected: []string{"Hash Table"}},
		{name: "contextual library topics", slug: recallLetterCase, topics: []string{"Строка", "Перебор", "Backtracking"}, expected: []string{"String", "Backtracking"}},
		{name: "unknown legacy topic", slug: recallPhone, topics: []string{"Особая тема"}, expected: []string{"Особая тема"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := recallSetup(t)
			ids := recallAddLibrary(t, pool, slices.Concat(due, []recallProblem{{slug: tt.slug}}))
			topics, err := json.Marshal(tt.topics)
			if err != nil {
				t.Fatal(err)
			}
			exec(t, pool, "UPDATE problems SET topics=$1,summary='Мои заметки',progress_version=9,understanding='with_help' WHERE slug=$2", string(topics), tt.slug)
			before := homeworkRow(t, pool, "SELECT * FROM problems WHERE slug=$1", tt.slug)
			picks := recallPicksOf(t, pool, s, "2026-10-11")
			index := slices.IndexFunc(picks, func(pick recallPick) bool { return pick.Slug == tt.slug })
			if index < 0 {
				t.Fatalf("picks = %v, want a pick of %s", picks, tt.slug)
			}
			text, err := json.Marshal(picks[index])
			if err != nil {
				t.Fatal(err)
			}
			var pick any
			if err := json.Unmarshal(text, &pick); err != nil {
				t.Fatal(err)
			}
			problemsMatchBody(t, pick, map[string]any{"problemId": ids[tt.slug], "slug": tt.slug, "topics": tt.expected, "note": "Мои заметки"})
			if after := homeworkRow(t, pool, "SELECT * FROM problems WHERE slug=$1", tt.slug); !reflect.DeepEqual(after, before) {
				t.Errorf("problem after the picks = %v, want %v", after, before)
			}
		})
	}
}

func TestRecallPicksOfTheDay(t *testing.T) {
	tests := []struct {
		name  string
		other string
		same  bool
	}{
		{name: "same day gives the same picks", other: "2026-10-11", same: true},
		{name: "next day gives new picks", other: "2026-10-12", same: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := recallSetup(t)
			recallAddLibrary(t, pool, []recallProblem{{slug: recallXor}, {slug: recallLetterCase}, {slug: recallPhone}})
			slugs := func(today string) string {
				var names []string
				for _, pick := range recallPicksOf(t, pool, s, today) {
					names = append(names, pick.Slug)
				}
				slices.Sort(names)
				return strings.Join(names, ",")
			}
			if picks := recallPicksOf(t, pool, s, "2026-10-11"); len(picks) != 5 {
				t.Errorf("%d picks, want 5", len(picks))
			}
			if got, today := slugs(tt.other), slugs("2026-10-11"); (got == today) != tt.same {
				t.Errorf("picks of %s = %s, picks of 2026-10-11 = %s, want the same picks: %v", tt.other, got, today, tt.same)
			}
		})
	}
}

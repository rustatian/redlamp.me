package notebook

import (
	"encoding/json"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The route cases of test/sets.test.ts. The content checks are in content_checks_test.go.

// The prepared statements of test/sets.test.ts.
const (
	setsInsertProblem = "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,'Problem','medium','[]','','','now','now')"
	setsInsertAttempt = "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,$3,'{}','',$4,'with_help','now','now')"
	setsInsertBundle  = "INSERT INTO homework (id,title,instructions,created_at,updated_at) VALUES ($1,$2,'','now','now')"
	setsInsertTask    = "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,$4,'now','now')"
)

// setsSetup returns an empty database, a server with the learning sets of the Worker and the cookie of the student, as beforeEach of test/sets.test.ts.
func setsSetup(t *testing.T) (*pgxpool.Pool, *Server, string) {
	t.Helper()
	pool := database(t)
	s := newServerWithSets(t, pool, goldenSets(t))
	return pool, s, loginAs(t, s, "student")
}

// setsAddProblem adds a problem as addProblem of test/sets.test.ts and returns its ID. state is "" for a problem without attempts.
func setsAddProblem(t *testing.T, pool *pgxpool.Pool, slug, state, acceptance string) string {
	t.Helper()
	id := newID()
	exec(t, pool, setsInsertProblem, id, slug)
	if state != "" {
		exec(t, pool, setsInsertAttempt, newID(), id, state, acceptance)
	}
	return id
}

// setsFind returns the learning set with the slug.
func setsFind(t *testing.T, content []sets.LearningSet, slug string) sets.LearningSet {
	t.Helper()
	i := slices.IndexFunc(content, func(set sets.LearningSet) bool { return set.Slug == slug })
	if i < 0 {
		t.Fatalf("the content has no learning set %q", slug)
	}
	return content[i]
}

// setsJSON returns a value in the form of a decoded JSON body: numbers are float64, slices are []any and structs are maps.
func setsJSON(t *testing.T, value any) any {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	return decoded
}

// setsEqual fails the test unless got and want have the same JSON form, as toEqual of Vitest.
func setsEqual(t *testing.T, what string, got, want any) {
	t.Helper()
	if got, want := setsJSON(t, got), setsJSON(t, want); !reflect.DeepEqual(got, want) {
		gotText, _ := json.Marshal(got)
		wantText, _ := json.Marshal(want)
		t.Errorf("%s = %s, want %s", what, gotText, wantText)
	}
}

// setsDetails is the body of GET /sets/{slug}, with the fields of each chapter as JSON values.
type setsDetails struct {
	Chapters []struct {
		Number  any              `json:"number"`
		Title   any              `json:"title"`
		Summary any              `json:"summary"`
		Links   any              `json:"links"`
		Tasks   []map[string]any `json:"tasks"`
		Quiz    any              `json:"quiz"`
	} `json:"chapters"`
}

// setsLessonBody is { html: chapter.lesson, figures: chapter.figures } of test/sets.test.ts.
// toEqual ignores a field with the value undefined, so a chapter without figures gives no figures field.
func setsLessonBody(chapter sets.Chapter) map[string]any {
	body := map[string]any{"html": chapter.Lesson}
	if chapter.Figures != nil {
		body["figures"] = chapter.Figures
	}
	return body
}

func TestSetsSummary(t *testing.T) {
	type problem struct{ slug, state, acceptance string }
	tests := []struct {
		name          string
		problems      []problem
		acceptedCount int
	}{
		{name: "no library problems", problems: []problem{}, acceptedCount: 0},
		{name: "accepted set task", problems: []problem{{slug: "subsets", state: "saved", acceptance: "accepted"}}, acceptedCount: 1},
		{name: "two accepted set tasks", problems: []problem{{slug: "subsets", state: "saved", acceptance: "accepted"}, {slug: "n-queens", state: "saved", acceptance: "accepted"}}, acceptedCount: 2},
		{name: "accepted problem outside the set", problems: []problem{{slug: "two-sum", state: "saved", acceptance: "accepted"}}, acceptedCount: 0},
		{name: "saved attempt that is not accepted", problems: []problem{{slug: "subsets", state: "saved", acceptance: "not_accepted"}}, acceptedCount: 0},
		{name: "accepted draft", problems: []problem{{slug: "subsets", state: "draft", acceptance: "accepted"}}, acceptedCount: 0},
		{name: "library problem without attempts", problems: []problem{{slug: "subsets", state: "", acceptance: "accepted"}}, acceptedCount: 0},
	}
	content := goldenSets(t)
	set := setsFind(t, content, "permutations-and-combinations")
	slugs := []string{}
	for _, item := range content {
		slugs = append(slugs, item.Slug)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := setsSetup(t)
			for _, p := range tt.problems {
				setsAddProblem(t, pool, p.slug, p.state, p.acceptance)
			}
			w := request(t, s, "/sets", http.MethodGet, nil, student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			items := decode[[]map[string]any](t, w)
			got := []any{}
			var item map[string]any
			for _, value := range items {
				got = append(got, value["slug"])
				if item == nil && value["slug"] == set.Slug {
					item = value
				}
			}
			setsEqual(t, "slugs", got, slugs)
			setsEqual(t, "set summary", item, map[string]any{"slug": set.Slug, "title": set.Title, "summary": set.Summary, "topics": set.Topics, "chapterCount": 4, "taskCount": len(set.Tasks()), "acceptedCount": tt.acceptedCount})
		})
	}
}

// acceptance and taskState are "" for null.
func TestSetsDetails(t *testing.T) {
	tests := []struct {
		name       string
		inLibrary  bool
		acceptance string
		taskState  string
		accepted   bool
		active     bool
		assignable bool
	}{
		{name: "task outside the library", inLibrary: false, acceptance: "", taskState: "", accepted: false, active: false, assignable: true},
		{name: "library task without an accepted attempt", inLibrary: true, acceptance: "not_accepted", taskState: "", accepted: false, active: false, assignable: true},
		{name: "accepted library task", inLibrary: true, acceptance: "accepted", taskState: "", accepted: true, active: false, assignable: false},
		{name: "task with an active homework task", inLibrary: true, acceptance: "", taskState: "changes_requested", accepted: false, active: true, assignable: false},
		{name: "completed homework task is not active", inLibrary: true, acceptance: "accepted", taskState: "completed", accepted: true, active: false, assignable: false},
		{name: "completed homework without acceptance", inLibrary: true, acceptance: "", taskState: "completed", accepted: false, active: false, assignable: false},
		{name: "cancelled unsolved homework", inLibrary: true, acceptance: "", taskState: "cancelled", accepted: false, active: false, assignable: true},
	}
	set := setsFind(t, goldenSets(t), "permutations-and-combinations")
	links := []any{}
	for _, chapter := range set.Chapters {
		links = append(links, chapter.Links)
	}
	i := slices.IndexFunc(set.Chapters[1].Tasks, func(task sets.Task) bool { return task.Slug() == "subsets" })
	if i < 0 {
		t.Fatal("chapter 2 has no task subsets")
	}
	subsets := set.Chapters[1].Tasks[i]
	outline := []any{
		map[string]any{"number": 1, "title": "Считаем варианты", "summary": "Считаем операции и память: Linear Scan, Product Rule, Sum Rule, Factorial, Permutations, Combinations и Subsets. Отличаем число вариантов от размера ответа.", "tasks": []string{"cf-1821-a", "cf-1743-a", "cf-629-a", "cf-1674-b", "number-of-good-pairs"}},
		map[string]any{"number": 2, "title": "Backtracking: Subsets, Permutations, Combinations", "summary": "Backtracking: выбрать, пойти глубже, отменить. start задаёт возрастающие индексы, used отмечает занятые. Считаем вызовы, проверки, копии и память Call Stack.", "tasks": []string{"letter-case-permutation", "subsets", "permutations", "combinations", "letter-combinations-of-a-phone-number"}},
		map[string]any{"number": 3, "title": "Повторы и Pruning", "summary": "Sorting и пропуск убирают повторные ответы Backtracking. Pruning останавливает бесполезные ветки. Считаем вызовы, проверки, копии и память префиксов.", "tasks": []string{"subsets-ii", "permutations-ii", "combination-sum", "generate-parentheses", "cf-1097-b"}},
		map[string]any{"number": 4, "title": "Порядок и номер", "summary": "Next Permutation, Dynamic Programming для путей, Permutation Sequence и N-Queens. Учитываем копии Python, вызовы Factorial и формат ответа.", "tasks": []string{"next-permutation", "unique-paths", "permutation-sequence", "n-queens"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := setsSetup(t)
			other, homeworkID, taskID := setsAddProblem(t, pool, "two-sum", "saved", "accepted"), newID(), newID()
			exec(t, pool, setsInsertBundle, homeworkID, "Week 41: backtracking")
			exec(t, pool, setsInsertTask, newID(), homeworkID, other, "assigned")
			var problemID any
			if tt.inLibrary {
				state := ""
				if tt.acceptance != "" {
					state = "saved"
				}
				problemID = setsAddProblem(t, pool, "subsets", state, tt.acceptance)
			}
			if tt.taskState != "" {
				exec(t, pool, setsInsertTask, taskID, homeworkID, problemID, tt.taskState)
			}
			w := request(t, s, "/sets/"+set.Slug, http.MethodGet, nil, student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"slug": set.Slug, "title": set.Title, "summary": set.Summary, "topics": set.Topics})
			gotLinks, gotOutline, withState := []any{}, []any{}, []any{}
			var row map[string]any
			for _, chapter := range decode[setsDetails](t, w).Chapters {
				gotLinks = append(gotLinks, chapter.Links)
				slugs := []any{}
				for _, task := range chapter.Tasks {
					slugs = append(slugs, task["slug"])
					if row == nil && task["slug"] == "subsets" {
						row = task
					}
					if id, _ := task["problemId"].(string); id != "" || task["accepted"] == true || task["activeTask"] != nil {
						withState = append(withState, task["slug"])
					}
				}
				gotOutline = append(gotOutline, map[string]any{"number": chapter.Number, "title": chapter.Title, "summary": chapter.Summary, "tasks": slugs})
			}
			setsEqual(t, "chapter links", gotLinks, links)
			setsEqual(t, "chapters", gotOutline, outline)
			want, _ := setsJSON(t, subsets).(map[string]any)
			want["slug"], want["label"], want["problemId"], want["accepted"], want["assignable"], want["activeTask"] = "subsets", "LeetCode 78", problemID, tt.accepted, tt.assignable, nil
			if tt.active {
				want["activeTask"] = map[string]any{"state": tt.taskState}
			}
			setsEqual(t, "subsets row", row, want)
			wantState := []string{}
			if tt.inLibrary {
				wantState = []string{"subsets"}
			}
			setsEqual(t, "rows with a problem, an acceptance or an active task", withState, wantState)
		})
	}
}

// best, bestTotal and nextMinutes are nil for null.
func TestSetsChapterQuiz(t *testing.T) {
	type attempt struct {
		deck                            string
		correct, total, passed, minutes int
	}
	tests := []struct {
		name        string
		chapter     int
		cardCount   int
		attempts    []attempt
		best        *int
		bestTotal   *int
		nextMinutes *int
	}{
		{name: "no quiz attempt", chapter: 1, cardCount: 8, attempts: []attempt{}, best: nil, bestTotal: nil, nextMinutes: nil},
		{name: "recent failed quiz attempt", chapter: 1, cardCount: 8, attempts: []attempt{{deck: "permutations-and-combinations-1", correct: 5, total: 8, passed: 0, minutes: 30}}, best: nil, bestTotal: nil, nextMinutes: new(30)},
		{name: "recent passing quiz attempt", chapter: 1, cardCount: 8, attempts: []attempt{{deck: "permutations-and-combinations-1", correct: 7, total: 8, passed: 1, minutes: 30}}, best: new(7), bestTotal: new(8), nextMinutes: new(30)},
		{name: "expired quiz pause", chapter: 1, cardCount: 8, attempts: []attempt{{deck: "permutations-and-combinations-1", correct: 7, total: 8, passed: 1, minutes: 61}}, best: new(7), bestTotal: new(8), nextMinutes: nil},
		{name: "failed quiz attempt after an older pass", chapter: 1, cardCount: 8, attempts: []attempt{{deck: "permutations-and-combinations-1", correct: 7, total: 8, passed: 1, minutes: 180}, {deck: "permutations-and-combinations-1", correct: 5, total: 8, passed: 0, minutes: 30}}, best: new(7), bestTotal: new(8), nextMinutes: new(30)},
		{name: "quiz attempt in another chapter", chapter: 1, cardCount: 8, attempts: []attempt{{deck: "permutations-and-combinations-2", correct: 8, total: 8, passed: 1, minutes: 30}}, best: nil, bestTotal: nil, nextMinutes: nil},
		{name: "quiz attempt in another learning set", chapter: 1, cardCount: 8, attempts: []attempt{{deck: "two-pointers-1", correct: 8, total: 8, passed: 1, minutes: 30}}, best: nil, bestTotal: nil, nextMinutes: nil},
		{name: "passing quiz attempt from a larger deck", chapter: 3, cardCount: 7, attempts: []attempt{{deck: "permutations-and-combinations-3", correct: 8, total: 8, passed: 1, minutes: 120}}, best: new(8), bestTotal: new(8), nextMinutes: nil},
	}
	set := setsFind(t, goldenSets(t), "permutations-and-combinations")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := setsSetup(t)
			now := time.Now().UTC().Truncate(time.Millisecond)
			for _, a := range tt.attempts {
				exec(t, pool, "INSERT INTO quiz_attempts VALUES ($1,$2,$3,$4,$5,$6)", newID(), a.deck, a.correct, a.total, a.passed, now.Add(-time.Duration(a.minutes)*time.Minute).Format("2006-01-02T15:04:05.000Z"))
			}
			w := request(t, s, "/sets/"+set.Slug, http.MethodGet, nil, student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			chapters := decode[setsDetails](t, w).Chapters
			if len(chapters) < tt.chapter {
				t.Fatalf("the set has %d chapters, want chapter %d", len(chapters), tt.chapter)
			}
			var nextAt any
			if tt.nextMinutes != nil {
				nextAt = now.Add(time.Duration(*tt.nextMinutes) * time.Minute).Format("2006-01-02T15:04:05.000Z")
			}
			number := strconv.Itoa(tt.chapter)
			setsEqual(t, "quiz", chapters[tt.chapter-1].Quiz, map[string]any{
				"slug": set.Slug + "-" + number, "title": set.Title + ", глава " + number,
				"chapter":   map[string]any{"set": set.Slug, "setTitle": set.Title, "number": tt.chapter, "title": set.Chapters[tt.chapter-1].Title},
				"cardCount": tt.cardCount, "best": tt.best, "bestTotal": tt.bestTotal,
				"nextAt": nextAt,
			})
		})
	}
}

// body is nil for a missing lesson.
func TestSetsLesson(t *testing.T) {
	content := goldenSets(t)
	set, twoPointers := setsFind(t, content, "permutations-and-combinations"), setsFind(t, content, "two-pointers")
	tests := []struct {
		name    string
		role    string
		slug    string
		chapter string
		body    map[string]any
	}{
		{name: "student chapter 1 lesson", role: "student", slug: set.Slug, chapter: "1", body: setsLessonBody(set.Chapters[0])},
		{name: "parent chapter 4 lesson", role: "parent", slug: set.Slug, chapter: "4", body: setsLessonBody(set.Chapters[3])},
		{name: "lesson with figures", role: "student", slug: twoPointers.Slug, chapter: "1", body: setsLessonBody(twoPointers.Chapters[0])},
		{name: "chapter 5 is missing", role: "parent", slug: set.Slug, chapter: "5", body: nil},
		{name: "chapter 0 is missing", role: "parent", slug: set.Slug, chapter: "0", body: nil},
		{name: "chapter x is missing", role: "parent", slug: set.Slug, chapter: "x", body: nil},
		{name: "lesson without a chapter number is missing", role: "parent", slug: set.Slug, chapter: "", body: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, s, student := setsSetup(t)
			cookie := student
			if tt.role == "parent" {
				cookie = loginAs(t, s, "parent")
			}
			path := "/sets/" + tt.slug + "/lesson"
			if tt.chapter != "" {
				path += "/" + tt.chapter
			}
			w := request(t, s, path, http.MethodGet, nil, cookie)
			if tt.body == nil {
				if w.Code != http.StatusNotFound || errorCode(w) != "not_found" {
					t.Errorf("status = %d, body %s; want 404 and the code not_found", w.Code, w.Body.String())
				}
				return
			}
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			setsEqual(t, "lesson", decode[any](t, w), tt.body)
		})
	}
}

func TestSetsRouteErrors(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		signedIn bool
		status   int
		code     string
	}{
		{name: "unknown set", path: "/sets/missing-set", signedIn: true, status: 404, code: "not_found"},
		{name: "unknown set lesson", path: "/sets/missing-set/lesson/1", signedIn: true, status: 404, code: "not_found"},
		{name: "anonymous set list", path: "/sets", signedIn: false, status: 401, code: "unauthorized"},
		{name: "anonymous set", path: "/sets/permutations-and-combinations", signedIn: false, status: 401, code: "unauthorized"},
		{name: "anonymous set lesson", path: "/sets/permutations-and-combinations/lesson/1", signedIn: false, status: 401, code: "unauthorized"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, s, student := setsSetup(t)
			cookie := ""
			if tt.signedIn {
				cookie = student
			}
			w := request(t, s, tt.path, http.MethodGet, nil, cookie)
			if w.Code != tt.status || errorCode(w) != tt.code {
				t.Errorf("status = %d, body %s; want %d and the code %s", w.Code, w.Body.String(), tt.status, tt.code)
			}
		})
	}
}

// owner is "problems" or "homework". before adds the link of the set before the request.
func TestSetsLinks(t *testing.T) {
	linkedSet := map[string]any{"slug": "permutations-and-combinations", "title": "Permutations и Combinations"}
	tests := []struct {
		name            string
		owner           string
		role            string
		method          string
		slug            string
		before          bool
		version         int
		expected        int
		linked          []map[string]any
		homeworkVersion int64
	}{
		{name: "parent links a set to a problem", owner: "problems", role: "parent", method: http.MethodPut, slug: "permutations-and-combinations", before: false, version: 1, expected: 204, linked: []map[string]any{linkedSet}, homeworkVersion: 1},
		{name: "parent links a set to homework", owner: "homework", role: "parent", method: http.MethodPut, slug: "permutations-and-combinations", before: false, version: 1, expected: 204, linked: []map[string]any{linkedSet}, homeworkVersion: 2},
		{name: "stale homework set link", owner: "homework", role: "parent", method: http.MethodPut, slug: "permutations-and-combinations", before: false, version: 2, expected: 409, linked: []map[string]any{}, homeworkVersion: 1},
		{name: "unknown set", owner: "problems", role: "parent", method: http.MethodPut, slug: "missing-set", before: false, version: 1, expected: 404, linked: []map[string]any{}, homeworkVersion: 1},
		{name: "student cannot link a set", owner: "problems", role: "student", method: http.MethodPut, slug: "permutations-and-combinations", before: false, version: 1, expected: 403, linked: []map[string]any{}, homeworkVersion: 1},
		{name: "parent unlinks a set from a problem", owner: "problems", role: "parent", method: http.MethodDelete, slug: "permutations-and-combinations", before: true, version: 1, expected: 204, linked: []map[string]any{}, homeworkVersion: 1},
		{name: "parent unlinks a set from homework", owner: "homework", role: "parent", method: http.MethodDelete, slug: "permutations-and-combinations", before: true, version: 1, expected: 204, linked: []map[string]any{}, homeworkVersion: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := setsSetup(t)
			parent, problemID, homeworkID := loginAs(t, s, "parent"), setsAddProblem(t, pool, "subsets", "", ""), newID()
			exec(t, pool, setsInsertBundle, homeworkID, "Week 41")
			ownerID, table, body, cookie := problemID, "problem_sets", map[string]any{}, student
			if tt.owner == "homework" {
				ownerID, table, body = homeworkID, "homework_sets", map[string]any{"version": tt.version}
			}
			if tt.before {
				exec(t, pool, "INSERT INTO "+table+" VALUES ($1,$2)", ownerID, tt.slug)
			}
			if tt.role == "parent" {
				cookie = parent
			}
			if w := request(t, s, "/"+tt.owner+"/"+ownerID+"/sets/"+tt.slug, tt.method, body, cookie); w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			linked := decode[map[string]any](t, request(t, s, "/"+tt.owner+"/"+ownerID+"/sets", http.MethodGet, nil, student))
			setsEqual(t, "linked sets", linked["sets"], tt.linked)
			if got := queryValue[int64](t, pool, "SELECT version FROM homework WHERE id=$1", homeworkID); got != tt.homeworkVersion {
				t.Errorf("homework version = %d, want %d", got, tt.homeworkVersion)
			}
		})
	}
}

// owner is "problems" for the first problem of slugs, or "homework" for the homework with a task of each problem.
func TestSetsChapters(t *testing.T) {
	chapter2 := map[string]any{"set": "permutations-and-combinations", "setTitle": "Permutations и Combinations", "number": 2, "title": "Backtracking: Subsets, Permutations, Combinations"}
	chapter3 := map[string]any{"set": "permutations-and-combinations", "setTitle": "Permutations и Combinations", "number": 3, "title": "Повторы и Pruning"}
	tests := []struct {
		name     string
		owner    string
		slugs    []string
		chapters []map[string]any
	}{
		{name: "problem in a chapter", owner: "problems", slugs: []string{"subsets"}, chapters: []map[string]any{chapter2}},
		{name: "problem in no chapter", owner: "problems", slugs: []string{"two-sum"}, chapters: []map[string]any{}},
		{name: "homework with tasks in two chapters", owner: "homework", slugs: []string{"subsets-ii", "two-sum", "subsets", "permutations"}, chapters: []map[string]any{chapter3, chapter2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := setsSetup(t)
			homeworkID, problemIDs := newID(), []string{}
			for _, slug := range tt.slugs {
				problemIDs = append(problemIDs, setsAddProblem(t, pool, slug, "", ""))
			}
			exec(t, pool, setsInsertBundle, homeworkID, "Week 41")
			// Reverse insert order, so only created_at gives the task order.
			for i := len(problemIDs) - 1; i >= 0; i-- {
				exec(t, pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,'assigned',$4,'now')", newID(), homeworkID, problemIDs[i], "2026-10-1"+strconv.Itoa(i))
			}
			ownerID := problemIDs[0]
			if tt.owner == "homework" {
				ownerID = homeworkID
			}
			w := request(t, s, "/"+tt.owner+"/"+ownerID+"/sets", http.MethodGet, nil, student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			setsEqual(t, "body", decode[any](t, w), map[string]any{"chapters": tt.chapters, "sets": []any{}})
		})
	}
}

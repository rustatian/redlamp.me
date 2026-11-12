package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// The cases of test/problems.test.ts.

// problemsSetup returns an empty database, a server on it and the cookie of the student, as beforeEach of test/problems.test.ts.
func problemsSetup(t *testing.T) (*pgxpool.Pool, *Server, string) {
	t.Helper()
	pool := database(t)
	s := newServer(t, pool)
	return pool, s, loginAs(t, s, "student")
}

// problemsInput returns the input of test/problems.test.ts with the changed fields.
func problemsInput(changes map[string]any) map[string]any {
	input := map[string]any{"url": "https://leetcode.com/problems/two-sum/", "number": 1, "title": "Two Sum", "difficulty": "easy", "topics": []string{"Arrays"}, "summary": ""}
	maps.Copy(input, changes)
	return input
}

// problemsPost adds a problem as addProblem does and returns the whole response body.
func problemsPost(t *testing.T, s *Server, cookie string, body any) map[string]any {
	t.Helper()
	w := request(t, s, "/problems", http.MethodPost, body, cookie)
	if w.Code != http.StatusCreated {
		t.Fatalf("add problem: status %d, body %s", w.Code, w.Body.String())
	}
	return decode[map[string]any](t, w)
}

// problemsAddTask adds the homework "Week 41" with one task of the problem, as addTask of test/problems.test.ts.
func problemsAddTask(t *testing.T, pool *pgxpool.Pool, problemID, state string) {
	t.Helper()
	exec(t, pool, "INSERT INTO homework (id,title,instructions,created_at,updated_at) VALUES ('77777777-7777-4777-8777-777777777777','Week 41','Explain','now','now')")
	exec(t, pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,'77777777-7777-4777-8777-777777777777',$2,$3,'now','now')", newID(), problemID, state)
}

// problemsAddAttempt adds an attempt of the problem, as insertAttempt of test/problems.test.ts.
func problemsAddAttempt(t *testing.T, pool *pgxpool.Pool, problemID, state, acceptance string) {
	t.Helper()
	exec(t, pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,$3,'{}','',$4,'independent','now','now')", newID(), problemID, state, acceptance)
}

// problemsTableRows returns the rows of a query as maps from the column name to the value.
func problemsTableRows(t *testing.T, pool *pgxpool.Pool, sql string) []map[string]any {
	t.Helper()
	rows, err := pool.Query(context.Background(), sql)
	if err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	values, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	return values
}

// problemsMatch reports whether got has the fields of want, as toMatchObject of Vitest.
// An object matches when it has each field of want with a matching value. An array matches when it has the same length and each element matches.
func problemsMatch(got, want any) bool {
	switch want := want.(type) {
	case map[string]any:
		object, ok := got.(map[string]any)
		if !ok {
			return false
		}
		for key, value := range want {
			field, ok := object[key]
			if !ok || !problemsMatch(field, value) {
				return false
			}
		}
		return true
	case []any:
		array, ok := got.([]any)
		if !ok || len(array) != len(want) {
			return false
		}
		for i := range want {
			if !problemsMatch(array[i], want[i]) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(got, want)
	}
}

// problemsMatchBody fails the test unless the decoded JSON value got matches want, as toMatchObject of Vitest.
// want goes through JSON first, so its numbers become float64 and its slices become []any.
func problemsMatchBody(t *testing.T, got, want any) {
	t.Helper()
	text, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var expected any
	if err := json.Unmarshal(text, &expected); err != nil {
		t.Fatal(err)
	}
	if !problemsMatch(got, expected) {
		body, _ := json.Marshal(got)
		t.Errorf("body = %s, want the fields of %s", body, text)
	}
}

// problemsListBody is the Page<Problem> of a library read.
type problemsListBody struct {
	Items []struct {
		ID     string   `json:"id"`
		Slug   string   `json:"slug"`
		Topics []string `json:"topics"`
	} `json:"items"`
	NextOffset *int `json:"nextOffset"`
}

func TestProblemsURL(t *testing.T) {
	twoSum := map[string]any{"slug": "two-sum", "url": "https://leetcode.com/problems/two-sum/"}
	theatreSquare := map[string]any{"slug": "cf-1-a", "url": "https://codeforces.com/problemset/problem/1/A"}
	tests := []struct {
		name string
		url  string
		want map[string]any
	}{
		{name: "canonical link", url: "https://leetcode.com/problems/two-sum/", want: twoSum},
		{name: "description link", url: "https://leetcode.com/problems/two-sum/description/?envType=daily-question#notes", want: twoSum},
		{name: "foreign host", url: "https://leetcode.com.evil.example/problems/two-sum/", want: nil},
		{name: "missing slug", url: "https://leetcode.com/problems/", want: nil},
		{name: "insecure URL", url: "http://leetcode.com/problems/two-sum/", want: nil},
		{name: "URL with credentials", url: "https://user@leetcode.com/problems/two-sum/", want: nil},
		{name: "problemset form", url: "https://codeforces.com/problemset/problem/1/A", want: theatreSquare},
		{name: "contest form", url: "https://codeforces.com/contest/1/problem/A", want: theatreSquare},
		{name: "lower case index", url: "https://codeforces.com/problemset/problem/1/a", want: theatreSquare},
		{name: "index with digit", url: "https://codeforces.com/contest/1234/problem/b1", want: map[string]any{"slug": "cf-1234-b1", "url": "https://codeforces.com/problemset/problem/1234/B1"}},
		{name: "codeforces http", url: "http://codeforces.com/problemset/problem/1/A", want: nil},
		{name: "codeforces missing index", url: "https://codeforces.com/problemset/problem/1", want: nil},
		{name: "codeforces extra path", url: "https://codeforces.com/problemset/problem/1/A/?locale=ru#statement", want: theatreSquare},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, s, cookie := problemsSetup(t)
			w := request(t, s, "/problems", http.MethodPost, problemsInput(map[string]any{"url": tt.url}), cookie)
			status := http.StatusCreated
			if tt.want == nil {
				status = http.StatusBadRequest
			}
			if w.Code != status {
				t.Fatalf("status = %d, want %d; body %s", w.Code, status, w.Body.String())
			}
			if tt.want != nil {
				problemsMatchBody(t, decode[any](t, w), tt.want)
			}
		})
	}
}

func TestProblemsDuplicateURL(t *testing.T) {
	pool, s, cookie := problemsSetup(t)
	first := problemsPost(t, s, cookie, problemsInput(nil))
	w := request(t, s, "/problems", http.MethodPost, problemsInput(map[string]any{"url": "https://leetcode.com/problems/two-sum/description/"}), cookie)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if got := decode[map[string]any](t, w); !reflect.DeepEqual(got, first) {
		t.Errorf("body = %v, want the first problem %v", got, first)
	}
	if count := queryValue[int64](t, pool, "SELECT count(*) AS count FROM problems"); count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
}

// The cases "stores canonical topics for $name" of test/problems.test.ts.
func TestProblemsCanonicalTopics(t *testing.T) {
	tests := []struct {
		name   string
		slug   string
		values []string
		want   []string
	}{
		{name: "manual aliases", slug: "two-sum", values: []string{"HashTable", "Хеш-таблица", "Массивы"}, want: []string{"Hash Table", "Array"}},
		{name: "Backtracking context", slug: "subsets", values: []string{"Перебор", "Backtracking"}, want: []string{"Backtracking"}},
		{name: "Dynamic Programming context", slug: "unique-paths", values: []string{"Таблица", "Сочетания"}, want: []string{"Dynamic Programming", "Combinations"}},
		{name: "new English topic", slug: "custom-problem", values: []string{"Custom Technique"}, want: []string{"Custom Technique"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, cookie := problemsSetup(t)
			w := request(t, s, "/problems", http.MethodPost, problemsInput(map[string]any{"url": "https://leetcode.com/problems/" + tt.slug + "/", "topics": tt.values}), cookie)
			if w.Code != http.StatusCreated {
				t.Fatalf("status = %d, want 201; body %s", w.Code, w.Body.String())
			}
			problem := decode[struct {
				ID     string   `json:"id"`
				Topics []string `json:"topics"`
			}](t, w)
			if !slices.Equal(problem.Topics, tt.want) {
				t.Errorf("topics = %q, want %q", problem.Topics, tt.want)
			}
			stored, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatal(err)
			}
			var gotTopics, gotSearch string
			if err := pool.QueryRow(context.Background(), "SELECT topics,search_text FROM problems WHERE id=$1", problem.ID).Scan(&gotTopics, &gotSearch); err != nil {
				t.Fatal(err)
			}
			wantSearch := "two sum 1 " + strings.ToLower(strings.Join(tt.want, " ")) + " "
			if gotTopics != string(stored) || gotSearch != wantSearch {
				t.Errorf("topics = %s, search_text = %q; want %s, %q", gotTopics, gotSearch, stored, wantSearch)
			}
		})
	}
}

func TestProblemsMetadataEditKeepsWorkAndProgress(t *testing.T) {
	ctx := context.Background()
	pool, s, cookie := problemsSetup(t)
	problem := problemsPost(t, s, cookie, problemsInput(nil))
	id, _ := problem["id"].(string)
	const summary = "Моя Хеш-таблица"
	const document = `{"notes":"Перебор и Таблица","approaches":[]}`
	exec(t, pool, `UPDATE problems SET topics='["Хеш-таблица"]',summary=$1,progress_version=7,understanding='with_help',next_review_date='2026-11-01' WHERE id=$2`, summary, id)
	exec(t, pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'saved',$3,'перебор и таблица','accepted','with_help','now','now')", newID(), id, document)
	exec(t, pool, "INSERT INTO reviews (id,problem_id,result,note,reviewed_on,created_at) VALUES ($1,$2,'with_help','Таблица','2026-10-27','now')", newID(), id)
	problemsAddTask(t, pool, id, "completed")
	attempts := problemsTableRows(t, pool, "SELECT * FROM attempts")
	reviews := problemsTableRows(t, pool, "SELECT * FROM reviews")
	tasks := problemsTableRows(t, pool, "SELECT * FROM homework_tasks")
	w := request(t, s, "/problems/"+id, http.MethodPatch, problemsInput(map[string]any{"topics": []string{"HashTable", "Массив", "Хеш-таблица"}, "summary": summary, "version": 1}), cookie)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	problemsMatchBody(t, decode[any](t, w), map[string]any{"id": id, "slug": "two-sum", "number": 1, "version": 2, "progressVersion": 7, "solved": true, "understanding": "with_help", "nextReviewDate": "2026-11-01", "createdAt": problem["createdAt"], "summary": summary, "topics": []string{"Hash Table", "Array"}})
	var gotTopics, gotSearch string
	if err := pool.QueryRow(ctx, "SELECT topics,search_text FROM problems WHERE id=$1", id).Scan(&gotTopics, &gotSearch); err != nil {
		t.Fatal(err)
	}
	if want := `["Hash Table","Array"]`; gotTopics != want {
		t.Errorf("topics = %s, want %s", gotTopics, want)
	}
	if want := "two sum 1 hash table array моя хеш-таблица"; gotSearch != want {
		t.Errorf("search_text = %q, want %q", gotSearch, want)
	}
	if got := problemsTableRows(t, pool, "SELECT * FROM attempts"); !reflect.DeepEqual(got, attempts) {
		t.Errorf("attempts = %v, want %v", got, attempts)
	}
	if got := problemsTableRows(t, pool, "SELECT * FROM reviews"); !reflect.DeepEqual(got, reviews) {
		t.Errorf("reviews = %v, want %v", got, reviews)
	}
	if got := problemsTableRows(t, pool, "SELECT * FROM homework_tasks"); !reflect.DeepEqual(got, tasks) {
		t.Errorf("homework_tasks = %v, want %v", got, tasks)
	}
}

// The cases "normalizes legacy topics in $name" of test/problems.test.ts. A read normalizes the topics, and the stored text stays.
func TestProblemsLegacyTopics(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		version  int
		archived bool
	}{
		{name: "single problem read", action: "read", version: 1, archived: false},
		{name: "library read", action: "list", version: 1, archived: false},
		{name: "duplicate URL", action: "duplicate", version: 1, archived: false},
		{name: "archive response", action: "archive", version: 2, archived: false},
		{name: "restore response", action: "restore", version: 2, archived: true},
	}
	const stored = `["HashTable","Хеш-таблица","Hash Table","Особая тема"]`
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, cookie := problemsSetup(t)
			problem := problemsPost(t, s, cookie, problemsInput(nil))
			id, _ := problem["id"].(string)
			var archivedAt any
			if tt.archived {
				archivedAt = "now"
			}
			exec(t, pool, "UPDATE problems SET topics=$1,summary='Мои заметки',progress_version=8,understanding='independent',next_review_date='2026-11-01',archived_at=$2 WHERE id=$3", stored, archivedAt, id)
			path := "/problems/" + id
			var w *httptest.ResponseRecorder
			switch tt.action {
			case "list":
				w = request(t, s, "/problems", http.MethodGet, nil, cookie)
			case "duplicate":
				w = request(t, s, "/problems", http.MethodPost, problemsInput(nil), cookie)
			case "read":
				w = request(t, s, path, http.MethodGet, nil, cookie)
			default:
				w = request(t, s, path+"/"+tt.action, http.MethodPost, map[string]any{"version": 1}, cookie)
			}
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			body := decode[any](t, w)
			if tt.action == "list" {
				page, _ := body.(map[string]any)
				items, _ := page["items"].([]any)
				if len(items) == 0 {
					t.Fatalf("body = %s, want a problem in items", w.Body.String())
				}
				body = items[0]
			}
			problemsMatchBody(t, body, map[string]any{"id": id, "slug": "two-sum", "version": tt.version, "progressVersion": 8, "understanding": "independent", "nextReviewDate": "2026-11-01", "createdAt": problem["createdAt"], "summary": "Мои заметки", "topics": []string{"Hash Table", "Особая тема"}})
			if got := queryValue[string](t, pool, "SELECT topics FROM problems WHERE id=$1", id); got != stored {
				t.Errorf("topics = %s, want %s", got, stored)
			}
		})
	}
}

func TestProblemsRejectNonEnglishTopics(t *testing.T) {
	tests := []struct {
		name   string
		method string
		query  string
		values []string
	}{
		{name: "POST rejects a new non-English topic", method: http.MethodPost, query: "", values: []string{"Особая тема"}},
		{name: "PATCH rejects a new non-English topic", method: http.MethodPatch, query: "", values: []string{"Особая тема"}},
		{name: "topic filter needs unambiguous context", method: http.MethodGet, query: "?" + url.Values{"topic": {"Перебор"}}.Encode(), values: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, s, cookie := problemsSetup(t)
			id := addProblem(t, s, cookie, problemsInput(nil))
			path := "/problems" + tt.query
			if tt.method == http.MethodPatch {
				path = "/problems/" + id
			}
			var body any
			if tt.method != http.MethodGet {
				changes := map[string]any{"topics": tt.values}
				if tt.method == http.MethodPatch {
					changes["version"] = 1
				}
				body = problemsInput(changes)
			}
			w := request(t, s, path, tt.method, body, cookie)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", w.Code)
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"error": map[string]any{"code": "invalid_input", "message": "Укажите название темы на английском языке."}})
		})
	}
}

// The cases "filters $name" of test/problems.test.ts.
func TestProblemsFilters(t *testing.T) {
	tests := []struct {
		name       string
		query      url.Values
		slugs      []string
		nextOffset *int
	}{
		{name: "canonical topic on the first page", query: url.Values{"topic": {"Hash Table"}, "limit": {"1"}}, slugs: []string{"legacy-first"}, nextOffset: new(1)},
		{name: "Russian alias on the second page", query: url.Values{"topic": {"Хеш-таблица"}, "limit": {"1"}, "offset": {"1"}}, slugs: []string{"mixed"}, nextOffset: new(2)},
		{name: "compact alias on the last page", query: url.Values{"topic": {"HashTable"}, "limit": {"2"}, "offset": {"2"}}, slugs: []string{"canonical", "legacy-last"}, nextOffset: nil},
		{name: "topic case and whitespace", query: url.Values{"topic": {" hash  TABLE "}}, slugs: []string{"legacy-first", "mixed", "canonical", "legacy-last"}, nextOffset: nil},
		{name: "canonical text finds legacy metadata", query: url.Values{"q": {"HASH TABLE"}, "limit": {"1"}}, slugs: []string{"legacy-first"}, nextOffset: new(1)},
		{name: "canonical text on the last page", query: url.Values{"q": {"Hash Table"}, "limit": {"1"}, "offset": {"3"}}, slugs: []string{"legacy-last"}, nextOffset: nil},
		{name: "raw metadata text stays searchable", query: url.Values{"q": {"HashTable"}}, slugs: []string{"legacy-first", "mixed"}, nextOffset: nil},
		{name: "topic and scalar filters before pagination", query: url.Values{"topic": {"HashTable"}, "difficulty": {"medium"}, "understanding": {"independent"}, "due": {"2026-10-27"}, "limit": {"1"}, "offset": {"1"}}, slugs: []string{"legacy-last"}, nextOffset: nil},
		{name: "literal attempt text and topic filter", query: url.Values{"topic": {"Hash Table"}, "q": {"МОЯ ХЕШ-ТАБЛИЦА 100%_READY"}}, slugs: []string{"legacy-last"}, nextOffset: nil},
		{name: "canonical metadata respects topic filter", query: url.Values{"topic": {"Array"}, "q": {"Hash Table"}}, slugs: []string{}, nextOffset: nil},
		{name: "archived canonical text", query: url.Values{"archived": {"1"}, "q": {"Hash Table"}}, slugs: []string{"archived"}, nextOffset: nil},
	}
	rows := []struct {
		name          string
		slug          string
		topics        []string
		date          string
		difficulty    string
		understanding string
		archived      any
	}{
		{name: "last legacy match", slug: "legacy-last", topics: []string{"Хеш-таблица"}, date: "2026-10-01", difficulty: "medium", understanding: "independent", archived: nil},
		{name: "canonical match", slug: "canonical", topics: []string{"Hash Table"}, date: "2026-10-02", difficulty: "easy", understanding: "independent", archived: nil},
		{name: "several aliases on one problem", slug: "mixed", topics: []string{"HashTable", "Хеш-таблица", "Hash Table"}, date: "2026-10-03", difficulty: "medium", understanding: "with_help", archived: nil},
		{name: "first legacy match", slug: "legacy-first", topics: []string{"HashTable"}, date: "2026-10-04", difficulty: "medium", understanding: "independent", archived: nil},
		{name: "newer unrelated problem", slug: "unrelated", topics: []string{"Array"}, date: "2026-10-05", difficulty: "medium", understanding: "independent", archived: nil},
		{name: "archived legacy match", slug: "archived", topics: []string{"Хеш-таблица"}, date: "2026-10-06", difficulty: "medium", understanding: "independent", archived: "now"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, cookie := problemsSetup(t)
			for _, row := range rows {
				id := newID()
				stored, err := json.Marshal(row.topics)
				if err != nil {
					t.Fatal(err)
				}
				exec(t, pool, "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,understanding,next_review_date,archived_at,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,'',$6,$7,'2026-10-20',$8,'now',$9)", id, row.slug, row.name, row.difficulty, string(stored), strings.ToLower(strings.Join(row.topics, " ")), row.understanding, row.archived, row.date)
				if row.slug == "legacy-last" || row.slug == "unrelated" {
					search := "моя хеш-таблица 100xready"
					if row.slug == "legacy-last" {
						search = "моя хеш-таблица 100%_ready"
					}
					exec(t, pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'draft','{}',$3,'not_submitted','needs_practice','now','now')", newID(), id, search)
				}
			}
			w := request(t, s, "/problems?"+tt.query.Encode(), http.MethodGet, nil, cookie)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			page := decode[problemsListBody](t, w)
			slugs := []string{}
			for _, item := range page.Items {
				slugs = append(slugs, item.Slug)
				if !slices.Equal(item.Topics, []string{"Hash Table"}) {
					t.Errorf("topics of %s = %q, want [Hash Table]", item.Slug, item.Topics)
				}
			}
			if !slices.Equal(slugs, tt.slugs) {
				t.Errorf("slugs = %q, want %q", slugs, tt.slugs)
			}
			got, _ := json.Marshal(page.NextOffset)
			want, _ := json.Marshal(tt.nextOffset)
			if string(got) != string(want) {
				t.Errorf("nextOffset = %s, want %s", got, want)
			}
		})
	}
}

func TestProblemsUnicodeNotesAsLiteralText(t *testing.T) {
	pool, s, cookie := problemsSetup(t)
	first, second := newID(), newID()
	const problem = "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,'Search fixture','medium','[]','','','now','now')"
	const attempt = "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'draft','{}',$3,'not_submitted','needs_practice','now','now')"
	exec(t, pool, problem, first, "one")
	exec(t, pool, problem, second, "two")
	exec(t, pool, attempt, newID(), first, "инвариант 100%_ready")
	exec(t, pool, attempt, newID(), second, "инвариант 100xready")
	w := request(t, s, "/problems?"+url.Values{"q": {"ИНВАРИАНТ 100%_ready"}}.Encode(), http.MethodGet, nil, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	ids := []string{}
	for _, item := range decode[problemsListBody](t, w).Items {
		ids = append(ids, item.ID)
	}
	if !slices.Equal(ids, []string{first}) {
		t.Errorf("ids = %q, want [%s]", ids, first)
	}
}

func TestProblemsEditVersion(t *testing.T) {
	tests := []struct {
		name   string
		stored int
		status int
		title  string
	}{
		{name: "current edit", stored: 1, status: 200, title: "Updated"},
		{name: "stale edit", stored: 2, status: 409, title: "Two Sum"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, cookie := problemsSetup(t)
			id := addProblem(t, s, cookie, problemsInput(nil))
			exec(t, pool, "UPDATE problems SET version=$1", tt.stored)
			if w := request(t, s, "/problems/"+id, http.MethodPatch, problemsInput(map[string]any{"version": 1, "title": "Updated"}), cookie); w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, request(t, s, "/problems/"+id, http.MethodGet, nil, cookie)), map[string]any{"title": tt.title, "version": 2})
		})
	}
}

// The cases "homework picker filters $name" of test/problems.test.ts. An empty task or attempt means none.
func TestProblemsHomeworkPicker(t *testing.T) {
	tests := []struct {
		name      string
		task      string
		attempt   string
		available bool
	}{
		{name: "unassigned problem", task: "", attempt: "", available: true},
		{name: "assigned problem", task: "assigned", attempt: "", available: false},
		{name: "problem in progress", task: "in_progress", attempt: "", available: false},
		{name: "submitted problem", task: "submitted", attempt: "", available: false},
		{name: "problem returned for changes", task: "changes_requested", attempt: "", available: false},
		{name: "completed problem without acceptance", task: "completed", attempt: "", available: false},
		{name: "cancelled unsolved problem", task: "cancelled", attempt: "", available: true},
		{name: "accepted problem without homework", task: "", attempt: "saved", available: false},
		{name: "accepted problem with cancelled homework", task: "cancelled", attempt: "saved", available: false},
		{name: "draft acceptance is not a solved problem", task: "", attempt: "draft", available: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, cookie := problemsSetup(t)
			id := addProblem(t, s, cookie, problemsInput(nil))
			if tt.task != "" {
				problemsAddTask(t, pool, id, tt.task)
			}
			if tt.attempt != "" {
				problemsAddAttempt(t, pool, id, tt.attempt, "accepted")
			}
			w := request(t, s, "/problems?assignable=1&q=Two", http.MethodGet, nil, cookie)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			ids := []string{}
			for _, item := range decode[problemsListBody](t, w).Items {
				ids = append(ids, item.ID)
			}
			want := []string{}
			if tt.available {
				want = []string{id}
			}
			if !slices.Equal(ids, want) {
				t.Errorf("assignable ids = %q, want %q", ids, want)
			}
			if items := decode[problemsListBody](t, request(t, s, "/problems", http.MethodGet, nil, cookie)).Items; len(items) != 1 {
				t.Errorf("library has %d problems, want 1", len(items))
			}
		})
	}
}

// The URL change cases of test/problems.test.ts. An empty attempt or task means none.
func TestProblemsURLChange(t *testing.T) {
	const moveZeroes = "https://leetcode.com/problems/move-zeroes/"
	tests := []struct {
		name    string
		attempt string
		task    string
		url     string
		status  int
		slug    string
	}{
		{name: "URL change of a problem without work", attempt: "", task: "", url: moveZeroes, status: 200, slug: "move-zeroes"},
		{name: "URL change of a problem with a draft", attempt: "draft", task: "", url: moveZeroes, status: 200, slug: "move-zeroes"},
		{name: "URL change of a problem with an accepted attempt", attempt: "accepted", task: "", url: moveZeroes, status: 409, slug: "two-sum"},
		{name: "URL change of a problem with an assigned task", attempt: "", task: "assigned", url: moveZeroes, status: 409, slug: "two-sum"},
		{name: "URL change of a problem with a completed task", attempt: "", task: "completed", url: moveZeroes, status: 409, slug: "two-sum"},
		{name: "URL change of a problem with a cancelled task", attempt: "", task: "cancelled", url: moveZeroes, status: 200, slug: "move-zeroes"},
		{name: "other form of the same URL of a problem with an accepted attempt", attempt: "accepted", task: "", url: "https://leetcode.com/problems/two-sum/description/", status: 200, slug: "two-sum"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, cookie := problemsSetup(t)
			id := addProblem(t, s, cookie, problemsInput(nil))
			switch tt.attempt {
			case "draft":
				problemsAddAttempt(t, pool, id, "draft", "not_submitted")
			case "accepted":
				problemsAddAttempt(t, pool, id, "saved", "accepted")
			}
			if tt.task != "" {
				problemsAddTask(t, pool, id, tt.task)
			}
			w := request(t, s, "/problems/"+id, http.MethodPatch, problemsInput(map[string]any{"url": tt.url, "version": 1, "title": "Updated"}), cookie)
			if w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			if tt.status == http.StatusConflict {
				problemsMatchBody(t, decode[any](t, w), map[string]any{"error": map[string]any{"code": "slug_locked"}})
			}
			title := "Two Sum"
			if tt.status == http.StatusOK {
				title = "Updated"
			}
			problemsMatchBody(t, decode[any](t, request(t, s, "/problems/"+id, http.MethodGet, nil, cookie)), map[string]any{"slug": tt.slug, "title": title})
		})
	}
}

func TestProblemsArchiveAfterActiveHomeworkCompletes(t *testing.T) {
	pool, s, cookie := problemsSetup(t)
	id := addProblem(t, s, cookie, problemsInput(nil))
	path := "/problems/" + id
	problemsAddTask(t, pool, id, "assigned")
	if w := request(t, s, path+"/archive", http.MethodPost, map[string]any{"version": 1}, cookie); w.Code != http.StatusConflict {
		t.Errorf("archive with an assigned task: status = %d, want 409", w.Code)
	}
	exec(t, pool, "UPDATE homework_tasks SET state='completed'")
	if w := request(t, s, path+"/archive", http.MethodPost, map[string]any{"version": 1}, cookie); w.Code != http.StatusOK {
		t.Errorf("archive with a completed task: status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	problemsMatchBody(t, decode[any](t, request(t, s, "/problems", http.MethodGet, nil, cookie)), map[string]any{"items": []any{}})
	problemsMatchBody(t, decode[any](t, request(t, s, path, http.MethodGet, nil, cookie)), map[string]any{"title": "Two Sum"})
}

func TestProblemsArchiveConflicts(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		archivedAt any
		homework   string
		stored     int64
		wantError  map[string]any
	}{
		{name: "archive with active homework", action: "archive", archivedAt: nil, homework: "assigned", stored: 1, wantError: map[string]any{"code": "active_homework", "message": "Завершите или отмените активное задание, прежде чем отправлять задачу в архив."}},
		{name: "stale archive", action: "archive", archivedAt: nil, homework: "completed", stored: 2, wantError: map[string]any{"code": "conflict"}},
		{name: "stale restore", action: "restore", archivedAt: "now", homework: "completed", stored: 2, wantError: map[string]any{"code": "conflict"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, cookie := problemsSetup(t)
			id := addProblem(t, s, cookie, problemsInput(nil))
			exec(t, pool, "UPDATE problems SET archived_at=$1,version=$2", tt.archivedAt, tt.stored)
			problemsAddTask(t, pool, id, tt.homework)
			w := request(t, s, "/problems/"+id+"/"+tt.action, http.MethodPost, map[string]any{"version": 1}, cookie)
			if w.Code != http.StatusConflict {
				t.Errorf("status = %d, want 409", w.Code)
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"error": tt.wantError})
			var archivedAt any
			var version int64
			if err := pool.QueryRow(context.Background(), "SELECT archived_at AS archivedAt,version FROM problems").Scan(&archivedAt, &version); err != nil {
				t.Fatal(err)
			}
			if archivedAt != tt.archivedAt || version != tt.stored {
				t.Errorf("archived_at = %v, version = %d; want %v, %d", archivedAt, version, tt.archivedAt, tt.stored)
			}
		})
	}
}

// The pagination and filter validation cases of test/problems.test.ts. A nil body means no body check.
func TestProblemsPages(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		status int
		body   map[string]any
	}{
		{name: "page with more results", query: "limit=1", status: 200, body: map[string]any{"nextOffset": 1}},
		{name: "last page", query: "limit=1&offset=1", status: 200, body: map[string]any{"nextOffset": nil}},
		{name: "combined filters without a match", query: "difficulty=hard&topic=Arrays", status: 200, body: map[string]any{"items": []any{}}},
		{name: "page size above the limit", query: "limit=51", status: 400, body: nil},
		{name: "impossible due date", query: "due=2026-02-30", status: 400, body: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, s, cookie := problemsSetup(t)
			addProblem(t, s, cookie, problemsInput(nil))
			addProblem(t, s, cookie, problemsInput(map[string]any{"url": "https://leetcode.com/problems/three-sum/", "difficulty": "medium"}))
			w := request(t, s, "/problems?"+tt.query, http.MethodGet, nil, cookie)
			if w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			if tt.body != nil {
				problemsMatchBody(t, decode[any](t, w), tt.body)
			}
		})
	}
}

// The problem details cases of test/problems.test.ts. The fake LeetCode takes the place of the fetch mock.
// An empty slug means no LeetCode request.
func TestProblemsDetails(t *testing.T) {
	const question = `{"questionFrontendId":"1","title":"Two Sum","difficulty":"Easy","topicTags":[{"name":"Array"},{"name":"Hash Table"}]}`
	answer := func(value string) func(string, string) (json.RawMessage, error) {
		return func(string, string) (json.RawMessage, error) { return json.RawMessage(value), nil }
	}
	fail := func(message string) func(string, string) (json.RawMessage, error) {
		return func(string, string) (json.RawMessage, error) { return nil, errors.New(message) }
	}
	tests := []struct {
		name   string
		url    string
		reply  func(titleSlug, fields string) (json.RawMessage, error)
		status int
		body   map[string]any
		slug   string
	}{
		{name: "details from a description link", url: "https://leetcode.com/problems/two-sum/description/", reply: answer(question), status: 200, body: map[string]any{"number": 1, "title": "Two Sum", "difficulty": "easy", "topics": []string{"Array", "Hash Table"}}, slug: "two-sum"},
		{name: "details without a numeric problem ID", url: "https://leetcode.com/problems/two-sum/", reply: answer(`{"questionFrontendId":"LCP 01","title":"Two Sum","difficulty":"Easy","topicTags":[{"name":"Array"},{"name":"Hash Table"}]}`), status: 200, body: map[string]any{"number": nil, "title": "Two Sum"}, slug: "two-sum"},
		{name: "details normalize aliases and retain new English topics", url: "https://leetcode.com/problems/two-sum/", reply: answer(`{"questionFrontendId":"1","title":"Two Sum","difficulty":"Easy","topicTags":[{"name":"HashTable"},{"name":"Хеш-таблица"},{"name":"Arrays"},{"name":"Custom Technique"}]}`), status: 200, body: map[string]any{"topics": []string{"Hash Table", "Array", "Custom Technique"}}, slug: "two-sum"},
		{name: "unknown LeetCode problem", url: "https://leetcode.com/problems/two-sum/", reply: func(string, string) (json.RawMessage, error) { return nil, nil }, status: 404, body: map[string]any{"error": map[string]any{"code": "not_found"}}, slug: "two-sum"},
		{name: "LeetCode error response", url: "https://leetcode.com/problems/two-sum/", reply: fail("LeetCode returned 503"), status: 502, body: map[string]any{"error": map[string]any{"code": "lookup_failed"}}, slug: "two-sum"},
		{name: "unreachable LeetCode", url: "https://leetcode.com/problems/two-sum/", reply: fail("Network connection lost."), status: 502, body: map[string]any{"error": map[string]any{"code": "lookup_failed"}}, slug: "two-sum"},
		{name: "details for a URL outside LeetCode", url: "https://example.com/problems/two-sum/", reply: nil, status: 400, body: map[string]any{"error": map[string]any{"code": "invalid_url"}}, slug: ""},
		{name: "details for a Codeforces URL", url: "https://codeforces.com/problemset/problem/1/A", reply: nil, status: 400, body: map[string]any{"error": map[string]any{"code": "lookup_unavailable"}}, slug: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeLeetCode{reply: tt.reply}
			s := newServerWithLeetCode(t, database(t), fake)
			cookie := loginAs(t, s, "student")
			w := request(t, s, "/problem-details?url="+url.QueryEscape(tt.url), http.MethodGet, nil, cookie)
			if w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), tt.body)
			slugs := []string{}
			for _, c := range fake.calls {
				slugs = append(slugs, c.titleSlug)
			}
			want := []string{}
			if tt.slug != "" {
				want = []string{tt.slug}
			}
			if !slices.Equal(slugs, want) {
				t.Errorf("LeetCode requests for %q, want %q", slugs, want)
			}
		})
	}
}

// The request and the answers of question() in src/leetcode/problems.ts.
// The details cases of test/problems.test.ts send these answers through the fetch mock.
func TestProblemsDetailsLeetCodeRequest(t *testing.T) {
	const fields = "questionFrontendId title difficulty topicTags { name }"
	const question = `{"questionFrontendId":"1","title":"Two Sum","difficulty":"Easy","topicTags":[{"name":"Array"},{"name":"Hash Table"}]}`
	tests := []struct {
		name    string
		status  int
		reply   string
		want    string
		wantErr bool
	}{
		{name: "details from a description link", status: 200, reply: `{"data":{"question":` + question + `}}`, want: question},
		{name: "unknown LeetCode problem", status: 200, reply: `{"data":{"question":null}}`, want: ""},
		{name: "LeetCode error response", status: 503, reply: "Busy", wantErr: true},
	}
	wantBody := map[string]any{
		"query":     "query question($titleSlug: String!) { question(titleSlug: $titleSlug) { " + fields + " } }",
		"variables": map[string]any{"titleSlug": "two-sum"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var method, contentType, referer string
			var body map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method, contentType, referer = r.Method, r.Header.Get("Content-Type"), r.Header.Get("Referer")
				_ = json.NewDecoder(r.Body).Decode(&body)
				w.WriteHeader(tt.status)
				_, _ = io.WriteString(w, tt.reply)
			}))
			defer server.Close()
			got, err := GraphQL{URL: server.URL}.Question(context.Background(), "two-sum", fields)
			if (err != nil) != tt.wantErr || string(got) != tt.want {
				t.Errorf("Question = %s, %v; want %q and an error: %v", got, err, tt.want, tt.wantErr)
			}
			if method != http.MethodPost || contentType != "application/json" || referer != "https://leetcode.com/problems/two-sum/" {
				t.Errorf("request = %s, Content-Type %q, Referer %q; want POST, application/json, https://leetcode.com/problems/two-sum/", method, contentType, referer)
			}
			if !reflect.DeepEqual(body, wantBody) {
				t.Errorf("request body = %v, want %v", body, wantBody)
			}
		})
	}
}

func TestProblemsRestoreToLibraryAndHomework(t *testing.T) {
	_, s, cookie := problemsSetup(t)
	id := addProblem(t, s, cookie, problemsInput(nil))
	parent := loginAs(t, s, "parent")
	if w := request(t, s, "/problems/"+id+"/archive", http.MethodPost, map[string]any{"version": 1}, cookie); w.Code != http.StatusOK {
		t.Fatalf("archive: status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	restored := request(t, s, "/problems/"+id+"/restore", http.MethodPost, map[string]any{"version": 2}, parent)
	if restored.Code != http.StatusOK {
		t.Errorf("restore: status = %d, want 200; body %s", restored.Code, restored.Body.String())
	}
	problemsMatchBody(t, decode[any](t, restored), map[string]any{"id": id, "archivedAt": nil, "version": 3})
	problemsMatchBody(t, decode[any](t, request(t, s, "/problems", http.MethodGet, nil, cookie)), map[string]any{"items": []any{map[string]any{"id": id}}})
	homework := decode[struct {
		ID string `json:"id"`
	}](t, request(t, s, "/homework", http.MethodPost, map[string]any{"title": "Week 41", "instructions": "", "dueDate": nil}, parent))
	if w := request(t, s, "/tasks", http.MethodPost, map[string]any{"homeworkId": homework.ID, "problemId": id}, parent); w.Code != http.StatusCreated {
		t.Errorf("task: status = %d, want 201; body %s", w.Code, w.Body.String())
	}
}

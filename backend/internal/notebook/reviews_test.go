package notebook

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The cases of test/reviews.test.ts.

// The prepared statements of test/reviews.test.ts.
const (
	reviewsInsertAttempt = "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,$3,'{}','',$4,$5,'now','now')"
	reviewsSelectDates   = "SELECT next_review_date,progress_version FROM problems WHERE id=$1"
	reviewsCount         = "SELECT count(*) AS count FROM reviews"
)

// reviewsFixture is the state that beforeEach of test/reviews.test.ts makes.
type reviewsFixture struct {
	pool                              *pgxpool.Pool
	s                                 *Server
	cookie, first, taskID, homeworkID string
}

// as returns the student cookie of beforeEach, or a new parent session for the role "parent".
func (f reviewsFixture) as(t *testing.T, role string) string {
	t.Helper()
	if role == "parent" {
		return loginAs(t, f.s, "parent")
	}
	return f.cookie
}

// reviewsSetup adds three problems with one saved attempt each, due on 2026-10-07, and the homework "Week 41" with an assigned task of the first problem.
func reviewsSetup(t *testing.T) reviewsFixture {
	t.Helper()
	pool := database(t)
	s := newServerWithSets(t, pool, goldenSets(t))
	f := reviewsFixture{pool: pool, s: s, cookie: loginAs(t, s, "student")}
	rows := []struct{ name, slug, topics, acceptance, understanding string }{
		{name: "accepted independently", slug: "one", topics: `["Arrays"]`, acceptance: "accepted", understanding: "independent"},
		{name: "accepted with help", slug: "two", topics: `["Arrays","Graphs"]`, acceptance: "accepted", understanding: "with_help"},
		{name: "unaccepted independently", slug: "three", topics: `["Graphs"]`, acceptance: "not_accepted", understanding: "independent"},
	}
	for _, row := range rows {
		id := newID()
		if row.slug == "one" {
			f.first = id
		}
		exec(t, pool, "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,understanding,next_review_date,created_at,updated_at) VALUES ($1,$2,$3,'medium',$4,'','',$5,'2026-10-07','now','now')", id, row.slug, row.name, row.topics, row.understanding)
		exec(t, pool, reviewsInsertAttempt, newID(), id, "saved", row.acceptance, row.understanding)
	}
	f.taskID, f.homeworkID = newID(), newID()
	exec(t, pool, "INSERT INTO homework (id,title,instructions,due_date,created_at,updated_at) VALUES ($1,'Week 41','Explain','2026-10-07','now','now')", f.homeworkID)
	exec(t, pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,'assigned','now','now')", f.taskID, f.homeworkID, f.first)
	return f
}

// reviewsDashboard returns the body of GET /dashboard?today=<today>. The test stops unless the status is 200.
func reviewsDashboard(t *testing.T, f reviewsFixture, today, cookie string) map[string]any {
	t.Helper()
	w := request(t, f.s, "/dashboard?today="+today, http.MethodGet, nil, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("dashboard: status %d, body %s", w.Code, w.Body.String())
	}
	return decode[map[string]any](t, w)
}

// reviewsItems returns the elements of a JSON array of objects. The test stops when value is not such an array.
func reviewsItems(t *testing.T, value any) []map[string]any {
	t.Helper()
	array, ok := value.([]any)
	if !ok {
		t.Fatalf("value = %v, want a JSON array", value)
	}
	items := make([]map[string]any, len(array))
	for i, item := range array {
		if items[i], ok = item.(map[string]any); !ok {
			t.Fatalf("element %d = %v, want a JSON object", i, item)
		}
	}
	return items
}

// reviewsEqualJSON fails the test unless got and want have the same JSON value, as toEqual of Vitest.
func reviewsEqualJSON(t *testing.T, got, want any) {
	t.Helper()
	normal := func(value any) (any, string) {
		text, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var decoded any
		if err := json.Unmarshal(text, &decoded); err != nil {
			t.Fatal(err)
		}
		return decoded, string(text)
	}
	gotValue, gotText := normal(got)
	wantValue, wantText := normal(want)
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("value = %s, want %s", gotText, wantText)
	}
}

// reviewsNull returns nil when set is false. A nil argument becomes NULL.
func reviewsNull(set bool, value any) any {
	if !set {
		return nil
	}
	return value
}

func TestReviewsDashboardContentHomework(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		state   string
		visible bool
		href    string
	}{
		{name: "assigned chapter", kind: "chapter", state: "assigned", visible: true, href: "/leetcode/sets/permutations-and-combinations/1"},
		{name: "quiz in progress", kind: "quiz", state: "in_progress", visible: true, href: "/leetcode/quiz/permutations-and-combinations-1"},
		{name: "completed chapter", kind: "chapter", state: "completed", visible: false, href: "/leetcode/sets/permutations-and-combinations/1"},
		{name: "cancelled quiz", kind: "quiz", state: "cancelled", visible: false, href: "/leetcode/quiz/permutations-and-combinations-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			id := newID()
			exec(t, f.pool, "INSERT INTO homework_tasks (id,homework_id,kind,set_slug,chapter,state,created_at,updated_at,completed_at) VALUES ($1,$2,$3,'permutations-and-combinations',1,$4,'now','now',$5)", id, f.homeworkID, tt.kind, tt.state, reviewsNull(tt.state == "completed", "2026-10-07"))
			dashboard := reviewsDashboard(t, f, "2026-10-07", f.cookie)
			problemsMatchBody(t, dashboard["counts"], map[string]any{"recorded": 3, "solved": 2, "waitingReview": 0})
			homework := reviewsItems(t, dashboard["homework"])
			var mine []map[string]any
			for _, task := range homework {
				if task["id"] == id {
					mine = append(mine, task)
				}
			}
			switch {
			case !tt.visible && len(mine) != 0:
				t.Errorf("content task = %v, want no task", mine)
			case tt.visible && len(mine) != 1:
				t.Errorf("content tasks = %v, want one task", mine)
			case tt.visible:
				title, _ := mine[0]["title"].(string)
				if !problemsMatch(mine[0], map[string]any{"kind": tt.kind, "state": tt.state, "href": tt.href, "unread": false}) || !strings.Contains(title, "Считаем варианты") {
					t.Errorf("content task = %v, want kind %s, state %s, href %s, unread false and a title with «Считаем варианты»", mine[0], tt.kind, tt.state, tt.href)
				}
			}
			want := map[string]any{"id": f.taskID, "kind": "problem", "title": "accepted independently", "href": "/leetcode/problems/" + f.first}
			if !slices.ContainsFunc(homework, func(task map[string]any) bool { return problemsMatch(task, want) }) {
				t.Errorf("homework = %v, want a task with the fields %v", homework, want)
			}
		})
	}
}

func TestReviewsDashboardTenTaskLimit(t *testing.T) {
	tests := []struct {
		name   string
		role   string
		order  []string
		unread []string
	}{
		{name: "parent sees submissions first, then the oldest updates", role: "parent", order: []string{"01", "02", "11", "10", "09", "08", "07", "06", "05", "04"}, unread: []string{"11"}},
		{name: "student sees work first, then submissions, by due date", role: "student", order: []string{"03", "04", "05", "07", "08", "09", "10", "06", "11", "02"}, unread: []string{"05"}},
	}
	const prefix = "00000000-0000-4000-8000-0000000000"
	// chapter and due are nil for NULL.
	rows := []struct {
		name, id, kind, state string
		chapter, due          any
		updated               string
		archived              bool
	}{
		{name: "older submission", id: "01", kind: "problem", state: "submitted", chapter: nil, due: "2026-10-20", updated: "2026-10-01", archived: false},
		{name: "newer submission", id: "02", kind: "problem", state: "submitted", chapter: nil, due: "2026-10-01", updated: "2026-10-02", archived: false},
		{name: "reading with the first due date", id: "03", kind: "chapter", state: "assigned", chapter: 1, due: "2026-10-01", updated: "2026-10-09", archived: false},
		{name: "quiz with the same due date", id: "04", kind: "quiz", state: "in_progress", chapter: 1, due: "2026-10-01", updated: "2026-10-08", archived: false},
		{name: "problem returned for changes", id: "05", kind: "problem", state: "changes_requested", chapter: nil, due: "2026-10-02", updated: "2026-10-07", archived: false},
		{name: "reading without a due date", id: "06", kind: "chapter", state: "assigned", chapter: 2, due: nil, updated: "2026-10-06", archived: false},
		{name: "later quiz", id: "07", kind: "quiz", state: "assigned", chapter: 2, due: "2026-10-03", updated: "2026-10-05", archived: false},
		{name: "problem in progress", id: "08", kind: "problem", state: "in_progress", chapter: nil, due: "2026-10-04", updated: "2026-10-04", archived: false},
		{name: "later reading", id: "09", kind: "chapter", state: "assigned", chapter: 3, due: "2026-10-05", updated: "2026-10-03", archived: false},
		{name: "quiz with the last due date", id: "10", kind: "quiz", state: "assigned", chapter: 3, due: "2026-10-06", updated: "2026-10-02", archived: false},
		{name: "problem without a due date", id: "11", kind: "problem", state: "assigned", chapter: nil, due: nil, updated: "2026-10-01", archived: false},
		{name: "archived submission", id: "12", kind: "problem", state: "submitted", chapter: nil, due: "2026-09-01", updated: "2026-09-01", archived: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			exec(t, f.pool, "UPDATE homework_tasks SET state='cancelled'")
			for _, row := range rows {
				id := prefix + row.id
				problem := row.kind == "problem"
				if problem {
					exec(t, f.pool, "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,archived_at,created_at,updated_at) VALUES ($1,$2,$3,'easy','[]','','',$4,'now','now')", id, id, row.name, reviewsNull(row.archived, "now"))
				}
				exec(t, f.pool, "INSERT INTO homework (id,title,instructions,due_date,created_at,updated_at) VALUES ($1,$2,'',$3,'now','now')", id, row.name, row.due)
				exec(t, f.pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,kind,set_slug,chapter,state,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,'now',$8)", id, id, reviewsNull(problem, id), row.kind, reviewsNull(!problem, "permutations-and-combinations"), row.chapter, row.state, row.updated)
			}
			exec(t, f.pool, "INSERT INTO feedback (id,problem_id,author,kind,body,created_at) VALUES ($1,$2,'parent','reply','Add a proof.','2026-10-07')", newID(), prefix+"05")
			exec(t, f.pool, "INSERT INTO feedback (id,problem_id,author,kind,body,created_at) VALUES ($1,$2,'student','reply','Please check.','2026-10-07')", newID(), prefix+"11")
			homework := reviewsItems(t, reviewsDashboard(t, f, "2026-10-07", f.as(t, tt.role))["homework"])
			want := []any{}
			for _, id := range tt.order {
				want = append(want, map[string]any{"id": prefix + id, "unread": slices.Contains(tt.unread, id)})
			}
			reviewsEqualJSON(t, homeworkFields(homework, "id", "unread"), want)
		})
	}
}

func TestReviewsAddCalendarDays(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		days     int
		expected string
	}{
		{name: "month boundary", date: "2026-01-31", days: 1, expected: "2026-02-01"},
		{name: "leap day", date: "2028-02-28", days: 1, expected: "2028-02-29"},
		{name: "year boundary", date: "2026-12-29", days: 7, expected: "2027-01-05"},
		{name: "daylight-saving calendar interval", date: "2026-03-07", days: 3, expected: "2026-03-10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addDays(tt.date, tt.days)
			if err != nil || got != tt.expected {
				t.Errorf("addDays(%s, %d) = %q, %v, want %q", tt.date, tt.days, got, err, tt.expected)
			}
		})
	}
}

func TestReviewsCountsAcceptedSeparatelyFromUnderstanding(t *testing.T) {
	f := reviewsSetup(t)
	problemsMatchBody(t, reviewsDashboard(t, f, "2026-10-07", f.cookie), map[string]any{"counts": map[string]any{"recorded": 3, "solved": 2, "independent": 1}, "topics": []any{map[string]any{"topic": "Array", "solved": 2, "independent": 1}, map[string]any{"topic": "Graph Theory", "solved": 1, "independent": 0}}})
	value := map[string]any{"version": 1, "result": "needs_practice", "note": "Missed the invariant.", "reviewedOn": "2026-10-07", "nextReviewDate": "2026-10-08"}
	if w := request(t, f.s, "/problems/"+f.first+"/reviews", http.MethodPost, value, f.cookie); w.Code != http.StatusCreated {
		t.Errorf("review: status = %d, want 201; body %s", w.Code, w.Body.String())
	}
	problemsMatchBody(t, reviewsDashboard(t, f, "2026-10-07", f.cookie), map[string]any{"counts": map[string]any{"solved": 2, "independent": 0}})
	problemsMatchBody(t, decode[any](t, request(t, f.s, "/problems/"+f.first, http.MethodGet, nil, f.cookie)), map[string]any{"solved": true, "understanding": "needs_practice", "nextReviewDate": "2026-10-08"})
}

func TestReviewsCountsSolvedOnceAfterTopicNormalization(t *testing.T) {
	f := reviewsSetup(t)
	rows := []struct {
		name, topics, understanding, acceptance, state string
		archived                                       bool
	}{
		{name: "canonical", topics: `["Hash Table"]`, understanding: "independent", acceptance: "accepted", state: "saved", archived: false},
		{name: "compact", topics: `["HashTable"]`, understanding: "independent", acceptance: "accepted", state: "saved", archived: false},
		{name: "Russian", topics: `["Хеш-таблица"]`, understanding: "independent", acceptance: "accepted", state: "saved", archived: false},
		{name: "several aliases", topics: `["Hash Table","HashTable","Хеш-таблица"]`, understanding: "with_help", acceptance: "accepted", state: "saved", archived: false},
		{name: "archived", topics: `["HashTable"]`, understanding: "independent", acceptance: "accepted", state: "saved", archived: true},
		{name: "unaccepted", topics: `["Хеш-таблица"]`, understanding: "independent", acceptance: "not_accepted", state: "saved", archived: false},
		{name: "draft", topics: `["Hash Table"]`, understanding: "independent", acceptance: "accepted", state: "draft", archived: false},
	}
	for _, row := range rows {
		id := newID()
		exec(t, f.pool, "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,understanding,archived_at,next_review_date,created_at,updated_at) VALUES ($1,$2,$3,'easy',$4,'','',$5,$6,'2026-10-27','now','now')", id, id, row.name, row.topics, row.understanding, reviewsNull(row.archived, "now"))
		exec(t, f.pool, reviewsInsertAttempt, newID(), id, row.state, row.acceptance, row.understanding)
		if row.name == "several aliases" {
			exec(t, f.pool, reviewsInsertAttempt, newID(), id, "saved", "accepted", "independent")
		}
	}
	dashboard := reviewsDashboard(t, f, "2026-10-27", f.cookie)
	reviewsEqualJSON(t, dashboard["topics"], []any{map[string]any{"topic": "Array", "solved": 2, "independent": 1}, map[string]any{"topic": "Graph Theory", "solved": 1, "independent": 0}, map[string]any{"topic": "Hash Table", "solved": 4, "independent": 3}})
	due := reviewsItems(t, dashboard["dueReviews"])
	index := slices.IndexFunc(due, func(problem map[string]any) bool { return problem["title"] == "several aliases" })
	if index < 0 {
		t.Fatalf("dueReviews = %v, want the problem «several aliases»", due)
	}
	problemsMatchBody(t, due[index], map[string]any{"topics": []string{"Hash Table"}, "solved": true, "understanding": "with_help"})
}

func TestReviewsTopicGroupsBeforeLimit(t *testing.T) {
	f := reviewsSetup(t)
	rows := []struct{ name, topics string }{
		{name: "first ten", topics: `["Topic 01","Topic 02","Topic 03","Topic 04","Topic 05","Topic 06","Topic 07","Topic 08","Topic 09","Topic 10"]`},
		{name: "second ten", topics: `["Topic 11","Topic 12","Topic 13","Topic 14","Topic 15","Topic 16","Topic 17","Topic 18","Topic 19","Topic 20"]`},
		{name: "third ten", topics: `["Topic 21","Topic 22","Topic 23","Topic 24","Topic 25","Topic 26","Topic 27","Topic 28","Topic 29","Topic 30"]`},
		{name: "fourth ten", topics: `["Topic 31","Topic 32","Topic 33","Topic 34","Topic 35","Topic 36","Topic 37","Topic 38","Topic 39","Topic 40"]`},
		{name: "last nine", topics: `["Topic 41","Topic 42","Topic 43","Topic 44","Topic 45","Topic 46","Topic 47","Topic 48","Topic 49"]`},
		{name: "legacy group after raw limit", topics: `["Хеш-таблица"]`},
	}
	for _, row := range rows {
		id := newID()
		exec(t, f.pool, "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,understanding,created_at,updated_at) VALUES ($1,$2,$3,'easy',$4,'','','independent','now','now')", id, id, row.name, row.topics)
		exec(t, f.pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'saved','{}','','accepted','independent','now','now')", newID(), id)
	}
	topics := reviewsItems(t, reviewsDashboard(t, f, "2026-10-27", f.cookie)["topics"])
	if len(topics) != 50 {
		t.Fatalf("%d topics, want 50", len(topics))
	}
	reviewsEqualJSON(t, topics[:3], []any{map[string]any{"topic": "Array", "solved": 2, "independent": 1}, map[string]any{"topic": "Graph Theory", "solved": 1, "independent": 0}, map[string]any{"topic": "Hash Table", "solved": 1, "independent": 1}})
	reviewsEqualJSON(t, topics[49], map[string]any{"topic": "Topic 47", "solved": 1, "independent": 1})
}

func TestReviewsDueReviewCount(t *testing.T) {
	tests := []struct {
		name       string
		today      string
		dueReviews int
	}{
		{name: "before due day", today: "2026-10-06", dueReviews: 0},
		{name: "same-day revision", today: "2026-10-07", dueReviews: 3},
		{name: "revision overdue the next day", today: "2026-10-08", dueReviews: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			reviewsEqualJSON(t, reviewsDashboard(t, f, tt.today, f.cookie)["counts"], map[string]any{"recorded": 3, "solved": 2, "independent": 1, "waitingReview": 0, "dueReviews": tt.dueReviews})
		})
	}
}

func TestReviewsDashboardHomeworkByState(t *testing.T) {
	tests := []struct {
		name  string
		state string
		rows  int
	}{
		{name: "active task with its homework", state: "assigned", rows: 1},
		{name: "submitted task waits for review", state: "submitted", rows: 1},
		{name: "completed task leaves the dashboard", state: "completed", rows: 0},
		{name: "cancelled task leaves the dashboard", state: "cancelled", rows: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			exec(t, f.pool, "UPDATE homework_tasks SET state=$1", tt.state)
			homework := reviewsDashboard(t, f, "2026-10-07", f.cookie)["homework"]
			want := []any{map[string]any{"id": f.taskID, "homeworkId": f.homeworkID, "homeworkTitle": "Week 41", "instructions": "Explain", "dueDate": "2026-10-07", "postponedTo": nil, "problemId": f.first, "title": "accepted independently", "href": "/leetcode/problems/" + f.first, "kind": "problem", "setSlug": nil, "chapter": nil, "completedAt": nil, "difficulty": "medium", "state": tt.state, "submissionId": nil, "submissionAttemptId": nil, "version": 1, "createdAt": "now", "updatedAt": "now", "unread": false}}
			reviewsEqualJSON(t, homework, want[:tt.rows])
		})
	}
}

func TestReviewsDashboardInvalidToday(t *testing.T) {
	tests := []struct {
		name  string
		today string
	}{
		{name: "impossible day", today: "2026-02-30"},
		{name: "missing local date", today: ""},
		{name: "timestamp in place of date", today: "2026-10-07T00:00:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			if w := request(t, f.s, "/dashboard?today="+tt.today, http.MethodGet, nil, f.cookie); w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400; body %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestReviewsRejectedReview(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		nextReviewDate string
		expected       int
	}{
		{name: "parent review", role: "parent", nextReviewDate: "2026-10-08", expected: 403},
		{name: "impossible next review date", role: "student", nextReviewDate: "2026-02-30", expected: 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			value := map[string]any{"version": 1, "result": "with_help", "note": "", "reviewedOn": "2026-10-07", "nextReviewDate": tt.nextReviewDate}
			if w := request(t, f.s, "/problems/"+f.first+"/reviews", http.MethodPost, value, f.as(t, tt.role)); w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			if got := queryValue[int64](t, f.pool, reviewsCount); got != 0 {
				t.Errorf("review count = %d, want 0", got)
			}
		})
	}
}

func TestReviewsDashboardCounts(t *testing.T) {
	tests := []struct {
		name   string
		change string
		counts map[string]any
	}{
		{name: "repeated accepted attempt counts one solved problem", change: "attempt", counts: map[string]any{"recorded": 3, "solved": 2, "independent": 1, "dueReviews": 3}},
		{name: "archived problem leaves the counts", change: "archive", counts: map[string]any{"recorded": 2, "solved": 2, "independent": 1, "dueReviews": 2}},
		{name: "submitted homework waits for review", change: "submit", counts: map[string]any{"waitingReview": 1}},
		{name: "archived problem hides submitted homework", change: "submit archived", counts: map[string]any{"waitingReview": 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			if tt.change == "attempt" {
				exec(t, f.pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'saved','{}','','accepted','independent','now','now')", newID(), f.first)
			}
			if tt.change == "archive" {
				exec(t, f.pool, "UPDATE problems SET archived_at='now' WHERE slug='three'")
			}
			if strings.HasPrefix(tt.change, "submit") {
				exec(t, f.pool, "UPDATE homework_tasks SET state='submitted'")
			}
			if tt.change == "submit archived" {
				exec(t, f.pool, "UPDATE problems SET archived_at='now' WHERE id=$1", f.first)
			}
			problemsMatchBody(t, reviewsDashboard(t, f, "2026-10-07", f.cookie), map[string]any{"counts": tt.counts})
		})
	}
}

func TestReviewsProgressVersion(t *testing.T) {
	tests := []struct {
		name     string
		earlier  string
		later    string
		expected int
		fields   map[string]any
		reviews  int64
	}{
		{name: "student review after a parent edit", earlier: "edit", later: "review", expected: 201, fields: map[string]any{"version": 2, "progressVersion": 2, "nextReviewDate": nil}, reviews: 1},
		{name: "second review with the same progress version", earlier: "review", later: "review", expected: 409, fields: map[string]any{"version": 1, "progressVersion": 2}, reviews: 1},
		{name: "parent edit after the student marks a draft solved", earlier: "solved", later: "edit", expected: 200, fields: map[string]any{"version": 2, "progressVersion": 2}, reviews: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			draftID, parent := newID(), loginAs(t, f.s, "parent")
			exec(t, f.pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'draft',$3,'','accepted','with_help','now','now')", draftID, f.first, `{"notes":"Explain the invariant.","approaches":[]}`)
			problem := decode[map[string]any](t, request(t, f.s, "/problems/"+f.first, http.MethodGet, nil, f.cookie))
			// write returns the status of the change.
			write := func(kind string) int {
				switch kind {
				case "solved":
					// The route /problems/{id}/solved is not in Go. homeworkSolve writes its rows.
					homeworkSolve(t, homeworkFixture{pool: f.pool, problemID: f.first, attemptID: draftID})
					return http.StatusOK
				case "edit":
					return request(t, f.s, "/problems/"+f.first, http.MethodPatch, map[string]any{"url": "https://leetcode.com/problems/one/", "number": nil, "title": "Edited", "difficulty": "medium", "topics": []string{}, "summary": "", "version": problem["version"]}, parent).Code
				default:
					return request(t, f.s, "/problems/"+f.first+"/reviews", http.MethodPost, map[string]any{"version": problem["progressVersion"], "result": "with_help", "note": "", "reviewedOn": "2026-10-07", "nextReviewDate": nil}, f.cookie).Code
				}
			}
			if code := write(tt.earlier); code < 200 || code > 299 {
				t.Errorf("%s: status = %d, want a 2xx status", tt.earlier, code)
			}
			if code := write(tt.later); code != tt.expected {
				t.Errorf("%s: status = %d, want %d", tt.later, code, tt.expected)
			}
			problemsMatchBody(t, decode[any](t, request(t, f.s, "/problems/"+f.first, http.MethodGet, nil, f.cookie)), tt.fields)
			if got := queryValue[int64](t, f.pool, reviewsCount); got != tt.reviews {
				t.Errorf("review count = %d, want %d", got, tt.reviews)
			}
		})
	}
}

// reviewsWindow returns the 30 days after 2026-10-07, as window30 of test/reviews.test.ts.
func reviewsWindow() []string {
	start := time.Date(2026, time.October, 7, 0, 0, 0, 0, time.UTC)
	days := make([]string, 30)
	for i := range days {
		days[i] = start.AddDate(0, 0, i+1).Format(time.DateOnly)
	}
	return days
}

// reviewsFresh returns count problems p0, p1, ... without a review date, as fresh of test/reviews.test.ts.
func reviewsFresh(count int) []spreadProblem {
	problems := make([]spreadProblem, count)
	for i := range problems {
		problems[i] = spreadProblem{ID: fmt.Sprintf("p%d", i)}
	}
	return problems
}

// reviewsDays gives the problem p<i> the day days[i].
func reviewsDays(days []string) map[string]string {
	result := map[string]string{}
	for i, day := range days {
		result[fmt.Sprintf("p%d", i)] = day
	}
	return result
}

func TestReviewsSpreadDays(t *testing.T) {
	first := func(int) int { return 0 }
	last := func(count int) int { return count - 1 }
	window := reviewsWindow()
	tests := []struct {
		name     string
		problems []spreadProblem
		random   func(count int) int
		days     map[string]string
	}{
		{name: "one problem gets the first free day", problems: reviewsFresh(1), random: first, days: map[string]string{"p0": "2026-10-08"}},
		{name: "the last free day is day 30", problems: reviewsFresh(1), random: last, days: map[string]string{"p0": "2026-11-06"}},
		{name: "30 problems get 30 distinct days", problems: reviewsFresh(30), random: first, days: reviewsDays(window)},
		{name: "31 problems use one day twice", problems: reviewsFresh(31), random: first, days: reviewsDays(slices.Concat(window, []string{"2026-10-08"}))},
		{name: "a date in the window stays and its day is used", problems: []spreadProblem{{ID: "kept", NextReviewDate: new("2026-10-08")}, {ID: "new"}}, random: first, days: map[string]string{"new": "2026-10-09"}},
		{name: "day 30 stays", problems: []spreadProblem{{ID: "kept", NextReviewDate: new("2026-11-06")}}, random: first, days: map[string]string{}},
		{name: "a past date gets a new day", problems: []spreadProblem{{ID: "past", NextReviewDate: new("2026-10-01")}}, random: first, days: map[string]string{"past": "2026-10-08"}},
		{name: "today gets a new day", problems: []spreadProblem{{ID: "today", NextReviewDate: new("2026-10-07")}}, random: first, days: map[string]string{"today": "2026-10-08"}},
		{name: "a date after day 30 gets a new day", problems: []spreadProblem{{ID: "late", NextReviewDate: new("2026-11-07")}}, random: last, days: map[string]string{"late": "2026-11-06"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := spreadDays("2026-10-07", tt.problems, tt.random)
			if err != nil || !maps.Equal(got, tt.days) {
				t.Errorf("spreadDays = %v, %v, want %v", got, err, tt.days)
			}
		})
	}
}

func TestReviewsDistribute(t *testing.T) {
	tests := []struct {
		name            string
		role            string
		state           string
		archived        bool
		today           string
		expected        int
		body            map[string]any
		date            string
		progressVersion int64
	}{
		{name: "student spreads a completed task", role: "student", state: "completed", archived: false, today: "2026-10-07", expected: 200, body: map[string]any{"count": 1, "from": "2026-10-08", "to": "2026-10-08"}, date: "2026-10-08", progressVersion: 2},
		{name: "parent spreads a completed task", role: "parent", state: "completed", archived: false, today: "2026-10-07", expected: 200, body: map[string]any{"count": 1, "from": "2026-10-08", "to": "2026-10-08"}, date: "2026-10-08", progressVersion: 2},
		{name: "archived problem is skipped", role: "student", state: "completed", archived: true, today: "2026-10-07", expected: 200, body: map[string]any{"count": 0, "from": nil, "to": nil}, date: "2026-10-07", progressVersion: 1},
		{name: "no completed tasks", role: "student", state: "assigned", archived: false, today: "2026-10-07", expected: 200, body: map[string]any{"count": 0, "from": nil, "to": nil}, date: "2026-10-07", progressVersion: 1},
		{name: "impossible today", role: "student", state: "completed", archived: false, today: "2026-02-30", expected: 400, body: map[string]any{"error": map[string]any{"code": "invalid_input"}}, date: "2026-10-07", progressVersion: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			// getRandomValues returns the unchanged array, so the random value is 0.
			f.s.random = func(int) int { return 0 }
			exec(t, f.pool, "UPDATE homework_tasks SET state=$1", tt.state)
			exec(t, f.pool, "UPDATE problems SET archived_at=$1 WHERE id=$2", reviewsNull(tt.archived, "now"), f.first)
			w := request(t, f.s, "/reviews/distribute", http.MethodPost, map[string]any{"today": tt.today}, f.as(t, tt.role))
			if w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), tt.body)
			if got, want := homeworkRow(t, f.pool, reviewsSelectDates, f.first), map[string]any{"next_review_date": tt.date, "progress_version": tt.progressVersion}; !reflect.DeepEqual(got, want) {
				t.Errorf("problem = %v, want %v", got, want)
			}
		})
	}
}

// The test of TypeScript mocks the read of the problems and records a recall after it.
// The Go handler calls s.random after the read and before the write, so the recall runs in s.random.
func TestReviewsDistributeKeepsConcurrentRecall(t *testing.T) {
	tests := []struct {
		name    string
		another bool
	}{
		{name: "preserves a concurrent recall and reports only applied dates (another problem: false)", another: false},
		{name: "preserves a concurrent recall and reports only applied dates (another problem: true)", another: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			exec(t, f.pool, "UPDATE homework_tasks SET state='completed'")
			if tt.another {
				exec(t, f.pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) SELECT $1,$2,id,'completed','now','now' FROM problems WHERE slug='two'", newID(), f.homeworkID)
			}
			calls := 0
			f.s.random = func(int) int {
				calls++
				if calls == 1 {
					recall := request(t, f.s, "/problems/"+f.first+"/reviews", http.MethodPost, map[string]any{"version": 1, "result": "with_help", "note": "", "reviewedOn": "2026-10-07", "nextReviewDate": "2026-10-25"}, f.cookie)
					if recall.Code != http.StatusCreated {
						t.Errorf("recall: status = %d, want 201; body %s", recall.Code, recall.Body.String())
					}
				}
				return 0
			}
			w := request(t, f.s, "/reviews/distribute", http.MethodPost, map[string]any{"today": "2026-10-07"}, loginAs(t, f.s, "parent"))
			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			if calls == 0 {
				t.Error("the handler did not call s.random, so no recall ran between the read and the write")
			}
			if got, want := homeworkRow(t, f.pool, reviewsSelectDates, f.first), map[string]any{"next_review_date": "2026-10-25", "progress_version": int64(2)}; !reflect.DeepEqual(got, want) {
				t.Errorf("problem = %v, want %v", got, want)
			}
			var date any
			if tt.another {
				two := queryValue[string](t, f.pool, "SELECT next_review_date FROM problems WHERE slug='two'")
				if !regexp.MustCompile(`^2026-10-0[89]$`).MatchString(two) {
					t.Errorf("next review date of two = %q, want 2026-10-08 or 2026-10-09", two)
				}
				date = two
			}
			count := 0
			if tt.another {
				count = 1
			}
			reviewsEqualJSON(t, decode[any](t, w), map[string]any{"count": count, "from": date, "to": date})
		})
	}
}

// The window is [today + 15, today + 45]. The random value 30 gives the last day.
func TestReviewsPractice(t *testing.T) {
	tests := []struct {
		name   string
		role   string
		random int
		today  string
		status int
		body   map[string]any
		after  map[string]any
	}{
		{name: "student gets the first day of the window", role: "student", random: 0, today: "2026-10-07", status: 200, body: map[string]any{"nextReviewDate": "2026-10-22"}, after: map[string]any{"understanding": "needs_practice", "next_review_date": "2026-10-22", "progress_version": int64(2)}},
		{name: "parent gets the last day of the window", role: "parent", random: 30, today: "2026-10-07", status: 200, body: map[string]any{"nextReviewDate": "2026-11-21"}, after: map[string]any{"understanding": "needs_practice", "next_review_date": "2026-11-21", "progress_version": int64(2)}},
		{name: "practice over a year boundary", role: "student", random: 30, today: "2026-12-20", status: 200, body: map[string]any{"nextReviewDate": "2027-02-03"}, after: map[string]any{"understanding": "needs_practice", "next_review_date": "2027-02-03", "progress_version": int64(2)}},
		{name: "impossible today", role: "student", random: 0, today: "2026-02-30", status: 400, body: map[string]any{"error": map[string]any{"code": "invalid_input"}}, after: map[string]any{"understanding": "independent", "next_review_date": "2026-10-07", "progress_version": int64(1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			account := f.as(t, tt.role)
			f.s.random = func(count int) int { return tt.random % count }
			w := request(t, f.s, "/problems/"+f.first+"/practice", http.MethodPost, map[string]any{"today": tt.today}, account)
			if w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), tt.body)
			if got := homeworkRow(t, f.pool, "SELECT understanding,next_review_date,progress_version FROM problems WHERE id=$1", f.first); !reflect.DeepEqual(got, tt.after) {
				t.Errorf("problem = %v, want %v", got, tt.after)
			}
		})
	}
}

func TestReviewsDueReviewFromHomework(t *testing.T) {
	tests := []struct {
		name         string
		state        string
		fromHomework bool
	}{
		{name: "due problem with a completed task opens a pick", state: "completed", fromHomework: true},
		{name: "due problem with an active task opens the review", state: "assigned", fromHomework: false},
		{name: "due problem with a cancelled task opens the review", state: "cancelled", fromHomework: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			exec(t, f.pool, "UPDATE homework_tasks SET state=$1", tt.state)
			got := map[string]any{}
			for _, problem := range reviewsItems(t, reviewsDashboard(t, f, "2026-10-07", f.cookie)["dueReviews"]) {
				slug, _ := problem["slug"].(string)
				got[slug] = problem["fromHomework"]
			}
			reviewsEqualJSON(t, got, map[string]any{"one": tt.fromHomework, "two": false, "three": false})
		})
	}
}

// The attempts of beforeEach are saved. A draft of a problem with an accepted attempt is not a next step.
func TestReviewsDrafts(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		archived bool
		titles   []string
	}{
		{name: "draft of an unaccepted problem is a next step", slug: "three", archived: false, titles: []string{"unaccepted independently"}},
		{name: "draft of an accepted problem is not a next step", slug: "one", archived: false, titles: []string{}},
		{name: "draft of an archived problem is not a next step", slug: "three", archived: true, titles: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reviewsSetup(t)
			exec(t, f.pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) SELECT $1,id,'draft','{}','','not_submitted','needs_practice','now','now' FROM problems WHERE slug=$2", newID(), tt.slug)
			exec(t, f.pool, "UPDATE problems SET archived_at=$1 WHERE slug=$2", reviewsNull(tt.archived, "now"), tt.slug)
			titles := []any{}
			for _, draft := range reviewsItems(t, reviewsDashboard(t, f, "2026-10-07", f.cookie)["drafts"]) {
				titles = append(titles, draft["problemTitle"])
			}
			reviewsEqualJSON(t, titles, tt.titles)
		})
	}
}

package notebook

import (
	"context"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// The cases of test/homework.test.ts. The cases that read /dashboard are in homework_dashboard_test.go.

// The prepared statements of test/homework.test.ts.
const (
	homeworkInsertProblem      = "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,$3,'medium','[]','','','now','now')"
	homeworkInsertBundle       = "INSERT INTO homework (id,title,instructions,due_date,created_at,updated_at) VALUES ($1,$2,'Explain the invariant',$3,$4,$5)"
	homeworkInsertTask         = "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,'now')"
	homeworkInsertFeedback     = "INSERT INTO feedback (id,problem_id,author,kind,body,created_at) VALUES ($1,$2,$3,$4,'Note',$5)"
	homeworkInsertRead         = "INSERT INTO discussion_reads (role,problem_id,read_at) VALUES ($1,$2,$3)"
	homeworkInsertPostponement = "INSERT INTO postponements VALUES ($1,'2026-10-09','2026-10-10','now')"
)

// homeworkFixture is the state that beforeEach of test/homework.test.ts makes.
type homeworkFixture struct {
	pool                                                   *pgxpool.Pool
	s                                                      *Server
	student, parent                                        string
	problemID, attemptID, homeworkID, taskID, submissionID string
}

// as returns the session cookie of the role.
func (f homeworkFixture) as(role string) string {
	if role == "parent" {
		return f.parent
	}
	return f.student
}

// homeworkSetup adds the problem Two Sum with a draft, and the homework "Two Sum practice" with one assigned task of the problem.
func homeworkSetup(t *testing.T) homeworkFixture {
	t.Helper()
	pool := database(t)
	s := newServerWithSets(t, pool, homeworkSets)
	f := homeworkFixture{pool: pool, s: s, student: loginAs(t, s, "student"), parent: loginAs(t, s, "parent")}
	f.problemID = addProblem(t, s, f.student, map[string]any{"url": "https://leetcode.com/problems/two-sum/", "number": 1, "title": "Two Sum", "difficulty": "easy", "topics": []string{}, "summary": ""})
	w := request(t, s, "/problems/"+f.problemID+"/attempts", http.MethodPost, "{}", f.student)
	if w.Code != http.StatusCreated {
		t.Fatalf("new draft: status %d, body %s", w.Code, w.Body.String())
	}
	f.attemptID = decode[struct {
		ID string `json:"id"`
	}](t, w).ID
	f.homeworkID, f.taskID, f.submissionID = newID(), newID(), newID()
	exec(t, pool, "UPDATE attempts SET document=jsonb_set(document::jsonb,'{notes}',to_jsonb($1::text))::text WHERE id=$2", "Explain the counting invariant.", f.attemptID)
	exec(t, pool, homeworkInsertBundle, f.homeworkID, "Two Sum practice", "2026-10-10", "2026-10-04", "2026-10-04")
	exec(t, pool, homeworkInsertTask, f.taskID, f.homeworkID, f.problemID, "assigned", "2026-10-04")
	return f
}

// homeworkPage is a Page of the API with its items as JSON objects.
type homeworkPage struct {
	Items []map[string]any `json:"items"`
}

// homeworkDetails is the body of GET /homework/{id}.
type homeworkDetails struct {
	Homework map[string]any   `json:"homework"`
	Tasks    []map[string]any `json:"tasks"`
}

// homeworkTasks returns the tasks of GET /homework/{id} for the session cookie.
func homeworkTasks(t *testing.T, s *Server, homeworkID, cookie string) []map[string]any {
	t.Helper()
	return decode[homeworkDetails](t, request(t, s, "/homework/"+homeworkID, http.MethodGet, nil, cookie)).Tasks
}

// homeworkFirst returns the first item of a list. The test stops when the list is empty.
func homeworkFirst(t *testing.T, items []map[string]any) map[string]any {
	t.Helper()
	if len(items) == 0 {
		t.Fatal("the list is empty, want at least one item")
	}
	return items[0]
}

// homeworkFields returns the listed fields of each item, as the map callbacks of test/homework.test.ts. A missing field becomes nil.
func homeworkFields(items []map[string]any, keys ...string) []any {
	picked := []any{}
	for _, item := range items {
		fields := map[string]any{}
		for _, key := range keys {
			fields[key] = item[key]
		}
		picked = append(picked, fields)
	}
	return picked
}

// homeworkRow returns the first row of a query as a map from the column name to the value.
func homeworkRow(t *testing.T, pool *pgxpool.Pool, sql string, args ...any) map[string]any {
	t.Helper()
	rows, err := pool.Query(context.Background(), sql, args...)
	if err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	values, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil || len(values) == 0 {
		t.Fatalf("%s: %d rows, %v", sql, len(values), err)
	}
	return values[0]
}

// homeworkSolve writes the rows of POST /problems/{id}/solved of src/leetcode/judge.ts and returns the ID of the submission.
// The draft becomes the accepted attempt. The assigned task gets a submission of the draft and becomes submitted.
func homeworkSolve(t *testing.T, f homeworkFixture) string {
	t.Helper()
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	submissionID := newID()
	// No tutor message exists, so the understanding is independent.
	exec(t, f.pool, "UPDATE attempts SET state='saved',acceptance='accepted',understanding='independent',version=version+1,saved_at=$1,updated_at=$1 WHERE id=$2 AND problem_id=$3 AND state='draft'", now, f.attemptID, f.problemID)
	exec(t, f.pool, "UPDATE problems SET understanding=(SELECT understanding FROM attempts WHERE id=$1),progress_version=progress_version+1,updated_at=$2 WHERE id=$3", f.attemptID, now, f.problemID)
	exec(t, f.pool, "INSERT INTO submissions (id,task_id,attempt_id,task_version,created_at) SELECT $1,id,$2,version+1,$3 FROM homework_tasks WHERE problem_id=$4 AND state IN ('assigned','in_progress','changes_requested')", submissionID, f.attemptID, now, f.problemID)
	exec(t, f.pool, "UPDATE homework_tasks SET state='submitted',current_submission_id=$1,version=version+1,updated_at=$2 WHERE id=(SELECT task_id FROM submissions WHERE id=$1)", submissionID, now)
	return queryValue[string](t, f.pool, "SELECT id FROM submissions")
}

// homeworkBundles adds the problems, the homework and the tasks of bundles() of test/homework.test.ts.
func homeworkBundles(t *testing.T, f homeworkFixture) {
	t.Helper()
	two, three := newID(), newID()
	arrays, backtracking, graphs := newID(), newID(), newID()
	exec(t, f.pool, homeworkInsertProblem, two, "three-sum", "3Sum")
	exec(t, f.pool, homeworkInsertProblem, three, "subsets", "Subsets")
	exec(t, f.pool, homeworkInsertBundle, arrays, "Week 40: arrays", nil, "2026-10-01", "2026-10-01")
	exec(t, f.pool, homeworkInsertBundle, backtracking, "Week 41: backtracking", nil, "2026-10-02", "2026-10-02")
	exec(t, f.pool, homeworkInsertBundle, graphs, "Week 42: graphs", nil, "2026-10-03", "2026-10-03")
	exec(t, f.pool, homeworkInsertTask, newID(), arrays, f.problemID, "completed", "2026-10-01")
	exec(t, f.pool, homeworkInsertTask, newID(), arrays, f.problemID, "cancelled", "2026-10-01")
	exec(t, f.pool, homeworkInsertTask, newID(), backtracking, two, "submitted", "2026-10-02")
	exec(t, f.pool, homeworkInsertTask, newID(), backtracking, three, "changes_requested", "2026-10-02")
	exec(t, f.pool, homeworkInsertTask, newID(), backtracking, f.problemID, "completed", "2026-10-02")
}

func TestHomeworkCreate(t *testing.T) {
	invalid := map[string]any{"error": map[string]any{"code": "invalid_input"}}
	tests := []struct {
		name     string
		role     string
		body     map[string]any
		expected int
		stored   int64
		fields   map[string]any
	}{
		{name: "parent creates homework", role: "parent", body: map[string]any{"title": " Week 41: backtracking ", "instructions": "Explain each choice.", "dueDate": "2026-10-12"}, expected: 201, stored: 2, fields: map[string]any{"title": "Week 41: backtracking", "instructions": "Explain each choice.", "dueDate": "2026-10-12", "version": 1}},
		{name: "title at 200 characters", role: "parent", body: map[string]any{"title": strings.Repeat("x", 200), "instructions": "", "dueDate": nil}, expected: 201, stored: 2, fields: map[string]any{"title": strings.Repeat("x", 200), "dueDate": nil}},
		{name: "title over 200 characters", role: "parent", body: map[string]any{"title": strings.Repeat("x", 201), "instructions": "", "dueDate": nil}, expected: 400, stored: 1, fields: invalid},
		{name: "blank title", role: "parent", body: map[string]any{"title": "  ", "instructions": "", "dueDate": nil}, expected: 400, stored: 1, fields: invalid},
		{name: "impossible due date", role: "parent", body: map[string]any{"title": "Week 41", "instructions": "", "dueDate": "2026-02-30"}, expected: 400, stored: 1, fields: invalid},
		{name: "student cannot create homework", role: "student", body: map[string]any{"title": "Week 41", "instructions": "", "dueDate": nil}, expected: 403, stored: 1, fields: map[string]any{"error": map[string]any{"code": "forbidden"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			w := request(t, f.s, "/homework", http.MethodPost, tt.body, f.as(tt.role))
			if w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), tt.fields)
			if got := queryValue[int64](t, f.pool, "SELECT count(*) AS count FROM homework"); got != tt.stored {
				t.Errorf("homework count = %d, want %d", got, tt.stored)
			}
		})
	}
}

func TestHomeworkEdit(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		stored   int
		expected int
		row      map[string]any
	}{
		{name: "parent edits homework", role: "parent", stored: 1, expected: 200, row: map[string]any{"title": "Week 41", "instructions": "Compare two approaches.", "due_date": nil, "version": int64(2)}},
		{name: "stale homework edit", role: "parent", stored: 2, expected: 409, row: map[string]any{"title": "Two Sum practice", "instructions": "Explain the invariant", "due_date": "2026-10-10", "version": int64(2)}},
		{name: "student cannot edit homework", role: "student", stored: 1, expected: 403, row: map[string]any{"title": "Two Sum practice", "instructions": "Explain the invariant", "due_date": "2026-10-10", "version": int64(1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			exec(t, f.pool, "UPDATE homework SET version=$1", tt.stored)
			w := request(t, f.s, "/homework/"+f.homeworkID, http.MethodPatch, map[string]any{"version": 1, "title": "Week 41", "instructions": "Compare two approaches.", "dueDate": nil}, f.as(tt.role))
			if w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			if got := homeworkRow(t, f.pool, "SELECT title,instructions,due_date,version FROM homework"); !reflect.DeepEqual(got, tt.row) {
				t.Errorf("homework = %v, want %v", got, tt.row)
			}
		})
	}
}

// The homework of homeworkSetup is due on 2026-10-10.
func TestHomeworkEditResetsPostponement(t *testing.T) {
	tests := []struct {
		name        string
		stored      int
		dueDate     any
		expected    int
		postponedTo any
	}{
		{name: "new due date resets Сила Бомжа", stored: 1, dueDate: "2026-10-20", expected: 200, postponedTo: nil},
		{name: "removed due date resets Сила Бомжа", stored: 1, dueDate: nil, expected: 200, postponedTo: nil},
		{name: "same due date keeps Сила Бомжа", stored: 1, dueDate: "2026-10-10", expected: 200, postponedTo: "2026-10-10"},
		{name: "stale edit keeps Сила Бомжа", stored: 2, dueDate: "2026-10-20", expected: 409, postponedTo: "2026-10-10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			exec(t, f.pool, "UPDATE homework SET version=$1", tt.stored)
			exec(t, f.pool, homeworkInsertPostponement, f.homeworkID)
			w := request(t, f.s, "/homework/"+f.homeworkID, http.MethodPatch, map[string]any{"version": 1, "title": "Two Sum practice", "instructions": "Explain the invariant", "dueDate": tt.dueDate}, f.parent)
			if w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, request(t, f.s, "/homework/"+f.homeworkID, http.MethodGet, nil, f.parent)), map[string]any{"homework": map[string]any{"postponedTo": tt.postponedTo}})
		})
	}
}

// The server takes today from its own UTC clock.
func TestHomeworkPostpone(t *testing.T) {
	now := time.Now().UTC()
	today, tomorrow, yesterday := now.Format(time.DateOnly), now.AddDate(0, 0, 1).Format(time.DateOnly), now.AddDate(0, 0, -1).Format(time.DateOnly)
	tests := []struct {
		name        string
		role        string
		due         any
		used        bool
		version     int
		expected    int
		code        string
		dueAfter    any
		postponedTo any
	}{
		{name: "student moves the due date by one day", role: "student", due: today, used: false, version: 1, expected: 200, code: "", dueAfter: tomorrow, postponedTo: tomorrow},
		{name: "second use of Сила Бомжа", role: "student", due: today, used: true, version: 1, expected: 409, code: "postponed", dueAfter: today, postponedTo: "2026-10-10"},
		{name: "parent cannot postpone", role: "parent", due: today, used: false, version: 1, expected: 403, code: "forbidden", dueAfter: today, postponedTo: nil},
		{name: "postpone after the due date", role: "student", due: yesterday, used: false, version: 1, expected: 400, code: "invalid_input", dueAfter: yesterday, postponedTo: nil},
		{name: "postpone without a due date", role: "student", due: nil, used: false, version: 1, expected: 400, code: "invalid_input", dueAfter: nil, postponedTo: nil},
		{name: "stale postpone", role: "student", due: today, used: false, version: 2, expected: 409, code: "conflict", dueAfter: today, postponedTo: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			exec(t, f.pool, "UPDATE homework SET due_date=$1", tt.due)
			if tt.used {
				exec(t, f.pool, homeworkInsertPostponement, f.homeworkID)
			}
			w := request(t, f.s, "/homework/"+f.homeworkID+"/postpone", http.MethodPost, map[string]any{"version": tt.version}, f.as(tt.role))
			if w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			want := map[string]any{"dueDate": tt.dueAfter, "postponedTo": tt.postponedTo, "version": 2}
			if tt.code != "" {
				want = map[string]any{"error": map[string]any{"code": tt.code}}
			}
			problemsMatchBody(t, decode[any](t, w), want)
			dates := map[string]any{"dueDate": tt.dueAfter, "postponedTo": tt.postponedTo}
			problemsMatchBody(t, decode[any](t, request(t, f.s, "/homework/"+f.homeworkID, http.MethodGet, nil, f.parent)), map[string]any{"homework": dates, "tasks": []any{dates}})
		})
	}
}

func TestHomeworkList(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		titles []string
	}{
		{name: "submitted homework first, then newest", query: "", titles: []string{"Week 41: backtracking", "Two Sum practice", "Week 42: graphs", "Week 40: arrays"}},
		{name: "active homework keeps a bundle without tasks", query: "active=1", titles: []string{"Week 41: backtracking", "Two Sum practice", "Week 42: graphs"}},
		{name: "title search ignores ASCII case", query: "q=WEEK%204", titles: []string{"Week 41: backtracking", "Week 42: graphs", "Week 40: arrays"}},
		{name: "title search with active homework", query: "q=week&active=1", titles: []string{"Week 41: backtracking", "Week 42: graphs"}},
		{name: "title search without a match", query: "q=heap", titles: []string{}},
		{name: "second page", query: "limit=2&offset=2", titles: []string{"Week 42: graphs", "Week 40: arrays"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			homeworkBundles(t, f)
			w := request(t, f.s, "/homework?"+tt.query, http.MethodGet, nil, f.student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			titles := []string{}
			for _, item := range decode[homeworkPage](t, w).Items {
				title, _ := item["title"].(string)
				titles = append(titles, title)
			}
			if !slices.Equal(titles, tt.titles) {
				t.Errorf("titles = %q, want %q", titles, tt.titles)
			}
		})
	}
}

func TestHomeworkListCountsTasksWithoutCancelled(t *testing.T) {
	f := homeworkSetup(t)
	homeworkBundles(t, f)
	items := decode[homeworkPage](t, request(t, f.s, "/homework", http.MethodGet, nil, f.parent)).Items
	problemsMatchBody(t, homeworkFields(items, "title", "taskCount", "completedCount", "submittedCount", "requestedCount"), []any{
		map[string]any{"title": "Week 41: backtracking", "taskCount": 3, "completedCount": 1, "submittedCount": 1, "requestedCount": 1},
		map[string]any{"title": "Two Sum practice", "taskCount": 1, "completedCount": 0, "submittedCount": 0, "requestedCount": 0},
		map[string]any{"title": "Week 42: graphs", "taskCount": 0, "completedCount": 0, "submittedCount": 0, "requestedCount": 0},
		map[string]any{"title": "Week 40: arrays", "taskCount": 1, "completedCount": 1, "submittedCount": 0, "requestedCount": 0},
	})
}

// The problem and the bundle "first" are the ones of homeworkSetup. "other" is a new problem or bundle, and "unknown" is a new ID.
func TestHomeworkAddTask(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		problem    string
		bundle     string
		firstState string
		accepted   bool
		archived   bool
		expected   int
		error      map[string]any
	}{
		{name: "parent adds a task", role: "parent", problem: "other", bundle: "other", archived: false, expected: 201, error: nil},
		{name: "completed problem in another homework", role: "parent", problem: "first", bundle: "other", firstState: "completed", archived: false, expected: 409, error: map[string]any{"code": "completed_homework"}},
		{name: "completed problem in the same homework", role: "parent", problem: "first", bundle: "first", firstState: "completed", archived: false, expected: 409, error: map[string]any{"code": "completed_homework"}},
		{name: "cancelled unsolved problem can be assigned again", role: "parent", problem: "first", bundle: "other", firstState: "cancelled", archived: false, expected: 201, error: nil},
		{name: "accepted problem without homework", role: "parent", problem: "other", bundle: "other", accepted: true, archived: false, expected: 409, error: map[string]any{"code": "completed_homework"}},
		{name: "student cannot add a task", role: "student", problem: "other", bundle: "other", archived: false, expected: 403, error: map[string]any{"code": "forbidden"}},
		{name: "problem with an active task in another bundle", role: "parent", problem: "first", bundle: "other", archived: false, expected: 409, error: map[string]any{"code": "active_homework", "message": "У этой задачи уже есть активное задание. Откройте его, чтобы продолжить работу."}},
		{name: "problem with an active task in the same bundle", role: "parent", problem: "first", bundle: "first", archived: false, expected: 409, error: map[string]any{"code": "active_homework"}},
		{name: "task for an archived problem", role: "parent", problem: "other", bundle: "other", archived: true, expected: 409, error: map[string]any{"code": "archived", "message": "Верните задачу из архива, прежде чем назначать домашнюю работу."}},
		{name: "task for an unknown problem", role: "parent", problem: "unknown", bundle: "other", archived: false, expected: 404, error: map[string]any{"code": "not_found"}},
		{name: "task in an unknown bundle", role: "parent", problem: "other", bundle: "unknown", archived: false, expected: 404, error: map[string]any{"code": "not_found"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			other, otherBundle := newID(), newID()
			var archivedAt any
			if tt.archived {
				archivedAt = "now"
			}
			exec(t, f.pool, homeworkInsertProblem, other, "subsets", "Subsets")
			exec(t, f.pool, homeworkInsertBundle, otherBundle, "Week 41: backtracking", "2026-10-12", "2026-10-05", "2026-10-05")
			exec(t, f.pool, "UPDATE problems SET archived_at=$1 WHERE id=$2", archivedAt, other)
			if tt.firstState != "" {
				exec(t, f.pool, "UPDATE homework_tasks SET state=$1", tt.firstState)
			}
			if tt.accepted {
				exec(t, f.pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'saved','{}','','accepted','independent','now','now')", newID(), other)
			}
			chosen := map[string]string{"first": f.problemID, "other": other, "unknown": newID()}[tt.problem]
			target := map[string]string{"first": f.homeworkID, "other": otherBundle, "unknown": newID()}[tt.bundle]
			w := request(t, f.s, "/tasks", http.MethodPost, map[string]any{"homeworkId": target, "problemId": chosen}, f.as(tt.role))
			if w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			want, tasks := map[string]any{"error": tt.error}, int64(1)
			if tt.error == nil {
				want, tasks = map[string]any{"homeworkId": otherBundle, "homeworkTitle": "Week 41: backtracking", "instructions": "Explain the invariant", "dueDate": "2026-10-12", "problemId": chosen, "state": "assigned", "submissionId": nil, "version": 1}, 2
			}
			problemsMatchBody(t, decode[any](t, w), want)
			if got := queryValue[int64](t, f.pool, "SELECT count(*) AS count FROM homework_tasks"); got != tasks {
				t.Errorf("task count = %d, want %d", got, tasks)
			}
		})
	}
}

func TestHomeworkTaskTransitions(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		from     string
		action   string
		expected int
		to       string
	}{
		{name: "student starts assignment", role: "student", from: "assigned", action: "start", expected: 200, to: "in_progress"},
		{name: "parent requests changes", role: "parent", from: "submitted", action: "changes_requested", expected: 200, to: "changes_requested"},
		{name: "parent completes review", role: "parent", from: "submitted", action: "completed", expected: 200, to: "completed"},
		{name: "student cannot complete", role: "student", from: "submitted", action: "completed", expected: 403, to: "submitted"},
		{name: "cancel assigned", role: "parent", from: "assigned", action: "cancel", expected: 200, to: "cancelled"},
		{name: "cancel progress", role: "parent", from: "in_progress", action: "cancel", expected: 200, to: "cancelled"},
		{name: "cancel submitted", role: "parent", from: "submitted", action: "cancel", expected: 200, to: "cancelled"},
		{name: "cancel requested changes", role: "parent", from: "changes_requested", action: "cancel", expected: 200, to: "cancelled"},
		{name: "completed homework stays complete", role: "parent", from: "completed", action: "cancel", expected: 409, to: "completed"},
		{name: "parent returns a completed task", role: "parent", from: "completed", action: "changes_requested", expected: 200, to: "changes_requested"},
		{name: "parent cannot complete a completed task again", role: "parent", from: "completed", action: "completed", expected: 409, to: "completed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			var current any
			if tt.from == "submitted" || tt.from == "completed" {
				exec(t, f.pool, "UPDATE attempts SET state='saved'")
				exec(t, f.pool, "INSERT INTO submissions VALUES ($1,$2,$3,$4,$5)", f.submissionID, f.taskID, f.attemptID, 1, "now")
				current = f.submissionID
			}
			exec(t, f.pool, "UPDATE homework_tasks SET state=$1,current_submission_id=$2", tt.from, current)
			path, body := "/tasks/"+f.taskID+"/"+tt.action, map[string]any{"version": 1}
			if tt.action == "changes_requested" || tt.action == "completed" {
				path, body = "/tasks/"+f.taskID+"/review", map[string]any{"version": 1, "submissionId": f.submissionID, "decision": tt.action, "body": "Explain the next step."}
			}
			w := request(t, f.s, path, http.MethodPost, body, f.as(tt.role))
			if w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			if got := homeworkRow(t, f.pool, "SELECT state FROM homework_tasks WHERE id=$1", f.taskID); !reflect.DeepEqual(got, map[string]any{"state": tt.to}) {
				t.Errorf("task = %v, want state %s", got, tt.to)
			}
		})
	}
}

func TestHomeworkReviewRejections(t *testing.T) {
	tests := []struct {
		name     string
		current  bool
		decision string
		body     string
		expected int
	}{
		{name: "blank change request", current: true, decision: "changes_requested", body: "  ", expected: 400},
		{name: "review of another submission", current: false, decision: "completed", body: "", expected: 409},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			submission := homeworkSolve(t, f)
			if !tt.current {
				submission = newID()
			}
			if w := request(t, f.s, "/tasks/"+f.taskID+"/review", http.MethodPost, map[string]any{"version": 2, "submissionId": submission, "decision": tt.decision, "body": tt.body}, f.parent); w.Code != tt.expected {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.expected, w.Body.String())
			}
			if got, want := homeworkRow(t, f.pool, "SELECT state,version FROM homework_tasks"), map[string]any{"state": "submitted", "version": int64(2)}; !reflect.DeepEqual(got, want) {
				t.Errorf("task = %v, want %v", got, want)
			}
			if got := queryValue[int64](t, f.pool, "SELECT count(*) AS count FROM feedback"); got != 0 {
				t.Errorf("feedback count = %d, want 0", got)
			}
		})
	}
}

func TestHomeworkRepeatedReviewDecision(t *testing.T) {
	f := homeworkSetup(t)
	submission := homeworkSolve(t, f)
	review := func() int {
		return request(t, f.s, "/tasks/"+f.taskID+"/review", http.MethodPost, map[string]any{"version": 2, "submissionId": submission, "decision": "changes_requested", "body": "Add a proof."}, f.parent).Code
	}
	if got := review(); got != http.StatusOK {
		t.Errorf("first review: status = %d, want 200", got)
	}
	if got := review(); got != http.StatusConflict {
		t.Errorf("second review: status = %d, want 409", got)
	}
	if got := queryValue[int64](t, f.pool, "SELECT count(*) AS count FROM feedback"); got != 1 {
		t.Errorf("feedback count = %d, want 1", got)
	}
}

func TestHomeworkTasksOfProblemWithSubmission(t *testing.T) {
	f := homeworkSetup(t)
	submission := homeworkSolve(t, f)
	w := request(t, f.s, "/tasks?problemId="+f.problemID, http.MethodGet, nil, f.parent)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	problemsMatchBody(t, decode[any](t, w), map[string]any{"items": []any{map[string]any{"id": f.taskID, "homeworkId": f.homeworkID, "homeworkTitle": "Two Sum practice", "instructions": "Explain the invariant", "dueDate": "2026-10-10", "problemId": f.problemID, "title": "Two Sum", "difficulty": "easy", "state": "submitted", "submissionId": submission, "version": 2}}})
}

func TestHomeworkBundleTasksInCreationOrder(t *testing.T) {
	f := homeworkSetup(t)
	other := newID()
	exec(t, f.pool, homeworkInsertProblem, other, "subsets", "Subsets")
	exec(t, f.pool, homeworkInsertTask, newID(), f.homeworkID, other, "assigned", "2026-10-03")
	w := request(t, f.s, "/homework/"+f.homeworkID, http.MethodGet, nil, f.student)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	problemsMatchBody(t, decode[any](t, w), map[string]any{
		"homework": map[string]any{"id": f.homeworkID, "title": "Two Sum practice"},
		"tasks":    []any{map[string]any{"title": "Subsets", "homeworkId": f.homeworkID}, map[string]any{"id": f.taskID, "title": "Two Sum", "homeworkId": f.homeworkID}},
	})
}

func TestHomeworkUnknown(t *testing.T) {
	f := homeworkSetup(t)
	w := request(t, f.s, "/homework/"+newID(), http.MethodGet, nil, f.student)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
	problemsMatchBody(t, decode[any](t, w), map[string]any{"error": map[string]any{"code": "not_found"}})
}

func TestHomeworkDiscussionReplyRejected(t *testing.T) {
	f := homeworkSetup(t)
	w := request(t, f.s, "/problems/"+f.problemID+"/feedback", http.MethodPost, map[string]any{"taskId": f.taskID, "body": "Follow-up question."}, f.student)
	if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "GET" {
		t.Errorf("status = %d, Allow = %q; want 405 and GET", w.Code, w.Header().Get("Allow"))
	}
	if got := queryValue[int64](t, f.pool, "SELECT count(*) AS count FROM feedback"); got != 0 {
		t.Errorf("feedback count = %d, want 0", got)
	}
}

// The student writes each message. earlier is the read marker of the parent before the request, or "" for none.
func TestHomeworkDiscussionReadMarkers(t *testing.T) {
	tests := []struct {
		name     string
		viewer   string
		offset   int
		messages []string
		earlier  string
		reads    []map[string]any
		unread   bool
	}{
		{name: "first page marks the newest message read", viewer: "parent", offset: 0, messages: []string{"2026-10-05T10:00:00.000Z", "2026-10-05T09:00:00.000Z"}, earlier: "", reads: []map[string]any{{"role": "parent", "read_at": "2026-10-05T10:00:00.000Z"}}, unread: false},
		{name: "first page moves an earlier read marker", viewer: "parent", offset: 0, messages: []string{"2026-10-05T09:00:00.000Z", "2026-10-05T10:00:00.000Z"}, earlier: "2026-10-05T09:00:00.000Z", reads: []map[string]any{{"role": "parent", "read_at": "2026-10-05T10:00:00.000Z"}}, unread: false},
		{name: "offset 50 does not mark read", viewer: "parent", offset: 50, messages: []string{"2026-10-05T10:00:00.000Z"}, earlier: "", reads: []map[string]any{}, unread: true},
		{name: "first page without messages records no read marker", viewer: "parent", offset: 0, messages: []string{}, earlier: "", reads: []map[string]any{}, unread: false},
		{name: "student read keeps the parent read marker", viewer: "student", offset: 0, messages: []string{"2026-10-05T10:00:00.000Z"}, earlier: "2026-10-05T09:00:00.000Z", reads: []map[string]any{{"role": "parent", "read_at": "2026-10-05T09:00:00.000Z"}, {"role": "student", "read_at": "2026-10-05T10:00:00.000Z"}}, unread: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			for _, at := range tt.messages {
				exec(t, f.pool, homeworkInsertFeedback, newID(), f.problemID, "student", "reply", at)
			}
			if tt.earlier != "" {
				exec(t, f.pool, homeworkInsertRead, "parent", f.problemID, tt.earlier)
			}
			if w := request(t, f.s, "/problems/"+f.problemID+"/feedback?offset="+strconv.Itoa(tt.offset), http.MethodGet, nil, f.as(tt.viewer)); w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			if got := problemsTableRows(t, f.pool, "SELECT role,read_at FROM discussion_reads ORDER BY role"); !reflect.DeepEqual(got, tt.reads) {
				t.Errorf("discussion_reads = %v, want %v", got, tt.reads)
			}
			tasks := decode[homeworkPage](t, request(t, f.s, "/tasks", http.MethodGet, nil, f.parent)).Items
			if got := homeworkFirst(t, tasks)["unread"]; got != tt.unread {
				t.Errorf("unread = %v, want %v", got, tt.unread)
			}
		})
	}
}

// section is "" for a task that is not on the board.
func TestHomeworkBoardSections(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		state   string
		section string
	}{
		{name: "assigned task is in progress", role: "student", state: "assigned", section: "progress"},
		{name: "started task is in progress", role: "student", state: "in_progress", section: "progress"},
		{name: "submitted task is in progress", role: "student", state: "submitted", section: "progress"},
		{name: "task with requested changes is in progress", role: "student", state: "changes_requested", section: "progress"},
		{name: "completed task is accepted", role: "student", state: "completed", section: "accepted"},
		{name: "parent sees the same board", role: "parent", state: "completed", section: "accepted"},
		{name: "cancelled task is not on the board", role: "student", state: "cancelled", section: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			exec(t, f.pool, "UPDATE homework_tasks SET state=$1", tt.state)
			w := request(t, f.s, "/tasks?view=board", http.MethodGet, nil, f.as(tt.role))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			task := map[string]any{"id": f.taskID, "homeworkId": f.homeworkID, "homeworkTitle": "Two Sum practice", "dueDate": "2026-10-10", "problemId": f.problemID, "title": "Two Sum", "state": tt.state, "unread": false}
			want := map[string]any{"progress": []any{}, "future": []any{}, "accepted": []any{}}
			if tt.section != "" {
				want[tt.section] = []any{task}
			}
			body := decode[map[string]any](t, w)
			if len(body) != len(want) {
				t.Errorf("board sections = %v, want progress, future and accepted", body)
			}
			problemsMatchBody(t, body, want)
		})
	}
}

func TestHomeworkBoardKeepsEveryTask(t *testing.T) {
	tests := []struct {
		name    string
		state   string
		section string
	}{
		{name: "keeps every assigned task on a board with more than 200 tasks", state: "assigned", section: "progress"},
		{name: "keeps every completed task on a board with more than 200 tasks", state: "completed", section: "accepted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			exec(t, f.pool, "UPDATE homework_tasks SET state=$1 WHERE id=$2", tt.state, f.taskID)
			want := []string{f.taskID}
			for range 200 {
				id, problemID := newID(), newID()
				exec(t, f.pool, homeworkInsertProblem, problemID, problemID, "Board problem")
				exec(t, f.pool, homeworkInsertTask, id, f.homeworkID, problemID, tt.state, "now")
				want = append(want, id)
			}
			w := request(t, f.s, "/tasks?view=board", http.MethodGet, nil, f.student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			ids := []string{}
			for _, task := range decode[map[string][]map[string]any](t, w)[tt.section] {
				id, _ := task["id"].(string)
				ids = append(ids, id)
			}
			slices.Sort(ids)
			slices.Sort(want)
			if !slices.Equal(ids, want) {
				t.Errorf("%s has %d tasks, want the %d tasks of the test", tt.section, len(ids), len(want))
			}
		})
	}
}

// homeworkBoardBundle is a homework of TestHomeworkBoardOrder. due is nil for no due date.
type homeworkBoardBundle struct {
	title string
	due   any
}

// homeworkBoardTask is a task of TestHomeworkBoardOrder. Its problem has the title of the task. at is the time of creation and of the last change.
type homeworkBoardTask struct {
	title, homework, state, at string
}

// The active homework with the earliest due date is in progress. The other active homework is in the future, by due date with no date last.
func TestHomeworkBoardOrder(t *testing.T) {
	tests := []struct {
		name     string
		homework []homeworkBoardBundle
		tasks    []homeworkBoardTask
		board    map[string][]string
	}{
		{
			name:     "earliest homework in progress, the other active homework in the future",
			homework: []homeworkBoardBundle{{title: "No date", due: nil}, {title: "Later", due: "2026-10-12"}, {title: "Sooner", due: "2026-10-11"}},
			tasks:    []homeworkBoardTask{{title: "No date", homework: "No date", state: "in_progress", at: "now"}, {title: "Later", homework: "Later", state: "submitted", at: "now"}, {title: "Sooner", homework: "Sooner", state: "assigned", at: "now"}},
			board:    map[string][]string{"progress": {"Sooner"}, "future": {"Later", "No date"}, "accepted": {}},
		},
		{
			name:     "homework without a due date in progress when it is the only active homework",
			homework: []homeworkBoardBundle{{title: "No date", due: nil}, {title: "Done", due: "2026-10-01"}},
			tasks:    []homeworkBoardTask{{title: "No date", homework: "No date", state: "assigned", at: "now"}, {title: "Done", homework: "Done", state: "completed", at: "now"}},
			board:    map[string][]string{"progress": {"No date"}, "future": {}, "accepted": {"Done"}},
		},
		{
			name:     "tasks of the current homework in creation order",
			homework: []homeworkBoardBundle{{title: "Week", due: "2026-10-11"}},
			tasks:    []homeworkBoardTask{{title: "Second", homework: "Week", state: "assigned", at: "2026-10-02T10:00:00.000Z"}, {title: "First", homework: "Week", state: "in_progress", at: "2026-10-01T10:00:00.000Z"}},
			board:    map[string][]string{"progress": {"First", "Second"}, "future": {}, "accepted": {}},
		},
		{
			name:     "accepted by completion time with the newest first",
			homework: []homeworkBoardBundle{{title: "Week", due: nil}, {title: "Middle", due: "2026-10-01"}},
			tasks:    []homeworkBoardTask{{title: "First", homework: "Week", state: "completed", at: "2026-10-05T10:00:00.000Z"}, {title: "Last", homework: "Week", state: "completed", at: "2026-10-07T10:00:00.000Z"}, {title: "Middle", homework: "Middle", state: "completed", at: "2026-10-06T10:00:00.000Z"}},
			board:    map[string][]string{"progress": {}, "future": {}, "accepted": {"Last", "Middle", "First"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			exec(t, f.pool, "UPDATE homework_tasks SET state='cancelled'")
			bundles := map[string]string{}
			for _, h := range tt.homework {
				bundles[h.title] = newID()
				exec(t, f.pool, homeworkInsertBundle, bundles[h.title], h.title, h.due, "now", "now")
			}
			for _, task := range tt.tasks {
				id := newID()
				exec(t, f.pool, homeworkInsertProblem, id, task.title, task.title)
				exec(t, f.pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)", newID(), bundles[task.homework], id, task.state, task.at, task.at)
			}
			board := map[string][]string{}
			for section, rows := range decode[map[string][]map[string]any](t, request(t, f.s, "/tasks?view=board", http.MethodGet, nil, f.student)) {
				titles := []string{}
				for _, row := range rows {
					title, _ := row["title"].(string)
					titles = append(titles, title)
				}
				board[section] = titles
			}
			if !reflect.DeepEqual(board, tt.board) {
				t.Errorf("board = %q, want %q", board, tt.board)
			}
		})
	}
}

func TestHomeworkBundleUnreadWithUnreadTask(t *testing.T) {
	f := homeworkSetup(t)
	homeworkBundles(t, f)
	exec(t, f.pool, "INSERT INTO feedback (id,problem_id,author,kind,body,created_at) SELECT $1,id,'student','reply','Question','2026-10-05T10:00:00.000Z' FROM problems WHERE slug='subsets'", newID())
	items := decode[homeworkPage](t, request(t, f.s, "/homework", http.MethodGet, nil, f.parent)).Items
	problemsMatchBody(t, homeworkFields(items, "title", "unread"), []any{
		map[string]any{"title": "Week 41: backtracking", "unread": true},
		map[string]any{"title": "Two Sum practice", "unread": false},
		map[string]any{"title": "Week 42: graphs", "unread": false},
		map[string]any{"title": "Week 40: arrays", "unread": false},
	})
}

func TestHomeworkAttemptOfCurrentSubmission(t *testing.T) {
	f := homeworkSetup(t)
	path := "/tasks?problemId=" + f.problemID
	before := decode[homeworkPage](t, request(t, f.s, path, http.MethodGet, nil, f.parent)).Items
	problemsMatchBody(t, homeworkFirst(t, before), map[string]any{"submissionId": nil, "submissionAttemptId": nil})
	exec(t, f.pool, "INSERT INTO submissions (id,task_id,attempt_id,task_version,created_at) VALUES ($1,$2,$3,2,$4)", f.submissionID, f.taskID, f.attemptID, time.Now().UTC().Format(time.DateOnly))
	exec(t, f.pool, "UPDATE homework_tasks SET state='submitted',current_submission_id=$1,version=2 WHERE id=$2", f.submissionID, f.taskID)
	after := decode[homeworkPage](t, request(t, f.s, path, http.MethodGet, nil, f.parent)).Items
	problemsMatchBody(t, homeworkFirst(t, after), map[string]any{"submissionId": f.submissionID, "submissionAttemptId": f.attemptID})
}

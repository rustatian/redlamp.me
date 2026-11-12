package notebook

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

var regexpV4 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

const raceProblem = `{"url":"https://leetcode.com/problems/two-sum/","number":1,"title":"Two Sum","difficulty":"easy","topics":[],"summary":""}`

// Each case sends two requests that are in progress at the same time. The rule of the case must hold for the two together.
func TestConcurrentRequestsOfProblemsAndAttempts(t *testing.T) {
	tests := []struct {
		name string
		// setup prepares the data and returns the write that both requests wait for, and the request.
		setup        func(t *testing.T, pool *pgxpool.Pool, s *Server, cookie string) (hold string, args []any, send func() *httptest.ResponseRecorder)
		wantStatuses []int
		count        string
		wantCount    int64
	}{
		{
			name: "one_draft_per_problem gives one draft to two requests that start an attempt",
			setup: func(t *testing.T, _ *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				id := addProblem(t, s, cookie, raceProblem)
				return "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ('held',$1,'draft','{}','','not_submitted','needs_practice','t0','t0')", []any{id},
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodPost, "/leetcode/api/problems/"+id+"/attempts", `{}`, cookie)
					}
			},
			wantStatuses: []int{201, 409},
			count:        "SELECT count(*) FROM attempts WHERE state='draft'",
			wantCount:    1,
		},
		{
			name: "the unique slug gives one problem to two requests that add the same problem",
			setup: func(_ *testing.T, _ *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				return "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ('held','two-sum','Held','easy','[]','','','t0','t0')", nil,
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodPost, "/leetcode/api/problems", raceProblem, cookie)
					}
			},
			wantStatuses: []int{200, 201},
			count:        "SELECT count(*) FROM problems",
			wantCount:    1,
		},
		{
			name: "a problem with 19 own tests gets one more from two requests",
			setup: func(t *testing.T, pool *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				id := addProblem(t, s, cookie, raceProblem)
				exec(t, pool, "INSERT INTO problem_tests (id,problem_id,input,created_at) SELECT 'test-' || n,$1,'[1]','t0' FROM generate_series(1,19) AS n", id)
				return "UPDATE problems SET version=version WHERE id=$1", []any{id},
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodPost, "/leetcode/api/problems/"+id+"/tests", `{"input":"[2]"}`, cookie)
					}
			},
			wantStatuses: []int{201, 409},
			count:        "SELECT count(*) FROM problem_tests",
			wantCount:    20,
		},
		{
			name: "a problem gets one recall from two requests with one progress version",
			setup: func(t *testing.T, _ *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				id := addProblem(t, s, cookie, raceProblem)
				return "UPDATE problems SET version=version WHERE id=$1", []any{id},
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodPost, "/leetcode/api/problems/"+id+"/reviews", `{"version":1,"result":"with_help","note":"","reviewedOn":"2026-10-07","nextReviewDate":null}`, cookie)
					}
			},
			wantStatuses: []int{201, 409},
			count:        "SELECT count(*) FROM reviews",
			wantCount:    1,
		},
		{
			name: "two reads of the game on one day record one visit",
			setup: func(_ *testing.T, _ *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				return "INSERT INTO visits (day) VALUES ($1)", []any{time.Now().UTC().Format(time.DateOnly)},
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodGet, "/leetcode/api/game", "", cookie)
					}
			},
			wantStatuses: []int{200, 200},
			count:        "SELECT count(*) FROM visits",
			wantCount:    1,
		},
		{
			name: "a problem gets one edit from two requests with one version",
			setup: func(t *testing.T, _ *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				id := addProblem(t, s, cookie, raceProblem)
				return "UPDATE problems SET version=version WHERE id=$1", []any{id},
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodPatch, "/leetcode/api/problems/"+id, `{"url":"https://leetcode.com/problems/two-sum/","number":1,"title":"Updated","difficulty":"easy","topics":[],"summary":"","version":1}`, cookie)
					}
			},
			wantStatuses: []int{200, 409},
			count:        "SELECT version FROM problems",
			wantCount:    2,
		},
		{
			name: "a draft gets one update from two requests with one version",
			setup: func(t *testing.T, _ *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				id := addProblem(t, s, cookie, raceProblem)
				attempt := decode[Attempt](t, request(t, s, "/problems/"+id+"/attempts", http.MethodPost, `{}`, cookie)).ID
				return "UPDATE attempts SET version=version WHERE id=$1", []any{attempt},
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodPut, "/leetcode/api/problems/"+id+"/attempts/"+attempt, `{"version":1,"document":`+attemptsStored+`}`, cookie)
					}
			},
			wantStatuses: []int{200, 409},
			count:        "SELECT version FROM attempts",
			wantCount:    2,
		},
		{
			name: "two reads of one discussion keep one read mark",
			setup: func(t *testing.T, pool *pgxpool.Pool, s *Server, cookie string) (string, []any, func() *httptest.ResponseRecorder) {
				id := addProblem(t, s, cookie, raceProblem)
				exec(t, pool, homeworkInsertFeedback, newID(), id, "parent", "reply", "2026-10-05T10:00:00.000Z")
				return homeworkInsertRead, []any{"student", id, "t0"},
					func() *httptest.ResponseRecorder {
						return call(s, http.MethodGet, "/leetcode/api/problems/"+id+"/feedback", "", cookie)
					}
			},
			wantStatuses: []int{200, 200},
			count:        "SELECT count(*) FROM discussion_reads",
			wantCount:    1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServer(t, pool)
			cookie := loginAs(t, s, "student")
			hold, args, send := tt.setup(t, pool, s, cookie)
			results := race(t, pool, hold, args, send, send)
			statuses := []int{results[0].Code, results[1].Code}
			slices.Sort(statuses)
			if !slices.Equal(statuses, tt.wantStatuses) {
				t.Errorf("statuses = %v, want %v; bodies %s and %s", statuses, tt.wantStatuses, results[0].Body.String(), results[1].Body.String())
			}
			if got := queryValue[int64](t, pool, tt.count); got != tt.wantCount {
				t.Errorf("%s = %d, want %d", tt.count, got, tt.wantCount)
			}
		})
	}
}

// A UUID of version 4 has the version digit 4 and a variant digit of 8, 9, a or b: https://www.rfc-editor.org/rfc/rfc9562#name-uuid-version-4
func TestNewID(t *testing.T) {
	first, second := newID(), newID()
	if !regexpV4.MatchString(first) || !regexpV4.MatchString(second) {
		t.Errorf("newID = %q and %q, want UUIDs of version 4", first, second)
	}
	if first == second {
		t.Errorf("two calls of newID gave %q", first)
	}
	if _, err := uuid(first); err != nil {
		t.Errorf("the identifier check rejects %q: %v", first, err)
	}
}

// The task of the second chapter is the problem of raceProblem.
var raceSets = []sets.LearningSet{{Slug: "race-set", Title: "Race set", Chapters: []sets.Chapter{
	{Title: "First chapter", Cards: []sets.Card{{ID: "race-card", Type: "choice", Prompt: "?", Options: []string{"a", "b"}, Answer: 0, Explain: "a"}}},
	{Title: "Second chapter", Tasks: []sets.Task{{URL: "https://leetcode.com/problems/two-sum/", Number: new(int64(1)), Title: "Two Sum", Difficulty: "easy"}}},
}}}

// raceActors are the data of one case of TestConcurrentRequestsOfHomework.
type raceActors struct {
	pool      *pgxpool.Pool
	s         *Server
	student   string
	parent    string
	problemID string
	homework  string
}

// Each case sends two requests that are in progress at the same time. The rule of the case must hold for the two together.
func TestConcurrentRequestsOfHomework(t *testing.T) {
	send := func(a raceActors, method, path, body, cookie string) func() *httptest.ResponseRecorder {
		return func() *httptest.ResponseRecorder {
			return call(a.s, method, "/leetcode/api"+path, body, cookie)
		}
	}
	post := func(a raceActors, path, body, cookie string) func() *httptest.ResponseRecorder {
		return send(a, http.MethodPost, path, body, cookie)
	}
	tests := []struct {
		name string
		// setup prepares the data and returns the write that both requests wait for, and the two requests.
		setup         func(t *testing.T, a raceActors) (hold string, args []any, first, second func() *httptest.ResponseRecorder)
		wantConflicts int
		// allowedConflicts replaces wantConflicts when each order of the two requests is legal. It lists the legal numbers of 409 answers.
		allowedConflicts []int
		// wantPaused is the number of requests that get 429, the answer of a quiz deck in its pause.
		wantPaused int
		count      string
		wantCount  int64
	}{
		{
			name: "a problem gets one active homework task from two requests that assign it",
			setup: func(_ *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				assign := post(a, "/tasks", `{"homeworkId":"`+a.homework+`","problemId":"`+a.problemID+`"}`, a.parent)
				return "UPDATE problems SET version=version WHERE id=$1", []any{a.problemID}, assign, assign
			},
			wantConflicts: 1,
			count:         "SELECT count(*) FROM homework_tasks WHERE state='assigned'",
			wantCount:     1,
		},
		{
			name: "an archive and a new homework task of one problem do not both succeed",
			setup: func(_ *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				return "UPDATE problems SET version=version WHERE id=$1", []any{a.problemID},
					post(a, "/problems/"+a.problemID+"/archive", `{"version":1}`, a.parent),
					post(a, "/tasks", `{"homeworkId":"`+a.homework+`","problemId":"`+a.problemID+`"}`, a.parent)
			},
			wantConflicts: 1,
			count:         "SELECT count(*) FROM problems p WHERE p.archived_at IS NOT NULL AND EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state='assigned')",
			wantCount:     0,
		},
		{
			name: "an archive and a return of a completed task for changes do not both succeed",
			setup: func(t *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				attempt := decode[Attempt](t, request(t, a.s, "/problems/"+a.problemID+"/attempts", http.MethodPost, `{}`, a.student)).ID
				task, submission := newID(), newID()
				exec(t, a.pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,version,created_at,updated_at) VALUES ($1,$2,$3,'completed',3,'t0','t0')", task, a.homework, a.problemID)
				exec(t, a.pool, "INSERT INTO submissions VALUES ($1,$2,$3,2,'t0')", submission, task, attempt)
				exec(t, a.pool, "UPDATE homework_tasks SET current_submission_id=$1 WHERE id=$2", submission, task)
				return "UPDATE problems SET version=version WHERE id=$1", []any{a.problemID},
					post(a, "/problems/"+a.problemID+"/archive", `{"version":1}`, a.parent),
					post(a, "/tasks/"+task+"/review", `{"version":3,"submissionId":"`+submission+`","decision":"changes_requested","body":"Fix the loop."}`, a.parent)
			},
			wantConflicts: 1,
			count:         "SELECT count(*) FROM problems p WHERE p.archived_at IS NOT NULL AND EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state='changes_requested')",
			wantCount:     0,
		},
		{
			name: "a submission gets one review from two requests of the parent",
			setup: func(t *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				attempt := decode[Attempt](t, request(t, a.s, "/problems/"+a.problemID+"/attempts", http.MethodPost, `{}`, a.student)).ID
				task, submission := newID(), newID()
				exec(t, a.pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,version,created_at,updated_at) VALUES ($1,$2,$3,'submitted',2,'t0','t0')", task, a.homework, a.problemID)
				exec(t, a.pool, "INSERT INTO submissions VALUES ($1,$2,$3,2,'t0')", submission, task, attempt)
				exec(t, a.pool, "UPDATE homework_tasks SET current_submission_id=$1 WHERE id=$2", submission, task)
				review := post(a, "/tasks/"+task+"/review", `{"version":2,"submissionId":"`+submission+`","decision":"completed","body":""}`, a.parent)
				return "UPDATE homework_tasks SET version=version WHERE id=$1", []any{task}, review, review
			},
			wantConflicts: 1,
			count:         "SELECT count(*) FROM feedback",
			wantCount:     1,
		},
		{
			name: "a homework gets one postponement from two requests of the student",
			setup: func(t *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				exec(t, a.pool, "UPDATE homework SET due_date=$1 WHERE id=$2", time.Now().UTC().AddDate(0, 0, 1).Format(time.DateOnly), a.homework)
				postpone := post(a, "/homework/"+a.homework+"/postpone", `{"version":1}`, a.student)
				return "UPDATE homework SET version=version WHERE id=$1", []any{a.homework}, postpone, postpone
			},
			wantConflicts: 1,
			count:         "SELECT count(*) FROM homework h JOIN postponements p ON p.homework_id=h.id WHERE h.due_date=p.to_date AND h.due_date=$1",
			wantCount:     1,
		},
		{
			name: "a quiz deck gets one attempt from two submits of the student",
			setup: func(_ *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				submit := post(a, "/quiz/race-set-1/attempts", `{"answers":{"race-card":0}}`, a.student)
				return "INSERT INTO quiz_decks (deck,last_attempt_at) VALUES ('race-set-1','t0')", nil, submit, submit
			},
			wantConflicts: 0,
			wantPaused:    1,
			count:         "SELECT count(*) FROM quiz_attempts",
			wantCount:     1,
		},
		{
			name: "one_active_content gives one active chapter task to two requests that assign one chapter",
			setup: func(_ *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				assign := post(a, "/homework/"+a.homework+"/content", `{"kind":"chapter","setSlug":"race-set","chapter":1}`, a.parent)
				return "INSERT INTO homework_tasks (id,homework_id,kind,set_slug,chapter,state,created_at,updated_at) VALUES ('held',$1,'chapter','race-set',1,'assigned','t0','t0')", []any{a.homework}, assign, assign
			},
			wantConflicts: 0,
			count:         "SELECT count(*) FROM homework_tasks WHERE kind='chapter' AND state='assigned'",
			wantCount:     1,
		},
		{
			name: "the unique slug gives a new slug to one of two problems",
			setup: func(t *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				other := addProblem(t, a.s, a.student, `{"url":"https://leetcode.com/problems/three-sum/","number":15,"title":"3Sum","difficulty":"medium","topics":[],"summary":""}`)
				edit := func(id string) func() *httptest.ResponseRecorder {
					return send(a, http.MethodPatch, "/problems/"+id, `{"url":"https://leetcode.com/problems/move-zeroes/","number":283,"title":"Move Zeroes","difficulty":"easy","topics":[],"summary":"","version":1}`, a.student)
				}
				return "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ('held','move-zeroes','Held','easy','[]','','','t0','t0')", nil, edit(a.problemID), edit(other)
			},
			wantConflicts: 1,
			count:         "SELECT count(*) FROM problems WHERE slug='move-zeroes'",
			wantCount:     1,
		},
		{
			name: "an archive and a chapter assignment of one problem do not both succeed",
			setup: func(_ *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				return "UPDATE problems SET version=version WHERE id=$1", []any{a.problemID},
					post(a, "/problems/"+a.problemID+"/archive", `{"version":1}`, a.parent),
					post(a, "/homework/"+a.homework+"/content", `{"kind":"chapter","setSlug":"race-set","chapter":2}`, a.parent)
			},
			// An archive that comes first makes the assignment skip the problem with 201. An assignment that comes first makes the archive fail with 409.
			allowedConflicts: []int{0, 1},
			count:            "SELECT count(*) FROM problems p WHERE p.archived_at IS NOT NULL AND EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state NOT IN ('completed','cancelled'))",
			wantCount:        0,
		},
		{
			name: "a homework gets one edit from two requests with one version",
			setup: func(_ *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				edit := send(a, http.MethodPatch, "/homework/"+a.homework, `{"version":1,"title":"Week 42","instructions":"","dueDate":null}`, a.parent)
				return "UPDATE homework SET version=version WHERE id=$1", []any{a.homework}, edit, edit
			},
			wantConflicts: 1,
			count:         "SELECT version FROM homework",
			wantCount:     2,
		},
		{
			name: "a task gets one start from two requests with one version",
			setup: func(t *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				task := decode[HomeworkTask](t, request(t, a.s, "/tasks", http.MethodPost, `{"homeworkId":"`+a.homework+`","problemId":"`+a.problemID+`"}`, a.parent)).ID
				start := post(a, "/tasks/"+task+"/start", `{"version":1}`, a.student)
				return "UPDATE homework_tasks SET version=version WHERE id=$1", []any{task}, start, start
			},
			wantConflicts: 1,
			count:         "SELECT version FROM homework_tasks",
			wantCount:     2,
		},
		{
			name: "a homework gets one set link from two requests with one version",
			setup: func(_ *testing.T, a raceActors) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				link := send(a, http.MethodPut, "/homework/"+a.homework+"/sets/race-set", `{"version":1}`, a.parent)
				return "UPDATE homework SET version=version WHERE id=$1", []any{a.homework}, link, link
			},
			wantConflicts: 1,
			count:         "SELECT version FROM homework",
			wantCount:     2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithSets(t, pool, raceSets)
			a := raceActors{pool: pool, s: s, student: loginAs(t, s, "student"), parent: loginAs(t, s, "parent")}
			a.problemID = addProblem(t, s, a.student, raceProblem)
			a.homework = decode[Homework](t, request(t, s, "/homework", http.MethodPost, `{"title":"Week 41","instructions":"","dueDate":null}`, a.parent)).ID
			hold, args, first, second := tt.setup(t, a)
			results := race(t, pool, hold, args, first, second)
			conflicts, paused := 0, 0
			for _, w := range results {
				switch w.Code {
				case http.StatusOK, http.StatusCreated, http.StatusNoContent:
				case http.StatusConflict:
					conflicts++
				case http.StatusTooManyRequests:
					paused++
				default:
					t.Errorf("status = %d, body %s; want a success, 409 or 429", w.Code, w.Body.String())
				}
			}
			wantConflicts := []int{tt.wantConflicts}
			if tt.allowedConflicts != nil {
				wantConflicts = tt.allowedConflicts
			}
			if !slices.Contains(wantConflicts, conflicts) || paused != tt.wantPaused {
				t.Errorf("%d requests got 409 and %d got 429, want one of %v and %d; bodies %s and %s", conflicts, paused, wantConflicts, tt.wantPaused, results[0].Body.String(), results[1].Body.String())
			}
			var countArgs []any
			if strings.Contains(tt.count, "$1") {
				// The due date after one postponement: tomorrow and one day.
				countArgs = []any{time.Now().UTC().AddDate(0, 0, 2).Format(time.DateOnly)}
			}
			if got := queryValue[int64](t, pool, tt.count, countArgs...); got != tt.wantCount {
				t.Errorf("%s = %d, want %d", tt.count, got, tt.wantCount)
			}
		})
	}
}

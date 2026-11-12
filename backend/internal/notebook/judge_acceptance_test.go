package notebook

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// acceptanceContent is Two Sum II with one hidden test. The reference answers with its argument, so the right output of the test is 7.
func acceptanceContent(t *testing.T) Config {
	t.Helper()
	return Config{
		Sets:       goldenSets(t),
		References: map[string]sets.Reference{judgeTwoSumII: {Arity: 1, Call: func(args []any) any { return args[0] }}},
		Generators: map[string]sets.Generator{judgeTwoSumII: func(func() float64) []string { return []string{"[7]"} }},
	}
}

var (
	acceptanceRight = `[{"value":"7"}]`
	acceptanceWrong = `[{"value":"8"}]`
)

// acceptancePost returns a judge request of the student for a seed, with the draft version 1.
func acceptancePost(f judgeFixture, seed, outputs string) func() *httptest.ResponseRecorder {
	return func() *httptest.ResponseRecorder {
		return call(f.s, http.MethodPost, "/leetcode/api/problems/"+f.problemID+"/judge", `{"seed":`+seed+`,"attemptId":"`+f.attemptID+`","version":1,"outputs":`+outputs+`}`, f.student)
	}
}

// Each case sends two requests that are in progress at the same time. The rule of the case must hold for the two together.
func TestJudgeConcurrentRequests(t *testing.T) {
	tests := []struct {
		name string
		// setup prepares the data and returns the write that both requests wait for, and the two requests.
		setup func(t *testing.T, f judgeFixture) (hold string, args []any, first, second func() *httptest.ResponseRecorder)
		// wantStatuses are the pairs of statuses that keep the rule, each in ascending order.
		wantStatuses [][]int
		// wantCode is the error code of a request that gets 409.
		wantCode  string
		count     string
		wantCount int64
	}{
		{
			name: "a seed gives one verdict to two requests with wrong outputs",
			setup: func(_ *testing.T, f judgeFixture) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				send := acceptancePost(f, "1", acceptanceWrong)
				return "UPDATE judge_runs SET created_at=created_at WHERE problem_id=$1", []any{f.problemID}, send, send
			},
			wantStatuses: [][]int{{200, 409}},
			wantCode:     "stale_seed",
			count:        "SELECT count(*) FROM attempts WHERE state='draft'",
			wantCount:    1,
		},
		{
			name: "a draft gets one acceptance from two requests with two seeds",
			setup: func(t *testing.T, f judgeFixture) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				exec(t, f.pool, "INSERT INTO judge_runs VALUES ($1,2,'now')", f.problemID)
				problemsAddTask(t, f.pool, f.problemID, "assigned")
				return "UPDATE attempts SET version=version WHERE id=$1", []any{f.attemptID}, acceptancePost(f, "1", acceptanceRight), acceptancePost(f, "2", acceptanceRight)
			},
			wantStatuses: [][]int{{200, 409}},
			wantCode:     "conflict",
			count:        "SELECT count(*) FROM submissions",
			wantCount:    1,
		},
		{
			// The acceptance completes a task that was assigned before it, and a solved problem gets no new task.
			name: "an acceptance and a new homework task leave no active task on a solved problem",
			setup: func(t *testing.T, f judgeFixture) (string, []any, func() *httptest.ResponseRecorder, func() *httptest.ResponseRecorder) {
				parent := loginAs(t, f.s, "parent")
				homework := decode[Homework](t, request(t, f.s, "/homework", http.MethodPost, `{"title":"Week 41","instructions":"","dueDate":null}`, parent)).ID
				return "UPDATE problems SET version=version WHERE id=$1", []any{f.problemID}, acceptancePost(f, "1", acceptanceRight),
					func() *httptest.ResponseRecorder {
						return call(f.s, http.MethodPost, "/leetcode/api/tasks", `{"homeworkId":"`+homework+`","problemId":"`+f.problemID+`"}`, parent)
					}
			},
			wantStatuses: [][]int{{200, 201}, {200, 409}},
			wantCode:     "completed_homework",
			count:        "SELECT count(*) FROM homework_tasks WHERE state NOT IN ('completed','cancelled')",
			wantCount:    0,
		},
	}
	content := acceptanceContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			hold, args, first, second := tt.setup(t, f)
			results := race(t, f.pool, hold, args, first, second)
			statuses := []int{results[0].Code, results[1].Code}
			slices.Sort(statuses)
			if !slices.ContainsFunc(tt.wantStatuses, func(want []int) bool { return slices.Equal(want, statuses) }) {
				t.Errorf("statuses = %v, want one of %v; bodies %s and %s", statuses, tt.wantStatuses, results[0].Body.String(), results[1].Body.String())
			}
			for _, w := range results {
				if w.Code == http.StatusConflict && errorCode(w) != tt.wantCode {
					t.Errorf("error code = %q, want %q", errorCode(w), tt.wantCode)
				}
			}
			if got := queryValue[int64](t, f.pool, tt.count); got != tt.wantCount {
				t.Errorf("%s = %d, want %d", tt.count, got, tt.wantCount)
			}
		})
	}
}

// A failure of the homework task change leaves the draft, the problem and the task as they were.
// The same request without the failure changes them together.
func TestJudgeAcceptanceIsAtomic(t *testing.T) {
	const state = `SELECT a.state AS attempt,a.acceptance,p.progress_version AS progress,t.state AS task,(SELECT count(*) FROM submissions) AS submissions,(SELECT count(*) FROM feedback) AS feedback FROM attempts a JOIN problems p ON p.id=a.problem_id JOIN homework_tasks t ON t.problem_id=p.id`
	tests := []struct {
		name string
		// judged is true for a problem with hidden tests: the request is an acceptance check. For false, the request is the solved route.
		judged       bool
		wantTask     string
		wantFeedback int
	}{
		{name: "acceptance check that completes the task", judged: true, wantTask: "completed", wantFeedback: 1},
		{name: "solved route that submits the task", judged: false, wantTask: "submitted", wantFeedback: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := acceptanceContent(t)
			if !tt.judged {
				content.Generators = nil
			}
			f := judgeSetup(t, content)
			exec(t, f.pool, judgeWithNotes, "Two pointers.")
			exec(t, f.pool, "INSERT INTO judge_runs VALUES ($1,2,'now')", f.problemID)
			problemsAddTask(t, f.pool, f.problemID, "assigned")
			// The first acceptance check uses its seed also when the transaction fails, so the second one has its own seed.
			send := func(seed string) *httptest.ResponseRecorder {
				if tt.judged {
					return acceptancePost(f, seed, acceptanceRight)()
				}
				return request(t, f.s, "/problems/"+f.problemID+"/solved", http.MethodPost, map[string]any{"attemptId": f.attemptID}, f.student)
			}
			exec(t, f.pool, "CREATE FUNCTION fail_task() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'the task change failed'; END $$")
			exec(t, f.pool, "CREATE TRIGGER fail_task BEFORE UPDATE ON homework_tasks FOR EACH ROW EXECUTE FUNCTION fail_task()")
			// The router answers an error of the database with 503, as the Worker does.
			if w := send("1"); w.Code != http.StatusServiceUnavailable {
				t.Fatalf("status with the failure = %d, want 503; body %s", w.Code, w.Body.String())
			}
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, state), map[string]any{"attempt": "draft", "acceptance": "not_submitted", "progress": 1, "task": "assigned", "submissions": 0, "feedback": 0})
			exec(t, f.pool, "DROP TRIGGER fail_task ON homework_tasks")
			if w := send("2"); w.Code != http.StatusOK {
				t.Fatalf("status without the failure = %d, want 200; body %s", w.Code, w.Body.String())
			}
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, state), map[string]any{"attempt": "saved", "acceptance": "accepted", "progress": 2, "task": tt.wantTask, "submissions": 1, "feedback": tt.wantFeedback})
		})
	}
}

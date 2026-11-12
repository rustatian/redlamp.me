package notebook

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
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

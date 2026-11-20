package notebook

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/bank"
)

// The expected values follow the rules of the contests: easy 800, medium 900 to 1100, hard 1200;
// places easy 1, medium 2, hard 1; points easy 3, medium 4, hard 6; a contest lasts 90 minutes;
// 5 minutes for each wrong submission before the first accepted one of a solved problem;
// the rating starts at 800 and changes by round(32 × Σ(S − E)) with E = 1 / (1 + 10^((R_problem − R) / 400)).

// contestNow is the clock of the contest tests.
var contestNow = time.Date(2026, 11, 20, 12, 0, 0, 0, time.UTC)

// contestBank is the bank of the contest tests: two problems at 800, two at 1000, one at 1100 and one at 1200.
var contestBank = []bank.Problem{
	{Slug: "cf-101-a", Title: "Easy One", Rating: 800, Topics: []string{"Math", "Brute Force"}},
	{Slug: "cf-102-a", Title: "Easy Two", Rating: 800, Topics: []string{"String"}},
	{Slug: "cf-103-b", Title: "Medium One", Rating: 1000, Topics: []string{"Greedy"}},
	{Slug: "cf-104-b", Title: "Medium Two", Rating: 1000, Topics: []string{"Sorting", "Greedy"}},
	{Slug: "cf-105-c", Title: "Medium Three", Rating: 1100, Topics: []string{"Math"}},
	{Slug: "cf-106-d", Title: "Hard One", Rating: 1200, Topics: []string{"Dynamic Programming", "Math"}},
}

// contestSeed is a contest that a test writes with SQL. start is the time from contestNow to the start of the contest.
type contestSeed struct {
	start    time.Duration
	problems []problemSeed
}

// problemSeed is a problem of a seeded contest. An empty slug gives a slug outside the bank.
// accepted and wrong are the times of the submissions from the start of the contest.
type problemSeed struct {
	slug       string
	rating     int
	difficulty string
	points     int
	accepted   []time.Duration
	wrong      []time.Duration
}

// contestID is the ID of the seeded contest at index i.
func contestID(i int) string { return fmt.Sprintf("00000000-0000-4000-8000-%012d", i+1) }

// seedContests writes the contests with their problems and submissions. Each problem has the title "Problem <position>" and the topic Math.
func seedContests(t *testing.T, pool *pgxpool.Pool, contests ...contestSeed) {
	t.Helper()
	for i, c := range contests {
		start := contestNow.Add(c.start)
		exec(t, pool, "INSERT INTO contests VALUES ($1,$2,$3)", contestID(i), start.Format(timeLayout), start.Add(90*time.Minute).Format(timeLayout))
		for j, p := range c.problems {
			slug := p.slug
			if slug == "" {
				slug = fmt.Sprintf("seed-%d-%d", i+1, j+1)
			}
			exec(t, pool, "INSERT INTO contest_problems (contest_id,position,slug,title,rating,topic,difficulty,points) VALUES ($1,$2,$3,$4,$5,'Math',$6,$7)",
				contestID(i), j+1, slug, fmt.Sprintf("Problem %d", j+1), p.rating, p.difficulty, p.points)
			for verdict, times := range map[string][]time.Duration{"accepted": p.accepted, "wrong": p.wrong} {
				for _, at := range times {
					exec(t, pool, "INSERT INTO contest_submissions (id,contest_id,position,verdict,code,created_at) VALUES ($1,$2,$3,$4,'',$5)",
						newID(), contestID(i), j+1, verdict, start.Add(at).Format(timeLayout))
				}
			}
		}
	}
}

// contestServer returns an empty database and a server with contestBank and the clock contestNow.
func contestServer(t *testing.T) (*pgxpool.Pool, *Server) {
	t.Helper()
	pool := database(t)
	s := newServerWithBank(t, pool, contestBank)
	s.now = func() time.Time { return contestNow }
	return pool, s
}

// contestStart sends a start with the problems, each a map with the slug and the topic.
func contestStart(t *testing.T, s *Server, cookie string, problems ...map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	if problems == nil {
		problems = []map[string]string{}
	}
	return request(t, s, "/contests", http.MethodPost, map[string]any{"problems": problems}, cookie)
}

func TestContestStart(t *testing.T) {
	type pick = map[string]string
	tests := []struct {
		name     string
		role     string
		contests []contestSeed
		// library is the slug of a problem in the library, or "".
		library      string
		problems     []pick
		wantStatus   int
		wantCode     string
		wantContests int64
		wantBody     map[string]any
		wantRows     []map[string]any
	}{
		{
			name:         "student starts with one of each place",
			role:         "student",
			problems:     []pick{{"slug": "cf-101-a", "topic": "Brute Force"}, {"slug": "cf-103-b", "topic": "Greedy"}, {"slug": "cf-105-c", "topic": "Math"}, {"slug": "cf-106-d", "topic": "Dynamic Programming"}},
			wantStatus:   http.StatusCreated,
			wantContests: 1,
			wantBody: map[string]any{
				"startedAt": "2026-11-20T12:00:00.000Z", "endsAt": "2026-11-20T13:30:00.000Z", "now": "2026-11-20T12:00:00.000Z", "active": true,
				"solved": 0, "total": 4, "score": 0, "finishSeconds": 0, "ratingChange": nil, "ratingAfter": nil, "best": false,
				"problems": []map[string]any{
					{"position": 1, "title": "Easy One", "difficulty": "easy", "points": 3, "topic": "Brute Force", "rating": nil, "solved": false, "wrong": 0, "acceptedSeconds": nil},
					{"position": 2, "title": "Medium One", "difficulty": "medium", "points": 4, "topic": "Greedy", "rating": nil, "solved": false, "wrong": 0, "acceptedSeconds": nil},
					{"position": 3, "title": "Medium Three", "difficulty": "medium", "points": 4, "topic": "Math", "rating": nil, "solved": false, "wrong": 0, "acceptedSeconds": nil},
					{"position": 4, "title": "Hard One", "difficulty": "hard", "points": 6, "topic": "Dynamic Programming", "rating": nil, "solved": false, "wrong": 0, "acceptedSeconds": nil},
				},
			},
			wantRows: []map[string]any{
				{"position": 1, "slug": "cf-101-a", "title": "Easy One", "rating": 800, "topic": "Brute Force", "difficulty": "easy", "points": 3},
				{"position": 2, "slug": "cf-103-b", "title": "Medium One", "rating": 1000, "topic": "Greedy", "difficulty": "medium", "points": 4},
				{"position": 3, "slug": "cf-105-c", "title": "Medium Three", "rating": 1100, "topic": "Math", "difficulty": "medium", "points": 4},
				{"position": 4, "slug": "cf-106-d", "title": "Hard One", "rating": 1200, "topic": "Dynamic Programming", "difficulty": "hard", "points": 6},
			},
		},
		{name: "parent cannot start", role: "parent", problems: []pick{{"slug": "cf-101-a", "topic": "Math"}}, wantStatus: http.StatusForbidden, wantCode: "forbidden"},
		{name: "no problems", role: "student", wantStatus: http.StatusBadRequest, wantCode: "invalid_input"},
		{
			name:       "five problems",
			role:       "student",
			problems:   []pick{{"slug": "cf-101-a", "topic": "Math"}, {"slug": "cf-103-b", "topic": "Greedy"}, {"slug": "cf-104-b", "topic": "Sorting"}, {"slug": "cf-105-c", "topic": "Math"}, {"slug": "cf-106-d", "topic": "Math"}},
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_input",
		},
		{name: "two easy problems", role: "student", problems: []pick{{"slug": "cf-101-a", "topic": "Math"}, {"slug": "cf-102-a", "topic": "String"}}, wantStatus: http.StatusBadRequest, wantCode: "contest_slots"},
		{
			name:       "three medium problems",
			role:       "student",
			problems:   []pick{{"slug": "cf-103-b", "topic": "Greedy"}, {"slug": "cf-104-b", "topic": "Sorting"}, {"slug": "cf-105-c", "topic": "Math"}},
			wantStatus: http.StatusBadRequest,
			wantCode:   "contest_slots",
		},
		{name: "unknown slug", role: "student", problems: []pick{{"slug": "cf-999-a", "topic": "Math"}}, wantStatus: http.StatusBadRequest, wantCode: "invalid_input"},
		{name: "the same slug twice", role: "student", problems: []pick{{"slug": "cf-104-b", "topic": "Sorting"}, {"slug": "cf-104-b", "topic": "Greedy"}}, wantStatus: http.StatusBadRequest, wantCode: "invalid_input"},
		{name: "topic that the problem does not have", role: "student", problems: []pick{{"slug": "cf-102-a", "topic": "Math"}}, wantStatus: http.StatusBadRequest, wantCode: "invalid_input"},
		{
			name:         "problem of an earlier contest",
			role:         "student",
			contests:     []contestSeed{{start: -3 * time.Hour, problems: []problemSeed{{slug: "cf-101-a", rating: 800, difficulty: "easy", points: 3}}}},
			problems:     []pick{{"slug": "cf-103-b", "topic": "Greedy"}, {"slug": "cf-101-a", "topic": "Math"}},
			wantStatus:   http.StatusConflict,
			wantCode:     "problem_used",
			wantContests: 1,
		},
		{name: "problem in the library", role: "student", library: "cf-101-a", problems: []pick{{"slug": "cf-103-b", "topic": "Greedy"}, {"slug": "cf-101-a", "topic": "Math"}}, wantStatus: http.StatusConflict, wantCode: "problem_used"},
		{
			name:         "second start while a contest is active",
			role:         "student",
			contests:     []contestSeed{{start: -30 * time.Minute}},
			problems:     []pick{{"slug": "cf-102-a", "topic": "String"}},
			wantStatus:   http.StatusConflict,
			wantCode:     "contest_active",
			wantContests: 1,
		},
		{
			name:         "start at the end of the last contest",
			role:         "student",
			contests:     []contestSeed{{start: -90 * time.Minute}},
			problems:     []pick{{"slug": "cf-102-a", "topic": "String"}},
			wantStatus:   http.StatusCreated,
			wantContests: 2,
			wantBody:     map[string]any{"startedAt": "2026-11-20T12:00:00.000Z", "endsAt": "2026-11-20T13:30:00.000Z", "active": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := contestServer(t)
			seedContests(t, pool, tt.contests...)
			if tt.library != "" {
				exec(t, pool, "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,'Library','easy','[]','','','t0','t0')", newID(), tt.library)
			}
			w := contestStart(t, s, loginAs(t, s, tt.role), tt.problems...)
			if w.Code != tt.wantStatus || errorCode(w) != tt.wantCode {
				t.Fatalf("status = %d, code %q, want %d, %q; body %s", w.Code, errorCode(w), tt.wantStatus, tt.wantCode, w.Body.String())
			}
			if got := queryValue[int64](t, pool, "SELECT count(*) FROM contests"); got != tt.wantContests {
				t.Errorf("contests = %d, want %d", got, tt.wantContests)
			}
			if tt.wantBody != nil {
				problemsMatchBody(t, decode[any](t, w), tt.wantBody)
			}
			if tt.wantRows != nil {
				reviewsEqualJSON(t, problemsTableRows(t, pool, "SELECT position,slug,title,rating,topic,difficulty,points FROM contest_problems ORDER BY position"), tt.wantRows)
			}
		})
	}
}

// Two starts at the same time with different problems: the guard row of contest_starts lets only one of them make a contest.
func TestContestStartRace(t *testing.T) {
	pool, s := contestServer(t)
	student := loginAs(t, s, "student")
	start := func(slug, topic string) func() *httptest.ResponseRecorder {
		return func() *httptest.ResponseRecorder {
			return call(s, http.MethodPost, "/leetcode/api/contests", fmt.Sprintf(`{"problems":[{"slug":%q,"topic":%q}]}`, slug, topic), student)
		}
	}
	results := race(t, pool, "INSERT INTO contest_starts VALUES (1,'2000-01-01T00:00:00.000Z') ON CONFLICT (id) DO UPDATE SET started_at=excluded.started_at", nil,
		start("cf-101-a", "Math"), start("cf-102-a", "String"))
	statuses := []int{results[0].Code, results[1].Code}
	slices.Sort(statuses)
	if !slices.Equal(statuses, []int{http.StatusCreated, http.StatusConflict}) {
		t.Fatalf("statuses = %v, want [201 409]; bodies %s and %s", statuses, results[0].Body.String(), results[1].Body.String())
	}
	for _, w := range results {
		if w.Code == http.StatusConflict && errorCode(w) != "contest_active" {
			t.Errorf("error code = %q, want contest_active", errorCode(w))
		}
	}
	if got := queryValue[int64](t, pool, "SELECT count(*) FROM contests"); got != 1 {
		t.Errorf("contests = %d, want 1", got)
	}
}

func TestContestList(t *testing.T) {
	const minute = time.Minute
	easy := problemSeed{rating: 800, difficulty: "easy", points: 3}
	medium := problemSeed{rating: 1000, difficulty: "medium", points: 4}
	hard := problemSeed{rating: 1200, difficulty: "hard", points: 6}
	solved := func(p problemSeed, at time.Duration, wrong ...time.Duration) problemSeed {
		p.accepted, p.wrong = []time.Duration{at}, wrong
		return p
	}
	tests := []struct {
		name     string
		contests []contestSeed
		want     map[string]any
	}{
		{
			name: "no contests",
			want: map[string]any{"rating": 800, "bestRating": 800, "best": nil, "active": nil, "items": []any{}},
		},
		{
			// One easy problem unsolved at 800: round(32 × (0 − 0.5)) = −16.
			name:     "finished by time",
			contests: []contestSeed{{start: -90*minute - time.Second, problems: []problemSeed{easy}}},
			want: map[string]any{"rating": 784, "bestRating": 800, "active": nil, "items": []any{
				map[string]any{"id": contestID(0), "endsAt": "2026-11-20T11:59:59.000Z", "active": false, "solved": 0, "total": 1, "ratingChange": -16, "ratingAfter": 784},
			}},
		},
		{
			// The contests of TestRatingHistory: +30 to 830, then −17 to 813.
			name: "rating from history",
			contests: []contestSeed{
				{start: -5 * time.Hour, problems: []problemSeed{solved(easy, 10*minute), solved(medium, 30*minute), medium, hard}},
				{start: -3 * time.Hour, problems: []problemSeed{easy}},
			},
			want: map[string]any{"rating": 813, "bestRating": 830, "items": []any{
				map[string]any{"id": contestID(1), "ratingChange": -17, "ratingAfter": 813},
				map[string]any{"id": contestID(0), "ratingChange": 30, "ratingAfter": 830},
			}},
		},
		{
			name:     "active contest",
			contests: []contestSeed{{start: -30 * minute, problems: []problemSeed{solved(easy, 10*minute), medium}}},
			want: map[string]any{
				"rating": 800, "bestRating": 800, "best": nil,
				"active": map[string]any{"id": contestID(0), "endsAt": "2026-11-20T13:00:00.000Z", "now": "2026-11-20T12:00:00.000Z"},
				"items":  []any{map[string]any{"id": contestID(0), "active": true, "solved": 1, "total": 2, "score": 3, "finishSeconds": 600, "ratingChange": nil, "ratingAfter": nil, "best": false}},
			},
		},
		{
			name: "best is the higher score, then the shorter finish",
			contests: []contestSeed{
				{start: -7 * time.Hour, problems: []problemSeed{solved(easy, 10*minute)}},
				{start: -5 * time.Hour, problems: []problemSeed{solved(easy, 12*minute, 5*minute), solved(medium, 40*minute)}},
				{start: -3 * time.Hour, problems: []problemSeed{solved(easy, 10*minute), solved(medium, 40*minute)}},
			},
			want: map[string]any{
				"best": map[string]any{"id": contestID(2), "score": 7, "finishSeconds": 2400, "best": true},
				"items": []any{
					map[string]any{"id": contestID(2), "score": 7, "finishSeconds": 2400, "best": true},
					map[string]any{"id": contestID(1), "score": 7, "finishSeconds": 2700, "best": false},
					map[string]any{"id": contestID(0), "score": 3, "finishSeconds": 600, "best": false},
				},
			},
		},
		{
			name: "an active contest is never best",
			contests: []contestSeed{
				{start: -5 * time.Hour, problems: []problemSeed{solved(easy, 10*minute), solved(medium, 40*minute)}},
				{start: -3 * time.Hour, problems: []problemSeed{solved(easy, 10*minute)}},
				{start: -30 * minute, problems: []problemSeed{solved(easy, 5*minute), solved(hard, 20*minute)}},
			},
			want: map[string]any{
				"best":   map[string]any{"id": contestID(0), "score": 7, "best": true},
				"active": map[string]any{"id": contestID(2)},
				"items": []any{
					map[string]any{"id": contestID(2), "active": true, "score": 9, "best": false},
					map[string]any{"id": contestID(1), "active": false, "score": 3, "best": false},
					map[string]any{"id": contestID(0), "active": false, "score": 7, "best": true},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := contestServer(t)
			seedContests(t, pool, tt.contests...)
			w := request(t, s, "/contests", http.MethodGet, nil, loginAs(t, s, "parent"))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), tt.want)
		})
	}
}

func TestContestActive(t *testing.T) {
	tests := []struct {
		name     string
		contests []contestSeed
		want     any
	}{
		{name: "no active contest", contests: []contestSeed{{start: -90 * time.Minute}}, want: nil},
		{name: "active contest", contests: []contestSeed{{start: -3 * time.Hour}, {start: -89 * time.Minute}}, want: map[string]any{"id": contestID(1), "endsAt": "2026-11-20T12:01:00.000Z", "now": "2026-11-20T12:00:00.000Z"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := contestServer(t)
			seedContests(t, pool, tt.contests...)
			w := request(t, s, "/contests/active", http.MethodGet, nil, loginAs(t, s, "student"))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			reviewsEqualJSON(t, decode[any](t, w), tt.want)
		})
	}
}

func TestContestBankIndex(t *testing.T) {
	tests := []struct {
		name     string
		contests []contestSeed
		library  string
		want     map[string]any
	}{
		{
			name:     "used and library problems are left out",
			contests: []contestSeed{{start: -3 * time.Hour, problems: []problemSeed{{slug: "cf-101-a", rating: 800, difficulty: "easy", points: 3}}}},
			library:  "cf-105-c",
			want: map[string]any{
				"topics": []string{"Brute Force", "Dynamic Programming", "Greedy", "Math", "Sorting", "String"},
				"problems": []map[string]any{
					{"slug": "cf-102-a", "rating": 800, "topics": []string{"String"}, "difficulty": "easy", "points": 3},
					{"slug": "cf-103-b", "rating": 1000, "topics": []string{"Greedy"}, "difficulty": "medium", "points": 4},
					{"slug": "cf-104-b", "rating": 1000, "topics": []string{"Sorting", "Greedy"}, "difficulty": "medium", "points": 4},
					{"slug": "cf-106-d", "rating": 1200, "topics": []string{"Dynamic Programming", "Math"}, "difficulty": "hard", "points": 6},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := contestServer(t)
			seedContests(t, pool, tt.contests...)
			exec(t, pool, "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,'Library','medium','[]','','','t0','t0')", newID(), tt.library)
			w := request(t, s, "/contests/bank", http.MethodGet, nil, loginAs(t, s, "student"))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			reviewsEqualJSON(t, decode[any](t, w), tt.want)
		})
	}
}

func TestContestDetail(t *testing.T) {
	const minute = time.Minute
	tests := []struct {
		name       string
		contests   []contestSeed
		wantStatus int
		wantCode   string
		want       map[string]any
	}{
		{
			// The finish is the accepted submission at 40 minutes plus 5 minutes for each of the 2 wrong submissions before the accepted one at 12 minutes.
			// The rating change at 800: round(32 × ((1 − 0.5) + (1 − 0.2403) − 0.2403)) = round(32.62) = 33.
			name: "score from submissions",
			contests: []contestSeed{{start: -2 * time.Hour, problems: []problemSeed{
				{rating: 800, difficulty: "easy", points: 3, accepted: []time.Duration{12 * minute}, wrong: []time.Duration{5 * minute, 8 * minute, 20 * minute}},
				{rating: 1000, difficulty: "medium", points: 4, accepted: []time.Duration{40 * minute}},
				{rating: 1000, difficulty: "medium", points: 4, wrong: []time.Duration{10 * minute, 15 * minute}},
			}}},
			wantStatus: http.StatusOK,
			want: map[string]any{
				"id": contestID(0), "startedAt": "2026-11-20T10:00:00.000Z", "endsAt": "2026-11-20T11:30:00.000Z", "now": "2026-11-20T12:00:00.000Z", "active": false,
				"solved": 2, "total": 3, "score": 7, "finishSeconds": 3000, "ratingChange": 33, "ratingAfter": 833, "best": true,
				"problems": []map[string]any{
					{"position": 1, "title": "Problem 1", "difficulty": "easy", "points": 3, "topic": "Math", "rating": 800, "solved": true, "wrong": 2, "acceptedSeconds": 720},
					{"position": 2, "title": "Problem 2", "difficulty": "medium", "points": 4, "topic": "Math", "rating": 1000, "solved": true, "wrong": 0, "acceptedSeconds": 2400},
					{"position": 3, "title": "Problem 3", "difficulty": "medium", "points": 4, "topic": "Math", "rating": 1000, "solved": false, "wrong": 2, "acceptedSeconds": nil},
				},
			},
		},
		{
			name:       "removed problem",
			contests:   []contestSeed{{start: -2 * time.Hour, problems: []problemSeed{{slug: "cf-999-z", rating: 1200, difficulty: "hard", points: 6}}}},
			wantStatus: http.StatusOK,
			want: map[string]any{"problems": []map[string]any{
				{"position": 1, "title": "Problem 1", "difficulty": "hard", "points": 6, "topic": "Math", "rating": 1200, "solved": false, "wrong": 0, "acceptedSeconds": nil},
			}},
		},
		{name: "unknown contest", wantStatus: http.StatusNotFound, wantCode: "not_found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := contestServer(t)
			seedContests(t, pool, tt.contests...)
			w := request(t, s, "/contests/"+contestID(0), http.MethodGet, nil, loginAs(t, s, "parent"))
			if w.Code != tt.wantStatus || errorCode(w) != tt.wantCode {
				t.Fatalf("status = %d, code %q, want %d, %q; body %s", w.Code, errorCode(w), tt.wantStatus, tt.wantCode, w.Body.String())
			}
			if tt.want != nil {
				problemsMatchBody(t, decode[any](t, w), tt.want)
			}
		})
	}
}

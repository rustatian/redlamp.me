package notebook

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/bank"
	"github.com/rustatian/redlamp.me/backend/internal/editorials"
	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The expected values follow the rules of the contests: easy 800, medium 900 to 1100, hard 1200;
// places easy 1, medium 2, hard 1; points easy 3, medium 4, hard 6; a contest lasts 90 minutes;
// 5 minutes for each wrong submission before the first accepted one of a solved problem;
// the rating starts at 800 and changes by round(32 × Σ(S − E)) with E = 1 / (1 + 10^((R_problem − R) / 400)).

// contestNow is the clock of the contest tests.
var contestNow = time.Date(2026, 11, 20, 12, 0, 0, 0, time.UTC)

// contestBank is the bank of the contest tests: two problems at 800, two at 1000, one at 1100 and one at 1200.
// Only cf-101-a has a statement, an example and hidden tests.
var contestBank = []bank.Problem{
	{
		Slug: "cf-101-a", Title: "Easy One", Rating: 800, Topics: []string{"Math", "Brute Force"}, Statement: "<p>Выведите ответ.</p>", IgnoreCase: true,
		Examples: []bank.Test{{Input: "1\n", Output: "1\n"}},
		Tests:    []bank.Test{{Input: "hidden-input-7731\n", Output: "hidden-answer-5519"}, {Input: "2\n", Output: "4"}},
	},
	{Slug: "cf-102-a", Title: "Easy Two", Rating: 800, Topics: []string{"String"}},
	{Slug: "cf-103-b", Title: "Medium One", Rating: 1000, Topics: []string{"Greedy"}},
	{Slug: "cf-104-b", Title: "Medium Two", Rating: 1000, Topics: []string{"Sorting", "Greedy"}},
	{Slug: "cf-105-c", Title: "Medium Three", Rating: 1100, Topics: []string{"Math"}},
	{Slug: "cf-106-d", Title: "Hard One", Rating: 1200, Topics: []string{"Dynamic Programming", "Math"}},
}

// contestEditorial is the Editorial of cf-101-a. contestEditorialBody is its JSON in the problem view.
var (
	contestEditorial = editorials.Editorial{HTML: "<p>Разбор задачи.</p>", Figures: []sets.Figure{{ID: "answer", Title: "Ответ", Frames: []sets.FigureFrame{
		{Caption: "Ответ 1.", Rows: []sets.FigureRow{{Cells: []any{1}, Marked: []int{0}}}},
	}}}}
	contestEditorialBody = map[string]any{"html": "<p>Разбор задачи.</p>", "figures": []any{map[string]any{"id": "answer", "title": "Ответ", "frames": []any{
		map[string]any{"caption": "Ответ 1.", "rows": []any{map[string]any{"cells": []any{1}, "marked": []any{0}}}},
	}}}}
)

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

// contestServer returns an empty database and a server with contestBank, the Editorial of cf-101-a and the clock contestNow.
func contestServer(t *testing.T) (*pgxpool.Pool, *Server) {
	t.Helper()
	pool := database(t)
	s := newServerWithBank(t, pool, contestBank)
	s.editorials = map[string]editorials.Editorial{"cf-101-a": contestEditorial}
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

func TestContestVerdict(t *testing.T) {
	value := func(s string) judgeOutput { return judgeOutput{Value: &s} }
	failure := func(s string) judgeOutput { return judgeOutput{Error: &s} }
	timeout := judgeOutput{Timeout: true}
	tests := []struct {
		name        string
		ignoreCase  bool
		outputs     []judgeOutput
		wantVerdict string
		wantTest    int
		wantMessage string
	}{
		{name: "all outputs match", outputs: []judgeOutput{value("1"), value("2"), value("YES")}, wantVerdict: "accepted"},
		{name: "wrong second test", outputs: []judgeOutput{value("1"), value("3"), value("YES")}, wantVerdict: "wrong", wantTest: 2},
		{name: "case differs with ignoreCase", ignoreCase: true, outputs: []judgeOutput{value("1"), value("2"), value("yes")}, wantVerdict: "accepted"},
		{name: "case differs without ignoreCase", outputs: []judgeOutput{value("1"), value("2"), value("yes")}, wantVerdict: "wrong", wantTest: 3},
		{
			name:        "error on the first test",
			outputs:     []judgeOutput{failure("ZeroDivisionError: division by zero"), timeout, timeout},
			wantVerdict: "error",
			wantTest:    1,
			wantMessage: "ZeroDivisionError: division by zero",
		},
		{name: "long error message is cut", outputs: []judgeOutput{failure(strings.Repeat("e", 2001)), timeout, timeout}, wantVerdict: "error", wantTest: 1, wantMessage: strings.Repeat("e", 2000) + "…"},
		{name: "timeout on the third test", outputs: []judgeOutput{value("1"), value("2"), timeout}, wantVerdict: "timeout", wantTest: 3},
		{name: "extra spaces and newlines", outputs: []judgeOutput{value("1 \n\n"), value("2"), value("YES")}, wantVerdict: "accepted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := bank.Problem{IgnoreCase: tt.ignoreCase, Tests: []bank.Test{{Input: "a\n", Output: "1"}, {Input: "b\n", Output: "2"}, {Input: "c\n", Output: "YES"}}}
			verdict, test, message := contestVerdict(p, tt.outputs)
			if verdict != tt.wantVerdict || test != tt.wantTest || message != tt.wantMessage {
				t.Errorf("contestVerdict = %q, %d, %q, want %q, %d, %q", verdict, test, message, tt.wantVerdict, tt.wantTest, tt.wantMessage)
			}
		})
	}
}

func TestContestTimeRules(t *testing.T) {
	right := map[string]any{"code": "print(1)", "outputs": []map[string]any{{"value": "hidden-answer-5519"}, {"value": "4"}}}
	save := map[string]any{"code": "print(2)"}
	tests := []struct {
		name string
		// at is the time from the end of the contest.
		at         time.Duration
		method     string
		action     string
		body       any
		wantStatus int
		wantCode   string
		// wantStored is contest_problems.code after the request.
		wantStored      string
		wantSubmissions int64
	}{
		{name: "check one second before the end", at: -time.Second, method: http.MethodPost, action: "check", wantStatus: http.StatusOK},
		{name: "check at the end", method: http.MethodPost, action: "check", wantStatus: http.StatusConflict, wantCode: "contest_over"},
		{name: "submission 59 seconds after the end", at: 59 * time.Second, method: http.MethodPost, action: "submissions", body: right, wantStatus: http.StatusCreated, wantStored: "print(1)", wantSubmissions: 1},
		{name: "submission 61 seconds after the end", at: 61 * time.Second, method: http.MethodPost, action: "submissions", body: right, wantStatus: http.StatusConflict, wantCode: "contest_over"},
		{name: "code save one second before the end", at: -time.Second, method: http.MethodPut, action: "code", body: save, wantStatus: http.StatusNoContent, wantStored: "print(2)"},
		{name: "code save at the end", method: http.MethodPut, action: "code", body: save, wantStatus: http.StatusConflict, wantCode: "contest_over"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := contestServer(t)
			// The contest ends at contestNow.
			seedContests(t, pool, contestSeed{start: -90 * time.Minute, problems: []problemSeed{{slug: "cf-101-a", rating: 800, difficulty: "easy", points: 3}}})
			s.now = func() time.Time { return contestNow.Add(tt.at) }
			w := request(t, s, "/contests/"+contestID(0)+"/problems/1/"+tt.action, tt.method, tt.body, loginAs(t, s, "student"))
			if w.Code != tt.wantStatus || errorCode(w) != tt.wantCode {
				t.Fatalf("status = %d, code %q, want %d, %q; body %s", w.Code, errorCode(w), tt.wantStatus, tt.wantCode, w.Body.String())
			}
			if got := queryValue[string](t, pool, "SELECT code FROM contest_problems"); got != tt.wantStored {
				t.Errorf("stored code = %q, want %q", got, tt.wantStored)
			}
			if got := queryValue[int64](t, pool, "SELECT count(*) FROM contest_submissions"); got != tt.wantSubmissions {
				t.Errorf("submissions = %d, want %d", got, tt.wantSubmissions)
			}
		})
	}
}

func TestContestSubmission(t *testing.T) {
	easy := problemSeed{slug: "cf-101-a", rating: 800, difficulty: "easy", points: 3}
	solved := problemSeed{slug: "cf-101-a", rating: 800, difficulty: "easy", points: 3, accepted: []time.Duration{10 * time.Minute}}
	removed := problemSeed{slug: "cf-999-z", rating: 800, difficulty: "easy", points: 3}
	// active is a contest of the problem p that started 30 minutes ago.
	active := func(p problemSeed) []contestSeed {
		return []contestSeed{{start: -30 * time.Minute, problems: []problemSeed{p}}}
	}
	submit := func(code string, outputs ...map[string]any) map[string]any {
		return map[string]any{"code": code, "outputs": outputs}
	}
	right := []map[string]any{{"value": "hidden-answer-5519"}, {"value": "4"}}
	const zeroDivision = "ZeroDivisionError: integer division or modulo by zero"
	tests := []struct {
		name     string
		contests []contestSeed
		role     string
		method   string
		// path is the path after /contests/<id>.
		path       string
		body       any
		wantStatus int
		wantCode   string
		// wantBody is the whole body. wantView has fields of the problem view after the request.
		wantBody        any
		wantView        map[string]any
		wantSubmissions int64
	}{
		{
			name: "problem view before the end", contests: active(easy), role: "parent", method: http.MethodGet, path: "/problems/1", wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"position": 1, "title": "Problem 1", "difficulty": "easy", "points": 3, "topic": "Math", "statement": "<p>Выведите ответ.</p>", "ignoreCase": true,
				"examples": []map[string]any{{"input": "1\n", "output": "1\n"}}, "code": "", "solved": false, "submissions": []any{}, "url": nil, "editorial": nil,
				"endsAt": "2026-11-20T13:00:00.000Z", "now": "2026-11-20T12:00:00.000Z",
			},
		},
		{
			name:       "link after the end",
			contests:   []contestSeed{{start: -2 * time.Hour, problems: []problemSeed{easy}}},
			role:       "parent",
			method:     http.MethodGet,
			path:       "/problems/1",
			wantStatus: http.StatusOK,
			wantView:   map[string]any{"url": "https://codeforces.com/problemset/problem/101/A", "editorial": contestEditorialBody},
		},
		{
			name: "no Editorial one millisecond before the end", contests: []contestSeed{{start: -(90*time.Minute - time.Millisecond), problems: []problemSeed{easy}}}, role: "parent", method: http.MethodGet, path: "/problems/1", wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"position": 1, "title": "Problem 1", "difficulty": "easy", "points": 3, "topic": "Math", "statement": "<p>Выведите ответ.</p>", "ignoreCase": true,
				"examples": []map[string]any{{"input": "1\n", "output": "1\n"}}, "code": "", "solved": false, "submissions": []any{}, "url": nil, "editorial": nil,
				"endsAt": "2026-11-20T12:00:00.001Z", "now": "2026-11-20T12:00:00.000Z",
			},
			wantView: map[string]any{"url": nil, "editorial": nil},
		},
		{
			name: "Editorial at the end", contests: []contestSeed{{start: -90 * time.Minute, problems: []problemSeed{easy}}}, role: "parent", method: http.MethodGet, path: "/problems/1", wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"position": 1, "title": "Problem 1", "difficulty": "easy", "points": 3, "topic": "Math", "statement": "<p>Выведите ответ.</p>", "ignoreCase": true,
				"examples": []map[string]any{{"input": "1\n", "output": "1\n"}}, "code": "", "solved": false, "submissions": []any{}, "url": "https://codeforces.com/problemset/problem/101/A", "editorial": contestEditorialBody,
				"endsAt": "2026-11-20T12:00:00.000Z", "now": "2026-11-20T12:00:00.000Z",
			},
			wantView: map[string]any{"url": "https://codeforces.com/problemset/problem/101/A", "editorial": contestEditorialBody},
		},
		{name: "check gives the hidden inputs", contests: active(easy), role: "student", method: http.MethodPost, path: "/problems/1/check", wantStatus: http.StatusOK, wantBody: map[string]any{"inputs": []string{"hidden-input-7731\n", "2\n"}}},
		{
			name:            "accepted submission solves the problem",
			contests:        active(easy),
			role:            "student",
			method:          http.MethodPost,
			path:            "/problems/1/submissions",
			body:            submit("print(answer)", right...),
			wantStatus:      http.StatusCreated,
			wantBody:        map[string]any{"verdict": "accepted", "test": nil, "error": nil},
			wantView:        map[string]any{"code": "print(answer)", "solved": true, "submissions": []map[string]any{{"verdict": "accepted", "test": nil, "error": nil, "createdAt": "2026-11-20T12:00:00.000Z"}}},
			wantSubmissions: 1,
		},
		{
			// The refused submission keeps the code of the problem.
			name:            "submission after an accepted one",
			contests:        active(solved),
			role:            "student",
			method:          http.MethodPost,
			path:            "/problems/1/submissions",
			body:            submit("print(2)", right...),
			wantStatus:      http.StatusConflict,
			wantCode:        "solved",
			wantView:        map[string]any{"code": "", "solved": true},
			wantSubmissions: 1,
		},
		{name: "check of a solved problem", contests: active(solved), role: "student", method: http.MethodPost, path: "/problems/1/check", wantStatus: http.StatusConflict, wantCode: "solved", wantSubmissions: 1},
		{name: "parent submission", contests: active(easy), role: "parent", method: http.MethodPost, path: "/problems/1/submissions", body: submit("print(1)", right...), wantStatus: http.StatusForbidden, wantCode: "forbidden"},
		{name: "wrong number of outputs", contests: active(easy), role: "student", method: http.MethodPost, path: "/problems/1/submissions", body: submit("print(1)", right[0]), wantStatus: http.StatusBadRequest, wantCode: "invalid_input"},
		{
			// The exact body has no hidden input and no hidden answer.
			name:            "wrong answer hides the test",
			contests:        active(easy),
			role:            "student",
			method:          http.MethodPost,
			path:            "/problems/1/submissions",
			body:            submit("print('x')", map[string]any{"value": "x"}, map[string]any{"value": "4"}),
			wantStatus:      http.StatusCreated,
			wantBody:        map[string]any{"verdict": "wrong", "test": 1, "error": nil},
			wantView:        map[string]any{"code": "print('x')", "solved": false, "submissions": []map[string]any{{"verdict": "wrong", "test": 1, "error": nil}}},
			wantSubmissions: 1,
		},
		{
			name:            "error message is kept",
			contests:        active(easy),
			role:            "student",
			method:          http.MethodPost,
			path:            "/problems/1/submissions",
			body:            submit("print(1 // 0)", map[string]any{"error": zeroDivision}, map[string]any{"timeout": true}),
			wantStatus:      http.StatusCreated,
			wantBody:        map[string]any{"verdict": "error", "test": 1, "error": zeroDivision},
			wantView:        map[string]any{"submissions": []map[string]any{{"verdict": "error", "test": 1, "error": zeroDivision}}},
			wantSubmissions: 1,
		},
		{
			name: "removed problem", contests: active(removed), role: "student", method: http.MethodGet, path: "/problems/1", wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"position": 1, "title": "Problem 1", "difficulty": "easy", "points": 3, "topic": "Math", "statement": nil, "ignoreCase": false,
				"examples": []any{}, "code": "", "solved": false, "submissions": []any{}, "url": nil, "editorial": nil, "endsAt": "2026-11-20T13:00:00.000Z", "now": "2026-11-20T12:00:00.000Z",
			},
		},
		{name: "check of a removed problem", contests: active(removed), role: "student", method: http.MethodPost, path: "/problems/1/check", wantStatus: http.StatusNotFound, wantCode: "problem_removed"},
		{name: "submission for a position without a problem", contests: active(easy), role: "student", method: http.MethodPost, path: "/problems/2/submissions", body: submit("print(1)", right...), wantStatus: http.StatusNotFound, wantCode: "not_found"},
		{name: "unknown contest", role: "student", method: http.MethodGet, path: "/problems/1", wantStatus: http.StatusNotFound, wantCode: "not_found"},
		{name: "code of 100,001 characters", contests: active(easy), role: "student", method: http.MethodPut, path: "/problems/1/code", body: map[string]any{"code": strings.Repeat("a", 100_001)}, wantStatus: http.StatusBadRequest, wantCode: "invalid_input"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s := contestServer(t)
			seedContests(t, pool, tt.contests...)
			w := request(t, s, "/contests/"+contestID(0)+tt.path, tt.method, tt.body, loginAs(t, s, tt.role))
			if w.Code != tt.wantStatus || errorCode(w) != tt.wantCode {
				t.Fatalf("status = %d, code %q, want %d, %q; body %s", w.Code, errorCode(w), tt.wantStatus, tt.wantCode, w.Body.String())
			}
			if tt.wantBody != nil {
				reviewsEqualJSON(t, decode[any](t, w), tt.wantBody)
			}
			if tt.wantView != nil {
				view := request(t, s, "/contests/"+contestID(0)+"/problems/1", http.MethodGet, nil, loginAs(t, s, "student"))
				problemsMatchBody(t, decode[any](t, view), tt.wantView)
			}
			if got := queryValue[int64](t, pool, "SELECT count(*) FROM contest_submissions"); got != tt.wantSubmissions {
				t.Errorf("submissions = %d, want %d", got, tt.wantSubmissions)
			}
		})
	}
}

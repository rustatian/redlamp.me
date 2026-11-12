package notebook

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The cases of test/game.test.ts.

// The prepared statements of test/game.test.ts.
const (
	gameInsertProblem    = "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,archived_at,created_at,updated_at) VALUES ($1,$2,'Problem',$3,'[]','','',$4,'now','now')"
	gameInsertAttempt    = "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,saved_at,created_at,updated_at) VALUES ($1,$2,$3,'{}','',$4,$5,$6,'now','now')"
	gameInsertReview     = "INSERT INTO reviews (id,problem_id,result,note,reviewed_on,created_at) VALUES ($1,$2,$3,'',$4,'now')"
	gameInsertBundle     = "INSERT INTO homework (id,title,instructions,due_date,created_at,updated_at) VALUES ($1,'Week 41','',$2,'now','now')"
	gameInsertTask       = "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,'now')"
	gameInsertCard       = "INSERT INTO card_answers (id,card_id,correct,interval_days,answered_on,next_date,created_at) VALUES ($1,'pc1-factorial',$2,3,$3,$4,'now')"
	gameInsertQuiz       = "INSERT INTO quiz_attempts (id,deck,correct,total,passed,created_at) VALUES ($1,$2,$3,$4,$5,$6)"
	gameInsertSubmission = "INSERT INTO submissions (id,task_id,attempt_id,task_version,created_at) VALUES ($1,$2,$3,1,$4)"
	gameInsertRun        = "INSERT INTO test_runs (id,problem_id,attempt_id,approach_id,code_hash,passed,total,run_on,created_at) VALUES ($1,$2,$3,'approach','hash',$4,$5,$6,$7)"
	gameInsertVisit      = "INSERT INTO visits (day) VALUES ($1)"
)

// gameOrder is the badge order of a Game.
var gameOrder = []string{"first-accept", "three-day-streak", "first-hard", "ten-accepted", "independent-five", "set-cleared", "week-no-miss", "month-no-miss", "twenty-five-accepted", "fifty-accepted", "fifty-cards", "ten-recalls", "first-codeforces", "first-try", "two-sets", "on-time-five", "level-five"}

// gameAttempt is an attempt of a Spec. An empty state is "saved", an empty understanding is "with_help", an empty savedAt is NULL.
type gameAttempt struct{ state, acceptance, understanding, savedAt string }

type gameReview struct{ result, on string }

type gameRun struct {
	passed, total int
	on            string
}

// gameSpec is Spec of the TypeScript test: one problem and its rows. An empty slug is the problem ID. An empty due is NULL.
type gameSpec struct {
	difficulty, slug string
	archived         bool
	attempts         []gameAttempt
	reviews          []gameReview
	task, due        string
	submitted, added string
	postpone         bool
	runs             []gameRun
}

type gameCard struct {
	correct bool
	on      string
}

type gameQuizAttempt struct {
	deck           string
	correct, total int
	passed         bool
}

// gameBadge is one badge of a Game.
type gameBadge struct {
	ID     string `json:"id"`
	Earned bool   `json:"earned"`
}

// gameView is the part of a Game that the tests read.
type gameView struct {
	Level  int         `json:"level"`
	Badges []gameBadge `json:"badges"`
}

// gameNull returns nil for "", so the column gets NULL.
func gameNull(value string) any {
	if value == "" {
		return nil
	}
	return value
}

// gameBit is Number(value) of a boolean.
func gameBit(value bool) int {
	if value {
		return 1
	}
	return 0
}

// gameAccepted is the attempts of accepted(): one saved accepted attempt.
func gameAccepted(understanding string) []gameAttempt {
	return []gameAttempt{{acceptance: "accepted", understanding: understanding}}
}

// gameRunOf is run(): a test run with 5 tests.
func gameRunOf(passed int, on string) gameRun { return gameRun{passed: passed, total: 5, on: on} }

// gameQuiz is quiz(): an attempt of the deck of a chapter of permutations-and-combinations.
func gameQuiz(chapter, correct, total int, passed bool) gameQuizAttempt {
	return gameQuizAttempt{deck: fmt.Sprintf("permutations-and-combinations-%d", chapter), correct: correct, total: total, passed: passed}
}

// gameDay is day(): the UTC date offset days from today.
func gameDay(offset int) string {
	return time.Now().UTC().AddDate(0, 0, offset).Format(time.DateOnly)
}

// gameDays is days(): the day offsets from from to to.
func gameDays(from, to int) []int {
	var offsets []int
	for offset := from; offset <= to; offset++ {
		offsets = append(offsets, offset)
	}
	return offsets
}

// gameCardDays is cardDays(): count correct answers on consecutive days from 2026-01-01.
func gameCardDays(count int) []gameCard {
	cards := make([]gameCard, count)
	for i := range cards {
		cards[i] = gameCard{correct: true, on: time.Date(2026, 1, 1+i, 0, 0, 0, 0, time.UTC).Format(time.DateOnly)}
	}
	return cards
}

// gameRecalls is recalls(): count "with_help" recalls from 2026-09-01.
func gameRecalls(count int) []gameReview {
	reviews := make([]gameReview, count)
	for i := range reviews {
		reviews[i] = gameReview{result: "with_help", on: fmt.Sprintf("2026-09-%02d", i+1)}
	}
	return reviews
}

// gameOnTime is onTime(): a task submitted on the due date, or the day after it when late is true.
func gameOnTime(late bool) gameSpec {
	submitted := "2026-10-01T12:00:00.000Z"
	if late {
		submitted = "2026-10-02T12:00:00.000Z"
	}
	return gameSpec{difficulty: "medium", task: "submitted", due: "2026-10-01", submitted: submitted}
}

// gameTaskSlug is taskSlug(): the problem slug of the URL of a set task.
func gameTaskSlug(task SetTask) string {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(task.URL, "https://"), "/"), "/")
	if parts[0] == "codeforces.com" {
		return "cf-" + parts[3] + "-" + strings.ToLower(parts[4])
	}
	return parts[2]
}

// gameSetProblems is an accepted problem for each task of the named sets, in the order of sets, without the task with the slug skip.
func gameSetProblems(sets []LearningSet, skip string, slugs ...string) []gameSpec {
	var specs []gameSpec
	for _, set := range sets {
		if !slices.Contains(slugs, set.Slug) {
			continue
		}
		for _, chapter := range set.Chapters {
			for _, task := range chapter.Tasks {
				if slug := gameTaskSlug(task); slug != skip {
					specs = append(specs, gameSpec{difficulty: task.Difficulty, slug: slug, attempts: gameAccepted("with_help")})
				}
			}
		}
	}
	return specs
}

// gameSeed inserts the rows of one Spec in the order of the TypeScript batch.
func gameSeed(t *testing.T, pool *pgxpool.Pool, spec gameSpec) {
	t.Helper()
	id, homeworkID, taskID, attemptID := newID(), newID(), newID(), newID()
	var archived any
	if spec.archived {
		archived = "now"
	}
	exec(t, pool, gameInsertProblem, id, cmp.Or(spec.slug, id), spec.difficulty, archived)
	for _, attempt := range spec.attempts {
		exec(t, pool, gameInsertAttempt, newID(), id, cmp.Or(attempt.state, "saved"), attempt.acceptance, cmp.Or(attempt.understanding, "with_help"), gameNull(attempt.savedAt))
	}
	for _, review := range spec.reviews {
		exec(t, pool, gameInsertReview, newID(), id, review.result, review.on)
	}
	if spec.task != "" {
		exec(t, pool, gameInsertBundle, homeworkID, gameNull(spec.due))
		exec(t, pool, gameInsertTask, taskID, homeworkID, id, spec.task, cmp.Or(spec.added, "2026-09-01T00:00:00.000Z"))
	}
	if spec.postpone {
		// date(due_date,'+1 day') of the TypeScript statement.
		due, err := time.Parse(time.DateOnly, spec.due)
		if err != nil {
			t.Fatal(err)
		}
		exec(t, pool, "UPDATE homework SET due_date=$1 WHERE id=$2", due.AddDate(0, 0, 1).Format(time.DateOnly), homeworkID)
	}
	if len(spec.runs) > 0 || spec.submitted != "" {
		exec(t, pool, gameInsertAttempt, attemptID, id, "draft", "not_submitted", "needs_practice", nil)
	}
	// The list order is the time order of the runs.
	for i, run := range spec.runs {
		exec(t, pool, gameInsertRun, newID(), id, attemptID, run.passed, run.total, run.on, fmt.Sprintf("%sT12:00:00.%03dZ", run.on, i))
	}
	if spec.submitted != "" {
		exec(t, pool, gameInsertSubmission, newID(), taskID, attemptID, spec.submitted)
	}
}

func TestGameContentHomework(t *testing.T) {
	today := gameDay(0)
	// An empty due or completed is NULL.
	tests := []struct {
		name, kind, state, due, created, completed string
		xp, late                                   int
	}{
		{name: "chapter completed before the due date", kind: "chapter", state: "completed", due: "2026-10-01", created: "2026-09-01", completed: "2026-09-30T12:00:00.000Z", xp: 15, late: 0},
		{name: "quiz completed on the due date", kind: "quiz", state: "completed", due: "2026-10-01", created: "2026-09-01", completed: "2026-10-01T23:59:00.000Z", xp: 15, late: 0},
		{name: "quiz completed after the due date", kind: "quiz", state: "completed", due: "2026-10-01", created: "2026-09-01", completed: "2026-10-02T00:01:00.000Z", xp: 0, late: 1},
		{name: "chapter still assigned after the due date", kind: "chapter", state: "assigned", due: "2026-10-01", created: "2026-09-01", completed: "", xp: 0, late: 1},
		{name: "quiz still assigned on the due date", kind: "quiz", state: "assigned", due: today, created: "2026-09-01", completed: "", xp: 0, late: 0},
		{name: "cancelled chapter after the due date", kind: "chapter", state: "cancelled", due: "2026-10-01", created: "2026-09-01", completed: "", xp: 0, late: 0},
		{name: "chapter assigned after the due date", kind: "chapter", state: "assigned", due: "2026-10-01", created: "2026-10-02", completed: "", xp: 0, late: 0},
		{name: "chapter completed without a due date", kind: "chapter", state: "completed", due: "", created: "2026-09-01", completed: "2026-10-02T12:00:00.000Z", xp: 15, late: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithSets(t, pool, goldenSets(t))
			student := loginAs(t, s, "student")
			homeworkID := newID()
			exec(t, pool, gameInsertBundle, homeworkID, gameNull(tt.due))
			exec(t, pool, "INSERT INTO homework_tasks (id,homework_id,kind,set_slug,chapter,state,created_at,updated_at,completed_at) VALUES ($1,$2,$3,'complexity',1,$4,$5,$6,$7)",
				newID(), homeworkID, tt.kind, tt.state, tt.created, cmp.Or(tt.completed, tt.created), gameNull(tt.completed))
			w := request(t, s, "/game", http.MethodGet, nil, student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, body %s", w.Code, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"xp": tt.xp, "late": tt.late, "missed": 0})
		})
	}
}

func TestGameCountsFiveContentTasksOnTime(t *testing.T) {
	pool := database(t)
	s := newServerWithSets(t, pool, goldenSets(t))
	student := loginAs(t, s, "student")
	homeworkID := newID()
	content := []struct {
		name, kind, set string
		chapter         int
	}{
		{name: "complexity reading", kind: "chapter", set: "complexity", chapter: 1},
		{name: "complexity quiz", kind: "quiz", set: "complexity", chapter: 1},
		{name: "counting reading", kind: "chapter", set: "permutations-and-combinations", chapter: 1},
		{name: "counting quiz", kind: "quiz", set: "permutations-and-combinations", chapter: 1},
		{name: "two pointers reading", kind: "chapter", set: "two-pointers", chapter: 1},
	}
	exec(t, pool, gameInsertBundle, homeworkID, "2026-10-01")
	for _, task := range content {
		exec(t, pool, "INSERT INTO homework_tasks (id,homework_id,kind,set_slug,chapter,state,created_at,updated_at,completed_at) VALUES ($1,$2,$3,$4,$5,'completed','2026-09-01','2026-10-01','2026-10-01')",
			newID(), homeworkID, task.kind, task.set, task.chapter)
	}
	w := request(t, s, "/game", http.MethodGet, nil, student)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", w.Code, w.Body.String())
	}
	problemsMatchBody(t, decode[any](t, w), map[string]any{"xp": 75, "late": 0})
	if badges := decode[gameView](t, w).Badges; !slices.Contains(badges, gameBadge{ID: "on-time-five", Earned: true}) {
		t.Errorf("badges = %v, want on-time-five earned", badges)
	}
}

func TestGameXPLevelAndBadges(t *testing.T) {
	sets := goldenSets(t)
	today := gameDay(0)
	tests := []struct {
		name                        string
		role                        string
		problems                    []gameSpec
		cards                       []gameCard
		quizzes                     []gameQuizAttempt
		visits                      []int
		xp, late, missed            int
		level, levelXp, nextLevelXp int
		earned                      []string
	}{
		{name: "empty notebook", problems: nil, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "accepted easy", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, xp: 10, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "two accepted attempts count once", problems: []gameSpec{{difficulty: "easy", attempts: []gameAttempt{{acceptance: "accepted"}, {acceptance: "accepted"}}}}, xp: 10, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "parent reads the same game", role: "parent", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, xp: 10, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "accepted hard", problems: []gameSpec{{difficulty: "hard", attempts: gameAccepted("with_help")}}, xp: 40, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept", "first-hard"}},
		{name: "saved but not accepted", problems: []gameSpec{{difficulty: "medium", attempts: []gameAttempt{{acceptance: "not_accepted"}}}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "accepted draft", problems: []gameSpec{{difficulty: "medium", attempts: []gameAttempt{{state: "draft", acceptance: "accepted"}}}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "independent attempt", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("independent")}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "independent recall after help", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), reviews: []gameReview{{result: "independent", on: "2026-10-01"}}}}, xp: 25, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "later needs-practice recall keeps bonus", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("independent"), reviews: []gameReview{{result: "needs_practice", on: "2026-10-02"}}}}, xp: 25, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "independent attempt and independent recall count once", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("independent"), reviews: []gameReview{{result: "independent", on: "2026-10-01"}}}}, xp: 25, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "two recalls on one day", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), reviews: []gameReview{{result: "with_help", on: "2026-10-01"}, {result: "needs_practice", on: "2026-10-01"}}}}, xp: 15, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "recalls of two problems on one day", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), reviews: []gameReview{{result: "with_help", on: "2026-10-01"}}}, {difficulty: "easy", attempts: gameAccepted("with_help"), reviews: []gameReview{{result: "with_help", on: "2026-10-01"}}}}, xp: 30, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "completed task", problems: []gameSpec{{difficulty: "medium", task: "completed"}}, xp: 15, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "cancelled task", problems: []gameSpec{{difficulty: "medium", task: "cancelled"}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "archived accepted problem", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), archived: true}}, xp: 10, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "95 XP stays level 1", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help")}, {difficulty: "medium", attempts: gameAccepted("with_help")}, {difficulty: "medium", attempts: gameAccepted("with_help")}, {difficulty: "medium", attempts: gameAccepted("with_help"), reviews: []gameReview{{result: "with_help", on: "2026-10-01"}, {result: "with_help", on: "2026-10-02"}, {result: "needs_practice", on: "2026-10-03"}}}}, xp: 95, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "exactly 100 XP is level 2", problems: slices.Repeat([]gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help")}}, 5), xp: 100, level: 2, levelXp: 100, nextLevelXp: 300, earned: []string{"first-accept"}},
		{name: "ten accepted", problems: slices.Repeat([]gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, 10), xp: 100, level: 2, levelXp: 100, nextLevelXp: 300, earned: []string{"first-accept", "ten-accepted"}},
		{name: "five independent", problems: slices.Repeat([]gameSpec{{difficulty: "easy", attempts: gameAccepted("independent")}}, 5), xp: 100, level: 2, levelXp: 100, nextLevelXp: 300, earned: []string{"first-accept", "independent-five"}},
		{name: "three save days in a row", problems: []gameSpec{{difficulty: "medium", attempts: []gameAttempt{{acceptance: "not_accepted", savedAt: "2026-10-01T23:59:00.000Z"}, {acceptance: "not_accepted", savedAt: "2026-10-02T00:01:00.000Z"}, {acceptance: "not_accepted", savedAt: "2026-10-03T12:00:00.000Z"}}}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"three-day-streak"}},
		{name: "gap day", problems: []gameSpec{{difficulty: "medium", attempts: []gameAttempt{{acceptance: "not_accepted", savedAt: "2026-10-01T12:00:00.000Z"}, {acceptance: "not_accepted", savedAt: "2026-10-02T12:00:00.000Z"}, {acceptance: "not_accepted", savedAt: "2026-10-04T12:00:00.000Z"}}}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "set cleared", problems: gameSetProblems(sets, "", "permutations-and-combinations"), xp: 370, level: 3, levelXp: 300, nextLevelXp: 600, earned: []string{"first-accept", "first-hard", "ten-accepted", "set-cleared", "first-codeforces"}},
		{name: "set one short", problems: gameSetProblems(sets, "next-permutation", "permutations-and-combinations"), xp: 350, level: 3, levelXp: 300, nextLevelXp: 600, earned: []string{"first-accept", "first-hard", "ten-accepted", "first-codeforces"}},
		{name: "correct card", cards: []gameCard{{correct: true, on: "2026-10-01"}}, xp: 2, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "wrong card", cards: []gameCard{{correct: false, on: "2026-10-01"}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "same card on two days", cards: []gameCard{{correct: true, on: "2026-10-01"}, {correct: true, on: "2026-10-02"}}, xp: 4, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "first test pass", problems: []gameSpec{{difficulty: "easy", runs: []gameRun{gameRunOf(5, "2026-10-01")}}}, xp: 10, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-try"}},
		{name: "second test pass on the same problem", problems: []gameSpec{{difficulty: "easy", runs: []gameRun{gameRunOf(5, "2026-10-01"), gameRunOf(5, "2026-10-02")}}}, xp: 10, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-try"}},
		{name: "test pass after a failed first run", problems: []gameSpec{{difficulty: "easy", runs: []gameRun{gameRunOf(4, "2026-10-01"), gameRunOf(5, "2026-10-02")}}}, xp: 8, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "test pass after a failed run on the same day", problems: []gameSpec{{difficulty: "easy", runs: []gameRun{gameRunOf(4, "2026-10-01"), gameRunOf(5, "2026-10-01")}}}, xp: 8, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "test run without a check", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), runs: []gameRun{{passed: 0, total: 0, on: "2026-10-01"}}}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "three failed runs on one day", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), runs: slices.Repeat([]gameRun{gameRunOf(4, "2026-10-01")}, 3)}}, xp: 14, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "four failed runs on one day", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), runs: slices.Repeat([]gameRun{gameRunOf(4, "2026-10-01")}, 4)}}, xp: 14, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "failed runs on two days", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), runs: slices.Concat(slices.Repeat([]gameRun{gameRunOf(4, "2026-10-01")}, 3), slices.Repeat([]gameRun{gameRunOf(4, "2026-10-02")}, 3))}}, xp: 8, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "failed runs of two problems on one day", problems: slices.Repeat([]gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), runs: slices.Repeat([]gameRun{gameRunOf(4, "2026-10-01")}, 3)}}, 2), xp: 28, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "submitted on the due date", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "submitted", due: "2026-10-01", submitted: "2026-10-01T23:59:00.000Z"}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "submitted the day after the due date", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "submitted", due: "2026-10-01", submitted: "2026-10-02T00:01:00.000Z"}}, xp: 5, late: 1, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "not submitted and overdue", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "assigned", due: "2026-10-01"}}, xp: 5, late: 1, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "not submitted and due today", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "assigned", due: today}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "cancelled and overdue", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "cancelled", due: "2026-10-01"}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "postponed and submitted on the new date", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "submitted", due: "2026-10-01", postpone: true, submitted: "2026-10-02T12:00:00.000Z"}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "added after the due date and not submitted", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "assigned", due: "2026-10-01", added: "2026-10-02T09:00:00.000Z"}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "added on the due date and not submitted", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), task: "assigned", due: "2026-10-01", added: "2026-10-01T09:00:00.000Z"}}, xp: 5, late: 1, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "late penalty below 0 XP", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), task: "assigned", due: "2026-10-01"}}, xp: 0, late: 1, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "visit only today", visits: []int{0}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "two missed days", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help")}}, visits: []int{-3, 0}, xp: 10, missed: 2, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "yesterday missed before the visit of today", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help")}}, visits: []int{-2}, xp: 15, missed: 1, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "missed days below 0 XP", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, visits: []int{-10, 0}, xp: 0, missed: 9, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "six visit days in a row", visits: gameDays(-5, 0), xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "seven visit days in a row", visits: gameDays(-6, 0), xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"week-no-miss"}},
		{name: "seven visit days with a gap", visits: slices.Concat(gameDays(-7, -4), gameDays(-2, 0)), xp: 0, missed: 1, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "seven visit days that ended before", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help")}}, visits: append(gameDays(-9, -3), 0), xp: 10, missed: 2, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept", "week-no-miss"}},
		{name: "twenty-nine visit days in a row", visits: gameDays(-28, 0), xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"week-no-miss"}},
		{name: "thirty visit days in a row", visits: gameDays(-29, 0), xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"week-no-miss", "month-no-miss"}},
		{name: "24 accepted", problems: slices.Repeat([]gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, 24), xp: 240, level: 2, levelXp: 100, nextLevelXp: 300, earned: []string{"first-accept", "ten-accepted"}},
		{name: "25 accepted", problems: slices.Repeat([]gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, 25), xp: 250, level: 2, levelXp: 100, nextLevelXp: 300, earned: []string{"first-accept", "ten-accepted", "twenty-five-accepted"}},
		{name: "49 accepted", problems: slices.Repeat([]gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, 49), xp: 490, level: 3, levelXp: 300, nextLevelXp: 600, earned: []string{"first-accept", "ten-accepted", "twenty-five-accepted"}},
		{name: "50 accepted", problems: slices.Repeat([]gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help")}}, 50), xp: 500, level: 3, levelXp: 300, nextLevelXp: 600, earned: []string{"first-accept", "ten-accepted", "twenty-five-accepted", "fifty-accepted"}},
		{name: "49 correct cards", cards: gameCardDays(49), xp: 98, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "50 correct cards", cards: gameCardDays(50), xp: 100, level: 2, levelXp: 100, nextLevelXp: 300, earned: []string{"fifty-cards"}},
		{name: "50 cards with one wrong", cards: append(gameCardDays(49), gameCard{correct: false, on: "2026-06-01"}), xp: 98, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "failed quiz attempt gives no XP", quizzes: []gameQuizAttempt{gameQuiz(1, 6, 8, false)}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "passed quiz attempt gives 2 XP for each correct answer", quizzes: []gameQuizAttempt{gameQuiz(1, 7, 8, true)}, xp: 14, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "two passed attempts of one deck count the best", quizzes: []gameQuizAttempt{gameQuiz(1, 7, 8, true), gameQuiz(1, 8, 8, true)}, xp: 16, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "a worse second pass keeps the best", quizzes: []gameQuizAttempt{gameQuiz(1, 8, 8, true), gameQuiz(1, 7, 8, true)}, xp: 16, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "passed attempts of two decks add up", quizzes: []gameQuizAttempt{gameQuiz(1, 8, 8, true), gameQuiz(3, 7, 7, true)}, xp: 30, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "49 from passed attempts and old answers", cards: gameCardDays(26), quizzes: []gameQuizAttempt{gameQuiz(1, 8, 8, true), gameQuiz(2, 8, 8, true), gameQuiz(3, 7, 7, true)}, xp: 98, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "50 from passed attempts and old answers", cards: gameCardDays(27), quizzes: []gameQuizAttempt{gameQuiz(1, 8, 8, true), gameQuiz(2, 8, 8, true), gameQuiz(3, 7, 7, true)}, xp: 100, level: 2, levelXp: 100, nextLevelXp: 300, earned: []string{"fifty-cards"}},
		{name: "failed attempt does not count for 50 cards", cards: gameCardDays(49), quizzes: []gameQuizAttempt{gameQuiz(1, 6, 8, false)}, xp: 98, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "nine recalls", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), reviews: gameRecalls(9)}}, xp: 55, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
		{name: "ten recalls", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), reviews: gameRecalls(10)}}, xp: 60, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept", "ten-recalls"}},
		{name: "accepted Codeforces problem", problems: []gameSpec{{difficulty: "medium", attempts: gameAccepted("with_help"), slug: "cf-1538-c"}}, xp: 20, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept", "first-codeforces"}},
		{name: "Codeforces problem not accepted", problems: []gameSpec{{difficulty: "medium", slug: "cf-1538-c", attempts: []gameAttempt{{acceptance: "not_accepted"}}}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "two sets cleared", problems: gameSetProblems(sets, "", "permutations-and-combinations", "two-pointers"), xp: 630, level: 4, levelXp: 600, nextLevelXp: 1000, earned: []string{"first-accept", "first-hard", "ten-accepted", "set-cleared", "twenty-five-accepted", "two-sets", "first-codeforces"}},
		{name: "four tasks on time", problems: slices.Repeat([]gameSpec{gameOnTime(false)}, 4), xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "five tasks on time", problems: slices.Repeat([]gameSpec{gameOnTime(false)}, 5), xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"on-time-five"}},
		{name: "four on time and one late", problems: append(slices.Repeat([]gameSpec{gameOnTime(false)}, 4), gameOnTime(true)), xp: 0, late: 1, level: 1, levelXp: 0, nextLevelXp: 100, earned: nil},
		{name: "level 4 at 960 XP", problems: slices.Repeat([]gameSpec{{difficulty: "hard", attempts: gameAccepted("with_help")}}, 24), xp: 960, level: 4, levelXp: 600, nextLevelXp: 1000, earned: []string{"first-accept", "first-hard", "ten-accepted"}},
		{name: "level 5 at 1000 XP", problems: slices.Repeat([]gameSpec{{difficulty: "hard", attempts: gameAccepted("with_help")}}, 25), xp: 1000, level: 5, levelXp: 1000, nextLevelXp: 1500, earned: []string{"first-accept", "first-hard", "ten-accepted", "twenty-five-accepted", "level-five"}},
		{name: "failed runs below 0 XP", problems: []gameSpec{{difficulty: "easy", attempts: gameAccepted("with_help"), runs: slices.Concat(slices.Repeat([]gameRun{gameRunOf(0, "2026-10-01")}, 3), slices.Repeat([]gameRun{gameRunOf(0, "2026-10-02")}, 3))}}, xp: 0, level: 1, levelXp: 0, nextLevelXp: 100, earned: []string{"first-accept"}},
	}
	// The quizzes list order is the time order of the attempts.
	quizStart := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithSets(t, pool, sets)
			student := loginAs(t, s, "student")
			for _, problem := range tt.problems {
				gameSeed(t, pool, problem)
			}
			for _, card := range tt.cards {
				exec(t, pool, gameInsertCard, newID(), gameBit(card.correct), card.on, card.on)
			}
			for i, quiz := range tt.quizzes {
				exec(t, pool, gameInsertQuiz, newID(), quiz.deck, quiz.correct, quiz.total, gameBit(quiz.passed), quizStart.Add(time.Duration(i)*time.Minute).Format("2006-01-02T15:04:05.000Z"))
			}
			for _, offset := range tt.visits {
				exec(t, pool, gameInsertVisit, gameDay(offset))
			}
			cookie := student
			if tt.role != "" {
				cookie = loginAs(t, s, tt.role)
			}
			w := request(t, s, "/game", http.MethodGet, nil, cookie)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, body %s", w.Code, w.Body.String())
			}
			badges := make([]map[string]any, len(gameOrder))
			for i, id := range gameOrder {
				badges[i] = map[string]any{"id": id, "earned": slices.Contains(tt.earned, id)}
			}
			// json.Marshal sorts the keys of a map, so equal values give equal text.
			got, _ := json.Marshal(decode[any](t, w))
			want, _ := json.Marshal(map[string]any{"xp": tt.xp, "level": tt.level, "levelXp": tt.levelXp, "nextLevelXp": tt.nextLevelXp, "late": tt.late, "missed": tt.missed, "badges": badges})
			if string(got) != string(want) {
				t.Errorf("game =\n%s\nwant\n%s", got, want)
			}
			// COLLATE "C" gives the code unit order of the JavaScript sort.
			stored := queryValue[[]string](t, pool, `SELECT coalesce(array_agg(id ORDER BY id COLLATE "C"),'{}') FROM earned_badges`)
			wantStored := slices.Sorted(slices.Values(slices.DeleteFunc(slices.Clone(tt.earned), func(id string) bool { return id == "level-five" })))
			if !slices.Equal(stored, wantStored) {
				t.Errorf("stored awards = %v, want %v", stored, wantStored)
			}
		})
	}
}

// The layout reads the game on each page change, also without a reload, so a tab that stays open overnight records the new day.
func TestGameRecordsTheVisitDay(t *testing.T) {
	today := gameDay(0)
	tests := []struct {
		name string
		role string
		days []string
	}{
		{name: "the student game records the day", role: "student", days: []string{today}},
		{name: "the parent game records nothing", role: "parent", days: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithSets(t, pool, goldenSets(t))
			loginAs(t, s, "student")
			cookie := loginAs(t, s, tt.role)
			for range 2 {
				if w := request(t, s, "/game", http.MethodGet, nil, cookie); w.Code != http.StatusOK {
					t.Fatalf("status = %d, body %s", w.Code, w.Body.String())
				}
			}
			if days := queryValue[[]string](t, pool, "SELECT coalesce(array_agg(day),'{}') FROM visits"); !slices.Equal(days, tt.days) {
				t.Errorf("visits = %v, want %v", days, tt.days)
			}
		})
	}
}

// Earned badges stay after edits.
func TestGameKeepsTheOnTimeAward(t *testing.T) {
	// An empty dueDate is null.
	tests := []struct {
		name         string
		dueDate      string
		cancel, read bool
	}{
		{name: "removes the due date before the first game read", dueDate: "", cancel: false, read: false},
		{name: "moves the due date before the submissions", dueDate: "2026-09-30", cancel: false, read: false},
		{name: "cancels a submitted task before the first game read", dueDate: "2026-10-01", cancel: true, read: false},
		{name: "keeps an award recorded by an earlier game read", dueDate: "", cancel: false, read: true},
	}
	onTimeFive := gameBadge{ID: "on-time-five", Earned: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithSets(t, pool, goldenSets(t))
			student := loginAs(t, s, "student")
			parent, homeworkID := loginAs(t, s, "parent"), newID()
			taskIDs := make([]string, 5)
			exec(t, pool, gameInsertBundle, homeworkID, "2026-10-01")
			for i := range taskIDs {
				problemID, attemptID := newID(), newID()
				taskIDs[i] = newID()
				exec(t, pool, gameInsertProblem, problemID, problemID, "medium", nil)
				exec(t, pool, gameInsertAttempt, attemptID, problemID, "saved", "not_accepted", "with_help", "2026-10-01T12:00:00.000Z")
				exec(t, pool, gameInsertTask, taskIDs[i], homeworkID, problemID, "submitted", "2026-09-01T00:00:00.000Z")
				exec(t, pool, gameInsertSubmission, newID(), taskIDs[i], attemptID, "2026-10-01T12:00:00.000Z")
			}
			if tt.read {
				if badges := decode[gameView](t, request(t, s, "/game", http.MethodGet, nil, student)).Badges; !slices.Contains(badges, onTimeFive) {
					t.Errorf("first read: badges = %v, want on-time-five earned", badges)
				}
			}
			var w *httptest.ResponseRecorder
			if tt.cancel {
				w = request(t, s, "/tasks/"+taskIDs[0]+"/cancel", http.MethodPost, map[string]any{"version": 1}, parent)
			} else {
				w = request(t, s, "/homework/"+homeworkID, http.MethodPatch, map[string]any{"version": 1, "title": "Week 41", "instructions": "", "dueDate": gameNull(tt.dueDate)}, parent)
			}
			if w.Code != http.StatusOK {
				t.Fatalf("edit: status = %d, body %s", w.Code, w.Body.String())
			}
			for _, cookie := range []string{student, parent} {
				if badges := decode[gameView](t, request(t, s, "/game", http.MethodGet, nil, cookie)).Badges; !slices.Contains(badges, onTimeFive) {
					t.Errorf("badges = %v, want on-time-five earned", badges)
				}
			}
		})
	}
}

// The level badge follows the current XP.
func TestGameKeepsBadgesAfterAProblemEdit(t *testing.T) {
	tests := []struct {
		name        string
		count       int
		read        bool
		levelBefore int
	}{
		{name: "keeps a hard-problem badge before the first game read", count: 1, read: false, levelBefore: 1},
		{name: "keeps earned badges when the current level falls", count: 25, read: true, levelBefore: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithSets(t, pool, goldenSets(t))
			student := loginAs(t, s, "student")
			ids := make([]string, tt.count)
			for i := range ids {
				ids[i] = newID()
				exec(t, pool, gameInsertProblem, ids[i], ids[i], "hard", nil)
				exec(t, pool, gameInsertAttempt, newID(), ids[i], "saved", "accepted", "with_help", gameDay(0))
			}
			if tt.read {
				game := decode[gameView](t, request(t, s, "/game", http.MethodGet, nil, student))
				if game.Level != tt.levelBefore {
					t.Errorf("level before the edit = %d, want %d", game.Level, tt.levelBefore)
				}
				if !slices.Contains(game.Badges, gameBadge{ID: "level-five", Earned: true}) {
					t.Errorf("badges before the edit = %v, want level-five earned", game.Badges)
				}
			}
			w := request(t, s, "/problems/"+ids[0], http.MethodPatch, map[string]any{"version": 1, "url": "https://leetcode.com/problems/" + ids[0] + "/", "number": nil, "title": "Problem", "difficulty": "easy", "topics": []string{}, "summary": ""}, student)
			if w.Code != http.StatusOK {
				t.Fatalf("edit: status = %d, body %s", w.Code, w.Body.String())
			}
			badges := decode[gameView](t, request(t, s, "/game", http.MethodGet, nil, student)).Badges
			if !slices.Contains(badges, gameBadge{ID: "first-hard", Earned: true}) {
				t.Errorf("badges = %v, want first-hard earned", badges)
			}
			if !slices.Contains(badges, gameBadge{ID: "level-five", Earned: false}) {
				t.Errorf("badges = %v, want level-five not earned", badges)
			}
		})
	}
}

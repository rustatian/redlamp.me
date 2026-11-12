package notebook

import (
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The cases of test/quiz.test.ts.

const (
	quizSlug = "permutations-and-combinations-1"
	// quizTime is the format of Date.toISOString.
	quizTime = "2006-01-02T15:04:05.000Z"
)

// quizSetup returns an empty database, a server with the whole content and the cookie of the student, as beforeEach of test/quiz.test.ts.
func quizSetup(t *testing.T) (*pgxpool.Pool, *Server, string) {
	t.Helper()
	pool := database(t)
	s := newServerWithSets(t, pool, goldenSets(t))
	return pool, s, loginAs(t, s, "student")
}

// quizChapters returns the chapters of the set permutations-and-combinations.
func quizChapters(t *testing.T) []sets.Chapter {
	t.Helper()
	for _, item := range goldenSets(t) {
		if item.Slug == "permutations-and-combinations" {
			return item.Chapters
		}
	}
	t.Fatal("the content has no set permutations-and-combinations")
	return nil
}

// quizAgo returns the time the given minutes before now.
func quizAgo(minutes int) string {
	return time.Now().UTC().Add(-time.Duration(minutes) * time.Minute).Format(quizTime)
}

// quizHourAfter returns the time one hour after value.
func quizHourAfter(t *testing.T, value string) string {
	t.Helper()
	moment, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return moment.Add(time.Hour).UTC().Format(quizTime)
}

// quizWrong returns a wrong answer to the card.
func quizWrong(card sets.Card) any {
	if card.Type == "short" {
		return "нет"
	}
	return (card.Answer.(int) + 1) % len(card.Options)
}

// quizAnswers gives the right answer to the first right cards of the chapter and a wrong answer to the other cards.
func quizAnswers(chapters []sets.Chapter, chapter, right int) map[string]any {
	answers := map[string]any{}
	for i, card := range chapters[chapter-1].Cards {
		if i < right {
			answers[card.ID] = card.Answer
		} else {
			answers[card.ID] = quizWrong(card)
		}
	}
	return answers
}

// quizWith returns a copy of the answers with the answer value to the card id.
func quizWith(answers map[string]any, id string, value any) map[string]any {
	changed := maps.Clone(answers)
	changed[id] = value
	return changed
}

// quizWithout returns a copy of the answers without the card id.
func quizWithout(answers map[string]any, id string) map[string]any {
	changed := maps.Clone(answers)
	delete(changed, id)
	return changed
}

// quizSubmit sends an attempt of a chapter of permutations-and-combinations.
func quizSubmit(t *testing.T, s *Server, chapter int, body any, cookie string) *httptest.ResponseRecorder {
	t.Helper()
	return request(t, s, fmt.Sprintf("/quiz/permutations-and-combinations-%d/attempts", chapter), http.MethodPost, body, cookie)
}

// quizRows returns the stored attempts in the order of their time.
func quizRows(t *testing.T, pool *pgxpool.Pool) []map[string]any {
	t.Helper()
	return problemsTableRows(t, pool, "SELECT deck,correct,total,passed FROM quiz_attempts ORDER BY created_at")
}

func TestQuizPassRule(t *testing.T) {
	tests := []struct {
		name    string
		chapter int
		right   int
		total   int
		passed  bool
	}{
		{name: "8 of 8 passes", chapter: 1, right: 8, total: 8, passed: true},
		{name: "7 of 8 passes", chapter: 1, right: 7, total: 8, passed: true},
		{name: "6 of 8 fails", chapter: 1, right: 6, total: 8, passed: false},
		{name: "6 of 7 passes", chapter: 3, right: 6, total: 7, passed: true},
		{name: "5 of 7 fails", chapter: 3, right: 5, total: 7, passed: false},
		{name: "5 of 6 passes", chapter: 4, right: 5, total: 6, passed: true},
		{name: "4 of 6 fails", chapter: 4, right: 4, total: 6, passed: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := quizSetup(t)
			w := quizSubmit(t, s, tt.chapter, map[string]any{"answers": quizAnswers(quizChapters(t), tt.chapter, tt.right)}, student)
			if w.Code != http.StatusCreated {
				t.Errorf("status = %d, want 201; body %s", w.Code, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"correct": tt.right, "total": tt.total, "passed": tt.passed})
			reviewsEqualJSON(t, quizRows(t, pool), []map[string]any{{"deck": fmt.Sprintf("permutations-and-combinations-%d", tt.chapter), "correct": tt.right, "total": tt.total, "passed": gameBit(tt.passed)}})
		})
	}
}

func TestQuizPassedAttemptShowsAnswers(t *testing.T) {
	pool, s, student := quizSetup(t)
	chapters := quizChapters(t)
	w := quizSubmit(t, s, 1, map[string]any{"answers": quizWith(quizAnswers(chapters, 1, 7), "pc1-medals", "  126 ")}, student)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body %s", w.Code, w.Body.String())
	}
	createdAt := queryValue[string](t, pool, "SELECT created_at FROM quiz_attempts")
	results := []any{}
	for _, card := range chapters[0].Cards {
		answer := card.Answer
		if card.Type != "short" {
			answer = card.Options[card.Answer.(int)]
		}
		results = append(results, map[string]any{"cardId": card.ID, "correct": card.ID != "pc1-order-matters", "answer": answer, "explain": card.Explain})
	}
	reviewsEqualJSON(t, decode[any](t, w), map[string]any{"correct": 7, "total": 8, "passed": true, "nextAt": quizHourAfter(t, createdAt), "results": results})
}

func TestQuizFailedAttemptHidesAnswers(t *testing.T) {
	_, s, student := quizSetup(t)
	chapters := quizChapters(t)
	body := decode[map[string]any](t, quizSubmit(t, s, 1, map[string]any{"answers": quizAnswers(chapters, 1, 5)}, student))
	results := []any{}
	for i, card := range chapters[0].Cards {
		results = append(results, map[string]any{"cardId": card.ID, "correct": i < 5})
	}
	reviewsEqualJSON(t, body["results"], results)
}

func TestQuizRefusesAnswers(t *testing.T) {
	all := quizAnswers(quizChapters(t), 1, 8)
	tests := []struct {
		name string
		body any
	}{
		{name: "missing card", body: map[string]any{"answers": quizWithout(all, "pc1-team")}},
		{name: "unknown card id", body: map[string]any{"answers": quizWith(all, "pc-missing", 1)}},
		{name: "card from another deck", body: map[string]any{"answers": quizWith(quizWithout(all, "pc1-team"), "pc2-no-pop", 1)}},
		{name: "option index as text", body: map[string]any{"answers": quizWith(all, "pc1-product-rule", "0")}},
		{name: "option index out of range", body: map[string]any{"answers": quizWith(all, "pc1-product-rule", 4)}},
		{name: "empty short answer", body: map[string]any{"answers": quizWith(all, "pc1-factorial", " ")}},
		{name: "answers not an object", body: map[string]any{"answers": []int{0, 1, 2, 3, 4, 5, 6, 7}}},
		{name: "unknown field", body: map[string]any{"answers": all, "answeredOn": "2026-10-11"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := quizSetup(t)
			if w := quizSubmit(t, s, 1, tt.body, student); w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400; body %s", w.Code, w.Body.String())
			}
			if rows := quizRows(t, pool); len(rows) != 0 {
				t.Errorf("rows = %v, want none", rows)
			}
		})
	}
}

func TestQuizCooldown(t *testing.T) {
	tests := []struct {
		name    string
		deck    string
		minutes int
		correct int
		passed  int
		status  int
	}{
		{name: "failed attempt 59 min ago", deck: quizSlug, minutes: 59, correct: 5, passed: 0, status: 429},
		{name: "a passed attempt also starts the cooldown", deck: quizSlug, minutes: 59, correct: 7, passed: 1, status: 429},
		{name: "attempt 61 min ago", deck: quizSlug, minutes: 61, correct: 5, passed: 0, status: 201},
		{name: "attempt 1 min ago on another deck", deck: "permutations-and-combinations-2", minutes: 1, correct: 5, passed: 0, status: 201},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := quizSetup(t)
			createdAt := quizAgo(tt.minutes)
			exec(t, pool, gameInsertQuiz, "a", tt.deck, tt.correct, 8, tt.passed, createdAt)
			w := quizSubmit(t, s, 1, map[string]any{"answers": quizAnswers(quizChapters(t), 1, 8)}, student)
			if w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			want := 1
			if tt.status == http.StatusCreated {
				want = 2
			}
			if rows := quizRows(t, pool); len(rows) != want {
				t.Errorf("rows = %v, want %d rows", rows, want)
			}
			if tt.status == http.StatusTooManyRequests {
				body := decode[struct {
					Error struct {
						Fields any `json:"fields"`
					} `json:"error"`
				}](t, w)
				reviewsEqualJSON(t, body.Error.Fields, map[string]any{"nextAt": quizHourAfter(t, createdAt)})
			}
		})
	}
}

// D1 runs one request at a time. Here the two submits run at the same time in goroutines.
func TestQuizTwoQuickSubmits(t *testing.T) {
	pool, s, student := quizSetup(t)
	body := map[string]any{"answers": quizAnswers(quizChapters(t), 1, 8)}
	statuses := make([]int, 2)
	var wg sync.WaitGroup
	for i := range statuses {
		wg.Go(func() { statuses[i] = quizSubmit(t, s, 1, body, student).Code })
	}
	wg.Wait()
	slices.Sort(statuses)
	if want := []int{201, 429}; !slices.Equal(statuses, want) {
		t.Errorf("statuses = %v, want %v", statuses, want)
	}
	if rows := quizRows(t, pool); len(rows) != 1 {
		t.Errorf("rows = %v, want 1 row", rows)
	}
}

// previous is nil for no earlier attempt.
func TestQuizHomework(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		setSlug  string
		chapter  int
		right    int
		previous *int
		status   int
		expected string
		version  int
	}{
		{name: "pass completes an assigned quiz", state: "assigned", setSlug: "permutations-and-combinations", chapter: 1, right: 8, previous: nil, status: 201, expected: "completed", version: 2},
		{name: "pass completes a started quiz", state: "in_progress", setSlug: "permutations-and-combinations", chapter: 1, right: 7, previous: nil, status: 201, expected: "completed", version: 2},
		{name: "failure keeps the quiz assigned", state: "assigned", setSlug: "permutations-and-combinations", chapter: 1, right: 6, previous: nil, status: 201, expected: "assigned", version: 1},
		{name: "blocked pass keeps the quiz assigned", state: "assigned", setSlug: "permutations-and-combinations", chapter: 1, right: 8, previous: new(30), status: 429, expected: "assigned", version: 1},
		{name: "old pass does not complete a new assignment after failure", state: "assigned", setSlug: "permutations-and-combinations", chapter: 1, right: 6, previous: new(120), status: 201, expected: "assigned", version: 1},
		{name: "pass leaves another chapter assigned", state: "assigned", setSlug: "permutations-and-combinations", chapter: 2, right: 8, previous: nil, status: 201, expected: "assigned", version: 1},
		{name: "pass leaves another set assigned", state: "assigned", setSlug: "two-pointers", chapter: 1, right: 8, previous: nil, status: 201, expected: "assigned", version: 1},
		{name: "pass keeps a cancelled assignment cancelled", state: "cancelled", setSlug: "permutations-and-combinations", chapter: 1, right: 8, previous: nil, status: 201, expected: "cancelled", version: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := quizSetup(t)
			parent := loginAs(t, s, "parent")
			homeworkID := decode[struct {
				ID string `json:"id"`
			}](t, request(t, s, "/homework", http.MethodPost, map[string]any{"title": "Quiz homework", "instructions": "", "dueDate": nil}, parent)).ID
			if w := request(t, s, "/homework/"+homeworkID+"/content", http.MethodPost, map[string]any{"kind": "quiz", "setSlug": tt.setSlug, "chapter": tt.chapter}, parent); w.Code != http.StatusCreated {
				t.Fatalf("add quiz: status %d, want 201; body %s", w.Code, w.Body.String())
			}
			exec(t, pool, "UPDATE homework_tasks SET state=$1 WHERE homework_id=$2", tt.state, homeworkID)
			if tt.previous != nil {
				exec(t, pool, gameInsertQuiz, "old", quizSlug, 8, 8, 1, quizAgo(*tt.previous))
			}
			before := time.Now().UTC().Format(quizTime)
			if w := quizSubmit(t, s, 1, map[string]any{"answers": quizAnswers(quizChapters(t), 1, tt.right)}, student); w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			tasks := homeworkTasks(t, s, homeworkID, student)
			if len(tasks) != 1 {
				t.Fatalf("tasks = %v, want 1 task", tasks)
			}
			problemsMatchBody(t, tasks[0], map[string]any{"state": tt.expected, "version": tt.version})
			completedAt, ok := tasks[0]["completedAt"]
			if tt.expected == "completed" {
				text, _ := completedAt.(string)
				if after := time.Now().UTC().Format(quizTime); text < before || text > after {
					t.Errorf("completedAt = %v, want a time from %s to %s", completedAt, before, after)
				}
			} else if !ok || completedAt != nil {
				t.Errorf("completedAt = %v, want null", completedAt)
			}
		})
	}
}

// minutes is nil for no attempt.
func TestQuizRound(t *testing.T) {
	tests := []struct {
		name    string
		minutes *int
	}{
		{name: "no attempt", minutes: nil},
		{name: "attempt 30 min ago", minutes: new(30)},
		{name: "attempt 61 min ago", minutes: new(61)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, s, student := quizSetup(t)
			var nextAt any
			if tt.minutes != nil {
				createdAt := quizAgo(*tt.minutes)
				exec(t, pool, gameInsertQuiz, "a", quizSlug, 5, 8, 0, createdAt)
				if *tt.minutes == 30 {
					nextAt = quizHourAfter(t, createdAt)
				}
			}
			round := decode[map[string]any](t, request(t, s, "/quiz/"+quizSlug, http.MethodGet, nil, student))
			if want := "Permutations и Combinations, глава 1"; round["title"] != want {
				t.Errorf("title = %v, want %s", round["title"], want)
			}
			reviewsEqualJSON(t, round["chapter"], map[string]any{"set": "permutations-and-combinations", "setTitle": "Permutations и Combinations", "number": 1, "title": "Считаем варианты"})
			cards, _ := round["cards"].([]any)
			ids := []string{}
			for _, item := range cards {
				card, _ := item.(map[string]any)
				id, _ := card["id"].(string)
				ids = append(ids, id)
				_, answer := card["answer"]
				_, explain := card["explain"]
				if answer || explain {
					t.Errorf("card %s has the answer or the explanation", id)
				}
			}
			want := []string{}
			for _, card := range quizChapters(t)[0].Cards {
				want = append(want, card.ID)
			}
			if !slices.Equal(ids, want) {
				t.Errorf("card ids = %q, want %q", ids, want)
			}
			if got, ok := round["nextAt"]; !ok || got != nextAt {
				t.Errorf("nextAt = %v, want %v", got, nextAt)
			}
		})
	}
}

// The body of a case is nil for a request without a body.
func TestQuizAccess(t *testing.T) {
	all := quizAnswers(quizChapters(t), 1, 8)
	tests := []struct {
		name   string
		role   string
		path   string
		method string
		body   any
		status int
	}{
		{name: "parent cannot submit", role: "parent", path: "/quiz/" + quizSlug + "/attempts", method: http.MethodPost, body: map[string]any{"answers": all}, status: 403},
		{name: "parent reads the cards", role: "parent", path: "/quiz/" + quizSlug + "/cards", method: http.MethodGet, status: 200},
		{name: "student cannot read the cards", role: "student", path: "/quiz/" + quizSlug + "/cards", method: http.MethodGet, status: 403},
		{name: "unknown deck", role: "student", path: "/quiz/two-sum", method: http.MethodGet, status: 404},
		{name: "attempt of an unknown deck", role: "student", path: "/quiz/two-sum/attempts", method: http.MethodPost, body: map[string]any{"answers": all}, status: 404},
		{name: "set slug without a chapter number is not a deck", role: "student", path: "/quiz/permutations-and-combinations", method: http.MethodGet, status: 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, s, student := quizSetup(t)
			cookie := student
			if tt.role == "parent" {
				cookie = loginAs(t, s, "parent")
			}
			w := request(t, s, tt.path, tt.method, tt.body, cookie)
			if w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			if strings.HasSuffix(tt.path, "/cards") && tt.status == http.StatusOK {
				reviewsEqualJSON(t, decode[any](t, w), quizChapters(t)[0].Cards)
			}
		})
	}
}

func TestQuizDecks(t *testing.T) {
	pool, s, student := quizSetup(t)
	failedAt := quizAgo(30)
	exec(t, pool, gameInsertQuiz, "a", quizSlug, 7, 8, 1, quizAgo(180))
	exec(t, pool, gameInsertQuiz, "b", quizSlug, 5, 8, 0, failedAt)
	exec(t, pool, gameInsertQuiz, "c", "permutations-and-combinations-2", 8, 8, 1, quizAgo(120))
	exec(t, pool, gameInsertQuiz, "d", "permutations-and-combinations-2", 7, 8, 1, quizAgo(90))
	// The deck had 8 cards at this pass and has 7 now.
	exec(t, pool, gameInsertQuiz, "e", "permutations-and-combinations-3", 8, 8, 1, quizAgo(240))
	exec(t, pool, gameInsertQuiz, "f", "permutations-and-combinations-4", 4, 6, 0, quizAgo(120))
	exec(t, pool, gameInsertQuiz, "g", "two-pointers-2", 8, 8, 1, quizAgo(240))
	exec(t, pool, gameInsertQuiz, "h", "heaps-and-queues-1", 7, 8, 1, quizAgo(240))
	want := []any{
		map[string]any{"slug": "complexity-1", "title": "Как понять O(n): время и память, глава 1", "chapter": map[string]any{"set": "complexity", "setTitle": "Как понять O(n): время и память", "number": 1, "title": "Как понять O(n): время и память"}, "cardCount": 6, "best": nil, "bestTotal": nil, "nextAt": nil},
		map[string]any{"slug": quizSlug, "title": "Permutations и Combinations, глава 1", "chapter": map[string]any{"set": "permutations-and-combinations", "setTitle": "Permutations и Combinations", "number": 1, "title": "Считаем варианты"}, "cardCount": 8, "best": 7, "bestTotal": 8, "nextAt": quizHourAfter(t, failedAt)},
		map[string]any{"slug": "permutations-and-combinations-2", "title": "Permutations и Combinations, глава 2", "chapter": map[string]any{"set": "permutations-and-combinations", "setTitle": "Permutations и Combinations", "number": 2, "title": "Backtracking: Subsets, Permutations, Combinations"}, "cardCount": 8, "best": 8, "bestTotal": 8, "nextAt": nil},
		map[string]any{"slug": "permutations-and-combinations-3", "title": "Permutations и Combinations, глава 3", "chapter": map[string]any{"set": "permutations-and-combinations", "setTitle": "Permutations и Combinations", "number": 3, "title": "Повторы и Pruning"}, "cardCount": 7, "best": 8, "bestTotal": 8, "nextAt": nil},
		map[string]any{"slug": "permutations-and-combinations-4", "title": "Permutations и Combinations, глава 4", "chapter": map[string]any{"set": "permutations-and-combinations", "setTitle": "Permutations и Combinations", "number": 4, "title": "Порядок и номер"}, "cardCount": 6, "best": nil, "bestTotal": nil, "nextAt": nil},
		map[string]any{"slug": "two-pointers-1", "title": "Two Pointers, глава 1", "chapter": map[string]any{"set": "two-pointers", "setTitle": "Two Pointers", "number": 1, "title": "Навстречу"}, "cardCount": 8, "best": nil, "bestTotal": nil, "nextAt": nil},
		map[string]any{"slug": "two-pointers-2", "title": "Two Pointers, глава 2", "chapter": map[string]any{"set": "two-pointers", "setTitle": "Two Pointers", "number": 2, "title": "Читатель и писатель"}, "cardCount": 7, "best": 8, "bestTotal": 8, "nextAt": nil},
		map[string]any{"slug": "two-pointers-3", "title": "Two Pointers, глава 3", "chapter": map[string]any{"set": "two-pointers", "setTitle": "Two Pointers", "number": 3, "title": "Две последовательности"}, "cardCount": 8, "best": nil, "bestTotal": nil, "nextAt": nil},
		map[string]any{"slug": "two-pointers-4", "title": "Two Pointers, глава 4", "chapter": map[string]any{"set": "two-pointers", "setTitle": "Two Pointers", "number": 4, "title": "Карта Two Pointers"}, "cardCount": 8, "best": nil, "bestTotal": nil, "nextAt": nil},
		map[string]any{"slug": "heaps-and-queues-1", "title": "Heap (Priority Queue), Queue и Deque, глава 1", "chapter": map[string]any{"set": "heaps-and-queues", "setTitle": "Heap (Priority Queue), Queue и Deque", "number": 1, "title": "Queue и Deque"}, "cardCount": 8, "best": 7, "bestTotal": 8, "nextAt": nil},
		map[string]any{"slug": "heaps-and-queues-2", "title": "Heap (Priority Queue), Queue и Deque, глава 2", "chapter": map[string]any{"set": "heaps-and-queues", "setTitle": "Heap (Priority Queue), Queue и Deque", "number": 2, "title": "Binary Heap: приоритет и стоимость"}, "cardCount": 8, "best": nil, "bestTotal": nil, "nextAt": nil},
		map[string]any{"slug": "heaps-and-queues-3", "title": "Heap (Priority Queue), Queue и Deque, глава 3", "chapter": map[string]any{"set": "heaps-and-queues", "setTitle": "Heap (Priority Queue), Queue и Deque", "number": 3, "title": "Sliding Window и Monotonic Queue"}, "cardCount": 8, "best": nil, "bestTotal": nil, "nextAt": nil},
		map[string]any{"slug": "heaps-and-queues-4", "title": "Heap (Priority Queue), Queue и Deque, глава 4", "chapter": map[string]any{"set": "heaps-and-queues", "setTitle": "Heap (Priority Queue), Queue и Deque", "number": 4, "title": "Heap и Greedy"}, "cardCount": 8, "best": nil, "bestTotal": nil, "nextAt": nil},
	}
	reviewsEqualJSON(t, decode[any](t, request(t, s, "/quiz", http.MethodGet, nil, student)), want)
}

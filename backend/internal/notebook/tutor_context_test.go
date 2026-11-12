package notebook

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode/utf16"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

const (
	tutorContextProblem       = "11111111-1111-4111-8111-111111111111"
	tutorContextHomework      = "77777777-7777-4777-8777-777777777777"
	tutorContextInsertProblem = `INSERT INTO problems (id,slug,number,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,'two-sum',1,'Two Sum','easy','["Array"]','Find two numbers.','','now','now')`
	tutorContextInsertMessage = "INSERT INTO tutor_messages (id,context_kind,context_ref,author,body,quote,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	tutorContextInsertAttempt = "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at,saved_at) VALUES ($1,$2,$3,$4,'','not_submitted','needs_practice','now','now',$5)"
)

// tutorContextLessons are the lessons of the four chapters of tutorSet, as contextLessons of test/tutor.test.ts.
var tutorContextLessons = []string{
	"<p>Вводный пример.</p><h2>1. Linear Scan</h2><p>Просмотрите Array.</p>",
	"<p>Backtracking.</p><h2>1. Выбор</h2><p>Добавьте элемент.</p><h2>2. Call Stack</h2><p>Call Stack хранит незавершённые вызовы.</p>",
	"<p>Pruning.</p><h2>1. Pruning</h2><p>Закончите ветку.</p>",
	"<p>Порядок.</p><h2>1. Next Permutation</h2><p>Найдите следующий порядок.</p>",
}

var (
	tutorContextAnswered     = []any{map[string]any{"delta": "Hi"}, map[string]any{"done": true}}
	tutorContextBlockPattern = regexp.MustCompile("(?s)^<context>\n.*\n</context>$")
)

// tutorContextMinute is the time of minute i on 2026-10-09 in the ISO format of JavaScript, as minute() of the TypeScript test.
func tutorContextMinute(i int) string {
	return time.Date(2026, time.October, 9, 10, i, 0, 0, time.UTC).Format("2006-01-02T15:04:05.000Z")
}

// tutorContextSets returns the learning sets of the Worker with tutorContextLessons as the lessons of tutorSet.
// Decided difference: the TypeScript test changes the shared sets and restores them. goldenSets decodes a new value on each call, so the change stays in this value.
func tutorContextSets(t *testing.T) []sets.LearningSet {
	t.Helper()
	content := goldenSets(t)
	for i := range content {
		if content[i].Slug == tutorSet {
			for j := range content[i].Chapters {
				content[i].Chapters[j].Lesson = tutorContextLessons[j]
			}
		}
	}
	return content
}

// tutorContextWork is an attempt document with one approach, as the work of the TypeScript test.
func tutorContextWork(notes, code string) map[string]any {
	return map[string]any{"notes": notes, "approaches": []any{map[string]any{
		"id": "33333333-3333-4333-8333-333333333333", "label": "Hash map", "idea": "Remember seen values.", "correctness": "Each pair is checked once.",
		"timeComplexity": "O(n)", "spaceComplexity": "O(n)", "edgeCases": "Empty input", "mistakes": "", "language": "python", "code": code,
	}}}
}

// tutorContextJSON returns the JSON text of value.
func tutorContextJSON(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// tutorContextStatement is a LeetCode reply with the statement content. A nil content is the reply for a paid-only problem.
func tutorContextStatement(content any) func(string, string) (json.RawMessage, error) {
	return func(string, string) (json.RawMessage, error) { return json.Marshal(map[string]any{"content": content}) }
}

// tutorContextBusy is a LeetCode reply with the status 500.
func tutorContextBusy(string, string) (json.RawMessage, error) {
	return nil, errors.New("LeetCode returned 500")
}

// tutorContextSaid is the text of a turn with an optional quote, as said() of src/leetcode/tutor.ts. An empty quote is no quote.
func tutorContextSaid(message, quote string) string {
	if quote == "" {
		return message
	}
	return "Selected text:\n> " + quote + "\n\n" + message
}

// tutorContextMessage is one message of the Messages API with one text block.
// Decided difference: the Responses API item {role, content: text} becomes {role, content: [{type: "text", text}]}.
func tutorContextMessage(role, text string) map[string]any {
	return map[string]any{"role": role, "content": []any{map[string]any{"type": "text", "text": text}}}
}

// tutorContextBody returns the body of the one request to the fake upstream and the text of the message before the question: the context block.
func tutorContextBody(t *testing.T, u *tutorUpstream) (map[string]any, string) {
	t.Helper()
	sent := u.requests()
	if len(sent) != 1 {
		t.Fatalf("%d requests to the upstream, want 1", len(sent))
	}
	var body struct {
		Messages []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"messages"`
	}
	text := tutorContextJSON(t, sent[0].body)
	if err := json.Unmarshal([]byte(text), &body); err != nil {
		t.Fatal(err)
	}
	n := len(body.Messages)
	if n < 2 || len(body.Messages[n-2].Content) != 1 {
		t.Fatalf("body = %s, want a context block before the question", text)
	}
	return sent[0].body, body.Messages[n-2].Content[0].Text
}

// A role or an input that the server does not accept gets an error before the upstream and LeetCode get a request.
func TestTutorAccess(t *testing.T) {
	thread := "/tutor/messages?kind=problem&id=" + tutorContextProblem
	tests := []struct {
		name       string
		role       string
		method     string
		path       string
		body       any
		wantStatus int
	}{
		{name: "parent cannot send", role: "parent", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(nil), wantStatus: 403},
		{name: "student reads a thread", role: "student", method: http.MethodGet, path: thread, wantStatus: 200},
		{name: "parent cannot read a thread", role: "parent", method: http.MethodGet, path: thread, wantStatus: 403},
		{name: "parent cannot list the threads", role: "parent", method: http.MethodGet, path: "/tutor/threads", wantStatus: 404},
		{name: "unknown problem", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"kind": "problem", "id": newID()}), wantStatus: 404},
		{name: "lesson kind", role: "student", method: http.MethodGet, path: "/tutor/messages?kind=lesson&id=" + newID(), wantStatus: 400},
		{name: "unknown set", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"kind": "set", "id": "missing-set"}), wantStatus: 404},
		{name: "unknown kind", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"kind": "homework"}), wantStatus: 400},
		{name: "empty message", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"message": " "}), wantStatus: 400},
		{name: "message over 4000 characters", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"message": strings.Repeat("m", 4001)}), wantStatus: 400},
		{name: "quote over 2000 characters", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"quote": strings.Repeat("q", 2001)}), wantStatus: 400},
		{name: "extra field", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"model": "x"}), wantStatus: 400},
		{name: "parent cannot evaluate", role: "parent", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"kind": "problem", "id": tutorContextProblem, "evaluate": true}), wantStatus: 403},
		{name: "evaluate a set", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"evaluate": true}), wantStatus: 400},
		{name: "evaluate flag that is not true", role: "student", method: http.MethodPost, path: "/tutor/messages", body: tutorSend(map[string]any{"kind": "problem", "id": tutorContextProblem, "evaluate": "yes"}), wantStatus: 400},
	}
	content := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			exec(t, pool, tutorContextInsertProblem, tutorContextProblem)
			exec(t, pool, tutorContextInsertMessage, newID(), "problem", tutorContextProblem, "student", "Why?", nil, tutorContextMinute(0))
			exec(t, pool, tutorContextInsertMessage, newID(), "problem", tutorContextProblem, "assistant", "Look at the base case.", nil, tutorContextMinute(1))
			u := newTutorUpstream(t)
			leetcode := &fakeLeetCode{}
			s := tutorServer(t, Config{DB: pool, Sets: content, LeetCode: leetcode}, u)
			w := request(t, s, tt.path, tt.method, tt.body, loginAs(t, s, tt.role))
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			body := decode[map[string]any](t, w)
			if tt.wantStatus == http.StatusOK {
				// The ID of an item is any string.
				items, _ := body["items"].([]any)
				for _, item := range items {
					if object, ok := item.(map[string]any); ok {
						if _, ok := object["id"].(string); !ok {
							t.Errorf("id = %v, want a string", object["id"])
						}
						delete(object, "id")
					}
				}
				reviewsEqualJSON(t, body, map[string]any{"items": []any{
					map[string]any{"author": "assistant", "body": "Look at the base case.", "quote": nil, "createdAt": tutorContextMinute(1)},
					map[string]any{"author": "student", "body": "Why?", "quote": nil, "createdAt": tutorContextMinute(0)},
				}, "nextOffset": nil})
			} else {
				fields, _ := body["error"].(map[string]any)
				_, code := fields["code"].(string)
				_, message := fields["message"].(string)
				if len(body) != 1 || len(fields) != 2 || !code || !message {
					t.Errorf("body = %s, want an error with a code and a message", w.Body.String())
				}
			}
			if got := len(u.requests()) + len(leetcode.calls); got != 0 {
				t.Errorf("%d requests to the upstream and LeetCode, want none", got)
			}
			if got := len(tutorStored(t, pool)); got != 2 {
				t.Errorf("%d stored messages, want 2", got)
			}
		})
	}
}

// The upstream gets the newest turns of the thread that fit the limits, then the context block, then the question.
func TestTutorHistory(t *testing.T) {
	tests := []struct {
		name string
		// history is the number of stored turns. The turns alternate from the student, and the body of turn i is "m<i>" filled with dots to size characters.
		history int
		size    int
		// quote is the quote of the question. historyQuote is the quote of the first stored turn. An empty quote is no quote.
		quote        string
		historyQuote string
		// wantCount is the number of the newest stored turns in the request.
		wantCount int
	}{
		{name: "fixed fields", history: 0, size: 2, wantCount: 0},
		{name: "order", history: 2, size: 2, historyQuote: "A tree", wantCount: 2},
		{name: "quote format", history: 0, size: 2, quote: "Each circle is one call", wantCount: 0},
		{name: "history limit", history: 30, size: 2, wantCount: 20},
		// Decided difference: the TypeScript case keeps 2 turns, m1 of the assistant and m2 of the student.
		// The Messages API takes a user turn first, so the Go server drops m1 and keeps m2.
		{name: "history over the character budget", history: 3, size: 7000, wantCount: 1},
	}
	content := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			type turn struct{ author, body, quote string }
			seeded := make([]turn, tt.history)
			for i := range seeded {
				seeded[i] = turn{author: "student", body: fmt.Sprintf("m%d", i)}
				if i%2 == 1 {
					seeded[i].author = "assistant"
				}
				if i == 0 {
					seeded[i].quote = tt.historyQuote
				}
				seeded[i].body += strings.Repeat(".", max(0, tt.size-len(seeded[i].body)))
				var quote any
				if seeded[i].quote != "" {
					quote = seeded[i].quote
				}
				exec(t, pool, tutorContextInsertMessage, newID(), "set", tutorSet, seeded[i].author, seeded[i].body, quote, tutorContextMinute(i))
			}
			u := newTutorUpstream(t, tutorAnswer())
			s := tutorServer(t, Config{DB: pool, Sets: content}, u)
			var quote any
			if tt.quote != "" {
				quote = tt.quote
			}
			request(t, s, "/tutor/messages", http.MethodPost, tutorSend(map[string]any{"quote": quote}), loginAs(t, s, "student"))
			body, block := tutorContextBody(t, u)
			// Decided difference: the keys of the Messages API body. TestTutorRequest checks their values and the headers.
			if got, want := slices.Sorted(maps.Keys(body)), []string{"fallbacks", "max_tokens", "messages", "model", "output_config", "stream", "system", "thinking"}; !slices.Equal(got, want) {
				t.Errorf("keys = %v, want %v", got, want)
			}
			if !tutorContextBlockPattern.MatchString(block) {
				t.Errorf("context block = %q, want the text in <context> tags", block)
			}
			want := []any{}
			for _, item := range seeded[tt.history-tt.wantCount:] {
				role := "user"
				if item.author == "assistant" {
					role = "assistant"
				}
				want = append(want, tutorContextMessage(role, tutorContextSaid(item.body, item.quote)))
			}
			want = append(want, tutorContextMessage("user", block), tutorContextMessage("user", tutorContextSaid("Explain this part.", tt.quote)))
			reviewsEqualJSON(t, body["messages"], want)
		})
	}
}

// The Messages API takes a user turn first and takes no text block without text, so the request leaves out the stored turns that break these rules.
func TestTutorHistoryRules(t *testing.T) {
	type turn struct{ author, body string }
	// later are 9 questions with their answers: 18 turns.
	var later []turn
	for i := range 9 {
		later = append(later, turn{author: "student", body: fmt.Sprintf("question %d", i)}, turn{author: "assistant", body: fmt.Sprintf("answer %d", i)})
	}
	tests := []struct {
		name string
		// stored are the stored turns from the oldest to the newest.
		stored []turn
		// want are the stored turns of the request.
		want []turn
	}{
		// Two questions at the same time store their answers after the two questions. The newest 20 turns start with the two answers.
		{name: "two answers at the start of the newest 20 turns", stored: slices.Concat([]turn{{author: "student", body: "first"}, {author: "student", body: "second"}, {author: "assistant", body: "to the first"}, {author: "assistant", body: "to the second"}}, later), want: later},
		{name: "answer with an empty body", stored: []turn{{author: "student", body: "first"}, {author: "assistant", body: ""}, {author: "student", body: "second"}, {author: "assistant", body: "to the second"}}, want: []turn{{author: "student", body: "first"}, {author: "student", body: "second"}, {author: "assistant", body: "to the second"}}},
		{name: "answer with only white space", stored: []turn{{author: "student", body: "first"}, {author: "assistant", body: " \n"}, {author: "student", body: "second"}, {author: "assistant", body: "to the second"}}, want: []turn{{author: "student", body: "first"}, {author: "student", body: "second"}, {author: "assistant", body: "to the second"}}},
	}
	content := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			for i, item := range tt.stored {
				exec(t, pool, tutorContextInsertMessage, newID(), "set", tutorSet, item.author, item.body, nil, tutorContextMinute(i))
			}
			u := newTutorUpstream(t, tutorAnswer())
			s := tutorServer(t, Config{DB: pool, Sets: content}, u)
			request(t, s, "/tutor/messages", http.MethodPost, tutorSend(nil), loginAs(t, s, "student"))
			body, block := tutorContextBody(t, u)
			want := []any{}
			for _, item := range tt.want {
				role := "user"
				if item.author == "assistant" {
					role = "assistant"
				}
				want = append(want, tutorContextMessage(role, item.body))
			}
			want = append(want, tutorContextMessage("user", block), tutorContextMessage("user", "Explain this part."))
			reviewsEqualJSON(t, body["messages"], want)
		})
	}
}

// The evaluate flag selects the system prompt. The context has the normalized topics and the notes, and the attempt does not change.
func TestTutorPrompt(t *testing.T) {
	tests := []struct {
		name       string
		data       map[string]any
		wantSystem string
	}{
		{name: "tutor prompt without the evaluate flag", data: map[string]any{"kind": "problem", "id": tutorContextProblem}, wantSystem: tutorPrompt},
		{name: "evaluate prompt with the evaluate flag", data: map[string]any{"kind": "problem", "id": tutorContextProblem, "evaluate": true}, wantSystem: evaluatePrompt},
	}
	content := goldenSets(t)
	work := tutorContextWork("Хеш-таблица: моя заметка.", "return seen")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			exec(t, pool, tutorContextInsertProblem, tutorContextProblem)
			exec(t, pool, "UPDATE problems SET topics=$1 WHERE id=$2", `["HashTable","Хеш-таблица","Array","Hash Table"]`, tutorContextProblem)
			exec(t, pool, tutorContextInsertAttempt, newID(), tutorContextProblem, "draft", tutorContextJSON(t, work), nil)
			u := newTutorUpstream(t, tutorAnswer())
			s := tutorServer(t, Config{DB: pool, Sets: content, LeetCode: &fakeLeetCode{reply: tutorContextStatement("<p>Text</p>")}}, u)
			w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(tt.data), loginAs(t, s, "student"))
			if got := tutorEvents(t, w); !reflect.DeepEqual(got, tutorContextAnswered) {
				t.Errorf("events = %v, want %v", got, tutorContextAnswered)
			}
			body, block := tutorContextBody(t, u)
			// Decided difference: the Messages API has the prompt in system, in place of instructions.
			reviewsEqualJSON(t, body["system"], []any{map[string]any{"type": "text", "text": tt.wantSystem}})
			for _, text := range []string{"Topics: Hash Table, Array\n", "Notes: Хеш-таблица: моя заметка."} {
				if !strings.Contains(block, text) {
					t.Errorf("context block = %q, want %q in it", block, text)
				}
			}
			var stored any
			if err := json.Unmarshal([]byte(queryValue[string](t, pool, "SELECT document FROM attempts WHERE problem_id=$1", tutorContextProblem)), &stored); err != nil {
				t.Fatal(err)
			}
			reviewsEqualJSON(t, stored, work)
		})
	}
}

// tutorContextAttempt is a stored attempt of the problem. savedAt is nil for no time.
type tutorContextAttempt struct {
	state   string
	code    string
	savedAt any
}

// The context block has the problem, the statement, the homework, the chapter of the problem and the newest work, or the learning set.
func TestTutorContext(t *testing.T) {
	tests := []struct {
		name string
		kind string
		// task adds an in-progress homework task of the problem. slug is the slug of the problem in place of two-sum.
		task bool
		slug string
		// lessons gives tutorSet the lessons tutorContextLessons.
		lessons bool
		work    []tutorContextAttempt
		// statement is the LeetCode reply. A nil statement fails the request.
		statement func(string, string) (json.RawMessage, error)
		quote     string
		// leetcode is true when the server must ask LeetCode for the statement, as the leetcode flag of the TypeScript test.
		leetcode     bool
		wantIncludes []string
		// wantCutDraft is true for a block that ends with 12000 characters of the draft and the cut mark.
		wantCutDraft bool
		wantExcludes []string
	}{
		{name: "task instructions and code", kind: "problem", task: true, work: []tutorContextAttempt{{"draft", "return seen", nil}}, statement: tutorContextStatement("<p>Given an array <code>nums</code>.</p>"), leetcode: true,
			wantIncludes: []string{"Problem: 1. Two Sum (easy)", "LeetCode statement:\nGiven an array nums.", "Homework instructions:\nUse a hash map.", "Week 41", "return seen"}},
		{name: "newest saved attempt", kind: "problem", work: []tutorContextAttempt{{"saved", "old code", "2026-10-01T10:00:00.000Z"}, {"saved", "new code", "2026-10-02T10:00:00.000Z"}}, statement: tutorContextStatement("<p>Text</p>"), leetcode: true,
			wantIncludes: []string{"new code"}, wantExcludes: []string{"old code"}},
		{name: "draft before a newer saved attempt", kind: "problem", work: []tutorContextAttempt{{"draft", "draft code", nil}, {"saved", "saved code", "2026-10-02T10:00:00.000Z"}}, statement: tutorContextStatement("<p>Text</p>"), leetcode: true,
			wantIncludes: []string{"draft code"}, wantExcludes: []string{"saved code"}},
		{name: "statement fetch fails", kind: "problem", statement: tutorContextBusy, leetcode: true,
			wantIncludes: []string{"Problem: 1. Two Sum (easy)"}, wantExcludes: []string{"LeetCode statement"}},
		{name: "paid-only statement", kind: "problem", statement: tutorContextStatement(nil), leetcode: true,
			wantIncludes: []string{"Problem: 1. Two Sum (easy)"}, wantExcludes: []string{"LeetCode statement"}},
		{name: "Codeforces problem has no LeetCode statement", kind: "problem", slug: "cf-1-a",
			wantIncludes: []string{"Problem: 1. Two Sum (easy)"}, wantExcludes: []string{"LeetCode statement"}},
		{name: "local Codeforces statement supplies constraints and input/output", kind: "problem", slug: "cf-1526-c2",
			wantIncludes: []string{"Codeforces statement:\n", "В первой строке целое n.", "Выведите наибольшее число принятых зелий.", "1 ≤ n ≤ 200000", "-10^(9) ≤ ai ≤ 10^(9)"}, wantExcludes: []string{"LeetCode statement"}},
		{name: "local LeetCode statement wins over the remote statement", kind: "problem", slug: "subsets", statement: tutorContextStatement("<p>Remote statement</p>"),
			wantIncludes: []string{"LeetCode statement:\nДан Array целых чисел nums.", "1 ≤ nums.length ≤ 10"}, wantExcludes: []string{"Remote statement"}},
		{name: "chapter of the problem after the homework", kind: "problem", slug: "subsets", lessons: true, task: true, work: []tutorContextAttempt{{"draft", "return seen", nil}}, statement: tutorContextStatement("<p>Text</p>"),
			wantIncludes: []string{"Homework instructions:\nUse a hash map.\n\nChapter: Permutations и Combinations, глава 2. Backtracking: Subsets, Permutations, Combinations\nLesson outline:\n## 1. Выбор\n## 2. Call Stack\n\nLesson text (sections 1-3 of 3).\nBacktracking.\n## 1. Выбор\nДобавьте элемент.\n## 2. Call Stack\nCall Stack хранит незавершённые вызовы.\n\nCurrent draft:\n"},
			wantExcludes: []string{"Просмотрите Array.", "Learning set", "The rest is cut"}},
		{name: "quote selects a chapter section", kind: "problem", slug: "subsets", lessons: true, statement: tutorContextStatement("<p>Text</p>"), quote: "Call Stack  хранит\nнезавершённые вызовы.",
			wantIncludes: []string{"Lesson text (sections 3-3 of 3). The rest is cut.\n## 2. Call Stack\nCall Stack хранит незавершённые вызовы."}, wantExcludes: []string{"Добавьте элемент."}},
		{name: "problem in no chapter", kind: "problem", statement: tutorContextStatement("<p>Text</p>"), leetcode: true,
			wantIncludes: []string{"Problem: 1. Two Sum (easy)"}, wantExcludes: []string{"Chapter:", "Lesson"}},
		// The TypeScript pattern is /Current draft:\n[\s\S]{12000} \[cut\]\n<\/context>$/. RE2 takes no repeat count over 1000, so wantCutDraft checks it.
		{name: "field over its cap", kind: "problem", work: []tutorContextAttempt{{"draft", strings.Repeat("x", 13_000), nil}}, statement: tutorContextStatement(nil), leetcode: true, wantCutDraft: true},
		{name: "quote selects a section", kind: "set", lessons: true, quote: "Call Stack хранит незавершённые вызовы.",
			wantIncludes: []string{"Learning set: Permutations и Combinations", "Lesson text (sections 4-6 of 6). The rest is cut.\n## 2. Call Stack\nCall Stack хранит незавершённые вызовы.", "Закончите ветку.", "Найдите следующий порядок."},
			wantExcludes: []string{"Просмотрите Array.", "Добавьте элемент."}},
		{name: "set without a quote", kind: "set", lessons: true,
			wantIncludes: []string{"Learning set: Permutations и Combinations", "Lesson outline:\n# Глава 1. Считаем варианты\n## 1. Linear Scan\n", "\n# Глава 2. Backtracking: Subsets, Permutations, Combinations\n## 1. Выбор\n## 2. Call Stack\n", "Lesson text (sections 1-6 of 6).\n", "Chapters:\n1. Считаем варианты:", "Matching (Codeforces 1821A, easy).", "Просмотрите Array.", "Добавьте элемент.", "Call Stack хранит незавершённые вызовы.", "Закончите ветку.", "Найдите следующий порядок."},
			wantExcludes: []string{"The rest is cut"}},
	}
	golden, replaced := goldenSets(t), tutorContextSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			exec(t, pool, tutorContextInsertProblem, tutorContextProblem)
			slug := "two-sum"
			if tt.slug != "" {
				slug = tt.slug
				exec(t, pool, "UPDATE problems SET slug=$1 WHERE id=$2", slug, tutorContextProblem)
			}
			if tt.task {
				exec(t, pool, "INSERT INTO homework (id,title,instructions,due_date,created_at,updated_at) VALUES ($1,'Week 41','Use a hash map.','2026-10-12','now','now')", tutorContextHomework)
				exec(t, pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,'in_progress','now','now')", newID(), tutorContextHomework, tutorContextProblem)
			}
			for _, attempt := range tt.work {
				exec(t, pool, tutorContextInsertAttempt, newID(), tutorContextProblem, attempt.state, tutorContextJSON(t, tutorContextWork("My notes", attempt.code)), attempt.savedAt)
			}
			content := golden
			if tt.lessons {
				content = replaced
			}
			leetcode := &fakeLeetCode{reply: tt.statement}
			u := newTutorUpstream(t, tutorAnswer())
			s := tutorServer(t, Config{DB: pool, Sets: content, LeetCode: leetcode}, u)
			id := tutorContextProblem
			if tt.kind == "set" {
				id = tutorSet
			}
			var quote any
			if tt.quote != "" {
				quote = tt.quote
			}
			w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(map[string]any{"kind": tt.kind, "id": id, "quote": quote}), loginAs(t, s, "student"))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			if got := tutorEvents(t, w); !reflect.DeepEqual(got, tutorContextAnswered) {
				t.Errorf("events = %v, want %v", got, tutorContextAnswered)
			}
			if tt.leetcode {
				if len(leetcode.calls) != 1 || leetcode.calls[0].titleSlug != slug || !strings.Contains(leetcode.calls[0].fields, "content") {
					t.Errorf("LeetCode calls = %+v, want one call for %s with the field content", leetcode.calls, slug)
				}
			} else if len(leetcode.calls) != 0 {
				t.Errorf("LeetCode calls = %+v, want none", leetcode.calls)
			}
			_, block := tutorContextBody(t, u)
			for _, text := range tt.wantIncludes {
				if !strings.Contains(block, text) {
					t.Errorf("context block = %q, want %q in it", block, text)
				}
			}
			for _, text := range tt.wantExcludes {
				if strings.Contains(block, text) {
					t.Errorf("context block = %q, want no %q in it", block, text)
				}
			}
			if tt.wantCutDraft {
				_, draft, found := strings.Cut(block, "Current draft:\n")
				draft, cut := strings.CutSuffix(draft, " [cut]\n</context>")
				if length := len(utf16.Encode([]rune(draft))); !found || !cut || length != 12000 {
					t.Errorf("context block ends with %q, want 12000 characters of the draft and the cut mark", block[max(0, len(block)-200):])
				}
			}
		})
	}
}

// The whole lessons of a published set, or of the chapter of a problem, reach the tutor without a quote.
// Each case is "<name> reaches the tutor without a quote" of the TypeScript test.
func TestTutorWholeLessons(t *testing.T) {
	tests := []struct {
		name string
		set  string
		// problem is the slug of the problem of a problem chat, or "" for a set chat.
		problem string
	}{
		{name: "complete published counting set", set: "permutations-and-combinations"},
		{name: "complete published pointer set", set: "two-pointers"},
		{name: "complete published queue set", set: "heaps-and-queues"},
		{name: "complete published N-Queens chapter in problem chat", set: "permutations-and-combinations", problem: "n-queens"},
	}
	content := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			exec(t, pool, tutorContextInsertProblem, tutorContextProblem)
			data := map[string]any{"kind": "set", "id": tt.set}
			if tt.problem != "" {
				exec(t, pool, "UPDATE problems SET slug=$1 WHERE id=$2", tt.problem, tutorContextProblem)
				data = map[string]any{"kind": "problem", "id": tutorContextProblem}
			}
			u := newTutorUpstream(t, tutorAnswer())
			s := tutorServer(t, Config{DB: pool, Sets: content, LeetCode: &fakeLeetCode{reply: tutorContextStatement(nil)}}, u)
			w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(data), loginAs(t, s, "student"))
			if got := tutorEvents(t, w); !reflect.DeepEqual(got, tutorContextAnswered) {
				t.Errorf("events = %v, want %v", got, tutorContextAnswered)
			}
			_, block := tutorContextBody(t, u)
			var chapters []sets.Chapter
			for _, set := range content {
				if set.Slug != tt.set {
					continue
				}
				for _, chapter := range set.Chapters {
					if tt.problem == "" || slices.ContainsFunc(chapter.Tasks, func(task sets.Task) bool { return strings.HasSuffix(task.URL, "/"+tt.problem+"/") }) {
						chapters = append(chapters, chapter)
					}
				}
			}
			if len(chapters) == 0 {
				t.Fatalf("no chapter of the set %s for the problem %q", tt.set, tt.problem)
			}
			for _, chapter := range chapters {
				if !strings.Contains(block, htmlText(chapter.Lesson)) {
					t.Errorf("the context block has no whole lesson of the chapter %q", chapter.Title)
				}
			}
			if strings.Contains(block, "The rest is cut.") {
				t.Errorf("the context block has the cut mark")
			}
		})
	}
}

func TestTutorHTMLText(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{name: "script removed", html: "<p>a</p><script>x()</script>", want: "a"},
		{name: "sup becomes grouped power", html: "10<sup>4</sup>", want: "10^(4)"},
		{name: "sum in an exponent stays grouped", html: "2<sup>n + 1</sup> − 1", want: "2^(n + 1) − 1"},
		{name: "difference in an exponent stays grouped", html: "n · 2<sup>n − 1</sup>", want: "n · 2^(n − 1)"},
		{name: "entities decoded", html: "a &lt; b &amp;&amp; c", want: "a < b && c"},
		{name: "heading marked", html: "<h2>Base case</h2>", want: "## Base case"},
		{name: "code indentation kept", html: "<pre>if x:\n    y</pre>", want: "if x:\n    y"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := htmlText(tt.html); got != tt.want {
				t.Errorf("htmlText(%q) = %q, want %q", tt.html, got, tt.want)
			}
		})
	}
}

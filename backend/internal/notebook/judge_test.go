package notebook

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The cases of test/judge.test.ts and the expectedOf cases of test/tests.test.ts.

// The values come from the C code of the reference: https://gist.github.com/tommyettinger/46a874533244883189143505d203312c
func TestJudgeMulberry32(t *testing.T) {
	tests := []struct {
		name   string
		seed   uint32
		values []float64
	}{
		{name: "seed 0", seed: 0, values: []float64{1144304738, 1416247, 958946056}},
		{name: "seed 1", seed: 1, values: []float64{2693262067, 11749833, 2265367787}},
		{name: "largest seed", seed: 4294967295, values: []float64{3850105811, 813802916, 3073704848}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			random := sets.Mulberry32(tt.seed)
			got := []float64{random(), random(), random()}
			want := []float64{tt.values[0] / (1 << 32), tt.values[1] / (1 << 32), tt.values[2] / (1 << 32)}
			if !slices.Equal(got, want) {
				t.Errorf("values = %v, want %v", got, want)
			}
		})
	}
}

// judgeDouble is the reference that doubles the number of the input.
var judgeDouble = sets.Reference{Stdio: func(input string) string {
	n, _ := strconv.Atoi(strings.TrimSpace(input))
	return strconv.Itoa(2*n) + "\n"
}}

func TestJudgeVerdict(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []string
		outputs []judgeOutput
		want    map[string]any
	}{
		{name: "all outputs pass", inputs: []string{"1", "2", "3"}, outputs: []judgeOutput{{Value: new("2")}, {Value: new("4")}, {Value: new("6")}}, want: map[string]any{"verdict": "accepted", "failed": nil}},
		{name: "first failed test gives the verdict", inputs: []string{"1", "2", "3"}, outputs: []judgeOutput{{Value: new("2")}, {Value: new("5")}, {Timeout: true}}, want: map[string]any{"verdict": "wrong", "failed": map[string]any{"index": 1, "input": "2", "expected": "4\n", "output": "5"}}},
		{name: "timeout", inputs: []string{"1", "2", "3"}, outputs: []judgeOutput{{Timeout: true}, {Value: new("5")}, {Value: new("6")}}, want: map[string]any{"verdict": "timeout", "failed": map[string]any{"index": 0, "input": "1", "expected": "2\n", "output": ""}}},
		{name: "runtime error", inputs: []string{"1", "2", "3"}, outputs: []judgeOutput{{Value: new("2")}, {Value: new("4")}, {Error: new("ZeroDivisionError: division by zero")}}, want: map[string]any{"verdict": "error", "failed": map[string]any{"index": 2, "input": "3", "expected": "6\n", "output": "ZeroDivisionError: division by zero"}}},
		{name: "texts over 2,000 characters", inputs: []string{"1" + strings.Repeat(" ", 2500)}, outputs: []judgeOutput{{Value: new(strings.Repeat("3", 2500))}}, want: map[string]any{"verdict": "wrong", "failed": map[string]any{"index": 0, "input": "1" + strings.Repeat(" ", 1999) + "…", "expected": "2\n", "output": strings.Repeat("3", 2000) + "…"}}},
		{name: "text of 2,000 characters", inputs: []string{"1"}, outputs: []judgeOutput{{Value: new(strings.Repeat("3", 2000))}}, want: map[string]any{"verdict": "wrong", "failed": map[string]any{"index": 0, "input": "1", "expected": "2\n", "output": strings.Repeat("3", 2000)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewsEqualJSON(t, judge(judgeDouble, sets.TestSpec{Mode: "stdio"}, tt.inputs, tt.outputs), tt.want)
		})
	}
}

func TestJudgeHiddenAnswerAboveOwnTestLimit(t *testing.T) {
	answer := strings.Repeat("x", 100_001)
	reference := sets.Reference{Stdio: func(string) string { return answer }}
	reviewsEqualJSON(t, judge(reference, sets.TestSpec{Mode: "stdio"}, []string{"1"}, []judgeOutput{{Value: &answer}}), map[string]any{"verdict": "accepted", "failed": nil})
}

// The fake references check their arguments the way the set references check the statement constraints.
var (
	judgePrefix = sets.Reference{Arity: 2, Call: func(args []any) any {
		nums, isList := args[0].([]any)
		k, isNumber := args[1].(float64)
		sets.Need(isList && isNumber && k == math.Trunc(k))
		return nums[:min(int(k), len(nums))]
	}}
	judgeSum = sets.Reference{Stdio: func(stdin string) string {
		fields := strings.Fields(stdin)
		sets.Need(len(fields) >= 2)
		a, errA := strconv.Atoi(fields[0])
		b, errB := strconv.Atoi(fields[1])
		sets.Need(errA == nil && errB == nil)
		return strconv.Itoa(a+b) + "\n"
	}}
	// judgeBroken fails with a type error for an argument that is not a string.
	judgeBroken = sets.Reference{Arity: 1, Call: func(args []any) any { return strings.ToUpper(args[0].(string)) }}
	judgeEcho   = sets.Reference{Stdio: func(stdin string) string { return stdin }}
)

func TestJudgeExpectedOf(t *testing.T) {
	tests := []struct {
		name      string
		reference sets.Reference
		input     string
		want      map[string]any
	}{
		{name: "answer of a call reference", reference: judgePrefix, input: "[[1, 'a', True, None], 3]", want: map[string]any{"value": "[1, 'a', True]"}},
		{name: "answer of a stdio reference", reference: judgeSum, input: "2 3\n", want: map[string]any{"value": "5\n"}},
		{name: "call input that breaks the constraints", reference: judgePrefix, input: "[[1, 2], 1.5]", want: map[string]any{"error": "constraints"}},
		{name: "stdio input that breaks the constraints", reference: judgeSum, input: "a b", want: map[string]any{"error": "constraints"}},
		{name: "input that does not parse", reference: judgePrefix, input: "[[1, 2], 1", want: map[string]any{"error": "input"}},
		{name: "too few arguments", reference: judgePrefix, input: "[[1, 2]]", want: map[string]any{"error": "input"}},
		{name: "arguments not in a list", reference: judgePrefix, input: "1", want: map[string]any{"error": "input"}},
		{name: "reference error on a wrong type", reference: judgeBroken, input: "[1]", want: map[string]any{"error": "input"}},
		{name: "stdio answer of 100,000 characters", reference: judgeEcho, input: strings.Repeat("x", 100_000), want: map[string]any{"value": strings.Repeat("x", 100_000)}},
		{name: "stdio answer over 100,000 characters", reference: judgeEcho, input: strings.Repeat("x", 100_001), want: map[string]any{"error": "size"}},
		{name: "call answer over 100,000 characters", reference: judgePrefix, input: "[[" + strings.Repeat("1, ", 40_000) + "1], 40001]", want: map[string]any{"error": "size"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewsEqualJSON(t, expectedOf(tt.reference, tt.input), tt.want)
		})
	}
}

const (
	judgeTwoSumII = "two-sum-ii-input-array-is-sorted"
	// judgeWithNotes sets the notes of the draft. A draft with notes has work.
	judgeWithNotes     = "UPDATE attempts SET document=jsonb_set(document::jsonb,'{notes}',to_jsonb($1::text))::text"
	judgeInsertMessage = "INSERT INTO tutor_messages (id,context_kind,context_ref,author,body,created_at) VALUES ($1,$2,$3,$4,'Hint','now')"
)

// judgeTestGenerator is the test generator of Two Sum II. The last input takes a number from random.
func judgeTestGenerator(random func() float64) []string {
	n := 1 + int(random()*1000)
	return []string{"[[2, 7, 11, 15], 9]", "[[2, 3, 4], 6]", "[[0, " + strconv.Itoa(n) + "], " + strconv.Itoa(n) + "]"}
}

// judgeTwoSum is the reference of Two Sum II from src/leetcode/sets/two-pointers/solutions.ts.
// Task #239 replaces it with the ported reference of the content.
var judgeTwoSum = sets.Reference{Arity: 2, Call: func(args []any) any {
	number := func(value any) float64 {
		n, ok := value.(float64)
		sets.Need(ok && n == math.Trunc(n) && n >= -1000 && n <= 1000)
		return n
	}
	list, ok := args[0].([]any)
	sets.Need(ok && len(list) >= 2 && len(list) <= 30_000)
	a := make([]float64, len(list))
	for i, value := range list {
		a[i] = number(value)
	}
	sets.Need(slices.IsSorted(a))
	target := number(args[1])
	pairs, seen := 0, map[float64]int{}
	for _, x := range a {
		pairs += seen[target-x]
		seen[x]++
	}
	// The input must have exactly one solution.
	sets.Need(pairs == 1)
	l, r := 0, len(a)-1
	for a[l]+a[r] != target {
		if a[l]+a[r] < target {
			l++
		} else {
			r--
		}
	}
	return []any{float64(l + 1), float64(r + 1)}
}}

// judgeAnswers are the right outputs for seed 1. With seed 1 the last input is [[0, 628], 628].
var judgeAnswers = []map[string]any{{"value": "[1, 2]"}, {"value": "[1, 3]"}, {"value": "[1, 2]"}}

// judgeContent is the mock of contentOf in beforeEach: the learning sets of the Worker, and Two Sum II with its reference and the test generator.
func judgeContent(t *testing.T) Config {
	t.Helper()
	return Config{Sets: goldenSets(t), References: map[string]sets.Reference{judgeTwoSumII: judgeTwoSum}, Generators: map[string]sets.Generator{judgeTwoSumII: judgeTestGenerator}}
}

// judgeFixture is the state that beforeEach of test/judge.test.ts makes.
type judgeFixture struct {
	pool                          *pgxpool.Pool
	s                             *Server
	student, problemID, attemptID string
}

// judgeSetup makes a server with the content of cfg, the problem Two Sum II with a draft, and the open check of seed 1.
func judgeSetup(t *testing.T, cfg Config) judgeFixture {
	t.Helper()
	cfg.DB, cfg.PublicOrigin, cfg.Accounts, cfg.LeetCode = database(t), publicOrigin, testAccounts, &fakeLeetCode{}
	s, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	f := judgeFixture{pool: cfg.DB, s: s, student: loginAs(t, s, "student"), problemID: newID()}
	exec(t, f.pool, testsInsertProblem, f.problemID, judgeTwoSumII)
	w := request(t, s, "/problems/"+f.problemID+"/attempts", http.MethodPost, "{}", f.student)
	if w.Code != http.StatusCreated {
		t.Fatalf("new draft: status %d, body %s", w.Code, w.Body.String())
	}
	f.attemptID = decode[struct {
		ID string `json:"id"`
	}](t, w).ID
	exec(t, f.pool, "INSERT INTO judge_runs VALUES ($1,1,'now')", f.problemID)
	return f
}

// judgeStartSeed starts a check of the problem as the student and returns its seed.
func judgeStartSeed(t *testing.T, f judgeFixture) uint32 {
	t.Helper()
	w := request(t, f.s, "/problems/"+f.problemID+"/judge/start", http.MethodPost, nil, f.student)
	if w.Code != http.StatusOK {
		t.Fatalf("start: status %d, body %s", w.Code, w.Body.String())
	}
	return decode[struct {
		Seed uint32 `json:"seed"`
	}](t, w).Seed
}

// judgePostOutputs sends the outputs for the seed of the problem, with the draft version 1.
func judgePostOutputs(t *testing.T, f judgeFixture, seed uint32, outputs any, cookie string) *httptest.ResponseRecorder {
	t.Helper()
	return request(t, f.s, "/problems/"+f.problemID+"/judge", http.MethodPost, map[string]any{"seed": seed, "attemptId": f.attemptID, "version": 1, "outputs": outputs}, cookie)
}

// judgeWantError fails the test unless the response has the status and the error code.
func judgeWantError(t *testing.T, w *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if w.Code != status || errorCode(w) != code {
		t.Errorf("status %d, body %s; want %d with the error code %s", w.Code, w.Body.String(), status, code)
	}
}

// judgeFirstRow returns the first row of a query as a map from the column name to the value, or nil for no rows, as first() of D1.
func judgeFirstRow(t *testing.T, pool *pgxpool.Pool, sql string, args ...any) map[string]any {
	t.Helper()
	rows, err := pool.Query(context.Background(), sql, args...)
	if err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	values, err := pgx.CollectRows(rows, pgx.RowToMap)
	if err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	if len(values) == 0 {
		return nil
	}
	return values[0]
}

func TestJudgeStart(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		judged     bool
		wantStatus int
		wantCode   string
	}{
		{name: "student gets the inputs of a new seed", role: "student", judged: true, wantStatus: 200},
		{name: "parent cannot start the judge", role: "parent", judged: true, wantStatus: 403, wantCode: "forbidden"},
		{name: "problem without a judge", role: "student", judged: false, wantStatus: 404, wantCode: "no_judge"},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			id := f.problemID
			if !tt.judged {
				id = newID()
				exec(t, f.pool, testsInsertProblem, id, "two-sum")
			}
			w := request(t, f.s, "/problems/"+id+"/judge/start", http.MethodPost, nil, loginAs(t, f.s, tt.role))
			if tt.wantCode != "" {
				judgeWantError(t, w, tt.wantStatus, tt.wantCode)
				return
			}
			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			got := decode[struct {
				Seed   uint32   `json:"seed"`
				Inputs []string `json:"inputs"`
			}](t, w)
			if want := judgeTestGenerator(sets.Mulberry32(got.Seed)); !slices.Equal(got.Inputs, want) {
				t.Errorf("inputs = %q, want %q", got.Inputs, want)
			}
		})
	}
}

func TestJudgeOutputs(t *testing.T) {
	invalid := map[string]any{"error": map[string]any{"code": "invalid_input"}}
	tests := []struct {
		name       string
		role       string
		outputs    []map[string]any
		wantStatus int
		wantBody   map[string]any
	}{
		{name: "wrong answer", role: "student", outputs: []map[string]any{judgeAnswers[0], {"value": "[2, 3]"}, judgeAnswers[2]}, wantStatus: 200, wantBody: map[string]any{"verdict": "wrong", "failed": map[string]any{"index": 1, "input": "[[2, 3, 4], 6]", "expected": "[1, 3]", "output": "[2, 3]"}}},
		{name: "wrong answer of the generated input", role: "student", outputs: []map[string]any{judgeAnswers[0], judgeAnswers[1], {"value": "[1, 628]"}}, wantStatus: 200, wantBody: map[string]any{"verdict": "wrong", "failed": map[string]any{"index": 2, "input": "[[0, 628], 628]", "expected": "[1, 2]", "output": "[1, 628]"}}},
		{name: "output at 1,000,000 characters", role: "student", outputs: []map[string]any{judgeAnswers[0], judgeAnswers[1], {"value": strings.Repeat("9", 1_000_000)}}, wantStatus: 200, wantBody: map[string]any{"verdict": "wrong", "failed": map[string]any{"index": 2, "input": "[[0, 628], 628]", "expected": "[1, 2]", "output": strings.Repeat("9", 2000) + "…"}}},
		{name: "output over 1,000,000 characters", role: "student", outputs: []map[string]any{judgeAnswers[0], judgeAnswers[1], {"value": strings.Repeat("9", 1_000_001)}}, wantStatus: 400, wantBody: invalid},
		{name: "fewer outputs than inputs", role: "student", outputs: judgeAnswers[1:], wantStatus: 400, wantBody: invalid},
		{name: "output with a value and an error", role: "student", outputs: []map[string]any{{"value": "[1, 2]", "error": "ValueError"}, judgeAnswers[1], judgeAnswers[2]}, wantStatus: 400, wantBody: invalid},
		{name: "timeout that is false", role: "student", outputs: []map[string]any{{"timeout": false}, judgeAnswers[1], judgeAnswers[2]}, wantStatus: 400, wantBody: invalid},
		{name: "parent cannot judge", role: "parent", outputs: judgeAnswers, wantStatus: 403, wantBody: map[string]any{"error": map[string]any{"code": "forbidden"}}},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			w := judgePostOutputs(t, f, 1, tt.outputs, loginAs(t, f.s, tt.role))
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), tt.wantBody)
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT state,acceptance FROM attempts"), map[string]any{"state": "draft", "acceptance": "not_submitted"})
		})
	}
}

func TestJudgeSeed(t *testing.T) {
	tests := []struct {
		name       string
		start      bool
		replay     bool
		answered   string
		wantStatus int
		wantCode   string
	}{
		{name: "seed of a started check", start: true, replay: false, answered: "", wantStatus: 200},
		{name: "seed that the server did not give", start: false, replay: false, answered: "", wantStatus: 409, wantCode: "stale_seed"},
		{name: "seed used for a second verdict", start: true, replay: true, answered: "", wantStatus: 409, wantCode: "stale_seed"},
		{name: "seed whose inputs the reference answered as own tests", start: true, replay: false, answered: "same", wantStatus: 409, wantCode: "stale_seed"},
		{name: "seed invalidated by another problem without a reference", start: true, replay: false, answered: "other", wantStatus: 409, wantCode: "stale_seed"},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			seed := uint32(2)
			if tt.start {
				seed = judgeStartSeed(t, f)
			}
			inputs := judgeTestGenerator(sets.Mulberry32(seed))
			outputs := slices.Repeat([]map[string]any{{"value": "[9, 9]"}}, len(inputs))
			if tt.replay {
				judgePostOutputs(t, f, seed, outputs, f.student)
			}
			if tt.answered != "" {
				id := f.problemID
				if tt.answered == "other" {
					id = newID()
					exec(t, f.pool, testsInsertProblem, id, "two-sum")
				}
				if w := request(t, f.s, "/problems/"+id+"/tests/expected", http.MethodPost, map[string]any{"inputs": inputs}, f.student); w.Code != http.StatusOK {
					t.Fatalf("expected answers: status %d, body %s", w.Code, w.Body.String())
				}
			}
			w := judgePostOutputs(t, f, seed, outputs, f.student)
			if tt.wantCode != "" {
				judgeWantError(t, w, tt.wantStatus, tt.wantCode)
				return
			}
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"verdict": "wrong", "failed": map[string]any{"index": 0}})
		})
	}
}

// judgeStdioSets has the mocked tasks of the disclosure test: Two Sum II and Codeforces 1821A with stdio tests and no cases.
var judgeStdioSets = []sets.LearningSet{{Slug: "judge-set", Title: "Judge set", Chapters: []sets.Chapter{{Title: "First chapter", Tasks: []sets.Task{
	{URL: "https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/", Title: "Two Sum II - Input Array Is Sorted", Difficulty: "medium", Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{}}},
	{URL: "https://codeforces.com/problemset/problem/1821/A", Title: "Matching", Difficulty: "easy", Tests: &sets.TaskTests{TestSpec: sets.TestSpec{Mode: "stdio"}, Cases: []sets.TestCase{}}},
}}}}}

// Two Sum II has the reference judgeDouble and the hidden inputs 1 and 2. Codeforces 1821A has the same reference and no generator.
// The three checks start with the seeds 11, 12 and 13.
func TestJudgeDisclosure(t *testing.T) {
	tests := []struct {
		name         string
		target       string
		role         string
		wantExpected []any
	}{
		{name: "owns disclosure through the same problem and the next acceptance check", target: "same", role: "student", wantExpected: []any{map[string]any{"value": "4\n"}}},
		{name: "owns disclosure through another problem with a reference but no generator and the next acceptance check", target: "reference", role: "student", wantExpected: []any{map[string]any{"value": "4\n"}}},
		{name: "owns disclosure through another problem without a reference and the next acceptance check", target: "none", role: "student", wantExpected: []any{nil}},
		{name: "owns disclosure through a parent request and the next acceptance check", target: "same", role: "parent", wantExpected: []any{map[string]any{"value": "4\n"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, Config{
				Sets:       judgeStdioSets,
				References: map[string]sets.Reference{judgeTwoSumII: judgeDouble, "cf-1821-a": judgeDouble},
				Generators: map[string]sets.Generator{judgeTwoSumII: func(func() float64) []string { return []string{"1", "2"} }},
			})
			seeds := []int{11, 12, 13}
			f.s.random = func(count int) int {
				seed := seeds[0]
				seeds = seeds[1:]
				return seed % count
			}
			first, second := judgeStartSeed(t, f), judgeStartSeed(t, f)
			id := f.problemID
			if tt.target != "same" {
				id = newID()
				slug := "two-sum"
				if tt.target == "reference" {
					slug = "cf-1821-a"
				}
				exec(t, f.pool, testsInsertProblem, id, slug)
			}
			disclosure := request(t, f.s, "/problems/"+id+"/tests/expected", http.MethodPost, map[string]any{"inputs": []string{"2"}}, loginAs(t, f.s, tt.role))
			reviewsEqualJSON(t, decode[any](t, disclosure), map[string]any{"expected": tt.wantExpected})
			outputs := []map[string]any{{"value": "2\n"}, {"value": "4\n"}}
			judgeWantError(t, judgePostOutputs(t, f, first, outputs, f.student), http.StatusConflict, "stale_seed")
			judgeWantError(t, judgePostOutputs(t, f, second, outputs, f.student), http.StatusConflict, "stale_seed")
			fresh := judgeStartSeed(t, f)
			reviewsEqualJSON(t, decode[any](t, judgePostOutputs(t, f, fresh, outputs, f.student)), map[string]any{"verdict": "accepted", "failed": nil})
			judgeWantError(t, judgePostOutputs(t, f, fresh, outputs, f.student), http.StatusConflict, "stale_seed")
		})
	}
}

func TestJudgeDisclosureEnvelope(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
		stale      bool
	}{
		{name: "invalid envelope keeps open checks", body: map[string]any{"inputs": "1"}, wantStatus: 400, stale: false},
		{name: "too many inputs keep open checks", body: map[string]any{"inputs": slices.Repeat([]string{"1"}, 21)}, wantStatus: 400, stale: false},
		{name: "invalid entry invalidates open checks", body: map[string]any{"inputs": []any{1}}, wantStatus: 400, stale: true},
		{name: "empty disclosure invalidates open checks", body: map[string]any{"inputs": []string{}}, wantStatus: 200, stale: true},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			if w := request(t, f.s, "/problems/"+f.problemID+"/tests/expected", http.MethodPost, tt.body, f.student); w.Code != tt.wantStatus {
				t.Errorf("expected answers: status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			w := judgePostOutputs(t, f, 1, judgeAnswers, f.student)
			if tt.stale {
				judgeWantError(t, w, http.StatusConflict, "stale_seed")
				return
			}
			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"verdict": "accepted", "failed": nil})
		})
	}
}

func TestJudgeRejectsOutputsOfAChangedDraft(t *testing.T) {
	f := judgeSetup(t, judgeContent(t))
	problemsAddTask(t, f.pool, f.problemID, "in_progress")
	path := "/problems/" + f.problemID + "/attempts/" + f.attemptID
	draft := decode[struct {
		Version  int            `json:"version"`
		Document map[string]any `json:"document"`
	}](t, request(t, f.s, path, http.MethodGet, nil, f.student))
	approaches, _ := draft.Document["approaches"].([]any)
	if len(approaches) == 0 {
		t.Fatalf("document = %v, want an approach", draft.Document)
	}
	// Another tab changes the code of the tested draft.
	approaches[0].(map[string]any)["code"] = `raise RuntimeError("new code from another tab")`
	if w := request(t, f.s, path, http.MethodPut, map[string]any{"version": draft.Version, "document": draft.Document}, f.student); w.Code != http.StatusOK {
		t.Fatalf("edit: status %d, body %s", w.Code, w.Body.String())
	}
	w := request(t, f.s, "/problems/"+f.problemID+"/judge", http.MethodPost, map[string]any{"seed": 1, "attemptId": f.attemptID, "version": draft.Version, "outputs": judgeAnswers}, f.student)
	judgeWantError(t, w, http.StatusConflict, "conflict")
	problemsMatchBody(t, decode[any](t, request(t, f.s, path, http.MethodGet, nil, f.student)), map[string]any{"state": "draft", "acceptance": "not_submitted", "version": draft.Version + 1, "document": draft.Document})
	reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT understanding,progress_version FROM problems"), map[string]any{"understanding": "needs_practice", "progress_version": 1})
	reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT state,(SELECT count(*) FROM submissions) AS submissions,(SELECT count(*) FROM feedback) AS feedback FROM homework_tasks"), map[string]any{"state": "in_progress", "submissions": 0, "feedback": 0})
}

func TestJudgeDraftVersion(t *testing.T) {
	tests := []struct {
		name    string
		version any
	}{
		{name: "rejects an invalid tested draft version: undefined", version: nil},
		{name: "rejects an invalid tested draft version: 0", version: 0},
		{name: "rejects an invalid tested draft version: 1.5", version: 1.5},
		{name: "rejects an invalid tested draft version: 1", version: "1"},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			body := map[string]any{"seed": 1, "attemptId": f.attemptID, "outputs": judgeAnswers}
			// A nil version leaves the field out, as undefined does in JSON.stringify.
			if tt.version != nil {
				body["version"] = tt.version
			}
			judgeWantError(t, request(t, f.s, "/problems/"+f.problemID+"/judge", http.MethodPost, body, f.student), http.StatusBadRequest, "invalid_input")
		})
	}
}

func TestJudgeAccepted(t *testing.T) {
	completed := map[string]any{"state": "completed", "version": 2, "taskVersion": 2, "author": "student", "kind": "completed", "body": "Принято проверкой."}
	tests := []struct {
		name          string
		messageKind   string
		messageAuthor string
		task          string
		understanding string
		after         map[string]any
	}{
		{name: "accepted without tutor messages", messageKind: "", messageAuthor: "", task: "", understanding: "independent", after: nil},
		{name: "accepted after a student message to the tutor", messageKind: "problem", messageAuthor: "student", task: "", understanding: "with_help", after: nil},
		{name: "accepted after only tutor answers", messageKind: "problem", messageAuthor: "assistant", task: "", understanding: "independent", after: nil},
		{name: "accepted after a message in the set thread", messageKind: "set", messageAuthor: "student", task: "", understanding: "independent", after: nil},
		{name: "accepted task in progress", messageKind: "", messageAuthor: "", task: "in_progress", understanding: "independent", after: completed},
		{name: "accepted assigned task", messageKind: "", messageAuthor: "", task: "assigned", understanding: "independent", after: completed},
		{name: "accepted submitted task", messageKind: "", messageAuthor: "", task: "submitted", understanding: "independent", after: completed},
		{name: "accepted task with requested changes", messageKind: "", messageAuthor: "", task: "changes_requested", understanding: "independent", after: completed},
		{name: "completed task stays", messageKind: "", messageAuthor: "", task: "completed", understanding: "independent", after: map[string]any{"state": "completed", "version": 1, "taskVersion": nil, "author": nil, "kind": nil, "body": nil}},
		{name: "cancelled task stays", messageKind: "", messageAuthor: "", task: "cancelled", understanding: "independent", after: map[string]any{"state": "cancelled", "version": 1, "taskVersion": nil, "author": nil, "kind": nil, "body": nil}},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			if tt.messageKind != "" {
				thread := f.problemID
				if tt.messageKind == "set" {
					thread = "two-pointers"
				}
				exec(t, f.pool, judgeInsertMessage, newID(), tt.messageKind, thread, tt.messageAuthor)
			}
			if tt.task != "" {
				problemsAddTask(t, f.pool, f.problemID, tt.task)
			}
			w := judgePostOutputs(t, f, 1, judgeAnswers, f.student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			reviewsEqualJSON(t, decode[any](t, w), map[string]any{"verdict": "accepted", "failed": nil})
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT state,acceptance,understanding,version FROM attempts"), map[string]any{"state": "saved", "acceptance": "accepted", "understanding": tt.understanding, "version": 2})
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT understanding,progress_version FROM problems"), map[string]any{"understanding": tt.understanding, "progress_version": 2})
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, `SELECT t.state,t.version,s.task_version AS "taskVersion",f.author,f.kind,f.body FROM homework_tasks t LEFT JOIN submissions s ON s.id=t.current_submission_id AND s.attempt_id=$1 LEFT JOIN feedback f ON f.submission_id=s.id AND f.task_id=t.id`, f.attemptID), tt.after)
		})
	}
}

func TestJudgeAcceptedOutputsRejected(t *testing.T) {
	tests := []struct {
		name       string
		saved      bool
		known      bool
		wantStatus int
		wantCode   string
	}{
		{name: "accepted outputs of a saved attempt", saved: true, known: true, wantStatus: 409, wantCode: "conflict"},
		{name: "accepted outputs of an unknown attempt", saved: false, known: false, wantStatus: 404, wantCode: "not_found"},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			problemsAddTask(t, f.pool, f.problemID, "in_progress")
			if tt.saved {
				exec(t, f.pool, "UPDATE attempts SET state='saved'")
			}
			attemptID := f.attemptID
			if !tt.known {
				attemptID = newID()
			}
			w := request(t, f.s, "/problems/"+f.problemID+"/judge", http.MethodPost, map[string]any{"seed": 1, "attemptId": attemptID, "version": 1, "outputs": judgeAnswers}, f.student)
			judgeWantError(t, w, tt.wantStatus, tt.wantCode)
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT acceptance,understanding FROM attempts"), map[string]any{"acceptance": "not_submitted", "understanding": "needs_practice"})
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT t.state,(SELECT count(*) FROM submissions) AS submissions,(SELECT count(*) FROM feedback) AS feedback FROM homework_tasks t"), map[string]any{"state": "in_progress", "submissions": 0, "feedback": 0})
		})
	}
}

// The problem has no generator, as the mock of contentOf in these cases.
func TestJudgeSolved(t *testing.T) {
	tests := []struct {
		name          string
		message       bool
		task          string
		understanding string
		after         map[string]any
	}{
		{name: "solved without tutor messages", message: false, task: "", understanding: "independent", after: nil},
		{name: "solved after a student message to the tutor", message: true, task: "", understanding: "with_help", after: nil},
		{name: "solved assigned task", message: false, task: "assigned", understanding: "independent", after: map[string]any{"state": "submitted", "version": 2, "taskVersion": 2, "feedback": 0}},
		{name: "solved task in progress", message: false, task: "in_progress", understanding: "independent", after: map[string]any{"state": "submitted", "version": 2, "taskVersion": 2, "feedback": 0}},
		{name: "solved task with requested changes", message: false, task: "changes_requested", understanding: "independent", after: map[string]any{"state": "submitted", "version": 2, "taskVersion": 2, "feedback": 0}},
		{name: "submitted task stays", message: false, task: "submitted", understanding: "independent", after: map[string]any{"state": "submitted", "version": 1, "taskVersion": nil, "feedback": 0}},
		{name: "completed task stays", message: false, task: "completed", understanding: "independent", after: map[string]any{"state": "completed", "version": 1, "taskVersion": nil, "feedback": 0}},
		{name: "cancelled task stays", message: false, task: "cancelled", understanding: "independent", after: map[string]any{"state": "cancelled", "version": 1, "taskVersion": nil, "feedback": 0}},
	}
	content := judgeContent(t)
	content.Generators = nil
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := judgeSetup(t, content)
			exec(t, f.pool, judgeWithNotes, "Two pointers.")
			if tt.message {
				exec(t, f.pool, judgeInsertMessage, newID(), "problem", f.problemID, "student")
			}
			if tt.task != "" {
				problemsAddTask(t, f.pool, f.problemID, tt.task)
			}
			w := request(t, f.s, "/problems/"+f.problemID+"/solved", http.MethodPost, map[string]any{"attemptId": f.attemptID}, f.student)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, w), map[string]any{"id": f.attemptID, "state": "saved", "acceptance": "accepted", "understanding": tt.understanding, "version": 2})
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT understanding,progress_version FROM problems"), map[string]any{"understanding": tt.understanding, "progress_version": 2})
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, `SELECT t.state,t.version,s.task_version AS "taskVersion",(SELECT count(*) FROM feedback) AS feedback FROM homework_tasks t LEFT JOIN submissions s ON s.id=t.current_submission_id AND s.attempt_id=$1`, f.attemptID), tt.after)
		})
	}
}

func TestJudgeSolvedRejected(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		judged     bool
		notes      string
		attempt    string
		wantStatus int
		wantCode   string
	}{
		{name: "parent cannot mark a problem solved", role: "parent", judged: false, notes: "Two pointers.", attempt: "draft", wantStatus: 403, wantCode: "forbidden"},
		{name: "problem with a judge", role: "student", judged: true, notes: "Two pointers.", attempt: "draft", wantStatus: 409, wantCode: "judged"},
		{name: "empty work", role: "student", judged: false, notes: " ", attempt: "draft", wantStatus: 400, wantCode: "empty_attempt"},
		{name: "solved attempt again", role: "student", judged: false, notes: "Two pointers.", attempt: "saved", wantStatus: 409, wantCode: "conflict"},
		{name: "unknown attempt", role: "student", judged: false, notes: "Two pointers.", attempt: "unknown", wantStatus: 404, wantCode: "not_found"},
	}
	content := judgeContent(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := content
			if !tt.judged {
				cfg.Generators = nil
			}
			f := judgeSetup(t, cfg)
			exec(t, f.pool, judgeWithNotes, tt.notes)
			problemsAddTask(t, f.pool, f.problemID, "in_progress")
			if tt.attempt == "saved" {
				exec(t, f.pool, "UPDATE attempts SET state='saved'")
			}
			attemptID := f.attemptID
			if tt.attempt == "unknown" {
				attemptID = newID()
			}
			w := request(t, f.s, "/problems/"+f.problemID+"/solved", http.MethodPost, map[string]any{"attemptId": attemptID}, loginAs(t, f.s, tt.role))
			judgeWantError(t, w, tt.wantStatus, tt.wantCode)
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT acceptance,understanding FROM attempts"), map[string]any{"acceptance": "not_submitted", "understanding": "needs_practice"})
			reviewsEqualJSON(t, judgeFirstRow(t, f.pool, "SELECT t.state,(SELECT count(*) FROM submissions) AS submissions FROM homework_tasks t"), map[string]any{"state": "in_progress", "submissions": 0})
		})
	}
}

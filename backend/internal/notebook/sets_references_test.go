package notebook

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf16"

	"github.com/rustatian/redlamp.me/backend/internal/answers"
	"github.com/rustatian/redlamp.me/backend/internal/sets"
	"github.com/rustatian/redlamp.me/backend/internal/sets/content"
)

// The cases of test/sets.test.ts that need the references and the generators of the Go content.

// setReferencesServer returns a server with the Go content, as contentOf of the Worker. The tests send no requests, so the server has no database.
func setReferencesServer(t *testing.T) *Server {
	t.Helper()
	s, err := New(Config{PublicOrigin: publicOrigin, Accounts: testAccounts, LeetCode: &fakeLeetCode{}, Sets: content.Sets, References: content.References, Generators: content.Generators})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// setReferencesTask is a task of content.Sets. name is the task slug.
type setReferencesTask struct {
	name string
	task sets.Task
}

// setReferencesTasks returns the tasks of all learning sets of the Go content.
func setReferencesTasks() []setReferencesTask {
	var tests []setReferencesTask
	for _, set := range content.Sets {
		for _, task := range set.Tasks() {
			tests = append(tests, setReferencesTask{name: task.Slug(), task: task})
		}
	}
	return tests
}

// setReferencesMode returns the mode of a reference: "stdio", "call", or "" for a missing reference.
func setReferencesMode(reference sets.Reference) string {
	switch {
	case reference.Stdio != nil:
		return "stdio"
	case reference.Call != nil:
		return "call"
	}
	return ""
}

func TestSetReferencesResolve(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		spec    sets.TestSpec
		input   string
		answer  string
		chapter int64
		title   string
		lesson  string
	}{
		{name: "in-place problem with its ordering lesson", slug: "next-permutation", spec: sets.TestSpec{Mode: "call", Method: "nextPermutation", Check: "arg0"}, input: "[[1, 2, 3]]", answer: "[1, 3, 2]", chapter: 4, title: "Порядок и номер", lesson: "2. Next Permutation"},
		{name: "standard-input problem with its counting lesson", slug: "cf-1821-a", spec: sets.TestSpec{Mode: "stdio"}, input: "2\n??\n0\n", answer: "90\n0\n", chapter: 1, title: "Считаем варианты", lesson: "2. Product Rule"},
	}
	s := setReferencesServer(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, ok := s.setTask(tt.slug)
			if !ok || task.Tests == nil {
				t.Fatalf("task %s: found %v, tests %v; want a task with tests", tt.slug, ok, task.Tests)
			}
			if task.Tests.TestSpec != tt.spec {
				t.Errorf("test spec = %+v, want %+v", task.Tests.TestSpec, tt.spec)
			}
			reviewsEqualJSON(t, expectedOf(s.references[tt.slug], tt.input), map[string]any{"value": tt.answer})
			want := []chapterRef{{Set: "permutations-and-combinations", SetTitle: "Permutations и Combinations", Number: tt.chapter, Title: tt.title}}
			chapters := s.chaptersOf(tt.slug)
			if !slices.Equal(chapters, want) {
				t.Fatalf("chapters = %+v, want %+v", chapters, want)
			}
			if chapter, _ := s.chapter(chapters[0].Set, chapters[0].Number); !strings.Contains(chapter.Lesson, tt.lesson) {
				t.Errorf("lesson of chapter %d has no %q", chapters[0].Number, tt.lesson)
			}
		})
	}
}

func TestSetReferencesPassStandardTests(t *testing.T) {
	for _, tt := range setReferencesTasks() {
		t.Run(tt.name, func(t *testing.T) {
			if tt.task.Tests == nil {
				t.Fatal("the task has no tests")
			}
			tests, reference := tt.task.Tests, content.References[tt.name]
			if mode := setReferencesMode(reference); mode != tests.Mode {
				t.Fatalf("reference mode = %q, want %q", mode, tests.Mode)
			}
			spec := answers.Spec{Mode: tests.Mode, Order: tests.Order, IgnoreCase: tests.IgnoreCase}
			var failed []string
			for _, c := range tests.Cases {
				answer := expectedOf(reference, c.Input)
				if c.Expected == nil || answer.Value == nil || !answers.Same(spec, *c.Expected, *answer.Value) {
					failed = append(failed, c.Input)
				}
			}
			checkNone(t, "standard tests that the reference fails", failed)
		})
	}
}

// Each output letter can use either case: https://codeforces.com/problemset/problem/1097/B
func TestSetReferencesCodeforces1097BJudges(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		output  string
		verdict string
	}{
		{name: "lowercase YES", input: "3\n10\n20\n30\n", output: "yes\n", verdict: "accepted"},
		{name: "mixed-case YES", input: "3\n10\n20\n30\n", output: "yEs\n", verdict: "accepted"},
		{name: "mixed-case NO", input: "3\n10\n10\n10\n", output: "No\n", verdict: "accepted"},
		{name: "wrong answer", input: "3\n10\n20\n30\n", output: "no\n", verdict: "wrong"},
	}
	s := setReferencesServer(t)
	task, ok := s.setTask("cf-1097-b")
	if !ok || task.Tests == nil {
		t.Fatal("no task cf-1097-b with tests")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := judge(s.references["cf-1097-b"], task.Tests.TestSpec, []string{tt.input}, []judgeOutput{{Value: &tt.output}})
			if got.Verdict != tt.verdict {
				t.Errorf("verdict = %q, want %q", got.Verdict, tt.verdict)
			}
		})
	}
}

// setReferencesStack is one state of the bonus stack and the total power of the army.
type setReferencesStack struct {
	stack []int
	total int
}

// A bonus can enter the stack or be discarded. A hero can use only the top bonus: https://codeforces.com/problemset/problem/1800/C2
func TestSetReferencesCodeforces1800C2StackChoices(t *testing.T) {
	tests := []struct {
		name  string
		cards []int
		total int
	}{
		{name: "smaller bonus above a larger bonus", cards: []int{4, 1, 0, 0}, total: 5},
		{name: "bonus after the only hero", cards: []int{0, 9}, total: 0},
		{name: "new bonus between heroes", cards: []int{4, 1, 0, 7, 0, 0}, total: 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			states := []setReferencesStack{{}}
			for _, card := range tt.cards {
				var next []setReferencesStack
				for _, state := range states {
					if card > 0 {
						next = append(next, state, setReferencesStack{stack: append(slices.Clone(state.stack), card), total: state.total})
						continue
					}
					top, rest := 0, state.stack
					if n := len(state.stack); n > 0 {
						top, rest = state.stack[n-1], state.stack[:n-1]
					}
					next = append(next, setReferencesStack{stack: rest, total: state.total + top})
				}
				states = next
			}
			best := math.MinInt
			for _, state := range states {
				best = max(best, state.total)
			}
			if best != tt.total {
				t.Errorf("best total of the legal choices = %d, want %d", best, tt.total)
			}
			input := fmt.Sprintf("1\n%d\n%s\n", len(tt.cards), strings.Trim(fmt.Sprint(tt.cards), "[]"))
			reviewsEqualJSON(t, expectedOf(content.References["cf-1800-c2"], input), map[string]any{"value": strconv.Itoa(tt.total) + "\n"})
		})
	}
}

// setReferencesArgs returns the arguments of a call input, or nil for an input that is not a list.
func setReferencesArgs(input string) []any {
	value, _ := answers.Parse(input)
	args, _ := value.([]any)
	return args
}

// setReferencesHiddenTest returns the problem of a hidden test, or "" for a valid test.
// The input must have the argument count of the reference. The fastest of three runs must take less than 100 ms,
// and the answer must have 1,000,000 characters or fewer, as the judge route needs.
func setReferencesHiddenTest(reference sets.Reference, input string) (problem string) {
	defer func() {
		if r := recover(); r != nil {
			problem = fmt.Sprint("the reference panics: ", r)
		}
	}()
	args := setReferencesArgs(input)
	if reference.Stdio == nil && len(args) != reference.Arity {
		return fmt.Sprintf("%d arguments, want %d", len(args), reference.Arity)
	}
	// The fastest of three runs leaves out the pauses that other tests cause.
	answer, fastest := "", time.Duration(math.MaxInt64)
	for run := range 3 {
		// Each run gets new arguments, because some references change their input in place.
		if run > 0 && reference.Stdio == nil {
			args = setReferencesArgs(input)
		}
		start := time.Now()
		if reference.Stdio != nil {
			answer = reference.Stdio(input)
		} else {
			var err error
			if answer, err = answers.Repr(reference.Call(args)); err != nil {
				return err.Error()
			}
		}
		fastest = min(fastest, time.Since(start))
	}
	if !raceDetector && fastest >= 100*time.Millisecond {
		return fmt.Sprintf("the fastest run takes %v", fastest)
	}
	if n := len(utf16.Encode([]rune(answer))); n > 1_000_000 {
		return fmt.Sprintf("the answer has %d characters", n)
	}
	return ""
}

// The hidden tests of seeds 1 and 2^32 - 1.
func TestSetReferencesGeneratorsGiveValidHiddenTests(t *testing.T) {
	for _, tt := range setReferencesTasks() {
		t.Run(tt.name, func(t *testing.T) {
			var first, last []string
			if generate := content.Generators[tt.name]; generate != nil {
				first, last = generate(sets.Mulberry32(1)), generate(sets.Mulberry32(math.MaxUint32))
			}
			if len(first) < 10 || len(first) > 16 {
				t.Errorf("seed 1 gives %d inputs, want 10 to 16", len(first))
			}
			var bad []string
			for i, input := range slices.Concat(first, last) {
				if problem := setReferencesHiddenTest(content.References[tt.name], input); problem != "" {
					bad = append(bad, fmt.Sprintf("input %d: %s", i, problem))
				}
			}
			checkNone(t, "invalid hidden tests", bad)
		})
	}
}

// Wrong solutions of LeetCode 2462 that pass the examples of the statement.
func setReferencesCheapest(costs []int, k, _ int) int {
	total := 0
	for _, cost := range slices.Sorted(slices.Values(costs))[:k] {
		total += cost
	}
	return total
}

// setReferencesRightTie gives a tie of the two groups to the right group.
func setReferencesRightTie(costs []int, k, candidates int) int {
	var left, right []int
	least := func(group []int) int {
		best := group[0]
		for _, i := range group {
			if costs[i] < costs[best] || (costs[i] == costs[best] && i < best) {
				best = i
			}
		}
		return best
	}
	l, r, total := 0, len(costs)-1, 0
	for ; l < candidates && l <= r; l++ {
		left = append(left, l)
	}
	for ; len(costs)-1-r < candidates && l <= r; r-- {
		right = append(right, r)
	}
	for range k {
		fromLeft := len(right) == 0 || (len(left) > 0 && costs[least(left)] < costs[least(right)])
		group := &right
		if fromLeft {
			group = &left
		}
		i := least(*group)
		at := slices.Index(*group, i)
		*group = slices.Delete(*group, at, at+1)
		total += costs[i]
		if l <= r {
			if fromLeft {
				*group = append(*group, l)
				l++
			} else {
				*group = append(*group, r)
				r--
			}
		}
	}
	return total
}

// The hidden tests of every seed must reject the wrong solutions. The check uses the inputs of 1,000 characters or fewer.
func TestSetReferencesLeetCode2462RejectsWrongSolutions(t *testing.T) {
	const slug = "total-cost-to-hire-k-workers"
	tests := []struct {
		name  string
		wrong func(costs []int, k, candidates int) int
	}{
		{name: "sum of the k cheapest costs", wrong: setReferencesCheapest},
		{name: "tie that goes to the right group", wrong: setReferencesRightTie},
	}
	reference, generate := content.References[slug], content.Generators[slug]
	if reference.Call == nil || generate == nil {
		t.Fatalf("no reference or generator for %s", slug)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for seed := uint32(1); seed <= 50; seed++ {
				rejected := false
				for _, input := range generate(sets.Mulberry32(seed)) {
					if len(input) > 1000 {
						continue
					}
					args := setReferencesArgs(input)
					if len(args) != 3 {
						t.Fatalf("seed %d: the input %q does not have 3 arguments", seed, input)
					}
					list, _ := args[0].([]any)
					costs := make([]int, len(list))
					for i, cost := range list {
						costs[i] = int(cost.(float64))
					}
					if sets.Py(tt.wrong(costs, int(args[1].(float64)), int(args[2].(float64)))) != reference.Call(args) {
						rejected = true
						break
					}
				}
				if !rejected {
					t.Errorf("seed %d: the hidden tests accept the wrong solution", seed)
				}
			}
		})
	}
}

// Each task of the Go content has a reference and a generator, and each reference and generator has a task.
func TestSetReferencesCoverTheTasks(t *testing.T) {
	var slugs []string
	for _, task := range setReferencesTasks() {
		slugs = append(slugs, task.name)
	}
	slices.Sort(slugs)
	slugs = slices.Compact(slugs)
	tests := []struct {
		name string
		keys []string
	}{
		{name: "references", keys: slices.Sorted(maps.Keys(content.References))},
		{name: "generators", keys: slices.Sorted(maps.Keys(content.Generators))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !slices.Equal(tt.keys, slugs) {
				t.Errorf("keys = %q, want the task slugs %q", tt.keys, slugs)
			}
		})
	}
}

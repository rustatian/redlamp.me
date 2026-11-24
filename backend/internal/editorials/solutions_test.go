package editorials

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"maps"
	"math"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/rustatian/redlamp.me/backend/internal/answers"
	"github.com/rustatian/redlamp.me/backend/internal/bank"
	"github.com/rustatian/redlamp.me/backend/internal/sets"
	"github.com/rustatian/redlamp.me/backend/internal/sets/content"
)

// solutionDriver runs the code on each input with run() of web/harness.py, as the browser worker does, and prints one JSON result for each input.
// It reads all inputs before the first run, because run() replaces sys.stdin. The time includes the parse of the input, as in the browser.
const solutionDriver = `import json, resource, sys, time
resource.setrlimit(resource.RLIMIT_AS, (2 << 30, 2 << 30))
exec(open(sys.argv[1], encoding='utf-8').read())
code, spec, inputs = json.load(sys.stdin)
for value in inputs:
    start = time.perf_counter()
    result = json.loads(run(code, json.dumps(spec), json.dumps({'input': value, 'expected': None})))
    result['seconds'] = time.perf_counter() - start
    print(json.dumps(result), flush=True)
`

// solutionCode finds each solution of an Editorial. The browser shows the text of the code element, so the test decodes the HTML entities.
var solutionCode = regexp.MustCompile(`(?s)<pre class="no-run"><code class="language-python">(.*?)</code></pre>`)

// solutionLimit is the CPython time of one test on a shared CI runner, which runs several solutions at once. It equals the 3 s of web/acceptance.ts.
// The authoring gate is the margin for a slow browser: it allows 1 s in Pyodide on the dev machine.
const solutionLimit = 3 * time.Second

// solutionTest is one input with its expected answer. name is the source and the number of the input, for example "seed 1 test 3".
type solutionTest struct {
	name, input, want string
}

// solutionTests returns the spec and the tests of a problem. A set task has its standard cases and the hidden tests of the seeds 1 and 2^32 - 1
// with the answers of the reference, as the judge does. A bank problem has its examples and tests with the official answers, as a contest does.
func solutionTests(slug string) (sets.TestSpec, []solutionTest, bool) {
	var tests []solutionTest
	for _, set := range content.Sets {
		for _, task := range set.Tasks() {
			if task.Slug() != slug || task.Tests == nil {
				continue
			}
			reference := content.References[slug]
			add := func(source string, inputs []string) {
				for i, input := range inputs {
					tests = append(tests, solutionTest{name: fmt.Sprintf("%s %d", source, i+1), input: input, want: reference.Answer(input)})
				}
			}
			var standard []string
			for _, c := range task.Tests.Cases {
				standard = append(standard, c.Input)
			}
			add("standard test", standard)
			add("seed 1 test", content.Generators[slug](sets.Mulberry32(1)))
			add("seed 4294967295 test", content.Generators[slug](sets.Mulberry32(math.MaxUint32)))
			return task.Tests.TestSpec, tests, true
		}
	}
	for _, problem := range bank.Problems {
		if problem.Slug != slug {
			continue
		}
		for i, example := range problem.Examples {
			tests = append(tests, solutionTest{name: fmt.Sprintf("example %d", i+1), input: example.Input, want: example.Output})
		}
		for i, test := range problem.Tests {
			tests = append(tests, solutionTest{name: fmt.Sprintf("test %d", i+1), input: test.Input, want: test.Output})
		}
		return sets.TestSpec{Mode: "stdio", IgnoreCase: problem.IgnoreCase}, tests, true
	}
	return sets.TestSpec{}, nil, false
}

// solutionHead is the start of a long text for a failure message.
func solutionHead(text string) string {
	if runes := []rune(text); len(runes) > 100 {
		return string(runes[:100]) + "…"
	}
	return text
}

// solutionRun runs the code on all tests in one python3 process and returns the problem of each failed test.
func solutionRun(ctx context.Context, python, code string, spec sets.TestSpec, tests []solutionTest) []string {
	inputs := make([]string, len(tests))
	for i, test := range tests {
		inputs[i] = test.input
	}
	stdin, err := json.Marshal([]any{code, spec, inputs})
	if err != nil {
		return []string{err.Error()}
	}
	// The deadline stops an endless loop. A run within the limits never comes near it.
	ctx, cancel := context.WithTimeout(ctx, solutionLimit*time.Duration(len(tests))+30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, python, "-I", "-c", solutionDriver, "../../../web/harness.py")
	cmd.Stdin = bytes.NewReader(stdin)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	lines := slices.Collect(strings.Lines(string(out)))
	compare := answers.Spec{Mode: spec.Mode, Order: spec.Order, IgnoreCase: spec.IgnoreCase}
	var problems []string
	for i, test := range tests {
		where := fmt.Sprintf("%s, input %q", test.name, solutionHead(test.input))
		if i >= len(lines) {
			return append(problems, fmt.Sprintf("%s: no result, the run stopped: %v %v %s", where, err, ctx.Err(), solutionHead(stderr.String())))
		}
		var result struct {
			Output string `json:"output"`
			Error  *struct {
				Message string `json:"message"`
				Line    *int   `json:"line"`
			} `json:"error"`
			Seconds float64 `json:"seconds"`
		}
		if err := json.Unmarshal([]byte(lines[i]), &result); err != nil {
			return append(problems, fmt.Sprintf("%s: bad result %q: %v", where, solutionHead(lines[i]), err))
		}
		switch {
		case result.Error != nil:
			line := "no line"
			if result.Error.Line != nil {
				line = fmt.Sprintf("line %d", *result.Error.Line)
			}
			problems = append(problems, fmt.Sprintf("%s: error at %s: %s", where, line, result.Error.Message))
		case !answers.Same(compare, test.want, result.Output):
			problems = append(problems, fmt.Sprintf("%s: output %q, want %q", where, solutionHead(result.Output), solutionHead(test.want)))
		case result.Seconds > solutionLimit.Seconds():
			problems = append(problems, fmt.Sprintf("%s: %.2f s, want at most %v", where, result.Seconds, solutionLimit))
		}
	}
	return problems
}

func TestEditorialSolutions(t *testing.T) {
	python, err := exec.LookPath("python3")
	if err != nil {
		if os.Getenv("CI") != "" {
			t.Fatal("python3 is not in PATH")
		}
		t.Skip("python3 is not in PATH")
	}
	for _, slug := range slices.Sorted(maps.Keys(All)) {
		t.Run(slug, func(t *testing.T) {
			t.Parallel()
			spec, tests, ok := solutionTests(slug)
			if !ok {
				t.Fatal("no learning set task and no bank problem has this slug")
			}
			blocks := solutionCode.FindAllStringSubmatch(All[slug].HTML, -1)
			if n := strings.Count(All[slug].HTML, "<pre"); len(blocks) != n {
				t.Fatalf("%d of %d code blocks end with </code></pre>", len(blocks), n)
			}
			for i, m := range blocks {
				for _, problem := range solutionRun(t.Context(), python, html.UnescapeString(m[1]), spec, tests) {
					t.Errorf("solution %d, %s", i+1, problem)
				}
			}
		})
	}
}

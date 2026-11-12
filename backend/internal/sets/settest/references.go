package settest

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"unicode/utf16"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// references.json has the hidden test inputs and the reference answers of the Worker for each learning set task.
// The one-time generator in internal/notebook/testdata/golden wrote it before the port of the references.
//
//go:embed references.json
var references []byte

// Text is one text of references.json. A text of more than 512 characters is in the file as its SHA-256 and its length in UTF-16 code units.
type Text struct {
	Value  *string
	SHA256 string `json:"sha256"`
	Length int    `json:"length"`
}

func (g *Text) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		return json.Unmarshal(data, &g.Value)
	}
	type hashed Text
	return json.Unmarshal(data, (*hashed)(g))
}

// Equal reports whether text is the text of the golden file.
func (g Text) Equal(text string) bool {
	if g.Value != nil {
		return *g.Value == text
	}
	sum := sha256.Sum256([]byte(text))
	return g.SHA256 == hex.EncodeToString(sum[:]) && g.Length == len(utf16.Encode([]rune(text)))
}

// Test is one input with the answer of the reference of the Worker.
type Test struct {
	Input  Text `json:"input"`
	Answer Text `json:"answer"`
}

// Seed has the hidden tests that the generator of the Worker makes for one seed.
type Seed struct {
	Seed  uint32 `json:"seed"`
	Tests []Test `json:"tests"`
}

// Task has the tests of one learning set task: the standard tests of the task and the hidden tests of each seed.
type Task struct {
	// Mode is "call" or "stdio".
	Mode string `json:"mode"`
	// Arity is the number of arguments of a call reference.
	Arity int    `json:"arity"`
	Cases []Test `json:"cases"`
	Seeds []Seed `json:"seeds"`
}

// References returns the tests of the Worker by problem slug.
func References(t testing.TB) map[string]Task {
	t.Helper()
	var tasks map[string]Task
	decoder := json.NewDecoder(bytes.NewReader(references))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&tasks); err != nil {
		t.Fatal(err)
	}
	return tasks
}

func (g Text) String() string {
	if g.Value == nil {
		return fmt.Sprintf("a text of %d characters with the SHA-256 %s", g.Length, g.SHA256)
	}
	return short(*g.Value)
}

// short cuts a long text for a test message.
func short(text string) string {
	if len(text) > 300 {
		return text[:300] + "..."
	}
	return text
}

// Compare checks the references and the generators of one learning set against the tests of the Worker.
// Each task has a reference and a generator. A reference gives the answers of the Worker to the standard tests of the task.
// For each seed the generator gives the hidden test inputs of the Worker, and the reference gives the answers of the Worker to them.
func Compare(t *testing.T, set sets.LearningSet, references map[string]sets.Reference, generators map[string]sets.Generator) {
	t.Helper()
	golden := References(t)
	tasks := set.Tasks()
	if len(references) != len(tasks) || len(generators) != len(tasks) {
		t.Errorf("%d references and %d generators, want one of each for the %d tasks", len(references), len(generators), len(tasks))
	}
	for _, task := range tasks {
		t.Run(task.Slug(), func(t *testing.T) {
			want, ok := golden[task.Slug()]
			if !ok {
				t.Fatal("the Worker has no tests of this task")
			}
			reference, ok := references[task.Slug()]
			if !ok {
				t.Fatal("the task has no reference")
			}
			generate, ok := generators[task.Slug()]
			if !ok {
				t.Fatal("the task has no generator")
			}
			mode := "call"
			if reference.Stdio != nil {
				mode = "stdio"
			}
			if mode != want.Mode || reference.Arity != want.Arity {
				t.Fatalf("a %s reference with %d arguments, want a %s reference with %d", mode, reference.Arity, want.Mode, want.Arity)
			}
			if task.Tests == nil || len(task.Tests.Cases) != len(want.Cases) {
				t.Fatalf("the task does not have the %d standard tests of the Worker", len(want.Cases))
			}
			for i, test := range want.Cases {
				compare(t, reference, fmt.Sprintf("standard test %d", i+1), task.Tests.Cases[i].Input, test)
			}
			for _, seed := range want.Seeds {
				inputs := generate(sets.Mulberry32(seed.Seed))
				if len(inputs) != len(seed.Tests) {
					t.Errorf("seed %d: %d hidden tests, want %d", seed.Seed, len(inputs), len(seed.Tests))
					continue
				}
				for i, test := range seed.Tests {
					compare(t, reference, fmt.Sprintf("seed %d, hidden test %d", seed.Seed, i+1), inputs[i], test)
				}
			}
		})
	}
}

// compare checks one input and the answer of the reference to it.
func compare(t *testing.T, reference sets.Reference, name, input string, want Test) {
	t.Helper()
	if !want.Input.Equal(input) {
		t.Errorf("%s: input = %s, want %s", name, short(input), want.Input)
		return
	}
	defer func() {
		if cause := recover(); cause != nil {
			t.Errorf("%s: the reference panics with %v for the input %s", name, cause, short(input))
		}
	}()
	if answer := reference.Answer(input); !want.Answer.Equal(answer) {
		t.Errorf("%s: answer = %s, want %s for the input %s", name, short(answer), want.Answer, short(input))
	}
}

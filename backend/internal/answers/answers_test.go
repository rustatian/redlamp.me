package answers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"
)

// caseFile is shared/answers.cases.json. The web test of shared/answers.ts reads the same file, so the two copies of the code stay equal.
// The one-time script scripts/golden-answers.ts wrote the file from the TypeScript code.
type caseFile struct {
	Parse []struct {
		Name  string          `json:"name"`
		Text  string          `json:"text"`
		Value json.RawMessage `json:"value"`
	} `json:"parse"`
	RoundTrip []struct {
		Name string `json:"name"`
		Text string `json:"text"`
		Repr string `json:"repr"`
	} `json:"roundTrip"`
	ParseErrors []struct {
		Name string `json:"name"`
		Text string `json:"text"`
	} `json:"parseErrors"`
	Print []struct {
		Name  string          `json:"name"`
		Value json.RawMessage `json:"value"`
		Text  string          `json:"text"`
	} `json:"print"`
	Compare []struct {
		Name string `json:"name"`
		Spec struct {
			Mode       string `json:"mode"`
			Order      string `json:"order"`
			IgnoreCase bool   `json:"ignoreCase"`
		} `json:"spec"`
		Expected string `json:"expected"`
		Output   string `json:"output"`
		Equal    bool   `json:"equal"`
	} `json:"compare"`
}

func readCases(t *testing.T) caseFile {
	t.Helper()
	data, err := os.ReadFile("../../../shared/answers.cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var c caseFile
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	if len(c.Parse) == 0 || len(c.RoundTrip) == 0 || len(c.ParseErrors) == 0 || len(c.Print) == 0 || len(c.Compare) == 0 {
		t.Fatal("shared/answers.cases.json has an empty list of cases")
	}
	return c
}

// fromJSON reads a JSON value as a Python value. A JSON object becomes a dict with the keys in the order of the text.
func fromJSON(t *testing.T, raw json.RawMessage) any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(raw))
	var read func() any
	read = func() any {
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		switch value := token.(type) {
		case json.Delim:
			if value == '[' {
				items := []any{}
				for decoder.More() {
					items = append(items, read())
				}
				_, _ = decoder.Token()
				return items
			}
			dict := &Dict{}
			for decoder.More() {
				key, _ := decoder.Token()
				dict.Set(key.(string), read())
			}
			_, _ = decoder.Token()
			return dict
		default:
			return value
		}
	}
	return read()
}

// equal compares two Python values: the kinds, the items in order, and the dict keys in order.
func equal(a, b any) bool {
	switch x := a.(type) {
	case []any:
		y, ok := b.([]any)
		return ok && slices.EqualFunc(x, y, equal)
	case *Dict:
		y, ok := b.(*Dict)
		if !ok || !slices.Equal(x.Keys(), y.Keys()) {
			return false
		}
		for _, key := range x.Keys() {
			if !equal(x.Get(key), y.Get(key)) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}

func TestParse(t *testing.T) {
	for _, tt := range readCases(t).Parse {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := Parse(tt.Text)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.Text, err)
			}
			if want := fromJSON(t, tt.Value); !equal(got, want) {
				t.Errorf("Parse(%q) = %#v, want %s", tt.Text, got, tt.Value)
			}
		})
	}
}

func TestParseAndRepr(t *testing.T) {
	for _, tt := range readCases(t).RoundTrip {
		t.Run(tt.Name, func(t *testing.T) {
			value, err := Parse(tt.Text)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.Text, err)
			}
			got, err := Repr(value)
			if err != nil {
				t.Fatalf("Repr: %v", err)
			}
			if got != tt.Repr {
				t.Errorf("Repr(Parse(%q)) = %s, want %s", tt.Text, got, tt.Repr)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	for _, tt := range readCases(t).ParseErrors {
		t.Run(tt.Name, func(t *testing.T) {
			if got, err := Parse(tt.Text); err == nil {
				t.Errorf("Parse(%q) = %#v, want an error", tt.Text, got)
			}
		})
	}
}

func TestRepr(t *testing.T) {
	for _, tt := range readCases(t).Print {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := Repr(fromJSON(t, tt.Value))
			if err != nil {
				t.Fatalf("Repr(%s): %v", tt.Value, err)
			}
			if got != tt.Text {
				t.Errorf("Repr(%s) = %s, want %s", tt.Value, got, tt.Text)
			}
		})
	}
}

// A Go value of another kind has no Python literal, as undefined has none in shared/answers.ts.
func TestReprRejectsOtherKinds(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "an integer", value: 5},
		{name: "a list of strings", value: []string{"a"}},
		{name: "a function in a list", value: []any{fmt.Sprint}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Repr(tt.value); err == nil {
				t.Errorf("Repr = %s, want an error", got)
			}
		})
	}
}

func TestSame(t *testing.T) {
	for _, tt := range readCases(t).Compare {
		t.Run(tt.Name, func(t *testing.T) {
			spec := Spec{Mode: tt.Spec.Mode, Order: tt.Spec.Order, IgnoreCase: tt.Spec.IgnoreCase}
			if got := Same(spec, tt.Expected, tt.Output); got != tt.Equal {
				t.Errorf("Same(%+v, %q, %q) = %v, want %v", spec, tt.Expected, tt.Output, got, tt.Equal)
			}
		})
	}
}

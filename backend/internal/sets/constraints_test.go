package sets

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// constraints reports whether fn panics with ErrConstraints. Another panic fails the test.
func constraints(t *testing.T, fn func()) (failed bool) {
	t.Helper()
	defer func() {
		if cause := recover(); cause != nil {
			err, ok := cause.(error)
			if !ok || !errors.Is(err, ErrConstraints) {
				t.Fatalf("panic = %v, want ErrConstraints", cause)
			}
			failed = true
		}
	}()
	fn()
	return false
}

// The rules of src/leetcode/sets/constraints.ts: int takes a whole number of the range, ints takes a list of them with a length of the range, str takes a string that matches.
func TestConstraintChecks(t *testing.T) {
	lower := regexp.MustCompile(`^[a-z]+$`)
	tests := []struct {
		name  string
		check func()
		// wantFail is true for a value that breaks the constraint.
		wantFail bool
	}{
		{name: "whole number at the low bound", check: func() { Int(float64(-3), -3, 5) }},
		{name: "whole number at the high bound", check: func() { Int(float64(5), -3, 5) }},
		{name: "number below the range", check: func() { Int(float64(-4), -3, 5) }, wantFail: true},
		{name: "number above the range", check: func() { Int(float64(6), -3, 5) }, wantFail: true},
		{name: "number with a fraction", check: func() { Int(1.5, -3, 5) }, wantFail: true},
		{name: "boolean in place of a number", check: func() { Int(true, 0, 1) }, wantFail: true},
		{name: "string in place of a number", check: func() { Int("1", 0, 1) }, wantFail: true},
		{name: "list of the shortest length", check: func() { Ints([]any{float64(1)}, 1, 2, 0, 9) }},
		{name: "list above the longest length", check: func() { Ints([]any{float64(1), float64(2), float64(3)}, 1, 2, 0, 9) }, wantFail: true},
		{name: "empty list below the shortest length", check: func() { Ints([]any{}, 1, 2, 0, 9) }, wantFail: true},
		{name: "list with a value outside the range", check: func() { Ints([]any{float64(1), float64(10)}, 1, 2, 0, 9) }, wantFail: true},
		{name: "list with a string", check: func() { Ints([]any{"1"}, 1, 2, 0, 9) }, wantFail: true},
		{name: "number in place of a list", check: func() { Ints(float64(1), 1, 2, 0, 9) }, wantFail: true},
		{name: "string that matches", check: func() { Str("abc", lower) }},
		{name: "string that does not match", check: func() { Str("aBc", lower) }, wantFail: true},
		{name: "number in place of a string", check: func() { Str(float64(1), lower) }, wantFail: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constraints(t, tt.check); got != tt.wantFail {
				t.Errorf("constraint failure = %v, want %v", got, tt.wantFail)
			}
		})
	}
}

// The reader of constraints.ts: tokens between white space, a number of digits only, and no token after the last expected token.
func TestReader(t *testing.T) {
	tests := []struct {
		name  string
		stdin string
		// The test reads two numbers from 0 to 100 and the end of the text.
		wantFail bool
		want     []int
	}{
		{name: "two numbers on two lines", stdin: "3\n100\n", want: []int{3, 100}},
		{name: "no-break space and tabs between the tokens", stdin: "\t3   4 \r\n", want: []int{3, 4}},
		{name: "number with leading zeros", stdin: "007 4", want: []int{7, 4}},
		{name: "missing token", stdin: "3", wantFail: true},
		{name: "token after the last number", stdin: "3 4 5", wantFail: true},
		{name: "negative number", stdin: "3 -4", wantFail: true},
		{name: "number with a fraction", stdin: "3 4.0", wantFail: true},
		{name: "number above the range", stdin: "3 101", wantFail: true},
		{name: "number of 400 digits, above the float range", stdin: "3 " + strings.Repeat("9", 400), wantFail: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			failed := constraints(t, func() {
				r := NewReader(tt.stdin)
				got = []int{r.Int(0, 100), r.Int(0, 100)}
				r.End()
			})
			if failed != tt.wantFail {
				t.Fatalf("constraint failure = %v, want %v", failed, tt.wantFail)
			}
			if !failed && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("numbers = %v, want %v", got, tt.want)
			}
		})
	}
}

// A signed number of constraints.ts is a token of the form -?\d+ in the range.
func TestReaderSigned(t *testing.T) {
	tests := []struct {
		name     string
		stdin    string
		want     int
		wantFail bool
	}{
		{name: "negative number", stdin: "-12", want: -12},
		{name: "number without a sign", stdin: "12", want: 12},
		{name: "negative zero", stdin: "-0", want: 0},
		{name: "minus sign alone", stdin: "-", wantFail: true},
		{name: "two minus signs", stdin: "--12", wantFail: true},
		{name: "plus sign", stdin: "+12", wantFail: true},
		{name: "number below the range", stdin: "-101", wantFail: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int
			failed := constraints(t, func() { got = NewReader(tt.stdin).Signed(-100, 100) })
			if failed != tt.wantFail || got != tt.want {
				t.Errorf("number = %d with the constraint failure %v, want %d with %v", got, failed, tt.want, tt.wantFail)
			}
		})
	}
}

// Call equals JSON.stringify of the arguments, and the answer text of a reference is the Python literal of its result.
func TestCallInputAndAnswer(t *testing.T) {
	sum := Reference{Arity: 2, Call: func(args []any) any {
		total := 0
		for _, n := range Ints(args[0], 0, 10, -100, 100) {
			total += n
		}
		return Py([]any{total, Str(args[1], regexp.MustCompile(`^[ -~]*$`)), [][]int{{1, 2}, {}}, true, nil})
	}}
	tests := []struct {
		name       string
		args       []any
		wantInput  string
		wantAnswer string
	}{
		{name: "numbers and a string", args: []any{[]int{1, -2, 30}, "ab"}, wantInput: `[[1,-2,30],"ab"]`, wantAnswer: "[29, 'ab', [[1, 2], []], True, None]"},
		{name: "characters that HTML escapes", args: []any{[]int{}, "<a & b>"}, wantInput: `[[],"<a & b>"]`, wantAnswer: "[0, '<a & b>', [[1, 2], []], True, None]"},
		{name: "quote and backslash", args: []any{[]int{5}, `"\`}, wantInput: `[[5],"\"\\"]`, wantAnswer: `[5, '"\\', [[1, 2], []], True, None]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Call(tt.args...)
			if input != tt.wantInput {
				t.Errorf("input = %s, want %s", input, tt.wantInput)
			}
			if got := sum.Answer(input); got != tt.wantAnswer {
				t.Errorf("answer = %s, want %s", got, tt.wantAnswer)
			}
		})
	}
}

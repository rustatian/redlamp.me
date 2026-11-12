package sets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/answers"
)

// The checks of a reference for the constraints of the statement. A failed check panics with ErrConstraints.

// Int checks a Python value: a whole number from lo to hi.
func Int(value any, lo, hi int) int {
	number, ok := value.(float64)
	Need(ok && number == math.Trunc(number) && number >= float64(lo) && number <= float64(hi))
	return int(number)
}

// Ints checks a Python value: a list of min to max whole numbers from lo to hi.
func Ints(list any, min, max, lo, hi int) []int {
	values, ok := list.([]any)
	Need(ok && len(values) >= min && len(values) <= max)
	numbers := make([]int, len(values))
	for i, value := range values {
		numbers[i] = Int(value, lo, hi)
	}
	return numbers
}

// Str checks a Python value: a string that matches pattern.
func Str(value any, pattern *regexp.Regexp) string {
	text, ok := value.(string)
	Need(ok && pattern.MatchString(text))
	return text
}

// Reader reads a stdin text as tokens between white space.
type Reader struct {
	tokens []string
	at     int
}

func NewReader(stdin string) *Reader {
	return &Reader{tokens: strings.FieldsFunc(stdin, answers.IsSpace)}
}

// Word returns the next token. The token must match pattern.
func (r *Reader) Word(pattern *regexp.Regexp) string {
	Need(r.at < len(r.tokens))
	r.at++
	return Str(r.tokens[r.at-1], pattern)
}

// number returns the next token as a whole number from lo to hi. The token has digits only, or a minus sign and digits when signed is true.
// It tests the bytes of the token: a regular expression for each token of a large input makes a reference slow.
func (r *Reader) number(lo, hi int, signed bool) int {
	Need(r.at < len(r.tokens))
	token := r.tokens[r.at]
	r.at++
	digits := token
	if signed {
		digits = strings.TrimPrefix(token, "-")
	}
	Need(digits != "")
	for i := range len(digits) {
		Need(digits[i] >= '0' && digits[i] <= '9')
	}
	// A number above the float range is an infinity, which is outside each range.
	value, _ := strconv.ParseFloat(token, 64)
	return Int(value, lo, hi)
}

// Int returns the next token: a whole number from lo to hi, written with digits only.
func (r *Reader) Int(lo, hi int) int {
	return r.number(lo, hi, false)
}

// Signed returns the next token: a whole number from lo to hi, written with digits and a minus sign before a negative number.
func (r *Reader) Signed(lo, hi int) int {
	return r.number(lo, hi, true)
}

// End checks that the text has no token after the last token that the reference read.
func (r *Reader) End() {
	Need(r.at == len(r.tokens))
}

// Py makes a Python value of package answers from a Go value: an int becomes a number, and a slice becomes a list.
func Py(value any) any {
	switch v := value.(type) {
	case nil, bool, string, float64, *answers.Dict:
		return v
	case int:
		return float64(v)
	}
	slice := reflect.ValueOf(value)
	if slice.Kind() != reflect.Slice {
		panic(fmt.Sprintf("no Python value for %T", value))
	}
	list := make([]any, slice.Len())
	for i := range list {
		list[i] = Py(slice.Index(i).Interface())
	}
	return list
}

// Call makes the input of a call task: the Python literal list of the arguments.
// The text equals JSON.stringify of the arguments in JavaScript, for whole numbers, strings of printable ASCII, booleans and lists of them.
func Call(args ...any) string {
	var text bytes.Buffer
	encoder := json.NewEncoder(&text)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(args); err != nil {
		panic(err)
	}
	return strings.TrimSuffix(text.String(), "\n")
}

var errInput = errors.New("the input is not a list of the arguments")

// Answer runs the reference on one input of the task and returns the answer text: the stdout text, or the Python literal of the result.
// It panics when the reference does not accept the input.
func (r Reference) Answer(input string) string {
	if r.Stdio != nil {
		return r.Stdio(input)
	}
	value, err := answers.Parse(input)
	args, ok := value.([]any)
	if err != nil || !ok || len(args) != r.Arity {
		panic(errInput)
	}
	text, err := answers.Repr(r.Call(args))
	if err != nil {
		panic(err)
	}
	return text
}

// Package answers has the Python values of the tests: the parser of Python literals, repr, and the comparison of answers.
// It is the Go copy of shared/answers.ts. The file shared/answers.cases.json holds the cases of both copies.
//
// A value is nil (None), bool, float64, string, []any (a list or a tuple) or *Dict.
package answers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Dict is a Python dict with str keys. It has the key order of a JavaScript object, because the web client holds the same value as an object:
// first the keys that are list indexes (0 to 4294967294) in number order, then the other keys in the order of their first use.
type Dict struct {
	keys   []string
	values map[string]any
}

// Set stores a value. A key that exists keeps its place.
func (d *Dict) Set(key string, value any) {
	if d.values == nil {
		d.values = map[string]any{}
	}
	if _, ok := d.values[key]; !ok {
		d.keys = append(d.keys, key)
	}
	d.values[key] = value
}

// Get returns the value of a key, or nil.
func (d *Dict) Get(key string) any { return d.values[key] }

var indexKey = regexp.MustCompile(`^(?:0|[1-9][0-9]{0,9})$`)

// index returns the number of a key that is a list index.
func index(key string) (uint64, bool) {
	if !indexKey.MatchString(key) {
		return 0, false
	}
	n, _ := strconv.ParseUint(key, 10, 64)
	return n, n <= 4294967294
}

// Keys returns the keys in the order of a JavaScript object.
func (d *Dict) Keys() []string {
	var indexes, names []string
	for _, key := range d.keys {
		if _, ok := index(key); ok {
			indexes = append(indexes, key)
		} else {
			names = append(names, key)
		}
	}
	slices.SortFunc(indexes, func(a, b string) int {
		x, _ := index(a)
		y, _ := index(b)
		return int(x) - int(y)
	})
	return append(indexes, names...)
}

// space is the white space of a JavaScript regular expression.
const space = `\t\n\v\f\r \x{a0}\x{1680}\x{2000}-\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000}\x{feff}`

var (
	// token matches one token: a number, a string in single or double quotes, a name or a mark. White space can come before each token.
	token   = regexp.MustCompile(`^[` + space + `]*(?:(-?(?:\d+\.?\d*|\.\d+)(?:[eE][-+]?\d+)?)|'((?:[^'\\\n]|\\.)*)'|"((?:[^"\\\n]|\\.)*)"|(True|False|None)\b|([\[\](){},:]))`)
	escapes = map[rune]rune{'\\': '\\', '\'': '\'', '"': '"', 'n': '\n', 'r': '\r', 't': '\t'}
)

// isSpace reports whether a character is white space of JavaScript, as in the class space.
func isSpace(char rune) bool {
	switch char {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0xa0, 0x1680, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000, 0xfeff:
		return true
	}
	return char >= 0x2000 && char <= 0x200a
}

// match is one token. kind is 'n' (number), 's' (string), 'w' (name) or the mark itself.
type match struct {
	kind byte
	text string
}

type parser struct {
	text string
	at   int
}

func (p *parser) fail() error {
	return fmt.Errorf("no Python literal at byte %d", p.at+1)
}

func (p *parser) next() (match, error) {
	groups := token.FindStringSubmatchIndex(p.text[p.at:])
	if groups == nil {
		return match{}, p.fail()
	}
	group := func(i int) string { return p.text[p.at+groups[2*i] : p.at+groups[2*i+1]] }
	var m match
	switch {
	case groups[2] >= 0:
		m = match{'n', group(1)}
	case groups[4] >= 0:
		m = match{'s', group(2)}
	case groups[6] >= 0:
		m = match{'s', group(3)}
	case groups[8] >= 0:
		m = match{'w', group(4)}
	default:
		m = match{group(5)[0], ""}
	}
	p.at += groups[1]
	return m, nil
}

// sequence reads the items to the close mark. A comma can follow the last item.
func (p *parser) sequence(closing byte, item func(match) error) (comma bool, err error) {
	m, err := p.next()
	for err == nil && m.kind != closing {
		if err = item(m); err != nil {
			return false, err
		}
		if m, err = p.next(); err != nil {
			return false, err
		}
		if m.kind == ',' {
			comma = true
			m, err = p.next()
		} else if m.kind != closing {
			return false, p.fail()
		}
	}
	return comma, err
}

func (p *parser) list(closing byte) ([]any, bool, error) {
	items := []any{}
	comma, err := p.sequence(closing, func(m match) error {
		item, err := p.value(m)
		items = append(items, item)
		return err
	})
	return items, comma, err
}

func (p *parser) value(m match) (any, error) {
	switch m.kind {
	case 'n':
		// A number above the float range is an infinity, as in JavaScript.
		number, _ := strconv.ParseFloat(m.text, 64)
		return number, nil
	case 's':
		var text strings.Builder
		escaped := false
		for _, char := range m.text {
			switch {
			case escaped:
				// Other escapes of Python (\x41, octal, \N) are not Python literals of the tests.
				plain, ok := escapes[char]
				if !ok {
					return nil, p.fail()
				}
				text.WriteRune(plain)
				escaped = false
			case char == '\\':
				escaped = true
			default:
				text.WriteRune(char)
			}
		}
		return text.String(), nil
	case 'w':
		switch m.text {
		case "True":
			return true, nil
		case "False":
			return false, nil
		}
		return nil, nil
	case '[':
		items, _, err := p.list(']')
		return items, err
	case '(':
		// (5) is the number 5, and (5,) is a tuple.
		items, comma, err := p.list(')')
		if err == nil && len(items) == 1 && !comma {
			return items[0], nil
		}
		return items, err
	case '{':
		dict := &Dict{}
		_, err := p.sequence('}', func(m match) error {
			key, err := p.value(m)
			if err != nil {
				return err
			}
			name, ok := key.(string)
			if !ok {
				return p.fail()
			}
			if colon, err := p.next(); err != nil || colon.kind != ':' {
				return p.fail()
			}
			first, err := p.next()
			if err != nil {
				return err
			}
			value, err := p.value(first)
			dict.Set(name, value)
			return err
		})
		return dict, err
	}
	return nil, p.fail()
}

// Parse reads one Python literal: int, float, str, True, False, None, list, tuple (as a list) and dict with str keys.
func Parse(text string) (any, error) {
	p := &parser{text: text}
	first, err := p.next()
	if err != nil {
		return nil, err
	}
	value, err := p.value(first)
	if err != nil {
		return nil, err
	}
	if strings.TrimFunc(text[p.at:], isSpace) != "" {
		return nil, p.fail()
	}
	return value, nil
}

// number prints a float as JavaScript prints a number, so 2.0 prints as 2: https://tc39.es/ecma262/#sec-numeric-types-number-tostring
func number(f float64) string {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	case f == 0:
		return "0"
	}
	if abs := math.Abs(f); abs >= 1e-6 && abs < 1e21 {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	// JavaScript writes the exponent without a leading zero: 1e-7 and 1e+21.
	mantissa, exponent, _ := strings.Cut(strconv.FormatFloat(f, 'e', -1, 64), "e")
	return mantissa + "e" + exponent[:1] + strings.TrimLeft(exponent[1:], "0")
}

// Repr prints a value as Python repr prints it.
func Repr(value any) (string, error) {
	switch v := value.(type) {
	case nil:
		return "None", nil
	case bool:
		if v {
			return "True", nil
		}
		return "False", nil
	case float64:
		return number(v), nil
	case string:
		// Python takes double quotes when the text has a single quote and no double quote.
		quote := '\''
		if strings.Contains(v, "'") && !strings.Contains(v, `"`) {
			quote = '"'
		}
		var text strings.Builder
		text.WriteRune(quote)
		for _, char := range v {
			switch {
			case char == '\\' || char == quote:
				text.WriteRune('\\')
				text.WriteRune(char)
			case char == '\n':
				text.WriteString(`\n`)
			case char == '\r':
				text.WriteString(`\r`)
			case char == '\t':
				text.WriteString(`\t`)
			case char < 0x20 || char == 0x7f:
				fmt.Fprintf(&text, `\x%02x`, char)
			default:
				text.WriteRune(char)
			}
		}
		text.WriteRune(quote)
		return text.String(), nil
	case []any:
		items, err := reprs(v)
		return "[" + strings.Join(items, ", ") + "]", err
	case *Dict:
		keys := v.Keys()
		items := make([]string, len(keys))
		for i, key := range keys {
			name, _ := Repr(key)
			item, err := Repr(v.Get(key))
			if err != nil {
				return "", err
			}
			items[i] = name + ": " + item
		}
		return "{" + strings.Join(items, ", ") + "}", nil
	}
	return "", fmt.Errorf("no Python literal for %T", value)
}

func reprs(values []any) ([]string, error) {
	items := make([]string, len(values))
	for i, value := range values {
		item, err := Repr(value)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}
	return items, nil
}

// Spec is the part of a test spec that the comparison reads. Mode is "call" or "stdio". Order is "", "exact", "any" or "any-nested".
type Spec struct {
	Mode       string
	Order      string
	IgnoreCase bool
}

var errNotList = errors.New("the answer is not a list")

// sorted returns a text for the items of a list without their order.
func sorted(items []string) string {
	slices.Sort(items)
	// A list of strings always encodes.
	text, _ := json.Marshal(items)
	return string(text)
}

// canonical returns a text that is equal for equal answers. Stdout compares token by token. Other answers compare Python literals with the list order of the spec.
func canonical(spec Spec, answer string) (string, error) {
	if spec.Mode == "stdio" {
		tokens := strings.Join(strings.FieldsFunc(answer, isSpace), " ")
		if spec.IgnoreCase {
			// The lower case of Unicode without a language, as toLowerCase of JavaScript.
			tokens = cases.Lower(language.Und).String(tokens)
		}
		return tokens, nil
	}
	value, err := Parse(answer)
	if err != nil {
		return "", err
	}
	if spec.Order == "" || spec.Order == "exact" {
		return Repr(value)
	}
	list, ok := value.([]any)
	if !ok {
		return "", errNotList
	}
	if spec.Order == "any" {
		items, err := reprs(list)
		return sorted(items), err
	}
	outer := make([]string, len(list))
	for i, item := range list {
		inner, ok := item.([]any)
		if !ok {
			return "", errNotList
		}
		items, err := reprs(inner)
		if err != nil {
			return "", err
		}
		outer[i] = sorted(items)
	}
	return sorted(outer), nil
}

// Same reports whether the output equals the expected answer for a test spec.
// An output that is not a Python literal, or not a list for a free order, is a wrong answer.
func Same(spec Spec, expected, output string) bool {
	want, err := canonical(spec, expected)
	if err != nil {
		return false
	}
	got, err := canonical(spec, output)
	return err == nil && got == want
}

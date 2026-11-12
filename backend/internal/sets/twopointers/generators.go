package twopointers

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

const lower = "abcdefghijklmnopqrstuvwxyz"

// printable has the 95 printable ASCII characters from the space to the tilde.
var printable = func() string {
	var out strings.Builder
	for c := byte(32); c < 127; c++ {
		out.WriteByte(c)
	}
	return out.String()
}()

// between returns a random integer in [lo, hi].
func between(random func() float64, lo, hi int) int {
	return lo + int(math.Floor(random()*float64(hi-lo+1)))
}

func list(random func() float64, n, lo, hi int) []int {
	values := make([]int, n)
	for i := range values {
		values[i] = between(random, lo, hi)
	}
	return values
}

// asc sorts values in place and returns them.
func asc(values []int) []int {
	slices.Sort(values)
	return values
}

func charOf(random func() float64, chars string) string {
	i := between(random, 0, len(chars)-1)
	return chars[i : i+1]
}

func textOf(random func() float64, n int, chars string) string {
	var out strings.Builder
	for range n {
		out.WriteString(charOf(random, chars))
	}
	return out.String()
}

func call(cases [][]any) []string {
	inputs := make([]string, len(cases))
	for i, args := range cases {
		inputs[i] = sets.Call(args...)
	}
	return inputs
}

// Generators make the hidden tests of the set "two-pointers". Each generator gives the edge cases of the statement, then random medium inputs, then 2 inputs at the maximum sizes.
// A call input is a Python literal list of the arguments. JSON of integers and printable ASCII strings is a Python literal too.
var Generators = map[string]sets.Generator{
	"two-sum-ii-input-array-is-sorted": func(random func() float64) []string {
		// The pair x + y = t is the answer. Of each other pair of values with the sum t, the pool keeps one value, so the answer is unique.
		// spread keeps all values near t / 2, so a small spread gives many duplicates.
		input := func(n, t, spread int) []any {
			x := between(random, max(-1000, t-1000, int(math.Ceil(float64(t-spread)/2))), int(math.Floor(float64(t)/2)))
			y := t - x
			var pool []int
			for v := -1000; v <= 1000; v++ {
				w := t - v
				if v == x || v == y || 2*v == t || max(2*v-t, t-2*v) > spread {
					continue
				}
				if w < -1000 || w > 1000 {
					pool = append(pool, v)
				} else if v < w {
					if random() < 0.5 {
						pool = append(pool, v)
					} else {
						pool = append(pool, w)
					}
				}
			}
			nums := []int{x, y}
			for range n - 2 {
				nums = append(nums, pool[between(random, 0, len(pool)-1)])
			}
			return []any{asc(nums), t}
		}
		cases := [][]any{{[]int{-1000, 1000}, 0}, {[]int{-500, -500}, -1000}, {[]int{500, 500}, 1000}, {[]int{0, 0, 3, 4}, 0}, {[]int{1, 3, 5, 7, 100}, 101}, {[]int{-3, -1, 0, 2, 2, 2, 5}, 5}}
		for range 5 {
			cases = append(cases, input(between(random, 3, 60), between(random, -1000, 1000), between(random, 10, 4000)))
		}
		cases = append(cases, input(30000, between(random, -1000, 1000), 4000), input(30000, between(random, -1000, 1000), 40))
		return call(cases)
	},
	"valid-palindrome": func(random func() float64) []string {
		alnums := lower + strings.ToUpper(lower) + "0123456789"
		other := strings.Map(func(c rune) rune {
			if alnum(byte(c)) {
				return -1
			}
			return c
		}, printable)
		// The right half mirrors the letters and digits of the left half in random case, with other random characters between them. n is 2 or more.
		input := func(n int, palindrome bool) []any {
			left := make([]string, n>>1)
			for i := range left {
				chars := other
				if random() < 0.7 {
					chars = alnums
				}
				left[i] = charOf(random, chars)
			}
			right := make([]string, len(left))
			for i, c := range left {
				switch {
				case !alnum(c[0]):
					right[i] = charOf(random, other)
				case random() < 0.5:
					right[i] = strings.ToUpper(c)
				default:
					right[i] = strings.ToLower(c)
				}
			}
			slices.Reverse(right)
			if !palindrome {
				i := between(random, 0, len(right)-1)
				if strings.ToLower(right[i]) == "a" {
					right[i] = "b"
				} else {
					right[i] = "a"
				}
			}
			middle := ""
			if n%2 != 0 {
				middle = charOf(random, printable)
			}
			return []any{strings.Join(left, "") + middle + strings.Join(right, "")}
		}
		cases := [][]any{{" "}, {"a"}, {".,"}, {"0P"}, {"Aa"}, {"ab@#BA"}, {`\'"`}, {"1a2"}}
		for range 4 {
			cases = append(cases, input(between(random, 2, 60), random() < 0.5))
		}
		cases = append(cases, []any{textOf(random, between(random, 1, 60), printable)}, input(200000, true), input(200000, false))
		return call(cases)
	},
	"max-number-of-k-sum-pairs": func(random func() float64) []string {
		// Most values are x or k - x for a few random x. The other values are random noise up to hi.
		input := func(n, k, kinds, hi int) []any {
			xs := list(random, kinds, 1, k-1)
			nums := make([]int, n)
			for i := range nums {
				x := xs[between(random, 0, kinds-1)]
				switch {
				case random() < 0.2:
					nums[i] = between(random, 1, hi)
				case random() < 0.5:
					nums[i] = x
				default:
					nums[i] = k - x
				}
			}
			return []any{nums, k}
		}
		cases := [][]any{{[]int{5}, 10}, {[]int{1, 1}, 2}, {[]int{1, 1, 1}, 2}, {[]int{2, 2, 2, 2}, 4}, {[]int{1, 1}, 1}, {[]int{1e9, 1e9}, 1000000000}, {[]int{1, 999999999}, 1000000000}}
		for range 5 {
			cases = append(cases, input(between(random, 2, 60), between(random, 2, 20), between(random, 1, 4), 30))
		}
		cases = append(cases, input(100000, between(random, 2, 1e9), 50, 1e9), input(100000, 1000, 500, 1000))
		return call(cases)
	},
	"container-with-most-water": func(random func() float64) []string {
		cases := [][]any{{[]int{1, 1}}, {[]int{0, 0}}, {[]int{0, 10000}}, {[]int{10000, 10000}}, {[]int{1, 2, 3, 4, 5}}, {[]int{5, 4, 3, 2, 1}}, {[]int{2, 2, 2, 2}}, {[]int{10000, 0, 0, 0, 10000}}}
		for range 4 {
			cases = append(cases, []any{list(random, between(random, 3, 60), 0, between(random, 1, 10000))})
		}
		cases = append(cases, []any{list(random, 100000, 0, 10000)}, []any{slices.Repeat([]int{10000}, 100000)})
		return call(cases)
	},
	"boats-to-save-people": func(random func() float64) []string {
		cases := [][]any{{[]int{1}, 1}, {[]int{1, 1}, 2}, {[]int{1, 1}, 1}, {[]int{30000, 30000}, 30000}, {[]int{1, 29999}, 30000}, {[]int{2, 2, 2, 2}, 4}, {[]int{5, 1, 4, 2}, 6}}
		for range 4 {
			limit := between(random, 1, 100)
			cases = append(cases, []any{list(random, between(random, 2, 60), 1, limit), limit})
		}
		cases = append(cases, []any{list(random, 50000, 1, 30000), 30000}, []any{list(random, 50000, 15000, 30000), 30000})
		return call(cases)
	},
	"remove-duplicates-from-sorted-array": func(random func() float64) []string {
		all := make([]int, 201)
		for i := range all {
			all[i] = i - 100
		}
		cases := [][]any{{[]int{7}}, {[]int{-100}}, {[]int{100, 100}}, {[]int{-100, 100}}, {[]int{1, 1, 1, 1}}, {all}}
		for range 4 {
			lo := between(random, -100, 100)
			cases = append(cases, []any{asc(list(random, between(random, 2, 60), lo, min(100, lo+between(random, 0, 20))))})
		}
		cases = append(cases, []any{asc(list(random, 30000, -100, 100))}, []any{slices.Repeat([]int{100}, 30000)})
		return call(cases)
	},
	"remove-element": func(random func() float64) []string {
		// val is a value of nums in most inputs.
		input := func(n, hi int) []any {
			nums := list(random, n, 0, hi)
			if random() < 0.8 {
				return []any{nums, nums[between(random, 0, n-1)]}
			}
			return []any{nums, between(random, 0, 100)}
		}
		cases := [][]any{{[]int{}, 0}, {[]int{}, 100}, {[]int{1}, 1}, {[]int{1}, 2}, {[]int{50, 50, 50}, 50}, {[]int{0, 1, 2}, 100}, {[]int{0, 0, 1}, 0}}
		for range 4 {
			cases = append(cases, input(between(random, 1, 30), between(random, 0, 50)))
		}
		cases = append(cases, input(100, 50), []any{slices.Repeat([]int{37}, 100), 37})
		return call(cases)
	},
	"move-zeroes": func(random func() float64) []string {
		input := func(n int, zeros float64, lo, hi int) []any {
			nums := make([]int, n)
			for i := range nums {
				if random() >= zeros {
					nums[i] = between(random, lo, hi)
				}
			}
			return []any{nums}
		}
		cases := [][]any{{[]int{0}}, {[]int{1}}, {[]int{0, 0, 0}}, {[]int{1, 2, 3}}, {[]int{math.MinInt32, 0, math.MaxInt32}}, {[]int{0, 0, 1}}, {[]int{1, 0, 0}}, {[]int{0, -1, 0, -2}}}
		for range 4 {
			cases = append(cases, input(between(random, 1, 60), random(), -1000, 1000))
		}
		cases = append(cases, input(10000, 0.3, math.MinInt32, math.MaxInt32), []any{append(make([]int, 5000), list(random, 5000, 1, math.MaxInt32)...)})
		return call(cases)
	},
	"remove-duplicates-from-sorted-array-ii": func(random func() float64) []string {
		cases := [][]any{{[]int{5}}, {[]int{5, 5}}, {[]int{5, 5, 5}}, {[]int{-10000, -10000, -10000, 10000, 10000, 10000}}, {[]int{1, 2, 3}}, {[]int{0, 0, 0, 0, 0, 0}}, {[]int{1, 1, 2, 2, 2, 3, 3, 3, 3}}}
		for range 4 {
			lo := between(random, -10000, 9990)
			cases = append(cases, []any{asc(list(random, between(random, 2, 60), lo, lo+between(random, 0, 10)))})
		}
		cases = append(cases, []any{asc(list(random, 30000, -10000, 10000))}, []any{asc(list(random, 30000, -5, 5))})
		return call(cases)
	},
	"sort-colors": func(random func() float64) []string {
		cases := [][]any{{[]int{0}}, {[]int{1}}, {[]int{2}}, {[]int{2, 1, 0}}, {[]int{2, 2, 2}}, {[]int{0, 0, 1, 1}}, {[]int{1, 0}}, {[]int{2, 0, 2, 0}}}
		for range 4 {
			cases = append(cases, []any{list(random, between(random, 2, 60), 0, 2)})
		}
		cases = append(cases, []any{list(random, 300, 0, 2)}, []any{slices.Concat(slices.Repeat([]int{2}, 100), slices.Repeat([]int{1}, 100), make([]int, 100))})
		return call(cases)
	},
	"merge-sorted-array": func(random func() float64) []string {
		// nums1 holds m sorted values and n zeros at the end.
		input := func(m, n, hi int) []any {
			nums1 := append(asc(list(random, m, -hi, hi)), make([]int, n)...)
			return []any{nums1, m, asc(list(random, n, -hi, hi)), n}
		}
		cases := [][]any{
			{[]int{1}, 1, []int{}, 0}, {[]int{0}, 0, []int{1}, 1}, {[]int{4, 5, 6, 0, 0, 0}, 3, []int{1, 2, 3}, 3}, {[]int{1, 2, 3, 0, 0, 0}, 3, []int{4, 5, 6}, 3},
			{[]int{2, 2, 0, 0}, 2, []int{2, 2}, 2}, {[]int{-1e9, 1e9, 0}, 2, []int{0}, 1}, {[]int{0, 0, 0}, 0, []int{-5, -5, -5}, 3},
		}
		for range 4 {
			m := between(random, 0, 30)
			least := 0
			if m == 0 {
				least = 1
			}
			n := between(random, least, 30)
			hi := 1000000000
			if random() < 0.5 {
				hi = 20
			}
			cases = append(cases, input(m, n, hi))
		}
		cases = append(cases, input(100, 100, 1e9), input(1, 199, 1e9))
		return call(cases)
	},
	"is-subsequence": func(random func() float64) []string {
		// s takes size letters of t in order (selection sampling). A changed letter can break the subsequence.
		sample := func(t string, size int) string {
			var s []byte
			for j := 0; j < len(t) && len(s) < size; j++ {
				if random()*float64(len(t)-j) < float64(size-len(s)) {
					s = append(s, t[j])
				}
			}
			return string(s)
		}
		input := func(t string, size int, change bool) []any {
			s := sample(t, size)
			i := between(random, 0, len(s)-1)
			if change && s != "" {
				s = s[:i] + charOf(random, "abc") + s[i+1:]
			}
			return []any{s, t}
		}
		big := textOf(random, 10000, lower)
		noZ := textOf(random, 10000, lower[:25])
		cases := [][]any{{"", ""}, {"", "abc"}, {"a", ""}, {"abc", "abc"}, {"aaaa", "aaa"}, {"b", "abc"}, {"abc", "cba"}, {"z", "aaaaz"}}
		for range 4 {
			n := between(random, 0, 200)
			chars := lower
			if random() < 0.5 {
				chars = "abc"
			}
			t := textOf(random, n, chars)
			cases = append(cases, input(t, between(random, 0, min(10, len(t))), random() < 0.5))
		}
		cases = append(cases, []any{sample(big, 100), big}, []any{sample(noZ, 99) + "z", noZ})
		return call(cases)
	},
	"merge-strings-alternately": func(random func() float64) []string {
		cases := [][]any{{"a", "b"}, {"a", "bcdef"}, {"abcde", "f"}, {"ab", "ab"}, {strings.Repeat("a", 100), "b"}, {"x", strings.Repeat("y", 100)}}
		for range 4 {
			cases = append(cases, []any{textOf(random, between(random, 1, 30), lower), textOf(random, between(random, 1, 30), lower)})
		}
		cases = append(cases, []any{textOf(random, 100, lower), textOf(random, 100, lower)}, []any{textOf(random, 100, lower), textOf(random, 100, lower)})
		return call(cases)
	},
	"intersection-of-two-arrays": func(random func() float64) []string {
		cases := [][]any{{[]int{0}, []int{0}}, {[]int{1000}, []int{0}}, {[]int{1, 1, 1}, []int{1}}, {[]int{1, 2, 3}, []int{4, 5, 6}}, {[]int{5, 4, 3}, []int{3, 4, 5}}, {[]int{0, 1000}, []int{1000, 0, 0}}}
		for range 4 {
			hi := between(random, 0, 30)
			cases = append(cases, []any{list(random, between(random, 1, 40), 0, hi), list(random, between(random, 1, 40), 0, hi)})
		}
		cases = append(cases, []any{list(random, 1000, 0, 1000), list(random, 1000, 0, 1000)}, []any{list(random, 1000, 0, 30), list(random, 1000, 0, 30)})
		return call(cases)
	},
	"cf-1538-c": func(random func() float64) []string {
		// One test set: the line "n l r", then the line of a. l and r stay near the sums of two values.
		block := func(n, hi int) string {
			top := min(2*hi, 1e9)
			l := between(random, 1, top)
			r := between(random, l, top)
			return fmt.Sprintf("%d %d %d\n%s\n", n, l, r, strings.Trim(fmt.Sprint(list(random, n, 1, hi)), "[]"))
		}
		input := func(blocks []string) string { return strconv.Itoa(len(blocks)) + "\n" + strings.Join(blocks, "") }
		inputs := []string{
			"1\n1 1 1\n1\n", "1\n2 2 2\n1 1\n", "1\n2 1 1000000000\n1000000000 1000000000\n", "1\n4 2 1000000000\n1 1 1 1\n", "1\n5 1 1\n1 1 1 1 1\n",
			"3\n3 4 7\n5 1 2\n1 1 1\n7\n2 1000000000 1000000000\n500000000 500000000\n",
		}
		for range 4 {
			blocks := make([]string, between(random, 1, 5))
			for i := range blocks {
				blocks[i] = block(between(random, 1, 30), between(random, 1, 20))
			}
			inputs = append(inputs, input(blocks))
		}
		inputs = append(inputs, input([]string{block(200000, 1e9)}))
		blocks := make([]string, 10000)
		for i := range blocks {
			blocks[i] = block(20, 1e9)
		}
		return append(inputs, input(blocks))
	},
	"3sum": func(random func() float64) []string {
		cases := [][]any{{[]int{0, 0, 0}}, {[]int{0, 1, 1}}, {[]int{0, 0, 0, 0}}, {[]int{1, 2, 3}}, {[]int{-100000, 50000, 50000}}, {[]int{-100000, 0, 100000}}, {[]int{-2, 0, 1, 1, 2}}, {[]int{-1, -1, -1, 2, 2, 2}}}
		for range 4 {
			cases = append(cases, []any{list(random, between(random, 3, 60), -10, 10)})
		}
		cases = append(cases, []any{list(random, 3000, -100000, 100000)})
		steps := list(random, 3000, -2, 2)
		for i := range steps {
			steps[i] *= 50000
		}
		return call(append(cases, []any{steps}))
	},
	"string-compression": func(random func() float64) []string {
		// Groups of one random character with random lengths up to most. Two near groups can have the same character and join.
		input := func(n int, chars string, most int) []any {
			out := []string{}
			for len(out) < n {
				size := min(between(random, 1, most), n-len(out))
				out = append(out, slices.Repeat([]string{charOf(random, chars)}, size)...)
			}
			return []any{out}
		}
		cases := [][]any{{[]string{"a"}}, {[]string{"a", "a"}}, {slices.Repeat([]string{"a"}, 10)}, {[]string{"a", "b", "c"}}, {[]string{"1", "1", "1"}}, {[]string{"'", "'", `\`}}, {[]string{"a", "A", "a"}}, {slices.Repeat([]string{" "}, 100)}}
		for range 4 {
			cases = append(cases, input(between(random, 1, 60), "aAb1 '", 12))
		}
		cases = append(cases, []any{slices.Repeat([]string{"z"}, 2000)}, input(2000, printable, 150))
		return call(cases)
	},
}

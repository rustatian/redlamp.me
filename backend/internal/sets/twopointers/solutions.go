package twopointers

import (
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

var (
	printableText = regexp.MustCompile(`^[\x20-\x7e]*$`)
	printableChar = regexp.MustCompile(`^[\x20-\x7e]$`)
	lowerText     = regexp.MustCompile(`^[a-z]*$`)
)

// sorted checks that list does not decrease.
func sorted(list []int) []int {
	sets.Need(slices.IsSorted(list))
	return list
}

// alnum reports whether c is an ASCII letter or digit.
func alnum(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z'
}

// References are the reference solutions of the tasks of the set.
var References = map[string]sets.Reference{
	// Expected O(n) time and O(n) storage for uniqueness validation. Two Pointers then uses constant auxiliary state.
	"two-sum-ii-input-array-is-sorted": {Arity: 2, Call: func(args []any) any {
		a := sorted(sets.Ints(args[0], 2, 3e4, -1000, 1000))
		t := sets.Int(args[1], -1000, 1000)
		seen := map[int]int{}
		pairs, l, r := 0, 0, len(a)-1
		for _, x := range a {
			pairs += seen[t-x]
			seen[x]++
		}
		sets.Need(pairs == 1)
		for a[l]+a[r] != t {
			if a[l]+a[r] < t {
				l++
			} else {
				r--
			}
		}
		return sets.Py([]int{l + 1, r + 1})
	}},
	// O(n): l and r skip the characters that are not letters or digits and compare the rest without case.
	"valid-palindrome": {Arity: 1, Call: func(args []any) any {
		t := sets.Str(args[0], printableText)
		sets.Need(len(t) >= 1 && len(t) <= 200000)
		low := func(c byte) byte {
			if c >= 'A' && c <= 'Z' {
				return c + 32
			}
			return c
		}
		l, r := 0, len(t)-1
		for l < r {
			switch {
			case !alnum(t[l]):
				l++
			case !alnum(t[r]):
				r--
			case low(t[l]) != low(t[r]):
				return false
			default:
				l++
				r--
			}
		}
		return true
	}},
	// Expected O(n) time: a Hash Table keeps the numbers that wait for a pair.
	"max-number-of-k-sum-pairs": {Arity: 2, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1e5, 1, 1e9)
		sum := sets.Int(args[1], 1, 1e9)
		wait := map[int]int{}
		ops := 0
		for _, x := range a {
			if c := wait[sum-x]; c != 0 {
				wait[sum-x] = c - 1
				ops++
			} else {
				wait[x]++
			}
		}
		return sets.Py(ops)
	}},
	// O(n): move the lower line, because no narrower pair with it holds more water.
	"container-with-most-water": {Arity: 1, Call: func(args []any) any {
		h := sets.Ints(args[0], 2, 1e5, 0, 1e4)
		l, r, best := 0, len(h)-1, 0
		for l < r {
			best = max(best, min(h[l], h[r])*(r-l))
			if h[l] < h[r] {
				l++
			} else {
				r--
			}
		}
		return sets.Py(best)
	}},
	// O(n log n): sort, then the heaviest person takes the lightest person when the pair fits.
	"boats-to-save-people": {Arity: 2, Call: func(args []any) any {
		lim := sets.Int(args[1], 1, 3e4)
		p := sets.Ints(args[0], 1, 5e4, 1, lim)
		slices.Sort(p)
		l, boats := 0, 0
		for r := len(p) - 1; l <= r; r, boats = r-1, boats+1 {
			if p[l]+p[r] <= lim {
				l++
			}
		}
		return sets.Py(boats)
	}},
	// O(n) time and O(1) scan state. The returned prefix copy stores k output entries.
	"remove-duplicates-from-sorted-array": {Arity: 1, Call: func(args []any) any {
		a := sorted(sets.Ints(args[0], 1, 3e4, -100, 100))
		k := 1
		for i := 1; i < len(a); i++ {
			if a[i] != a[k-1] {
				a[k] = a[i]
				k++
			}
		}
		return sets.Py(a[:k])
	}},
	// O(n) time and O(1) scan state. The returned prefix copy stores k output entries.
	"remove-element": {Arity: 2, Call: func(args []any) any {
		a := sets.Ints(args[0], 0, 100, 0, 50)
		v := sets.Int(args[1], 0, 100)
		k := 0
		for i := range a {
			if a[i] != v {
				a[k] = a[i]
				k++
			}
		}
		return sets.Py(a[:k])
	}},
	// O(n): the write pointer keeps the values that are not zero in order, and zeros fill the tail.
	"move-zeroes": {Arity: 1, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1e4, math.MinInt32, math.MaxInt32)
		k := 0
		for i := range a {
			if a[i] != 0 {
				a[k] = a[i]
				k++
			}
		}
		clear(a[k:])
		return sets.Py(a)
	}},
	// O(n) time and O(1) scan state. The returned prefix copy stores k output entries.
	"remove-duplicates-from-sorted-array-ii": {Arity: 1, Call: func(args []any) any {
		a := sorted(sets.Ints(args[0], 1, 3e4, -1e4, 1e4))
		k := 0
		for i := range a {
			if k < 2 || a[i] != a[k-2] {
				a[k] = a[i]
				k++
			}
		}
		return sets.Py(a[:k])
	}},
	// O(n) in one pass: the Dutch national flag partition with three pointers.
	"sort-colors": {Arity: 1, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 300, 0, 2)
		lo, mid, hi := 0, 0, len(a)-1
		for mid <= hi {
			switch a[mid] {
			case 0:
				a[lo], a[mid] = a[mid], a[lo]
				lo++
				mid++
			case 1:
				mid++
			default:
				a[mid], a[hi] = a[hi], a[mid]
				hi--
			}
		}
		return sets.Py(a)
	}},
	// O(m + n) time. Validation slices use O(m + n) temporary entries; the backward Merge uses O(1) auxiliary state.
	"merge-sorted-array": {Arity: 4, Call: func(args []any) any {
		mm := sets.Int(args[1], 0, 200)
		nn := sets.Int(args[3], 0, 200)
		sets.Need(mm+nn >= 1 && mm+nn <= 200)
		a := sets.Ints(args[0], mm+nn, mm+nn, -1e9, 1e9)
		b := sorted(sets.Ints(args[2], nn, nn, -1e9, 1e9))
		sorted(a[:mm])
		// The statement puts zeros into the last n places of nums1.
		sets.Need(slices.Equal(a[mm:], make([]int, nn)))
		i, j := mm-1, nn-1
		for w := mm + nn - 1; j >= 0; w-- {
			if i >= 0 && a[i] > b[j] {
				a[w] = a[i]
				i--
			} else {
				a[w] = b[j]
				j--
			}
		}
		return sets.Py(a)
	}},
	// O(|t|): one pass over t moves the pointer in s on each match.
	"is-subsequence": {Arity: 2, Call: func(args []any) any {
		sub := sets.Str(args[0], lowerText)
		sets.Need(len(sub) <= 100)
		full := sets.Str(args[1], lowerText)
		sets.Need(len(full) <= 10000)
		i := 0
		for j := 0; j < len(full) && i < len(sub); j++ {
			if full[j] == sub[i] {
				i++
			}
		}
		return i == len(sub)
	}},
	// O(n + m) append operations and output characters. String-copy costs depend on the engine.
	"merge-strings-alternately": {Arity: 2, Call: func(args []any) any {
		a := sets.Str(args[0], lowerText)
		sets.Need(len(a) >= 1 && len(a) <= 100)
		b := sets.Str(args[1], lowerText)
		sets.Need(len(b) >= 1 && len(b) <= 100)
		var out strings.Builder
		for i := range max(len(a), len(b)) {
			if i < len(a) {
				out.WriteByte(a[i])
			}
			if i < len(b) {
				out.WriteByte(b[i])
			}
		}
		return out.String()
	}},
	// Expected O(n + m) time: a Hash Table of nums2 filters the unique values of nums1.
	"intersection-of-two-arrays": {Arity: 2, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1000, 0, 1000)
		b := map[int]bool{}
		for _, x := range sets.Ints(args[1], 1, 1000, 0, 1000) {
			b[x] = true
		}
		var out []int
		for _, x := range a {
			if b[x] {
				out = append(out, x)
				delete(b, x)
			}
		}
		return sets.Py(out)
	}},
	// O(n log n) for each test: sort, then count the pairs with a sum at most r minus the pairs with a sum at most l - 1, each with two pointers.
	"cf-1538-c": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			l := input.Int(1, 1e9)
			r := input.Int(l, 1e9)
			total += n
			sets.Need(total <= 2e5)
			a := make([]int, n)
			for i := range a {
				a[i] = input.Int(1, 1e9)
			}
			slices.Sort(a)
			atMost := func(x int) int {
				c := 0
				for i, j := 0, n-1; i < j; i++ {
					for i < j && a[i]+a[j] > x {
						j--
					}
					c += j - i
				}
				return c
			}
			out.WriteString(strconv.Itoa(atMost(r) - atMost(l-1)))
			out.WriteByte('\n')
		}
		input.End()
		return out.String()
	}},
	// O(n^2): sort, fix the first value, two pointers find the other two values, and equal values are skipped.
	"3sum": {Arity: 1, Call: func(args []any) any {
		a := sets.Ints(args[0], 3, 3000, -1e5, 1e5)
		slices.Sort(a)
		var out [][]int
		for i := 0; i < len(a)-2 && a[i] <= 0; i++ {
			if i > 0 && a[i] == a[i-1] {
				continue
			}
			for l, r := i+1, len(a)-1; l < r; {
				switch s := a[i] + a[l] + a[r]; {
				case s < 0:
					l++
				case s > 0:
					r--
				default:
					out = append(out, []int{a[i], a[l], a[r]})
					l++
					for l < r && a[l] == a[l-1] {
						l++
					}
					r--
				}
			}
		}
		return sets.Py(out)
	}},
	// O(n) time and storage for the copied input and output. Each group gives its character and the digits of its length.
	"string-compression": {Arity: 1, Call: func(args []any) any {
		chars, ok := args[0].([]any)
		sets.Need(ok && len(chars) >= 1 && len(chars) <= 2000)
		a := make([]string, len(chars))
		for i, char := range chars {
			a[i] = sets.Str(char, printableChar)
		}
		var out []string
		for i, j := 0, 0; i < len(a); i = j {
			for j < len(a) && a[j] == a[i] {
				j++
			}
			out = append(out, a[i])
			if j-i > 1 {
				out = append(out, strings.Split(strconv.Itoa(j-i), "")...)
			}
		}
		return sets.Py(out)
	}},
}

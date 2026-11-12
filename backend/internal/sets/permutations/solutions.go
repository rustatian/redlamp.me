package permutations

import (
	"math"
	"math/bits"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

var (
	numberPattern = regexp.MustCompile(`^[0-9?]{1,5}$`)
	cakeRow       = regexp.MustCompile(`^[.C]+$`)
	twoLetters    = regexp.MustCompile(`^[a-z]{2}$`)
	letterCase    = regexp.MustCompile(`^[a-zA-Z0-9]{1,12}$`)
	phoneDigits   = regexp.MustCompile(`^[2-9]{1,4}$`)
)

// unique checks that the numbers are different.
func unique(list []int) []int {
	seen := map[int]bool{}
	for _, x := range list {
		seen[x] = true
	}
	sets.Need(len(seen) == len(list))
	return list
}

// nextPermutation changes a to the next permutation in lexicographic order. It returns false when a goes back to the ascending order.
func nextPermutation(a []int) bool {
	i := len(a) - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i >= 0 {
		j := len(a) - 1
		for a[j] <= a[i] {
			j--
		}
		a[i], a[j] = a[j], a[i]
	}
	slices.Reverse(a[i+1:])
	return i >= 0
}

// Sorting plus Next Permutation takes O(n log n + n * P) time for P distinct permutations. Copies store n * P output entries.
func permutations(nums []int) [][]int {
	a := slices.Clone(nums)
	slices.Sort(a)
	out := [][]int{slices.Clone(a)}
	for nextPermutation(a) {
		out = append(out, slices.Clone(a))
	}
	return out
}

// References are the reference solutions of the tasks by problem slug.
var References = map[string]sets.Reference{
	// O(t * |s|): the first '?' gives 9 digits, every other '?' gives 10.
	"cf-1821-a": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		for t := input.Int(1, 20000); t > 0; t-- {
			s := input.Word(numberPattern)
			first := 1
			switch s[0] {
			case '0':
				first = 0
			case '?':
				first = 9
			}
			out = append(out, strconv.Itoa(first*int(math.Pow10(strings.Count(s[1:], "?")))))
		}
		input.End()
		return lines(out)
	}},
	// O(t * n): choose 2 of the 10 - n allowed digits, then place them in C(4, 2) = 6 ways.
	"cf-1743-a": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		for t := input.Int(1, 200); t > 0; t-- {
			n := input.Int(1, 8)
			// The statement gives the digits in increasing order.
			for j, last := 0, -1; j < n; j++ {
				last = input.Int(last+1, 9)
			}
			out = append(out, strconv.Itoa(3*(10-n)*(9-n)))
		}
		input.End()
		return lines(out)
	}},
	// O(n^2): a row or a column with c chocolates gives C(c, 2) pairs.
	"cf-629-a": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 100)
		rows, cols := make([]int, n), make([]int, n)
		for r := range n {
			line := input.Word(cakeRow)
			sets.Need(len(line) == n)
			for c := range n {
				if line[c] == 'C' {
					rows[r]++
					cols[c]++
				}
			}
		}
		input.End()
		pairs := 0
		for _, c := range slices.Concat(rows, cols) {
			pairs += c * (c - 1) / 2
		}
		return strconv.Itoa(pairs) + "\n"
	}},
	// O(t): each first letter owns 25 words, and the second letter skips the first letter.
	"cf-1674-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		for t := input.Int(1, 650); t > 0; t-- {
			s := input.Word(twoLetters)
			a, b := int(s[0]-'a'), int(s[1]-'a')
			sets.Need(a != b)
			if b > a {
				b--
			}
			out = append(out, strconv.Itoa(a*25+b+1))
		}
		input.End()
		return lines(out)
	}},
	// O(n): each value makes a pair with every equal value before it.
	"number-of-good-pairs": {Arity: 1, Call: func(args []any) any {
		var seen [101]int
		pairs := 0
		for _, x := range sets.Ints(args[0], 1, 100, 1, 100) {
			pairs += seen[x]
			seen[x]++
		}
		return sets.Py(pairs)
	}},
	// For n characters and l letters, at most n * 2^l prefix extensions produce n * 2^l output characters. String-copy costs depend on the engine.
	"letter-case-permutation": {Arity: 1, Call: func(args []any) any {
		out := []string{""}
		for _, c := range sets.Str(args[0], letterCase) {
			var next []string
			for _, p := range out {
				if '0' <= c && c <= '9' {
					next = append(next, p+string(c))
				} else {
					next = append(next, p+strings.ToLower(string(c)), p+strings.ToUpper(string(c)))
				}
			}
			out = next
		}
		return sets.Py(out)
	}},
	// O(n * 2^n): every number extends all the subsets before it.
	"subsets": {Arity: 1, Call: func(args []any) any {
		out := [][]int{{}}
		for _, x := range unique(sets.Ints(args[0], 1, 10, -10, 10)) {
			for i := range len(out) {
				out = append(out, append(slices.Clone(out[i]), x))
			}
		}
		return sets.Py(out)
	}},
	// O(n * n!).
	"permutations": {Arity: 1, Call: func(args []any) any {
		return sets.Py(permutations(unique(sets.Ints(args[0], 1, 6, -10, 10))))
	}},
	// O(k * C(n, k)): the search stops when too few numbers are left to fill the combination.
	"combinations": {Arity: 2, Call: func(args []any) any {
		top := sets.Int(args[0], 1, 20)
		size := sets.Int(args[1], 1, top)
		var out [][]int
		var cur []int
		var walk func(from int)
		walk = func(from int) {
			if len(cur) == size {
				out = append(out, slices.Clone(cur))
				return
			}
			for x := from; x <= top-size+len(cur)+1; x++ {
				cur = append(cur, x)
				walk(x + 1)
				cur = cur[:len(cur)-1]
			}
		}
		walk(1)
		return sets.Py(out)
	}},
	// O(4^n * n).
	"letter-combinations-of-a-phone-number": {Arity: 1, Call: func(args []any) any {
		keys := [...]string{"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"}
		out := []string{""}
		for _, d := range sets.Str(args[0], phoneDigits) {
			var next []string
			for _, p := range out {
				for _, c := range keys[d-'0'] {
					next = append(next, p+string(c))
				}
			}
			out = next
		}
		return sets.Py(out)
	}},
	// O(n * 2^n): after the sort, a repeated number extends only the subsets that its previous copy made.
	"subsets-ii": {Arity: 1, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 10, -10, 10)
		slices.Sort(a)
		out := [][]int{{}}
		previous := 0
		for i, x := range a {
			from := 0
			if i > 0 && x == a[i-1] {
				from = previous
			}
			previous = len(out)
			for j := from; j < previous; j++ {
				out = append(out, append(slices.Clone(out[j]), x))
			}
		}
		return sets.Py(out)
	}},
	// O(n log n + n * P) time for Sorting and P distinct permutations.
	"permutations-ii": {Arity: 1, Call: func(args []any) any {
		return sets.Py(permutations(sets.Ints(args[0], 1, 8, -10, 10)))
	}},
	// The search goes over the sorted candidates. A branch stops when the next candidate is larger than the rest of the target.
	"combination-sum": {Arity: 2, Call: func(args []any) any {
		a := unique(sets.Ints(args[0], 1, 30, 2, 40))
		slices.Sort(a)
		goal := sets.Int(args[1], 1, 40)
		var out [][]int
		var cur []int
		var walk func(from, left int)
		walk = func(from, left int) {
			if left == 0 {
				out = append(out, slices.Clone(cur))
				return
			}
			for i := from; i < len(a) && a[i] <= left; i++ {
				cur = append(cur, a[i])
				walk(i, left-a[i])
				cur = cur[:len(cur)-1]
			}
		}
		walk(0, goal)
		// The statement guarantees fewer than 150 combinations.
		sets.Need(len(out) < 150)
		return sets.Py(out)
	}},
	// O(4^n / sqrt(n)): add '(' while pairs are left, add ')' while it closes an open one.
	"generate-parentheses": {Arity: 1, Call: func(args []any) any {
		pairs := sets.Int(args[0], 1, 8)
		var out []string
		var walk func(s string, opened, closed int)
		walk = func(s string, opened, closed int) {
			if closed == pairs {
				out = append(out, s)
			}
			if opened < pairs {
				walk(s+"(", opened+1, closed)
			}
			if closed < opened {
				walk(s+")", opened, closed+1)
			}
		}
		walk("", 0, 0)
		return sets.Py(out)
	}},
	// O(n * 360): the set of reachable angles modulo 360 after each rotation.
	"cf-1097-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		reach := [360]bool{0: true}
		for i := input.Int(1, 15); i > 0; i-- {
			a := input.Int(1, 180)
			var next [360]bool
			for r, ok := range reach {
				if ok {
					next[(r+a)%360] = true
					next[(r-a+360)%360] = true
				}
			}
			reach = next
		}
		input.End()
		if reach[0] {
			return "YES\n"
		}
		return "NO\n"
	}},
	// O(n) time and O(1) extra memory.
	"next-permutation": {Arity: 1, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 100, 0, 100)
		nextPermutation(a)
		return sets.Py(a)
	}},
	// O(min(m, n)): the answer is C(m + n - 2, min(m, n) - 1).
	"unique-paths": {Arity: 2, Call: func(args []any) any {
		rows := sets.Int(args[0], 1, 100)
		cols := sets.Int(args[1], 1, 100)
		k := min(rows, cols) - 1
		rest := rows + cols - 2 - k
		paths := 1
		// paths = C(rest + i, i) grows with i. The statement guarantees an answer of at most 2 · 10^9, so a larger value fails at once.
		for i := 1; i <= k; i++ {
			paths = paths * (rest + i) / i
			sets.Need(paths <= 2e9)
		}
		return sets.Py(paths)
	}},
	// O(n^2): the factorial number system of k - 1 gives the digit at each position.
	"permutation-sequence": {Arity: 2, Call: func(args []any) any {
		size := sets.Int(args[0], 1, 9)
		factorial := []int{1}
		for i := 1; i <= size; i++ {
			factorial = append(factorial, factorial[i-1]*i)
		}
		left := span(1, size)
		rest := sets.Int(args[1], 1, factorial[size]) - 1
		var out strings.Builder
		for i := size - 1; i >= 0; i-- {
			at := rest / factorial[i]
			out.WriteString(strconv.Itoa(left[at]))
			left = slices.Delete(left, at, at+1)
			rest %= factorial[i]
		}
		return sets.Py(out.String())
	}},
	// With bounded masks, the search takes O(n!) time and O(n) auxiliary state. Constructing Q boards adds O(n^2 * Q) time and output characters.
	"n-queens": {Arity: 1, Call: func(args []any) any {
		size := sets.Int(args[0], 1, 9)
		full := 1<<size - 1
		var out [][]string
		cols := make([]int, size)
		var walk func(row, used, left, right int)
		walk = func(row, used, left, right int) {
			if row == size {
				board := make([]string, size)
				for i, c := range cols {
					board[i] = strings.Repeat(".", c) + "Q" + strings.Repeat(".", size-c-1)
				}
				out = append(out, board)
			}
			for free := full &^ (used | left | right); free != 0; free &= free - 1 {
				bit := free & -free
				cols[row] = bits.TrailingZeros(uint(bit))
				walk(row+1, used|bit, ((left|bit)<<1)&full, (right|bit)>>1)
			}
		}
		walk(0, 0, 0, 0)
		return sets.Py(out)
	}},
}

package prefixsums

import (
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// prefix returns the Prefix Sum of a: p[0] = 0 and p[i + 1] = p[i] + a[i], so the sum of a[l..r] is p[r + 1] - p[l].
func prefix(a []int) []int {
	p := make([]int, len(a)+1)
	for i, x := range a {
		p[i+1] = p[i] + x
	}
	return p
}

// prefix2 returns the 2D Prefix Sum of a: p[0][j] = p[i][0] = 0 and p[i][j] is the sum of the rows 0..i - 1 and the columns 0..j - 1.
func prefix2(a [][]int) [][]int {
	p := make([][]int, len(a)+1)
	p[0] = make([]int, len(a[0])+1)
	for i, row := range a {
		p[i+1] = make([]int, len(row)+1)
		for j, x := range row {
			p[i+1][j+1] = x + p[i][j+1] + p[i+1][j] - p[i][j]
		}
	}
	return p
}

// rectangle returns the sum of the rows r1..r2 and the columns c1..c2 from the 2D Prefix Sum p. An empty range r2 = r1 - 1 gives 0.
func rectangle(p [][]int, r1, c1, r2, c2 int) int {
	return p[r2+1][c2+1] - p[r1][c2+1] - p[r2+1][c1] + p[r1][c1]
}

// grid checks a Python value: a Matrix of 1 to maxRows rows of the same length 1 to maxColumns, with whole numbers from lo to hi.
func grid(value any, maxRows, maxColumns, lo, hi int) [][]int {
	rows, ok := value.([]any)
	sets.Need(ok && len(rows) >= 1 && len(rows) <= maxRows)
	a := make([][]int, len(rows))
	for i, row := range rows {
		a[i] = sets.Ints(row, 1, maxColumns, lo, hi)
		sets.Need(len(a[i]) == len(a[0]))
	}
	return a
}

// blocks returns f(sum, cells) for each cell (i, j) of a: the sum and the number of cells of the rows i - k..i + k and the columns j - k..j + k inside a.
func blocks(a [][]int, k int, f func(sum, cells int) int) [][]int {
	p := prefix2(a)
	m, n := len(a), len(a[0])
	out := make([][]int, m)
	for i := range out {
		out[i] = make([]int, n)
		r1, r2 := max(0, i-k), min(m-1, i+k)
		for j := range out[i] {
			c1, c2 := max(0, j-k), min(n-1, j+k)
			out[i][j] = f(rectangle(p, r1, c1, r2, c2), (r2-r1+1)*(c2-c1+1))
		}
	}
	return out
}

var decimalDigits = regexp.MustCompile(`^[0-9]+$`)

// References are the reference solutions of the tasks by problem slug.
var References = map[string]sets.Reference{
	// O(n): the sum right of i is total - left - nums[i]. The first i where it equals the left sum is the leftmost pivot.
	"find-pivot-index": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e4, -1000, 1000)
		total := prefix(nums)[len(nums)]
		left := 0
		for i, x := range nums {
			if left == total-left-x {
				return sets.Py(i)
			}
			left += x
		}
		return sets.Py(-1)
	}},
	// O(n + q) for each test: with p[i] = a_1 + ... + a_i, the new sum is total - (p[r] - p[l - 1]) + k (r - l + 1).
	"cf-1807-d": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		sumN, sumQ := 0, 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			q := input.Int(1, 2e5)
			sumN += n
			sumQ += q
			sets.Need(sumN <= 2e5 && sumQ <= 2e5)
			a := make([]int, n)
			for i := range a {
				a[i] = input.Int(1, 1e9)
			}
			p := prefix(a)
			for range q {
				l := input.Int(1, n)
				r := input.Int(l, n)
				k := input.Int(1, 1e9)
				if (p[n]-(p[r]-p[l-1])+k*(r-l+1))%2 == 1 {
					out = append(out, "YES")
				} else {
					out = append(out, "NO")
				}
			}
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n): the answer of i is the product before i times the product after i. A pass from the left writes the prefix products, a pass from the right multiplies by the suffix products.
	"product-of-array-except-self": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 2, 1e5, -30, 30)
		// left, right and ans[i] fit in 32 bits before each multiplication, so no product leaves the int64 range before its check.
		fits := func(x int) bool { return x >= math.MinInt32 && x <= math.MaxInt32 }
		ans := make([]int, len(nums))
		left := 1
		for i, x := range nums {
			ans[i] = left
			left *= x
			sets.Need(fits(left))
		}
		right := 1
		for i := len(nums) - 1; i >= 0; i-- {
			ans[i] *= right
			right *= nums[i]
			sets.Need(fits(ans[i]) && fits(right))
		}
		return sets.Py(ans)
	}},
	// O(n log n + m): the Prefix Sums of v and of its sorted copy u. A question is one subtraction.
	"cf-433-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 1e5)
		v := make([]int, n)
		for i := range v {
			v[i] = input.Int(1, 1e9)
		}
		sums := [3][]int{nil, prefix(v), prefix(slices.Sorted(slices.Values(v)))}
		var out []string
		for m := input.Int(1, 1e5); m > 0; m-- {
			p := sums[input.Int(1, 2)]
			l := input.Int(1, n)
			r := input.Int(l, n)
			out = append(out, strconv.Itoa(p[r]-p[l-1]))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// Expected O(n): seen keeps the values so far. A value that is already in seen is a repeat.
	"contains-duplicate": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, -1e9, 1e9)
		seen := map[int]bool{}
		for _, x := range nums {
			if seen[x] {
				return sets.Py(true)
			}
			seen[x] = true
		}
		return sets.Py(false)
	}},
	// Expected O(n): cnt counts the Prefix Sums seen so far, from the empty prefix 0. A new prefix s adds cnt[s - k] Subarrays, then cnt[s] increases.
	"subarray-sum-equals-k": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 2e4, -1000, 1000)
		k := sets.Int(args[1], -1e7, 1e7)
		cnt := map[int]int{0: 1}
		s, ans := 0, 0
		for _, x := range nums {
			s += x
			ans += cnt[s-k]
			cnt[s]++
		}
		return sets.Py(ans)
	}},
	// O(n): 1 adds 1 and 0 adds -1 to the balance. Equal counts of 0 and 1 make two equal balances. first keeps the first position of each balance, shifted by n.
	"contiguous-array": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, 0, 1)
		n := len(nums)
		first := slices.Repeat([]int{-1}, 2*n+1)
		first[n] = 0
		b, best := n, 0
		for i, x := range nums {
			b += 2*x - 1
			if first[b] < 0 {
				first[b] = i + 1
			} else {
				best = max(best, i+1-first[b])
			}
		}
		return sets.Py(best)
	}},
	// Expected O(n) for each test: odd glasses with +, even glasses with -. Equal juice on [l, r] is a repeated alternating Prefix Sum.
	"cf-1915-e": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		sumN := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			sumN += n
			sets.Need(sumN <= 2e5)
			seen := map[int]bool{0: true}
			s, found := 0, false
			for i := range n {
				if a := input.Int(1, 1e9); i%2 == 0 {
					s += a
				} else {
					s -= a
				}
				found = found || seen[s]
				seen[s] = true
			}
			if found {
				out = append(out, "YES")
			} else {
				out = append(out, "NO")
			}
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n) for each test: subtract 1 from each digit, then a good Subarray has the sum 0. The answer is the number of pairs of equal Prefix Sums.
	"cf-1398-c": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		sumN := 0
		for t := input.Int(1, 1000); t > 0; t-- {
			n := input.Int(1, 1e5)
			sumN += n
			sets.Need(sumN <= 1e5)
			a := input.Word(decimalDigits)
			sets.Need(len(a) == n)
			// The Prefix Sums of a_i - 1 are from -n to 8n. cnt counts them shifted by n.
			cnt := make([]int, 9*n+1)
			s, ans := n, 0
			cnt[s] = 1
			for i := range n {
				s += int(a[i]-'0') - 1
				ans += cnt[s]
				cnt[s]++
			}
			out = append(out, strconv.Itoa(ans))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n + Y), Y = 101 years: d[birth - 1950] += 1 and d[death - 1950] -= 1. The running sum of d is the population of each year, and the strict > keeps the earliest year.
	"maximum-population-year": {Arity: 1, Call: func(args []any) any {
		logs, ok := args[0].([]any)
		sets.Need(ok && len(logs) >= 1 && len(logs) <= 100)
		var d [101]int
		for _, item := range logs {
			person := sets.Ints(item, 2, 2, 1950, 2050)
			sets.Need(person[0] < person[1])
			d[person[0]-1950]++
			d[person[1]-1950]--
		}
		best, year, s := 0, 1950, 0
		for i, x := range d {
			s += x
			if s > best {
				best, year = s, 1950+i
			}
		}
		return sets.Py(year)
	}},
	// O(n + b): a Difference Array of length n + 1 gets +seats at first - 1 and -seats at last. Its running sum gives the seats of each flight.
	"corporate-flight-bookings": {Arity: 2, Call: func(args []any) any {
		bookings, ok := args[0].([]any)
		sets.Need(ok && len(bookings) >= 1 && len(bookings) <= 2e4)
		n := sets.Int(args[1], 1, 2e4)
		d := make([]int, n+1)
		for _, item := range bookings {
			booking, ok := item.([]any)
			sets.Need(ok && len(booking) == 3)
			first := sets.Int(booking[0], 1, n)
			last := sets.Int(booking[1], first, n)
			seats := sets.Int(booking[2], 1, 1e4)
			d[first-1] += seats
			d[last] -= seats
		}
		return sets.Py(prefix(d[:n])[1:])
	}},
	// O(n + T + q), T = 2 * 10^5 temperatures: a Difference Array gives the number of recipes of each temperature. A Prefix Sum of the marks "at least k recipes" answers a question with one subtraction.
	"cf-816-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 2e5)
		k := input.Int(1, n)
		q := input.Int(1, 2e5)
		const top = 200000
		d := make([]int, top+2)
		for range n {
			l := input.Int(1, top)
			r := input.Int(l, top)
			d[l]++
			d[r+1]--
		}
		admissible := make([]int, top+1)
		recipes := 0
		for t := 1; t <= top; t++ {
			recipes += d[t]
			if recipes >= k {
				admissible[t] = 1
			}
		}
		c := prefix(admissible)
		out := make([]string, q)
		for i := range out {
			a := input.Int(1, top)
			b := input.Int(a, top)
			out[i] = strconv.Itoa(c[b+1] - c[a])
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n + m + k): a Difference Array over the operation numbers counts how many queries apply each operation. A second Difference Array over the elements adds d * cnt on [l, r] of each operation.
	"cf-295-a": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 1e5)
		m := input.Int(1, 1e5)
		k := input.Int(1, 1e5)
		a := make([]int, n)
		for i := range a {
			a[i] = input.Int(0, 1e5)
		}
		type operation struct{ l, r, d int }
		operations := make([]operation, m)
		for i := range operations {
			l := input.Int(1, n)
			r := input.Int(l, n)
			operations[i] = operation{l, r, input.Int(0, 1e5)}
		}
		cd := make([]int, m+1)
		for range k {
			x := input.Int(1, m)
			y := input.Int(x, m)
			cd[x-1]++
			cd[y]--
		}
		input.End()
		ad := make([]int, n+1)
		cnt := 0
		for i, op := range operations {
			cnt += cd[i]
			ad[op.l-1] += op.d * cnt
			ad[op.r] -= op.d * cnt
		}
		add := prefix(ad[:n])
		for i := range a {
			a[i] += add[i+1]
		}
		return line(a...)
	}},
	// O(mn): the 2D Prefix Sum gives the sum of the 3 x 3 window of each cell, clipped at the edges, in 4 lookups. The answer is the sum // the number of cells.
	"image-smoother": {Arity: 1, Call: func(args []any) any {
		img := grid(args[0], 200, 200, 0, 255)
		return sets.Py(blocks(img, 1, func(sum, cells int) int { return sum / cells }))
	}},
	// O(mn): the 2D Prefix Sum gives the sum of the rows max(0, i - k)..min(m - 1, i + k) and the same columns in 4 lookups.
	"matrix-block-sum": {Arity: 2, Call: func(args []any) any {
		mat := grid(args[0], 100, 100, 1, 100)
		k := sets.Int(args[1], 1, 100)
		return sets.Py(blocks(mat, k, func(sum, _ int) int { return sum }))
	}},
	// O(q + n^2): a 2D Difference Array of size (n + 1) x (n + 1) gets 4 marks for each query. A Prefix Sum by rows, then by columns, gives the values.
	"increment-submatrices-by-one": {Arity: 2, Call: func(args []any) any {
		n := sets.Int(args[0], 1, 500)
		queries, ok := args[1].([]any)
		sets.Need(ok && len(queries) >= 1 && len(queries) <= 1e4)
		d := make([][]int, n+1)
		for i := range d {
			d[i] = make([]int, n+1)
		}
		for _, item := range queries {
			query := sets.Ints(item, 4, 4, 0, n-1)
			r1, c1, r2, c2 := query[0], query[1], query[2], query[3]
			sets.Need(r1 <= r2 && c1 <= c2)
			d[r1][c1]++
			d[r1][c2+1]--
			d[r2+1][c1]--
			d[r2+1][c2+1]++
		}
		mat := make([][]int, n)
		for i := range mat {
			mat[i] = prefix(d[i][:n])[1:]
			if i > 0 {
				for j := range mat[i] {
					mat[i][j] += mat[i-1][j]
				}
			}
		}
		return sets.Py(mat)
	}},
	// O(C^2 + n + q) for each test, C = 1000: area[h][w] sums the areas of the rectangles h x w. The strict bounds give the rows h_s + 1..h_b - 1 and the columns w_s + 1..w_b - 1 of its 2D Prefix Sum.
	"cf-1722-e": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		const side = 1000
		var out []string
		sumN, sumQ := 0, 0
		for t := input.Int(1, 100); t > 0; t-- {
			n := input.Int(1, 1e5)
			q := input.Int(1, 1e5)
			sumN += n
			sumQ += q
			sets.Need(sumN <= 1e5 && sumQ <= 1e5)
			area := make([][]int, side+1)
			for i := range area {
				area[i] = make([]int, side+1)
			}
			for range n {
				h := input.Int(1, side)
				w := input.Int(1, side)
				area[h][w] += h * w
			}
			p := prefix2(area)
			for range q {
				hs := input.Int(1, side-1)
				ws := input.Int(1, side-1)
				hb := input.Int(hs+1, side)
				wb := input.Int(ws+1, side)
				out = append(out, strconv.Itoa(rectangle(p, hs+1, ws+1, hb-1, wb-1)))
			}
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
}

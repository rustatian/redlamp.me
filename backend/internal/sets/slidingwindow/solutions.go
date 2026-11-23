package slidingwindow

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

var (
	cells          = regexp.MustCompile(`^[WB]+$`)
	lowercase      = regexp.MustCompile(`^[a-z]+$`)
	englishLetters = regexp.MustCompile(`^[A-Za-z]+$`)
	abcLetters     = regexp.MustCompile(`^[abc]+$`)
	lightColor     = regexp.MustCompile(`^[ryg]$`)
	lightColors    = regexp.MustCompile(`^[ryg]+$`)
)

// letters checks a Python value: a string of 1 to 30 000 lowercase English letters.
func letters(value any) string {
	text := sets.Str(value, lowercase)
	sets.Need(len(text) <= 3e4)
	return text
}

// shortestWindow returns the start and the length of the first shortest window of s that has every letter of t with its repeats. The length is 0 when no window exists.
// count is 1 plus the number of the later windows of this length with a different text. Each window of the shortest length is checked once inside the while loop, so count is 1 only when the answer text is unique.
func shortestWindow(s, t string) (start, length, count int) {
	var need [128]int
	for i := range len(t) {
		need[t[i]]++
	}
	missing, l := len(t), 0
	for r := range len(s) {
		if need[s[r]] > 0 {
			missing--
		}
		need[s[r]]--
		for missing == 0 {
			switch size := r - l + 1; {
			case length == 0 || size < length:
				start, length, count = l, size, 1
			case size == length && s[l:r+1] != s[start:start+length]:
				count++
			}
			need[s[l]]++
			if need[s[l]] > 0 {
				missing++
			}
			l++
		}
	}
	return start, length, count
}

// References are the reference solutions of the tasks by problem slug.
var References = map[string]sets.Reference{
	// O(n) for each test: the count of W in the window of length k changes by at most 1 at each shift.
	"cf-1690-d": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			k := input.Int(1, n)
			s := input.Word(cells)
			total += n
			sets.Need(len(s) == n && total <= 2e5)
			white := strings.Count(s[:k], "W")
			best := white
			for r := k; r < n; r++ {
				if s[r] == 'W' {
					white++
				}
				if s[r-k] == 'W' {
					white--
				}
				best = min(best, white)
			}
			out = append(out, strconv.Itoa(best))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// Expected O(n): a set holds the last k numbers. A number that is already in the set has an equal number at a distance of at most k.
	"contains-duplicate-ii": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, -1e9, 1e9)
		k := sets.Int(args[1], 0, 1e5)
		window := map[int]bool{}
		for i, x := range nums {
			if window[x] {
				return sets.Py(true)
			}
			window[x] = true
			if len(window) > k {
				delete(window, nums[i-k])
			}
		}
		return sets.Py(false)
	}},
	// O(n + m): the letter counts of the window of length m change for two letters at each shift. A window is good when its counts equal the counts of p.
	"find-all-anagrams-in-a-string": {Arity: 2, Call: func(args []any) any {
		s, p := letters(args[0]), letters(args[1])
		m := len(p)
		var need, have [26]int
		for i := range len(p) {
			need[p[i]-'a']++
		}
		starts := []int{}
		for r := range len(s) {
			have[s[r]-'a']++
			if r >= m {
				have[s[r-m]-'a']--
			}
			if r >= m-1 && have == need {
				starts = append(starts, r-m+1)
			}
		}
		return sets.Py(starts)
	}},
	// O(n): the counts of the values and the sum of the window of length k. A window is good when it has k different values.
	"maximum-sum-of-distinct-subarrays-with-length-k": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, 1, 1e5)
		k := sets.Int(args[1], 1, len(nums))
		count := make([]int, 1e5+1)
		different, sum, best := 0, 0, 0
		for r, x := range nums {
			if count[x] == 0 {
				different++
			}
			count[x]++
			sum += x
			if r >= k {
				y := nums[r-k]
				count[y]--
				if count[y] == 0 {
					different--
				}
				sum -= y
			}
			if r >= k-1 && different == k {
				best = max(best, sum)
			}
		}
		return sets.Py(best)
	}},
	// O(n): the theorems of the awake minutes, plus the largest sum of the sleep minutes in a window of length k.
	"cf-961-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 1e5)
		k := input.Int(1, n)
		a := make([]int, n)
		for i := range a {
			a[i] = input.Int(1, 1e4)
		}
		// gain[i] is the count of theorems that the technique adds in minute i.
		base, gain := 0, make([]int, n)
		for i := range n {
			if input.Int(0, 1) == 1 {
				base += a[i]
			} else {
				gain[i] = a[i]
			}
		}
		input.End()
		window := 0
		for i := range k {
			window += gain[i]
		}
		best := window
		for r := k; r < n; r++ {
			window += gain[r] - gain[r-k]
			best = max(best, window)
		}
		return strconv.Itoa(base+best) + "\n"
	}},
	// O(n): a window of ones. A 0 moves the left border past itself.
	"max-consecutive-ones": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, 0, 1)
		l, best := 0, 0
		for r, x := range nums {
			if x == 0 {
				l = r + 1
			}
			best = max(best, r-l+1)
		}
		return sets.Py(best)
	}},
	// O(n): the window has at most k zeros. While it has more, the left border moves right. Each border makes at most n steps.
	"max-consecutive-ones-iii": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, 0, 1)
		k := sets.Int(args[1], 0, len(nums))
		zeros, l, best := 0, 0, 0
		for r, x := range nums {
			zeros += 1 - x
			for zeros > k {
				zeros -= 1 - nums[l]
				l++
			}
			best = max(best, r-l+1)
		}
		return sets.Py(best)
	}},
	// O(n): the counts of the fruit types in the window. While the window has 3 types, the left border moves right. A type with the count 0 leaves the map.
	"fruit-into-baskets": {Arity: 1, Call: func(args []any) any {
		fruits := sets.Ints(args[0], 1, 1e5, 0, 1e5-1)
		sets.Need(slices.Max(fruits) < len(fruits))
		count := map[int]int{}
		l, best := 0, 0
		for r, x := range fruits {
			count[x]++
			for len(count) > 2 {
				y := fruits[l]
				count[y]--
				if count[y] == 0 {
					delete(count, y)
				}
				l++
			}
			best = max(best, r-l+1)
		}
		return sets.Py(best)
	}},
	// O(n) for each test: the removals leave one window. The answer is n minus the length of the longest window with a sum of at most s. When the total is at least s, that window has the sum s.
	"cf-1692-e": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			s := input.Int(1, 2e5)
			total += n
			sets.Need(total <= 2e5)
			a := make([]int, n)
			sum := 0
			for i := range a {
				a[i] = input.Int(0, 1)
				sum += a[i]
			}
			if sum < s {
				out = append(out, "-1")
				continue
			}
			window, l, best := 0, 0, 0
			for r, x := range a {
				window += x
				for window > s {
					window -= a[l]
					l++
				}
				best = max(best, r-l+1)
			}
			out = append(out, strconv.Itoa(n-best))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(m + n): the window grows until it has all letters of t, then it shrinks while it still has them. The statement guarantees that the answer is unique.
	"minimum-window-substring": {Arity: 2, Call: func(args []any) any {
		s, t := sets.Str(args[0], englishLetters), sets.Str(args[1], englishLetters)
		sets.Need(len(s) <= 1e5 && len(t) <= 1e5)
		start, length, count := shortestWindow(s, t)
		sets.Need(count <= 1)
		return sets.Py(s[start : start+length])
	}},
	// O(n) for each test: run is the count of the suitable days in a row up to day r. The windows of at least k days that end at r add max(0, run - k + 1).
	"cf-1840-c": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			k := input.Int(1, n)
			q := input.Signed(-1e9, 1e9)
			total += n
			sets.Need(total <= 2e5)
			run, ways := 0, 0
			for range n {
				if input.Signed(-1e9, 1e9) <= q {
					run++
				} else {
					run = 0
				}
				ways += max(0, run-k+1)
			}
			out = append(out, strconv.Itoa(ways))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n): the window keeps the product of its numbers below k. Each end r adds the r - l + 1 windows that end at r.
	"subarray-product-less-than-k": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 3e4, 1, 1000)
		k := sets.Int(args[1], 0, 1e6)
		product, l, windows := 1, 0, 0
		for r, x := range nums {
			product *= x
			for l <= r && product >= k {
				product /= nums[l]
				l++
			}
			windows += r - l + 1
		}
		return sets.Py(windows)
	}},
	// O(n): the window shrinks while it has a, b and c. Then the starts 0..l - 1 give the good substrings that end at r: add l.
	"number-of-substrings-containing-all-three-characters": {Arity: 1, Call: func(args []any) any {
		s := sets.Str(args[0], abcLetters)
		sets.Need(len(s) >= 3 && len(s) <= 5e4)
		var count [3]int
		l, substrings := 0, 0
		for r := range len(s) {
			count[s[r]-'a']++
			for count[0] > 0 && count[1] > 0 && count[2] > 0 {
				count[s[l]-'a']--
				l++
			}
			substrings += l
		}
		return sets.Py(substrings)
	}},
	// O(n) for each test: a[r] is at the place r - l + 1 of a window that starts at l, so a good window needs l ≥ r - a[r] + 1. The left border only grows, and each end r adds r - l + 1.
	"cf-1736-c1": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 2e5); t > 0; t-- {
			n := input.Int(1, 2e5)
			total += n
			sets.Need(total <= 2e5)
			l, pairs := 0, 0
			for r := range n {
				l = max(l, r-input.Int(1, n)+1)
				pairs += r - l + 1
			}
			out = append(out, strconv.Itoa(pairs))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n): the windows with exactly k different values are the windows with at most k of them minus the windows with at most k - 1. Each count adds r - l + 1 for each end r.
	"subarrays-with-k-different-integers": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 2e4, 1, 2e4)
		sets.Need(slices.Max(nums) <= len(nums))
		k := sets.Int(args[1], 1, len(nums))
		atMost := func(m int) int {
			count := make([]int, len(nums)+1)
			different, l, windows := 0, 0, 0
			for r, x := range nums {
				if count[x] == 0 {
					different++
				}
				count[x]++
				for different > m {
					y := nums[l]
					count[y]--
					if count[y] == 0 {
						different--
					}
					l++
				}
				windows += r - l + 1
			}
			return windows
		}
		return sets.Py(atMost(k) - atMost(k-1))
	}},
	// O(n) for each test: the string is written twice. A scan from the right keeps the place of the nearest g, so each c of the first half gets its wait.
	"cf-1744-c": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			c := input.Word(lightColor)[0]
			s := input.Word(lightColors)
			total += n
			sets.Need(len(s) == n && total <= 2e5 && strings.IndexByte(s, 'g') >= 0 && strings.IndexByte(s, c) >= 0)
			double := s + s
			next, wait := 0, 0
			for i := 2*n - 1; i >= 0; i-- {
				if double[i] == 'g' {
					next = i
				} else if i < n && double[i] == c {
					wait = max(wait, next-i)
				}
			}
			out = append(out, strconv.Itoa(wait))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n): the prefix sums p of the doubled Array. The best sum that ends at j is p[j] minus the smallest p[i] for i in [j - n, j - 1]. A Monotonic Deque keeps that minimum.
	"maximum-sum-circular-subarray": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 3e4, -3e4, 3e4)
		n := len(nums)
		p := make([]int, 2*n+1)
		for i := range 2 * n {
			p[i+1] = p[i] + nums[i%n]
		}
		// dq holds the starts of the window with increasing p.
		dq := []int{0}
		best := nums[0]
		for j := 1; j <= 2*n; j++ {
			if dq[0] < j-n {
				dq = dq[1:]
			}
			best = max(best, p[j]-p[dq[0]])
			for len(dq) > 0 && p[dq[len(dq)-1]] >= p[j] {
				dq = dq[:len(dq)-1]
			}
			dq = append(dq, j)
		}
		return sets.Py(best)
	}},
	// O(n): dp[i] = nums[i] + max(dp[i - k..i - 1]). A Monotonic Deque keeps the places of the window with decreasing dp.
	"jump-game-vi": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, -1e4, 1e4)
		k := sets.Int(args[1], 1, 1e5)
		dp := make([]int, len(nums))
		dp[0] = nums[0]
		dq := []int{0}
		for i := 1; i < len(nums); i++ {
			if dq[0] < i-k {
				dq = dq[1:]
			}
			dp[i] = nums[i] + dp[dq[0]]
			for len(dq) > 0 && dp[dq[len(dq)-1]] <= dp[i] {
				dq = dq[:len(dq)-1]
			}
			dq = append(dq, i)
		}
		return sets.Py(dp[len(dp)-1])
	}},
	// O(n · m) for each test: for each row dp[j] = a[j] + 1 + min(dp[j - d - 1..j - 1]) with a Monotonic Deque. Then a window of k neighbor rows gives the smallest sum of the row costs.
	"cf-1941-e": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e3); t > 0; t-- {
			n := input.Int(1, 100)
			m := input.Int(3, 2e5)
			k := input.Int(1, n)
			d := input.Int(1, m)
			total += n * m
			sets.Need(total <= 2e5)
			cost, dp := make([]int, n), make([]int, m)
			for i := range cost {
				input.Int(0, 0)
				dp[0] = 1
				dq := []int{0}
				for j := 1; j < m; j++ {
					depth := input.Int(0, 1e6)
					sets.Need(j < m-1 || depth == 0)
					if dq[0] < j-d-1 {
						dq = dq[1:]
					}
					dp[j] = depth + 1 + dp[dq[0]]
					for len(dq) > 0 && dp[dq[len(dq)-1]] >= dp[j] {
						dq = dq[:len(dq)-1]
					}
					dq = append(dq, j)
				}
				cost[i] = dp[m-1]
			}
			sum := 0
			for i := range k {
				sum += cost[i]
			}
			best := sum
			for i := k; i < n; i++ {
				sum += cost[i] - cost[i-k]
				best = min(best, sum)
			}
			out = append(out, strconv.Itoa(best))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n): the prefix sums p and a Monotonic Deque of starts with increasing p. A start leaves the front when p[j] - p[start] ≥ k: a later end gives it a longer window.
	"shortest-subarray-with-sum-at-least-k": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, -1e5, 1e5)
		k := sets.Int(args[1], 1, 1e9)
		n := len(nums)
		p := make([]int, n+1)
		for i, x := range nums {
			p[i+1] = p[i] + x
		}
		best := n + 1
		var dq []int
		for j := range p {
			for len(dq) > 0 && p[j]-p[dq[0]] >= k {
				best = min(best, j-dq[0])
				dq = dq[1:]
			}
			for len(dq) > 0 && p[dq[len(dq)-1]] >= p[j] {
				dq = dq[:len(dq)-1]
			}
			dq = append(dq, j)
		}
		if best > n {
			return sets.Py(-1)
		}
		return sets.Py(best)
	}},
}

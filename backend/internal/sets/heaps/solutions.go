package heaps

import (
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// minHeap is a binary min-heap of numbers in an array: the children of place i are places 2i + 1 and 2i + 2.
type minHeap []int

func (h *minHeap) push(x int) {
	*h = append(*h, x)
	a := *h
	i := len(a) - 1
	for parent := (i - 1) >> 1; i > 0 && a[parent] > x; i, parent = parent, (parent-1)>>1 {
		a[i] = a[parent]
	}
	a[i] = x
}

func (h *minHeap) pop() int {
	a := *h
	top, last := a[0], a[len(a)-1]
	a = a[:len(a)-1]
	*h = a
	if len(a) == 0 {
		return top
	}
	i := 0
	for child := 1; child < len(a); i, child = child, 2*child+1 {
		if child+1 < len(a) && a[child+1] < a[child] {
			child++
		}
		if a[child] >= last {
			break
		}
		a[i] = a[child]
	}
	a[i] = last
	return top
}

// For each right border r, return the smallest l with max - min <= limit on a[l..r]. Two index Deques take O(n) time and O(n) storage for retained prefixes and all left borders.
func lefts(a []int, limit int) []int {
	var hi, lo, out []int
	l, h, m := 0, 0, 0
	for r := range a {
		for len(hi) > h && a[hi[len(hi)-1]] <= a[r] {
			hi = hi[:len(hi)-1]
		}
		for len(lo) > m && a[lo[len(lo)-1]] >= a[r] {
			lo = lo[:len(lo)-1]
		}
		hi = append(hi, r)
		lo = append(lo, r)
		for a[hi[h]]-a[lo[m]] > limit {
			l++
			if hi[h] < l {
				h++
			}
			if lo[m] < l {
				m++
			}
		}
		out = append(out, l)
	}
	return out
}

// join writes the numbers with one space between them.
func join(a []int) string {
	var text strings.Builder
	for i, x := range a {
		if i > 0 {
			text.WriteByte(' ')
		}
		text.WriteString(strconv.Itoa(x))
	}
	return text.String()
}

var (
	lowercase = regexp.MustCompile(`^[a-z]{1,10}$`)
)

// References are the reference solutions of the tasks by problem slug.
var References = map[string]sets.Reference{
	// O(n): before person k ends, a person in front of k or k buys min(tickets[i], tickets[k]) tickets, and a person after k buys min(tickets[i], tickets[k] - 1).
	"time-needed-to-buy-tickets": {Arity: 2, Call: func(args []any) any {
		t := sets.Ints(args[0], 1, 100, 1, 100)
		at := sets.Int(args[1], 0, len(t)-1)
		time := 0
		for i, x := range t {
			if i <= at {
				time += min(x, t[at])
			} else {
				time += min(x, t[at]-1)
			}
		}
		return sets.Py(time)
	}},
	// O(n): the queue order does not matter, only the count of each kind. The line stops at the first sandwich that no student wants.
	"number-of-students-unable-to-eat-lunch": {Arity: 2, Call: func(args []any) any {
		s := sets.Ints(args[0], 1, 100, 0, 1)
		w := sets.Ints(args[1], len(s), len(s), 0, 1)
		want := [2]int{}
		for _, x := range s {
			want[x]++
		}
		for i, x := range w {
			if want[x] == 0 {
				return sets.Py(len(w) - i)
			}
			want[x]--
		}
		return sets.Py(0)
	}},
	// O(n) for each test: a number smaller than the first number of the deque goes to the front, any other number goes to the back.
	"cf-1579-e1": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1000); t > 0; t-- {
			n := input.Int(1, 2e5)
			seen := make([]bool, n+1)
			var front, back []int
			total += n
			sets.Need(total <= 2e5)
			first := math.MaxInt
			for range n {
				x := input.Int(1, n)
				sets.Need(!seen[x])
				seen[x] = true
				if x < first {
					front = append(front, x)
					first = x
				} else {
					back = append(back, x)
				}
			}
			slices.Reverse(front)
			out = append(out, join(append(front, back...)))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n): the Josephus recurrence. The winner place among m friends is (the winner place among m - 1 friends + k) mod m.
	"find-the-winner-of-the-circular-game": {Arity: 2, Call: func(args []any) any {
		m := sets.Int(args[0], 1, 500)
		step := sets.Int(args[1], 1, m)
		winner := 0
		for size := 2; size <= m; size++ {
			winner = (winner + step) % size
		}
		return sets.Py(winner + 1)
	}},
	// O(n log n): sort, then repeat the reveal on a queue of the places. The next smallest card goes to the place that opens next.
	"reveal-cards-in-increasing-order": {Arity: 1, Call: func(args []any) any {
		d := sets.Ints(args[0], 1, 1000, 1, 1e6)
		places := from(len(d), func(i int) int { return i })
		out := make([]int, len(d))
		unique := map[int]bool{}
		for _, x := range d {
			unique[x] = true
		}
		sets.Need(len(unique) == len(d))
		slices.Sort(d)
		head := 0
		for _, x := range d {
			out[places[head]] = x
			head++
			if head < len(places) {
				places = append(places, places[head])
				head++
			}
		}
		return sets.Py(out)
	}},
	// O(n^2 log n) for n <= 30: smash the two heaviest stones and put the difference back in sorted order.
	"last-stone-weight": {Arity: 1, Call: func(args []any) any {
		s := sets.Ints(args[0], 1, 30, 1, 1000)
		slices.Sort(s)
		for len(s) > 1 {
			y, x := s[len(s)-1], s[len(s)-2]
			s = s[:len(s)-2]
			if y > x {
				s = append(s, y-x)
				slices.Sort(s)
			}
		}
		if len(s) == 0 {
			return sets.Py(0)
		}
		return sets.Py(s[0])
	}},
	// O(k n): k times multiply the first minimum.
	"final-array-state-after-k-multiplication-operations-i": {Arity: 3, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 100, 1, 100)
		times := sets.Int(args[1], 1, 10)
		m := sets.Int(args[2], 1, 5)
		for range times {
			a[slices.Index(a, slices.Min(a))] *= m
		}
		return sets.Py(a)
	}},
	// O(n log n) for each test. A min-heap stores negated bonuses. Each zero adds the largest unused earlier bonus to the total.
	"cf-1800-c2": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(1, 2e5)
			var bonuses minHeap
			total += n
			sets.Need(total <= 2e5)
			army := 0
			for range n {
				if s := input.Int(0, 1e9); s != 0 {
					bonuses.push(-s)
				} else if len(bonuses) > 0 {
					army -= bonuses.pop()
				}
			}
			out = append(out, strconv.Itoa(army))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n log n): sort in descending order. The k-th largest counts the repeats.
	"kth-largest-element-in-an-array": {Arity: 2, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1e5, -1e4, 1e4)
		at := sets.Int(args[1], 1, len(a))
		slices.SortFunc(a, func(x, y int) int { return y - x })
		return sets.Py(a[at-1])
	}},
	// For n words, u distinct words, and bounded word lengths: expected O(n + u log(u + 1)) time to count and sort.
	"top-k-frequent-words": {Arity: 2, Call: func(args []any) any {
		words, ok := args[0].([]any)
		sets.Need(ok && len(words) >= 1 && len(words) <= 500)
		count := map[string]int{}
		var unique []string
		for _, w := range words {
			word := sets.Str(w, lowercase)
			if count[word] == 0 {
				unique = append(unique, word)
			}
			count[word]++
		}
		top := sets.Int(args[1], 1, len(unique))
		slices.SortFunc(unique, func(a, b string) int {
			if count[a] != count[b] {
				return count[b] - count[a]
			}
			return strings.Compare(a, b)
		})
		return sets.Py(unique[:top])
	}},
	// O(n) time and O(n) physical storage. The head index skips expired entries but retains their array prefix.
	"sliding-window-maximum": {Arity: 2, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1e5, -1e4, 1e4)
		size := sets.Int(args[1], 1, len(a))
		var dq, out []int
		head := 0
		for i := range a {
			for len(dq) > head && a[dq[len(dq)-1]] <= a[i] {
				dq = dq[:len(dq)-1]
			}
			dq = append(dq, i)
			if dq[head] <= i-size {
				head++
			}
			if i >= size-1 {
				out = append(out, a[dq[head]])
			}
		}
		return sets.Py(out)
	}},
	// O(n): the longest window r - l + 1 over the left borders of lefts.
	"longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit": {Arity: 2, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1e5, 1, 1e9)
		best := 0
		for r, l := range lefts(a, sets.Int(args[1], 0, 1e9)) {
			best = max(best, r-l+1)
		}
		return sets.Py(best)
	}},
	// O(n): for each right border r, the subarrays l..r, ..., r..r fit, so add r - l + 1.
	"continuous-subarrays": {Arity: 1, Call: func(args []any) any {
		count := 0
		for r, l := range lefts(sets.Ints(args[0], 1, 1e5, 1, 1e9), 2) {
			count += r - l + 1
		}
		return sets.Py(count)
	}},
	// O(n) time with a maximum Deque and a running sum. Retained array prefixes use O(n) physical storage.
	"maximum-number-of-robots-within-budget": {Arity: 3, Call: func(args []any) any {
		charge := sets.Ints(args[0], 1, 5e4, 1, 1e5)
		cost := sets.Ints(args[1], len(charge), len(charge), 1, 1e5)
		limit := sets.Int(args[2], 1, 1e15)
		var dq []int
		l, head, sum, best := 0, 0, 0, 0
		for r := range charge {
			for len(dq) > head && charge[dq[len(dq)-1]] <= charge[r] {
				dq = dq[:len(dq)-1]
			}
			dq = append(dq, r)
			sum += cost[r]
			for l <= r && charge[dq[head]]+(r-l+1)*sum > limit {
				sum -= cost[l]
				l++
				if dq[head] < l {
					head++
				}
			}
			best = max(best, r-l+1)
		}
		return sets.Py(best)
	}},
	// O(n^2 log n): sort all n^2 values. The k-th smallest counts the repeats.
	"kth-smallest-element-in-a-sorted-matrix": {Arity: 2, Call: func(args []any) any {
		matrix, ok := args[0].([]any)
		sets.Need(ok && len(matrix) >= 1 && len(matrix) <= 300)
		n := len(matrix)
		m := make([][]int, n)
		for i, row := range matrix {
			m[i] = sets.Ints(row, n, n, -1e9, 1e9)
		}
		for i, row := range m {
			for j, x := range row {
				sets.Need((j == 0 || row[j-1] <= x) && (i == 0 || m[i-1][j] <= x))
			}
		}
		flat := slices.Concat(m...)
		slices.Sort(flat)
		return sets.Py(flat[sets.Int(args[1], 1, n*n)-1])
	}},
	// O(n log n): go over the pairs by nums2 from large to small. A min-heap keeps the k largest nums1 so far and their sum. sum * nums2[i] is at most the real score of the k numbers in the heap, and at the last pair of the best choice it equals the best score.
	"maximum-subsequence-score": {Arity: 3, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1e5, 0, 1e5)
		b := sets.Ints(args[1], len(a), len(a), 0, 1e5)
		size := sets.Int(args[2], 1, len(a))
		var chosen minHeap
		sum, best := 0, 0
		order := from(len(a), func(i int) int { return i })
		slices.SortStableFunc(order, func(x, y int) int { return b[y] - b[x] })
		for _, i := range order {
			chosen.push(a[i])
			sum += a[i]
			if len(chosen) > size {
				sum -= chosen.pop()
			}
			if len(chosen) == size {
				best = max(best, sum*b[i])
			}
		}
		return sets.Py(best)
	}},
	// O(n log n): drink every potion and keep the drunk potions in a min-heap. When the health is below 0, drop the most harmful drunk potion.
	"cf-1526-c2": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 2e5)
		var drunk minHeap
		health := 0
		for range n {
			a := input.Signed(-1e9, 1e9)
			drunk.push(a)
			health += a
			if health < 0 {
				health -= drunk.pop()
			}
		}
		input.End()
		return strconv.Itoa(len(drunk)) + "\n"
	}},
	// O(n log n): sort the task numbers by the enqueue time. A pointer puts the arrived tasks in a heap of (processing time, number). An empty heap moves the time to the next arrival.
	"single-threaded-cpu": {Arity: 1, Call: func(args []any) any {
		tasks, ok := args[0].([]any)
		sets.Need(ok && len(tasks) >= 1 && len(tasks) <= 1e5)
		t := make([][]int, len(tasks))
		for i, task := range tasks {
			t[i] = sets.Ints(task, 2, 2, 1, 1e9)
		}
		var ready minHeap
		var order []int
		// The key enqueue time * 10^5 + number orders arrivals, then breaks ties by job number.
		byArrival := from(len(t), func(i int) int { return t[i][0]*1e5 + i })
		slices.Sort(byArrival)
		for i := range byArrival {
			byArrival[i] %= 1e5
		}
		time, next := 0, 0
		for len(order) < len(t) {
			if len(ready) == 0 {
				time = max(time, t[byArrival[next]][0])
			}
			// The key is processing time * 10^5 + number: a shorter task wins, then a smaller number.
			for ; next < len(t) && t[byArrival[next]][0] <= time; next++ {
				ready.push(t[byArrival[next]][1]*1e5 + byArrival[next])
			}
			i := ready.pop() % 1e5
			order = append(order, i)
			time += t[i][1]
		}
		return sets.Py(order)
	}},
	// For c candidates and k hires, h = min(n, 2c): O(n + (h + k) log(h + 1)) time and O(h) Heap storage. Initialization uses h separate pushes.
	// Keys order by cost, then worker number. The pointers refill the group of the hired worker while l <= r.
	"total-cost-to-hire-k-workers": {Arity: 3, Call: func(args []any) any {
		c := sets.Ints(args[0], 1, 1e5, 1, 1e5)
		hires := sets.Int(args[1], 1, len(c))
		m := sets.Int(args[2], 1, len(c))
		var workers minHeap
		// The key is cost * 10^5 + number.
		add := func(i int) { workers.push(c[i]*1e5 + i) }
		l, r, total := 0, len(c)-1, 0
		for ; l < m && l <= r; l++ {
			add(l)
		}
		for ; len(c)-1-r < m && l <= r; r-- {
			add(r)
		}
		for range hires {
			i := workers.pop() % 1e5
			total += c[i]
			if l <= r {
				if i < l {
					add(l)
					l++
				} else {
					add(r)
					r--
				}
			}
		}
		return sets.Py(total)
	}},
}

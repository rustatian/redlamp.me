package binarysearch

import (
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// References are the reference solutions of the tasks by problem slug.
var References = map[string]sets.Reference{
	// O(log n): Binary Search keeps the index of target, if it is present, between lo and hi.
	"binary-search": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e4, -9999, 9999)
		for i := 1; i < len(nums); i++ {
			sets.Need(nums[i-1] < nums[i])
		}
		if i, found := slices.BinarySearch(nums, sets.Int(args[1], -9999, 9999)); found {
			return sets.Py(i)
		}
		return sets.Py(-1)
	}},
	// O(log num): Binary Search on x in [1, 46340] for x * x == num.
	"valid-perfect-square": {Arity: 1, Call: func(args []any) any {
		num := sets.Int(args[0], 1, math.MaxInt32)
		x := sort.Search(46341, func(x int) bool { return x*x >= num })
		return sets.Py(x*x == num)
	}},
	// O(n + m): Two Pointers move the pointer of the smaller value; the first equal pair is the answer.
	"minimum-common-value": {Arity: 2, Call: func(args []any) any {
		a := sets.Ints(args[0], 1, 1e5, 1, 1e9)
		b := sets.Ints(args[1], 1, 1e5, 1, 1e9)
		sets.Need(slices.IsSorted(a) && slices.IsSorted(b))
		for i, j := 0, 0; i < len(a) && j < len(b); {
			switch {
			case a[i] == b[j]:
				return sets.Py(a[i])
			case a[i] < b[j]:
				i++
			default:
				j++
			}
		}
		return sets.Py(-1)
	}},
	// O(sqrt(n)): for each triangular T(i) < n, n - T(i) is triangular when 8(n - T(i)) + 1 is a perfect square.
	"cf-192-a": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 1e9)
		input.End()
		for i := 1; i*(i+1)/2 < n; i++ {
			// rest = j(j + 1) / 2 with j >= 1 when 8 rest + 1 = (2j + 1)^2. math.Sqrt is exact for a square below 2^53.
			d := 8*(n-i*(i+1)/2) + 1
			if r := int(math.Sqrt(float64(d))); r*r == d {
				return "YES\n"
			}
		}
		return "NO\n"
	}},
	// O(log n): Lower Bound, the first index with nums[i] >= target, or n.
	"search-insert-position": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e4, -1e4, 1e4)
		for i := 1; i < len(nums); i++ {
			sets.Need(nums[i-1] < nums[i])
		}
		i, _ := slices.BinarySearch(nums, sets.Int(args[1], -1e4, 1e4))
		return sets.Py(i)
	}},
	// O((n + q) log n): sort the prices; the answer of a day is the count of prices <= m (upper bound).
	"cf-706-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		x := from(input.Int(1, 1e5), func(int) int { return input.Int(1, 1e5) })
		slices.Sort(x)
		var out strings.Builder
		for q := input.Int(1, 1e5); q > 0; q-- {
			// The first index with a price >= m + 1 is the count of prices <= m.
			count, _ := slices.BinarySearch(x, input.Int(1, 1e9)+1)
			out.WriteString(strconv.Itoa(count) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(log n): lower bound and upper bound of target; the range is [lower, upper - 1] when it is not empty.
	"find-first-and-last-position-of-element-in-sorted-array": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 0, 1e5, -1e9, 1e9)
		sets.Need(slices.IsSorted(nums))
		target := sets.Int(args[1], -1e9, 1e9)
		lower, found := slices.BinarySearch(nums, target)
		if !found {
			return sets.Py([]int{-1, -1})
		}
		upper, _ := slices.BinarySearch(nums, target+1)
		return sets.Py([]int{lower, upper - 1})
	}},
	// O((n + m) log m): sort potions; a pair is successful when potion >= ceil(success / spell).
	"successful-pairs-of-spells-and-potions": {Arity: 3, Call: func(args []any) any {
		spells := sets.Ints(args[0], 1, 1e5, 1, 1e5)
		potions := sets.Ints(args[1], 1, 1e5, 1, 1e5)
		success := sets.Int(args[2], 1, 1e10)
		slices.Sort(potions)
		return sets.Py(from(len(spells), func(i int) int {
			weak, _ := slices.BinarySearch(potions, (success+spells[i]-1)/spells[i])
			return len(potions) - weak
		}))
	}},
	// O(n + m log n): prefix sums of the piles; the pile of label q is the first prefix sum >= q.
	"cf-474-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		sum := 0
		p := from(input.Int(1, 1e5), func(int) int {
			sum += input.Int(1, 1e3)
			return sum
		})
		sets.Need(sum <= 1e6)
		var out strings.Builder
		for m := input.Int(1, 1e5); m > 0; m-- {
			pile, _ := slices.BinarySearch(p, input.Int(1, sum))
			out.WriteString(strconv.Itoa(pile+1) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(log n): Binary Search on k in [0, 65536] for the first k with k(k + 1) / 2 > n; the answer is the k before it.
	"arranging-coins": {Arity: 1, Call: func(args []any) any {
		n := sets.Int(args[0], 1, math.MaxInt32)
		return sets.Py(sort.Search(65537, func(k int) bool { return k*(k+1)/2 > n }) - 1)
	}},
	// O(n log(max(a) + x)) for each test: the largest h with the water sum(max(0, h - a_i)) <= x.
	"cf-1873-e": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n, x := input.Int(1, 2e5), input.Int(1, 1e9)
			total += n
			sets.Need(total <= 2e5)
			a := from(n, func(int) int { return input.Int(1, 1e9) })
			// The water of h = max(a) + x + 1 is more than x, so the first h with too much water is at most this value.
			h := sort.Search(slices.Max(a)+x+1, func(h int) bool {
				water := 0
				for _, v := range a {
					water += max(0, h-v)
				}
				return water > x
			}) - 1
			out.WriteString(strconv.Itoa(h) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(n log M), M = max(piles): the smallest speed k in [1, M] with sum(ceil(p / k)) <= h.
	"koko-eating-bananas": {Arity: 2, Call: func(args []any) any {
		piles := sets.Ints(args[0], 1, 1e4, 1, 1e9)
		h := sets.Int(args[1], len(piles), 1e9)
		// The speed k = i + 1. At k = M each pile takes 1 hour, and len(piles) <= h.
		return sets.Py(1 + sort.Search(slices.Max(piles)-1, func(i int) bool {
			hours := 0
			for _, p := range piles {
				hours += (p + i) / (i + 1)
			}
			return hours <= h
		}))
	}},
	// O(n log S), S = sum(weights): the smallest capacity in [max(weights), S] for which the greedy loading in order needs at most days days.
	"capacity-to-ship-packages-within-d-days": {Arity: 2, Call: func(args []any) any {
		weights := sets.Ints(args[0], 1, 5e4, 1, 500)
		days := sets.Int(args[1], 1, len(weights))
		low, sum := slices.Max(weights), 0
		for _, w := range weights {
			sum += w
		}
		// The capacity is low + i. At the capacity sum all packages go in 1 day.
		return sets.Py(low + sort.Search(sum-low, func(i int) bool {
			need, load := 1, 0
			for _, w := range weights {
				if load+w > low+i {
					need++
					load = 0
				}
				load += w
			}
			return need <= days
		}))
	}},
	// O(n log n + n log k): sort; the median m is reachable when sum(max(0, m - a[i])) over the upper half is <= k.
	"cf-1201-c": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 2e5)
		sets.Need(n%2 == 1)
		k := input.Int(1, 1e9)
		a := from(n, func(int) int { return input.Int(1, 1e9) })
		input.End()
		slices.Sort(a)
		upper := a[n/2:]
		// The median is upper[0] + i. It costs 0 for i = 0 and more than k for i = k + 1.
		i := sort.Search(k+1, func(i int) bool {
			cost := 0
			for _, v := range upper {
				cost += max(0, upper[0]+i-v)
			}
			return cost > k
		}) - 1
		return strconv.Itoa(upper[0]+i) + "\n"
	}},
	// O(m + n): walk from the bottom-left corner; a negative cell makes the rest of its row negative.
	"count-negative-numbers-in-a-sorted-matrix": {Arity: 1, Call: func(args []any) any {
		rows, ok := args[0].([]any)
		sets.Need(ok && len(rows) >= 1 && len(rows) <= 100)
		grid := make([][]int, len(rows))
		for i, row := range rows {
			grid[i] = sets.Ints(row, 1, 100, -100, 100)
			sets.Need(len(grid[i]) == len(grid[0]))
			for j, x := range grid[i] {
				sets.Need((j == 0 || grid[i][j-1] >= x) && (i == 0 || grid[i-1][j] >= x))
			}
		}
		n, count := len(grid[0]), 0
		for i, j := len(grid)-1, 0; i >= 0 && j < n; {
			if grid[i][j] < 0 {
				count += n - j
				i--
			} else {
				j++
			}
		}
		return sets.Py(count)
	}},
	// O(log n): the first index with arr[i] > arr[i + 1] is the peak.
	"peak-index-in-a-mountain-array": {Arity: 1, Call: func(args []any) any {
		arr := sets.Ints(args[0], 3, 1e5, 0, 1e6)
		peak := sort.Search(len(arr)-1, func(i int) bool { return arr[i] > arr[i+1] })
		// A mountain rises strictly up to the peak and falls strictly after it, and the peak is not at an end.
		sets.Need(peak > 0 && peak < len(arr)-1)
		for i := 1; i < len(arr); i++ {
			sets.Need(arr[i-1] != arr[i] && (arr[i-1] < arr[i]) == (i <= peak))
		}
		return sets.Py(peak)
	}},
	// O(log n) for each test: d(i) = 2 S(i) - S(n) increases; the answer is min(|d(i - 1)|, |d(i)|) at the first i with d(i) >= 0.
	"cf-2009-e": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		for t := input.Int(1, 1e4); t > 0; t-- {
			n, k := input.Int(2, 1e9), input.Int(2, 1e9)
			// S(i) = a_1 + ... + a_i. 2 S(n) is at most 3 * 10^18, below the int64 limit.
			s := func(i int) int { return i*k + i*(i-1)/2 }
			d := func(i int) int { return 2*s(i) - s(n) }
			// d(1) = 2k - S(n) < 0 and d(n) = S(n) > 0, so the first i with d(i) >= 0 is in [2, n], and d(i - 1) < 0 <= d(i).
			i := 1 + sort.Search(n-1, func(j int) bool { return d(j+1) >= 0 })
			out.WriteString(strconv.Itoa(min(-d(i-1), d(i))) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(log n): the minimum is at the first index with nums[i] <= nums[n - 1].
	"find-minimum-in-rotated-sorted-array": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 5000, -5000, 5000)
		return sets.Py(nums[rotation(nums)])
	}},
	// O(log n): find the rotation point p as the first index with nums[i] <= nums[n - 1], then search target in the sorted part [p, n - 1] or [0, p - 1].
	"search-in-rotated-sorted-array": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 5000, -1e4, 1e4)
		target := sets.Int(args[1], -1e4, 1e4)
		p := rotation(nums)
		part, start := nums[:p], 0
		if target <= nums[len(nums)-1] {
			part, start = nums[p:], p
		}
		if i, found := slices.BinarySearch(part, target); found {
			return sets.Py(start + i)
		}
		return sets.Py(-1)
	}},
}

// rotation checks that nums is an increasing Array of different numbers after a rotation, and returns the index of its minimum: the first index with nums[i] <= nums[n - 1].
func rotation(nums []int) int {
	n, descents := len(nums), 0
	for i := range n - 1 {
		sets.Need(nums[i] != nums[i+1])
		if nums[i] > nums[i+1] {
			descents++
		}
	}
	// A rotated increasing Array falls at one place at most. Then its last number is below its first number.
	sets.Need(descents == 0 || descents == 1 && nums[n-1] < nums[0])
	return sort.Search(n-1, func(i int) bool { return nums[i] <= nums[n-1] })
}

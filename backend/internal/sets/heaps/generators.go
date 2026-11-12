package heaps

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The hidden tests of the set "heaps-and-queues". Each generator gives the edge cases of the statement, then random medium inputs, then the inputs at the maximum sizes.
// A call input is a Python literal list of the arguments. JSON of integers and lowercase words is a Python literal too.

// from returns n values f(0), ..., f(n - 1), as Array.from({ length: n }, (_, i) => f(i)) in JavaScript.
func from[T any](n int, f func(i int) T) []T {
	values := make([]T, n)
	for i := range values {
		values[i] = f(i)
	}
	return values
}

// A random integer in [lo, hi].
func between(random func() float64, lo, hi int) int {
	return lo + int(math.Floor(random()*float64(hi-lo+1)))
}

func list(random func() float64, n, lo, hi int) []int {
	return from(n, func(int) int { return between(random, lo, hi) })
}

// The numbers 1..n in a random order (Fisher-Yates).
func permutation(random func() float64, n int) []int {
	p := from(n, func(i int) int { return i + 1 })
	for i := n - 1; i > 0; i-- {
		j := between(random, 0, i)
		p[i], p[j] = p[j], p[i]
	}
	return p
}

// n different random integers in [lo, hi], in the order of their first draw.
func distinct(random func() float64, n, lo, hi int) []int {
	seen := map[int]bool{}
	var values []int
	for len(values) < n {
		if x := between(random, lo, hi); !seen[x] {
			seen[x] = true
			values = append(values, x)
		}
	}
	return values
}

func call(cases [][]any) []string {
	return from(len(cases), func(i int) string { return sets.Call(cases[i]...) })
}

// block writes a Codeforces array to stdin: the length is on one line, and the numbers are on the next line. A multi-test input starts with the count of the tests.
func block(a []int) string { return strconv.Itoa(len(a)) + "\n" + join(a) + "\n" }

func tests(blocks []string) string {
	return strconv.Itoa(len(blocks)) + "\n" + strings.Join(blocks, "")
}

// A random word of 1..10 lowercase letters.
func word(random func() float64) string {
	return string(from(between(random, 1, 10), func(int) byte { return byte(between(random, 97, 122)) }))
}

// n words from a vocabulary of `size` words, so the words repeat. The second argument k is from 1 to the number of unique words.
func words(random func() float64, n, size int) []any {
	vocabulary := from(size, func(int) string { return word(random) })
	picked := from(n, func(int) string { return vocabulary[between(random, 0, size-1)] })
	unique := map[string]bool{}
	for _, w := range picked {
		unique[w] = true
	}
	return []any{picked, between(random, 1, len(unique))}
}

// An n x n matrix with non-decreasing rows and columns: each value is the larger of its upper and left neighbours plus a random step from 0 to `step`. The largest value is start + (2n - 1) * step.
func sortedMatrix(random func() float64, n, start, step int) [][]int {
	m := make([][]int, n)
	for i := range m {
		m[i] = make([]int, n)
		for j := range m[i] {
			up, left := start, start
			if i > 0 {
				up = m[i-1][j]
			}
			if j > 0 {
				left = m[i][j-1]
			}
			m[i][j] = max(up, left) + between(random, 0, step)
		}
	}
	return m
}

// n tasks [enqueueTime, processingTime] with enqueue times in [1, last] and processing times in [1, longest].
func cpuTasks(random func() float64, n, last, longest int) [][]int {
	return from(n, func(int) []int { return []int{between(random, 1, last), between(random, 1, longest)} })
}

// Generators are the hidden test generators of the tasks by problem slug.
var Generators = map[string]sets.Generator{
	"time-needed-to-buy-tickets": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 0}, {[]int{100}, 0}, {[]int{1, 100}, 0}, {[]int{100, 1}, 1}, {[]int{2, 3, 2}, 2}, {[]int{5, 1, 1, 1}, 0}, {slices.Repeat([]int{1}, 100), 99}, {slices.Repeat([]int{100}, 100), 0}},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{list(random, n, 1, between(random, 1, 100)), between(random, 0, n-1)}
			}),
			[][]any{{list(random, 100, 1, 100), between(random, 0, 99)}, {slices.Repeat([]int{100}, 100), 99}},
		))
	},
	"number-of-students-unable-to-eat-lunch": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{0}, []int{0}}, {[]int{0}, []int{1}}, {[]int{1}, []int{1}}, {[]int{1}, []int{0}}, {[]int{0, 1}, []int{1, 0}}, {[]int{1, 1, 0, 0}, []int{0, 1, 0, 1}}, {[]int{1, 0, 0, 0}, []int{0, 0, 1, 1}}, {[]int{1, 1, 1, 0, 0, 1}, []int{1, 0, 0, 0, 1, 1}}},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{list(random, n, 0, 1), list(random, n, 0, 1)}
			}),
			[][]any{{list(random, 100, 0, 1), list(random, 100, 0, 1)}, {slices.Repeat([]int{1}, 100), slices.Repeat([]int{0}, 100)}},
		))
	},
	"cf-1579-e1": func(random func() float64) []string {
		// The answer of n = 2 * 10^5 has about 1.3 million characters, more than the judge route accepts. n = 150000 gives about 940 thousand.
		// In the decreasing permutation each number goes to the front, so list.insert(0, x) makes about 1.1 * 10^10 moves.
		return slices.Concat(
			[]string{"1\n1\n1\n", "1\n2\n1 2\n", "1\n2\n2 1\n", "1\n5\n1 2 3 4 5\n", "1\n5\n5 4 3 2 1\n", "3\n4\n3 1 2 4\n1\n1\n6\n4 6 2 5 1 3\n"},
			from(4, func(int) string {
				return tests(from(between(random, 1, 5), func(int) string { return block(permutation(random, between(random, 1, 30))) }))
			}),
			[]string{
				tests(from(1000, func(int) string { return block(permutation(random, 200)) })), tests([]string{block(permutation(random, 150000))}),
				tests([]string{block(from(150000, func(i int) int { return 150000 - i }))}),
			},
		)
	},
	"find-the-winner-of-the-circular-game": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{1, 1}, {2, 1}, {2, 2}, {5, 2}, {6, 5}, {500, 1}, {500, 500}},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{n, between(random, 1, n)}
			}),
			[][]any{{500, between(random, 2, 499)}, {500, 499}},
		))
	},
	"reveal-cards-in-increasing-order": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}}, {[]int{1000000}}, {[]int{1, 1000000}}, {[]int{1000000, 1}}, {[]int{1, 2, 3, 4}}, {[]int{17, 13, 11, 2, 3, 5, 7}}},
			from(4, func(int) []any {
				return []any{distinct(random, between(random, 2, 60), 1, between(random, 60, 1000000))}
			}),
			[][]any{{distinct(random, 1000, 1, 1000000)}, {permutation(random, 1000)}},
		))
	},
	"last-stone-weight": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}}, {[]int{1000}}, {[]int{1, 1}}, {[]int{1000, 1}}, {[]int{2, 2, 2}}, {[]int{5, 3, 2}}, {slices.Repeat([]int{1000}, 30)}, {from(30, func(i int) int { return i + 1 })}},
			from(4, func(int) []any { return []any{list(random, between(random, 2, 20), 1, between(random, 2, 1000))} }),
			[][]any{{list(random, 30, 1, 1000)}, {list(random, 30, 990, 1000)}},
		))
	},
	"final-array-state-after-k-multiplication-operations-i": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 1, 1}, {[]int{100}, 10, 5}, {[]int{1, 1}, 10, 2}, {[]int{7, 7, 7}, 2, 3}, {[]int{2, 1, 3}, 10, 1}, {[]int{100, 1}, 10, 5}, {slices.Repeat([]int{1}, 100), 10, 5}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 2, 60), 1, between(random, 2, 100)), between(random, 1, 10), between(random, 1, 5)}
			}),
			[][]any{{list(random, 100, 1, 100), 10, between(random, 2, 5)}, {slices.Repeat([]int{100}, 100), 10, 5}},
		))
	},
	"cf-1800-c2": func(random func() float64) []string {
		// A card is a hero with the chance of one half, else a bonus up to hi.
		deck := func(n, hi int) []int {
			return from(n, func(int) int {
				if random() < 0.5 {
					return 0
				}
				return between(random, 1, hi)
			})
		}
		return slices.Concat(
			[]string{"1\n1\n0\n", "1\n1\n1000000000\n", "1\n2\n0 5\n", "1\n2\n5 0\n", "1\n5\n0 0 0 0 0\n", "1\n6\n1000000000 1000000000 1000000000 0 0 0\n"},
			from(4, func(int) string {
				return tests(from(between(random, 1, 5), func(int) string { return block(deck(between(random, 1, 30), between(random, 1, 1e9))) }))
			}),
			[]string{
				tests(from(10000, func(int) string { return block(deck(20, 1e9)) })), tests([]string{block(deck(200000, 1e9))}),
				// All bonuses come before all heroes: a scan of the unused bonuses for each hero reads about 5 * 10^9 numbers.
				tests([]string{block(slices.Concat(list(random, 100000, 1, 1e9), make([]int, 100000)))}),
			},
		)
	},
	"kth-largest-element-in-an-array": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 1}, {[]int{-10000}, 1}, {[]int{10000, -10000}, 1}, {[]int{10000, -10000}, 2}, {[]int{5, 5, 5}, 2}, {[]int{1, 2, 3, 4, 5}, 5}},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{list(random, n, -between(random, 0, 20), between(random, 0, 20)), between(random, 1, n)}
			}),
			[][]any{{list(random, 100000, -10000, 10000), between(random, 1, 100000)}, {list(random, 100000, -3, 3), between(random, 1, 100000)}},
		))
	},
	"top-k-frequent-words": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]string{"a"}, 1}, {[]string{"abcdefghij"}, 1}, {[]string{"b", "a"}, 2}, {[]string{"a", "aa", "a", "aa", "b"}, 2}, {[]string{"z", "y", "x", "y"}, 3}, {slices.Repeat([]string{"word"}, 500), 1}},
			from(4, func(int) []any { return words(random, between(random, 2, 60), between(random, 1, 15)) }),
			[][]any{words(random, 500, 50), words(random, 500, 500)},
		))
	},
	"sliding-window-maximum": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 1}, {[]int{-10000}, 1}, {[]int{10000, -10000}, 1}, {[]int{10000, -10000}, 2}, {[]int{5, 5, 5, 5}, 2}, {[]int{1, 2, 3, 4, 5}, 3}, {[]int{5, 4, 3, 2, 1}, 3}},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{list(random, n, -between(random, 0, 20), between(random, 0, 20)), between(random, 1, n)}
			}),
			// k = n / 2 makes a scan of each window cost 2.5 * 10^9 reads. In descending order the maximum leaves each window, so a new scan only when the maximum leaves costs the same.
			[][]any{
				{slices.SortedFunc(slices.Values(list(random, 100000, -10000, 10000)), func(x, y int) int { return y - x }), 50000},
				{list(random, 100000, -10000, 10000), between(random, 1, 100000)},
			},
		))
	},
	"longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 0}, {[]int{1000000000}, 1000000000}, {[]int{1, 1000000000}, 999999999}, {[]int{1, 1000000000}, 999999998}, {[]int{7, 7, 7, 7}, 0}, {[]int{1, 2, 1, 2}, 0}, {[]int{5, 1, 5, 1}, 4}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 2, 60), 1, between(random, 2, 30)), between(random, 0, 10)}
			}),
			[][]any{{list(random, 100000, 1, 1000000000), between(random, 0, 1000000000)}, {list(random, 100000, 1, 10), 9}},
		))
	},
	"continuous-subarrays": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}}, {[]int{1000000000}}, {[]int{1, 3}}, {[]int{1, 4}}, {[]int{2, 2, 2, 2}}, {[]int{1, 1000000000, 1}}},
			from(4, func(int) []any { return []any{list(random, between(random, 2, 60), 1, between(random, 3, 8))} }),
			// Values in a range of 2 make every subarray continuous: 100000 * 100001 / 2 = 5000050000.
			[][]any{{list(random, 100000, 1, 1000000000)}, {list(random, 100000, 999999998, 1000000000)}},
		))
	},
	"maximum-number-of-robots-within-budget": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, []int{1}, 1}, {[]int{1}, []int{1}, 2}, {[]int{100000}, []int{100000}, 1000000000000000}, {[]int{3, 1}, []int{1, 1}, 4}, {[]int{1, 1, 1, 1}, []int{1, 1, 1, 1}, 17}, {[]int{10, 1, 1, 1}, []int{1, 1, 1, 1}, 13}},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{list(random, n, 1, 100000), list(random, n, 1, 100000), between(random, 1, 100000+n*n*100000)}
			}),
			[][]any{
				{list(random, 50000, 1, 100000), list(random, 50000, 1, 100000), between(random, 1e9, 1e12)},
				{slices.Repeat([]int{100000}, 50000), slices.Repeat([]int{100000}, 50000), 1000000000000000},
			},
		))
	},
	"kth-smallest-element-in-a-sorted-matrix": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[][]int{{-5}}, 1}, {[][]int{{1000000000}}, 1}, {[][]int{{-1000000000}}, 1}, {[][]int{{1, 1}, {1, 1}}, 3}, {[][]int{{1, 2}, {3, 4}}, 4}, {[][]int{{-1000000000, 1000000000}, {1000000000, 1000000000}}, 1}},
			from(4, func(int) []any {
				n := between(random, 2, 40)
				return []any{sortedMatrix(random, n, between(random, -50, 50), between(random, 0, 3)), between(random, 1, n*n)}
			}),
			// 599 steps of at most 3338898 keep the values in [-10^9, 10^9].
			[][]any{{sortedMatrix(random, 300, -1000000000, 3338898), between(random, 1, 90000)}, {sortedMatrix(random, 300, 0, 1), 45000}},
		))
	},
	"maximum-subsequence-score": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{0}, []int{0}, 1}, {[]int{100000}, []int{100000}, 1}, {[]int{1, 2}, []int{2, 1}, 2}, {[]int{0, 0, 0}, []int{1, 2, 3}, 2}, {[]int{3, 3, 3}, []int{3, 3, 3}, 3}, {[]int{9, 1, 1}, []int{1, 9, 9}, 2}},
			from(4, func(int) []any {
				n, hi := between(random, 2, 60), between(random, 1, 100000)
				return []any{list(random, n, 0, hi), list(random, n, 0, hi), between(random, 1, n)}
			}),
			// A new sum(h) on each step with k numbers in the heap reads about 2.5 * 10^9 numbers for k = n / 2. The second input has many equal nums2.
			[][]any{
				{list(random, 100000, 0, 100000), list(random, 100000, 0, 100000), 50000},
				{list(random, 100000, 0, 100000), list(random, 100000, 1, 10), 70000},
			},
		))
	},
	"cf-1526-c2": func(random func() float64) []string {
		return slices.Concat(
			[]string{"1\n0\n", "1\n-1000000000\n", "1\n1000000000\n", "2\n-1 1\n", "2\n1 -1\n", "5\n-1 -2 -3 -4 -5\n", "3\n1000000000 -1000000000 -1000000000\n"},
			from(4, func(int) string { return block(list(random, between(random, 2, 60), -20, 20)) }),
			// The pattern 2, -2, -2 drops one potion in each three, and the drunk list grows by two: a scan of the list for the worst potion reads about 4 * 10^9 numbers.
			[]string{
				block(list(random, 200000, -1000000000, 1000000000)),
				block(from(200000, func(i int) int {
					if i%3 != 0 {
						return -2
					}
					return 2
				})),
			},
		)
	},
	"single-threaded-cpu": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[][]int{{1, 1}}}, {[][]int{{1000000000, 1000000000}}}, {[][]int{{1, 5}, {1, 5}, {1, 5}}}, {[][]int{{5, 1}, {1, 10}}}, {[][]int{{1, 1}, {1000000000, 1}}}, {[][]int{{2, 3}, {1, 3}, {2, 1}}}},
			from(2, func(int) []any { return []any{cpuTasks(random, between(random, 2, 60), 30, 10)} }),
			from(2, func(int) []any { return []any{cpuTasks(random, between(random, 2, 60), 1000, 20)} }),
			// All tasks arrive at 1: a scan of the arrived tasks on each step costs about 5 * 10^9 reads. Arrivals up to 10^9 leave idle gaps: a clock that steps by 1 makes about 5 * 10^8 idle steps.
			[][]any{{cpuTasks(random, 100000, 1, 10)}, {cpuTasks(random, 100000, 1000000000, 10000)}},
		))
	},
	"total-cost-to-hire-k-workers": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{
				{[]int{1}, 1, 1}, {[]int{100000}, 1, 1}, {[]int{1, 2}, 2, 1}, {[]int{2, 1}, 1, 1}, {[]int{3, 2, 7, 7, 1, 2}, 2, 2}, {[]int{1, 1, 1, 1}, 4, 4}, {[]int{5, 4, 3, 2, 1}, 5, 1}, {[]int{1, 2}, 2, 2},
				// In the first input both windows hold a 3: the tie rule hires worker 0, and then the 1 enters the left window. The answer is 4, the sum of the k cheapest costs is 3, and a tie to the right group gives 5.
				// In the second input the windows hold [2, 2] and [2, 3]: a tie to the right group gives 4, and the answer is 3.
				{[]int{3, 1, 2, 3}, 2, 1}, {[]int{2, 2, 1, 3, 3, 2}, 2, 2},
			},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{list(random, n, 1, 20), between(random, 1, n), between(random, 1, n)}
			}),
			// A scan of both groups on each hire costs about 3.75 * 10^9 reads on each of the two inputs. The second input has 2 * candidates = n, so the groups meet at once.
			[][]any{{list(random, 100000, 1, 100000), 100000, 25000}, {list(random, 100000, 1, 10), 50000, 50000}},
		))
	},
}

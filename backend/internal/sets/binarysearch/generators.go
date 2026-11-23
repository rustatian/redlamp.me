package binarysearch

import (
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The hidden tests of the set "binary-search". Each generator gives the edge cases of the statement, then random medium inputs, then the inputs at the maximum sizes.
// A call input is a Python literal list of the arguments. JSON of integers is a Python literal too.

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

// n random integers in [lo, hi] in non-decreasing order.
func sortedValues(random func() float64, n, lo, hi int) []int {
	return slices.Sorted(slices.Values(list(random, n, lo, hi)))
}

// n different random integers in [lo, hi] in increasing order.
func sortedDistinct(random func() float64, n, lo, hi int) []int {
	seen := map[int]bool{}
	for len(seen) < n {
		seen[between(random, lo, hi)] = true
	}
	return slices.Sorted(maps.Keys(seen))
}

// absent returns a random integer of [lo, hi] that nums does not have.
func absent(random func() float64, nums []int, lo, hi int) int {
	for {
		if x := between(random, lo, hi); !slices.Contains(nums, x) {
			return x
		}
	}
}

// memberOrAny returns a member of nums or a random integer of [lo, hi], with the chance of one half each.
func memberOrAny(random func() float64, nums []int, lo, hi int) int {
	if random() < 0.5 {
		return nums[between(random, 0, len(nums)-1)]
	}
	return between(random, lo, hi)
}

func call(cases [][]any) []string {
	return from(len(cases), func(i int) string { return sets.Call(cases[i]...) })
}

// join writes the numbers with one space between them.
func join(a []int) string {
	return strings.Join(from(len(a), func(i int) string { return strconv.Itoa(a[i]) }), " ")
}

// lines writes each number on its own line.
func lines(a []int) string {
	return strings.Join(from(len(a), func(i int) string { return strconv.Itoa(a[i]) + "\n" }), "")
}

// triangular returns the triangular number T(k) = 1 + 2 + ... + k.
func triangular(k int) int { return k * (k + 1) / 2 }

// block writes a Codeforces array to stdin: the length is on one line, and the numbers are on the next line.
func block(a []int) string { return strconv.Itoa(len(a)) + "\n" + join(a) + "\n" }

// blockWith writes a Codeforces array with one more number v: the line "n v", and the numbers on the next line.
func blockWith(v int, a []int) string {
	return strconv.Itoa(len(a)) + " " + strconv.Itoa(v) + "\n" + join(a) + "\n"
}

// rotate returns a[r:] + a[:r]: the first r numbers move to the end.
func rotate(a []int, r int) []int { return slices.Concat(a[r:], a[:r]) }

// mountain returns n numbers that rise strictly up to the index peak and fall strictly after it.
// The peak is at least 5n, and each step away from it is 1 to 5 lower, so every number stays in [1, 10^6].
func mountain(random func() float64, n, peak int) []int {
	arr := make([]int, n)
	arr[peak] = between(random, 5*n, 1e6)
	for i := peak - 1; i >= 0; i-- {
		arr[i] = arr[i+1] - between(random, 1, 5)
	}
	for i := peak + 1; i < n; i++ {
		arr[i] = arr[i-1] - between(random, 1, 5)
	}
	return arr
}

// nonIncreasingMatrix returns an m x n Matrix with the first cell first. Each other cell is min(up, left) minus 0 to step, at least -100, so the rows and the columns do not increase.
func nonIncreasingMatrix(random func() float64, m, n, first, step int) [][]int {
	grid := make([][]int, m)
	for i := range grid {
		grid[i] = make([]int, n)
		for j := range grid[i] {
			switch {
			case i == 0 && j == 0:
				grid[i][j] = first
			case i == 0:
				grid[i][j] = grid[i][j-1]
			case j == 0:
				grid[i][j] = grid[i-1][j]
			default:
				grid[i][j] = min(grid[i-1][j], grid[i][j-1])
			}
			if i+j > 0 {
				grid[i][j] = max(-100, grid[i][j]-between(random, 0, step))
			}
		}
	}
	return grid
}

// Generators are the hidden test generators of the tasks by problem slug.
var Generators = map[string]sets.Generator{
	"binary-search": func(random func() float64) []string {
		first, last, missing := sortedDistinct(random, 1e4, -9999, 9999), sortedDistinct(random, 1e4, -9999, 9999), sortedDistinct(random, 1e4, -9999, 9999)
		return call(slices.Concat(
			[][]any{{[]int{5}, 5}, {[]int{5}, -9999}, {[]int{-9999}, 9999}, {[]int{1, 3}, 0}, {[]int{1, 3}, 3}, {[]int{1, 3}, 2}},
			from(5, func(int) []any {
				nums := sortedDistinct(random, between(random, 2, 2000), -9999, 9999)
				return []any{nums, memberOrAny(random, nums, -9999, 9999)}
			}),
			// 10^4 numbers leave 9999 of the 19999 values of the range for the absent target.
			[][]any{{first, first[0]}, {last, last[len(last)-1]}, {missing, absent(random, missing, -9999, 9999)}},
		))
	},
	"valid-perfect-square": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{1}, {2}, {3}, {4}, {2147395600}, {2147395599}, {2147395601}, {math.MaxInt32}},
			// A square, the number before a square, the number after a square, then any number.
			from(4, func(i int) []any {
				if i == 3 {
					return []any{between(random, 1, math.MaxInt32)}
				}
				x := between(random, 2, 46340)
				return []any{x*x + []int{0, -1, 1}[i]}
			}),
			[][]any{{46339 * 46339}, {46340*46340 + 46340}},
		))
	},
	"minimum-common-value": func(random func() float64) []string {
		// 2v + 1 in nums1 and 2v + 2 in nums2 are odd and even, so 10^9 at the end is the only common value.
		halves := func(n, add int) []int {
			values := sortedDistinct(random, n, 0, 5e8-add)
			for i := range values {
				values[i] = 2*values[i] + add
			}
			return append(values, 1e9)
		}
		odd, even := halves(99999, 1), halves(99999, 2)
		one := sortedValues(random, 1e5, 1, 1e9)
		return call(slices.Concat(
			[][]any{{[]int{1}, []int{1}}, {[]int{1}, []int{2}}, {[]int{2}, []int{1}}, {[]int{1, 1, 1}, []int{1}}, {[]int{1, 2, 3}, []int{4, 5, 6}}, {[]int{1000000000}, []int{1000000000}}},
			// A small top value gives common values, a large one usually gives -1.
			from(5, func(int) []any {
				hi := between(random, 10, 1e9)
				return []any{sortedValues(random, between(random, 1, 5000), 1, hi), sortedValues(random, between(random, 1, 5000), 1, hi)}
			}),
			// A check x in nums2 on a list reads 10^10 numbers for the first two inputs.
			[][]any{
				{from(1e5, func(i int) int { return 2*i + 1 }), from(1e5, func(i int) int { return 2*i + 2 })},
				{odd, even},
				{[]int{one[between(random, 0, len(one)-1)]}, one},
			},
		))
	},
	"cf-192-a": func(random func() float64) []string {
		return slices.Concat(
			[]string{"1\n", "2\n", "3\n", "4\n", "5\n", "6\n", "256\n", "512\n"},
			from(2, func(int) string { return strconv.Itoa(between(random, 1, 1e6)) + "\n" }),
			[]string{strconv.Itoa(triangular(between(random, 1, 1000))+triangular(between(random, 1, 1000))) + "\n"},
			// All pairs of the 44720 triangular numbers below 10^9 are about 2 * 10^9 sums. A NO answer makes a solution check every T(k) < n.
			[]string{"1000000000\n", "999950883\n", "999999999\n"},
		)
	},
	"search-insert-position": func(random func() float64) []string {
		low, high, middle := sortedDistinct(random, 1e4, -9999, 1e4), sortedDistinct(random, 1e4, -1e4, 9999), sortedDistinct(random, 1e4, -1e4, 1e4)
		return call(slices.Concat(
			[][]any{{[]int{7}, 7}, {[]int{7}, -10000}, {[]int{7}, 10000}, {[]int{-10000, 10000}, 0}, {[]int{-10000, 10000}, -10000}, {[]int{-10000, 10000}, 10000}},
			from(5, func(int) []any {
				nums := sortedDistinct(random, between(random, 2, 2000), -1e4, 1e4)
				return []any{nums, memberOrAny(random, nums, -1e4, 1e4)}
			}),
			// The target before all numbers, after all numbers, and absent near the middle of the range.
			[][]any{{low, -10000}, {high, 10000}, {middle, absent(random, middle, -1000, 1000)}},
		))
	},
	"cf-706-b": func(random func() float64) []string {
		drink := func(x, m []int) string { return block(x) + strconv.Itoa(len(m)) + "\n" + lines(m) }
		same := slices.Repeat([]int{1e5}, 1e5)
		return slices.Concat(
			[]string{"5\n3 10 8 6 11\n4\n1\n10\n3\n11\n", "1\n5\n3\n4\n5\n6\n", "4\n2 2 2 2\n3\n1\n2\n1000000000\n", "1\n100000\n1\n1\n", "1\n1\n1\n1000000000\n", "5\n100000 100000 100000 100000 100000\n2\n99999\n100000\n"},
			// Most days have m near the prices; one day in ten can spend up to 10^9 coins.
			from(4, func(int) string {
				x := list(random, between(random, 1, 3000), 1, between(random, 1, 1e5))
				return drink(x, from(between(random, 1, 3000), func(int) int {
					if random() < 0.1 {
						return between(random, 1, 1e9)
					}
					return between(random, 1, 1e5+10)
				}))
			}),
			// A Linear Scan of the prices for each day is n * q = 10^10 comparisons.
			[]string{
				drink(list(random, 1e5, 1, 1e5), list(random, 1e5, 1, 1e5)),
				drink(same, from(1e5, func(i int) int { return 1e5 - 1 + i%2 })),
				drink(list(random, 1e5, 1, 1e5), list(random, 1e5, 1, 1e9)),
			},
		)
	},
	"find-first-and-last-position-of-element-in-sorted-array": func(random func() float64) []string {
		wide := sortedValues(random, 1e5, -1e9, 1e9)
		return call(slices.Concat(
			[][]any{{[]int{}, 0}, {[]int{}, -1000000000}, {[]int{1}, 1}, {[]int{1}, 0}, {[]int{2, 2, 2, 2, 2}, 2}, {[]int{1, 2, 3}, 4}},
			// Values from a small range give long runs of equal numbers.
			from(5, func(int) []any {
				return []any{sortedValues(random, between(random, 1, 3000), -20, 20), between(random, -22, 22)}
			}),
			// One search and a walk to both sides takes 10^5 steps on the first input.
			[][]any{{slices.Repeat([]int{7}, 1e5), 7}, {wide, wide[between(random, 0, len(wide)-1)]}, {wide, absent(random, wide, -1e9, 1e9)}},
		))
	},
	"successful-pairs-of-spells-and-potions": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, []int{1}, 1}, {[]int{1}, []int{1}, 2}, {[]int{100000}, []int{100000}, 10000000000}, {[]int{99999}, []int{100000}, 10000000000}, {[]int{1, 2, 3}, []int{1, 2, 3}, 1}},
			// success is a random value or the product of a spell and a potion of the input, which tests the exact division.
			from(5, func(i int) []any {
				spells, potions := list(random, between(random, 1, 3000), 1, 1e5), list(random, between(random, 1, 3000), 1, 1e5)
				success := between(random, 1, 1e10)
				if i%2 == 1 {
					success = spells[between(random, 0, len(spells)-1)] * potions[between(random, 0, len(potions)-1)]
				}
				return []any{spells, potions, success}
			}),
			// All pairs are 10^10 products. success = 1 gives the longest answer: 10^5 counts of 100000, about 8 * 10^5 characters.
			[][]any{
				{list(random, 1e5, 1, 1e5), list(random, 1e5, 1, 1e5), between(random, 1, 1e10)},
				{list(random, 1e5, 1, 1e5), list(random, 1e5, 1, 1e5), 10000000000},
				{list(random, 1e5, 1, 1e5), list(random, 1e5, 1, 1e5), 1},
			},
		))
	},
	"cf-474-b": func(random func() float64) []string {
		// Random labels of the worms of the piles a.
		labels := func(a []int, m int) []int {
			sum := 0
			for _, x := range a {
				sum += x
			}
			return list(random, m, 1, sum)
		}
		heavy := slices.Repeat([]int{1000}, 1000)
		return slices.Concat(
			[]string{"5\n2 7 3 4 9\n3\n1 25 11\n", "1\n1\n1\n1\n", "3\n1 1 1\n3\n3 2 1\n", "1\n1000\n2\n1 1000\n", "4\n3 1 4 2\n5\n3 4 5 8 9\n"},
			// At most 1000 piles of at most 1000 worms keep the sum at most 10^6.
			from(5, func(int) string {
				a := list(random, between(random, 1, 1000), 1, between(random, 1, 1000))
				return block(a) + block(labels(a, between(random, 1, 2000)))
			}),
			// A walk over the piles for each label is up to n * m = 10^10 steps for the first and the third input.
			[]string{
				func() string { a := list(random, 1e5, 1, 10); return block(a) + block(labels(a, 1e5)) }(),
				block(heavy) + block(labels(heavy, 1e5)),
				block(slices.Repeat([]int{1}, 1e5)) + block(from(1e5, func(i int) int { return 1e5 - i })),
			},
		)
	},
	"arranging-coins": func(random func() float64) []string {
		return call(slices.Concat(
			// T(65535) = 2147450880 is the largest triangular number of the range.
			[][]any{{1}, {2}, {3}, {5}, {6}, {8}, {math.MaxInt32}, {2147450880}, {2147450879}},
			// A random n and the number before a triangular number.
			[][]any{{between(random, 1, 1e6)}, {triangular(between(random, 2, 1400)) - 1}},
			[][]any{{between(random, 2e9, math.MaxInt32)}, {triangular(65000)}},
		))
	},
	"cf-1873-e": func(random func() float64) []string {
		// tests writes t random tests of n() columns each; a test is the line "n x" and the heights.
		tests := func(t int, n func() int) string {
			return strconv.Itoa(t) + "\n" + strings.Join(from(t, func(int) string {
				return blockWith(between(random, 1, 1e9), list(random, n(), 1, between(random, 1, 1e9)))
			}), "")
		}
		return slices.Concat(
			[]string{"1\n1 1\n1\n", "1\n1 1000000000\n1\n", "1\n1 1\n1000000000\n", "1\n3 1\n2 2 2\n", "1\n2 1000000000\n1000000000 1000000000\n", "1\n" + blockWith(1, slices.Repeat([]int{1e9}, 2e5))},
			from(4, func(int) string { return tests(between(random, 1, 50), func() int { return between(random, 1, 200) }) }),
			// The answers are near 3 * 10^6, up to 6 * 10^8 and near 10^8. A Linear Scan over h = 1, 2, ... with a check of all columns makes more than 10^11 steps.
			[]string{
				"1\n" + blockWith(1e9, list(random, 2e5, 1, 1e9)),
				tests(1e4, func() int { return 20 }),
				"1\n" + blockWith(1e9, list(random, 2e5, 1e8, 1e9)),
			},
		)
	},
	"koko-eating-bananas": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 1}, {[]int{1000000000}, 1}, {[]int{1000000000}, 1000000000}, {[]int{1000000000}, 2}, {[]int{3, 6, 7, 11}, 4}, {[]int{2, 2, 2}, 1000000000}},
			from(4, func(int) []any {
				n := between(random, 1, 500)
				return []any{list(random, n, 1, between(random, 1, 1e9)), between(random, n, 1e9)}
			}),
			// h = n gives the answer max(piles), near 10^9, and h = 2n about 3 * 10^8. A Linear Scan over k = 1, 2, ... makes about 10^13 divisions.
			[][]any{{list(random, 1e4, 1, 1e9), 10000}, {list(random, 1e4, 1, 1e9), 20000}, {list(random, 1e4, 1, 1e9), between(random, 1e4, 1e9)}},
		))
	},
	"capacity-to-ship-packages-within-d-days": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{7}, 1}, {[]int{500}, 1}, {[]int{1, 2, 3, 4}, 4}, {[]int{5, 5, 5, 5}, 1}, {[]int{1, 500, 1}, 2}, {[]int{3, 2, 2, 4, 1, 4}, 3}},
			from(4, func(int) []any {
				n := between(random, 1, 1000)
				return []any{list(random, n, 1, 500), between(random, 1, n)}
			}),
			// days = 1 gives the answer sum(weights), about 1.25 * 10^7. A Linear Scan over the capacity from max(weights) makes up to 10^7 checks of 5 * 10^4 weights.
			[][]any{{list(random, 5e4, 1, 500), 1}, {list(random, 5e4, 1, 500), 50}, {list(random, 5e4, 1, 500), 50000}},
		))
	},
	"cf-1201-c": func(random func() float64) []string {
		return slices.Concat(
			[]string{"1 1\n1\n", "1 1000000000\n1000000000\n", "3 1\n5 5 5\n", "3 1000000000\n1 1 1\n", "5 3\n1 1 9 9 9\n"},
			from(5, func(int) string {
				return blockWith(between(random, 1, 1e9), list(random, 2*between(random, 0, 1500)+1, 1, between(random, 1, 1e9)))
			}),
			// The median grows by about 3 * 10^6 for the random numbers and by 10^4 for the equal numbers. Raising it by 1 at a time with a check of 10^5 numbers makes at least 10^9 steps.
			[]string{blockWith(1e9, list(random, 199999, 1, 1e9)), blockWith(1e9, slices.Repeat([]int{1}, 199999)), blockWith(1e9, slices.Repeat([]int{1e9}, 199999))},
		)
	},
	"count-negative-numbers-in-a-sorted-matrix": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[][]int{{-1}}}, {[][]int{{100}}}, {[][]int{{0}}}, {[][]int{{5, 1, 0}}}, {[][]int{{3}, {0}, {-2}, {-5}}}, {slices.Repeat([][]int{slices.Repeat([]int{-100}, 100)}, 100)}},
			from(4, func(int) []any {
				return []any{nonIncreasingMatrix(random, between(random, 1, 40), between(random, 1, 40), between(random, -100, 100), between(random, 0, 10))}
			}),
			// From 100 with steps of 0 to 1 the border of the negative cells goes through the middle; steps of 0 to 2 make most cells negative.
			[][]any{{nonIncreasingMatrix(random, 100, 100, 100, 1)}, {nonIncreasingMatrix(random, 100, 100, 100, 2)}},
		))
	},
	"peak-index-in-a-mountain-array": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{0, 1, 0}}, {[]int{0, 1000000, 0}}, {[]int{0, 1, 2, 3, 1}}, {[]int{0, 5, 4, 3, 2, 1}}, {[]int{0, 1, 2, 3, 4, 0}}},
			from(4, func(int) []any {
				n := between(random, 3, 2000)
				return []any{mountain(random, n, between(random, 1, n-2))}
			}),
			// The peak next to the first end, next to the last end, and at a random place.
			[][]any{{mountain(random, 1e5, 1)}, {mountain(random, 1e5, 1e5-2)}, {mountain(random, 1e5, between(random, 1, 1e5-2))}},
		))
	},
	"cf-2009-e": func(random func() float64) []string {
		// tests writes t tests "n k" with n from n() and k from k().
		tests := func(t int, n, k func() int) string {
			return strconv.Itoa(t) + "\n" + strings.Join(from(t, func(int) string { return strconv.Itoa(n()) + " " + strconv.Itoa(k()) + "\n" }), "")
		}
		small := func() int { return between(random, 2, 1000) }
		large := func() int { return between(random, 2, 1e9) }
		return slices.Concat(
			[]string{"6\n2 2\n2 1000000000\n1000000000 2\n1000000000 1000000000\n3 2\n4 10\n", "1\n2 2\n", "1\n1000000000 2\n", "1\n2 1000000000\n"},
			from(4, func(int) string { return tests(between(random, 1, 100), small, small) }),
			[]string{tests(1000, large, large), tests(1000, large, large)},
			// A loop over i takes up to 10^9 steps for one test with n = 10^9.
			[]string{tests(1e4, large, large), tests(1e4, func() int { return 1e9 }, large)},
		)
	},
	"find-minimum-in-rotated-sorted-array": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{7}}, {[]int{2, 1}}, {[]int{1, 2}}, {[]int{-5000, 5000}}, {[]int{5000, -5000}}},
			from(4, func(int) []any {
				n := between(random, 3, 2000)
				return []any{rotate(sortedDistinct(random, n, -5000, 5000), between(random, 0, n-1))}
			}),
			// No rotation, the minimum at the last index, and the minimum at index 1.
			[][]any{{sortedDistinct(random, 5000, -5000, 5000)}, {rotate(sortedDistinct(random, 5000, -5000, 5000), 1)}, {rotate(sortedDistinct(random, 5000, -5000, 5000), 4999)}},
		))
	},
	"search-in-rotated-sorted-array": func(random func() float64) []string {
		// A rotated Array of n different numbers, and the index of its minimum.
		rotated := func(n, r int) ([]int, int) {
			return rotate(sortedDistinct(random, n, -1e4, 1e4), r), (n - r) % n
		}
		atMin, p := rotated(5000, between(random, 0, 4999))
		beforeMin, q := rotated(5000, between(random, 1, 4999))
		missing, _ := rotated(5000, between(random, 0, 4999))
		return call(slices.Concat(
			[][]any{{[]int{1}, 0}, {[]int{1}, 1}, {[]int{1, 3}, 3}, {[]int{3, 1}, 1}, {[]int{3, 1}, 2}, {[]int{5, 1, 3}, 5}},
			from(4, func(int) []any {
				n := between(random, 3, 2000)
				nums, _ := rotated(n, between(random, 0, n-1))
				return []any{nums, memberOrAny(random, nums, -1e4, 1e4)}
			}),
			// The target is the minimum, the maximum just before the minimum, or absent.
			[][]any{{atMin, atMin[p]}, {beforeMin, beforeMin[q-1]}, {missing, absent(random, missing, -1e4, 1e4)}},
		))
	},
}

package slidingwindow

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The hidden tests of the set "sliding-window". Each generator gives the edge cases of the statement, then random medium inputs, then the inputs at the maximum sizes.
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

// n random integers in [lo, hi].
func list(random func() float64, n, lo, hi int) []int {
	return from(n, func(int) int { return between(random, lo, hi) })
}

// n random values 0 or 1. A value is 1 with the chance one.
func bits(random func() float64, n int, one float64) []int {
	return from(n, func(int) int {
		if random() < one {
			return 1
		}
		return 0
	})
}

// A random string of n letters of alphabet.
func word(random func() float64, n int, alphabet string) string {
	return string(from(n, func(int) byte { return alphabet[between(random, 0, len(alphabet)-1)] }))
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

// join writes the numbers with one space between them.
func join(a []int) string {
	return strings.Trim(fmt.Sprint(a), "[]")
}

// tests writes a multi-test input: the count of the tests on the first line, then the tests.
func tests(blocks []string) string {
	return strconv.Itoa(len(blocks)) + "\n" + strings.Join(blocks, "")
}

// stripe writes one test of Codeforces 1690D: n and k, then n cells. A cell is W with the chance white.
func stripe(random func() float64, n, k int, white float64) string {
	cells := from(n, func(int) byte {
		if random() < white {
			return 'W'
		}
		return 'B'
	})
	return fmt.Sprintf("%d %d\n%s\n", n, k, cells)
}

// lecture writes the input of Codeforces 961B: n and k, the theorems a, then the awake minutes t.
func lecture(k int, a, t []int) string {
	return fmt.Sprintf("%d %d\n%s\n%s\n", len(a), k, join(a), join(t))
}

// binaryDeque writes one test of Codeforces 1692E: n and s, then the n numbers a.
func binaryDeque(s int, a []int) string {
	return fmt.Sprintf("%d %d\n%s\n", len(a), s, join(a))
}

// skiResort writes one test of Codeforces 1840C: n, k and q, then the n temperatures a.
func skiResort(k, q int, a []int) string {
	return fmt.Sprintf("%d %d %d\n%s\n", len(a), k, q, join(a))
}

// goodSubarrays writes one test of Codeforces 1736C1: n, then the n numbers a.
func goodSubarrays(a []int) string {
	return fmt.Sprintf("%d\n%s\n", len(a), join(a))
}

// trafficLight writes one test of Codeforces 1744C: n and the current color c, then the colors s.
func trafficLight(c byte, s string) string {
	return fmt.Sprintf("%d %c\n%s\n", len(s), c, s)
}

// randomLight writes one test of Codeforces 1744C with n random colors. One of them is g, and the current color is a random color of s.
func randomLight(random func() float64, n int) string {
	s := []byte(word(random, n, "ryg"))
	s[between(random, 0, n-1)] = 'g'
	return trafficLight(s[between(random, 0, n-1)], string(s))
}

// bridges writes one test of Codeforces 1941E: n, m, k and d, then the depths of the n rows.
func bridges(k, d int, rows [][]int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d %d\n", len(rows), len(rows[0]), k, d)
	for _, row := range rows {
		b.WriteString(join(row) + "\n")
	}
	return b.String()
}

// river returns n rows of m random depths from 0 to deep. The banks in the first and the last column have the depth 0.
func river(random func() float64, n, m, deep int) [][]int {
	return from(n, func(int) []int {
		row := list(random, m, 0, deep)
		row[0], row[m-1] = 0, 0
		return row
	})
}

// countOnes returns the count of the values 1 in a list of the values 0 and 1.
func countOnes(a []int) int {
	return strings.Count(join(a), "1")
}

// Generators are the hidden test generators of the tasks by problem slug.
var Generators = map[string]sets.Generator{
	"cf-1690-d": func(random func() float64) []string {
		// All cells are W except one B.
		white := []byte(strings.Repeat("W", 2e5))
		white[between(random, 0, 2e5-1)] = 'B'
		return slices.Concat(
			[]string{"1\n1 1\nW\n", "1\n1 1\nB\n", "1\n10 5\nBBBBBBBBBB\n", "1\n8 8\nWWWWWWWW\n"},
			from(5, func(int) string {
				return tests(from(between(random, 50, 200), func(int) string {
					n := between(random, 1, 60)
					return stripe(random, n, between(random, 1, n), []float64{0.2, 0.5, 0.8}[between(random, 0, 2)])
				}))
			}),
			// A new count of each window for n = 2 · 10^5 and k = 10^5 reads (n - k + 1) · k ≈ 10^10 cells.
			[]string{
				tests(from(1e4, func(int) string { return stripe(random, 20, between(random, 1, 20), 0.5) })),
				tests([]string{stripe(random, 2e5, 1e5, 0.5)}),
				tests([]string{stripe(random, 2e5, 2e5-1, 0.5)}),
				tests([]string{fmt.Sprintf("200000 1\n%s\n", white)}),
			},
		)
	},
	"contains-duplicate-ii": func(random func() float64) []string {
		// n different values with one equal pair at the distance gap.
		pair := func(gap int) []int {
			a := distinct(random, 1e5, -1e9, 1e9)
			i := between(random, 0, len(a)-1-gap)
			a[i+gap] = a[i]
			return a
		}
		return call(slices.Concat(
			[][]any{{[]int{5}, 0}, {[]int{5}, 100000}, {[]int{7, 7}, 1}, {[]int{5, 6, 5}, 1}},
			from(4, func(i int) []any {
				hi := []int{5, 50}[i%2]
				return []any{list(random, between(random, 100, 1000), -hi, hi), between(random, 0, 20)}
			}),
			// A comparison of each pair at a distance of at most k makes about 5 · 10^9 comparisons for the first input.
			[][]any{{distinct(random, 1e5, -1e9, 1e9), 100000}, {pair(50000), 50000}, {pair(50001), 50000}, {list(random, 1e5, 1, 1000), 0}},
		))
	},
	"find-all-anagrams-in-a-string": func(random func() float64) []string {
		s := word(random, 3e4, "abcdefghijklmnopqrstuvwxyz")
		at := between(random, 0, 3e4-300)
		return call(slices.Concat(
			[][]any{{"a", "a"}, {"a", "b"}, {"ab", "abc"}, {"ba", "ab"}},
			from(4, func(i int) []any {
				alphabet := []string{"ab", "abc"}[i%2]
				return []any{word(random, between(random, 100, 2000), alphabet), word(random, between(random, 1, 10), alphabet)}
			}),
			// Sorting of each of the 15 001 windows of 15 000 letters takes seconds.
			[][]any{{word(random, 3e4, "ab"), word(random, 15000, "ab")}, {s, s[at : at+300]}, {strings.Repeat("a", 3e4), strings.Repeat("a", 15000)}},
		))
	},
	"maximum-sum-of-distinct-subarrays-with-length-k": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 1}, {[]int{2, 2}, 1}, {[]int{1, 1}, 2}, {[]int{100000, 99999, 100000}, 2}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 50, 2000), 1, 20), between(random, 1, 8)}
			}),
			// A new set for each of the 50 001 windows of 50 000 numbers makes 2.5 · 10^9 steps.
			[][]any{{list(random, 1e5, 1, 1e5), 300}, {permutation(random, 1e5), 50000}, {list(random, 1e5, 1, 1000), 2000}, {list(random, 1e5, 1, 1e5), 100000}},
		))
	},
	"cf-961-b": func(random func() float64) []string {
		medium := func(int) string {
			n := between(random, 50, 2000)
			return lecture(between(random, 1, n), list(random, n, 1, 1e4), bits(random, n, 0.5))
		}
		// t = 1, except one block of 30 000 sleep minutes.
		block := slices.Repeat([]int{1}, 1e5)
		at := between(random, 0, 1e5-3e4)
		for i := at; i < at+3e4; i++ {
			block[i] = 0
		}
		return slices.Concat(
			[]string{"1 1\n1\n1\n", "1 1\n7\n0\n", lecture(5, list(random, 5, 1, 1e4), make([]int, 5)), lecture(1, list(random, 5, 1, 1e4), slices.Repeat([]int{1}, 5))},
			from(4, medium),
			// A new sum of each of the 50 001 windows of 50 000 minutes makes 2.5 · 10^9 additions.
			[]string{
				lecture(5e4, list(random, 1e5, 1, 1e4), bits(random, 1e5, 0.5)), lecture(1, list(random, 1e5, 1, 1e4), bits(random, 1e5, 0.5)),
				lecture(1e5, list(random, 1e5, 1, 1e4), bits(random, 1e5, 0.5)), lecture(3e4, list(random, 1e5, 1, 1e4), block),
			},
		)
	},
	"max-consecutive-ones": func(random func() float64) []string {
		ones := slices.Repeat([]int{1}, 1e5)
		middle := slices.Clone(ones)
		middle[5e4] = 0
		return call(slices.Concat(
			[][]any{{[]int{0}}, {[]int{1}}, {[]int{0, 0, 0}}, {slices.Repeat([]int{1}, 10)}},
			from(4, func(i int) []any {
				return []any{bits(random, between(random, 100, 5000), []float64{0.3, 0.7, 0.95}[i%3])}
			}),
			// A count of the ones from each start of 10^5 ones makes 5 000 050 000 steps.
			[][]any{{ones}, {bits(random, 1e5, 0.5)}, {middle}, {from(1e5, func(i int) int { return 1 - i%2 })}},
		))
	},
	"max-consecutive-ones-iii": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{0}, 0}, {[]int{0}, 1}, {[]int{1}, 0}, {[]int{0, 0, 0, 0}, 0}},
			from(5, func(int) []any {
				n := between(random, 100, 5000)
				return []any{bits(random, n, random()), between(random, 0, n/10)}
			}),
			// From each start, a walk to the right while the window has at most k zeros makes about 3.75 · 10^9 steps for k = 50 000 on zeros and 5 · 10^9 steps for k = n.
			[][]any{{bits(random, 1e5, 0.5), 1000}, {make([]int, 1e5), 50000}, {bits(random, 1e5, 0.5), 0}, {bits(random, 1e5, 0.9), 100000}},
		))
	},
	"fruit-into-baskets": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{0}}, {[]int{0, 1}}, {[]int{0, 1, 2}}, {[]int{1, 1, 1}}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 100, 5000), 0, []int{2, 3, 5}[between(random, 0, 2)])}
			}),
			// From each start, a walk to the right until a third type makes 5 000 050 000 steps on 10^5 trees of two types.
			[][]any{
				{list(random, 1e5, 0, 1)}, {list(random, 1e5, 0, 2)}, {list(random, 1e5, 0, 1e5-1)},
				// Types 0 and 1 alternate in the first half, types 2 and 3 in the second half.
				{from(1e5, func(i int) int { return i/50000*2 + i%2 })},
			},
		))
	},
	"cf-1692-e": func(random func() float64) []string {
		// A test with s from 1 to the count of ones + 1. Only s above the count of ones gives -1.
		test := func(n int, one float64) string {
			a := bits(random, n, one)
			return binaryDeque(between(random, 1, countOnes(a)+1), a)
		}
		medium := func(int) string {
			return tests(from(100, func(int) string { return test(between(random, 1, 40), random()) }))
		}
		half := bits(random, 2e5, 0.5)
		// Ones only in the first and the last 1000 places: the answer removes elements at both ends.
		ends := make([]int, 2e5)
		for i := range 1000 {
			ends[i], ends[2e5-1-i] = 1, 1
		}
		return slices.Concat(
			[]string{"1\n1 1\n0\n", "1\n1 1\n1\n", "1\n3 3\n1 1 1\n", "1\n4 2\n1 1 1 1\n"},
			from(4, medium),
			// Each pair of the removal counts at the two ends makes about 2 · 10^10 pairs for n = 2 · 10^5.
			[]string{
				tests([]string{binaryDeque(countOnes(half)/2, half)}),
				tests([]string{binaryDeque(1, slices.Repeat([]int{1}, 2e5))}),
				tests(from(1e4, func(int) string { return test(20, 0.5) })),
				tests([]string{binaryDeque(1000, ends)}),
			},
		)
	},
	"minimum-window-substring": func(random func() float64) []string {
		// s has one Z among 3 to 6 lowercase letters, t has the Z and up to 4 more letters. A draw whose shortest window text is not unique is drawn again.
		medium := func(int) []any {
			for {
				alphabet := "abcdef"[:between(random, 3, 6)]
				s := []byte(word(random, between(random, 50, 2000), alphabet))
				s[between(random, 0, len(s)-1)] = 'Z'
				t := "Z" + word(random, between(random, 0, 4), alphabet)
				if _, _, count := shortestWindow(string(s), t); count == 1 {
					return []any{string(s), t}
				}
			}
		}
		// A, B and C once each in a block of 10 letters.
		abc := []byte(word(random, 1e5, "abcdefghijklmnopqrstuvwxyz"))
		at := between(random, 0, 1e5-10)
		for i, place := range permutation(random, 10)[:3] {
			abc[at+place-1] = "ABC"[i]
		}
		// Two letters z once in a String of the letters a..y.
		zz := []byte(word(random, 1e5, "abcdefghijklmnopqrstuvwxy"))
		at = between(random, 0, 1e5-2)
		zz[at], zz[at+1] = 'z', 'z'
		all := word(random, 1e5, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
		shuffled := make([]byte, 1e5)
		for i, j := range permutation(random, 1e5) {
			shuffled[i] = all[j-1]
		}
		return call(slices.Concat(
			[][]any{{"a", "b"}, {"ab", "ba"}, {"aa", "aa"}, {"ABC", "abc"}},
			from(4, medium),
			// From each start, a walk to the right until the window has all letters of t makes about 5 · 10^9 steps for the starts after the planted letters.
			[][]any{
				{string(abc), "CAB"}, {strings.Repeat("a", 5e4) + "b" + strings.Repeat("c", 5e4-1), "aab"},
				{all, string(shuffled)}, {string(zz), "zz"},
			},
		))
	},
	"cf-1840-c": func(random func() float64) []string {
		// A test of n days with the values and q in [-5, 5].
		small := func(n int) string {
			return skiResort(between(random, 1, n), between(random, -5, 5), list(random, n, -5, 5))
		}
		q := between(random, -1e9, 1e9)
		return slices.Concat(
			[]string{"1\n1 1 5\n5\n", "1\n1 1 4\n5\n", "1\n3 3 0\n0 0 0\n", "1\n3 1 -1000000000\n1000000000 -1000000000 1000000000\n"},
			from(4, func(int) string {
				return tests(from(100, func(int) string { return small(between(random, 1, 30)) }))
			}),
			// A check of each pair of borders makes about 2 · 10^10 checks for n = 2 · 10^5.
			[]string{
				tests([]string{skiResort(1, q, slices.Repeat([]int{q}, 2e5))}),
				// A day is too warm with the chance 1/2000, so the blocks of suitable days are often longer than k.
				tests([]string{skiResort(1000, 999000000, list(random, 2e5, -1e9, 1e9))}),
				tests(from(1e4, func(int) string { return small(20) })),
				tests([]string{skiResort(1e5, 1e9, list(random, 2e5, -1e9, 1e9))}),
			},
		)
	},
	"subarray-product-less-than-k": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 0}, {[]int{1}, 2}, {[]int{1000}, 1000000}, {[]int{1000, 1000}, 1000000}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 100, 3000), 1, 10), between(random, 0, 1e4)}
			}),
			// The product from each start of 3 · 10^4 ones stays below k: 450 015 000 multiplications for the first input.
			[][]any{
				{slices.Repeat([]int{1}, 3e4), 1000000},
				{from(3e4, func(int) int { return []int{1, 1, 1, 2}[between(random, 0, 3)] }), 1000000},
				{list(random, 3e4, 1, 1000), 1000000},
				{list(random, 3e4, 1, 1000), 1},
			},
		))
	},
	"number-of-substrings-containing-all-three-characters": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{"abc"}, {"cba"}, {"aab"}, {"abcc"}},
			from(4, func(int) []any { return []any{word(random, between(random, 100, 5000), "abc")} }),
			// From each start, a walk to the right until the window has a, b and c makes about 1.25 · 10^9 steps for the second and the last input.
			[][]any{
				{word(random, 5e4, "abc")},
				{strings.Repeat("a", 49998) + "bc"},
				{strings.Repeat("abc", 16666) + "ab"},
				{strings.Repeat("a", 25000) + strings.Repeat("b", 24999) + "c"},
			},
		))
	},
	"cf-1736-c1": func(random func() float64) []string {
		medium := func(int) string {
			return tests(from(100, func(int) string {
				n := between(random, 1, 40)
				return goodSubarrays(list(random, n, 1, n))
			}))
		}
		return slices.Concat(
			[]string{"1\n1\n1\n", "1\n2\n1 1\n", "1\n2\n2 2\n", "1\n3\n3 2 1\n"},
			from(4, medium),
			// From each start, a walk to the right while the window is good makes about 2 · 10^10 steps for the first input.
			[]string{
				tests([]string{goodSubarrays(slices.Repeat([]int{2e5}, 2e5))}),
				tests([]string{goodSubarrays(list(random, 2e5, 1, 2e5))}),
				tests(from(2e5, func(int) string { return "1\n1\n" })),
				tests([]string{goodSubarrays(list(random, 2e5, 1, 50))}),
			},
		)
	},
	"subarrays-with-k-different-integers": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 1}, {[]int{1, 2}, 1}, {[]int{1, 2}, 2}, {[]int{2, 2, 2}, 1}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 100, 2000), 1, 6), between(random, 1, 6)}
			}),
			// A set from each start makes n(n + 1)/2 = 200 010 000 steps for the last input and about 10^8 steps for the Permutation.
			[][]any{
				{list(random, 2e4, 1, 200), 100},
				{list(random, 2e4, 1, 3), 2},
				{permutation(random, 2e4), 10000},
				{slices.Repeat([]int{1}, 2e4), 1},
			},
		))
	},
	"cf-1744-c": func(random func() float64) []string {
		// Colors r and y, and g at about 1 % of the places.
		rare := from(2e5, func(int) byte {
			if random() < 0.01 {
				return 'g'
			}
			return "ry"[between(random, 0, 1)]
		})
		return slices.Concat(
			[]string{"1\n1 g\ng\n", "1\n2 r\ngr\n", "1\n3 y\nygy\n", "1\n5 r\ngrrrr\n"},
			from(4, func(int) string {
				return tests(from(100, func(int) string { return randomLight(random, between(random, 1, 30)) }))
			}),
			// A walk from each r to the next g makes about 2 · 10^10 steps for one g and n = 2 · 10^5.
			[]string{
				tests([]string{trafficLight('r', "g"+strings.Repeat("r", 2e5-1))}),
				tests([]string{trafficLight('y', string(rare))}),
				tests(from(1e4, func(int) string { return randomLight(random, 20) })),
				tests([]string{trafficLight('r', strings.Repeat("r", 2e5-1)+"g")}),
			},
		)
	},
	"maximum-sum-circular-subarray": func(random func() float64) []string {
		// Values 1..100 with a block of 1000 values -30 000 in the middle: the best Subarray goes through the end of the circle.
		wrap := list(random, 3e4, 1, 100)
		for i := 14500; i < 15500; i++ {
			wrap[i] = -3e4
		}
		return call(slices.Concat(
			[][]any{{[]int{1}}, {[]int{-1}}, {[]int{0, 0, 0}}, {[]int{30000, 30000, 30000}}},
			from(4, func(int) []any { return []any{list(random, between(random, 100, 3000), -100, 100)} }),
			// Each start and each length makes 9 · 10^8 additions for n = 3 · 10^4.
			[][]any{
				{list(random, 3e4, -3e4, 3e4)},
				{slices.Repeat([]int{-3e4}, 3e4)},
				{wrap},
				{from(3e4, func(i int) int { return []int{30000, -29999}[i%2] })},
			},
		))
	},
	"jump-game-vi": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{5}, 1}, {[]int{1, 2}, 1}, {[]int{-1, -2}, 1}, {[]int{1, -1, 1}, 5}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 100, 3000), -100, 100), between(random, 1, 50)}
			}),
			// The maximum of the k last dp values makes 3 749 975 000 reads for n = 10^5 and k = 5 · 10^4.
			[][]any{
				{list(random, 1e5, -1e4, 1e4), 50000},
				{list(random, 1e5, -1e4, 1e4), 100000},
				{slices.Repeat([]int{-1e4}, 1e5), 3},
				{list(random, 1e5, -1e4, 1e4), 1},
			},
		))
	},
	"cf-1941-e": func(random func() float64) []string {
		medium := func(int) string {
			return tests(from(20, func(int) string {
				n, m := between(random, 1, 5), between(random, 3, 50)
				return bridges(between(random, 1, n), between(random, 1, m), river(random, n, m, 20))
			}))
		}
		return slices.Concat(
			[]string{"1\n1 3 1 1\n0 0 0\n", "1\n1 3 1 3\n0 7 0\n", "1\n2 4 2 1\n0 1 1 0\n0 1 1 0\n", "1\n3 5 1 5\n0 9 9 9 0\n0 8 8 8 0\n0 7 7 7 0\n"},
			from(4, medium),
			// The minimum of the d + 1 last dp values makes about 1.5 · 10^10 reads for m = 2 · 10^5 and d = 10^5.
			[]string{
				tests([]string{bridges(1, 1e5, river(random, 1, 2e5, 1e6))}),
				tests([]string{bridges(50, 1000, river(random, 100, 2000, 1e6))}),
				tests([]string{bridges(1, 2, river(random, 2, 1e5, 1e6))}),
				tests(from(1000, func(int) string { return bridges(1, between(random, 1, 100), river(random, 2, 100, 1e6)) })),
			},
		)
	},
	"shortest-subarray-with-sum-at-least-k": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[]int{1}, 2}, {[]int{5}, 5}, {[]int{-5}, 1}, {[]int{1, 1, 1, 1, 1, 1000}, 1000}},
			from(4, func(int) []any {
				return []any{list(random, between(random, 100, 3000), -100, 100), between(random, 1, 500)}
			}),
			// A check of each pair of borders makes about 5 · 10^9 additions for n = 10^5.
			[][]any{
				{list(random, 1e5, -1e5, 1e5), 1000000},
				{list(random, 1e5, -1e5, 1e5), 1000000000},
				// The mean value 45 gives a total of about 4.5 · 10^6, so the answer is a long Subarray.
				{list(random, 1e5, -10, 100), 3000000},
				{slices.Repeat([]int{1e5}, 1e5), 1000000000},
			},
		))
	},
}

package modular

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The hidden tests of the set "modular-arithmetic". Each generator gives the edge cases of the statement, then random medium inputs, then the inputs at the maximum sizes.
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

// shuffle puts the values of a in a random order (Fisher-Yates).
func shuffle[T any](random func() float64, a []T) {
	for i := len(a) - 1; i > 0; i-- {
		j := between(random, 0, i)
		a[i], a[j] = a[j], a[i]
	}
}

func call(cases [][]any) []string {
	return from(len(cases), func(i int) string { return sets.Call(cases[i]...) })
}

// join writes the numbers with one space between them.
func join(a []int) string {
	parts := from(len(a), func(i int) string { return strconv.Itoa(a[i]) })
	return strings.Join(parts, " ")
}

// block writes a Codeforces array to stdin: the length is on one line, and the numbers are on the next line. A multi-test input starts with the count of the tests.
func block(a []int) string { return strconv.Itoa(len(a)) + "\n" + join(a) + "\n" }

func tests(blocks []string) string {
	return strconv.Itoa(len(blocks)) + "\n" + strings.Join(blocks, "")
}

// pairs returns n numbers in [-10^9, 10^9] in a random order. They make n / 2 pairs, and the sum of each pair divides by k.
func pairs(random func() float64, n, k int) []int {
	a := make([]int, 0, n)
	for range n / 2 {
		x := between(random, -1e9, 1e9)
		r := remainder(-x, k)
		a = append(a, x, k*between(random, -((1e9+r)/k), (1e9-r)/k)+r)
	}
	shuffle(random, a)
	return a
}

// broken changes one number of a by 1. For k >= 2 the number goes to another class of remainders, and the answer of LeetCode 1497 becomes False.
func broken(random func() float64, a []int) []int {
	a = slices.Clone(a)
	i := between(random, 0, len(a)-1)
	if a[i] == 1e9 {
		a[i]--
	} else {
		a[i]++
	}
	return a
}

// line writes the numbers of one input line.
func line(values ...int) string { return join(values) + "\n" }

// digits returns the n decimal digits of a random positive number without leading zeros.
func digits(random func() float64, n int) []int {
	d := list(random, n, 0, 9)
	d[0] = between(random, 1, 9)
	return d
}

// four returns n random items [a, b, c, m] of LeetCode 2961 with values from 1 to 1000. A value of fixed replaces the random value at its place.
func four(random func() float64, n int, fixed map[int]int) [][]int {
	return from(n, func(int) []int {
		v := list(random, 4, 1, 1000)
		for i, x := range fixed {
			v[i] = x
		}
		return v
	})
}

// divisible writes one test of Codeforces 1931D: the line n x y, then the numbers.
func divisible(x, y int, a []int) string {
	return strconv.Itoa(len(a)) + " " + strconv.Itoa(x) + " " + strconv.Itoa(y) + "\n" + join(a) + "\n"
}

// Generators are the hidden test generators of the tasks by problem slug.
var Generators = map[string]sets.Generator{
	"binary-prefix-divisible-by-5": func(random func() float64) []string {
		// The full number x without % has up to 10^5 binary digits. The statement limit lets it pass, so the tests check only the answers.
		return call(slices.Concat(
			[][]any{{[]int{0}}, {[]int{1}}, {[]int{1, 0, 1, 0}}, {slices.Repeat([]int{1}, 20)}, {slices.Repeat([]int{0}, 50)}, {from(31, func(i int) int { return 1 - i%2 })}},
			from(4, func(i int) []any { return []any{list(random, []int{100, 1000, 10000, 50000}[i], 0, 1)} }),
			[][]any{{list(random, 100000, 0, 1)}, {slices.Repeat([]int{1}, 100000)}},
		))
	},
	"cf-1490-b": func(random func() float64) []string {
		// All numbers of one class need 2n / 3 + n / 3 = n moves. One move at a time with a recount of the Array makes 30 000 · 30 000 steps for each such test.
		return slices.Concat(
			[]string{
				"4\n6\n0 2 5 5 4 8\n6\n2 0 2 1 0 0\n9\n7 1 3 4 2 10 3 9 6\n6\n0 1 2 3 4 5\n",
				tests([]string{block([]int{0, 0, 0}), block([]int{1, 1, 1}), block([]int{2, 2, 2})}),
				tests([]string{block([]int{100, 100, 99})}),
				tests([]string{block([]int{100, 100, 100})}),
				tests([]string{block([]int{0, 1, 2, 3, 4, 5, 6, 7, 8})}),
			},
			[]string{
				tests(from(10, func(int) string { return block(list(random, 3*between(random, 1, 10), 0, 100)) })),
				tests(from(100, func(int) string { return block(list(random, 3*between(random, 1, 100), 0, 100)) })),
				tests([]string{block(from(300, func(int) int { return 3*between(random, 0, 33) + 1 }))}),
				tests(from(10000, func(int) string { return block(list(random, 15, 0, 100)) })),
			},
			[]string{
				tests(from(5, func(int) string { return block(from(30000, func(int) int { return 3 * between(random, 0, 33) })) })),
				tests(from(5, func(int) string { return block(list(random, 30000, 0, 100)) })),
				tests([]string{block(from(30000, func(int) int { return 3*between(random, 0, 32) + 2 }))}),
			},
		)
	},
	"check-if-array-pairs-are-divisible-by-k": func(random func() float64) []string {
		// A search of a partner for each number makes about n² / 2 = 5 · 10^9 checks at n = 10^5.
		k := between(random, 2, 100)
		medium := pairs(random, 1000, k)
		return call(slices.Concat(
			[][]any{{[]int{1, 2, 3, 4, 5, 10, 6, 7, 8, 9}, 5}, {[]int{1, 2, 3, 4, 5, 6}, 7}, {[]int{1, 2, 3, 4, 5, 6}, 10}, {[]int{-7, 5}, 1}, {slices.Repeat([]int{0}, 6), 7}, {[]int{3, 3, 3, 0, 1, 5}, 6}, {[]int{-1, -9, -2, -8, -3, -7, -5, -10}, 5}},
			[][]any{{medium, k}, {broken(random, medium), k}},
			[][]any{{pairs(random, 100000, 100000), 100000}, {list(random, 100000, -1e9, 1e9), 2}, {broken(random, pairs(random, 100000, 99999)), 99999}},
		))
	},
	"subarray-sums-divisible-by-k": func(random func() float64) []string {
		// Brute Force with a running sum checks n(n + 1) / 2 = 450 015 000 Subarray at n = 30 000. Even numbers with k = 2 give this largest answer.
		return call(slices.Concat(
			[][]any{{[]int{4, 5, 0, -2, -3, 1}, 5}, {[]int{5}, 9}, {[]int{0, 0, 0}, 2}, {slices.Repeat([]int{7}, 5), 7}, {from(20, func(i int) int { return (1 - 2*(i%2)) * between(random, 1, 10000) }), between(random, 2, 20)}, {[]int{between(random, -10000, 10000)}, 10000}},
			from(3, func(i int) []any {
				return []any{list(random, []int{100, 1000, 10000}[i], -10000, 10000), between(random, 2, []int{20, 200, 10000}[i])}
			}),
			[][]any{{from(30000, func(int) int { return 2 * between(random, -5000, 5000) }), 2}, {list(random, 30000, -10000, 10000), 10000}, {slices.Repeat([]int{-10000}, 30000), 9999}},
		))
	},
	"cf-1931-d": func(random func() float64) []string {
		// A check of all pairs makes n(n - 1) / 2 = 19 999 900 000 checks at n = 2 · 10^5.
		return slices.Concat(
			[]string{
				"7\n6 5 2\n1 2 7 4 9 6\n7 9 5\n1 10 15 3 8 12 15\n9 4 10\n14 10 2 2 11 11 13 5 6\n9 5 6\n10 7 6 7 9 7 7 10 10\n9 6 2\n4 9 7 1 2 2 13 3 15\n9 2 3\n14 6 1 15 12 15 8 2 15\n10 5 7\n13 3 3 2 12 11 3 7 13 14\n",
				tests([]string{divisible(1, 1, []int{1, 1})}),
				tests([]string{divisible(2, 3, slices.Repeat([]int{between(random, 1, 1e9)}, 6))}),
				tests([]string{divisible(1e9, 1e9, []int{1e9, 1, 1e9, 999999999, 1e9})}),
				tests([]string{divisible(3, 1, []int{3, 6, 9, 12})}),
			},
			[]string{
				tests(from(100, func(int) string {
					return divisible(between(random, 1, 10), between(random, 1, 10), list(random, between(random, 2, 20), 1, 100))
				})),
				tests(from(10, func(int) string {
					return divisible(between(random, 1, 100), between(random, 1, 100), list(random, between(random, 2, 2000), 1, 1e9))
				})),
				tests(from(10000, func(int) string {
					return divisible(between(random, 1, 20), between(random, 1, 20), list(random, 20, 1, 1e9))
				})),
			},
			[]string{
				tests([]string{divisible(1e9, 999999999, list(random, 200000, 1, 1e9))}),
				tests([]string{divisible(5, 7, list(random, 200000, 1, 1e9))}),
				tests(from(2, func(int) string { return divisible(2, 3, list(random, 100000, 1, 1e9)) })),
			},
		)
	},
	"cf-1594-b": func(random func() float64) []string {
		// A list of the special numbers in increasing order up to the k-th makes up to 10^9 steps for each test.
		return slices.Concat(
			[]string{
				"3\n3 4\n2 12\n105 564\n",
				tests([]string{line(2, 1)}),
				tests([]string{line(2, 1e9)}),
				tests([]string{line(10, 7)}),
				tests([]string{line(10, 8)}),
				tests([]string{line(1e9, 1e9)}),
				tests([]string{line(1e9, 1<<29)}),
				tests([]string{line(1e9, 999999999)}),
			},
			[]string{tests(from(100, func(int) string { return line(between(random, 2, 10), between(random, 1, 1000)) }))},
			[]string{
				tests(from(10000, func(int) string { return line(between(random, 2, 1e9), between(random, 1, 1e9)) })),
				tests(from(10000, func(int) string { return line(2, between(random, 1, 1e9)) })),
				tests(from(10000, func(int) string { return line(between(random, 2, 1e9), 1<<29-1) })),
			},
		)
	},
	"double-modular-exponentiation": func(random func() float64) []string {
		// a ** b has at most 3001 digits here, so the exact power is fast too. The tests check the answers and the order of the two remainders.
		return call(slices.Concat(
			[][]any{
				{[][]int{{2, 3, 3, 10}, {3, 3, 3, 1}, {6, 1, 1, 4}}, 2},
				{[][]int{{39, 3, 1000, 1000}}, 17},
				{four(random, 5, map[int]int{3: 1}), 0},
				{four(random, 5, nil), 1000},
				{slices.Repeat([][]int{{1000, 1000, 1000, 1000}}, 3), 0},
				{from(20, func(i int) []int {
					return []int{10*between(random, 1, 99) + []int{0, 1, 5, 6}[i%4], between(random, 1, 1000), between(random, 1, 1000), between(random, 1, 10)}
				}), between(random, 0, 9)},
			},
			[][]any{
				{four(random, 10, nil), between(random, 0, 9)},
				{four(random, 30, nil), 1},
				{four(random, 50, map[int]int{3: 10}), between(random, 0, 9)},
			},
			[][]any{
				{slices.Repeat([][]int{{1000, 1000, 1000, 1000}}, 100), 0},
				{four(random, 100, nil), 1000},
				{four(random, 100, map[int]int{0: 2}), between(random, 0, 9)},
			},
		))
	},
	"count-good-numbers": func(random func() float64) []string {
		// A loop of n multiplications makes 10^15 steps at the largest n.
		return call(slices.Concat(
			[][]any{{1}, {2}, {3}, {4}, {50}},
			from(4, func(int) []any { return []any{between(random, 1e6, 1e15)} }),
			[][]any{{1 << 49}, {1e15 - 1}, {1e15}},
		))
	},
	"super-pow": func(random func() float64) []string {
		// A loop up to the exponent makes up to 10^2000 steps. The int of the digits with pow is fast too: the exponent has about 6644 binary digits.
		return call(slices.Concat(
			[][]any{{2, []int{3}}, {2, []int{1, 0}}, {1, []int{4, 3, 3, 8, 5, 2}}, {1337, digits(random, 5)}, {1338, digits(random, 5)}, {2674, digits(random, 5)}},
			from(3, func(i int) []any {
				return []any{between(random, 1, math.MaxInt32), digits(random, []int{10, 100, 1000}[i])}
			}),
			[][]any{{between(random, 1, math.MaxInt32), digits(random, 2000)}, {math.MaxInt32, slices.Repeat([]int{9}, 2000)}, {2, append([]int{1}, make([]int, 1999)...)}},
		))
	},
	"cf-1236-b": func(random func() float64) []string {
		// A loop of n or m multiplications makes up to 10^9 steps. 2 has the order 500 000 003 by the modulus 10^9 + 7, so m = 500 000 003 gives 2^m - 1 = 0.
		return slices.Concat(
			[]string{line(1, 3), line(2, 2), line(1, 1), line(5, 2), line(1, 500000003)},
			from(4, func(int) string { return line(between(random, 1, 1e9), between(random, 1, 1e9)) }),
			[]string{line(1e9, 1), line(1, 1e9), line(1e9, 1e9)},
		)
	},
	"prime-arrangements": func(random func() float64) []string {
		// Backtracking over all n! Permutations stops working at n ≥ 13. The limit n ≤ 100 lets each correct method pass, so the tests check the answers.
		return call(slices.Concat(
			[][]any{{1}, {2}, {3}, {4}, {5}, {10}, {23}, {97}},
			from(3, func(int) []any { return []any{between(random, 11, 96)} }),
			[][]any{{98}, {99}, {100}},
		))
	},
	"cf-2008-f": func(random func() float64) []string {
		// A check of all pairs makes n(n - 1) / 2 = 19 999 900 000 products at n = 2 · 10^5. 10^9 = M - 7, so the product of two values 10^9 gives 49.
		return slices.Concat(
			[]string{
				"3\n3\n3 2 3\n4\n2 2 2 4\n5\n1 2 3 4 5\n",
				tests([]string{block([]int{0, 0})}),
				tests([]string{block([]int{1e9, 1e9})}),
				tests([]string{block([]int{1, 1, 2})}),
				tests([]string{block([]int{0, 1e9, 0, 1e9})}),
			},
			[]string{
				tests(from(100, func(int) string { return block(list(random, between(random, 2, 20), 0, 10)) })),
				tests(from(10, func(int) string { return block(list(random, between(random, 2, 2000), 0, 1e9)) })),
				tests(from(2, func(int) string {
					a := slices.Concat(make([]int, 500), list(random, 500, 1, 1e9))
					shuffle(random, a)
					return block(a)
				})),
				tests(from(10000, func(int) string { return block(list(random, 20, 0, 1e9)) })),
			},
			[]string{
				tests([]string{block(list(random, 200000, 0, 1e9))}),
				tests([]string{block(list(random, 200000, 1e9-10, 1e9))}),
				tests([]string{block(slices.Repeat([]int{1e9}, 200000))}),
			},
		)
	},
	"cf-1569-c": func(random func() float64) []string {
		// The exact n! by a loop of multiplications takes seconds at n = 2 · 10^5, but math.factorial passes. The n = 13 test and the k = 5 test reject the integer division of the remainder n! % M by k + 1.
		// uniqueMax returns one test with one maximum top, k numbers top - 1 and n - k - 1 random numbers below top - 1, in a random order.
		uniqueMax := func(n, top, k int) string {
			a := slices.Concat([]int{top}, slices.Repeat([]int{top - 1}, k), list(random, n-k-1, 1, top-2))
			shuffle(random, a)
			return block(a)
		}
		return slices.Concat(
			[]string{
				"4\n2\n1 2\n3\n5 5 5\n4\n1 3 3 7\n6\n3 4 2 1 3 3\n",
				tests([]string{block([]int{7, 7})}),
				tests([]string{block([]int{1, 1, 3})}),
				tests([]string{block([]int{2, 2, 3})}),
				tests([]string{block([]int{1e9, 999999999, 1, 1, 1})}),
				tests([]string{block(slices.Concat([]int{5, 4, 4, 4}, slices.Repeat([]int{1}, 9)))}),
			},
			[]string{
				tests(from(100, func(int) string { return block(list(random, between(random, 2, 8), 1, 5)) })),
				tests(from(10000, func(int) string { return block(list(random, 20, 1, 4)) })),
				tests(from(2, func(int) string { return block(list(random, 100000, 1, 1e9)) })),
			},
			[]string{
				tests([]string{uniqueMax(200000, 1e9, 5)}),
				tests([]string{block(slices.Concat([]int{1e9, 1e9}, list(random, 199998, 1, 1e9)))}),
				tests([]string{uniqueMax(200000, 1e9, 199999)}),
			},
		)
	},
	"count-anagrams": func(random func() float64) []string {
		// Backtracking over the orders of the letters stops working at once. The exact math.factorial with the integer division takes about 0.4 s in CPython at n = 10^5 and passes, so the tests check the answers.
		alphabet := []byte("abcdefghijklmnopqrstuvwxyz")
		// word returns n random letters from the first letters of the alphabet.
		word := func(n, letters int) string {
			return string(from(n, func(int) byte { return alphabet[between(random, 0, letters-1)] }))
		}
		// words returns random words of length 1 to longest with one space between them, with a total length of at most n.
		words := func(n, longest, letters int) string {
			var parts []string
			for size := 0; ; {
				w := word(between(random, 1, longest), letters)
				if size+len(w)+len(parts) > n {
					return strings.Join(parts, " ")
				}
				size += len(w)
				parts = append(parts, w)
			}
		}
		letters := slices.Clone(alphabet)
		shuffle(random, letters)
		return call(slices.Concat(
			[][]any{{"too hot"}, {"aa"}, {"a"}, {"abc def"}, {"aab bba ccc"}, {strings.Repeat("a", 15) + strings.Repeat("b", 15)}, {string(letters)}},
			[][]any{{words(9999, 1, 26)}, {words(10000, 10, 3)}, {words(30000, 100, 26)}},
			[][]any{
				{word(100000, 2)},
				{word(100000, 26)},
				{strings.Join(from(5, func(int) string { return word(19999, 26) }), " ")},
			},
		))
	},
}

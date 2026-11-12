package permutations

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The hidden test generators of the set. Each one gets a seeded random() in [0, 1) and returns the inputs in the test format of the task:
// the edge cases of the statement first, then random medium inputs, then 2 inputs at the maximum constraint sizes.

// randomInt returns an integer in [lo, hi].
func randomInt(random func() float64, lo, hi int) int {
	return lo + int(math.Floor(random()*float64(hi-lo+1)))
}

func randomInts(random func() float64, n, lo, hi int) []int {
	return repeat(n, func() int { return randomInt(random, lo, hi) })
}

// shuffle is the Fisher-Yates shuffle in place.
func shuffle[T any](random func() float64, list []T) []T {
	for i := len(list) - 1; i > 0; i-- {
		j := randomInt(random, 0, i)
		list[i], list[j] = list[j], list[i]
	}
	return list
}

// span returns the integers from lo to hi in ascending order.
func span(lo, hi int) []int {
	list := make([]int, hi-lo+1)
	for i := range list {
		list[i] = lo + i
	}
	return list
}

// sample returns n different integers in [lo, hi] in random order.
func sample(random func() float64, n, lo, hi int) []int {
	return shuffle(random, span(lo, hi))[:n]
}

// randomText returns a string of n random characters of the alphabet.
func randomText(random func() float64, n int, alphabet string) string {
	var out strings.Builder
	for range n {
		out.WriteByte(alphabet[randomInt(random, 0, len(alphabet)-1)])
	}
	return out.String()
}

func repeat[T any](times int, next func() T) []T {
	out := make([]T, times)
	for i := range out {
		out[i] = next()
	}
	return out
}

// lines returns the stdin text with one line for each row.
func lines(rows []string) string {
	return strings.Join(rows, "\n") + "\n"
}

// join writes the numbers with sep between them.
func join(list []int, sep string) string {
	words := make([]string, len(list))
	for i, x := range list {
		words[i] = strconv.Itoa(x)
	}
	return strings.Join(words, sep)
}

// Generators are the hidden test generators of the tasks by problem slug.
var Generators = map[string]sets.Generator{
	"cf-1821-a": func(random func() float64) []string {
		cases := func(t, lo, hi int) string {
			return lines(append([]string{strconv.Itoa(t)}, repeat(t, func() string {
				return randomText(random, randomInt(random, lo, hi), "???0123456789")
			})...))
		}
		return slices.Concat(
			[]string{"1\n?\n", "1\n0\n", "1\n?????\n", "1\n99999\n", "4\n0????\n00000\n0\n0?\n", "3\n10000\n?0000\n9\n"},
			repeat(6, func() string { return cases(randomInt(random, 2, 30), 1, 5) }),
			[]string{cases(20000, 5, 5), cases(20000, 1, 5)},
		)
	},
	"cf-1743-a": func(random func() float64) []string {
		cases := func(t, lo, hi int) string {
			return lines(append([]string{strconv.Itoa(t)}, repeat(t, func() string {
				n := randomInt(random, lo, hi)
				digits := sample(random, n, 0, 9)
				slices.Sort(digits)
				return strconv.Itoa(n) + "\n" + join(digits, " ")
			})...))
		}
		return slices.Concat(
			[]string{"1\n1\n0\n", "1\n1\n9\n", "1\n8\n0 1 2 3 4 5 6 7\n", "1\n8\n2 3 4 5 6 7 8 9\n", "2\n2\n0 9\n7\n1 2 3 4 5 6 7\n"},
			repeat(6, func() string { return cases(randomInt(random, 2, 20), 1, 8) }),
			[]string{cases(200, 8, 8), cases(200, 1, 8)},
		)
	},
	"cf-629-a": func(random func() float64) []string {
		// A cake of side n where each cell has a chocolate with the probability p.
		cake := func(n int, p float64) string {
			return lines(append([]string{strconv.Itoa(n)}, repeat(n, func() string {
				return strings.Join(repeat(n, func() string {
					if random() < p {
						return "C"
					}
					return "."
				}), "")
			})...))
		}
		return slices.Concat(
			[]string{"1\n.\n", "1\nC\n", "2\nCC\nCC\n", "3\n...\n...\n...\n", "4\nCCCC\n....\n....\n....\n", "4\nC...\n.C..\n..C.\n...C\n"},
			repeat(6, func() string { return cake(randomInt(random, 2, 30), random()) }),
			[]string{cake(100, 1), cake(100, random())},
		)
	},
	"cf-1674-b": func(random func() float64) []string {
		var words []string
		for a := 'a'; a <= 'z'; a++ {
			for b := 'a'; b <= 'z'; b++ {
				if b != a {
					words = append(words, string([]rune{a, b}))
				}
			}
		}
		cases := func(list []string) string { return lines(append([]string{strconv.Itoa(len(list))}, list...)) }
		word := func() string { return words[randomInt(random, 0, 649)] }
		return slices.Concat(
			[]string{"1\nab\n", "1\nzy\n", "1\naz\n", "1\nba\n", "1\nza\n", "2\nyz\nmn\n"},
			repeat(6, func() string { return cases(repeat(randomInt(random, 2, 30), word)) }),
			[]string{cases(shuffle(random, slices.Clone(words))), cases(repeat(650, word))},
		)
	},
	"number-of-good-pairs": func(random func() float64) []string {
		return slices.Concat(
			[]string{sets.Call([]int{1}), sets.Call([]int{100}), sets.Call([]int{7, 7}), sets.Call([]int{1, 2}), sets.Call([]int{100, 1, 100, 1, 100}), sets.Call(span(1, 100))},
			repeat(6, func() string {
				return sets.Call(randomInts(random, randomInt(random, 2, 60), 1, randomInt(random, 2, 100)))
			}),
			[]string{sets.Call(slices.Repeat([]int{100}, 100)), sets.Call(randomInts(random, 100, 1, 100))},
		)
	},
	"letter-case-permutation": func(random func() float64) []string {
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		return slices.Concat(
			[]string{sets.Call("a"), sets.Call("Z"), sets.Call("7"), sets.Call("012345678901"), sets.Call("aB3")},
			repeat(6, func() string { return sets.Call(randomText(random, randomInt(random, 2, 9), letters+"0123456789")) }),
			[]string{sets.Call(randomText(random, 12, letters)), sets.Call(randomText(random, 12, letters))},
		)
	},
	"subsets": func(random func() float64) []string {
		return slices.Concat(
			[]string{sets.Call([]int{0}), sets.Call([]int{-10}), sets.Call([]int{10}), sets.Call([]int{-10, 10}), sets.Call([]int{1, 2, 3, 4}), sets.Call([]int{5, 4, 3, 2, 1})},
			repeat(6, func() string { return sets.Call(sample(random, randomInt(random, 2, 8), -10, 10)) }),
			[]string{sets.Call(sample(random, 10, -10, 10)), sets.Call(sample(random, 10, -10, 10))},
		)
	},
	"permutations": func(random func() float64) []string {
		return slices.Concat(
			[]string{sets.Call([]int{0}), sets.Call([]int{-10}), sets.Call([]int{10, -10}), sets.Call([]int{1, 2, 3}), sets.Call([]int{3, 2, 1})},
			repeat(6, func() string { return sets.Call(sample(random, randomInt(random, 2, 5), -10, 10)) }),
			[]string{sets.Call(sample(random, 6, -10, 10)), sets.Call(sample(random, 6, -10, 10))},
		)
	},
	"combinations": func(random func() float64) []string {
		return slices.Concat(
			[]string{sets.Call(1, 1), sets.Call(2, 1), sets.Call(2, 2), sets.Call(20, 1), sets.Call(20, 20), sets.Call(6, 6)},
			repeat(6, func() string {
				n := randomInt(random, 3, 12)
				return sets.Call(n, randomInt(random, 1, n))
			}),
			// The answer for k from 6 to 15 is longer than the output limit of the judge (1,000,000 characters).
			[]string{sets.Call(20, 5), sets.Call(20, 16)},
		)
	},
	"letter-combinations-of-a-phone-number": func(random func() float64) []string {
		return slices.Concat(
			[]string{sets.Call("2"), sets.Call("7"), sets.Call("9"), sets.Call("22"), sets.Call("79")},
			repeat(6, func() string { return sets.Call(randomText(random, randomInt(random, 2, 3), "23456789")) }),
			[]string{sets.Call(randomText(random, 4, "23456789")), sets.Call(randomText(random, 4, "79"))},
		)
	},
	"subsets-ii": func(random func() float64) []string {
		return slices.Concat(
			[]string{sets.Call([]int{0}), sets.Call([]int{5, 5}), sets.Call([]int{-10, -10, -10}), sets.Call([]int{1, 2, 2}), sets.Call([]int{3, 1, 2}), sets.Call([]int{2, 1, 2, 1}), sets.Call(slices.Repeat([]int{10}, 10))},
			repeat(6, func() string { return sets.Call(randomInts(random, randomInt(random, 3, 8), -3, 3)) }),
			[]string{sets.Call(sample(random, 10, -10, 10)), sets.Call(randomInts(random, 10, -10, 10))},
		)
	},
	"permutations-ii": func(random func() float64) []string {
		tests := slices.Concat(
			[]string{sets.Call([]int{0}), sets.Call([]int{1, 1}), sets.Call([]int{2, 2, 2}), sets.Call([]int{1, 1, 2}), sets.Call([]int{-10, 10}), sets.Call(slices.Repeat([]int{-3}, 8))},
			repeat(6, func() string { return sets.Call(randomInts(random, randomInt(random, 3, 6), -2, 2)) }),
		)
		// 8 different values give 40,320 permutations, longer than the output limit of the judge (1,000,000 characters). One repeated value halves the count. 8 values in [-3, 3] always repeat a value.
		values := sample(random, 7, -10, 10)
		return append(tests, sets.Call(append([]int{values[0]}, values...)), sets.Call(randomInts(random, 8, -3, 3)))
	},
	"combination-sum": func(random func() float64) []string {
		// The count of combinations with the sum target (coin change). The statement allows fewer than 150.
		count := func(candidates []int, target int) int {
			ways := make([]int, target+1)
			ways[0] = 1
			for _, c := range candidates {
				for s := c; s <= target; s++ {
					ways[s] += ways[s-c]
				}
			}
			return ways[target]
		}
		// The largest target up to top with fewer than 150 combinations. The target 1 always fits, because each candidate is at least 2.
		fit := func(candidates []int, top int) string {
			target := top
			for count(candidates, target) >= 150 {
				target--
			}
			return sets.Call(candidates, target)
		}
		return slices.Concat(
			[]string{sets.Call([]int{2}, 1), sets.Call([]int{40}, 40), sets.Call([]int{2}, 40), sets.Call([]int{3}, 40), sets.Call([]int{39, 40}, 1), sets.Call([]int{2, 3}, 7)},
			repeat(6, func() string { return fit(sample(random, randomInt(random, 2, 10), 2, 20), randomInt(random, 1, 40)) }),
			// The candidates 11 to 40 give 19 combinations for the target 40.
			[]string{sets.Call(sample(random, 30, 11, 40), 40), fit(sample(random, 30, 2, 40), 40)},
		)
	},
	"generate-parentheses": func(random func() float64) []string {
		tests := []string{sets.Call(1), sets.Call(2)}
		for _, n := range shuffle(random, []int{3, 4, 5, 6, 7}) {
			tests = append(tests, sets.Call(n))
		}
		return append(tests, sets.Call(randomInt(random, 3, 7)), sets.Call(8), sets.Call(8))
	},
	"cf-1097-b": func(random func() float64) []string {
		turns := func(angles []int) string { return lines([]string{strconv.Itoa(len(angles)), join(angles, "\n")}) }
		tests := slices.Concat(
			[]string{"1\n1\n", "1\n180\n", "2\n180\n180\n", "2\n90\n90\n", "4\n90\n90\n90\n90\n", "3\n120\n120\n120\n", "2\n1\n179\n"},
			repeat(6, func() string { return turns(randomInts(random, randomInt(random, 2, 10), 1, 180)) }),
		)
		// 15 odd angles always give an odd sum, so the answer is NO after the full search.
		odd := randomInts(random, 15, 0, 89)
		for i, a := range odd {
			odd[i] = 2*a + 1
		}
		return append(tests, turns(odd), turns(randomInts(random, 15, 1, 180)))
	},
	"next-permutation": func(random func() float64) []string {
		tests := slices.Concat(
			[]string{sets.Call([]int{0}), sets.Call([]int{100}), sets.Call([]int{1, 2}), sets.Call([]int{2, 1}), sets.Call([]int{5, 5, 5}), sets.Call([]int{3, 2, 1}), sets.Call([]int{1, 2, 3, 4, 5}), sets.Call([]int{0, 100, 100, 0})},
			repeat(6, func() string {
				return sets.Call(randomInts(random, randomInt(random, 3, 20), 0, randomInt(random, 1, 100)))
			}),
			[]string{sets.Call(randomInts(random, 100, 0, 100))},
		)
		descending := randomInts(random, 100, 0, 100)
		slices.Sort(descending)
		slices.Reverse(descending)
		return append(tests, sets.Call(descending))
	},
	"unique-paths": func(random func() float64) []string {
		// C(m + n - 2, m - 1). The statement allows at most 2 · 10^9. Each step value C(n - 1 + i, i) is an integer below 2^53.
		paths := func(m, n int) int {
			out := 1
			for i := 1; i < m && out <= 2e9; i++ {
				out = out * (n - 1 + i) / i
			}
			return out
		}
		// A random n in [1, 100] with at most 2 · 10^9 paths for this m.
		fit := func(m int) string {
			top := 100
			for paths(m, top) > 2e9 {
				top--
			}
			return sets.Call(m, randomInt(random, 1, top))
		}
		return slices.Concat(
			[]string{sets.Call(1, 1), sets.Call(1, 100), sets.Call(100, 1), sets.Call(2, 2), sets.Call(2, 100), sets.Call(100, 2), sets.Call(17, 18), sets.Call(18, 17)},
			repeat(6, func() string { return fit(randomInt(random, 1, 100)) }),
			// C(105, 6) = 1609344100.
			[]string{sets.Call(100, 7), sets.Call(7, 100)},
		)
	},
	"permutation-sequence": func(random func() float64) []string {
		factorial := func(n int) int {
			out := 1
			for i := 2; i <= n; i++ {
				out *= i
			}
			return out
		}
		pick := func(n int) string { return sets.Call(n, randomInt(random, 1, factorial(n))) }
		return slices.Concat(
			[]string{sets.Call(1, 1), sets.Call(2, 1), sets.Call(2, 2), sets.Call(3, 6), sets.Call(4, 24), sets.Call(9, 1), sets.Call(9, 362880)},
			repeat(6, func() string { return pick(randomInt(random, 3, 8)) }),
			[]string{pick(9), pick(9)},
		)
	},
	"n-queens": func(random func() float64) []string {
		tests := []string{sets.Call(1), sets.Call(2), sets.Call(3)}
		for _, n := range shuffle(random, []int{4, 5, 6, 7, 8}) {
			tests = append(tests, sets.Call(n))
		}
		return append(tests, sets.Call(9), sets.Call(9))
	},
}

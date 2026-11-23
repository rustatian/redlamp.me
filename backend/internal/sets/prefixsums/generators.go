package prefixsums

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The hidden tests of the set "prefix-sums". Each generator gives the edge cases of the statement, then random medium inputs, then the inputs at the largest sizes.
// A call input is a Python literal list of the arguments. JSON of integers is a Python literal too.

// from returns n values f(0), ..., f(n - 1).
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

func call(cases [][]any) []string {
	return from(len(cases), func(i int) string { return sets.Call(cases[i]...) })
}

// line writes the numbers on one line with one space between them.
func line(a ...int) string {
	text := make([]string, len(a))
	for i, x := range a {
		text[i] = strconv.Itoa(x)
	}
	return strings.Join(text, " ") + "\n"
}

// signs returns n values 1 or -1. The values of big go to random places, at most one value at each place.
func signs(random func() float64, n int, big []int) []int {
	a := from(n, func(int) int { return 2*between(random, 0, 1) - 1 })
	places := from(n, func(i int) int { return i })
	for i, x := range big {
		j := between(random, i, n-1)
		places[i], places[j] = places[j], places[i]
		a[places[i]] = x
	}
	return a
}

// oddQueries writes one test of Codeforces 1807D: n values from 1 to hi and q random queries. whole makes every query cover the whole Array.
func oddQueries(random func() float64, n, q, hi int, whole bool) string {
	var text strings.Builder
	text.WriteString(line(n, q) + line(list(random, n, 1, hi)...))
	for range q {
		l, r := 1, n
		if !whole {
			l = between(random, 1, n)
			r = between(random, l, n)
		}
		text.WriteString(line(l, r, between(random, 1, hi)))
	}
	return text.String()
}

// stones writes the input of Codeforces 433B: the values and m questions. whole makes every question cover the whole Array.
func stones(random func() float64, v []int, m int, whole bool) string {
	var text strings.Builder
	text.WriteString(line(len(v)) + line(v...) + line(m))
	for range m {
		l, r := 1, len(v)
		if !whole {
			l = between(random, 1, len(v))
			r = between(random, l, len(v))
		}
		text.WriteString(line(between(random, 1, 2), l, r))
	}
	return text.String()
}

// distinct returns n different random integers from lo to hi.
func distinct(random func() float64, n, lo, hi int) []int {
	seen := map[int]bool{}
	values := make([]int, 0, n)
	for len(values) < n {
		if x := between(random, lo, hi); !seen[x] {
			seen[x] = true
			values = append(values, x)
		}
	}
	return values
}

// glasses writes one test of Codeforces 1915E: the number of glasses and the amounts of juice.
func glasses(a []int) string {
	return line(len(a)) + line(a...)
}

// increasing returns n strictly increasing values from 1 to 5000n: the gaps are random from 1 to 5000.
func increasing(random func() float64, n int) []int {
	a := make([]int, n)
	last := 0
	for i := range a {
		last += between(random, 1, 5000)
		a[i] = last
	}
	return a
}

// digits writes one test of Codeforces 1398C: n and a string of n random digits from 0 to hi.
func digits(random func() float64, n, hi int) string {
	text := make([]byte, n)
	for i := range text {
		text[i] = byte('0' + between(random, 0, hi))
	}
	return line(n) + string(text) + "\n"
}

// tests writes t tests of one kind after the line with t.
func tests(t int, test func() string) string {
	return line(t) + strings.Join(from(t, func(int) string { return test() }), "")
}

// segment returns a random segment [l, r] of 1..size with at most length numbers.
func segment(random func() float64, size, length int) (int, int) {
	l := between(random, 1, size)
	return l, between(random, l, min(size, l+length-1))
}

// bookings returns count random bookings of LeetCode 1109 on n flights: at most length flights and from 1 to seats seats each.
func bookings(random func() float64, count, n, length, seats int) [][]int {
	return from(count, func(int) []int {
		first, last := segment(random, n, length)
		return []int{first, last, between(random, 1, seats)}
	})
}

// coffee writes the input of Codeforces 816B: n recipes of at most length temperatures from 1 to hi, and q random questions. whole makes every question [1, 200000].
func coffee(random func() float64, n, k, q, hi, length int, whole bool) string {
	var text strings.Builder
	text.WriteString(line(n, k, q))
	for range n {
		text.WriteString(line(segment(random, hi, length)))
	}
	for range q {
		a, b := 1, 200000
		if !whole {
			a, b = segment(random, hi, hi)
		}
		text.WriteString(line(a, b))
	}
	return text.String()
}

// operationsInput writes the input of Codeforces 295A: the values a, m operations of at most opLength elements with d from 0 to hiD, and k queries of at most queryLength operations. whole makes every query [1, m].
func operationsInput(random func() float64, a []int, m, k, opLength, hiD, queryLength int, whole bool) string {
	var text strings.Builder
	text.WriteString(line(len(a), m, k) + line(a...))
	for range m {
		l, r := segment(random, len(a), opLength)
		text.WriteString(line(l, r, between(random, 0, hiD)))
	}
	for range k {
		x, y := 1, m
		if !whole {
			x, y = segment(random, m, queryLength)
		}
		text.WriteString(line(x, y))
	}
	return text.String()
}

// matrix returns m rows of n random values from lo to hi.
func matrix(random func() float64, m, n, lo, hi int) [][]int {
	return from(m, func(int) []int { return list(random, n, lo, hi) })
}

// filled returns m rows of n values x.
func filled(m, n, x int) [][]int {
	return from(m, func(int) []int { return slices.Repeat([]int{x}, n) })
}

// submatrices returns count random queries [row1, col1, row2, col2] of LeetCode 2536 on an n x n Matrix, with at most length rows.
func submatrices(random func() float64, count, n, length int) [][]int {
	return from(count, func(int) []int {
		r1, r2 := segment(random, n, length)
		c1, c2 := segment(random, n, n)
		return []int{r1 - 1, c1 - 1, r2 - 1, c2 - 1}
	})
}

// rectangles writes one test of Codeforces 1722E: n rectangles with sides from 1 to hi and q random queries with sides up to min(hi + 1, 1000). whole makes every query 1 1 1000 1000.
func rectangles(random func() float64, n, q, hi int, whole bool) string {
	var text strings.Builder
	text.WriteString(line(n, q))
	for range n {
		text.WriteString(line(between(random, 1, hi), between(random, 1, hi)))
	}
	top := min(hi+1, 1000)
	for range q {
		if whole {
			text.WriteString(line(1, 1, 1000, 1000))
			continue
		}
		hs, ws := between(random, 1, top-1), between(random, 1, top-1)
		text.WriteString(line(hs, ws, between(random, hs+1, top), between(random, ws+1, top)))
	}
	return text.String()
}

// Generators are the hidden test generators of the tasks by problem slug.
var Generators = map[string]sets.Generator{
	"find-pivot-index": func(random func() float64) []string {
		// Only index 9998 is a pivot: 4999 pairs (v, -v) with v > 0, then x < 0, then 0. Index 2j has the left sum 0 and the right sum x - v, index 2j + 1 has v and x.
		late := slices.Concat(slices.Concat(from(4999, func(int) []int {
			v := between(random, 1, 1000)
			return []int{v, -v}
		})...), []int{between(random, -1000, -1), 0})
		return call(slices.Concat(
			[][]any{{[]int{0}}, {[]int{1000}}, {[]int{-1000, 1000}}, {[]int{0, 0}}, {[]int{1, -1, 3}}, {slices.Repeat([]int{1}, 9999)}, {slices.Repeat([]int{1}, 10000)}},
			// Values from -3 to 3 give frequent pivots.
			from(4, func(int) []any { return []any{list(random, between(random, 2, 60), -3, 3)} }),
			// The sums of slices run in C and can pass at n = 10^4. Nested Python loops make about 10^8 additions, near the 3 s limit in Pyodide.
			[][]any{{list(random, 10000, -1000, 1000)}, {late}},
		))
	},
	"cf-1807-d": func(random func() float64) []string {
		// 5 * 10^4 queries over the whole Array make a new sum for each query cost 2.5 * 10^9 additions. The answer has at most 5 * 10^4 lines of 4 characters.
		return slices.Concat(
			[]string{"1\n1 1\n1\n1 1 1\n", "1\n1 1\n1000000000\n1 1 1000000000\n", "1\n2 2\n1 2\n1 2 1\n2 2 1\n", tests(10000, func() string { return oddQueries(random, 1, 1, 1e9, false) })},
			from(4, func(i int) string {
				hi := []int{10, 1e9}[i%2]
				return tests(between(random, 1, 5), func() string { return oddQueries(random, between(random, 1, 30), between(random, 1, 30), hi, false) })
			}),
			[]string{
				tests(1, func() string { return oddQueries(random, 50000, 50000, 1e9, false) }),
				tests(1, func() string { return oddQueries(random, 50000, 50000, 1e9, true) }),
				tests(2, func() string { return oddQueries(random, 25000, 25000, 1e9, false) }),
			},
		)
	},
	"product-of-array-except-self": func(random func() float64) []string {
		// Every input keeps the product of the absolute values of its nonzero elements at most 2^31 - 1, so every prefix, suffix and answer fits in 32 bits.
		return call(slices.Concat(
			[][]any{{[]int{0, 0}}, {[]int{2, 3}}, {[]int{30, -30}}, {[]int{0, 5}}, {slices.Repeat([]int{1}, 100000)}, {slices.Repeat([]int{-1}, 100000)}},
			from(4, func(int) []any {
				n := between(random, 2, 60)
				return []any{signs(random, n, list(random, between(random, 0, min(n, 6)), -30, 30))}
			}),
			// The answers of the first input are at most 21^3 = 9261 in absolute value, so the answer text stays below 10^6 characters. A product of all other numbers for each i makes 10^10 multiplications.
			[][]any{{signs(random, 100000, list(random, 3, 2, 21))}, {signs(random, 100000, slices.Concat([]int{0}, list(random, 6, 2, 30)))}, {signs(random, 100000, []int{0, 0})}},
		))
	},
	"cf-433-b": func(random func() float64) []string {
		// A new sum for each question of the first largest input adds 10^9 numbers. Its answers have about 15 characters, about 375 000 characters in all.
		var pairs []string
		for kind := 1; kind <= 2; kind++ {
			for l := 1; l <= 5; l++ {
				for r := l; r <= 5; r++ {
					pairs = append(pairs, line(kind, l, r))
				}
			}
		}
		return slices.Concat(
			[]string{"1\n1\n2\n1 1 1\n2 1 1\n", "1\n1000000000\n2\n2 1 1\n1 1 1\n", line(5) + line(5, 4, 3, 2, 1) + line(len(pairs)) + strings.Join(pairs, ""), stones(random, slices.Repeat([]int{7}, 6), 10, false)},
			from(4, func(i int) string {
				return stones(random, list(random, between(random, 1, 30), 1, []int{10, 1e9}[i%2]), between(random, 1, 30), false)
			}),
			[]string{stones(random, list(random, 40000, 1, 1e9), 25000, true), stones(random, list(random, 40000, 1, 10), 40000, false)},
		)
	},
	"contains-duplicate": func(random func() float64) []string {
		// Different values make a set solution read the whole input. Brute Force over all pairs makes n(n - 1)/2, about 5 * 10^9 comparisons at n = 10^5.
		repeatAtEnd := func() []int {
			a := distinct(random, 100000, -1e9, 1e9)
			a[len(a)-1] = a[0]
			return a
		}
		return call(slices.Concat(
			[][]any{{[]int{0}}, {[]int{1e9, -1e9}}, {[]int{-1e9, -1e9}}, {[]int{1, 2, 1}}, {[]int{3, 1, 2}}},
			from(4, func(i int) []any {
				hi := []int{5, 1e9}[i%2]
				return []any{list(random, between(random, 1, 60), -hi, hi)}
			}),
			[][]any{{distinct(random, 100000, -1e9, 1e9)}, {repeatAtEnd()}, {from(100000, func(i int) int { return 1e9 - i })}},
		))
	},
	"subarray-sum-equals-k": func(random func() float64) []string {
		// Brute Force over all starts with a running sum makes n(n + 1)/2 = 2 * 10^8 additions at n = 2 * 10^4.
		return call(slices.Concat(
			[][]any{{[]int{1}, 1}, {[]int{-1000}, -1000}, {[]int{1000}, 0}, {slices.Repeat([]int{0}, 20000), 0}, {slices.Repeat([]int{1000}, 20000), 1e7}, {slices.Repeat([]int{-1000, 1000}, 10000), 0}},
			from(4, func(int) []any { return []any{list(random, between(random, 1, 60), -3, 3), between(random, -5, 5)} }),
			[][]any{{list(random, 20000, -1000, 1000), between(random, -3000, 3000)}, {list(random, 20000, -2, 2), 1}},
		))
	},
	"contiguous-array": func(random func() float64) []string {
		// Brute Force over all Subarrays with running counts makes n(n + 1)/2 = 5 * 10^9 steps at n = 10^5.
		return call(slices.Concat(
			[][]any{{[]int{0}}, {[]int{1}}, {[]int{0, 1}}, {[]int{1, 1, 0, 0}}, {slices.Repeat([]int{0}, 100000)}, {slices.Repeat([]int{0, 1}, 50000)}, {slices.Concat(slices.Repeat([]int{1}, 50000), slices.Repeat([]int{0}, 50000))}},
			from(3, func(int) []any { return []any{list(random, between(random, 1, 60), 0, 1)} }),
			// The second largest input has about 60 % ones.
			[][]any{{list(random, 100000, 0, 1)}, {from(100000, func(int) int { return min(1, between(random, 0, 4)/2) })}},
		))
	},
	"cf-1915-e": func(random func() float64) []string {
		// Strictly increasing amounts make the odd prefixes positive and increasing and the even prefixes negative and decreasing: no prefix repeats, so the answer is NO after a full scan.
		// With the last amount equal to the amount before it, the last prefix equals the prefix two places before it: the first repeat is at the end.
		// Brute Force over all pairs l, r makes 2 * 10^10 steps at n = 2 * 10^5.
		late := func() []int {
			a := increasing(random, 200000)
			a[len(a)-1] = a[len(a)-2]
			return a
		}
		return slices.Concat(
			[]string{"1\n1\n1\n", "1\n2\n1 1\n", "1\n2\n1000000000 1\n", "1\n3\n2 4 2\n", tests(10000, func() string { return glasses(list(random, 20, 1, 3)) })},
			from(4, func(i int) string {
				hi := []int{5, 1e9}[i%2]
				return tests(between(random, 1, 5), func() string { return glasses(list(random, between(random, 1, 30), 1, hi)) })
			}),
			[]string{tests(1, func() string { return glasses(increasing(random, 200000)) }), tests(1, func() string { return glasses(late()) }), tests(2, func() string { return glasses(list(random, 100000, 1, 3)) })},
		)
	},
	"cf-1398-c": func(random func() float64) []string {
		// Brute Force over all Subarrays makes n(n + 1)/2 = 5 * 10^9 steps at n = 10^5. 10^5 ones give the largest answer 5 000 050 000.
		return slices.Concat(
			[]string{"1\n1\n1\n", "1\n1\n0\n", "1\n1\n9\n", "1\n100000\n" + strings.Repeat("1", 100000) + "\n", "1\n100000\n" + strings.Repeat("0", 100000) + "\n", tests(1000, func() string { return digits(random, 100, 2) })},
			from(4, func(i int) string {
				hi := []int{9, 2}[i%2]
				return tests(between(random, 1, 5), func() string { return digits(random, between(random, 1, 30), hi) })
			}),
			[]string{tests(1, func() string { return digits(random, 100000, 2) }), tests(1, func() string { return digits(random, 100000, 9) })},
		)
	},
	"maximum-population-year": func(random func() float64) []string {
		// The limits are small: a check of every year against every person makes 100 * 101 steps and also passes. This is the easy task of the chapter.
		people := func(n, lo int) [][]int {
			return from(n, func(int) []int {
				birth := between(random, lo, 2049)
				return []int{birth, between(random, birth+1, 2050)}
			})
		}
		return call(slices.Concat(
			[][]any{{[][]int{{1950, 2050}}}, {[][]int{{2049, 2050}}}, {[][]int{{1950, 1951}, {2049, 2050}}}, {slices.Repeat([][]int{{1950, 2050}}, 100)}},
			from(4, func(int) []any { return []any{people(between(random, 2, 60), 1950)} }),
			[][]any{{people(100, 1950)}, {people(100, 2030)}},
		))
	},
	"corporate-flight-bookings": func(random func() float64) []string {
		// A loop over the flights of each booking makes 4 * 10^8 additions on the first largest input. Its answers are 2 * 10^8 each, about 220 000 characters.
		return call(slices.Concat(
			[][]any{{[][]int{{1, 1, 1}}, 1}, {[][]int{{1, 1, 10000}}, 1}, {[][]int{{1, 20000, 10000}}, 20000}, {[][]int{{20000, 20000, 1}}, 20000}},
			from(4, func(int) []any {
				n := between(random, 1, 60)
				return []any{bookings(random, between(random, 1, 60), n, n, 10000), n}
			}),
			[][]any{{slices.Repeat([][]int{{1, 20000, 10000}}, 20000), 20000}, {bookings(random, 20000, 20000, 20000, 10000), 20000}, {bookings(random, 20000, 20000, 1, 10000), 20000}},
		))
	},
	"cf-816-b": func(random func() float64) []string {
		// A count over each temperature of each question makes q * T = 8 * 10^9 steps on the first largest input. 4 * 10^4 answers of at most 7 characters keep the answer below 10^6 characters.
		// Recipes of at most 40 temperatures cover a temperature about 5 times, so k = 5 makes about half of the temperatures admissible.
		return slices.Concat(
			[]string{"1 1 1\n1 1\n1 1\n", "1 1 1\n1 200000\n1 200000\n", "2 2 1\n1 1\n2 2\n1 2\n", "3 1 2\n5 5\n5 5\n5 5\n5 5\n4 6\n"},
			from(4, func(int) string {
				n := between(random, 1, 30)
				return coffee(random, n, between(random, 1, n), between(random, 1, 30), 100, 100, false)
			}),
			[]string{coffee(random, 50000, 2, 40000, 200000, 200000, true), coffee(random, 50000, 1, 40000, 200000, 200000, false), coffee(random, 50000, 5, 40000, 200000, 40, false)},
		)
	},
	"cf-295-a": func(random func() float64) []string {
		// A loop over the operations of each query makes k * m = 9 * 10^8 steps on the second largest input. The values of the first largest input are below 10^14, so its answer has at most 15 characters for each of 3 * 10^4 numbers.
		return slices.Concat(
			[]string{"1 1 1\n0\n1 1 0\n1 1\n", "1 1 1\n100000\n1 1 100000\n1 1\n", "3 2 1\n0 0 0\n1 3 5\n2 2 7\n1 2\n"},
			from(4, func(i int) string {
				hi := []int{10, 1e5}[i%2]
				n, m := between(random, 1, 30), between(random, 1, 30)
				return operationsInput(random, list(random, n, 0, hi), m, between(random, 1, 30), n, hi, m, false)
			}),
			[]string{
				operationsInput(random, list(random, 30000, 0, 1e5), 30000, 30000, 30000, 1e5, 30000, false),
				operationsInput(random, list(random, 30000, 0, 1e5), 30000, 30000, 10, 1e5, 0, true),
				operationsInput(random, list(random, 30000, 0, 1e5), 30000, 30000, 10, 9, 10, false),
			},
		)
	},
	"image-smoother": func(random func() float64) []string {
		// A window has at most 9 cells, so a direct sum over the neighbours is O(mn) too and passes. This is the easy task of the chapter.
		return call(slices.Concat(
			[][]any{{[][]int{{0}}}, {[][]int{{255}}}, {filled(1, 200, 255)}, {matrix(random, 200, 1, 0, 255)}, {[][]int{{0, 255}}}, {filled(200, 200, 255)}},
			from(4, func(int) []any { return []any{matrix(random, between(random, 1, 20), between(random, 1, 20), 0, 255)} }),
			[][]any{{matrix(random, 200, 200, 0, 255)}, {matrix(random, 200, 1, 0, 255)}},
		))
	},
	"matrix-block-sum": func(random func() float64) []string {
		// The limits of the statement let the Brute Force over each window pass: 10^8 additions at 100 x 100 and k = 100, about 0.35 s with slice sums in CPython.
		return call(slices.Concat(
			[][]any{{[][]int{{1}}, 1}, {[][]int{{100}}, 100}, {filled(1, 100, 100), 1}, {matrix(random, 100, 1, 1, 100), 50}, {filled(100, 100, 1), 1}},
			from(4, func(int) []any {
				return []any{matrix(random, between(random, 1, 20), between(random, 1, 20), 1, 100), between(random, 1, 5)}
			}),
			[][]any{{matrix(random, 100, 100, 1, 100), 100}, {filled(100, 100, 100), 50}, {matrix(random, 100, 100, 1, 100), 1}},
		))
	},
	"increment-submatrices-by-one": func(random func() float64) []string {
		// A loop over the cells of each query makes q * n^2 = 9 * 10^8 additions on the second largest input. Its answer has 90 000 values 10000, about 630 000 characters.
		return call(slices.Concat(
			[][]any{{1, [][]int{{0, 0, 0, 0}}}, {500, [][]int{{0, 0, 499, 499}}}, {2, [][]int{{0, 1, 1, 1}, {1, 0, 1, 1}}}, {5, slices.Repeat([][]int{{2, 2, 2, 2}}, 10000)}},
			from(4, func(int) []any {
				n := between(random, 1, 20)
				return []any{n, submatrices(random, between(random, 1, 30), n, n)}
			}),
			[][]any{{350, submatrices(random, 10000, 350, 350)}, {300, slices.Repeat([][]int{{0, 0, 299, 299}}, 10000)}, {350, submatrices(random, 10000, 350, 1)}},
		))
	},
	"cf-1722-e": func(random func() float64) []string {
		// A check of every rectangle for every query makes n * q = 6 * 10^8 steps on the first two largest inputs. 2 * 10^4 answers of at most 11 digits keep the answer below 10^6 characters.
		// Every input has t <= 3: a Python solution builds a 1001 x 1001 table for each test.
		return slices.Concat(
			[]string{"1\n1 1\n1 1\n1 1 2 2\n", "1\n1 1\n2 2\n1 1 3 3\n", "1\n1 1\n1000 1000\n999 999 1000 1000\n", "1\n2 2\n500 500\n500 501\n499 499 501 501\n499 499 501 502\n"},
			from(4, func(int) string {
				return tests(between(random, 1, 3), func() string { return rectangles(random, between(random, 1, 50), between(random, 1, 50), 20, false) })
			}),
			[]string{tests(1, func() string { return rectangles(random, 30000, 20000, 1000, false) }), tests(1, func() string { return rectangles(random, 30000, 20000, 1000, true) }), tests(2, func() string { return rectangles(random, 10000, 10000, 1000, false) })},
		)
	},
}

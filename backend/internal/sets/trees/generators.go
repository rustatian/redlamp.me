package trees

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// The hidden tests of the set "trees". Each generator gives the edge cases of the statement, then random medium inputs, then the inputs at the maximum sizes.
// A call input is a Python literal list of the arguments. JSON of integers and lowercase letters is a Python literal too.
// A tree shape is a parent list with the root 0: the parent of vertex v is smaller than v.

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

func call(cases [][]any) []string {
	return from(len(cases), func(i int) string { return sets.Call(cases[i]...) })
}

// randomShape gives each vertex v >= 1 a random parent from 0..v - 1.
func randomShape(random func() float64, n int) []int {
	return from(n, func(v int) int {
		if v == 0 {
			return -1
		}
		return between(random, 0, v-1)
	})
}

// pathShape is the path 0 - 1 - ... - (n - 1).
func pathShape(n int) []int {
	return from(n, func(v int) int { return v - 1 })
}

// starShape joins the vertices 1..n - 1 to the center 0.
func starShape(n int) []int {
	return from(n, func(v int) int { return min(v, 1) - 1 })
}

// chainWithLeaves is a path 0 - 1 - ... - (m - 1) of m = ceil(n / 2) vertices with one leaf at each path vertex until n vertices.
func chainWithLeaves(n int) []int {
	m := (n + 1) / 2
	return from(n, func(v int) int {
		if v < m {
			return v - 1
		}
		return v - m
	})
}

// labels gives the vertices 0..n - 1 new numbers from base in a random order (Fisher-Yates). With keepRoot, vertex 0 keeps the number base.
func labels(random func() float64, n, base int, keepRoot bool) []int {
	label := from(n, func(v int) int { return v + base })
	first := 0
	if keepRoot {
		first = 1
	}
	for i := n - 1; i > first; i-- {
		j := between(random, first, i)
		label[i], label[j] = label[j], label[i]
	}
	return label
}

// edgesOf returns the edges [parent, child] of a shape with the numbers of label.
func edgesOf(parent, label []int) [][]int {
	var edges [][]int
	for v, p := range parent {
		if p != -1 {
			edges = append(edges, []int{label[p], label[v]})
		}
	}
	return edges
}

// shuffle puts the edges in a random order. With flip, each edge also turns with the chance of one half.
func shuffle(random func() float64, edges [][]int, flip bool) [][]int {
	for i := len(edges) - 1; i >= 0; i-- {
		j := between(random, 0, i)
		edges[i], edges[j] = edges[j], edges[i]
		if flip && random() < 0.5 {
			edges[i][0], edges[i][1] = edges[i][1], edges[i][0]
		}
	}
	return edges
}

// toward turns each edge [parent, child] into [child, parent].
func toward(edges [][]int) [][]int {
	for _, e := range edges {
		e[0], e[1] = e[1], e[0]
	}
	return edges
}

// tree returns the edges of a shape with random numbers from base, in a random order and orientation. With keepRoot, the root keeps the number base.
func tree(random func() float64, parent []int, base int, keepRoot bool) [][]int {
	return shuffle(random, edgesOf(parent, labels(random, len(parent), base, keepRoot)), true)
}

// edgeLines writes the n - 1 edges of a Codeforces tree after the line with n.
func edgeLines(n int, edges [][]int) string {
	var text strings.Builder
	text.WriteString(strconv.Itoa(n) + "\n")
	for _, e := range edges {
		text.WriteString(strconv.Itoa(e[0]) + " " + strconv.Itoa(e[1]) + "\n")
	}
	return text.String()
}

// star returns the edges of a star with the vertices 1..n and the given center, in a random order and orientation.
func star(random func() float64, n, center int) [][]int {
	var edges [][]int
	for v := 1; v <= n; v++ {
		if v != center {
			edges = append(edges, []int{center, v})
		}
	}
	return shuffle(random, edges, true)
}

// company returns the arguments of LeetCode 1376 for a shape with random numbers: n, the head, the manager list and informTime.
// A manager gets a random informTime from least to most, an employee without subordinates gets 0.
func company(random func() float64, parent []int, least, most int) []any {
	n := len(parent)
	label := labels(random, n, 0, false)
	manager := make([]int, n)
	boss := make([]bool, n)
	for v, p := range parent {
		manager[label[v]] = -1
		if p != -1 {
			manager[label[v]], boss[p] = label[p], true
		}
	}
	inform := make([]int, n)
	for v, b := range boss {
		if b {
			inform[label[v]] = between(random, least, most)
		}
	}
	return []any{n, label[0], manager, inform}
}

// bigNumber is a random whole number in [1, 10^16]. One draw has 32 random bits, so the number takes two draws.
func bigNumber(random func() float64) int {
	return between(random, 0, 1e8-1)*1e8 + between(random, 0, 1e8-1) + 1
}

// heapShape is the complete Binary Tree of n vertices: the parent of vertex v is (v - 1) / 2.
func heapShape(n int) []int {
	return from(n, func(v int) int { return (v+1)/2 - 1 })
}

// randomBinary gives each vertex v >= 1 a random parent among the vertices 0..v - 1 with fewer than 2 children.
func randomBinary(random func() float64, n int) []int {
	parent := make([]int, n)
	parent[0] = -1
	children := make([]int, n)
	free := []int{0}
	for v := 1; v < n; v++ {
		i := between(random, 0, len(free)-1)
		p := free[i]
		parent[v] = p
		if children[p]++; children[p] == 2 {
			free[i] = free[len(free)-1]
			free = free[:len(free)-1]
		}
		free = append(free, v)
	}
	return parent
}

// parentsOf returns the parent list of a shape with the numbers of label: the parent of label[v] is label[parent[v]], and the root has -1.
func parentsOf(parent, label []int) []int {
	parents := make([]int, len(parent))
	for v, p := range parent {
		parents[label[v]] = -1
		if p != -1 {
			parents[label[v]] = label[p]
		}
	}
	return parents
}

// countTests writes a Codeforces input of the tests: the number of tests, then the tests.
func countTests(tests []string) string {
	return strconv.Itoa(len(tests)) + "\n" + strings.Join(tests, "")
}

// colorTest writes one test of Codeforces 1676G: n, the parents a_2..a_n of a shape with the numbers from 1 and the colors of the vertices.
func colorTest(parent []int, colors string) string {
	var text strings.Builder
	text.WriteString(strconv.Itoa(len(parent)) + "\n")
	for v := 1; v < len(parent); v++ {
		if v > 1 {
			text.WriteByte(' ')
		}
		text.WriteString(strconv.Itoa(parent[v] + 1))
	}
	text.WriteString("\n" + colors + "\n")
	return text.String()
}

// appleTest writes one test of Codeforces 1843D: the tree, then q queries. Each vertex of a query is a random vertex of pick.
func appleTest(n int, edges [][]int, q int, pick func() int) string {
	var text strings.Builder
	text.WriteString(edgeLines(n, edges) + strconv.Itoa(q) + "\n")
	for range q {
		text.WriteString(strconv.Itoa(pick()) + " " + strconv.Itoa(pick()) + "\n")
	}
	return text.String()
}

// divisible changes one value so that the sum of the values is divisible by k. Each value stays from 0 to 10^9.
func divisible(random func() float64, values []int, k int) []int {
	total := 0
	for _, x := range values {
		total += x
	}
	r := total % k
	if r == 0 {
		return values
	}
	for i, x := range values {
		if x >= r {
			values[i] -= r
			return values
		}
	}
	// Each value is below r, so the new value is below k <= 10^9.
	values[between(random, 0, len(values)-1)] += k - r
	return values
}

// randomText returns n random letters of alphabet.
func randomText(random func() float64, n int, alphabet string) string {
	return string(from(n, func(int) byte { return alphabet[between(random, 0, len(alphabet)-1)] }))
}

// binaryTest writes one test of Codeforces 1900C: n, the letters and the children "l r" of each vertex. The shape has at most 2 children at each vertex.
// The root keeps the number 1, and the other vertices get random numbers. side(v) is 0 when v is the left child and 1 when v is the right child. A second child goes to the free side.
func binaryTest(random func() float64, parent []int, side func(v int) int, letters string) string {
	n := len(parent)
	label := labels(random, n, 1, true)
	child := make([][2]int, n+1)
	for v, p := range parent {
		if p != -1 {
			s := side(v)
			if child[label[p]][s] != 0 {
				s = 1 - s
			}
			child[label[p]][s] = label[v]
		}
	}
	var text strings.Builder
	text.WriteString(strconv.Itoa(n) + "\n" + letters + "\n")
	for _, c := range child[1:] {
		text.WriteString(strconv.Itoa(c[0]) + " " + strconv.Itoa(c[1]) + "\n")
	}
	return text.String()
}

// segmentTest writes one test of Codeforces 1528A: n, the segment "l r" of each vertex from segment, then the edges.
func segmentTest(edges [][]int, segment func() [2]int) string {
	n := len(edges) + 1
	var text strings.Builder
	text.WriteString(strconv.Itoa(n) + "\n")
	for range n {
		s := segment()
		text.WriteString(strconv.Itoa(s[0]) + " " + strconv.Itoa(s[1]) + "\n")
	}
	for _, e := range edges {
		text.WriteString(strconv.Itoa(e[0]) + " " + strconv.Itoa(e[1]) + "\n")
	}
	return text.String()
}

// numbers writes a Codeforces input of t tests with one number each.
func numbers(values []int) string {
	var text strings.Builder
	text.WriteString(strconv.Itoa(len(values)) + "\n")
	for _, x := range values {
		text.WriteString(strconv.Itoa(x) + "\n")
	}
	return text.String()
}

// Generators are the hidden test generators of the tasks by problem slug.
var Generators = map[string]sets.Generator{
	"find-center-of-star-graph": func(random func() float64) []string {
		return call(slices.Concat(
			[][]any{{[][]int{{2, 1}, {3, 2}}}, {[][]int{{1, 3}, {2, 3}}}, {[][]int{{1, 2}, {3, 1}, {1, 4}}}, {[][]int{{4, 1}, {2, 4}, {4, 3}}}},
			from(4, func(int) []any {
				n := between(random, 5, 1000)
				return []any{star(random, n, between(random, 1, n))}
			}),
			// Counting each vertex over all edges reads n(n - 1) = 2.5 * 10^9 numbers for n = 5 * 10^4.
			[][]any{{star(random, 10000, between(random, 1, 10000))}, {star(random, 50000, 1)}, {star(random, 50000, 50000)}, {star(random, 50000, between(random, 1, 50000))}},
		))
	},
	"cf-1843-c": func(random func() float64) []string {
		// mixed first draws a number of digits to cut, so small and large n meet in one input.
		mixed := func() int { return (bigNumber(random)-1)/int(math.Pow10(between(random, 0, 15))) + 1 }
		return slices.Concat(
			[]string{
				"1\n1\n", "1\n10000000000000000\n",
				numbers(from(54, func(k int) int { return 1 << k })),
				numbers(from(53, func(k int) int { return 1<<(k+1) - 1 })),
				// Above 2^53 a float64 division by 2 rounds: int(n / 2) gives wrong sums.
				numbers([]int{1<<53 + 1, 1<<53 + 3, 9999999999999999}),
				numbers(from(100, func(int) int { return between(random, 1, 1000) })),
				numbers(from(1000, func(int) int { return between(random, 1, 1e9) })),
			},
			[]string{
				numbers(from(10000, func(int) int { return bigNumber(random) })), numbers(from(10000, func(int) int { return bigNumber(random) })),
				numbers(from(10000, func(int) int { return mixed() })), numbers(from(10000, func(int) int { return mixed() })),
			},
		)
	},
	"time-needed-to-inform-all-employees": func(random func() float64) []string {
		// A chain with the head at its end: manager[i] = i + 1.
		chain := from(10, func(i int) int { return i + 1 })
		chain[9] = -1
		return call(slices.Concat(
			[][]any{
				{1, 0, []int{-1}, []int{0}},
				{8, 0, []int{-1, 0, 0, 0, 0, 0, 0, 0}, []int{1000, 0, 0, 0, 0, 0, 0, 0}},
				{8, 7, []int{7, 7, 7, 7, 7, 7, 7, -1}, []int{0, 0, 0, 0, 0, 0, 0, 5}},
				{10, 9, chain, from(10, func(i int) int { return min(i, 1) * between(random, 1, 1000) })},
				company(random, randomShape(random, 10), 1, 9),
			},
			from(4, func(int) []any { return company(random, randomShape(random, between(random, 100, 5000)), 0, 1000) }),
			// On the chain-with-leaves, a walk up from each employee to the head makes about 1.25 * 10^9 steps. The chain has the depth 10^5, too deep for Recursion.
			[][]any{
				{100000, 0, pathShape(100000), from(100000, func(i int) int { return min(99999-i, 1) * 1000 })},
				company(random, chainWithLeaves(100000), 0, 1000), company(random, randomShape(random, 100000), 0, 1000), company(random, randomShape(random, 100000), 1000, 1000),
			},
		))
	},
	"reorder-routes-to-make-all-paths-lead-to-the-city-zero": func(random func() float64) []string {
		// away keeps city 0 and each road from the city nearer to 0.
		away := func(parent []int) [][]int {
			return shuffle(random, edgesOf(parent, labels(random, len(parent), 0, true)), false)
		}
		return call(slices.Concat(
			[][]any{{2, [][]int{{0, 1}}}, {2, [][]int{{1, 0}}}, {10, toward(away(starShape(10)))}, {10, away(starShape(10))}, {10, away(pathShape(10))}},
			from(4, func(int) []any {
				n := between(random, 10, 5000)
				return []any{n, tree(random, randomShape(random, n), 0, false)}
			}),
			// One search for each road costs about 2.5 * 10^9 steps. City 0 is an end of the path, so the path has the depth 5 * 10^4, too deep for Recursion.
			[][]any{
				{50000, tree(random, pathShape(50000), 0, true)}, {50000, tree(random, randomShape(random, 50000), 0, false)},
				{50000, tree(random, chainWithLeaves(50000), 0, true)}, {50000, away(randomShape(random, 50000))},
			},
		))
	},
	"cf-862-b": func(random func() float64) []string {
		input := func(parent []int) string { return edgeLines(len(parent), tree(random, parent, 1, false)) }
		return slices.Concat(
			[]string{"1\n", "2\n2 1\n", input(starShape(10)), input(pathShape(10))},
			from(4, func(int) string { return input(randomShape(random, between(random, 10, 5000))) }),
			// A check of all pairs of vertices makes about 5 * 10^9 steps. The path has the answer 2 499 900 001, more than 32 bits.
			[]string{input(pathShape(100000)), input(randomShape(random, 100000)), input(chainWithLeaves(100000)), input(randomShape(random, 50000))},
		)
	},
	"cf-1676-g": func(random func() float64) []string {
		random4000 := func(shape func(n int) []int, color func(n int) string) string {
			return countTests(from(50, func(int) string { return colorTest(shape(4000), color(4000)) }))
		}
		colorsOf := func(n int) string { return randomText(random, n, "BW") }
		alternating := func(n int) string { return strings.Repeat("BW", n/2) }
		shapes := []func(n int) []int{pathShape, starShape, chainWithLeaves, heapShape, func(n int) []int { return randomShape(random, n) }}
		return slices.Concat(
			[]string{
				countTests([]string{colorTest(pathShape(2), "BB")}), countTests([]string{colorTest(pathShape(2), "BW")}), countTests([]string{colorTest(pathShape(2), "WW")}),
				countTests(from(10000, func(int) string { n := between(random, 2, 20); return colorTest(randomShape(random, n), colorsOf(n)) })),
			},
			from(3, func(int) string {
				return countTests(from(100, func(int) string { n := between(random, 2, 2000); return colorTest(randomShape(random, n), colorsOf(n)) }))
			}),
			// 40 tests of all shapes have at most 1.6 * 10^5 vertices.
			[]string{countTests(from(40, func(int) string {
				n := between(random, 2, 4000)
				return colorTest(shapes[between(random, 0, len(shapes)-1)](n), colorsOf(n))
			}))},
			// A walk up from each vertex to the root makes 8 * 10^6 steps on one path of 4000 vertices and 4 * 10^8 steps on 50 paths. A path is too deep for Recursion.
			[]string{
				random4000(pathShape, alternating), random4000(pathShape, colorsOf),
				random4000(func(n int) []int { return randomShape(random, n) }, colorsOf), random4000(starShape, colorsOf),
			},
		)
	},
	"count-nodes-with-the-highest-score": func(random func() float64) []string {
		input := func(parent []int) []any { return []any{parentsOf(parent, labels(random, len(parent), 0, true))} }
		return call(slices.Concat(
			[][]any{{[]int{-1, 0}}, {[]int{-1, 0, 0}}, input(pathShape(5)), {heapShape(7)}, input(randomBinary(random, 10))},
			from(4, func(int) []any { return input(randomBinary(random, between(random, 100, 5000))) }),
			// Removing each vertex with a new search for the parts makes about 10^10 steps. The path is too deep for Recursion.
			[][]any{input(pathShape(100000)), input(heapShape(100000)), input(randomBinary(random, 100000))},
		))
	},
	"cf-1843-d": func(random func() float64) []string {
		// rooted gives the tree of a shape with the root 1 and random numbers of the other vertices.
		rooted := func(parent []int) [][]int { return tree(random, parent, 1, true) }
		vertex := func(n int) func() int { return func() int { return between(random, 1, n) } }
		return slices.Concat(
			[]string{
				"1\n2\n1 2\n4\n1 1\n1 2\n2 1\n2 2\n",
				"1\n3\n2 1\n3 2\n9\n1 1\n1 2\n1 3\n2 1\n2 2\n2 3\n3 1\n3 2\n3 3\n",
				countTests(from(10000, func(int) string {
					n := between(random, 2, 10)
					return appleTest(n, tree(random, randomShape(random, n), 1, false), between(random, 1, 10), vertex(n))
				})),
			},
			from(3, func(int) string {
				return countTests(from(between(random, 2, 5), func(int) string {
					n := between(random, 10, 5000)
					return appleTest(n, tree(random, randomShape(random, n), 1, false), between(random, 1, 5000), vertex(n))
				}))
			}),
			// A new search for each query makes about 2.5 * 10^9 steps. The path is too deep for Recursion. The star has the answer 49 999^2 = 2 499 900 001 for "1 1".
			[]string{
				countTests([]string{appleTest(50000, rooted(pathShape(50000)), 50000, vertex(50000))}),
				// In the star, half of the query vertices are the center 1.
				countTests([]string{appleTest(50000, rooted(starShape(50000)), 50000, func() int { return max(between(random, -50000, 50000), 1) })}),
				countTests([]string{appleTest(50000, tree(random, randomShape(random, 50000), 1, false), 50000, vertex(50000))}),
				countTests([]string{appleTest(50000, rooted(chainWithLeaves(50000)), 50000, vertex(50000))}),
			},
		)
	},
	"minimum-fuel-cost-to-report-to-the-capital": func(random func() float64) []string {
		// capital gives the roads of a shape with the capital 0 and random numbers of the other cities.
		capital := func(parent []int) [][]int { return tree(random, parent, 0, true) }
		return call(slices.Concat(
			[][]any{
				{[][]int{}, between(random, 1, 100000)}, {[][]int{{1, 0}}, 1}, {capital(starShape(10)), 1}, {capital(starShape(10)), 100000}, {capital(pathShape(10)), 1},
			},
			from(4, func(int) []any {
				n := between(random, 10, 5000)
				return []any{tree(random, randomShape(random, n), 0, false), between(random, 1, n)}
			}),
			// A walk of each representative up to the capital makes about 1.25 * 10^9 steps on the path. The path has the answer 1 249 975 000 and is too deep for Recursion.
			[][]any{
				{capital(pathShape(50000)), 1}, {tree(random, randomShape(random, 50000), 0, false), between(random, 1, 100000)}, {capital(chainWithLeaves(50000)), 2},
			},
		))
	},
	"maximum-number-of-k-divisible-components": func(random func() float64) []string {
		// input gives the arguments for a shape with random numbers and random values from 0 to most.
		input := func(parent []int, most, k int) []any {
			n := len(parent)
			values := from(n, func(int) int { return between(random, 0, most) })
			return []any{n, tree(random, parent, 0, false), divisible(random, values, k), k}
		}
		return call(slices.Concat(
			[][]any{
				{1, [][]int{}, []int{between(random, 0, 1e9)}, 1}, input(randomShape(random, 8), 1e9, 1), input(randomShape(random, 10), 0, between(random, 1, 1e9)),
				input(randomShape(random, 10), 20, between(random, 2, 10)),
			},
			from(4, func(int) []any {
				return input(randomShape(random, between(random, 100, 5000)), 1e9, between(random, 2, 100))
			}),
			// A new search for the sum of each part makes about 9 * 10^8 steps. The path is too deep for Recursion.
			[][]any{
				input(pathShape(30000), 5, 3), input(randomShape(random, 30000), 1e9, 1e9), input(chainWithLeaves(30000), 13, 7), input(randomShape(random, 30000), 0, 1e9),
			},
		))
	},
	"cf-1900-c": func(random func() float64) []string {
		letters := func(n int) string { return randomText(random, n, "ULR") }
		anySide := func(int) int { return between(random, 0, 1) }
		left := func(int) int { return 0 }
		small := func(n int) string { return binaryTest(random, randomBinary(random, n), anySide, letters(n)) }
		return slices.Concat(
			[]string{
				// n = 2: the child on the left and on the right, with each letter of the root.
				countTests(from(6, func(i int) string {
					return binaryTest(random, pathShape(2), func(int) int { return i % 2 }, string("ULR"[i/2])+letters(1))
				})),
				countTests(from(50000, func(int) string { return small(2) })),
				countTests(from(1000, func(int) string { return small(between(random, 2, 10)) })),
			},
			from(3, func(int) string {
				return countTests(from(100, func(int) string { return small(between(random, 2, 3000)) }))
			}),
			// A walk up from each leaf to the root makes about 1.25 * 10^9 steps on the chain with leaves. The chains are too deep for Recursion. The left chain with U has the answer 99 999.
			[]string{
				countTests([]string{binaryTest(random, pathShape(100000), left, strings.Repeat("U", 100000))}),
				countTests([]string{binaryTest(random, pathShape(100000), func(v int) int { return v % 2 }, letters(100000))}),
				countTests([]string{binaryTest(random, chainWithLeaves(100000), left, letters(100000))}),
				countTests([]string{small(100000)}),
				countTests([]string{binaryTest(random, heapShape(100000), anySide, letters(100000))}),
				countTests(from(2, func(int) string { return small(50000) })),
			},
		)
	},
	"maximum-score-after-applying-operations-on-a-tree": func(random func() float64) []string {
		// input gives the arguments for a shape with the root 0 and random values from 1 to most.
		input := func(parent []int, most int) []any {
			return []any{tree(random, parent, 0, true), from(len(parent), func(int) int { return between(random, 1, most) })}
		}
		// starValues gives a star with the value root at the center 0 and random values from 1 to most at the leaves.
		starValues := func(root, most int) []any {
			values := from(10, func(int) int { return between(random, 1, most) })
			values[0] = root
			return []any{tree(random, starShape(10), 0, true), values}
		}
		return call(slices.Concat(
			[][]any{input(pathShape(2), 1e9), starValues(1, 1e9), starValues(1e9, 10), input(pathShape(10), 1e9), input(randomShape(random, 10), 10)},
			from(4, func(int) []any { return input(randomShape(random, between(random, 100, 5000)), 1e9) }),
			// A new search for the kept sum of each subtree makes about 2 * 10^8 steps on the path. The path is too deep for Recursion.
			[][]any{input(pathShape(20000), 1e9), input(randomShape(random, 20000), 1e9), input(chainWithLeaves(20000), 3)},
		))
	},
	"minimum-increments-to-equalize-leaf-paths": func(random func() float64) []string {
		// input gives the arguments for a shape with the root 0 and a cost from cost for each vertex.
		input := func(parent []int, cost func() int) []any {
			n := len(parent)
			return []any{n, tree(random, parent, 0, true), from(n, func(int) int { return cost() })}
		}
		upTo := func(most int) func() int { return func() int { return between(random, 1, most) } }
		equal := between(random, 1, 1e9)
		return call(slices.Concat(
			[][]any{
				input(pathShape(2), upTo(1e9)), input(starShape(10), func() int { return equal }), input(starShape(10), upTo(1e9)), input(pathShape(10), upTo(1e9)),
				input(randomShape(random, 10), upTo(3)),
			},
			from(4, func(i int) []any {
				return input(randomShape(random, between(random, 100, 5000)), upTo([]int{3, 3, 1e9, 1e9}[i]))
			}),
			// A new search for the longest path down from each vertex makes about 1.25 * 10^9 steps on the path. The path is too deep for Recursion.
			[][]any{input(pathShape(50000), upTo(1e9)), input(randomShape(random, 50000), upTo(3)), input(chainWithLeaves(50000), upTo(1e9))},
		))
	},
	"cf-1528-a": func(random func() float64) []string {
		// segment draws a random segment [l, r] with 1 <= l <= r <= most.
		segment := func(most int) func() [2]int {
			return func() [2]int {
				l := between(random, 1, most)
				return [2]int{l, between(random, l, most)}
			}
		}
		point := func() [2]int {
			x := between(random, 1, 1e9)
			return [2]int{x, x}
		}
		same := point()
		randomTree := func(n int) [][]int { return tree(random, randomShape(random, n), 1, false) }
		return slices.Concat(
			[]string{
				countTests([]string{segmentTest(tree(random, pathShape(2), 1, false), func() [2]int { return same })}),
				countTests([]string{segmentTest(tree(random, pathShape(2), 1, false), func() [2]int { return [2]int{1, 1e9} })}),
				countTests([]string{segmentTest(star(random, 4, between(random, 1, 4)), segment(100))}),
				countTests([]string{segmentTest(tree(random, pathShape(3), 1, false), segment(100))}),
				countTests(from(250, func(int) string { return segmentTest(randomTree(between(random, 2, 20)), segment(20)) })),
			},
			from(3, func(int) string {
				return countTests(from(10, func(int) string { return segmentTest(randomTree(between(random, 10, 5000)), segment(1e9)) }))
			}),
			// Both ends of each vertex give 2^n choices. Vertex 1 is an end of the path and the top of the chain, so they are too deep for Recursion.
			[]string{
				countTests([]string{segmentTest(tree(random, pathShape(30000), 1, true), segment(1e9)), segmentTest(randomTree(30000), segment(1e9))}),
				countTests([]string{segmentTest(star(random, 60000, 1), segment(1e9))}),
				countTests([]string{segmentTest(tree(random, chainWithLeaves(60000), 1, true), segment(1e9))}),
				countTests(from(2, func(int) string { return segmentTest(randomTree(30000), point) })),
			},
		)
	},
	"longest-path-with-different-adjacent-characters": func(random func() float64) []string {
		// input gives the arguments for a shape with the root 0: the parent list and the letter of each vertex of the shape.
		input := func(parent []int, letter func(v int) byte) []any {
			label := labels(random, len(parent), 0, true)
			s := make([]byte, len(parent))
			for v := range parent {
				s[label[v]] = letter(v)
			}
			return []any{parentsOf(parent, label), string(s)}
		}
		anyOf := func(alphabet string) func(int) byte {
			return func(int) byte { return randomText(random, 1, alphabet)[0] }
		}
		alternate := func(v int) byte { return "ab"[v%2] }
		return call(slices.Concat(
			[][]any{
				{[]int{-1}, randomText(random, 1, "abcdefghijklmnopqrstuvwxyz")}, input(pathShape(10), func(int) byte { return 'q' }), input(pathShape(10), alternate),
				// The center z and the leaves a and b: the answer is 3.
				input(starShape(10), func(v int) byte { return "zab"[min(v, 1)*between(random, 1, 2)] }),
				input(randomShape(random, 10), anyOf("ab")),
			},
			from(4, func(i int) []any {
				return input(randomShape(random, between(random, 100, 5000)), anyOf([]string{"ab", "ab", "abc", "abc"}[i]))
			}),
			// A search from every vertex makes about 10^10 steps. The path is too deep for Recursion and has the answer 10^5.
			[][]any{
				input(pathShape(100000), alternate), input(randomShape(random, 100000), anyOf("ab")), input(chainWithLeaves(100000), anyOf("abcdefghijklmnopqrstuvwxyz")),
			},
		))
	},
	"cf-1905-b": func(random func() float64) []string {
		// input writes one test of a shape with random numbers from 1.
		input := func(parent []int) string { return edgeLines(len(parent), tree(random, parent, 1, false)) }
		shape := func(n int) []int { return randomShape(random, n) }
		return slices.Concat(
			[]string{
				countTests([]string{input(pathShape(2))}), countTests([]string{edgeLines(10, star(random, 10, between(random, 1, 10)))}), countTests([]string{input(pathShape(10))}),
				countTests(from(10000, func(int) string { return input(shape(between(random, 2, 10))) })),
			},
			from(4, func(int) string {
				return countTests(from(between(random, 2, 5), func(int) string { return input(shape(between(random, 10, 5000))) }))
			}),
			// A check of the paths between all pairs of leaves makes about 5 * 10^9 steps on the star. The star has the answer 50 000, and the path has the answer 1.
			[]string{
				countTests(from(10, func(int) string { return input(shape(10000)) })),
				countTests([]string{edgeLines(100000, star(random, 100000, between(random, 1, 100000)))}), countTests([]string{input(pathShape(100000))}), countTests([]string{input(shape(100000))}),
			},
		)
	},
	"cf-690-c2": func(random func() float64) []string {
		// input writes the first line "n m" with m = n - 1, then the edges.
		input := func(n int, edges [][]int) string { return strconv.Itoa(n) + " " + edgeLines(n-1, edges) }
		shaped := func(parent []int) string { return input(len(parent), tree(random, parent, 1, false)) }
		return slices.Concat(
			[]string{
				shaped(pathShape(2)), input(10, star(random, 10, 1)), input(10, star(random, 10, between(random, 2, 10))), shaped(pathShape(10)),
				"7 6\n2 1\n1 3\n4 2\n2 5\n6 3\n5 7\n", shaped(randomShape(random, 10)),
			},
			from(3, func(int) string { return shaped(randomShape(random, between(random, 100, 5000))) }),
			// BFS from every brain makes about 10^10 steps. The path is too deep for Recursion and has the answer 99 999. The chain with leaves has the answer 50 001.
			[]string{shaped(pathShape(100000)), shaped(randomShape(random, 100000)), shaped(chainWithLeaves(100000))},
		)
	},
	"minimum-height-trees": func(random func() float64) []string {
		input := func(parent []int) []any { return []any{len(parent), tree(random, parent, 0, false)} }
		return call(slices.Concat(
			[][]any{{1, [][]int{}}, input(pathShape(2)), input(starShape(10)), input(pathShape(9)), input(pathShape(10)), input(randomShape(random, 10))},
			from(4, func(int) []any { return input(randomShape(random, between(random, 100, 5000))) }),
			// BFS from every vertex makes about 4 * 10^8 steps. The path is too deep for Recursion.
			[][]any{input(pathShape(20000)), input(randomShape(random, 20000)), input(chainWithLeaves(20000))},
		))
	},
	"most-profitable-path-in-a-tree": func(random func() float64) []string {
		// game gives the arguments for a shape with the root 0, Bob at the vertex at of the shape and random even amounts.
		game := func(parent []int, at int) []any {
			label := labels(random, len(parent), 0, true)
			amount := from(len(parent), func(int) int { return 2 * between(random, -5000, 5000) })
			return []any{shuffle(random, edgesOf(parent, label), true), label[at], amount}
		}
		// The last vertex of a shape is a leaf, because each parent is smaller than its child. On a path and on the chain with leaves it is the deepest vertex.
		leaf := func(parent []int) []any { return game(parent, len(parent)-1) }
		return call(slices.Concat(
			[][]any{leaf(pathShape(2)), leaf(pathShape(5)), leaf(pathShape(6)), game(starShape(10), between(random, 1, 9)), game(randomShape(random, 10), between(random, 1, 9))},
			from(4, func(int) []any {
				n := between(random, 100, 5000)
				return game(randomShape(random, n), between(random, 1, n-1))
			}),
			// A walk from 0 to each leaf makes about 3 * 10^8 steps on the chain with leaves. The path is too deep for Recursion.
			[][]any{leaf(pathShape(50000)), leaf(randomShape(random, 50000)), leaf(chainWithLeaves(50000))},
		))
	},
	"find-minimum-diameter-after-merging-two-trees": func(random func() float64) []string {
		edges := func(parent []int) [][]int { return tree(random, parent, 0, false) }
		randomEdges := func(n int) [][]int { return edges(randomShape(random, n)) }
		return call(slices.Concat(
			[][]any{
				{[][]int{}, [][]int{}}, {[][]int{}, edges(pathShape(10))}, {edges(starShape(between(random, 2, 10))), edges(starShape(between(random, 2, 10)))},
				// The paths of 5 and 7 vertices have the diameters 4 and 6, the paths of 4 and 6 vertices have 3 and 5.
				{edges(pathShape(5)), edges(pathShape(7))}, {edges(pathShape(4)), edges(pathShape(6))}, {randomEdges(10), randomEdges(10)},
			},
			from(3, func(int) []any {
				return []any{randomEdges(between(random, 100, 5000)), randomEdges(between(random, 100, 5000))}
			}),
			// BFS from every vertex to find a center makes about 9 * 10^8 steps. The paths are too deep for Recursion.
			[][]any{{edges(pathShape(30000)), edges(pathShape(30000))}, {randomEdges(30000), randomEdges(30000)}, {edges(pathShape(30000)), edges(starShape(30000))}},
		))
	},
}

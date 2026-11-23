package trees

import (
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// edgePairs checks a Python list of exactly count edges. Each edge is a list of two different whole numbers from lo to hi.
func edgePairs(value any, count, lo, hi int) [][2]int {
	list, ok := value.([]any)
	sets.Need(ok && len(list) == count)
	edges := make([][2]int, count)
	for i, item := range list {
		pair := sets.Ints(item, 2, 2, lo, hi)
		sets.Need(pair[0] != pair[1])
		edges[i] = [2]int{pair[0], pair[1]}
	}
	return edges
}

// adjacency returns the neighbour lists of the vertices 0..n - 1. base is the number of the first vertex in the edges: 0 for LeetCode, 1 for Codeforces.
func adjacency(n int, edges [][2]int, base int) [][]int {
	adj := make([][]int, n)
	for _, e := range edges {
		a, b := e[0]-base, e[1]-base
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	return adj
}

// bfsOrder returns the BFS order from root and the parent of each vertex, with -1 for the root.
// A vertex that the BFS finds again shows a cycle or a repeated edge. With n - 1 edges, a BFS that reaches all n vertices proves a tree.
func bfsOrder(adj [][]int, root int) (order, parent []int) {
	parent = make([]int, len(adj))
	seen := make([]bool, len(adj))
	parent[root], seen[root] = -1, true
	order = append(make([]int, 0, len(adj)), root)
	for i := 0; i < len(order); i++ {
		v := order[i]
		for _, u := range adj[v] {
			if u != parent[v] {
				sets.Need(!seen[u])
				seen[u], parent[u] = true, v
				order = append(order, u)
			}
		}
	}
	sets.Need(len(order) == len(adj))
	return order, parent
}

// farthest returns a vertex at the largest distance from s, the distances from s and the parents of the BFS from s. The last vertex of the BFS order is a farthest vertex.
func farthest(adj [][]int, s int) (vertex int, dist, parent []int) {
	order, parent := bfsOrder(adj, s)
	dist = make([]int, len(adj))
	for _, v := range order[1:] {
		dist[v] = dist[parent[v]] + 1
	}
	return order[len(order)-1], dist, parent
}

// diameter returns the number of edges of the longest path of a tree. The farthest vertex a from vertex 0 is an end of such a path, and the farthest vertex from a is its other end.
func diameter(adj [][]int) int {
	a, _, _ := farthest(adj, 0)
	b, dist, _ := farthest(adj, a)
	return dist[b]
}

// readTree reads the n - 1 edges "u v" of a Codeforces tree with the vertices 1..n and returns the neighbour lists of the vertices 0..n - 1.
func readTree(input *sets.Reader, n int) [][]int {
	edges := make([][2]int, n-1)
	for i := range edges {
		edges[i] = [2]int{input.Int(1, n), input.Int(1, n)}
		sets.Need(edges[i][0] != edges[i][1])
	}
	return adjacency(n, edges, 1)
}

// childrenOf returns the children lists of a parent list and an order with each parent before its children.
// The parent of root is -1, and the parent of each other vertex is a vertex. All vertices must be in the tree of root.
func childrenOf(parent []int, root int) (children [][]int, order []int) {
	children = make([][]int, len(parent))
	for v, p := range parent {
		sets.Need((v == root) == (p == -1))
		if p != -1 {
			children[p] = append(children[p], v)
		}
	}
	for stack := []int{root}; len(stack) > 0; {
		v := stack[len(stack)-1]
		stack = append(stack[:len(stack)-1], children[v]...)
		order = append(order, v)
	}
	sets.Need(len(order) == len(parent))
	return children, order
}

// decimal is a whole number from 1 with at most 17 digits. sets.Reader.Int reads through float64, which is not exact above 2^53.
var decimal = regexp.MustCompile(`^[1-9][0-9]{0,16}$`)

// colors are the colors of the vertices of Codeforces 1676G: B for black and W for white.
var colors = regexp.MustCompile(`^[BW]+$`)

// moves are the letters of the vertices of Codeforces 1900C: U for the parent, L for the left child and R for the right child.
var moves = regexp.MustCompile(`^[ULR]+$`)

// lowercase is a string of lowercase English letters.
var lowercase = regexp.MustCompile(`^[a-z]+$`)

// abs returns the absolute value of x.
func abs(x int) int {
	return max(x, -x)
}

// References are the reference solutions of the tasks by problem slug.
var References = map[string]sets.Reference{
	// O(n) for the checks of a star, O(1) for the answer: the center is in every edge, so it is the common vertex of the first two edges.
	"find-center-of-star-graph": {Arity: 1, Call: func(args []any) any {
		list, ok := args[0].([]any)
		sets.Need(ok && len(list) >= 2 && len(list) < 1e5)
		n := len(list) + 1
		edges := edgePairs(list, n-1, 1, n)
		degree := make([]int, n+1)
		for _, e := range edges {
			degree[e[0]]++
			degree[e[1]]++
		}
		centers := 0
		for _, d := range degree[1:] {
			if d == n-1 {
				centers++
			} else {
				sets.Need(d == 1)
			}
		}
		sets.Need(centers == 1)
		if first, second := edges[0], edges[1]; first[0] == second[0] || first[0] == second[1] {
			return sets.Py(first[0])
		}
		return sets.Py(edges[0][1])
	}},
	// O(log n) for each test: the parent of vertex v is v / 2, so the path to the root is n, n / 2, ..., 1. The sum is below 2 * 10^16 and fits in int64.
	"cf-1843-c": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		for t := input.Int(1, 1e4); t > 0; t-- {
			n, err := strconv.Atoi(input.Word(decimal))
			sets.Need(err == nil && n <= 1e16)
			sum := 0
			for ; n > 0; n /= 2 {
				sum += n
			}
			out.WriteString(strconv.Itoa(sum) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(n): the children lists of manager, then the order from headID. A subordinate learns the news at the time of the manager plus informTime of the manager.
	"time-needed-to-inform-all-employees": {Arity: 4, Call: func(args []any) any {
		n := sets.Int(args[0], 1, 1e5)
		head := sets.Int(args[1], 0, n-1)
		manager := sets.Ints(args[2], n, n, -1, n-1)
		inform := sets.Ints(args[3], n, n, 0, 1000)
		children, order := childrenOf(manager, head)
		time := make([]int, n)
		best := 0
		for _, v := range order {
			sets.Need(len(children[v]) > 0 || inform[v] == 0)
			for _, c := range children[v] {
				time[c] = time[v] + inform[v]
			}
			best = max(best, time[v])
		}
		return sets.Py(best)
	}},
	// O(n): BFS from city 0 gives the parent of each city. A road a -> b leads away from 0 when a is the parent of b, and only such roads change.
	"reorder-routes-to-make-all-paths-lead-to-the-city-zero": {Arity: 2, Call: func(args []any) any {
		n := sets.Int(args[0], 2, 5e4)
		roads := edgePairs(args[1], n-1, 0, n-1)
		_, parent := bfsOrder(adjacency(n, roads, 0), 0)
		count := 0
		for _, road := range roads {
			if parent[road[1]] == road[0] {
				count++
			}
		}
		return sets.Py(count)
	}},
	// O(n): BFS colours the vertices by the parity of the depth, a vertices of one colour and n - a of the other. A new edge can join only two colours: a(n - a) pairs, n - 1 of them are edges.
	"cf-862-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n := input.Int(1, 1e5)
		order, parent := bfsOrder(readTree(input, n), 0)
		input.End()
		odd := make([]bool, n)
		even := 0
		for _, v := range order {
			if parent[v] != -1 {
				odd[v] = !odd[parent[v]]
			}
			if !odd[v] {
				even++
			}
		}
		return strconv.Itoa(even*(n-even)-(n-1)) + "\n"
	}},
	// O(n) for each test: a_i < i, so the numbers n, n - 1, ..., 2 give each child before its parent. The balance of a subtree is the number of W minus the number of B, and a balanced subtree has the balance 0.
	"cf-1676-g": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		vertices := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(2, 4000)
			vertices += n
			sets.Need(vertices <= 2e5)
			parent := make([]int, n+1)
			for i := 2; i <= n; i++ {
				parent[i] = input.Int(1, i-1)
			}
			s := input.Word(colors)
			sets.Need(len(s) == n)
			balance := make([]int, n+1)
			for i := 1; i <= n; i++ {
				balance[i] = -1
				if s[i-1] == 'W' {
					balance[i] = 1
				}
			}
			for i := n; i >= 2; i-- {
				balance[parent[i]] += balance[i]
			}
			count := 0
			for _, b := range balance[1:] {
				if b == 0 {
					count++
				}
			}
			out.WriteString(strconv.Itoa(count) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(n): the subtree sizes from the end of an order with each parent before its children. Without vertex v, the parts are the subtrees of its children and the n - size[v] vertices above v.
	// A score is at most (n / 3)^3 < 4 * 10^13 and fits in int64.
	"count-nodes-with-the-highest-score": {Arity: 1, Call: func(args []any) any {
		parents := sets.Ints(args[0], 2, 1e5, -1, 1e5-1)
		n := len(parents)
		for _, p := range parents {
			sets.Need(p < n)
		}
		children, order := childrenOf(parents, 0)
		size := make([]int, n)
		for i := n - 1; i >= 0; i-- {
			v := order[i]
			sets.Need(len(children[v]) <= 2)
			size[v] = 1
			for _, c := range children[v] {
				size[v] += size[c]
			}
		}
		best, count := 0, 0
		for v := range n {
			score := max(n-size[v], 1)
			for _, c := range children[v] {
				score *= size[c]
			}
			if score > best {
				best, count = score, 1
			} else if score == best {
				count++
			}
		}
		return sets.Py(count)
	}},
	// O(n + q) for each test: BFS from vertex 1, then the leaf counts from the end of the order. A vertex without children is a leaf, and the answer of a query is leaves[x] * leaves[y] < 4 * 10^10.
	"cf-1843-d": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		vertices, queries := 0, 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(2, 2e5)
			vertices += n
			sets.Need(vertices <= 2e5)
			order, parent := bfsOrder(readTree(input, n), 0)
			leaves := make([]int, n)
			for i := n - 1; i > 0; i-- {
				v := order[i]
				leaves[v] = max(leaves[v], 1)
				leaves[parent[v]] += leaves[v]
			}
			q := input.Int(1, 2e5)
			queries += q
			sets.Need(queries <= 2e5)
			for ; q > 0; q-- {
				x, y := input.Int(1, n), input.Int(1, n)
				out.WriteString(strconv.Itoa(leaves[x-1]*leaves[y-1]) + "\n")
			}
		}
		input.End()
		return out.String()
	}},
	// O(n): BFS from the capital 0, then the subtree sizes from the end of the order. size[v] representatives go over the road above v and need ceil(size[v] / seats) cars.
	"minimum-fuel-cost-to-report-to-the-capital": {Arity: 2, Call: func(args []any) any {
		list, ok := args[0].([]any)
		sets.Need(ok && len(list) < 1e5)
		n := len(list) + 1
		roads := edgePairs(list, n-1, 0, n-1)
		seats := sets.Int(args[1], 1, 1e5)
		order, parent := bfsOrder(adjacency(n, roads, 0), 0)
		size := make([]int, n)
		fuel := 0
		for i := n - 1; i > 0; i-- {
			v := order[i]
			size[v]++
			fuel += (size[v] + seats - 1) / seats
			size[parent[v]] += size[v]
		}
		return sets.Py(fuel)
	}},
	// O(n): BFS from vertex 0, then the subtree sums from the end of the order. The edge above v can go when the sum of the subtree of v is divisible by k: the total is divisible by k, so the rest is divisible too.
	// Each such vertex, the root too, gives one part. A sum is at most 3 * 10^13 and fits in int64.
	"maximum-number-of-k-divisible-components": {Arity: 4, Call: func(args []any) any {
		n := sets.Int(args[0], 1, 3e4)
		edges := edgePairs(args[1], n-1, 0, n-1)
		sum := sets.Ints(args[2], n, n, 0, 1e9)
		k := sets.Int(args[3], 1, 1e9)
		order, parent := bfsOrder(adjacency(n, edges, 0), 0)
		parts := 0
		for i := n - 1; i >= 0; i-- {
			v := order[i]
			if sum[v]%k == 0 {
				parts++
			}
			if parent[v] != -1 {
				sum[parent[v]] += sum[v]
			}
		}
		sets.Need(sum[0]%k == 0)
		return sets.Py(parts)
	}},
	// O(n) for each test: dp[v] is the smallest number of changed letters on the way from v to a leaf. It is 0 at a leaf, else the minimum of dp[c] over the children c, plus 1 when the letter of v does not lead to c.
	// The children are checked before the BFS from vertex 1: vertex 1 is no child, and each other vertex is a child at most once. Then the BFS cannot visit a vertex two times.
	"cf-1900-c": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		vertices := 0
		for t := input.Int(1, 5e4); t > 0; t-- {
			n := input.Int(2, 3e5)
			vertices += n
			sets.Need(vertices <= 3e5)
			s := input.Word(moves)
			sets.Need(len(s) == n)
			child := make([][2]int, n+1)
			isChild := make([]bool, n+1)
			for v := 1; v <= n; v++ {
				for side := range child[v] {
					c := input.Int(0, n)
					if c != 0 {
						sets.Need(c != 1 && !isChild[c])
						isChild[c] = true
					}
					child[v][side] = c
				}
			}
			order := append(make([]int, 0, n), 1)
			for i := 0; i < len(order); i++ {
				for _, c := range child[order[i]] {
					if c != 0 {
						order = append(order, c)
					}
				}
			}
			sets.Need(len(order) == n)
			dp := make([]int, n+1)
			for i := n - 1; i >= 0; i-- {
				v := order[i]
				if child[v] == [2]int{} {
					continue
				}
				dp[v] = n
				for side, c := range child[v] {
					if c != 0 {
						change := 0
						if s[v-1] != "LR"[side] {
							change = 1
						}
						dp[v] = min(dp[v], dp[c]+change)
					}
				}
			}
			out.WriteString(strconv.Itoa(dp[1]) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(n): BFS from vertex 0, then from the end of the order. keep[v] is the smallest sum of the values that stay in the subtree of v, so that each path from v down to a leaf keeps a value.
	// keep[v] is values[v] at a leaf, else min(values[v], the sum of keep of the children). The answer is the total minus keep[0], at most 2 * 10^13.
	"maximum-score-after-applying-operations-on-a-tree": {Arity: 2, Call: func(args []any) any {
		list, ok := args[0].([]any)
		sets.Need(ok && len(list) >= 1 && len(list) < 2e4)
		n := len(list) + 1
		edges := edgePairs(list, n-1, 0, n-1)
		values := sets.Ints(args[1], n, n, 1, 1e9)
		order, parent := bfsOrder(adjacency(n, edges, 0), 0)
		below := make([]int, n)
		total, keep := 0, 0
		for i := n - 1; i >= 0; i-- {
			v := order[i]
			total += values[v]
			keep = values[v]
			if below[v] > 0 {
				keep = min(keep, below[v])
			}
			if parent[v] != -1 {
				below[parent[v]] += keep
			}
		}
		return sets.Py(total - keep)
	}},
	// O(n): BFS from vertex 0, then from the end of the order. best[v] is the largest sum of cost on a path from v down to a leaf, at most 10^14.
	// Each child c with best[c] below the largest best of its siblings needs one increase, and the children with the largest best need none.
	"minimum-increments-to-equalize-leaf-paths": {Arity: 3, Call: func(args []any) any {
		n := sets.Int(args[0], 2, 1e5)
		edges := edgePairs(args[1], n-1, 0, n-1)
		cost := sets.Ints(args[2], n, n, 1, 1e9)
		order, parent := bfsOrder(adjacency(n, edges, 0), 0)
		best := make([]int, n)
		top := make([]int, n)
		for i := n - 1; i >= 0; i-- {
			v := order[i]
			best[v] = cost[v] + top[v]
			if p := parent[v]; p != -1 {
				top[p] = max(top[p], best[v])
			}
		}
		count := 0
		for v, p := range parent {
			if p != -1 && best[v] < top[p] {
				count++
			}
		}
		return sets.Py(count)
	}},
	// O(n) for each test: some best choice has each a_v equal to l_v or r_v. dp[v][0] and dp[v][1] are the best sums in the subtree of v for a_v = l_v and for a_v = r_v.
	// From the end of the BFS order, a child c adds max(dp[c][0] + |x - l_c|, dp[c][1] + |x - r_c|) for the value x of its parent. A beauty is below 10^5 * 10^9 = 10^14.
	"cf-1528-a": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		vertices := 0
		for t := input.Int(1, 250); t > 0; t-- {
			n := input.Int(2, 1e5)
			vertices += n
			sets.Need(vertices <= 2e5)
			ends := make([][2]int, n)
			for v := range ends {
				ends[v][0] = input.Int(1, 1e9)
				ends[v][1] = input.Int(ends[v][0], 1e9)
			}
			order, parent := bfsOrder(readTree(input, n), 0)
			dp := make([][2]int, n)
			for i := n - 1; i > 0; i-- {
				v := order[i]
				p := parent[v]
				for j, x := range ends[p] {
					dp[p][j] += max(dp[v][0]+abs(x-ends[v][0]), dp[v][1]+abs(x-ends[v][1]))
				}
			}
			out.WriteString(strconv.Itoa(max(dp[0][0], dp[0][1])) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(n): the children lists, then from the end of the order. down[v] is the number of vertices of the longest chain from v down with different letters on each edge.
	// The longest path through v is 1 + the two largest down of the children with another letter: top[v][0] >= top[v][1], 0 when no such child is there.
	"longest-path-with-different-adjacent-characters": {Arity: 2, Call: func(args []any) any {
		parents := sets.Ints(args[0], 1, 1e5, -1, 1e5-1)
		n := len(parents)
		for _, p := range parents {
			sets.Need(p < n)
		}
		s := sets.Str(args[1], lowercase)
		sets.Need(len(s) == n)
		_, order := childrenOf(parents, 0)
		top := make([][2]int, n)
		best := 0
		for i := n - 1; i >= 0; i-- {
			v := order[i]
			best = max(best, 1+top[v][0]+top[v][1])
			if p := parents[v]; p != -1 && s[p] != s[v] {
				if down := 1 + top[v][0]; down > top[p][0] {
					top[p] = [2]int{down, top[p][0]}
				} else {
					top[p][1] = max(top[p][1], down)
				}
			}
		}
		return sets.Py(best)
	}},
	// O(n) for each test: a leaf is a vertex of degree 1. One operation changes only the ends of its path, so it removes at most 2 leaves. L leaves need ceil(L / 2) operations, and this number is enough.
	"cf-1905-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out strings.Builder
		vertices := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(2, 1e5)
			vertices += n
			sets.Need(vertices <= 1e5)
			adj := readTree(input, n)
			bfsOrder(adj, 0)
			leaves := 0
			for _, list := range adj {
				if len(list) == 1 {
					leaves++
				}
			}
			out.WriteString(strconv.Itoa((leaves+1)/2) + "\n")
		}
		input.End()
		return out.String()
	}},
	// O(n): the brain latency is the diameter of the tree, two BFS. With m = n - 1 connectors, a connected Graph is a tree, so the conditions (1) and (2) are true.
	"cf-690-c2": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n, m := input.Int(1, 1e5), input.Int(1, 1e5)
		sets.Need(m == n-1)
		adj := readTree(input, n)
		input.End()
		return strconv.Itoa(diameter(adj)) + "\n"
	}},
	// O(n): two BFS find the ends a and b of a diameter of d edges. The parents of the BFS from a lead from b to a. The center is the vertex d / 2 steps up from b, and for an odd d also its parent.
	"minimum-height-trees": {Arity: 2, Call: func(args []any) any {
		n := sets.Int(args[0], 1, 2e4)
		adj := adjacency(n, edgePairs(args[1], n-1, 0, n-1), 0)
		a, _, _ := farthest(adj, 0)
		b, dist, parent := farthest(adj, a)
		center := b
		for range dist[b] / 2 {
			center = parent[center]
		}
		centers := []int{center}
		if dist[b]%2 == 1 {
			centers = append(centers, parent[center])
		}
		return sets.Py(centers)
	}},
	// O(n): BFS from 0 gives the parent and the depth of each vertex, and Alice is at vertex v at the time depth[v]. Bob walks up from bob and is at the i-th vertex of his walk at the time i.
	// Alice gets amount[v] when she comes first, amount[v] / 2 at the same time and 0 later. Her income adds up from 0 down, and the answer is the best leaf v != 0.
	"most-profitable-path-in-a-tree": {Arity: 3, Call: func(args []any) any {
		list, ok := args[0].([]any)
		sets.Need(ok && len(list) >= 1 && len(list) < 1e5)
		n := len(list) + 1
		adj := adjacency(n, edgePairs(list, n-1, 0, n-1), 0)
		bob := sets.Int(args[1], 1, n-1)
		amount := sets.Ints(args[2], n, n, -1e4, 1e4)
		for _, a := range amount {
			sets.Need(a%2 == 0)
		}
		order, parent := bfsOrder(adj, 0)
		// Alice is at each vertex before the time n, so she comes first where Bob does not walk.
		bobTime := slices.Repeat([]int{n}, n)
		for v, t := bob, 0; v != -1; v, t = parent[v], t+1 {
			bobTime[v] = t
		}
		depth := make([]int, n)
		income := make([]int, n)
		best := math.MinInt
		for _, v := range order {
			if p := parent[v]; p != -1 {
				depth[v], income[v] = depth[p]+1, income[p]
			}
			switch {
			case depth[v] < bobTime[v]:
				income[v] += amount[v]
			case depth[v] == bobTime[v]:
				income[v] += amount[v] / 2
			}
			if v != 0 && len(adj[v]) == 1 {
				best = max(best, income[v])
			}
		}
		return sets.Py(best)
	}},
	// O(n + m): connect the centers of the two trees. The new diameter is max(d1, d2, ceil(d1 / 2) + ceil(d2 / 2) + 1), and no other edge gives a smaller diameter.
	"find-minimum-diameter-after-merging-two-trees": {Arity: 2, Call: func(args []any) any {
		var d [2]int
		for i, arg := range args {
			list, ok := arg.([]any)
			sets.Need(ok && len(list) < 1e5)
			n := len(list) + 1
			d[i] = diameter(adjacency(n, edgePairs(list, n-1, 0, n-1), 0))
		}
		return sets.Py(max(d[0], d[1], (d[0]+1)/2+(d[1]+1)/2+1))
	}},
}

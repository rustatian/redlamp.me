package modular

import (
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// mod is the prime modulus 10^9 + 7 of the answers.
const mod = 1_000_000_007

// remainder returns the remainder of x divided by m, from 0 to m - 1, as x % m in Python. The Go % keeps the sign of x.
func remainder(x, m int) int { return (x%m + m) % m }

// powMod returns a^b % m for a, b >= 0 by Binary Exponentiation: O(log b) multiplications. The remainders are below m <= 3 · 10^9, so a product of two of them fits in int.
func powMod(a, b, m int) int {
	r := 1 % m
	a %= m
	for ; b > 0; b /= 2 {
		if b%2 == 1 {
			r = r * a % m
		}
		a = a * a % m
	}
	return r
}

// factorials returns the tables fact[i] = i! % m and invFact[i] = (i!)^-1 % m for i from 0 to n, for a prime m > n. One powMod gives invFact[n], and invFact[i - 1] = invFact[i] · i because 1 / (i - 1)! = i / i!.
func factorials(n, m int) (fact, invFact []int) {
	fact = make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * i % m
	}
	invFact = make([]int, n+1)
	invFact[n] = powMod(fact[n], m-2, m)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * i % m
	}
	return fact, invFact
}

// count returns the count of the values of a that equal x.
func count(a []int, x int) int {
	n := 0
	for _, v := range a {
		if v == x {
			n++
		}
	}
	return n
}

// anagramText matches the s of LeetCode 2514: words of lower case letters with one space between them.
var anagramText = regexp.MustCompile(`^[a-z]+(?: [a-z]+)*$`)

// References are the reference solutions of the tasks by problem slug.
var References = map[string]sets.Reference{
	// O(n): keep only r = x % 5. A new bit makes x = 2x + bit, so r becomes (2r + bit) % 5.
	"binary-prefix-divisible-by-5": {Arity: 1, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 1e5, 0, 1)
		out := make([]bool, len(nums))
		r := 0
		for i, bit := range nums {
			r = (2*r + bit) % 5
			out[i] = r == 0
		}
		return sets.Py(out)
	}},
	// O(n) for each test: count the classes c0, c1 and c2. A move takes a number to the next class, so the excess of a class goes to the next class. Two rounds over the classes are enough.
	"cf-1490-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(3, 3e4)
			sets.Need(n%3 == 0)
			total += n
			sets.Need(total <= 150000)
			var c [3]int
			for range n {
				c[input.Int(0, 100)%3]++
			}
			moves := 0
			for range 2 {
				for i := range c {
					if extra := c[i] - n/3; extra > 0 {
						c[i] -= extra
						c[(i+1)%3] += extra
						moves += extra
					}
				}
			}
			out = append(out, strconv.Itoa(moves))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n + k): count the remainders. A number of the class r needs a partner of the class (k - r) % k, so cnt[0] is even and cnt[r] == cnt[k - r]. Then the class k / 2 is even too, because n is even.
	"check-if-array-pairs-are-divisible-by-k": {Arity: 2, Call: func(args []any) any {
		arr := sets.Ints(args[0], 1, 1e5, -1e9, 1e9)
		sets.Need(len(arr)%2 == 0)
		k := sets.Int(args[1], 1, 1e5)
		cnt := make([]int, k)
		for _, x := range arr {
			cnt[remainder(x, k)]++
		}
		ok := cnt[0]%2 == 0
		for r := 1; r < k; r++ {
			ok = ok && cnt[r] == cnt[k-r]
		}
		return sets.Py(ok)
	}},
	// O(n + k): the sum of nums[l..r] divides by k when the prefix sums p[l] and p[r + 1] have equal remainders. cnt[0] = 1 counts the empty prefix p[0] = 0.
	"subarray-sums-divisible-by-k": {Arity: 2, Call: func(args []any) any {
		nums := sets.Ints(args[0], 1, 3e4, -1e4, 1e4)
		k := sets.Int(args[1], 2, 1e4)
		cnt := make([]int, k)
		cnt[0] = 1
		p, count := 0, 0
		for _, x := range nums {
			p = remainder(p+x, k)
			count += cnt[p]
			cnt[p]++
		}
		return sets.Py(count)
	}},
	// O(n) expected for each test: a pair fits when a_i % x == (-a_j) % x and a_i % y == a_j % y. A map counts the keys (a % x, a % y) of the numbers on the left.
	"cf-1931-d": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n, x, y := input.Int(2, 2e5), input.Int(1, 1e9), input.Int(1, 1e9)
			total += n
			sets.Need(total <= 2e5)
			// Both remainders are below 2^30, so one int holds the pair.
			cnt := make(map[int]int, n)
			pairs := 0
			for range n {
				a := input.Int(1, 1e9)
				pairs += cnt[(x-a%x)%x<<30|a%y]
				cnt[a%x<<30|a%y]++
			}
			out = append(out, strconv.Itoa(pairs))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(log k) for each test: the base-n digits of the k-th special number are the binary digits of k. Bit i of k adds n^i.
	"cf-1594-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		for t := input.Int(1, 1e4); t > 0; t-- {
			n, k := input.Int(2, 1e9), input.Int(1, 1e9)
			sum := 0
			for p := 1; k > 0; k /= 2 {
				if k%2 == 1 {
					sum = (sum + p) % mod
				}
				p = p * n % mod
			}
			out = append(out, strconv.Itoa(sum))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(q log B) for q items and the largest exponent B: an index is good when pow(pow(a, b, 10), c, m) == target.
	"double-modular-exponentiation": {Arity: 2, Call: func(args []any) any {
		variables, ok := args[0].([]any)
		sets.Need(ok && len(variables) >= 1 && len(variables) <= 100)
		target := sets.Int(args[1], 0, 1e3)
		good := []int{}
		for i, item := range variables {
			v := sets.Ints(item, 4, 4, 1, 1e3)
			if powMod(powMod(v[0], v[1], 10), v[2], v[3]) == target {
				good = append(good, i)
			}
		}
		return sets.Py(good)
	}},
	// O(log n): 5 even digits for each of the (n + 1) / 2 even positions and 4 prime digits for each of the n / 2 odd positions. Product Rule.
	"count-good-numbers": {Arity: 1, Call: func(args []any) any {
		n := sets.Int(args[0], 1, 1e15)
		return sets.Py(powMod(5, (n+1)/2, mod) * powMod(4, n/2, mod) % mod)
	}},
	// O(L) for L digits: when b = 10x + d, a^b = (a^x)^10 · a^d. Go over the digits from the left.
	"super-pow": {Arity: 2, Call: func(args []any) any {
		a := sets.Int(args[0], 1, math.MaxInt32)
		b := sets.Ints(args[1], 1, 2000, 0, 9)
		sets.Need(b[0] != 0)
		r := 1
		for _, d := range b {
			r = powMod(r, 10, 1337) * powMod(a, d, 1337) % 1337
		}
		return sets.Py(r)
	}},
	// O(log n + log m): each kind of present goes into a non-empty Subset of the m boxes, 2^m - 1 ways. The kinds are independent, so the answer is (2^m - 1)^n.
	"cf-1236-b": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		n, m := input.Int(1, 1e9), input.Int(1, 1e9)
		input.End()
		return strconv.Itoa(powMod((powMod(2, m, mod)+mod-1)%mod, n, mod)) + "\n"
	}},
	// O(n√n + n): the p primes from 1 to n take the p prime places in any order, and the other n - p numbers take the other places. The answer is p! · (n - p)!.
	"prime-arrangements": {Arity: 1, Call: func(args []any) any {
		n := sets.Int(args[0], 1, 100)
		primes := 0
		for x := 2; x <= n; x++ {
			prime := true
			for d := 2; d*d <= x; d++ {
				prime = prime && x%d != 0
			}
			if prime {
				primes++
			}
		}
		fact, _ := factorials(n, mod)
		return sets.Py(fact[primes] * fact[n-primes] % mod)
	}},
	// O(n + log M) for each test: the sum of the products of the pairs is (S^2 - Q) / 2 for the sum S and the sum of squares Q, and there are n(n - 1) / 2 pairs. So the answer is (S^2 - Q) · (n(n - 1))^-1.
	"cf-2008-f": {Stdio: func(stdin string) string {
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(2, 2e5)
			total += n
			sets.Need(total <= 2e5)
			sum, squares := 0, 0
			for range n {
				v := input.Int(0, 1e9) % mod
				sum = (sum + v) % mod
				squares = (squares + v*v) % mod
			}
			p := remainder(sum*sum-squares, mod)
			out = append(out, strconv.Itoa(p*powMod(n*(n-1)%mod, mod-2, mod)%mod))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n + log M) for each test: with two maxima every order is nice. Else let k be the count of max - 1. An order is bad when the maximum is the last of these k + 1 members, so n! / (k + 1) orders are bad.
	"cf-1569-c": {Stdio: func(stdin string) string {
		const prime = 998244353
		input := sets.NewReader(stdin)
		var out []string
		total := 0
		for t := input.Int(1, 1e4); t > 0; t-- {
			n := input.Int(2, 2e5)
			total += n
			sets.Need(total <= 2e5)
			a := make([]int, n)
			for i := range a {
				a[i] = input.Int(1, 1e9)
			}
			top := slices.Max(a)
			fact, _ := factorials(n, prime)
			answer := fact[n]
			if count(a, top) == 1 {
				answer = remainder(answer-answer*powMod(count(a, top-1)+1, prime-2, prime)%prime, prime)
			}
			out = append(out, strconv.Itoa(answer))
		}
		input.End()
		return strings.Join(out, "\n") + "\n"
	}},
	// O(n + log M): a word of length L with the letter counts c1, c2, ... has L! / (c1! · c2! · ...) different orders of its letters. The tables fact and invFact replace the division.
	"count-anagrams": {Arity: 1, Call: func(args []any) any {
		s := sets.Str(args[0], anagramText)
		sets.Need(len(s) <= 1e5)
		fact, invFact := factorials(len(s), mod)
		answer := 1
		for word := range strings.SplitSeq(s, " ") {
			var letters [26]int
			for _, c := range word {
				letters[c-'a']++
			}
			answer = answer * fact[len(word)] % mod
			for _, c := range letters {
				answer = answer * invFact[c] % mod
			}
		}
		return sets.Py(answer)
	}},
}

package sets

import "errors"

// Reference is the reference solution of a learning set task. It stays on the server, so the student cannot read it.
// A call task has Call and Arity. A stdio task has Stdio.
//
// A reference panics with ErrConstraints for an input outside the constraints of the statement.
// For an input of a wrong shape or type it can panic with another value.
type Reference struct {
	// Arity is the number of arguments of a call task.
	Arity int
	// Call solves a call task. The arguments and the result are Python values of package answers.
	Call func(args []any) any
	// Stdio solves a stdio task: it returns the stdout text for a stdin text.
	Stdio func(stdin string) string
}

// Generator makes the hidden tests of a task: 10 to 16 inputs in the test format of the task. random is a seeded generator of numbers in [0, 1).
type Generator func(random func() float64) []string

// ErrConstraints is the panic value of a reference for an input outside the constraints of the statement.
var ErrConstraints = errors.New("constraints")

// Need panics with ErrConstraints when a constraint check fails.
func Need(ok bool) {
	if !ok {
		panic(ErrConstraints)
	}
}

// Mulberry32 is a seeded 32-bit generator of numbers in [0, 1): https://gist.github.com/tommyettinger/46a874533244883189143505d203312c
// It gives the same numbers as mulberry32 of the Worker, so a seed gives the same hidden tests.
func Mulberry32(seed uint32) func() float64 {
	return func() float64 {
		seed += 0x6D2B79F5
		z := seed
		z = (z ^ z>>15) * (z | 1)
		z ^= z + (z^z>>7)*(z|61)
		return float64(z^z>>14) / (1 << 32)
	}
}

// Package bank is the bank of Codeforces problems for the contests.
package bank

import (
	"embed"
	"encoding/json"
	"fmt"
)

// Test is the input of a test and its official answer.
type Test struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// Problem is a Codeforces problem with a Russian statement. The examples are visible, the tests are hidden.
type Problem struct {
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	Rating     int      `json:"rating"`
	Topics     []string `json:"topics"`
	Statement  string   `json:"statement"`
	IgnoreCase bool     `json:"ignoreCase"`
	Examples   []Test   `json:"examples"`
	Tests      []Test   `json:"tests"`
}

//go:embed problems-*.json
var files embed.FS

// Problems is all problems of the embedded files, files in name order.
// A bad file panics at start: TestBankContent reads the same files.
var Problems = func() []Problem {
	entries, err := files.ReadDir(".")
	if err != nil {
		panic(err)
	}
	var all []Problem
	for _, entry := range entries {
		data, err := files.ReadFile(entry.Name())
		if err != nil {
			panic(err)
		}
		var problems []Problem
		if err := json.Unmarshal(data, &problems); err != nil {
			panic(fmt.Sprintf("bank %s: %v", entry.Name(), err))
		}
		all = append(all, problems...)
	}
	return all
}()

// Difficulty returns "easy" for the rating 800, "medium" for 900 to 1100, "hard" for 1200, and "" for other ratings.
func Difficulty(rating int) string {
	switch rating {
	case 800:
		return "easy"
	case 900, 1000, 1100:
		return "medium"
	case 1200:
		return "hard"
	}
	return ""
}

// Points returns the points of a solved problem: 3 for easy, 4 for medium, 6 for hard, and 0 for other values.
func Points(difficulty string) int {
	switch difficulty {
	case "easy":
		return 3
	case "medium":
		return 4
	case "hard":
		return 6
	}
	return 0
}

package notebook

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/bank"
	"github.com/rustatian/redlamp.me/backend/internal/sets/content"
)

// The content checks of the problem bank.

var checkBankSlug = regexp.MustCompile(`^cf-\d+-[a-z]\d?$`)

// checkBankProblem is the check of one bank problem. taskSlugs is the slugs of the learning set tasks.
func checkBankProblem(t *testing.T, problem bank.Problem, taskSlugs map[string]bool) {
	t.Helper()
	if !checkBankSlug.MatchString(problem.Slug) {
		t.Errorf("slug %q has a wrong format", problem.Slug)
	}
	if taskSlugs[problem.Slug] {
		t.Errorf("a learning set has the task %q", problem.Slug)
	}
	if bank.Difficulty(problem.Rating) == "" {
		t.Errorf("rating %d is not 800 to 1200 in steps of 100", problem.Rating)
	}
	if len(problem.Topics) == 0 {
		t.Error("the problem has no topics")
	}
	var topics []string
	for _, topic := range problem.Topics {
		if topicNames[strings.ToLower(topic)] != topic {
			topics = append(topics, topic)
		}
	}
	checkNone(t, "topics that are not canonical names", topics)
	if checkBadStatement(problem.Statement) {
		t.Error("the statement has no example or constraints heading, or has a bad tag or an attribute")
	}
	headings := []string{"<h3>Входные данные</h3>", "<h3>Выходные данные</h3>"}
	for i := range problem.Examples {
		headings = append(headings, fmt.Sprintf("<h3>Пример %d</h3>", i+1))
	}
	var missing []string
	for _, heading := range headings {
		if !strings.Contains(problem.Statement, heading) {
			missing = append(missing, heading)
		}
	}
	checkNone(t, "statement headings that are missing", missing)
	if strings.Contains(problem.Statement, "$") {
		t.Error("the statement has a $ of LaTeX")
	}
	if len(problem.Examples) == 0 {
		t.Error("the problem has no examples")
	}
	var examples, tests []string
	for i, example := range problem.Examples {
		if !strings.HasSuffix(example.Input, "\n") {
			examples = append(examples, fmt.Sprint(i+1))
		}
	}
	checkNone(t, "examples with an input without a newline at the end", examples)
	if n := len(problem.Tests); n < 5 || n > 16 {
		t.Errorf("%d tests, want 5 to 16", n)
	}
	for i, test := range problem.Tests {
		if !strings.HasSuffix(test.Input, "\n") || len(test.Input) > 4096 || strings.TrimSpace(test.Output) == "" || len(test.Output) > 4096 {
			tests = append(tests, fmt.Sprint(i+1))
		}
	}
	checkNone(t, "tests with a wrong input or answer", tests)
}

func TestBankContent(t *testing.T) {
	taskSlugs := map[string]bool{}
	for _, set := range content.Sets {
		for _, task := range set.Tasks() {
			taskSlugs[task.Slug()] = true
		}
	}
	var slugs []string
	for _, problem := range bank.Problems {
		slugs = append(slugs, problem.Slug)
	}
	checkNone(t, "repeated problems", checkRepeated(slugs))
	for _, problem := range bank.Problems {
		t.Run(problem.Slug, func(t *testing.T) { checkBankProblem(t, problem, taskSlugs) })
	}
}

func TestBankDifficulty(t *testing.T) {
	tests := []struct {
		name   string
		rating int
		want   string
	}{
		{name: "800 is easy", rating: 800, want: "easy"},
		{name: "900 is medium", rating: 900, want: "medium"},
		{name: "1100 is medium", rating: 1100, want: "medium"},
		{name: "1200 is hard", rating: 1200, want: "hard"},
		{name: "700 is below the bank", rating: 700, want: ""},
		{name: "1300 is above the bank", rating: 1300, want: ""},
		{name: "1150 is not a step of 100", rating: 1150, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bank.Difficulty(tt.rating); got != tt.want {
				t.Errorf("Difficulty(%d) = %q, want %q", tt.rating, got, tt.want)
			}
		})
	}
}

// Package sets has the model of a learning set: chapters with a lesson, problems, quiz cards and figures.
// The JSON names are the names of shared/leetcode.ts.
package sets

import (
	"regexp"
	"strconv"
	"strings"
)

// LearningSet is an ordered sequence of chapters about one subject.
type LearningSet struct {
	Slug     string    `json:"slug"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Topics   []string  `json:"topics"`
	Chapters []Chapter `json:"chapters"`
}

// Chapter is one chapter of a learning set. Its number is its position in the set, from 1. Lesson is the HTML of the lesson.
type Chapter struct {
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Lesson  string   `json:"lesson"`
	Tasks   []Task   `json:"tasks"`
	Cards   []Card   `json:"cards"`
	Links   []Link   `json:"links"`
	Figures []Figure `json:"figures,omitempty"`
}

// Link is a link to more material about a chapter.
type Link struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source"`
}

// Task is a problem of a chapter. Number is nil for a problem without a number. Statement is HTML.
type Task struct {
	URL        string     `json:"url"`
	Number     *int64     `json:"number,omitempty"`
	Title      string     `json:"title"`
	Difficulty string     `json:"difficulty"`
	Topics     []string   `json:"topics"`
	Note       string     `json:"note"`
	Statement  string     `json:"statement,omitempty"`
	Tests      *TaskTests `json:"tests,omitempty"`
}

// TestSpec tells the runner how to call the code and how to compare an answer.
// Mode "call" is LeetCode style: the input is a Python list of the arguments, and the answer is a Python literal.
// Check is "return" (the return value), "arg0" (the first argument after the call) or "prefix" (the first argument cut to the returned length).
// Order is "exact", "any" (the order of the list is free) or "any-nested" (the order of each inner list and of the list is free).
// Mode "stdio" is Codeforces style: the input is the stdin text, and the answer is the stdout text, compared token by token.
type TestSpec struct {
	Mode       string `json:"mode"`
	Method     string `json:"method,omitempty"`
	Check      string `json:"check,omitempty"`
	Order      string `json:"order,omitempty"`
	IgnoreCase bool   `json:"ignoreCase,omitempty"`
}

// TestCase is one example. Expected is nil when the statement gives no output for it.
type TestCase struct {
	Input    string  `json:"input"`
	Expected *string `json:"expected"`
}

// TaskTests is the test spec of a task with its standard cases.
type TaskTests struct {
	TestSpec
	Cases []TestCase `json:"cases"`
}

// Card is a quiz card. Type is "choice", "short" or "code".
// Answer is the index of the right option (an int) for a choice card and a code card, and the right text for a short card.
type Card struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Prompt   string   `json:"prompt"`
	Language string   `json:"language,omitempty"`
	Code     string   `json:"code,omitempty"`
	Options  []string `json:"options,omitempty"`
	Answer   any      `json:"answer"`
	Explain  string   `json:"explain"`
}

// Figure is a lesson figure. The web client draws each frame as rows of cells, with indexes from 0.
type Figure struct {
	ID     string        `json:"id"`
	Title  string        `json:"title"`
	Frames []FigureFrame `json:"frames"`
}

// FigureFrame is one step of a figure.
type FigureFrame struct {
	Caption string      `json:"caption"`
	Rows    []FigureRow `json:"rows"`
}

// FigureRow is one row of cells. A cell is a number or a text. A tree row draws its cells as a binary tree in heap order.
type FigureRow struct {
	Label    string    `json:"label,omitempty"`
	Cells    []any     `json:"cells"`
	Marked   []int     `json:"marked,omitempty"`
	Dimmed   []int     `json:"dimmed,omitempty"`
	Pointers *Pointers `json:"pointers,omitempty"`
	Window   *[2]int   `json:"window,omitempty"`
	Tree     bool      `json:"tree,omitempty"`
}

var (
	leetcodePath   = regexp.MustCompile(`^/problems/([a-z0-9]+(?:-[a-z0-9]+)*)(?:/|$)`)
	codeforcesPath = regexp.MustCompile(`^/(?:problemset/problem/(\d{1,6})|contest/(\d{1,6})/problem)/([A-Za-z]\d?)(?:/|$)`)
	codeforcesSlug = regexp.MustCompile(`^cf-(\d+)-([a-z]\d?)$`)
)

// Slug returns the slug of a problem from the host and the escaped path of its URL.
// A Codeforces problem has the slug cf-<contest>-<index>, with the index in lower case.
func Slug(host, path string) (string, bool) {
	if match := leetcodePath.FindStringSubmatch(path); host == "leetcode.com" && match != nil {
		return match[1], true
	}
	if match := codeforcesPath.FindStringSubmatch(path); host == "codeforces.com" && match != nil {
		return "cf-" + match[1] + match[2] + "-" + strings.ToLower(match[3]), true
	}
	return "", false
}

// URL returns the canonical URL of a problem slug.
func URL(slug string) string {
	if match := codeforcesSlug.FindStringSubmatch(slug); match != nil {
		return "https://codeforces.com/problemset/problem/" + match[1] + "/" + strings.ToUpper(match[2])
	}
	return "https://leetcode.com/problems/" + slug + "/"
}

// Slug returns the slug of the task. The URL of a task is https://leetcode.com/... or https://codeforces.com/..., which a content test checks.
func (t Task) Slug() string {
	rest, _ := strings.CutPrefix(t.URL, "https://")
	host, path, _ := strings.Cut(rest, "/")
	slug, _ := Slug(host, "/"+path)
	return slug
}

// Label is the short name of the task: its source and its number.
func (t Task) Label() string {
	if rest, ok := strings.CutPrefix(t.Slug(), "cf-"); ok {
		return "Codeforces " + strings.ToUpper(strings.Replace(rest, "-", "", 1))
	}
	if t.Number != nil && *t.Number != 0 {
		return "LeetCode " + strconv.FormatInt(*t.Number, 10)
	}
	return "LeetCode"
}

// Tasks returns the problems of all chapters.
func (s LearningSet) Tasks() []Task {
	var tasks []Task
	for _, chapter := range s.Chapters {
		tasks = append(tasks, chapter.Tasks...)
	}
	return tasks
}

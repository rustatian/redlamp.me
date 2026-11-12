package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

type ownTest struct {
	ID        string `json:"id"`
	Input     string `json:"input"`
	CreatedAt string `json:"createdAt"`
}

type problemTests struct {
	Spec     *sets.TestSpec  `json:"spec"`
	Source   string          `json:"source"`
	Standard []sets.TestCase `json:"standard"`
	Own      []ownTest       `json:"own"`
	Skeleton string          `json:"skeleton"`
}

type testRun struct {
	ApproachID string `json:"approachId"`
	Passed     int64  `json:"passed"`
	Total      int64  `json:"total"`
	CodeHash   string `json:"codeHash"`
	CreatedAt  string `json:"createdAt"`
}

const (
	stdinSkeleton    = "def ans():\n    N = int(input())\n    a = list(map(int, input().split()))\n    # Посчитайте ответ и верните его.\n    return 0\n\nt = int(input())\na = [ans() for _ in range(t)]\nfor res in a:\n    print(res)\n"
	solutionSkeleton = "class Solution:\n    pass"
	maxOwnTests      = 20
)

var (
	pythonWords = map[string]string{"true": "True", "false": "False", "null": "None"}
	htmlWords   = strings.NewReplacer("&quot;", `"`, "&#39;", "'", "&lt;", "<", "&gt;", ">", "&amp;", "&", "&nbsp;", " ")
	// jsonToken matches a JSON string, or one of the words true, false and null outside a string.
	jsonToken = regexp.MustCompile(`"(?:[^"\\]|\\.)*"|\b(?:true|false|null)\b`)
	htmlTag   = regexp.MustCompile(`<[^>]*>`)
	// outputLine matches an "Output:" line of a LeetCode statement. The class is the white space of a JavaScript regular expression.
	outputLine    = regexp.MustCompile(`(?m)^[\s\p{Z}\x{FEFF}]*Output:[\s\p{Z}\x{FEFF}]*(.+)$`)
	upperCase     = regexp.MustCompile(`[A-Z]`)
	codeHash      = regexp.MustCompile(`^[a-f0-9]{64}$`)
	testsPath     = regexp.MustCompile(`^/leetcode/api/problems/([^/]+)/tests(?:/([^/]+))?$`)
	statementPath = regexp.MustCompile(`^/leetcode/api/problems/([^/]+)/statement$`)
	runsPath      = regexp.MustCompile(`^/leetcode/api/problems/([^/]+)/runs$`)
)

// python changes the JSON values of LeetCode to Python literals: the strings stay, and true, false and null change.
func python(value string) string {
	return jsonToken.ReplaceAllStringFunc(value, func(token string) string {
		if word, ok := pythonWords[token]; ok {
			return word
		}
		return token
	})
}

// lookupResult is the data of one LeetCode question for the runner. tests is nil for a problem without a runner.
type lookupResult struct {
	tests    *sets.TestSpec
	cases    []sets.TestCase
	skeleton *string
}

// lookup gets the examples and the Python 3 starter code of a LeetCode problem. It returns nil when LeetCode has no such problem.
// shortcut: the examples compare in exact order, so a correct answer in another order fails. Add an order rule when a student finds such a problem.
func (s *Server) lookup(ctx context.Context, titleSlug string) (*lookupResult, error) {
	raw, err := s.leetcode.Question(ctx, titleSlug, "metaData exampleTestcaseList content codeSnippets { langSlug code }")
	if err != nil || raw == nil {
		return nil, err
	}
	var question struct {
		MetaData            string   `json:"metaData"`
		ExampleTestcaseList []string `json:"exampleTestcaseList"`
		Content             *string  `json:"content"`
		CodeSnippets        []struct {
			LangSlug string `json:"langSlug"`
			Code     string `json:"code"`
		} `json:"codeSnippets"`
	}
	if err := json.Unmarshal(raw, &question); err != nil {
		return nil, err
	}
	result := &lookupResult{}
	for _, snippet := range question.CodeSnippets {
		if snippet.LangSlug == "python3" {
			result.skeleton = &snippet.Code
			break
		}
	}
	type typed struct {
		Type string `json:"type"`
	}
	var meta struct {
		Name   string   `json:"name"`
		Manual bool     `json:"manual"`
		Params *[]typed `json:"params"`
		Return *typed   `json:"return"`
		Output *struct {
			Paramindex *float64 `json:"paramindex"`
		} `json:"output"`
	}
	if err := json.Unmarshal([]byte(question.MetaData), &meta); err != nil {
		return nil, err
	}
	// The harness gives plain Python values, so a design problem (no name) or an interactive problem (manual) gets no runner.
	if meta.Name == "" || meta.Manual {
		return result, nil
	}
	if meta.Params == nil || meta.Return == nil {
		return nil, errors.New("the LeetCode metadata has no parameters or no return type")
	}
	// A type with an upper-case letter is a class such as TreeNode.
	for _, item := range append(*meta.Params, *meta.Return) {
		if upperCase.MatchString(item.Type) {
			return result, nil
		}
	}
	content := ""
	if question.Content != nil {
		content = *question.Content
	}
	statement := htmlWords.Replace(htmlTag.ReplaceAllString(content, ""))
	var outputs []string
	for _, match := range outputLine.FindAllStringSubmatch(statement, -1) {
		outputs = append(outputs, python(strings.TrimSpace(match[1])))
	}
	result.tests = &sets.TestSpec{Mode: "call", Method: meta.Name}
	if meta.Output != nil && meta.Output.Paramindex != nil && *meta.Output.Paramindex == 0 {
		result.tests.Check = "arg0"
	}
	for i, input := range question.ExampleTestcaseList {
		example := sets.TestCase{Input: python("[" + strings.ReplaceAll(input, "\n", ",") + "]")}
		if len(outputs) == len(question.ExampleTestcaseList) {
			example.Expected = &outputs[i]
		}
		result.cases = append(result.cases, example)
	}
	return result, nil
}

func (s *Server) handleTests(w http.ResponseWriter, r *http.Request, _ *SessionUser) (bool, error) {
	match := testsPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	return true, s.tests(w, r, match[1], match[2])
}

func (s *Server) tests(w http.ResponseWriter, r *http.Request, problemID, testID string) error {
	ctx := r.Context()
	allowed := []string{http.MethodGet, http.MethodPost}
	if testID != "" {
		allowed = []string{http.MethodDelete}
	}
	if err := methods(r, allowed...); err != nil {
		return err
	}
	problem, err := getProblem(ctx, s.db, problemID)
	if err != nil {
		return err
	}
	if testID != "" {
		if _, err := uuid(testID); err != nil {
			return err
		}
		tag, err := s.db.Exec(ctx, "DELETE FROM problem_tests WHERE id=$1 AND problem_id=$2", testID, problem.ID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такого теста нет."}
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	if r.Method == http.MethodPost {
		body, err := readJSON(r, 32_768)
		if err != nil {
			return err
		}
		data, err := object(body, "input")
		if err != nil {
			return err
		}
		input, err := text(data["input"], 10_000, true)
		if err != nil {
			return err
		}
		value := ownTest{ID: newID(), Input: input, CreatedAt: isoNow()}
		err = s.tx(ctx, func(tx pgx.Tx) error {
			// The count and the insert are one rule for the problem, so the transaction writes the problem row.
			if err := touchProblem(ctx, tx, problem.ID); err != nil {
				return err
			}
			tag, err := tx.Exec(ctx, "INSERT INTO problem_tests (id,problem_id,input,created_at) SELECT $1,$2,$3,$4 WHERE (SELECT count(*) FROM problem_tests WHERE problem_id=$2)<$5", value.ID, problem.ID, value.Input, value.CreatedAt, maxOwnTests)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return &HTTPError{Status: http.StatusConflict, Code: "too_many_tests", Message: "У задачи уже 20 своих тестов. Удалите лишний тест."}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusCreated, value)
	}
	codeforces := strings.HasPrefix(problem.Slug, "cf-")
	var found *lookupResult
	if !codeforces {
		// A failed lookup gives the page without examples.
		found, _ = s.lookup(ctx, problem.Slug)
	}
	result := problemTests{Source: "none", Standard: []sets.TestCase{}, Own: []ownTest{}}
	// The tests of a learning set task come first, then the examples of LeetCode.
	task, inSet := s.setTask(problem.Slug)
	switch {
	case inSet && task.Tests != nil:
		spec := task.Tests.TestSpec
		result.Spec, result.Source = &spec, "set"
		result.Standard = append(result.Standard, task.Tests.Cases...)
	case found != nil && found.tests != nil:
		result.Spec, result.Source = found.tests, "leetcode"
		result.Standard = append(result.Standard, found.cases...)
	case codeforces:
		result.Spec = &sets.TestSpec{Mode: "stdio"}
	}
	rows, err := s.db.Query(ctx, "SELECT id,input,created_at FROM problem_tests WHERE problem_id=$1 ORDER BY created_at,id", problem.ID)
	if err != nil {
		return err
	}
	own, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (ownTest, error) {
		var test ownTest
		return test, row.Scan(&test.ID, &test.Input, &test.CreatedAt)
	})
	if err != nil {
		return err
	}
	result.Own = append(result.Own, own...)
	switch {
	case codeforces:
		result.Skeleton = stdinSkeleton
	case found != nil && found.skeleton != nil:
		result.Skeleton = *found.skeleton
	case result.Spec != nil && result.Spec.Method != "":
		result.Skeleton = solutionSkeleton
	}
	return writeJSON(w, r, http.StatusOK, result)
}

func (s *Server) handleStatement(w http.ResponseWriter, r *http.Request, _ *SessionUser) (bool, error) {
	match := statementPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	if err := methods(r, http.MethodGet); err != nil {
		return true, err
	}
	problem, err := getProblem(r.Context(), s.db, match[1])
	if err != nil {
		return true, err
	}
	var statement struct {
		HTML *string `json:"html"`
		Text *string `json:"text"`
	}
	if task, ok := s.setTask(problem.Slug); ok && task.Statement != "" {
		statement.HTML = &task.Statement
	} else if !strings.HasPrefix(problem.Slug, "cf-") {
		if content := s.leetcodeContent(r.Context(), problem.Slug); content != "" {
			plain := htmlText(content)
			statement.Text = &plain
		}
	}
	return true, writeJSON(w, r, http.StatusOK, statement)
}

// leetcodeContent returns the statement HTML of a LeetCode problem. A failed lookup and a problem without content give no text.
func (s *Server) leetcodeContent(ctx context.Context, slug string) string {
	var question struct {
		Content string `json:"content"`
	}
	if raw, err := s.leetcode.Question(ctx, slug, "content"); err == nil && raw != nil {
		_ = json.Unmarshal(raw, &question)
	}
	return question.Content
}

func (s *Server) handleRuns(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	match := runsPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	return true, s.runs(w, r, user, match[1])
}

func (s *Server) runs(w http.ResponseWriter, r *http.Request, user *SessionUser, problemID string) error {
	ctx := r.Context()
	if err := methods(r, http.MethodGet, http.MethodPost); err != nil {
		return err
	}
	if _, err := uuid(problemID); err != nil {
		return err
	}
	if r.Method == http.MethodGet {
		var attemptID any
		if query := r.URL.Query(); query.Has("attemptId") {
			attemptID = query.Get("attemptId")
		}
		if _, err := uuid(attemptID); err != nil {
			return err
		}
		// The latest run of each approach of an attempt.
		rows, err := s.db.Query(ctx, "SELECT DISTINCT ON (approach_id) approach_id,passed,total,code_hash,created_at FROM test_runs WHERE problem_id=$1 AND attempt_id=$2 ORDER BY approach_id,created_at DESC", problemID, attemptID)
		if err != nil {
			return err
		}
		runs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (testRun, error) {
			var run testRun
			return run, row.Scan(&run.ApproachID, &run.Passed, &run.Total, &run.CodeHash, &run.CreatedAt)
		})
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, append([]testRun{}, runs...))
	}
	if err := requireRole(user, "student"); err != nil {
		return err
	}
	body, err := readJSON(r, 32_768)
	if err != nil {
		return err
	}
	data, err := object(body, "attemptId", "approachId", "codeHash", "passed", "total", "runOn")
	if err != nil {
		return err
	}
	var value testRun
	if value.Total, err = integer(data["total"], 0, 100); err != nil {
		return err
	}
	if value.ApproachID, err = uuid(data["approachId"]); err != nil {
		return err
	}
	if value.Passed, err = integer(data["passed"], 0, value.Total); err != nil {
		return err
	}
	if value.CodeHash, err = text(data["codeHash"], 64, false); err != nil {
		return err
	}
	value.CreatedAt = isoNow()
	if !codeHash.MatchString(value.CodeHash) {
		return invalid("Проверьте поля запроса.")
	}
	day, err := text(data["runOn"], 10, true)
	if err != nil {
		return err
	}
	runOn, err := calendarDate(day)
	if err != nil {
		return err
	}
	attemptID, err := uuid(data["attemptId"])
	if err != nil {
		return err
	}
	tag, err := s.db.Exec(ctx, "INSERT INTO test_runs (id,problem_id,attempt_id,approach_id,code_hash,passed,total,run_on,created_at) SELECT $1,problem_id,id,$2,$3,$4,$5,$6,$7 FROM attempts WHERE id=$8 AND problem_id=$9",
		newID(), value.ApproachID, value.CodeHash, value.Passed, value.Total, *runOn, value.CreatedAt, attemptID, problemID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой попытки нет."}
	}
	return writeJSON(w, r, http.StatusCreated, value)
}

var (
	hiddenElement = regexp.MustCompile(`(?i)<(script|style|template|svg)\b`)
	superscript   = regexp.MustCompile(`(?is)<sup\b[^>]*>(.*?)</sup\s*>`)
	heading       = regexp.MustCompile(`(?i)<h([1-4])\b[^>]*>`)
	blockTag      = regexp.MustCompile(`(?i)</?(?:address|article|aside|blockquote|br|dd|details|div|dl|dt|figcaption|figure|footer|h[1-6]|header|hr|li|main|nav|ol|p|pre|section|summary|table|td|th|tr|ul)\b[^>]*>`)
	innerSpaces   = regexp.MustCompile(`(\S)[ \t]+`)
	lineEndSpaces = regexp.MustCompile(`(?m)[ \t]+$`)
	lineBreaks    = regexp.MustCompile(`\n+`)
	htmlEntities  = strings.NewReplacer("&lt;", "<", "&gt;", ">", "&amp;", "&", "&quot;", `"`, "&#39;", "'", "&nbsp;", " ")
)

// htmlText changes HTML to plain text: it drops scripts and styles, keeps the line structure of the blocks and marks headings with #.
func htmlText(html string) string {
	// Drop each script, style, template and svg element with its content, from the left.
	var visible strings.Builder
	for {
		open := hiddenElement.FindStringSubmatchIndex(html)
		if open == nil {
			break
		}
		closing := regexp.MustCompile(`(?i)</` + html[open[2]:open[3]] + `\s*>`).FindStringIndex(html[open[1]:])
		if closing == nil {
			visible.WriteString(html[:open[1]])
			html = html[open[1]:]
			continue
		}
		visible.WriteString(html[:open[0]])
		html = html[open[1]+closing[1]:]
	}
	visible.WriteString(html)
	text := superscript.ReplaceAllString(visible.String(), "^($1)")
	text = heading.ReplaceAllStringFunc(text, func(tag string) string {
		return "\n" + strings.Repeat("#", int(heading.FindStringSubmatch(tag)[1][0]-'0')) + " "
	})
	text = htmlEntities.Replace(htmlTag.ReplaceAllString(blockTag.ReplaceAllString(text, "\n"), ""))
	text = lineBreaks.ReplaceAllString(lineEndSpaces.ReplaceAllString(innerSpaces.ReplaceAllString(text, "$1 "), ""), "\n")
	return strings.TrimSpace(text)
}

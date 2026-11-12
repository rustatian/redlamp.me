package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

const (
	testsInsertProblem = "INSERT INTO problems (id,slug,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,'Problem','medium','[]','','','now','now')"
	testsInsertOwn     = "INSERT INTO problem_tests (id,problem_id,input,created_at) VALUES ($1,$2,'1','2026-10-01T10:00:00.000Z')"
	testsInsertAttempt = "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,saved_at,created_at,updated_at) VALUES ($1,$2,$3,$4,'','not_submitted','needs_practice',NULL,$5,'now')"
	testsInsertRun     = "INSERT INTO test_runs (id,problem_id,attempt_id,approach_id,code_hash,passed,total,run_on,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)"
)

type testsCodeSnippet struct {
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
}

// testsSnippets are the code snippets of a LeetCode answer. The third one is the Python 3 skeleton.
var testsSnippets = []testsCodeSnippet{
	{LangSlug: "cpp", Code: "class Solution {\n};"},
	{LangSlug: "python", Code: "class Solution(object):\n    def solve(self):\n        "},
	{LangSlug: "python3", Code: "class Solution:\n    def solve(self) -> int:\n        "},
}

// The skeletons of a Codeforces problem and of a LeetCode problem without a Python 3 snippet.
const (
	testsStdin    = "def ans():\n    N = int(input())\n    a = list(map(int, input().split()))\n    # Посчитайте ответ и верните его.\n    return 0\n\nt = int(input())\na = [ans() for _ in range(t)]\nfor res in a:\n    print(res)\n"
	testsFallback = "class Solution:\n    pass"
)

// testsTwoSum is the metaData of the LeetCode problem two-sum.
const testsTwoSum = `{"name":"twoSum","params":[{"name":"nums","type":"integer[]"},{"name":"target","type":"integer"}],"return":{"type":"integer[]","size":2},"manual":false}`

// testsReply returns a fake LeetCode reply with the question object.
func testsReply(question map[string]any) func(string, string) (json.RawMessage, error) {
	return func(string, string) (json.RawMessage, error) { return json.Marshal(question) }
}

// testsQuestion is leetcode() of test/tests.test.ts: the fields of a GraphQL answer.
// The fixtures copy the shapes of LeetCode: <pre> examples (two-sum), <span class="example-io"> examples (valid-parentheses), null content and code snippets (a paid-only problem).
func testsQuestion(metaData string, exampleTestcaseList []string, content, codeSnippets any) func(string, string) (json.RawMessage, error) {
	return testsReply(map[string]any{"metaData": metaData, "exampleTestcaseList": exampleTestcaseList, "content": content, "codeSnippets": codeSnippets})
}

func testsPre(output string) string {
	return "<pre>\n<strong>Input:</strong> value\n<strong>Output:</strong> " + output + "\n<strong>Explanation:</strong> Output: text.\n</pre>"
}

func testsSpan(output string) string {
	return `<p><strong>Input:</strong> <span class="example-io">value</span></p>` + "\n" + `<p><strong>Output: </strong><span class="example-io" style="font-family: Menlo,sans-serif;">` + output + `</span></p>`
}

// The fake LeetCode replies for an unknown question, an HTTP error status, a network error and a request without an answer before the timeout.
func testsNoQuestion(string, string) (json.RawMessage, error) { return nil, nil }

func testsBusy(string, string) (json.RawMessage, error) {
	return nil, errors.New("LeetCode returned 503")
}

func testsUnreachable(string, string) (json.RawMessage, error) {
	return nil, errors.New("network connection lost")
}

func testsNoAnswer(string, string) (json.RawMessage, error) { return nil, context.DeadlineExceeded }

// testsSameJSON reports a difference between the JSON text got and the JSON form of want.
func testsSameJSON(t *testing.T, got []byte, want any) {
	t.Helper()
	text, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var gotValue, wantValue any
	if err := json.Unmarshal(text, &wantValue); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("body %q: %v", got, err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("body = %s, want %s", got, text)
	}
}

// The cases of test/tests.test.ts for GET /problems/{id}/tests.
// spec and standard are JSON text. own is always empty.
// subset is expect.arrayContaining: the response has each listed standard case and can have more.
func TestTestsStandard(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		reply    func(titleSlug, fields string) (json.RawMessage, error)
		spec     string
		source   string
		standard string
		subset   bool
		skeleton string
	}{
		{
			name:     "LeetCode task of a set",
			slug:     "two-sum-ii-input-array-is-sorted",
			reply:    testsQuestion(testsTwoSum, []string{"[2,7,11,15]\n9"}, testsPre("[0,1]"), testsSnippets),
			spec:     `{"mode":"call","method":"twoSum"}`,
			source:   "set",
			standard: `[{"input":"[[2,7,11,15],9]","expected":"[1,2]"}]`,
			subset:   true,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "LeetCode task of a set with unreachable LeetCode",
			slug:     "two-sum-ii-input-array-is-sorted",
			reply:    testsUnreachable,
			spec:     `{"mode":"call","method":"twoSum"}`,
			source:   "set",
			standard: `[{"input":"[[2,7,11,15],9]","expected":"[1,2]"}]`,
			subset:   true,
			skeleton: testsFallback,
		},
		{
			name:     "Codeforces task of a set",
			slug:     "cf-1821-a",
			reply:    nil,
			spec:     `{"mode":"stdio"}`,
			source:   "set",
			standard: `[{"input":"8\n??\n?\n0\n9\n03\n1??7\n?5?\n9??99\n","expected":"90\n9\n0\n1\n0\n100\n90\n100\n"}]`,
			subset:   true,
			skeleton: testsStdin,
		},
		{name: "Codeforces problem outside the sets", slug: "cf-1-a", reply: nil, spec: `{"mode":"stdio"}`, source: "none", standard: `[]`, skeleton: testsStdin},
		{
			name:     "LeetCode examples with two outputs",
			slug:     "two-sum",
			reply:    testsQuestion(testsTwoSum, []string{"[2,7,11,15]\n9", "[3,3]\n6"}, testsPre("[0,1]")+"\n"+testsPre("[0,1]"), testsSnippets),
			spec:     `{"mode":"call","method":"twoSum"}`,
			source:   "leetcode",
			standard: `[{"input":"[[2,7,11,15],9]","expected":"[0,1]"},{"input":"[[3,3],6]","expected":"[0,1]"}]`,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "LeetCode in-place examples",
			slug:     "rotate-array",
			reply:    testsQuestion(`{"name":"rotate","params":[{"name":"nums","type":"integer[]"},{"name":"k","type":"integer"}],"return":{"type":"void"},"output":{"paramindex":0}}`, []string{"[1,2,3,4,5,6,7]\n3", "[-1,-100,3,99]\n2"}, testsPre("[5,6,7,1,2,3,4]")+"\n"+testsPre("[3,99,-1,-100]"), testsSnippets),
			spec:     `{"mode":"call","method":"rotate","check":"arg0"}`,
			source:   "leetcode",
			standard: `[{"input":"[[1,2,3,4,5,6,7],3]","expected":"[5,6,7,1,2,3,4]"},{"input":"[[-1,-100,3,99],2]","expected":"[3,99,-1,-100]"}]`,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "LeetCode boolean examples",
			slug:     "kids-with-the-greatest-number-of-candies",
			reply:    testsQuestion(`{"name":"kidsWithCandies","params":[{"name":"candies","type":"integer[]"},{"name":"extraCandies","type":"integer"}],"return":{"type":"list<boolean>"}}`, []string{"[2,3,5,1,3]\n3", "[12,1,12]\n10"}, testsSpan("[true,true,true,false,true]")+"\n"+testsSpan("[true,false,true]"), testsSnippets),
			spec:     `{"mode":"call","method":"kidsWithCandies"}`,
			source:   "leetcode",
			standard: `[{"input":"[[2,3,5,1,3],3]","expected":"[True,True,True,False,True]"},{"input":"[[12,1,12],10]","expected":"[True,False,True]"}]`,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "LeetCode string examples",
			slug:     "longest-palindromic-substring",
			reply:    testsQuestion(`{"name":"longestPalindrome","params":[{"name":"s","type":"string"}],"return":{"type":"string"}}`, []string{`"babad"`, `"true"`}, testsPre("&quot;bab&quot;")+"\n"+testsPre("&quot;t&quot;"), testsSnippets),
			spec:     `{"mode":"call","method":"longestPalindrome"}`,
			source:   "leetcode",
			standard: `[{"input":"[\"babad\"]","expected":"\"bab\""},{"input":"[\"true\"]","expected":"\"t\""}]`,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "LeetCode problem without a Python 3 snippet",
			slug:     "two-sum",
			reply:    testsQuestion(testsTwoSum, []string{"[3,3]\n6"}, testsPre("[0,1]"), testsSnippets[:2]),
			spec:     `{"mode":"call","method":"twoSum"}`,
			source:   "leetcode",
			standard: `[{"input":"[[3,3],6]","expected":"[0,1]"}]`,
			skeleton: testsFallback,
		},
		{
			name:     "LeetCode examples with a missing output",
			slug:     "two-sum",
			reply:    testsQuestion(testsTwoSum, []string{"[2,7,11,15]\n9", "[3,3]\n6"}, testsPre("[0,1]"), testsSnippets),
			spec:     `{"mode":"call","method":"twoSum"}`,
			source:   "leetcode",
			standard: `[{"input":"[[2,7,11,15],9]","expected":null},{"input":"[[3,3],6]","expected":null}]`,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "paid-only LeetCode problem",
			slug:     "meeting-rooms",
			reply:    testsQuestion(`{"name":"canAttendMeetings","params":[{"name":"intervals","type":"integer[][]"}],"return":{"type":"boolean"}}`, []string{"[[0,30],[5,10],[15,20]]", "[[7,10],[2,4]]"}, nil, nil),
			spec:     `{"mode":"call","method":"canAttendMeetings"}`,
			source:   "leetcode",
			standard: `[{"input":"[[[0,30],[5,10],[15,20]]]","expected":null},{"input":"[[[7,10],[2,4]]]","expected":null}]`,
			skeleton: testsFallback,
		},
		{
			name:     "LeetCode tree problem",
			slug:     "binary-tree-inorder-traversal",
			reply:    testsQuestion(`{"name":"inorderTraversal","params":[{"name":"root","type":"TreeNode"}],"return":{"type":"list<integer>","dealloc":true}}`, []string{"[1,null,2,3]"}, testsSpan("[1,3,2]"), testsSnippets),
			spec:     `null`,
			source:   "none",
			standard: `[]`,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "LeetCode interactive problem",
			slug:     "first-bad-version",
			reply:    testsQuestion(`{"name":"firstBadVersion","params":[{"name":"n","type":"integer"},{"name":"bad","type":"integer"}],"return":{"type":"integer"},"manual":true}`, []string{"5\n4", "1\n1"}, testsPre("4")+"\n"+testsPre("1"), testsSnippets),
			spec:     `null`,
			source:   "none",
			standard: `[]`,
			skeleton: testsSnippets[2].Code,
		},
		{
			name:     "LeetCode design problem",
			slug:     "lru-cache",
			reply:    testsQuestion(`{"classname":"LRUCache","constructor":{"params":[{"type":"integer","name":"capacity"}]},"methods":[{"name":"get","params":[{"type":"integer","name":"key"}],"return":{"type":"integer"}}],"systemdesign":true}`, []string{"[\"LRUCache\",\"get\"]\n[[2],[1]]"}, "<pre>\n<strong>Output</strong>\n[null, -1]\n</pre>", testsSnippets),
			spec:     `null`,
			source:   "none",
			standard: `[]`,
			skeleton: testsSnippets[2].Code,
		},
		{name: "unknown LeetCode problem", slug: "two-sum", reply: testsNoQuestion, spec: `null`, source: "none", standard: `[]`, skeleton: ""},
		{name: "LeetCode error response", slug: "two-sum", reply: testsBusy, spec: `null`, source: "none", standard: `[]`, skeleton: ""},
		{name: "unreachable LeetCode", slug: "two-sum", reply: testsUnreachable, spec: `null`, source: "none", standard: `[]`, skeleton: ""},
		{name: "LeetCode without an answer", slug: "two-sum", reply: testsNoAnswer, spec: `null`, source: "none", standard: `[]`, skeleton: ""},
	}
	learningSets := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			fake := &fakeLeetCode{reply: tt.reply}
			s, err := New(Config{DB: pool, PublicOrigin: publicOrigin, Accounts: testAccounts, LeetCode: fake, Sets: learningSets})
			if err != nil {
				t.Fatal(err)
			}
			cookie := loginAs(t, s, "student")
			id := newID()
			exec(t, pool, testsInsertProblem, id, tt.slug)
			w := request(t, s, "/problems/"+id+"/tests", http.MethodGet, nil, cookie)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			got := w.Body.Bytes()
			if tt.subset {
				var listed []any
				if err := json.Unmarshal([]byte(tt.standard), &listed); err != nil {
					t.Fatal(err)
				}
				body := decode[map[string]any](t, w)
				standard, _ := body["standard"].([]any)
				for _, item := range listed {
					if !slices.ContainsFunc(standard, func(value any) bool { return reflect.DeepEqual(value, item) }) {
						t.Errorf("standard = %v, want the case %v", body["standard"], item)
					}
				}
				// The other fields are compared exactly.
				body["standard"] = listed
				if got, err = json.Marshal(body); err != nil {
					t.Fatal(err)
				}
			}
			testsSameJSON(t, got, map[string]any{"spec": json.RawMessage(tt.spec), "source": tt.source, "standard": json.RawMessage(tt.standard), "own": []any{}, "skeleton": tt.skeleton})
			if tt.reply == nil {
				if len(fake.calls) != 0 {
					t.Errorf("LeetCode calls = %+v, want none", fake.calls)
				}
				return
			}
			// One request asks for the examples and the code snippets.
			if len(fake.calls) != 1 || fake.calls[0].titleSlug != tt.slug || !strings.Contains(fake.calls[0].fields, "codeSnippets { langSlug code }") {
				t.Errorf("LeetCode calls = %+v, want one call for %s with the code snippets", fake.calls, tt.slug)
			}
		})
	}
}

// testsStatement is statement() of test/tests.test.ts: the statement of the set task with the slug.
func testsStatement(t *testing.T, learningSets []sets.LearningSet, slug string) string {
	t.Helper()
	for _, set := range learningSets {
		for _, task := range set.Tasks() {
			if task.Slug() == slug {
				return task.Statement
			}
		}
	}
	t.Fatalf("no task %s in the learning sets", slug)
	return ""
}

// The cases of test/tests.test.ts for GET /problems/{id}/statement.
func TestTestsStatement(t *testing.T) {
	const content = "<p>Find <code>nums[i] + nums[j] == target</code>.</p>\n<pre>\n<strong>Input:</strong> nums = [3,3], target = 6\n</pre>\n<ul><li><code>2 &lt;= nums.length</code></li></ul>"
	learningSets := goldenSets(t)
	tests := []struct {
		name     string
		slug     string
		reply    func(titleSlug, fields string) (json.RawMessage, error)
		wantHTML any
		wantText any
	}{
		{name: "statement of a LeetCode task of a set", slug: "two-sum-ii-input-array-is-sorted", reply: nil, wantHTML: testsStatement(t, learningSets, "two-sum-ii-input-array-is-sorted"), wantText: nil},
		{name: "statement of a Codeforces task of a set", slug: "cf-1821-a", reply: nil, wantHTML: testsStatement(t, learningSets, "cf-1821-a"), wantText: nil},
		{name: "statement of a Codeforces problem outside the sets", slug: "cf-1-a", reply: nil, wantHTML: nil, wantText: nil},
		{name: "statement of a LeetCode problem outside the sets", slug: "two-sum", reply: testsReply(map[string]any{"content": content}), wantHTML: nil, wantText: "Find nums[i] + nums[j] == target.\nInput: nums = [3,3], target = 6\n2 <= nums.length"},
		{name: "statement of a paid-only LeetCode problem", slug: "meeting-rooms", reply: testsReply(map[string]any{"content": nil}), wantHTML: nil, wantText: nil},
		{name: "statement of an unknown LeetCode problem", slug: "two-sum", reply: testsNoQuestion, wantHTML: nil, wantText: nil},
		{name: "statement with a LeetCode error response", slug: "two-sum", reply: testsBusy, wantHTML: nil, wantText: nil},
		{name: "statement with unreachable LeetCode", slug: "two-sum", reply: testsUnreachable, wantHTML: nil, wantText: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			fake := &fakeLeetCode{reply: tt.reply}
			s, err := New(Config{DB: pool, PublicOrigin: publicOrigin, Accounts: testAccounts, LeetCode: fake, Sets: learningSets})
			if err != nil {
				t.Fatal(err)
			}
			cookie := loginAs(t, s, "student")
			id := newID()
			exec(t, pool, testsInsertProblem, id, tt.slug)
			w := request(t, s, "/problems/"+id+"/statement", http.MethodGet, nil, cookie)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			testsSameJSON(t, w.Body.Bytes(), map[string]any{"html": tt.wantHTML, "text": tt.wantText})
			wantCalls := 0
			if tt.reply != nil {
				wantCalls = 1
			}
			if len(fake.calls) != wantCalls || (wantCalls == 1 && fake.calls[0].titleSlug != tt.slug) {
				t.Errorf("LeetCode calls = %+v, want %d for %s", fake.calls, wantCalls, tt.slug)
			}
		})
	}
}

func TestTestsUnknownProblemAndMethods(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		method     string
		wantStatus int
	}{
		{name: "unknown problem", path: "tests", method: http.MethodGet, wantStatus: 404},
		{name: "own test of an unknown problem", path: "tests", method: http.MethodPost, wantStatus: 404},
		{name: "change of the standard tests", path: "tests", method: http.MethodPut, wantStatus: 405},
		{name: "statement of an unknown problem", path: "statement", method: http.MethodGet, wantStatus: 404},
		{name: "change of the statement", path: "statement", method: http.MethodPut, wantStatus: 405},
		{name: "expected answers of an unknown problem", path: "tests/expected", method: http.MethodPost, wantStatus: 404},
		{name: "read of the expected answers", path: "tests/expected", method: http.MethodGet, wantStatus: 405},
	}
	s := newServerWithLeetCode(t, database(t), &fakeLeetCode{})
	cookie := loginAs(t, s, "student")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body any
			if tt.method != http.MethodGet {
				body = map[string]any{}
			}
			if w := request(t, s, "/problems/"+newID()+"/"+tt.path, tt.method, body, cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestTestsOwn(t *testing.T) {
	const input = "[[3, 3], 6]"
	tests := []struct {
		name       string
		role       string
		stored     int
		body       map[string]any
		wantStatus int
		wantCount  int
		wantInput  string
		wantCode   string
	}{
		{name: "student adds an own test", role: "student", stored: 0, body: map[string]any{"input": input}, wantStatus: 201, wantCount: 1, wantInput: input},
		{name: "parent adds an own test", role: "parent", stored: 0, body: map[string]any{"input": input}, wantStatus: 201, wantCount: 1, wantInput: input},
		{name: "twentieth own test", role: "student", stored: 19, body: map[string]any{"input": input}, wantStatus: 201, wantCount: 20, wantInput: input},
		{name: "twenty-first own test", role: "student", stored: 20, body: map[string]any{"input": input}, wantStatus: 409, wantCount: 20, wantCode: "too_many_tests"},
		{name: "input at 10,000 characters", role: "student", stored: 0, body: map[string]any{"input": strings.Repeat("ж", 10_000)}, wantStatus: 201, wantCount: 1, wantInput: strings.Repeat("ж", 10_000)},
		{name: "input over 10,000 characters", role: "student", stored: 0, body: map[string]any{"input": strings.Repeat("1", 10_001)}, wantStatus: 400, wantCount: 0, wantCode: "invalid_input"},
		{name: "blank input", role: "student", stored: 0, body: map[string]any{"input": " \n"}, wantStatus: 400, wantCount: 0, wantCode: "invalid_input"},
		{name: "own test with an expected answer", role: "student", stored: 0, body: map[string]any{"input": input, "expected": "[0, 1]"}, wantStatus: 400, wantCount: 0, wantCode: "invalid_input"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithLeetCode(t, pool, &fakeLeetCode{reply: testsNoQuestion})
			cookie := loginAs(t, s, "student")
			id := newID()
			exec(t, pool, testsInsertProblem, id, "two-sum")
			for range tt.stored {
				exec(t, pool, testsInsertOwn, newID(), id)
			}
			author := cookie
			if tt.role == "parent" {
				author = loginAs(t, s, "parent")
			}
			w := request(t, s, "/problems/"+id+"/tests", http.MethodPost, tt.body, author)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			created := decode[map[string]any](t, w)
			if tt.wantCode != "" && errorCode(w) != tt.wantCode {
				t.Errorf("error code = %q, want %q", errorCode(w), tt.wantCode)
			}
			if tt.wantCode == "" && created["input"] != tt.wantInput {
				t.Errorf("input = %v, want %q", created["input"], tt.wantInput)
			}
			if got := queryValue[int](t, pool, "SELECT count(*) FROM problem_tests"); got != tt.wantCount {
				t.Errorf("own tests = %d, want %d", got, tt.wantCount)
			}
			if tt.wantStatus != http.StatusCreated {
				return
			}
			// The runner numbers the own tests after the standard tests in the order of creation.
			own := decode[struct {
				Own []any `json:"own"`
			}](t, request(t, s, "/problems/"+id+"/tests", http.MethodGet, nil, cookie)).Own
			if len(own) == 0 || !reflect.DeepEqual(own[len(own)-1], created) {
				t.Errorf("last own test = %v, want the created test %v", own, created)
			}
		})
	}
}

// The cases of test/tests.test.ts for POST /problems/{id}/tests/expected without a reference of a learning set.
// 20 inputs of 10,000 Cyrillic characters are the largest request: two bytes for each character.
func TestTestsExpected(t *testing.T) {
	tests := []struct {
		name         string
		slug         string
		role         string
		inputs       any
		wantStatus   int
		wantExpected []any
		wantCode     string
	}{
		{name: "expected answers without a reference", slug: "two-sum", role: "student", inputs: []string{"[[3, 3], 6]", "1"}, wantStatus: 200, wantExpected: []any{nil, nil}},
		{name: "expected answers for the parent", slug: "two-sum", role: "parent", inputs: []string{"1"}, wantStatus: 200, wantExpected: []any{nil}},
		{name: "twenty inputs at 10,000 characters", slug: "two-sum", role: "student", inputs: slices.Repeat([]string{strings.Repeat("ж", 10_000)}, 20), wantStatus: 200, wantExpected: slices.Repeat([]any{nil}, 20)},
		{name: "twenty-one inputs", slug: "two-sum", role: "student", inputs: slices.Repeat([]string{"1"}, 21), wantStatus: 400, wantCode: "invalid_input"},
		{name: "input over 10,000 characters", slug: "two-sum", role: "student", inputs: []string{strings.Repeat("1", 10_001)}, wantStatus: 400, wantCode: "invalid_input"},
		{name: "input of a number", slug: "two-sum", role: "student", inputs: []any{1}, wantStatus: 400, wantCode: "invalid_input"},
		{name: "inputs of a string", slug: "two-sum", role: "student", inputs: "1", wantStatus: 400, wantCode: "invalid_input"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServerWithLeetCode(t, pool, &fakeLeetCode{})
			id := newID()
			exec(t, pool, testsInsertProblem, id, tt.slug)
			w := request(t, s, "/problems/"+id+"/tests/expected", http.MethodPost, map[string]any{"inputs": tt.inputs}, loginAs(t, s, tt.role))
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantCode != "" {
				if got := errorCode(w); got != tt.wantCode {
					t.Errorf("error code = %q, want %q", got, tt.wantCode)
				}
				return
			}
			if got := decode[map[string]any](t, w)["expected"]; !reflect.DeepEqual(got, tt.wantExpected) {
				t.Errorf("expected = %v, want %v", got, tt.wantExpected)
			}
		})
	}
}

func TestTestsDeleteOwn(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		target     string
		wantStatus int
		wantCount  int
	}{
		{name: "delete an own test", owner: "same", target: "stored", wantStatus: 204, wantCount: 0},
		{name: "delete an own test of another problem", owner: "other", target: "stored", wantStatus: 404, wantCount: 1},
		{name: "delete an unknown own test", owner: "same", target: "unknown", wantStatus: 404, wantCount: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServer(t, pool)
			cookie := loginAs(t, s, "student")
			id, other, stored := newID(), newID(), newID()
			exec(t, pool, testsInsertProblem, id, "two-sum")
			exec(t, pool, testsInsertProblem, other, "cf-1-a")
			exec(t, pool, testsInsertOwn, stored, id)
			problem, test := id, stored
			if tt.owner == "other" {
				problem = other
			}
			if tt.target == "unknown" {
				test = newID()
			}
			if w := request(t, s, "/problems/"+problem+"/tests/"+test, http.MethodDelete, nil, cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if got := queryValue[int](t, pool, "SELECT count(*) FROM problem_tests"); got != tt.wantCount {
				t.Errorf("own tests = %d, want %d", got, tt.wantCount)
			}
		})
	}
}

// The approach IDs of the runs and the SHA-256 (hex) of "abc".
const (
	testsX   = "11111111-1111-4111-8111-111111111111"
	testsY   = "22222222-2222-4222-8222-222222222222"
	testsZ   = "33333333-3333-4333-8333-333333333333"
	testsABC = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
)

// testsRunDocument has the approaches x and y. Both have the code "abc".
const testsRunDocument = `{"notes":"","approaches":[{"id":"11111111-1111-4111-8111-111111111111","label":"","idea":"","correctness":"","timeComplexity":"","spaceComplexity":"","edgeCases":"","mistakes":"","language":"Python","code":"abc"},{"id":"22222222-2222-4222-8222-222222222222","label":"","idea":"","correctness":"","timeComplexity":"","spaceComplexity":"","edgeCases":"","mistakes":"","language":"Python","code":"abc"}]}`

func TestTestsRecordRun(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		attempt    string
		codeHash   string
		passed     int
		total      int
		runOn      string
		wantStatus int
		wantCount  int
	}{
		{name: "student records a run", role: "student", attempt: "same", codeHash: testsABC, passed: 4, total: 5, runOn: "2026-10-12", wantStatus: 201, wantCount: 1},
		{name: "run without a check", role: "student", attempt: "same", codeHash: testsABC, passed: 0, total: 0, runOn: "2026-10-12", wantStatus: 201, wantCount: 1},
		{name: "run of 100 tests", role: "student", attempt: "same", codeHash: testsABC, passed: 100, total: 100, runOn: "2026-10-12", wantStatus: 201, wantCount: 1},
		{name: "parent records a run", role: "parent", attempt: "same", codeHash: testsABC, passed: 4, total: 5, runOn: "2026-10-12", wantStatus: 403, wantCount: 0},
		{name: "run of an attempt of another problem", role: "student", attempt: "other", codeHash: testsABC, passed: 4, total: 5, runOn: "2026-10-12", wantStatus: 404, wantCount: 0},
		{name: "more passed tests than tests", role: "student", attempt: "same", codeHash: testsABC, passed: 6, total: 5, runOn: "2026-10-12", wantStatus: 400, wantCount: 0},
		{name: "negative passed tests", role: "student", attempt: "same", codeHash: testsABC, passed: -1, total: 5, runOn: "2026-10-12", wantStatus: 400, wantCount: 0},
		{name: "run of 101 tests", role: "student", attempt: "same", codeHash: testsABC, passed: 0, total: 101, runOn: "2026-10-12", wantStatus: 400, wantCount: 0},
		{name: "short code hash", role: "student", attempt: "same", codeHash: testsABC[1:], passed: 4, total: 5, runOn: "2026-10-12", wantStatus: 400, wantCount: 0},
		{name: "run date out of the calendar", role: "student", attempt: "same", codeHash: testsABC, passed: 4, total: 5, runOn: "2026-02-30", wantStatus: 400, wantCount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServer(t, pool)
			id, otherID, attemptID := newID(), newID(), newID()
			exec(t, pool, testsInsertProblem, id, "cf-1-a")
			exec(t, pool, testsInsertProblem, otherID, "cf-1-b")
			owner := id
			if tt.attempt == "other" {
				owner = otherID
			}
			exec(t, pool, testsInsertAttempt, attemptID, owner, "draft", testsRunDocument, "now")
			body := map[string]any{"attemptId": attemptID, "approachId": testsX, "codeHash": tt.codeHash, "passed": tt.passed, "total": tt.total, "runOn": tt.runOn}
			w := request(t, s, "/problems/"+id+"/runs", http.MethodPost, body, loginAs(t, s, tt.role))
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.wantStatus == http.StatusCreated {
				createdAt, _ := decode[map[string]any](t, w)["createdAt"].(string)
				if createdAt == "" {
					t.Errorf("body = %s, want createdAt as a string", w.Body.String())
				}
				testsSameJSON(t, w.Body.Bytes(), map[string]any{"approachId": testsX, "codeHash": tt.codeHash, "passed": tt.passed, "total": tt.total, "createdAt": createdAt})
			}
			if got := queryValue[int](t, pool, "SELECT count(*) FROM test_runs"); got != tt.wantCount {
				t.Errorf("runs = %d, want %d", got, tt.wantCount)
			}
		})
	}
}

// testsRun is one run of test/tests.test.ts. other puts the run on the saved attempt. The JSON fields are the fields of a response item.
type testsRun struct {
	other      bool
	ApproachID string `json:"approachId"`
	Passed     int    `json:"passed"`
	Total      int    `json:"total"`
	CodeHash   string `json:"codeHash"`
	CreatedAt  string `json:"createdAt"`
}

func testsAt(hour int) string { return fmt.Sprintf("2026-10-12T%d:00:00.000Z", hour) }

// The latest run of each approach of the draft. The draft is created at 12:00, the saved attempt at 10:00.
func TestTestsLatestRuns(t *testing.T) {
	otherHash := strings.Repeat("f", 64)
	tests := []struct {
		name   string
		runs   []testsRun
		latest []testsRun
	}{
		{name: "no runs", runs: nil, latest: []testsRun{}},
		{
			name:   "latest run passes with the current code",
			runs:   []testsRun{{ApproachID: testsX, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}},
			latest: []testsRun{{ApproachID: testsX, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}},
		},
		{
			name:   "latest run passes with other code",
			runs:   []testsRun{{ApproachID: testsX, Passed: 5, Total: 5, CodeHash: otherHash, CreatedAt: testsAt(10)}},
			latest: []testsRun{{ApproachID: testsX, Passed: 5, Total: 5, CodeHash: otherHash, CreatedAt: testsAt(10)}},
		},
		{
			name:   "older run passes and latest run fails",
			runs:   []testsRun{{ApproachID: testsX, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}, {ApproachID: testsX, Passed: 4, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(11)}},
			latest: []testsRun{{ApproachID: testsX, Passed: 4, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(11)}},
		},
		{
			name:   "older run fails and latest run passes",
			runs:   []testsRun{{ApproachID: testsX, Passed: 4, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}, {ApproachID: testsX, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(11)}},
			latest: []testsRun{{ApproachID: testsX, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(11)}},
		},
		{
			name:   "latest run without a check",
			runs:   []testsRun{{ApproachID: testsX, Passed: 0, Total: 0, CodeHash: testsABC, CreatedAt: testsAt(10)}},
			latest: []testsRun{{ApproachID: testsX, Passed: 0, Total: 0, CodeHash: testsABC, CreatedAt: testsAt(10)}},
		},
		{
			name:   "second approach passes",
			runs:   []testsRun{{ApproachID: testsX, Passed: 4, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}, {ApproachID: testsY, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}},
			latest: []testsRun{{ApproachID: testsX, Passed: 4, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}, {ApproachID: testsY, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}},
		},
		{
			name:   "passed run of an approach outside the attempt",
			runs:   []testsRun{{ApproachID: testsZ, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}},
			latest: []testsRun{{ApproachID: testsZ, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}},
		},
		{
			name:   "passed run of another attempt",
			runs:   []testsRun{{other: true, ApproachID: testsX, Passed: 5, Total: 5, CodeHash: testsABC, CreatedAt: testsAt(10)}},
			latest: []testsRun{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServer(t, pool)
			cookie := loginAs(t, s, "student")
			id, attemptID, otherID := newID(), newID(), newID()
			exec(t, pool, testsInsertProblem, id, "cf-1-a")
			exec(t, pool, testsInsertAttempt, attemptID, id, "draft", testsRunDocument, testsAt(12))
			exec(t, pool, testsInsertAttempt, otherID, id, "saved", testsRunDocument, testsAt(10))
			for _, run := range tt.runs {
				owner := attemptID
				if run.other {
					owner = otherID
				}
				exec(t, pool, testsInsertRun, newID(), id, owner, run.ApproachID, run.CodeHash, run.Passed, run.Total, "2026-10-12", run.CreatedAt)
			}
			w := request(t, s, "/problems/"+id+"/runs?attemptId="+attemptID, http.MethodGet, nil, cookie)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			got := decode[[]map[string]any](t, w)
			slices.SortFunc(got, func(a, b map[string]any) int {
				return strings.Compare(fmt.Sprint(a["approachId"]), fmt.Sprint(b["approachId"]))
			})
			sorted, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			testsSameJSON(t, sorted, tt.latest)
		})
	}
}

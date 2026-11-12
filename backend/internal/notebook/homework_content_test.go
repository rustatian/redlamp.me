package notebook

import (
	"net/http"
	"reflect"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The cases of test/homework-content.test.ts. The case that reads /homework/{id}/sets is in homework_dashboard_test.go.

// homeworkSets is the content of src/leetcode/sets/content.ts that the homework tests need, in the same order.
// complexity has its only chapter, which has no tasks. two-pointers has chapter 1 with its five tasks.
var homeworkSets = []LearningSet{
	{Slug: "complexity", Title: "Как понять O(n): время и память", Chapters: []Chapter{{Title: "Как понять O(n): время и память", Tasks: []SetTask{}}}},
	{Slug: "two-pointers", Title: "Two Pointers", Chapters: []Chapter{{Title: "Навстречу", Tasks: []SetTask{
		{URL: "https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/", Number: new(int64(167)), Title: "Two Sum II - Input Array Is Sorted", Difficulty: "medium", Topics: []string{"Array", "Two Pointers"}, Note: "n - длина уже отсортированного Array, числа ограниченного размера. Время: не более n-1 вычислений суммы, O(n). Память: O(1) дополнительной, ответ - два индекса с единицы. При малой сумме исключайте левое число, при большой - правое. Brute Force проверяет до n(n-1)/2 пар; Hash Table требует O(u) записей для u разных значений."},
		{URL: "https://leetcode.com/problems/valid-palindrome/", Number: new(int64(125)), Title: "Valid Palindrome", Difficulty: "easy", Topics: []string{"String", "Two Pointers"}, Note: "n - длина String из ASCII. Пропускайте символы, кроме букв и цифр, и сравнивайте края без учёта регистра. Время: индексы суммарно проходят не более n позиций, O(n). Память: O(1) дополнительной и Boolean-ответ. Фильтрация с созданием String и её разворот тоже требуют O(n) времени, но O(n) дополнительной памяти."},
		{URL: "https://leetcode.com/problems/max-number-of-k-sum-pairs/", Number: new(int64(1679)), Title: "Max Number of K-Sum Pairs", Difficulty: "medium", Topics: []string{"Array", "Two Pointers", "Sorting"}, Note: "n - длина Array, u - число разных значений; числа ограниченного размера. Sorting допустим: нужны количество пар, а не исходные индексы. Время: O(n log n) на Sorting плюс O(n) на проход. Память: O(n) дополнительной в худшем случае Python Sorting; ответ - счётчик. Hash Table частот даёт ожидаемое O(n) время при хороших хешах и O(u) памяти. При сумме k используйте оба числа."},
		{URL: "https://leetcode.com/problems/container-with-most-water/", Number: new(int64(11)), Title: "Container With Most Water", Difficulty: "medium", Topics: []string{"Array", "Two Pointers"}, Note: "n - число линий, высоты ограниченного размера. Площадь равна min(h[l], h[r]) · (r-l). Исключайте меньшую высоту: с ней более узкая пара не улучшит результат. Время: n-1 вычислений площади и сужений, O(n). Память: O(1) дополнительной и одно число в ответе. Brute Force проверяет n(n-1)/2 пар за O(n²) времени с O(1) памятью."},
		{URL: "https://leetcode.com/problems/boats-to-save-people/", Number: new(int64(881)), Title: "Boats to Save People", Difficulty: "medium", Topics: []string{"Array", "Two Pointers", "Sorting"}, Note: "n - число людей, целые веса ограничены; в лодке не больше двух. После Sorting самый тяжёлый едет с самым лёгким, если они помещаются, иначе один. Время: O(n log n) на Sorting и O(n) на посадку. Память: O(n) дополнительной в худшем случае Python Sorting; ответ - счётчик лодок. Для весов 1..W частоты дают O(n+W) времени и O(W) памяти."},
	}}}},
}

// homeworkContentFixture is the state that beforeEach of test/homework-content.test.ts makes.
type homeworkContentFixture struct {
	pool                        *pgxpool.Pool
	s                           *Server
	parent, student, homeworkID string
}

// as returns the session cookie of the role.
func (f homeworkContentFixture) as(role string) string {
	if role == "parent" {
		return f.parent
	}
	return f.student
}

// homeworkContentSetup adds the empty homework "Chapter practice" as the parent.
func homeworkContentSetup(t *testing.T) homeworkContentFixture {
	t.Helper()
	pool := database(t)
	s := newServerWithSets(t, pool, homeworkSets)
	f := homeworkContentFixture{pool: pool, s: s, parent: loginAs(t, s, "parent"), student: loginAs(t, s, "student")}
	w := request(t, s, "/homework", http.MethodPost, map[string]any{"title": "Chapter practice", "instructions": "", "dueDate": nil}, f.parent)
	if w.Code != http.StatusCreated {
		t.Fatalf("new homework: status %d, body %s", w.Code, w.Body.String())
	}
	f.homeworkID = decode[struct {
		ID string `json:"id"`
	}](t, w).ID
	return f
}

func TestHomeworkContentRejected(t *testing.T) {
	tests := []struct {
		name   string
		role   string
		body   map[string]any
		status int
	}{
		{name: "student cannot assign content", role: "student", body: map[string]any{"kind": "chapter", "setSlug": "complexity", "chapter": 1}, status: 403},
		{name: "unknown set", role: "parent", body: map[string]any{"kind": "chapter", "setSlug": "missing", "chapter": 1}, status: 404},
		{name: "missing chapter", role: "parent", body: map[string]any{"kind": "chapter", "setSlug": "complexity", "chapter": 2}, status: 404},
		{name: "chapter zero", role: "parent", body: map[string]any{"kind": "chapter", "setSlug": "complexity", "chapter": 0}, status: 400},
		{name: "fractional chapter", role: "parent", body: map[string]any{"kind": "chapter", "setSlug": "complexity", "chapter": 1.5}, status: 400},
		{name: "unknown target type", role: "parent", body: map[string]any{"kind": "problem", "setSlug": "complexity", "chapter": 1}, status: 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkContentSetup(t)
			if w := request(t, f.s, "/homework/"+f.homeworkID+"/content", http.MethodPost, tt.body, f.as(tt.role)); w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			problemsMatchBody(t, decode[any](t, request(t, f.s, "/homework/"+f.homeworkID, http.MethodGet, nil, f.parent)), map[string]any{"tasks": []any{}})
		})
	}
}

// taskState is "" for no task of the problem in other homework.
func TestHomeworkContentSkippedProblem(t *testing.T) {
	tests := []struct {
		name      string
		archived  bool
		solved    bool
		taskState string
		reason    string
	}{
		{name: "archived problem is reported", archived: true, solved: false, taskState: "", reason: "Задача в архиве."},
		{name: "solved problem is reported", archived: false, solved: true, taskState: "", reason: "Задача уже выполнена."},
		{name: "completed homework problem is reported", archived: false, solved: false, taskState: "completed", reason: "Задача уже выполнена."},
		{name: "active problem stays in its homework", archived: false, solved: false, taskState: "assigned", reason: "У задачи уже есть активное задание."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkContentSetup(t)
			problem := addProblem(t, f.s, f.parent, map[string]any{"url": "https://leetcode.com/problems/valid-palindrome/", "number": 125, "title": "My palindrome", "difficulty": "easy", "topics": []string{"String"}, "summary": "My notes"})
			if tt.archived {
				request(t, f.s, "/problems/"+problem+"/archive", http.MethodPost, map[string]any{"version": 1}, f.parent)
			}
			if tt.solved {
				exec(t, f.pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'saved','{}','','accepted','independent','now','now')", newID(), problem)
			}
			if tt.taskState != "" {
				other := decode[struct {
					ID string `json:"id"`
				}](t, request(t, f.s, "/homework", http.MethodPost, map[string]any{"title": "Other homework", "instructions": ""}, f.parent)).ID
				exec(t, f.pool, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)", newID(), other, problem, tt.taskState, "now", "now")
			}
			w := request(t, f.s, "/homework/"+f.homeworkID+"/content", http.MethodPost, map[string]any{"kind": "chapter", "setSlug": "two-pointers", "chapter": 1}, f.parent)
			if w.Code != http.StatusCreated {
				t.Errorf("status = %d, want 201; body %s", w.Code, w.Body.String())
			}
			want := map[string]any{"added": float64(6), "skipped": []any{map[string]any{"title": "Valid Palindrome", "reason": tt.reason}}}
			if got := decode[map[string]any](t, w); !reflect.DeepEqual(got, want) {
				t.Errorf("body = %v, want %v", got, want)
			}
			tasks := homeworkTasks(t, f.s, f.homeworkID, f.parent)
			if len(tasks) != 6 {
				t.Errorf("homework has %d tasks, want 6", len(tasks))
			}
			for _, task := range tasks {
				if task["problemId"] == problem {
					t.Errorf("task %v has the skipped problem", task)
				}
			}
			problemsMatchBody(t, decode[any](t, request(t, f.s, "/problems/"+problem, http.MethodGet, nil, f.parent)), map[string]any{"title": "My palindrome", "summary": "My notes"})
		})
	}
}

func TestHomeworkContentComplete(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		version int
		kind    string
		status  int
		state   string
	}{
		{name: "student completes chapter reading", role: "student", version: 1, kind: "chapter", status: 200, state: "completed"},
		{name: "stale reading completion is rejected", role: "student", version: 2, kind: "chapter", status: 409, state: "assigned"},
		{name: "parent cannot complete student reading", role: "parent", version: 1, kind: "chapter", status: 403, state: "assigned"},
		{name: "quiz requires a passing attempt", role: "student", version: 1, kind: "quiz", status: 409, state: "assigned"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkContentSetup(t)
			request(t, f.s, "/homework/"+f.homeworkID+"/content", http.MethodPost, map[string]any{"kind": "chapter", "setSlug": "complexity", "chapter": 1}, f.parent)
			items := decode[homeworkPage](t, request(t, f.s, "/tasks?kind="+tt.kind+"&setSlug=complexity&chapter=1", http.MethodGet, nil, f.student)).Items
			if len(items) != 1 {
				t.Fatalf("tasks of kind %s = %v, want one task", tt.kind, items)
			}
			id, _ := items[0]["id"].(string)
			if w := request(t, f.s, "/tasks/"+id+"/complete", http.MethodPost, map[string]any{"version": tt.version}, f.as(tt.role)); w.Code != tt.status {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
			var task map[string]any
			for _, item := range homeworkTasks(t, f.s, f.homeworkID, f.parent) {
				if item["id"] == id {
					task = item
				}
			}
			problemsMatchBody(t, task, map[string]any{"state": tt.state})
			completedAt, ok := task["completedAt"]
			if tt.state == "completed" {
				if text, _ := completedAt.(string); text == "" {
					t.Errorf("completedAt = %v, want a time", completedAt)
				}
			} else if !ok || completedAt != nil {
				t.Errorf("completedAt = %v, want null", completedAt)
			}
			completed := 0
			if tt.state == "completed" {
				completed = 1
			}
			summary := decode[homeworkPage](t, request(t, f.s, "/homework", http.MethodGet, nil, f.parent)).Items
			problemsMatchBody(t, homeworkFirst(t, summary), map[string]any{"taskCount": 2, "completedCount": completed})
		})
	}
}

func TestHomeworkContentRollsBackFailedChapter(t *testing.T) {
	f := homeworkContentSetup(t)
	// PostgreSQL allows no subquery in the WHEN clause of a trigger, so the function body checks the problem.
	exec(t, f.pool, `CREATE FUNCTION fail_assignment() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.problem_id IN (SELECT id FROM problems WHERE slug='boats-to-save-people') THEN RAISE EXCEPTION 'test storage failure'; END IF; RETURN NEW; END $$`)
	exec(t, f.pool, "CREATE TRIGGER fail_assignment BEFORE INSERT ON homework_tasks FOR EACH ROW EXECUTE FUNCTION fail_assignment()")
	if w := request(t, f.s, "/homework/"+f.homeworkID+"/content", http.MethodPost, map[string]any{"kind": "chapter", "setSlug": "two-pointers", "chapter": 1}, f.parent); w.Code >= 200 && w.Code < 300 {
		t.Errorf("status = %d, want an error status", w.Code)
	}
	problemsMatchBody(t, decode[any](t, request(t, f.s, "/homework/"+f.homeworkID, http.MethodGet, nil, f.student)), map[string]any{"tasks": []any{}})
	problemsMatchBody(t, decode[any](t, request(t, f.s, "/problems", http.MethodGet, nil, f.student)), map[string]any{"items": []any{}})
}

func TestHomeworkContentAdded(t *testing.T) {
	tests := []struct {
		name     string
		kind     string
		setSlug  string
		count    int
		problems []string
	}{
		{name: "chapter includes its reading, quiz, and five problems", kind: "chapter", setSlug: "two-pointers", count: 7, problems: []string{"Boats to Save People", "Container With Most Water", "Max Number of K-Sum Pairs", "Two Sum II - Input Array Is Sorted", "Valid Palindrome"}},
		{name: "chapter without problems includes reading and quiz", kind: "chapter", setSlug: "complexity", count: 2, problems: []string{}},
		{name: "standalone quiz includes no reading or problems", kind: "quiz", setSlug: "two-pointers", count: 1, problems: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkContentSetup(t)
			body := map[string]any{"kind": tt.kind, "setSlug": tt.setSlug, "chapter": 1}
			w := request(t, f.s, "/homework/"+f.homeworkID+"/content", http.MethodPost, body, f.parent)
			if w.Code != http.StatusCreated {
				t.Errorf("status = %d, want 201; body %s", w.Code, w.Body.String())
			}
			if got, want := decode[map[string]any](t, w), map[string]any{"added": float64(tt.count), "skipped": []any{}}; !reflect.DeepEqual(got, want) {
				t.Errorf("body = %v, want %v", got, want)
			}
			tasks := homeworkTasks(t, f.s, f.homeworkID, f.student)
			if len(tasks) != tt.count {
				t.Errorf("homework has %d tasks, want %d", len(tasks), tt.count)
			}
			titles := []string{}
			first := map[string]any{}
			for _, task := range tasks {
				kind, _ := task["kind"].(string)
				if kind == "problem" {
					title, _ := task["title"].(string)
					titles = append(titles, title)
				}
				if _, ok := first[kind]; !ok {
					first[kind] = task
				}
			}
			slices.Sort(titles)
			if !slices.Equal(titles, tt.problems) {
				t.Errorf("problems = %q, want %q", titles, tt.problems)
			}
			problemsMatchBody(t, first["quiz"], map[string]any{"href": "/leetcode/quiz/" + tt.setSlug + "-1"})
			if tt.kind == "chapter" {
				problemsMatchBody(t, first["chapter"], map[string]any{"href": "/leetcode/sets/" + tt.setSlug + "/1"})
			}
			problemsMatchBody(t, decode[any](t, request(t, f.s, "/homework/"+f.homeworkID+"/content", http.MethodPost, body, f.parent)), map[string]any{"added": 0})
			if got := homeworkTasks(t, f.s, f.homeworkID, f.parent); len(got) != tt.count {
				t.Errorf("after the second request the homework has %d tasks, want %d", len(got), tt.count)
			}
		})
	}
}

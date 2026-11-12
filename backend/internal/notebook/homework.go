package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	postponedSQL    = "(SELECT to_date FROM postponements WHERE homework_id=h.id)"
	homeworkColumns = "h.id,h.title,h.instructions,h.due_date," + postponedSQL + ",h.version,h.created_at,h.updated_at"
	// unreadSQL is true for a task with feedback of the other role after the last read. $1 is the role of the viewer.
	unreadSQL   = "EXISTS(SELECT 1 FROM feedback f WHERE f.problem_id=t.problem_id AND f.author<>$1 AND f.created_at>coalesce((SELECT read_at FROM discussion_reads r WHERE r.role=$1 AND r.problem_id=t.problem_id),''))"
	taskColumns = "t.id,t.homework_id,h.title,h.instructions,h.due_date," + postponedSQL + ",t.problem_id,p.title,p.difficulty,t.kind,t.set_slug,t.chapter,t.completed_at,t.state,t.current_submission_id,s.attempt_id,t.version,t.created_at,t.updated_at," + unreadSQL
	taskTables  = "homework_tasks t JOIN homework h ON h.id=t.homework_id LEFT JOIN problems p ON p.id=t.problem_id LEFT JOIN submissions s ON s.id=t.current_submission_id"
	// onTimeSQL counts the homework tasks that were completed or submitted on or before the due date. The first 10 characters of a timestamp are its UTC date.
	onTimeSQL = "SELECT count(*) FROM homework_tasks t JOIN homework h ON h.id=t.homework_id WHERE t.state<>'cancelled' AND (substr(t.completed_at,1,10)<=h.due_date OR EXISTS(SELECT 1 FROM submissions s WHERE s.task_id=t.id AND substr(s.created_at,1,10)<=h.due_date))"
	// awardOnTimeSQL records the award before a change of a due date or a cancel removes its condition.
	awardOnTimeSQL  = "INSERT INTO earned_badges SELECT 'on-time-five' WHERE (" + onTimeSQL + ")>=5 ON CONFLICT DO NOTHING"
	feedbackColumns = "id,problem_id,task_id,submission_id,author,kind,body,created_at"
)

// Homework is a homework in the form of the API.
type Homework struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Instructions string  `json:"instructions"`
	DueDate      *string `json:"dueDate"`
	PostponedTo  *string `json:"postponedTo"`
	Version      int64   `json:"version"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type homeworkSummary struct {
	Homework
	TaskCount      int64 `json:"taskCount"`
	CompletedCount int64 `json:"completedCount"`
	SubmittedCount int64 `json:"submittedCount"`
	RequestedCount int64 `json:"requestedCount"`
	Unread         bool  `json:"unread"`
}

// HomeworkTask is a homework task in the form of the API.
type HomeworkTask struct {
	ID                  string  `json:"id"`
	HomeworkID          string  `json:"homeworkId"`
	HomeworkTitle       string  `json:"homeworkTitle"`
	Instructions        string  `json:"instructions"`
	DueDate             *string `json:"dueDate"`
	PostponedTo         *string `json:"postponedTo"`
	ProblemID           *string `json:"problemId"`
	Title               *string `json:"title"`
	Difficulty          *string `json:"difficulty"`
	Kind                string  `json:"kind"`
	SetSlug             *string `json:"setSlug"`
	Chapter             *int64  `json:"chapter"`
	CompletedAt         *string `json:"completedAt"`
	State               string  `json:"state"`
	SubmissionID        *string `json:"submissionId"`
	SubmissionAttemptID *string `json:"submissionAttemptId"`
	Version             int64   `json:"version"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
	Unread              bool    `json:"unread"`
	Href                string  `json:"href"`
}

type feedback struct {
	ID           string  `json:"id"`
	ProblemID    string  `json:"problemId"`
	TaskID       *string `json:"taskId"`
	SubmissionID *string `json:"submissionId"`
	Author       string  `json:"author"`
	Kind         string  `json:"kind"`
	Body         string  `json:"body"`
	CreatedAt    string  `json:"createdAt"`
}

func (h *Homework) fields() []any {
	return []any{&h.ID, &h.Title, &h.Instructions, &h.DueDate, &h.PostponedTo, &h.Version, &h.CreatedAt, &h.UpdatedAt}
}

func getHomework(ctx context.Context, q querier, id string) (Homework, error) {
	var h Homework
	if _, err := uuid(id); err != nil {
		return h, err
	}
	err := q.QueryRow(ctx, "SELECT "+homeworkColumns+" FROM homework h WHERE h.id=$1", id).Scan(h.fields()...)
	if errors.Is(err, pgx.ErrNoRows) {
		return h, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой домашней работы нет."}
	}
	return h, err
}

// contentTitle is the title of a chapter task or a quiz task.
func contentTitle(kind string, number int64, chapter Chapter) string {
	word := "Глава"
	if kind == "quiz" {
		word = "Квиз"
	}
	return fmt.Sprintf("%s %d. %s", word, number, chapter.Title)
}

// scanTask reads the columns of taskColumns and adds the link of the task. A chapter task and a quiz task get their title from the learning set.
func (s *Server) scanTask(row pgx.Row) (HomeworkTask, error) {
	var t HomeworkTask
	if err := row.Scan(&t.ID, &t.HomeworkID, &t.HomeworkTitle, &t.Instructions, &t.DueDate, &t.PostponedTo, &t.ProblemID, &t.Title, &t.Difficulty, &t.Kind, &t.SetSlug, &t.Chapter, &t.CompletedAt, &t.State, &t.SubmissionID, &t.SubmissionAttemptID, &t.Version, &t.CreatedAt, &t.UpdatedAt, &t.Unread); err != nil {
		return t, err
	}
	if t.Kind == "problem" {
		t.Href = "/leetcode/problems/" + *t.ProblemID
		return t, nil
	}
	chapter, ok := s.chapter(*t.SetSlug, *t.Chapter)
	if !ok {
		return t, fmt.Errorf("the homework task %s names chapter %d of the learning set %q, which is not in the content", t.ID, *t.Chapter, *t.SetSlug)
	}
	title := contentTitle(t.Kind, *t.Chapter, chapter)
	t.Title = &title
	t.Href = fmt.Sprintf("/leetcode/sets/%s/%d", *t.SetSlug, *t.Chapter)
	if t.Kind == "quiz" {
		t.Href = fmt.Sprintf("/leetcode/quiz/%s-%d", *t.SetSlug, *t.Chapter)
	}
	return t, nil
}

// tasks runs a query of taskColumns. The first argument is the role of the viewer.
func (s *Server) tasks(ctx context.Context, q querier, sql string, args ...any) ([]HomeworkTask, error) {
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	tasks, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (HomeworkTask, error) { return s.scanTask(row) })
	if err != nil {
		return nil, err
	}
	return append([]HomeworkTask{}, tasks...), nil
}

func (s *Server) getTask(ctx context.Context, q querier, id, role string) (HomeworkTask, error) {
	if _, err := uuid(id); err != nil {
		return HomeworkTask{}, err
	}
	t, err := s.scanTask(q.QueryRow(ctx, "SELECT "+taskColumns+" FROM "+taskTables+" WHERE t.id=$2", role, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return t, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такого задания нет."}
	}
	return t, err
}

var (
	contentPath    = regexp.MustCompile(`^/leetcode/api/homework/([^/]+)/content$`)
	discussionPath = regexp.MustCompile(`^/leetcode/api/problems/([^/]+)/feedback$`)
	homeworkPath   = regexp.MustCompile(`^/leetcode/api/homework(?:/([^/]+)(/postpone)?)?$`)
	tasksPath      = regexp.MustCompile(`^/leetcode/api/tasks(?:/([^/]+)/(start|review|cancel|complete))?$`)
)

func (s *Server) handleHomework(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	path := r.URL.Path
	if match := contentPath.FindStringSubmatch(path); match != nil {
		return true, s.assignContent(w, r, user, match[1])
	}
	if match := discussionPath.FindStringSubmatch(path); match != nil {
		return true, s.discussion(w, r, user, match[1])
	}
	if match := homeworkPath.FindStringSubmatch(path); match != nil {
		if match[1] != "" {
			if _, err := uuid(match[1]); err != nil {
				return true, err
			}
		}
		if match[2] != "" {
			return true, s.postpone(w, r, user, match[1])
		}
		return true, s.homework(w, r, user, match[1])
	}
	if match := tasksPath.FindStringSubmatch(path); match != nil {
		return true, s.homeworkTasks(w, r, user, match[1], match[2])
	}
	return false, nil
}

// readObject reads a JSON request of at most 32,768 bytes and checks that it is an object with no key outside keys.
func readObject(r *http.Request, keys ...string) (map[string]any, error) {
	body, err := readJSON(r, 32_768)
	if err != nil {
		return nil, err
	}
	return object(body, keys...)
}

// readVersion reads a request that has only the version of a record.
func readVersion(r *http.Request) (int64, error) {
	data, err := readObject(r, "version")
	if err != nil {
		return 0, err
	}
	return integer(data["version"], 1, maxSafeInteger)
}

// assignContent adds a chapter of a learning set to a homework: the chapter task, its quiz task and a task for each problem of the chapter.
// The kind "quiz" adds only the quiz task. The answer lists each entry that was not added, with the reason.
func (s *Server) assignContent(w http.ResponseWriter, r *http.Request, user *SessionUser, homeworkID string) error {
	ctx := r.Context()
	if err := methods(r, http.MethodPost); err != nil {
		return err
	}
	if err := requireRole(user, "parent"); err != nil {
		return err
	}
	if _, err := uuid(homeworkID); err != nil {
		return err
	}
	data, err := readObject(r, "kind", "setSlug", "chapter")
	if err != nil {
		return err
	}
	kind, err := choice(data["kind"], "chapter", "quiz")
	if err != nil {
		return err
	}
	setSlug, err := text(data["setSlug"], 100, true)
	if err != nil {
		return err
	}
	number, err := integer(data["chapter"], 1, 99)
	if err != nil {
		return err
	}
	chapter, ok := s.chapter(setSlug, number)
	if !ok {
		return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой главы нет."}
	}
	if _, err := getHomework(ctx, s.db, homeworkID); err != nil {
		return err
	}
	kinds := []string{"quiz"}
	var problems []SetTask
	if kind == "chapter" {
		kinds, problems = []string{"chapter", "quiz"}, chapter.Tasks
	}
	slugs := make([]string, len(problems))
	for i, problem := range problems {
		if slugs[i], err = problemSlug(problem.URL); err != nil {
			return fmt.Errorf("the task %q of the learning set %q: %w", problem.Title, setSlug, err)
		}
	}
	type skippedEntry struct {
		Title  string `json:"title"`
		Reason string `json:"reason"`
	}
	now := isoNow()
	var added int
	var skipped []skippedEntry
	err = s.tx(ctx, func(tx pgx.Tx) error {
		added, skipped = 0, []skippedEntry{}
		for _, kind := range kinds {
			// The unique index one_active_content allows one active task for each chapter and each quiz deck.
			tag, err := tx.Exec(ctx, "INSERT INTO homework_tasks (id,homework_id,kind,set_slug,chapter,state,created_at,updated_at) SELECT $1,$2,$3,$4,$5,'assigned',$6,$6 WHERE NOT EXISTS(SELECT 1 FROM homework_tasks WHERE homework_id=$2 AND kind=$3 AND set_slug=$4 AND chapter=$5 AND state<>'cancelled') ON CONFLICT DO NOTHING",
				newID(), homeworkID, kind, setSlug, number, now)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				skipped = append(skipped, skippedEntry{Title: contentTitle(kind, number, chapter), Reason: "Уже добавлено в домашнюю работу."})
				continue
			}
			added++
		}
		var notAdded []int
		for i, problem := range problems {
			topics := problem.Topics
			if topics == nil {
				topics = []string{}
			}
			names, err := json.Marshal(topics)
			if err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "INSERT INTO problems (id,slug,number,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9) ON CONFLICT (slug) DO NOTHING",
				newID(), slugs[i], problem.Number, problem.Title, problem.Difficulty, string(names), problem.Note, searchOf(problem.Title, problem.Number, topics, problem.Note), now); err != nil {
				return err
			}
			// The write of the problem row makes this check conflict with an archive or an accepted attempt of the same problem.
			if _, err := tx.Exec(ctx, "UPDATE problems SET version = version WHERE slug = $1", slugs[i]); err != nil {
				return err
			}
			tag, err := tx.Exec(ctx, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) SELECT $1,$2,p.id,'assigned',$3,$3 FROM problems p WHERE p.slug=$4 AND "+assignableSQL+" ON CONFLICT DO NOTHING",
				newID(), homeworkID, now, slugs[i])
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				notAdded = append(notAdded, i)
				continue
			}
			added++
		}
		for _, i := range notAdded {
			var reason string
			if err := tx.QueryRow(ctx, "SELECT CASE WHEN p.archived_at IS NOT NULL THEN 'Задача в архиве.' WHEN "+solvedSQL+" OR EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state='completed') THEN 'Задача уже выполнена.' ELSE 'У задачи уже есть активное задание.' END FROM problems p WHERE p.slug=$1", slugs[i]).Scan(&reason); err != nil {
				return err
			}
			skipped = append(skipped, skippedEntry{Title: problems[i].Title, Reason: reason})
		}
		return nil
	})
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusCreated, struct {
		Added   int            `json:"added"`
		Skipped []skippedEntry `json:"skipped"`
	}{added, skipped})
}

// discussion lists the feedback of a problem. A read of the first page marks the discussion as read for the role of the viewer.
func (s *Server) discussion(w http.ResponseWriter, r *http.Request, user *SessionUser, problemID string) error {
	ctx := r.Context()
	if err := methods(r, http.MethodGet); err != nil {
		return err
	}
	if _, err := getProblem(ctx, s.db, problemID); err != nil {
		return err
	}
	limit, offset, err := pagination(r.URL.Query())
	if err != nil {
		return err
	}
	var items []feedback
	err = s.tx(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, "SELECT "+feedbackColumns+" FROM feedback WHERE problem_id=$1 ORDER BY created_at DESC,id DESC LIMIT $2 OFFSET $3", problemID, limit+1, offset)
		if err != nil {
			return err
		}
		items, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (feedback, error) {
			var f feedback
			return f, row.Scan(&f.ID, &f.ProblemID, &f.TaskID, &f.SubmissionID, &f.Author, &f.Kind, &f.Body, &f.CreatedAt)
		})
		if err != nil || offset != 0 {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO discussion_reads (role,problem_id,read_at) SELECT $1,problem_id,max(created_at) FROM feedback WHERE problem_id=$2 GROUP BY problem_id ON CONFLICT (role,problem_id) DO UPDATE SET read_at=excluded.read_at", user.Role, problemID)
		return err
	})
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, page(items, limit, offset))
}

// postpone moves the due date of a homework by one day. The student can do it one time for each homework, on or before the due date.
func (s *Server) postpone(w http.ResponseWriter, r *http.Request, user *SessionUser, id string) error {
	ctx := r.Context()
	if err := methods(r, http.MethodPost); err != nil {
		return err
	}
	if err := requireRole(user, "student"); err != nil {
		return err
	}
	now := isoNow()
	version, err := readVersion(r)
	if err != nil {
		return err
	}
	err = s.tx(ctx, func(tx pgx.Tx) error {
		homework, err := getHomework(ctx, tx, id)
		if err != nil {
			return err
		}
		if homework.PostponedTo != nil {
			return &HTTPError{Status: http.StatusConflict, Code: "postponed", Message: "Сила Бомжа уже использована для этой домашней работы."}
		}
		// The UTC date, as in the rule for late tasks.
		if homework.DueDate == nil || now[:10] > *homework.DueDate {
			return invalid("Перенести можно только срок, который ещё не прошёл.")
		}
		due, err := time.Parse(time.DateOnly, *homework.DueDate)
		if err != nil {
			return err
		}
		next := due.AddDate(0, 0, 1).Format(time.DateOnly)
		// The primary key of postponements allows one row for each homework.
		tag, err := tx.Exec(ctx, "INSERT INTO postponements SELECT id,due_date,$1,$2 FROM homework WHERE id=$3 AND version=$4 ON CONFLICT DO NOTHING", next, now, id, version)
		if err != nil {
			return err
		}
		if err := changed(tag); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "UPDATE homework SET due_date=$1,version=version+1,updated_at=$2 WHERE id=$3 AND version=$4", next, now, id, version)
		return err
	})
	if err != nil {
		return err
	}
	homework, err := getHomework(ctx, s.db, id)
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, homework)
}

func (s *Server) homework(w http.ResponseWriter, r *http.Request, user *SessionUser, id string) error {
	ctx, query := r.Context(), r.URL.Query()
	allowed := []string{http.MethodGet, http.MethodPost}
	if id != "" {
		allowed = []string{http.MethodGet, http.MethodPatch}
	}
	if err := methods(r, allowed...); err != nil {
		return err
	}
	if r.Method == http.MethodGet && id != "" {
		homework, err := getHomework(ctx, s.db, id)
		if err != nil {
			return err
		}
		tasks, err := s.tasks(ctx, s.db, "SELECT "+taskColumns+" FROM "+taskTables+" WHERE t.homework_id=$2 ORDER BY t.created_at,t.id", user.Role, id)
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, struct {
			Homework Homework       `json:"homework"`
			Tasks    []HomeworkTask `json:"tasks"`
		}{homework, tasks})
	}
	if r.Method == http.MethodGet {
		limit, offset, err := pagination(query)
		if err != nil {
			return err
		}
		q, err := text(query.Get("q"), 200, false)
		if err != nil {
			return err
		}
		count := func(condition string) string {
			return "coalesce(sum(CASE WHEN " + condition + " THEN 1 ELSE 0 END),0)"
		}
		having := ""
		if query.Get("active") == "1" {
			having = " HAVING count(t.id)=0 OR " + count("t.state NOT IN ('completed','cancelled')") + ">0"
		}
		// lower() changes only the ASCII letters under the C collation.
		rows, err := s.db.Query(ctx, "SELECT "+homeworkColumns+","+count("t.state<>'cancelled'")+","+count("t.state='completed'")+","+count("t.state='submitted'")+","+count("t.state='changes_requested'")+","+count(unreadSQL)+">0"+
			" FROM homework h LEFT JOIN homework_tasks t ON t.homework_id=h.id WHERE strpos(lower(h.title),lower($2))>0 GROUP BY h.id,h.title,h.instructions,h.due_date,h.version,h.created_at,h.updated_at"+having+
			" ORDER BY ("+count("t.state='submitted'")+">0) DESC,h.created_at DESC,h.id LIMIT $3 OFFSET $4", user.Role, q, limit+1, offset)
		if err != nil {
			return err
		}
		items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (homeworkSummary, error) {
			var h homeworkSummary
			return h, row.Scan(append(h.fields(), &h.TaskCount, &h.CompletedCount, &h.SubmittedCount, &h.RequestedCount, &h.Unread)...)
		})
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, page(items, limit, offset))
	}
	if err := requireRole(user, "parent"); err != nil {
		return err
	}
	keys := []string{"title", "instructions", "dueDate"}
	if id != "" {
		keys = append(keys, "version")
	}
	data, err := readObject(r, keys...)
	if err != nil {
		return err
	}
	title, err := text(data["title"], 200, true)
	if err != nil {
		return err
	}
	instructions, err := text(data["instructions"], 8000, false)
	if err != nil {
		return err
	}
	dueDate, err := calendarDate(data["dueDate"])
	if err != nil {
		return err
	}
	title, instructions = strings.TrimSpace(title), strings.TrimSpace(instructions)
	now, status := isoNow(), http.StatusOK
	if id == "" {
		id, status = newID(), http.StatusCreated
		if _, err := s.db.Exec(ctx, "INSERT INTO homework (id,title,instructions,due_date,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$5)", id, title, instructions, dueDate, now); err != nil {
			return err
		}
	} else {
		if _, err := getHomework(ctx, s.db, id); err != nil {
			return err
		}
		version, err := integer(data["version"], 1, maxSafeInteger)
		if err != nil {
			return err
		}
		err = s.tx(ctx, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, awardOnTimeSQL); err != nil {
				return err
			}
			// A new due date gives the student the postponement again. The delete reads the old due date before the update.
			if _, err := tx.Exec(ctx, "DELETE FROM postponements WHERE homework_id=$1 AND EXISTS(SELECT 1 FROM homework WHERE id=$1 AND version=$2 AND due_date IS DISTINCT FROM $3)", id, version, dueDate); err != nil {
				return err
			}
			tag, err := tx.Exec(ctx, "UPDATE homework SET title=$1,instructions=$2,due_date=$3,version=version+1,updated_at=$4 WHERE id=$5 AND version=$6", title, instructions, dueDate, now, id, version)
			if err != nil {
				return err
			}
			return changed(tag)
		})
		if err != nil {
			return err
		}
	}
	homework, err := getHomework(ctx, s.db, id)
	if err != nil {
		return err
	}
	return writeJSON(w, r, status, homework)
}

func (s *Server) homeworkTasks(w http.ResponseWriter, r *http.Request, user *SessionUser, id, action string) error {
	ctx, query := r.Context(), r.URL.Query()
	if id != "" {
		if _, err := uuid(id); err != nil {
			return err
		}
	}
	allowed := []string{http.MethodGet, http.MethodPost}
	if action != "" {
		allowed = []string{http.MethodPost}
	}
	if err := methods(r, allowed...); err != nil {
		return err
	}
	if r.Method == http.MethodGet && query.Get("view") == "board" {
		// The active homework with the earliest due date is in progress. The other active homework is in the future.
		active, err := s.tasks(ctx, s.db, "SELECT "+taskColumns+" FROM "+taskTables+" WHERE t.state NOT IN ('completed','cancelled') ORDER BY h.due_date IS NULL,h.due_date,h.created_at,h.id,t.created_at,t.id", user.Role)
		if err != nil {
			return err
		}
		accepted, err := s.tasks(ctx, s.db, "SELECT "+taskColumns+" FROM "+taskTables+" WHERE t.state='completed' ORDER BY t.updated_at DESC,t.id", user.Role)
		if err != nil {
			return err
		}
		progress, future := []HomeworkTask{}, []HomeworkTask{}
		for _, task := range active {
			if task.HomeworkID == active[0].HomeworkID {
				progress = append(progress, task)
			} else {
				future = append(future, task)
			}
		}
		return writeJSON(w, r, http.StatusOK, struct {
			Progress []HomeworkTask `json:"progress"`
			Future   []HomeworkTask `json:"future"`
			Accepted []HomeworkTask `json:"accepted"`
		}{progress, future, accepted})
	}
	if r.Method == http.MethodGet {
		limit, offset, err := pagination(query)
		if err != nil {
			return err
		}
		args := []any{user.Role}
		where := "1=1"
		filter := func(column string, value any) {
			args = append(args, value)
			where += " AND " + column + "=$" + strconv.Itoa(len(args))
		}
		if query.Has("problemId") {
			value, err := uuid(query.Get("problemId"))
			if err != nil {
				return err
			}
			filter("t.problem_id", value)
		}
		if query.Has("kind") {
			value, err := choice(query.Get("kind"), "problem", "chapter", "quiz")
			if err != nil {
				return err
			}
			filter("t.kind", value)
		}
		if query.Has("setSlug") {
			value, err := text(query.Get("setSlug"), 100, true)
			if err != nil {
				return err
			}
			filter("t.set_slug", value)
		}
		if query.Has("chapter") {
			value, err := queryInteger(query, "chapter", 0, 1, 99)
			if err != nil {
				return err
			}
			filter("t.chapter", value)
		}
		args = append(args, limit+1, offset)
		tasks, err := s.tasks(ctx, s.db, fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY (t.state='submitted') DESC,t.created_at DESC,t.id LIMIT $%d OFFSET $%d", taskColumns, taskTables, where, len(args)-1, len(args)), args...)
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, page(tasks, limit, offset))
	}
	role := "parent"
	if action == "start" || action == "complete" {
		role = "student"
	}
	if err := requireRole(user, role); err != nil {
		return err
	}
	now := isoNow()
	if id == "" {
		return s.assignProblem(w, r, user, now)
	}
	task, err := s.getTask(ctx, s.db, id, user.Role)
	if err != nil {
		return err
	}
	switch action {
	case "review":
		err = s.review(r, task, now)
	case "complete":
		var version int64
		if version, err = readVersion(r); err == nil {
			err = s.changeTask(ctx, "", "UPDATE homework_tasks SET state='completed',completed_at=$1,version=version+1,updated_at=$1 WHERE id=$2 AND version=$3 AND kind='chapter' AND state IN ('assigned','in_progress')", now, id, version)
		}
	case "start":
		var version int64
		if version, err = readVersion(r); err == nil {
			err = s.changeTask(ctx, "", "UPDATE homework_tasks SET state='in_progress',version=version+1,updated_at=$1 WHERE id=$2 AND version=$3 AND state='assigned'", now, id, version)
		}
	default:
		var version int64
		if version, err = readVersion(r); err == nil {
			err = s.changeTask(ctx, awardOnTimeSQL, "UPDATE homework_tasks SET state='cancelled',version=version+1,updated_at=$1 WHERE id=$2 AND version=$3 AND state NOT IN ('completed','cancelled')", now, id, version)
		}
	}
	if err != nil {
		return err
	}
	if task, err = s.getTask(ctx, s.db, id, user.Role); err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, task)
}

// changeTask runs an update of a homework task that has the version in its condition. before is a statement that runs first, or "".
func (s *Server) changeTask(ctx context.Context, before, sql string, args ...any) error {
	return s.tx(ctx, func(tx pgx.Tx) error {
		if before != "" {
			if _, err := tx.Exec(ctx, before); err != nil {
				return err
			}
		}
		tag, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		return changed(tag)
	})
}

// assignProblem makes a homework task for a problem. The unique index one_active_task allows one active task for each problem.
func (s *Server) assignProblem(w http.ResponseWriter, r *http.Request, user *SessionUser, now string) error {
	ctx := r.Context()
	data, err := readObject(r, "homeworkId", "problemId")
	if err != nil {
		return err
	}
	homeworkID, err := uuid(data["homeworkId"])
	if err != nil {
		return err
	}
	problemID, err := uuid(data["problemId"])
	if err != nil {
		return err
	}
	if _, err := getHomework(ctx, s.db, homeworkID); err != nil {
		return err
	}
	id := newID()
	err = s.tx(ctx, func(tx pgx.Tx) error {
		// The write of the problem row makes this check conflict with an archive or an accepted attempt of the same problem.
		if err := touchProblem(ctx, tx, problemID); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, "INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) SELECT $1,$2,p.id,'assigned',$3,$3 FROM problems p WHERE p.id=$4 AND "+assignableSQL+" ON CONFLICT DO NOTHING", id, homeworkID, now, problemID)
		if err != nil || tag.RowsAffected() == 1 {
			return err
		}
		problem, err := getProblem(ctx, tx, problemID)
		if err != nil {
			return err
		}
		if problem.ArchivedAt != nil {
			return &HTTPError{Status: http.StatusConflict, Code: "archived", Message: "Верните задачу из архива, прежде чем назначать домашнюю работу."}
		}
		var completed bool
		if err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM homework_tasks WHERE problem_id=$1 AND state='completed')", problemID).Scan(&completed); err != nil {
			return err
		}
		if problem.Solved || completed {
			return &HTTPError{Status: http.StatusConflict, Code: "completed_homework", Message: "Эта задача уже выполнена. Для повторения откройте задачу."}
		}
		return &HTTPError{Status: http.StatusConflict, Code: "active_homework", Message: "У этой задачи уже есть активное задание. Откройте его, чтобы продолжить работу."}
	})
	if err != nil {
		return err
	}
	task, err := s.getTask(ctx, s.db, id, user.Role)
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusCreated, task)
}

// review stores the decision of the parent on the current submission of a problem task: completed, or returned for changes.
func (s *Server) review(r *http.Request, task HomeworkTask, now string) error {
	ctx := r.Context()
	if task.ProblemID == nil {
		return invalid("На проверку отправляют решение задачи.")
	}
	data, err := readObject(r, "version", "submissionId", "decision", "body")
	if err != nil {
		return err
	}
	decision, err := choice(data["decision"], "changes_requested", "completed")
	if err != nil {
		return err
	}
	body, err := text(data["body"], 8000, decision == "changes_requested")
	if err != nil {
		return err
	}
	body = strings.TrimSpace(body)
	version, err := integer(data["version"], 1, maxSafeInteger)
	if err != nil {
		return err
	}
	submissionID, err := uuid(data["submissionId"])
	if err != nil {
		return err
	}
	// The parent completes a submitted task. The parent returns a submitted or a completed task of a problem that is not in the archive.
	condition := "state='submitted'"
	if decision == "changes_requested" {
		condition = "state IN ('submitted','completed') AND EXISTS(SELECT 1 FROM problems WHERE id=homework_tasks.problem_id AND archived_at IS NULL)"
	}
	return s.tx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, "INSERT INTO feedback (id,problem_id,task_id,submission_id,author,kind,body,created_at) SELECT $1,problem_id,id,current_submission_id,'parent',$2,$3,$4 FROM homework_tasks WHERE id=$5 AND version=$6 AND current_submission_id=$7 AND "+condition,
			newID(), decision, body, now, task.ID, version, submissionID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			if decision == "changes_requested" {
				problem, err := getProblem(ctx, tx, *task.ProblemID)
				if err != nil {
					return err
				}
				if problem.ArchivedAt != nil {
					return &HTTPError{Status: http.StatusConflict, Code: "archived", Message: "Верните задачу из архива, прежде чем отправлять её на доработку."}
				}
			}
			return changed(tag)
		}
		// The update of the task row makes two reviews of one task conflict.
		_, err = tx.Exec(ctx, "UPDATE homework_tasks SET state=$1,version=version+1,updated_at=$2 WHERE id=$3", decision, now, task.ID)
		return err
	})
}

package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rustatian/redlamp.me/backend/internal/db"
)

// querier is the read and write interface of a pool and of a transaction.
type querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// tx runs fn in a transaction with the retry on SQLSTATE 40001. fn can run more than one time.
func (s *Server) tx(ctx context.Context, fn func(pgx.Tx) error) error {
	return db.Tx(ctx, s.db, fn)
}

// isoNow returns the current time in the form of the API, for example 2026-10-07T12:34:56.789Z.
func isoNow() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// touchProblem writes the row of a problem. Two transactions that check the data of one problem and then change it
// both write this row, so the one that commits second fails with SQLSTATE 40001 and runs again.
// A row lock gives no such error on PostgreSQL: https://www.postgresql.org/docs/16/transaction-iso.html#XACT-REPEATABLE-READ
func touchProblem(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, "UPDATE problems SET version = version WHERE id = $1", id)
	return err
}

const (
	solvedSQL      = "EXISTS(SELECT 1 FROM attempts a WHERE a.problem_id=p.id AND a.state='saved' AND a.acceptance='accepted')"
	problemColumns = "p.id,p.slug,p.number,p.title,p.difficulty,p.topics,p.summary,p.version,p.progress_version,p.understanding,p.next_review_date,p.archived_at,p.created_at,p.updated_at," + solvedSQL
	// activeTaskSQL is true for a problem with a homework task that is not completed and not cancelled.
	activeTaskSQL = "EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state NOT IN ('completed','cancelled'))"
	assignableSQL = "p.archived_at IS NULL AND NOT " + solvedSQL + " AND NOT EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state<>'cancelled')"
)

// Problem is a problem of the library in the form of the API.
type Problem struct {
	ID              string   `json:"id"`
	Slug            string   `json:"slug"`
	Number          *int64   `json:"number"`
	Title           string   `json:"title"`
	Difficulty      string   `json:"difficulty"`
	Topics          []string `json:"topics"`
	Summary         string   `json:"summary"`
	Version         int64    `json:"version"`
	ProgressVersion int64    `json:"progressVersion"`
	Understanding   string   `json:"understanding"`
	NextReviewDate  *string  `json:"nextReviewDate"`
	ArchivedAt      *string  `json:"archivedAt"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
	Solved          bool     `json:"solved"`
	URL             string   `json:"url"`
}

// scanProblem reads the columns of problemColumns and then the columns of extra.
func scanProblem(row pgx.Row, extra ...any) (Problem, error) {
	var p Problem
	var stored string
	if err := row.Scan(append([]any{&p.ID, &p.Slug, &p.Number, &p.Title, &p.Difficulty, &stored, &p.Summary, &p.Version, &p.ProgressVersion, &p.Understanding, &p.NextReviewDate, &p.ArchivedAt, &p.CreatedAt, &p.UpdatedAt, &p.Solved}, extra...)...); err != nil {
		return p, err
	}
	var names []string
	if err := json.Unmarshal([]byte(stored), &names); err != nil {
		return p, err
	}
	p.Topics = normalizeTopics(names, p.Slug)
	p.URL = problemURL(p.Slug)
	return p, nil
}

func getProblem(ctx context.Context, q querier, id string) (Problem, error) {
	if _, err := uuid(id); err != nil {
		return Problem{}, err
	}
	p, err := scanProblem(q.QueryRow(ctx, "SELECT "+problemColumns+" FROM problems p WHERE p.id=$1", id))
	if errors.Is(err, pgx.ErrNoRows) {
		return p, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой задачи нет."}
	}
	return p, err
}

var (
	leetcodePath     = regexp.MustCompile(`^/problems/([a-z0-9]+(?:-[a-z0-9]+)*)(?:/|$)`)
	codeforcesPath   = regexp.MustCompile(`^/(?:problemset/problem/(\d{1,6})|contest/(\d{1,6})/problem)/([A-Za-z]\d?)(?:/|$)`)
	codeforcesSlug   = regexp.MustCompile(`^cf-(\d+)-([a-z]\d?)$`)
	problemsPath     = regexp.MustCompile(`^/leetcode/api/problems(?:/([^/]+)(?:/(archive|restore))?)?$`)
	frontendIDDigits = regexp.MustCompile(`^\d+$`)
)

// problemSlug reads the slug from the URL of a LeetCode or Codeforces problem.
// A Codeforces problem is stored as cf-<contest>-<index>, with the index in lower case.
func problemSlug(value any) (string, error) {
	invalidURL := &HTTPError{Status: http.StatusBadRequest, Code: "invalid_url", Message: "Укажите ссылку на задачу LeetCode или Codeforces."}
	link, err := text(value, 2000, true)
	if err != nil {
		return "", invalidURL
	}
	u, err := url.Parse(strings.TrimSpace(link))
	if err != nil || u.Scheme != "https" || u.User != nil || (u.Port() != "" && u.Port() != "443") {
		return "", invalidURL
	}
	host, path := strings.ToLower(u.Hostname()), u.EscapedPath()
	if match := leetcodePath.FindStringSubmatch(path); host == "leetcode.com" && match != nil {
		return match[1], nil
	}
	if match := codeforcesPath.FindStringSubmatch(path); host == "codeforces.com" && match != nil {
		return "cf-" + match[1] + match[2] + "-" + strings.ToLower(match[3]), nil
	}
	return "", invalidURL
}

func problemURL(slug string) string {
	if match := codeforcesSlug.FindStringSubmatch(slug); match != nil {
		return "https://codeforces.com/problemset/problem/" + match[1] + "/" + strings.ToUpper(match[2])
	}
	return "https://leetcode.com/problems/" + slug + "/"
}

// problemInput is the checked data of a request that adds or changes a problem.
type problemInput struct {
	slug       string
	number     *int64
	title      string
	difficulty string
	topics     []string
	summary    string
}

func metadata(data map[string]any) (problemInput, error) {
	var in problemInput
	var err error
	if in.slug, err = problemSlug(data["url"]); err != nil {
		return in, err
	}
	if value, ok := data["number"]; !ok || value != nil {
		number, err := integer(value, 1, maxSafeInteger)
		if err != nil {
			return in, err
		}
		in.number = &number
	}
	if in.title, err = text(data["title"], 200, true); err != nil {
		return in, err
	}
	if in.difficulty, err = choice(data["difficulty"], "easy", "medium", "hard"); err != nil {
		return in, err
	}
	if in.topics, err = topics(data["topics"], in.slug); err != nil {
		return in, err
	}
	if in.summary, err = text(data["summary"], 8000, false); err != nil {
		return in, err
	}
	in.title, in.summary = strings.TrimSpace(in.title), strings.TrimSpace(in.summary)
	return in, nil
}

// searchOf is the search text of a problem: the title, the number, the topics and the summary.
func searchOf(title string, number *int64, topics []string, summary string) string {
	parts := []string{title, ""}
	if number != nil {
		parts[1] = strconv.FormatInt(*number, 10)
	}
	parts = append(append(parts, topics...), summary)
	return searchText(strings.Join(parts, " "))
}

// problemDetails is the data of a LeetCode problem for the form that adds a problem.
type problemDetails struct {
	Number     *int64   `json:"number"`
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"`
	Topics     []string `json:"topics"`
}

func (s *Server) details(ctx context.Context, titleSlug string) (problemDetails, error) {
	var details problemDetails
	raw, err := s.leetcode.Question(ctx, titleSlug, "questionFrontendId title difficulty topicTags { name }")
	if err != nil {
		return details, &HTTPError{Status: http.StatusBadGateway, Code: "lookup_failed", Message: "Данные с LeetCode недоступны. Заполните их сами."}
	}
	if raw == nil {
		return details, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "На LeetCode нет задачи по этой ссылке."}
	}
	var question struct {
		QuestionFrontendID string `json:"questionFrontendId"`
		Title              string `json:"title"`
		Difficulty         string `json:"difficulty"`
		TopicTags          []struct {
			Name any `json:"name"`
		} `json:"topicTags"`
	}
	if err := json.Unmarshal(raw, &question); err != nil {
		return details, err
	}
	if frontendIDDigits.MatchString(question.QuestionFrontendID) {
		if number, err := strconv.ParseInt(question.QuestionFrontendID, 10, 64); err == nil {
			details.Number = &number
		}
	}
	tags := make([]any, 0, 12)
	for _, tag := range question.TopicTags[:min(len(question.TopicTags), 12)] {
		tags = append(tags, tag.Name)
	}
	details.Title, details.Difficulty = question.Title, strings.ToLower(question.Difficulty)
	details.Topics, err = topics(tags, titleSlug)
	return details, err
}

func (s *Server) handleProblems(w http.ResponseWriter, r *http.Request, _ *SessionUser) (bool, error) {
	ctx, query := r.Context(), r.URL.Query()
	if r.URL.Path == "/leetcode/api/problem-details" {
		if err := methods(r, http.MethodGet); err != nil {
			return true, err
		}
		var link any
		if query.Has("url") {
			link = query.Get("url")
		}
		slug, err := problemSlug(link)
		if err != nil {
			return true, err
		}
		if strings.HasPrefix(slug, "cf-") {
			return true, &HTTPError{Status: http.StatusBadRequest, Code: "lookup_unavailable", Message: "Автозаполнение есть только для задач LeetCode."}
		}
		details, err := s.details(ctx, slug)
		if err != nil {
			return true, err
		}
		return true, writeJSON(w, r, http.StatusOK, details)
	}
	match := problemsPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	id, action := match[1], match[2]
	if id != "" {
		if _, err := uuid(id); err != nil {
			return true, err
		}
	}
	allowed := []string{http.MethodGet, http.MethodPost}
	if action != "" {
		allowed = []string{http.MethodPost}
	} else if id != "" {
		allowed = []string{http.MethodGet, http.MethodPatch}
	}
	if err := methods(r, allowed...); err != nil {
		return true, err
	}
	if r.Method == http.MethodGet {
		if id != "" {
			problem, err := getProblem(ctx, s.db, id)
			if err != nil {
				return true, err
			}
			return true, writeJSON(w, r, http.StatusOK, problem)
		}
		return true, s.listProblems(w, r, query)
	}
	keys := []string{"version"}
	if action == "" {
		keys = []string{"url", "number", "title", "difficulty", "topics", "summary"}
		if id != "" {
			keys = append(keys, "version")
		}
	}
	body, err := readJSON(r, 32_768)
	if err != nil {
		return true, err
	}
	data, err := object(body, keys...)
	if err != nil {
		return true, err
	}
	now := isoNow()
	// status and target are the results of the transaction. The function sets them again on each run.
	var status int
	var target string
	err = s.tx(ctx, func(tx pgx.Tx) error {
		status, target = http.StatusOK, id
		if id != "" {
			if _, err := getProblem(ctx, tx, id); err != nil {
				return err
			}
		}
		switch action {
		case "restore":
			version, err := integer(data["version"], 1, maxSafeInteger)
			if err != nil {
				return err
			}
			tag, err := tx.Exec(ctx, "UPDATE problems SET archived_at=NULL,version=version+1,updated_at=$1 WHERE id=$2 AND version=$3", now, id, version)
			if err != nil {
				return err
			}
			return changed(tag)
		case "archive":
			version, err := integer(data["version"], 1, maxSafeInteger)
			if err != nil {
				return err
			}
			tag, err := tx.Exec(ctx, "UPDATE problems p SET archived_at=$1,version=version+1,updated_at=$1 WHERE id=$2 AND version=$3 AND NOT "+activeTaskSQL, now, id, version)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				var active bool
				if err := tx.QueryRow(ctx, "SELECT "+activeTaskSQL+" FROM problems p WHERE p.id=$1", id).Scan(&active); err != nil {
					return err
				}
				if active {
					return &HTTPError{Status: http.StatusConflict, Code: "active_homework", Message: "Завершите или отмените активное задание, прежде чем отправлять задачу в архив."}
				}
			}
			return changed(tag)
		}
		value, err := metadata(data)
		if err != nil {
			return err
		}
		names, err := json.Marshal(value.topics)
		if err != nil {
			return err
		}
		search := searchOf(value.title, value.number, value.topics, value.summary)
		duplicate := &HTTPError{Status: http.StatusConflict, Code: "duplicate_problem", Message: "Эта задача уже есть в тетради."}
		if id == "" {
			tag, err := tx.Exec(ctx, "INSERT INTO problems (id,slug,number,title,difficulty,topics,summary,search_text,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9) ON CONFLICT (slug) DO NOTHING",
				newID(), value.slug, value.number, value.title, value.difficulty, string(names), value.summary, search, now)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 1 {
				status = http.StatusCreated
			}
			return tx.QueryRow(ctx, "SELECT id FROM problems WHERE slug=$1", value.slug).Scan(&target)
		}
		var other bool
		if err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM problems WHERE slug=$1 AND id<>$2)", value.slug, id).Scan(&other); err != nil {
			return err
		}
		if other {
			return duplicate
		}
		// The judge and the solved mark follow the slug, so an accepted problem or a problem in homework keeps its slug.
		var locked bool
		if err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM problems p WHERE id=$1 AND slug<>$2 AND ("+solvedSQL+" OR EXISTS(SELECT 1 FROM homework_tasks WHERE problem_id=p.id AND state<>'cancelled')))", id, value.slug).Scan(&locked); err != nil {
			return err
		}
		if locked {
			return &HTTPError{Status: http.StatusConflict, Code: "slug_locked", Message: "Ссылку принятой задачи или задачи из домашней работы изменить нельзя."}
		}
		version, err := integer(data["version"], 1, maxSafeInteger)
		if err != nil {
			return err
		}
		// The award for a hard problem stays when the difficulty of that problem changes.
		if _, err := tx.Exec(ctx, "INSERT INTO earned_badges SELECT 'first-hard' WHERE EXISTS(SELECT 1 FROM problems p WHERE p.difficulty='hard' AND "+solvedSQL+") ON CONFLICT DO NOTHING"); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, "UPDATE problems SET slug=$1,number=$2,title=$3,difficulty=$4,topics=$5,summary=$6,search_text=$7,version=version+1,updated_at=$8 WHERE id=$9 AND version=$10",
			value.slug, value.number, value.title, value.difficulty, string(names), value.summary, search, now, id, version)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Another request took the slug after the check of this transaction.
			return duplicate
		}
		if err != nil {
			return err
		}
		return changed(tag)
	})
	if err != nil {
		return true, err
	}
	problem, err := getProblem(ctx, s.db, target)
	if err != nil {
		return true, err
	}
	return true, writeJSON(w, r, status, problem)
}

func (s *Server) listProblems(w http.ResponseWriter, r *http.Request, query url.Values) error {
	limit, offset, err := pagination(query)
	if err != nil {
		return err
	}
	q, err := text(query.Get("q"), 100, false)
	if err != nil {
		return err
	}
	q = searchText(q)
	topic := ""
	if query.Has("topic") {
		names, err := topics([]any{query.Get("topic")}, "")
		if err != nil {
			return err
		}
		topic = names[0]
	}
	var args []any
	arg := func(value any) string {
		args = append(args, value)
		return "$" + strconv.Itoa(len(args))
	}
	where := []string{"p.archived_at IS NULL"}
	if query.Get("archived") == "1" {
		where[0] = "p.archived_at IS NOT NULL"
	}
	if query.Has("difficulty") {
		value, err := choice(query.Get("difficulty"), "easy", "medium", "hard")
		if err != nil {
			return err
		}
		where = append(where, "p.difficulty="+arg(value))
	}
	if query.Has("understanding") {
		value, err := choice(query.Get("understanding"), "needs_practice", "with_help", "independent")
		if err != nil {
			return err
		}
		where = append(where, "p.understanding="+arg(value))
	}
	if query.Has("solved") {
		value, err := choice(query.Get("solved"), "0", "1")
		if err != nil {
			return err
		}
		where = append(where, solvedSQL+"="+arg(value == "1"))
	}
	if query.Has("homework") {
		value, err := choice(query.Get("homework"), "0", "1")
		if err != nil {
			return err
		}
		where = append(where, activeTaskSQL+"="+arg(value == "1"))
	}
	if query.Get("assignable") == "1" {
		where = append(where, assignableSQL)
	}
	if query.Has("due") {
		date, err := calendarDate(query.Get("due"))
		if err != nil {
			return err
		}
		where = append(where, "p.next_review_date<="+arg(*date))
	}
	// A search or a topic filter reads all rows: the topic names and a part of the search text are checked in Go.
	filtered := q != "" || topic != ""
	columns := problemColumns + ",false"
	if q != "" {
		n := arg(q)
		columns = problemColumns + ",(strpos(p.search_text," + n + ")>0 OR EXISTS(SELECT 1 FROM attempts a WHERE a.problem_id=p.id AND strpos(a.search_text," + n + ")>0))"
	}
	sql := "SELECT " + columns + " FROM problems p WHERE " + strings.Join(where, " AND ") + " ORDER BY p.updated_at DESC,p.id"
	if !filtered {
		sql += " LIMIT " + arg(limit+1) + " OFFSET " + arg(offset)
	}
	rows, err := s.db.Query(r.Context(), sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := []Problem{}
	for rows.Next() {
		var searchMatch bool
		problem, err := scanProblem(rows, &searchMatch)
		if err != nil {
			return err
		}
		if topic != "" && !hasTopic(problem.Topics, topic) {
			continue
		}
		if q != "" && !searchMatch && !strings.Contains(searchOf(problem.Title, problem.Number, problem.Topics, problem.Summary), q) {
			continue
		}
		items = append(items, problem)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if filtered {
		items = items[min(offset, int64(len(items))):min(offset+limit+1, int64(len(items)))]
	}
	return writeJSON(w, r, http.StatusOK, page(items, limit, offset))
}

// hasTopic reports whether values has the topic value. The comparison uses the search form.
func hasTopic(values []string, value string) bool {
	for _, v := range values {
		if searchText(v) == searchText(value) {
			return true
		}
	}
	return false
}

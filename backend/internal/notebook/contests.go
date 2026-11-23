package notebook

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rustatian/redlamp.me/backend/internal/answers"
	"github.com/rustatian/redlamp.me/backend/internal/bank"
	"github.com/rustatian/redlamp.me/backend/internal/editorials"
	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// contestLength is the time of each contest.
const contestLength = 90 * time.Minute

// contestGrace is the time after the end in which a submission still counts, so a check that starts at the last second counts.
const contestGrace = time.Minute

// contestSummary is the result of a contest in the form of the API. RatingChange and RatingAfter are null while the contest is active.
// Best is true for the best finished result: the higher score, then the shorter finish.
type contestSummary struct {
	ID            string `json:"id"`
	StartedAt     string `json:"startedAt"`
	EndsAt        string `json:"endsAt"`
	Active        bool   `json:"active"`
	Solved        int    `json:"solved"`
	Total         int    `json:"total"`
	Score         int    `json:"score"`
	FinishSeconds int64  `json:"finishSeconds"`
	RatingChange  *int   `json:"ratingChange"`
	RatingAfter   *int   `json:"ratingAfter"`
	Best          bool   `json:"best"`
}

// contestProblem is a problem of a contest in the form of the API. Rating is null while the contest is active.
// Wrong is the number of wrong submissions before the first accepted one. AcceptedSeconds is the time from the start to the first accepted submission.
type contestProblem struct {
	Position        int    `json:"position"`
	Title           string `json:"title"`
	Difficulty      string `json:"difficulty"`
	Points          int    `json:"points"`
	Topic           string `json:"topic"`
	Rating          *int   `json:"rating"`
	Solved          bool   `json:"solved"`
	Wrong           int    `json:"wrong"`
	AcceptedSeconds *int64 `json:"acceptedSeconds"`
}

// contest is a contest with its problems in the form of the API. Now is the server time.
// results are the results of the problems, in the order of Problems.
type contest struct {
	contestSummary
	Now      string           `json:"now"`
	Problems []contestProblem `json:"problems"`
	started  time.Time
	results  []problemResult
}

// contestClock is the active contest for the clock of the web client. Now is the server time.
type contestClock struct {
	ID     string `json:"id"`
	EndsAt string `json:"endsAt"`
	Now    string `json:"now"`
}

// submissionVerdict is the verdict of a contest submission in the form of the API. It never has a hidden input or a hidden answer.
// Test is the 1-based number of the failed test, or null. Error is the message of an error verdict, else null.
type submissionVerdict struct {
	Verdict string  `json:"verdict"`
	Test    *int    `json:"test"`
	Error   *string `json:"error"`
}

// contestSubmission is a submission of a contest problem in the form of the API.
type contestSubmission struct {
	submissionVerdict
	CreatedAt string `json:"createdAt"`
}

// contestProblemView is a contest problem with its statement, the last code and the submissions in the form of the API.
// Statement is null and Examples is empty for a problem that is not in the bank. URL and Editorial are null while the contest is active. Now is the server time.
type contestProblemView struct {
	Position    int                   `json:"position"`
	Title       string                `json:"title"`
	Difficulty  string                `json:"difficulty"`
	Points      int                   `json:"points"`
	Topic       string                `json:"topic"`
	Statement   *string               `json:"statement"`
	IgnoreCase  bool                  `json:"ignoreCase"`
	Examples    []bank.Test           `json:"examples"`
	Code        string                `json:"code"`
	Solved      bool                  `json:"solved"`
	Submissions []contestSubmission   `json:"submissions"`
	URL         *string               `json:"url"`
	Editorial   *editorials.Editorial `json:"editorial"`
	EndsAt      string                `json:"endsAt"`
	Now         string                `json:"now"`
}

var contestsPath = regexp.MustCompile(`^/leetcode/api/contests(?:/([^/]+)(?:/problems/([1-4])(?:/(code|check|submissions))?)?)?$`)

func (s *Server) handleContests(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	match := contestsPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	if match[2] != "" {
		position, _ := strconv.Atoi(match[2])
		return true, s.contestProblemRoute(w, r, user, match[1], position, match[3])
	}
	return true, s.contests(w, r, user, match[1])
}

// contestProblemRoute serves the problem at position of the contest id: the view, the code autosave, the start of a check and the submission.
// Each action first reads the contest problem, so a submission is only written for a problem of the contest.
func (s *Server) contestProblemRoute(w http.ResponseWriter, r *http.Request, user *SessionUser, id string, position int, action string) error {
	ctx := r.Context()
	if _, err := uuid(id); err != nil {
		return err
	}
	method := http.MethodPost
	switch action {
	case "":
		method = http.MethodGet
	case "code":
		method = http.MethodPut
	}
	if err := methods(r, method); err != nil {
		return err
	}
	if action != "" {
		if err := requireRole(user, "student"); err != nil {
			return err
		}
	}
	now := s.now().UTC()
	view := contestProblemView{Position: position, Examples: []bank.Test{}, Submissions: []contestSubmission{}, Now: now.Format(timeLayout)}
	var slug string
	err := s.db.QueryRow(ctx, "SELECT c.ends_at,p.slug,p.title,p.difficulty,p.points,p.topic,p.code FROM contest_problems p JOIN contests c ON c.id=p.contest_id WHERE p.contest_id=$1 AND p.position=$2", id, position).
		Scan(&view.EndsAt, &slug, &view.Title, &view.Difficulty, &view.Points, &view.Topic, &view.Code)
	if errors.Is(err, pgx.ErrNoRows) {
		return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой задачи нет."}
	}
	if err != nil {
		return err
	}
	rows, err := s.db.Query(ctx, "SELECT verdict,failed_test,error,created_at FROM contest_submissions WHERE contest_id=$1 AND position=$2 ORDER BY created_at,id", id, position)
	if err != nil {
		return err
	}
	var item contestSubmission
	if _, err := pgx.ForEachRow(rows, []any{&item.Verdict, &item.Test, &item.Error, &item.CreatedAt}, func() error {
		view.Solved = view.Solved || item.Verdict == "accepted"
		view.Submissions = append(view.Submissions, item)
		return nil
	}); err != nil {
		return err
	}
	problem, inBank := s.bankBySlug[slug]
	active := view.Now < view.EndsAt
	over := &HTTPError{Status: http.StatusConflict, Code: "contest_over", Message: "Контест закончился."}
	removed := &HTTPError{Status: http.StatusNotFound, Code: "problem_removed", Message: "Задача убрана из банка."}
	solved := &HTTPError{Status: http.StatusConflict, Code: "solved", Message: "Эта задача уже решена."}
	switch action {
	case "code":
		if !active {
			return over
		}
		// 100,000 UTF-16 units of code fit in 1 MiB of JSON, also with escapes.
		body, err := readJSON(r, 1_048_576)
		if err != nil {
			return err
		}
		data, err := object(body, "code")
		if err != nil {
			return err
		}
		code, err := text(data["code"], 100_000, false)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(ctx, "UPDATE contest_problems SET code=$1 WHERE contest_id=$2 AND position=$3", code, id, position); err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	case "check":
		if !active {
			return over
		}
		if !inBank {
			return removed
		}
		if view.Solved {
			return solved
		}
		inputs := make([]string, len(problem.Tests))
		for i, test := range problem.Tests {
			inputs[i] = test.Input
		}
		return writeJSON(w, r, http.StatusOK, map[string][]string{"inputs": inputs})
	case "submissions":
		if now.Add(-contestGrace).Format(timeLayout) >= view.EndsAt {
			return over
		}
		if !inBank {
			return removed
		}
		// 16 outputs of 1,000,000 ASCII characters fit in 16 MiB, the code and the rest of the JSON in 1 MiB.
		body, err := readJSON(r, 17_825_792)
		if err != nil {
			return err
		}
		data, err := object(body, "code", "outputs")
		if err != nil {
			return err
		}
		code, err := text(data["code"], 100_000, false)
		if err != nil {
			return err
		}
		raw, ok := data["outputs"].([]any)
		if !ok || len(raw) != len(problem.Tests) {
			return invalid("Проверьте поля запроса.")
		}
		outputs := make([]judgeOutput, len(raw))
		for i, item := range raw {
			if outputs[i], err = outputOf(item); err != nil {
				return err
			}
		}
		verdict, test, message := contestVerdict(problem, outputs)
		result := submissionVerdict{Verdict: verdict}
		if test > 0 {
			result.Test = &test
		}
		if verdict == "error" {
			result.Error = &message
		}
		err = s.tx(ctx, func(tx pgx.Tx) error {
			// The write of the problem row makes two submissions at the same time conflict. The second one runs again and finds the accepted one.
			if _, err := tx.Exec(ctx, "UPDATE contest_problems SET code=$1 WHERE contest_id=$2 AND position=$3", code, id, position); err != nil {
				return err
			}
			tag, err := tx.Exec(ctx, "INSERT INTO contest_submissions (id,contest_id,position,verdict,failed_test,error,code,created_at) SELECT $1,$2,$3,$4,$5,$6,$7,$8 WHERE NOT EXISTS (SELECT 1 FROM contest_submissions WHERE contest_id=$2 AND position=$3 AND verdict='accepted')",
				newID(), id, position, verdict, result.Test, result.Error, code, view.Now)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return solved
			}
			return nil
		})
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusCreated, result)
	}
	view.IgnoreCase = problem.IgnoreCase
	view.Examples = append(view.Examples, problem.Examples...)
	if inBank {
		view.Statement = &problem.Statement
	}
	if !active {
		url := sets.URL(slug)
		view.URL = &url
		if editorial, ok := s.editorials[slug]; ok {
			view.Editorial = &editorial
		}
	}
	return writeJSON(w, r, http.StatusOK, view)
}

// contestVerdict compares the outputs with the answers of the hidden tests, test by test. The first failed test gives the verdict and its 1-based number.
// message is the cut error message of an error verdict.
func contestVerdict(p bank.Problem, outputs []judgeOutput) (verdict string, test int, message string) {
	spec := answers.Spec{Mode: "stdio", IgnoreCase: p.IgnoreCase}
	for i, t := range p.Tests {
		output := outputs[i]
		if output.Value != nil && answers.Same(spec, t.Output, *output.Value) {
			continue
		}
		switch {
		case output.Value != nil:
			return "wrong", i + 1, ""
		case output.Error != nil:
			return "error", i + 1, cut(*output.Error)
		}
		return "timeout", i + 1, ""
	}
	return "accepted", 0, ""
}

func (s *Server) contests(w http.ResponseWriter, r *http.Request, user *SessionUser, id string) error {
	ctx := r.Context()
	now := s.now().UTC()
	moment := now.Format(timeLayout)
	switch id {
	case "":
		if err := methods(r, http.MethodGet, http.MethodPost); err != nil {
			return err
		}
		if r.Method == http.MethodGet {
			return s.listContests(w, r, moment)
		}
		if err := requireRole(user, "student"); err != nil {
			return err
		}
		return s.startContest(w, r, now)
	case "active":
		if err := methods(r, http.MethodGet); err != nil {
			return err
		}
		clock := contestClock{Now: moment}
		err := s.db.QueryRow(ctx, "SELECT id,ends_at FROM contests WHERE ends_at>$1", moment).Scan(&clock.ID, &clock.EndsAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return writeJSON(w, r, http.StatusOK, nil)
		}
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, clock)
	case "bank":
		if err := methods(r, http.MethodGet); err != nil {
			return err
		}
		return s.bankIndex(w, r)
	}
	if _, err := uuid(id); err != nil {
		return err
	}
	if err := methods(r, http.MethodGet); err != nil {
		return err
	}
	c, err := s.contestByID(ctx, id, moment)
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, c)
}

// readContests returns the contests in the order of their start with their results, the current rating and the best rating.
// now is the current time in the form of timeLayout. A contest is active while now is before its end.
func (s *Server) readContests(ctx context.Context, now string) (contests []contest, rating, bestRating int, err error) {
	// One transaction gives the three reads one snapshot.
	err = s.tx(ctx, func(tx pgx.Tx) error {
		contests = nil
		index := map[string]int{}
		rows, err := tx.Query(ctx, "SELECT id,started_at,ends_at FROM contests ORDER BY started_at,id")
		if err != nil {
			return err
		}
		var c contest
		if _, err := pgx.ForEachRow(rows, []any{&c.ID, &c.StartedAt, &c.EndsAt}, func() error {
			started, err := time.Parse(timeLayout, c.StartedAt)
			c.Active, c.Now, c.started = now < c.EndsAt, now, started
			index[c.ID] = len(contests)
			contests = append(contests, c)
			return err
		}); err != nil {
			return err
		}
		var id string
		var p contestProblem
		var result problemResult
		if rows, err = tx.Query(ctx, "SELECT contest_id,position,title,rating,topic,difficulty,points FROM contest_problems ORDER BY contest_id,position"); err != nil {
			return err
		}
		if _, err := pgx.ForEachRow(rows, []any{&id, &p.Position, &p.Title, &result.Rating, &p.Topic, &p.Difficulty, &p.Points}, func() error {
			c := &contests[index[id]]
			result.Points = p.Points
			c.Problems, c.results = append(c.Problems, p), append(c.results, result)
			return nil
		}); err != nil {
			return err
		}
		var position int
		var verdict, created string
		if rows, err = tx.Query(ctx, "SELECT contest_id,position,verdict,created_at FROM contest_submissions ORDER BY created_at,id"); err != nil {
			return err
		}
		_, err = pgx.ForEachRow(rows, []any{&id, &position, &verdict, &created}, func() error {
			c := &contests[index[id]]
			result := &c.results[position-1]
			if result.Solved {
				return nil
			}
			if verdict != "accepted" {
				result.Wrong++
				return nil
			}
			at, err := time.Parse(timeLayout, created)
			result.Solved, result.Accepted = true, at.Sub(c.started)
			return err
		})
		return err
	})
	if err != nil {
		return nil, 0, 0, err
	}
	var finished [][]problemResult
	for i := range contests {
		c := &contests[i]
		score, finish := contestScore(c.results)
		c.Total, c.Score, c.FinishSeconds = len(c.results), score, int64(finish/time.Second)
		for j, result := range c.results {
			p := &c.Problems[j]
			p.Solved, p.Wrong = result.Solved, result.Wrong
			if result.Solved {
				c.Solved++
				seconds := int64(result.Accepted / time.Second)
				p.AcceptedSeconds = &seconds
			}
			if !c.Active {
				p.Rating = &c.results[j].Rating
			}
		}
		if !c.Active {
			finished = append(finished, c.results)
		}
	}
	history := ratingHistory(finished)
	rating, bestRating, best := startRating, startRating, -1
	for i := range contests {
		c := &contests[i]
		if c.Active {
			continue
		}
		change, after := history[0]-rating, history[0]
		history = history[1:]
		c.RatingChange, c.RatingAfter = &change, &after
		rating, bestRating = after, max(bestRating, after)
		if best < 0 || c.Score > contests[best].Score || c.Score == contests[best].Score && c.FinishSeconds < contests[best].FinishSeconds {
			best = i
		}
	}
	if best >= 0 {
		contests[best].Best = true
	}
	return contests, rating, bestRating, nil
}

// contestByID returns the contest id with its results.
func (s *Server) contestByID(ctx context.Context, id, now string) (contest, error) {
	contests, _, _, err := s.readContests(ctx, now)
	if err != nil {
		return contest{}, err
	}
	for _, c := range contests {
		if c.ID == id {
			return c, nil
		}
	}
	return contest{}, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такого контеста нет."}
}

// listContests writes the rating, the best rating, the best finished result, the active contest and the contests, newest first.
func (s *Server) listContests(w http.ResponseWriter, r *http.Request, now string) error {
	contests, rating, bestRating, err := s.readContests(r.Context(), now)
	if err != nil {
		return err
	}
	list := struct {
		Rating     int              `json:"rating"`
		BestRating int              `json:"bestRating"`
		Best       *contestSummary  `json:"best"`
		Active     *contestClock    `json:"active"`
		Items      []contestSummary `json:"items"`
	}{Rating: rating, BestRating: bestRating, Items: make([]contestSummary, len(contests))}
	for i, c := range contests {
		item := &list.Items[len(contests)-1-i]
		*item = c.contestSummary
		if c.Best {
			list.Best = item
		}
		if c.Active {
			list.Active = &contestClock{ID: c.ID, EndsAt: c.EndsAt, Now: now}
		}
	}
	return writeJSON(w, r, http.StatusOK, list)
}

// bankIndex writes all topics of the bank and the bank problems that are not in a contest and not in the library.
func (s *Server) bankIndex(w http.ResponseWriter, r *http.Request) error {
	rows, err := s.db.Query(r.Context(), "SELECT slug FROM contest_problems UNION SELECT slug FROM problems")
	if err != nil {
		return err
	}
	used := map[string]bool{}
	var slug string
	if _, err := pgx.ForEachRow(rows, []any{&slug}, func() error {
		used[slug] = true
		return nil
	}); err != nil {
		return err
	}
	type entry struct {
		Slug       string   `json:"slug"`
		Rating     int      `json:"rating"`
		Topics     []string `json:"topics"`
		Difficulty string   `json:"difficulty"`
		Points     int      `json:"points"`
	}
	index := struct {
		Topics   []string `json:"topics"`
		Problems []entry  `json:"problems"`
	}{Topics: []string{}, Problems: []entry{}}
	for _, p := range s.bank {
		index.Topics = append(index.Topics, p.Topics...)
		if !used[p.Slug] {
			difficulty := bank.Difficulty(p.Rating)
			index.Problems = append(index.Problems, entry{Slug: p.Slug, Rating: p.Rating, Topics: p.Topics, Difficulty: difficulty, Points: bank.Points(difficulty)})
		}
	}
	slices.Sort(index.Topics)
	index.Topics = slices.Compact(index.Topics)
	return writeJSON(w, r, http.StatusOK, index)
}

// startContest makes a contest of 1 to 4 bank problems that are not in an earlier contest and not in the library:
// at most 1 easy, 2 medium and 1 hard. The contest starts now and lasts contestLength. Only one contest is active at a time.
func (s *Server) startContest(w http.ResponseWriter, r *http.Request, now time.Time) error {
	ctx := r.Context()
	data, err := readObject(r, "problems")
	if err != nil {
		return err
	}
	items, ok := data["problems"].([]any)
	if !ok || len(items) == 0 || len(items) > 4 {
		return invalid("Выберите от 1 до 4 задач.")
	}
	type pick struct {
		problem           bank.Problem
		topic, difficulty string
	}
	picks := make([]pick, len(items))
	free := map[string]int{"easy": 1, "medium": 2, "hard": 1}
	for i, item := range items {
		fields, err := object(item, "slug", "topic")
		if err != nil {
			return err
		}
		slug, _ := fields["slug"].(string)
		topic, _ := fields["topic"].(string)
		p, ok := s.bankBySlug[slug]
		if !ok || !slices.Contains(p.Topics, topic) || slices.ContainsFunc(picks[:i], func(c pick) bool { return c.problem.Slug == slug }) {
			return invalid("Выберите каждую задачу из банка один раз и с её темой.")
		}
		difficulty := bank.Difficulty(p.Rating)
		free[difficulty]--
		if free[difficulty] < 0 {
			return &HTTPError{Status: http.StatusBadRequest, Code: "contest_slots", Message: "В контесте не больше 1 лёгкой, 2 средних и 1 сложной задачи."}
		}
		picks[i] = pick{problem: p, topic: topic, difficulty: difficulty}
	}
	id, started := newID(), now.Format(timeLayout)
	used := &HTTPError{Status: http.StatusConflict, Code: "problem_used", Message: "Эта задача уже была в контесте или есть в тетради."}
	err = s.tx(ctx, func(tx pgx.Tx) error {
		// The write of the guard row makes two starts at the same time conflict. The second one runs again and finds the active contest.
		if _, err := tx.Exec(ctx, "INSERT INTO contest_starts VALUES (1,$1) ON CONFLICT (id) DO UPDATE SET started_at=excluded.started_at", started); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, "INSERT INTO contests SELECT $1,$2,$3 WHERE NOT EXISTS (SELECT 1 FROM contests WHERE ends_at>$2)", id, started, now.Add(contestLength).Format(timeLayout))
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return &HTTPError{Status: http.StatusConflict, Code: "contest_active", Message: "Контест уже идёт."}
		}
		for i, p := range picks {
			tag, err := tx.Exec(ctx, "INSERT INTO contest_problems (contest_id,position,slug,title,rating,topic,difficulty,points) SELECT $1,$2,$3,$4,$5,$6,$7,$8 WHERE NOT EXISTS (SELECT 1 FROM problems WHERE slug=$3)",
				id, i+1, p.problem.Slug, p.problem.Title, p.problem.Rating, p.topic, p.difficulty, bank.Points(p.difficulty))
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				// The unique index contest_problems_slug: the problem was in an earlier contest.
				return used
			}
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return used
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	c, err := s.contestByID(ctx, id, started)
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusCreated, c)
}

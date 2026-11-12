package notebook

import (
	"cmp"
	"context"
	"encoding/json"
	"hash/fnv"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const reviewColumns = "id,problem_id,result,note,reviewed_on,next_review_date,created_at"

// review is one recall of a problem in the form of the API.
type review struct {
	ID             string  `json:"id"`
	ProblemID      string  `json:"problemId"`
	Result         string  `json:"result"`
	Note           string  `json:"note"`
	ReviewedOn     string  `json:"reviewedOn"`
	NextReviewDate *string `json:"nextReviewDate"`
	CreatedAt      string  `json:"createdAt"`
}

func scanReview(row pgx.Row) (review, error) {
	var v review
	return v, row.Scan(&v.ID, &v.ProblemID, &v.Result, &v.Note, &v.ReviewedOn, &v.NextReviewDate, &v.CreatedAt)
}

// addDays returns the calendar date that is days after date.
func addDays(date string, days int) (string, error) {
	day, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "", err
	}
	return day.AddDate(0, 0, days).Format(time.DateOnly), nil
}

// spreadProblem is a problem for the 30-day spread. NextReviewDate is nil for a problem without a recall date.
type spreadProblem struct {
	ID              string
	NextReviewDate  *string
	ProgressVersion int64
}

// spreadDays gives a recall date to each problem that has no date in the window [today + 1, today + 30].
// A date in the window stays. Each other problem gets a random day of the days with the fewest problems.
func spreadDays(today string, problems []spreadProblem, random func(count int) int) (map[string]string, error) {
	window := make([]string, 30)
	uses := map[string]int{}
	for i := range window {
		day, err := addDays(today, i+1)
		if err != nil {
			return nil, err
		}
		window[i], uses[day] = day, 0
	}
	inWindow := func(p spreadProblem) bool {
		if p.NextReviewDate == nil {
			return false
		}
		_, ok := uses[*p.NextReviewDate]
		return ok
	}
	for _, p := range problems {
		if inWindow(p) {
			uses[*p.NextReviewDate]++
		}
	}
	result := map[string]string{}
	for _, p := range problems {
		if inWindow(p) {
			continue
		}
		least := uses[window[0]]
		for _, day := range window {
			least = min(least, uses[day])
		}
		var free []string
		for _, day := range window {
			if uses[day] == least {
				free = append(free, day)
			}
		}
		day := free[random(len(free))]
		uses[day] = least + 1
		result[p.ID] = day
	}
	return result, nil
}

// recallPools are the sources of the recall picks of the Today page: a number of problems of a learning set, or of a list of slugs.
var recallPools = []struct {
	count int
	set   string
	slugs []string
}{
	{count: 3, set: "two-pointers"},
	{count: 2, slugs: []string{"sum-of-all-subset-xor-totals", "letter-case-permutation", "letter-combinations-of-a-phone-number"}},
}

// recallPick is a problem that the Today page offers for a recall.
type recallPick struct {
	URL        string   `json:"url"`
	Number     *int64   `json:"number,omitempty"`
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"`
	Topics     []string `json:"topics"`
	Note       string   `json:"note"`
	Slug       string   `json:"slug"`
	Label      string   `json:"label"`
	ProblemID  *string  `json:"problemId"`
}

// taskLabel is the short name of a problem: its source and its number.
func taskLabel(slug string, number *int64) string {
	if rest, ok := strings.CutPrefix(slug, "cf-"); ok {
		return "Codeforces " + strings.ToUpper(strings.Replace(rest, "-", "", 1))
	}
	if number != nil && *number != 0 {
		return "LeetCode " + strconv.FormatInt(*number, 10)
	}
	return "LeetCode"
}

// setTasks returns the problems of all chapters of a learning set.
func (s *Server) setTasks(slug string) []SetTask {
	var tasks []SetTask
	for _, set := range s.sets {
		if set.Slug == slug {
			for _, chapter := range set.Chapters {
				tasks = append(tasks, chapter.Tasks...)
			}
		}
	}
	return tasks
}

// recallPicks chooses the recall picks of a day. A problem in the archive, a problem that is due, and a problem that was recalled today are not picked.
// The order of a pool follows a hash of the day and the slug, so the picks stay the same for the whole day.
func (s *Server) recallPicks(ctx context.Context, q querier, today string) ([]recallPick, error) {
	type candidate struct {
		task SetTask
		slug string
	}
	pools := make([][]candidate, len(recallPools))
	var slugs []string
	for i, pool := range recallPools {
		if pool.slugs != nil {
			for _, slug := range pool.slugs {
				pools[i] = append(pools[i], candidate{slug: slug})
			}
		} else {
			for _, task := range s.setTasks(pool.set) {
				slug, err := problemSlug(task.URL)
				if err != nil {
					return nil, err
				}
				pools[i] = append(pools[i], candidate{task: task, slug: slug})
			}
		}
		for _, c := range pools[i] {
			slugs = append(slugs, c.slug)
		}
	}
	type libraryRow struct {
		id         string
		number     *int64
		title      string
		difficulty string
		topics     string
		summary    string
		skip       bool
	}
	library := map[string]libraryRow{}
	rows, err := q.Query(ctx, "SELECT p.slug,p.id,p.number,p.title,p.difficulty,p.topics,p.summary,coalesce(p.archived_at IS NOT NULL OR p.next_review_date<=$1 OR EXISTS(SELECT 1 FROM reviews r WHERE r.problem_id=p.id AND r.reviewed_on=$1),false) FROM problems p WHERE p.slug=ANY($2)", today, slugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var slug string
		var row libraryRow
		if err := rows.Scan(&slug, &row.id, &row.number, &row.title, &row.difficulty, &row.topics, &row.summary, &row.skip); err != nil {
			return nil, err
		}
		library[slug] = row
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	hash := func(slug string) uint32 {
		// 32-bit FNV-1a: https://www.rfc-editor.org/rfc/rfc9923.html#section-2
		h := fnv.New32a()
		_, _ = h.Write([]byte(today + ":" + slug))
		return h.Sum32()
	}
	picks := []recallPick{}
	for i, pool := range recallPools {
		var poolPicks []recallPick
		for _, c := range pools[i] {
			row, inLibrary := library[c.slug]
			task := c.task
			if pool.slugs != nil {
				// A pool of slugs offers only the problems of the library.
				if !inLibrary {
					continue
				}
				task = SetTask{URL: problemURL(c.slug), Number: row.number, Title: row.title, Difficulty: row.difficulty, Note: row.summary}
				if err := json.Unmarshal([]byte(row.topics), &task.Topics); err != nil {
					return nil, err
				}
			}
			if inLibrary && row.skip {
				continue
			}
			pick := recallPick{URL: task.URL, Number: task.Number, Title: task.Title, Difficulty: task.Difficulty, Topics: normalizeTopics(task.Topics, c.slug), Note: task.Note, Slug: c.slug, Label: taskLabel(c.slug, task.Number)}
			if inLibrary {
				pick.ProblemID = &row.id
			}
			poolPicks = append(poolPicks, pick)
		}
		slices.SortStableFunc(poolPicks, func(a, b recallPick) int { return cmp.Compare(hash(a.Slug), hash(b.Slug)) })
		picks = append(picks, poolPicks[:min(len(poolPicks), pool.count)]...)
	}
	return picks, nil
}

// dashboardHomework returns the number of tasks that wait for a review and the first 10 active tasks for the role.
func (s *Server) dashboardHomework(ctx context.Context, q querier, role string) (int64, []HomeworkTask, error) {
	var waiting int64
	if err := q.QueryRow(ctx, "SELECT count(*) FROM homework_tasks t JOIN problems p ON p.id=t.problem_id WHERE t.state='submitted' AND p.archived_at IS NULL").Scan(&waiting); err != nil {
		return 0, nil, err
	}
	// The parent sees the submitted tasks first. The student sees the tasks with the earliest due date first.
	order := "(t.state='submitted') ASC,h.due_date IS NULL,h.due_date ASC"
	if role == "parent" {
		order = "(t.state='submitted') DESC,t.updated_at ASC"
	}
	tasks, err := s.tasks(ctx, q, "SELECT "+taskColumns+" FROM "+taskTables+" WHERE p.archived_at IS NULL AND t.state NOT IN ('completed','cancelled') ORDER BY "+order+",t.id LIMIT 10", role)
	return waiting, tasks, err
}

var reviewsPath = regexp.MustCompile(`^/leetcode/api/problems/([^/]+)/(reviews|practice)$`)

func (s *Server) handleReviews(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	switch r.URL.Path {
	case "/leetcode/api/dashboard":
		return true, s.dashboard(w, r, user)
	case "/leetcode/api/reviews/distribute":
		return true, s.distribute(w, r)
	case "/leetcode/api/game":
		return true, s.game(w, r, user)
	}
	match := reviewsPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	return true, s.reviews(w, r, user, match[1], match[2] == "practice")
}

// readToday reads a request that has only the local date of the browser.
func readToday(r *http.Request) (string, error) {
	data, err := readObject(r, "today")
	if err != nil {
		return "", err
	}
	return todayValue(data["today"])
}

func todayValue(value any) (string, error) {
	day, err := text(value, 10, true)
	if err != nil {
		return "", err
	}
	date, err := calendarDate(day)
	if err != nil {
		return "", err
	}
	return *date, nil
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request, user *SessionUser) error {
	ctx := r.Context()
	if err := methods(r, http.MethodGet); err != nil {
		return err
	}
	var value any
	if query := r.URL.Query(); query.Has("today") {
		value = query.Get("today")
	}
	today, err := todayValue(value)
	if err != nil {
		return err
	}
	type topicCount struct {
		Topic       string `json:"topic"`
		Solved      int    `json:"solved"`
		Independent int    `json:"independent"`
	}
	type draft struct {
		Attempt
		ProblemTitle string `json:"problemTitle"`
	}
	type dueReview struct {
		Problem
		FromHomework bool `json:"fromHomework"`
	}
	var result struct {
		Counts struct {
			Recorded      int64 `json:"recorded"`
			Solved        int64 `json:"solved"`
			Independent   int64 `json:"independent"`
			DueReviews    int64 `json:"dueReviews"`
			WaitingReview int64 `json:"waitingReview"`
		} `json:"counts"`
		Topics      []topicCount   `json:"topics"`
		Homework    []HomeworkTask `json:"homework"`
		Drafts      []draft        `json:"drafts"`
		DueReviews  []dueReview    `json:"dueReviews"`
		RecallPicks []recallPick   `json:"recallPicks"`
	}
	count := func(condition string) string {
		return "coalesce(sum(CASE WHEN " + condition + " THEN 1 ELSE 0 END),0)"
	}
	if err := s.db.QueryRow(ctx, "SELECT count(*),"+count(solvedSQL)+","+count(solvedSQL+" AND p.understanding='independent'")+","+count("p.next_review_date<=$1")+" FROM problems p WHERE p.archived_at IS NULL", today).
		Scan(&result.Counts.Recorded, &result.Counts.Solved, &result.Counts.Independent, &result.Counts.DueReviews); err != nil {
		return err
	}
	rows, err := s.db.Query(ctx, "SELECT p.slug,p.topics,p.understanding FROM problems p WHERE p.archived_at IS NULL AND "+solvedSQL)
	if err != nil {
		return err
	}
	byTopic := map[string]*topicCount{}
	var slug, stored, understanding string
	_, err = pgx.ForEachRow(rows, []any{&slug, &stored, &understanding}, func() error {
		var names []string
		if err := json.Unmarshal([]byte(stored), &names); err != nil {
			return err
		}
		for _, topic := range normalizeTopics(names, slug) {
			counts := byTopic[topic]
			if counts == nil {
				counts = &topicCount{Topic: topic}
				byTopic[topic] = counts
			}
			counts.Solved++
			if understanding == "independent" {
				counts.Independent++
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	names := make([]string, 0, len(byTopic))
	for topic := range byTopic {
		names = append(names, topic)
	}
	slices.Sort(names)
	result.Topics = []topicCount{}
	for _, topic := range names[:min(len(names), 50)] {
		result.Topics = append(result.Topics, *byTopic[topic])
	}
	// The work window opens a draft of an accepted problem after the acceptance, so that draft is not a next step.
	rows, err = s.db.Query(ctx, "SELECT "+attemptColumns+",(SELECT title FROM problems p WHERE p.id=attempts.problem_id) FROM attempts WHERE state='draft' AND EXISTS(SELECT 1 FROM problems p WHERE p.id=attempts.problem_id AND p.archived_at IS NULL AND NOT "+solvedSQL+") ORDER BY updated_at DESC,id LIMIT 5")
	if err != nil {
		return err
	}
	drafts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (draft, error) {
		var d draft
		return d, row.Scan(&d.ID, &d.ProblemID, &d.State, &d.Acceptance, &d.Understanding, &d.Version, &d.CreatedAt, &d.UpdatedAt, &d.SavedAt, &d.ProblemTitle)
	})
	if err != nil {
		return err
	}
	result.Drafts = append([]draft{}, drafts...)
	rows, err = s.db.Query(ctx, "SELECT "+problemColumns+",EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state='completed') FROM problems p WHERE p.archived_at IS NULL AND p.next_review_date<=$1 ORDER BY p.next_review_date,p.id LIMIT 10", today)
	if err != nil {
		return err
	}
	due, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (dueReview, error) {
		var d dueReview
		var err error
		d.Problem, err = scanProblem(row, &d.FromHomework)
		return d, err
	})
	if err != nil {
		return err
	}
	result.DueReviews = append([]dueReview{}, due...)
	if result.Counts.WaitingReview, result.Homework, err = s.dashboardHomework(ctx, s.db, user.Role); err != nil {
		return err
	}
	if result.RecallPicks, err = s.recallPicks(ctx, s.db, today); err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, result)
}

// distribute spreads the recall dates of the problems with a completed homework task over the next 30 days.
func (s *Server) distribute(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	if err := methods(r, http.MethodPost); err != nil {
		return err
	}
	today, err := readToday(r)
	if err != nil {
		return err
	}
	const completed = "archived_at IS NULL AND EXISTS(SELECT 1 FROM homework_tasks t WHERE t.problem_id=p.id AND t.state='completed')"
	rows, err := s.db.Query(ctx, "SELECT id,next_review_date,progress_version FROM problems p WHERE "+completed)
	if err != nil {
		return err
	}
	problems, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (spreadProblem, error) {
		var p spreadProblem
		return p, row.Scan(&p.ID, &p.NextReviewDate, &p.ProgressVersion)
	})
	if err != nil {
		return err
	}
	days, err := spreadDays(today, problems, s.random)
	if err != nil {
		return err
	}
	now := isoNow()
	var dates []string
	err = s.tx(ctx, func(tx pgx.Tx) error {
		dates = nil
		for _, p := range problems {
			day, ok := days[p.ID]
			if !ok {
				continue
			}
			// The version check keeps a recall that was recorded after the read. The answer reports only the dates that this request changed.
			tag, err := tx.Exec(ctx, "UPDATE problems p SET next_review_date=$1,progress_version=progress_version+1,updated_at=$2 WHERE id=$3 AND progress_version=$4 AND "+completed, day, now, p.ID, p.ProgressVersion)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 1 {
				dates = append(dates, day)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	slices.Sort(dates)
	var distribution struct {
		Count int     `json:"count"`
		From  *string `json:"from"`
		To    *string `json:"to"`
	}
	if distribution.Count = len(dates); len(dates) > 0 {
		distribution.From, distribution.To = &dates[0], &dates[len(dates)-1]
	}
	return writeJSON(w, r, http.StatusOK, distribution)
}

// badgeIDs is the order of the awards in the answer of the game.
var badgeIDs = []string{"first-accept", "three-day-streak", "first-hard", "ten-accepted", "independent-five", "set-cleared", "week-no-miss", "month-no-miss", "twenty-five-accepted", "fifty-accepted", "fifty-cards", "ten-recalls", "first-codeforces", "first-try", "two-sets", "on-time-five", "level-five"}

// gameTotals are the counts of the notebook that give the XP and the awards.
type gameTotals struct {
	accepted, hard, acceptedXP, independent, completed, recalls, cards, passedTests, failedRuns, late, onTime, firstTry, codeforces int64
	// streak is true when attempts were saved on three days in a row.
	streak bool
	// missed is the number of days without a visit from the first visit to yesterday. longest is the longest run of visit days.
	missed, longest int64
	// cleared is the number of learning sets with all problems accepted.
	cleared int
}

type badge struct {
	ID     string `json:"id"`
	Earned bool   `json:"earned"`
}

// game is the answer of GET /leetcode/api/game.
type game struct {
	XP          int64   `json:"xp"`
	Level       int64   `json:"level"`
	LevelXP     int64   `json:"levelXp"`
	NextLevelXP int64   `json:"nextLevelXp"`
	Late        int64   `json:"late"`
	Missed      int64   `json:"missed"`
	Badges      []badge `json:"badges"`
}

// gameOf computes the XP, the level and the awards from the counts. Level n starts at 50*n*(n-1) XP.
func gameOf(t gameTotals) game {
	xp := max(0, t.acceptedXP+10*t.independent+15*t.completed+5*t.recalls+2*t.cards+10*t.passedTests-2*t.failedRuns-15*t.late-5*t.missed)
	level := int64(1)
	for xp >= 50*level*(level+1) {
		level++
	}
	earned := map[string]bool{
		"first-accept": t.accepted >= 1, "three-day-streak": t.streak, "first-hard": t.hard >= 1, "ten-accepted": t.accepted >= 10, "independent-five": t.independent >= 5, "set-cleared": t.cleared >= 1,
		"week-no-miss": t.longest >= 7, "month-no-miss": t.longest >= 30, "twenty-five-accepted": t.accepted >= 25, "fifty-accepted": t.accepted >= 50, "fifty-cards": t.cards >= 50, "ten-recalls": t.recalls >= 10,
		"first-codeforces": t.codeforces >= 1, "first-try": t.firstTry >= 1, "two-sets": t.cleared >= 2, "on-time-five": t.onTime >= 5, "level-five": level >= 5,
	}
	result := game{XP: xp, Level: level, LevelXP: 50 * level * (level - 1), NextLevelXP: 50 * level * (level + 1), Late: t.late, Missed: t.missed}
	for _, id := range badgeIDs {
		result.Badges = append(result.Badges, badge{ID: id, Earned: earned[id]})
	}
	return result
}

// hasStreak reports whether days has three calendar days in a row. A text that is not a date does not count.
func hasStreak(days []string) bool {
	saved := map[string]bool{}
	for _, day := range days {
		saved[day] = true
	}
	for _, day := range days {
		next, err := addDays(day, 1)
		if err != nil {
			continue
		}
		after, _ := addDays(day, 2)
		if saved[next] && saved[after] {
			return true
		}
	}
	return false
}

func (s *Server) game(w http.ResponseWriter, r *http.Request, user *SessionUser) error {
	ctx := r.Context()
	if err := methods(r, http.MethodGet); err != nil {
		return err
	}
	// A task is late after its due date in UTC. A task added after the due date is not late.
	today := s.now().UTC().Format(time.DateOnly)
	var setSlugs [][]string
	var all []string
	for _, set := range s.sets {
		var slugs []string
		for _, chapter := range set.Chapters {
			for _, task := range chapter.Tasks {
				slug, err := problemSlug(task.URL)
				if err != nil {
					return err
				}
				slugs = append(slugs, slug)
			}
		}
		setSlugs, all = append(setSlugs, slugs), append(all, slugs...)
	}
	count := func(condition string) string {
		return "coalesce(sum(CASE WHEN " + condition + " THEN 1 ELSE 0 END),0)"
	}
	var result game
	err := s.tx(ctx, func(tx pgx.Tx) error {
		// The web client reads the game on each page change, so this read records the visit day of the student.
		if user.Role == "student" {
			if _, err := tx.Exec(ctx, "INSERT INTO visits (day) VALUES ($1) ON CONFLICT DO NOTHING", today); err != nil {
				return err
			}
		}
		var t gameTotals
		// cards are the correct answers of the single-card quiz and the best passed attempt of each deck.
		err := tx.QueryRow(ctx, "SELECT count(*),"+count("p.difficulty='hard'")+",coalesce(sum(CASE p.difficulty WHEN 'easy' THEN 10 WHEN 'medium' THEN 20 ELSE 40 END),0),"+
			count("EXISTS(SELECT 1 FROM attempts a WHERE a.problem_id=p.id AND a.state='saved' AND a.understanding='independent') OR EXISTS(SELECT 1 FROM reviews r WHERE r.problem_id=p.id AND r.result='independent')")+","+
			"(SELECT count(*) FROM homework_tasks WHERE state='completed'),"+
			"(SELECT count(*) FROM (SELECT DISTINCT problem_id,reviewed_on FROM reviews) recalls),"+
			"(SELECT count(*) FROM card_answers WHERE correct=1)+(SELECT coalesce(sum(best),0)::bigint FROM (SELECT max(correct) AS best FROM quiz_attempts WHERE passed=1 GROUP BY deck) decks),"+
			"(SELECT count(DISTINCT problem_id) FROM test_runs WHERE passed=total AND total>0),"+
			"(SELECT coalesce(sum(least(fails,3)),0)::bigint FROM (SELECT count(*) AS fails FROM test_runs WHERE passed<total GROUP BY problem_id,run_on) days),"+
			"(SELECT count(*) FROM homework_tasks t JOIN homework h ON h.id=t.homework_id WHERE t.state<>'cancelled' AND h.due_date<$1 AND substr(t.created_at,1,10)<=h.due_date AND (t.completed_at IS NULL OR substr(t.completed_at,1,10)>h.due_date) AND NOT EXISTS(SELECT 1 FROM submissions s WHERE s.task_id=t.id AND substr(s.created_at,1,10)<=h.due_date)),"+
			"("+onTimeSQL+"),"+
			"(SELECT count(*) FROM test_runs r WHERE r.passed=r.total AND r.total>0 AND r.created_at=(SELECT min(created_at) FROM test_runs WHERE problem_id=r.problem_id)),"+
			count("p.slug LIKE 'cf-%'")+" FROM problems p WHERE "+solvedSQL, today).
			Scan(&t.accepted, &t.hard, &t.acceptedXP, &t.independent, &t.completed, &t.recalls, &t.cards, &t.passedTests, &t.failedRuns, &t.late, &t.onTime, &t.firstTry, &t.codeforces)
		if err != nil {
			return err
		}
		rows, err := tx.Query(ctx, "SELECT DISTINCT substr(saved_at,1,10) FROM attempts WHERE state='saved' AND saved_at IS NOT NULL")
		if err != nil {
			return err
		}
		days, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			return err
		}
		t.streak = hasStreak(days)
		// The missed days run from the first visit to yesterday. Days in a row have the same date minus row number.
		if err := tx.QueryRow(ctx, "SELECT coalesce(($1::text::date-min(day)::date)-"+count("day<$1")+",0)::bigint,(SELECT coalesce(max(n),0) FROM (SELECT count(*) AS n FROM (SELECT day::date-(row_number() OVER (ORDER BY day))::int AS run FROM visits) runs GROUP BY run) lengths) FROM visits", today).Scan(&t.missed, &t.longest); err != nil {
			return err
		}
		rows, err = tx.Query(ctx, "SELECT p.slug FROM problems p WHERE p.slug=ANY($1) AND "+solvedSQL, all)
		if err != nil {
			return err
		}
		accepted, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			return err
		}
		t.cleared = 0
		for _, slugs := range setSlugs {
			if len(slugs) > 0 && !slices.ContainsFunc(slugs, func(slug string) bool { return !slices.Contains(accepted, slug) }) {
				t.cleared++
			}
		}
		result = gameOf(t)
		// An award stays after its condition ends. The level award follows the current XP, so it has no row.
		for _, b := range result.Badges {
			if b.Earned && b.ID != "level-five" {
				if _, err := tx.Exec(ctx, "INSERT INTO earned_badges (id) VALUES ($1) ON CONFLICT DO NOTHING", b.ID); err != nil {
					return err
				}
			}
		}
		rows, err = tx.Query(ctx, "SELECT id FROM earned_badges")
		if err != nil {
			return err
		}
		stored, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			return err
		}
		for i := range result.Badges {
			if slices.Contains(stored, result.Badges[i].ID) {
				result.Badges[i].Earned = true
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, result)
}

func (s *Server) reviews(w http.ResponseWriter, r *http.Request, user *SessionUser, problemID string, practice bool) error {
	ctx := r.Context()
	if _, err := uuid(problemID); err != nil {
		return err
	}
	allowed := []string{http.MethodGet, http.MethodPost}
	if practice {
		allowed = []string{http.MethodPost}
	}
	if err := methods(r, allowed...); err != nil {
		return err
	}
	if _, err := getProblem(ctx, s.db, problemID); err != nil {
		return err
	}
	if practice {
		today, err := readToday(r)
		if err != nil {
			return err
		}
		// A random day in [today + 15, today + 45].
		next, err := addDays(today, 15+s.random(31))
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(ctx, "UPDATE problems SET understanding='needs_practice',next_review_date=$1,progress_version=progress_version+1,updated_at=$2 WHERE id=$3", next, isoNow(), problemID); err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, map[string]string{"nextReviewDate": next})
	}
	if r.Method == http.MethodGet {
		limit, offset, err := pagination(r.URL.Query())
		if err != nil {
			return err
		}
		rows, err := s.db.Query(ctx, "SELECT "+reviewColumns+" FROM reviews WHERE problem_id=$1 ORDER BY created_at DESC,id DESC LIMIT $2 OFFSET $3", problemID, limit+1, offset)
		if err != nil {
			return err
		}
		items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (review, error) { return scanReview(row) })
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, page(items, limit, offset))
	}
	if err := requireRole(user, "student"); err != nil {
		return err
	}
	data, err := readObject(r, "version", "result", "note", "reviewedOn", "nextReviewDate")
	if err != nil {
		return err
	}
	next, present := data["nextReviewDate"]
	if !present {
		return invalid("Укажите правильную дату.")
	}
	date, err := calendarDate(next)
	if err != nil {
		return err
	}
	version, err := integer(data["version"], 1, maxSafeInteger)
	if err != nil {
		return err
	}
	result, err := choice(data["result"], "needs_practice", "with_help", "independent")
	if err != nil {
		return err
	}
	note, err := text(data["note"], 8000, false)
	if err != nil {
		return err
	}
	reviewedOn, err := todayValue(data["reviewedOn"])
	if err != nil {
		return err
	}
	id, now := newID(), isoNow()
	var created review
	err = s.tx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, "INSERT INTO reviews SELECT $1,id,$2,$3,$4,$5,$6 FROM problems WHERE id=$7 AND progress_version=$8", id, result, strings.TrimSpace(note), reviewedOn, date, now, problemID, version)
		if err != nil {
			return err
		}
		if err := changed(tag); err != nil {
			return err
		}
		// The update of the problem row makes two recalls with one version conflict.
		if _, err := tx.Exec(ctx, "UPDATE problems SET understanding=$1,next_review_date=$2,progress_version=progress_version+1,updated_at=$3 WHERE id=$4", result, date, now, problemID); err != nil {
			return err
		}
		created, err = scanReview(tx.QueryRow(ctx, "SELECT "+reviewColumns+" FROM reviews WHERE id=$1", id))
		return err
	})
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusCreated, created)
}

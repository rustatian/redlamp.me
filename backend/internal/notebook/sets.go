package notebook

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// chapterRef names one chapter of a learning set.
type chapterRef struct {
	Set      string `json:"set"`
	SetTitle string `json:"setTitle"`
	Number   int64  `json:"number"`
	Title    string `json:"title"`
}

// chapter returns the chapter number of the learning set setSlug.
func (s *Server) chapter(setSlug string, number int64) (sets.Chapter, bool) {
	for _, set := range s.sets {
		if set.Slug == setSlug && number >= 1 && number <= int64(len(set.Chapters)) {
			return set.Chapters[number-1], true
		}
	}
	return sets.Chapter{}, false
}

// set returns the learning set with the slug.
func (s *Server) set(slug string) (sets.LearningSet, bool) {
	for _, set := range s.sets {
		if set.Slug == slug {
			return set, true
		}
	}
	return sets.LearningSet{}, false
}

// setTask returns the first task of the learning sets with the problem slug.
func (s *Server) setTask(slug string) (sets.Task, bool) {
	for _, set := range s.sets {
		for _, task := range set.Tasks() {
			if task.Slug() == slug {
				return task, true
			}
		}
	}
	return sets.Task{}, false
}

// chaptersOf returns the chapters that have a task with the problem slug.
func (s *Server) chaptersOf(slug string) []chapterRef {
	var chapters []chapterRef
	for _, set := range s.sets {
		for i, chapter := range set.Chapters {
			for _, task := range chapter.Tasks {
				if task.Slug() == slug {
					chapters = append(chapters, chapterRef{Set: set.Slug, SetTitle: set.Title, Number: int64(i + 1), Title: chapter.Title})
					break
				}
			}
		}
	}
	return chapters
}

// acceptedSetSlugs returns the slugs of the learning set problems that have an accepted attempt.
func (s *Server) acceptedSetSlugs(ctx context.Context, q querier) ([]string, error) {
	var all []string
	for _, set := range s.sets {
		for _, task := range set.Tasks() {
			all = append(all, task.Slug())
		}
	}
	rows, err := q.Query(ctx, "SELECT p.slug FROM problems p WHERE p.slug=ANY($1) AND "+solvedSQL, all)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

var (
	setLinksPath = regexp.MustCompile(`^/leetcode/api/(problems|homework)/([^/]+)/sets(?:/([^/]+))?$`)
	setsPath     = regexp.MustCompile(`^/leetcode/api/sets(?:/([^/]+)(?:/lesson/([1-9][0-9]?))?)?$`)
)

func (s *Server) handleSets(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	if match := setLinksPath.FindStringSubmatch(r.URL.Path); match != nil {
		return true, s.setLinks(w, r, user, match[1] == "homework", match[2], match[3])
	}
	match := setsPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	if err := methods(r, http.MethodGet); err != nil {
		return true, err
	}
	if match[1] == "" {
		return true, s.setList(w, r)
	}
	set, ok := s.set(match[1])
	if !ok {
		return true, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такого набора нет."}
	}
	if match[2] != "" {
		number, _ := strconv.ParseInt(match[2], 10, 64)
		chapter, ok := s.chapter(set.Slug, number)
		if !ok {
			return true, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой главы нет."}
		}
		return true, writeJSON(w, r, http.StatusOK, struct {
			HTML    string        `json:"html"`
			Figures []sets.Figure `json:"figures"`
		}{chapter.Lesson, append([]sets.Figure{}, chapter.Figures...)})
	}
	return true, s.setDetails(w, r, set)
}

// setLinks reads and changes the links from a problem or a homework to learning sets.
// The read also lists the chapters of the problems and the chapters that the homework assigns.
func (s *Server) setLinks(w http.ResponseWriter, r *http.Request, user *SessionUser, homework bool, ownerID, slug string) error {
	ctx := r.Context()
	if _, err := uuid(ownerID); err != nil {
		return err
	}
	allowed := []string{http.MethodGet}
	if slug != "" {
		allowed = []string{http.MethodPut, http.MethodDelete}
	}
	if err := methods(r, allowed...); err != nil {
		return err
	}
	table, field := "problem_sets", "problem_id"
	var slugs []string
	if homework {
		table, field = "homework_sets", "homework_id"
		if _, err := getHomework(ctx, s.db, ownerID); err != nil {
			return err
		}
	} else {
		problem, err := getProblem(ctx, s.db, ownerID)
		if err != nil {
			return err
		}
		slugs = []string{problem.Slug}
	}
	if slug == "" {
		rows, err := s.db.Query(ctx, "SELECT set_slug FROM "+table+" WHERE "+field+"=$1", ownerID)
		if err != nil {
			return err
		}
		linked, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			return err
		}
		var chapters []chapterRef
		if homework {
			rows, err := s.db.Query(ctx, "SELECT set_slug,chapter FROM homework_tasks WHERE homework_id=$1 AND set_slug IS NOT NULL AND state<>'cancelled' ORDER BY created_at,id", ownerID)
			if err != nil {
				return err
			}
			var setSlug string
			var number int64
			if _, err := pgx.ForEachRow(rows, []any{&setSlug, &number}, func() error {
				set, _ := s.set(setSlug)
				chapter, ok := s.chapter(setSlug, number)
				if !ok {
					return fmt.Errorf("the homework %s names chapter %d of the learning set %q, which is not in the content", ownerID, number, setSlug)
				}
				chapters = append(chapters, chapterRef{Set: set.Slug, SetTitle: set.Title, Number: number, Title: chapter.Title})
				return nil
			}); err != nil {
				return err
			}
			if rows, err = s.db.Query(ctx, "SELECT p.slug FROM homework_tasks t JOIN problems p ON p.id=t.problem_id WHERE t.homework_id=$1 ORDER BY t.created_at,t.id", ownerID); err != nil {
				return err
			}
			if slugs, err = pgx.CollectRows(rows, pgx.RowTo[string]); err != nil {
				return err
			}
		}
		for _, problem := range slugs {
			chapters = append(chapters, s.chaptersOf(problem)...)
		}
		type setLink struct {
			Slug  string `json:"slug"`
			Title string `json:"title"`
		}
		result := struct {
			Chapters []chapterRef `json:"chapters"`
			Sets     []setLink    `json:"sets"`
		}{[]chapterRef{}, []setLink{}}
		// A chapter shows one time.
		seen := map[string]bool{}
		for _, chapter := range chapters {
			if key := fmt.Sprintf("%s/%d", chapter.Set, chapter.Number); !seen[key] {
				seen[key] = true
				result.Chapters = append(result.Chapters, chapter)
			}
		}
		for _, set := range s.sets {
			for _, link := range linked {
				if link == set.Slug {
					result.Sets = append(result.Sets, setLink{set.Slug, set.Title})
				}
			}
		}
		return writeJSON(w, r, http.StatusOK, result)
	}
	if err := requireRole(user, "parent"); err != nil {
		return err
	}
	if _, ok := s.set(slug); !ok {
		return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такого набора нет."}
	}
	change := "INSERT INTO " + table + " VALUES ($1,$2) ON CONFLICT DO NOTHING"
	if r.Method == http.MethodDelete {
		change = "DELETE FROM " + table + " WHERE " + field + "=$1 AND set_slug=$2"
	}
	if !homework {
		if _, err := s.db.Exec(ctx, change, ownerID, slug); err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	version, err := readVersion(r)
	if err != nil {
		return err
	}
	// A change of the links is a change of the homework, so it needs the current version.
	err = s.tx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, "UPDATE homework SET version=version+1,updated_at=$1 WHERE id=$2 AND version=$3", isoNow(), ownerID, version)
		if err != nil {
			return err
		}
		if err := changed(tag); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, change, ownerID, slug)
		return err
	})
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) setList(w http.ResponseWriter, r *http.Request) error {
	accepted, err := s.acceptedSetSlugs(r.Context(), s.db)
	if err != nil {
		return err
	}
	type summary struct {
		Slug          string   `json:"slug"`
		Title         string   `json:"title"`
		Summary       string   `json:"summary"`
		Topics        []string `json:"topics"`
		ChapterCount  int      `json:"chapterCount"`
		TaskCount     int      `json:"taskCount"`
		AcceptedCount int      `json:"acceptedCount"`
	}
	list := []summary{}
	for _, set := range s.sets {
		item := summary{Slug: set.Slug, Title: set.Title, Summary: set.Summary, Topics: set.Topics, ChapterCount: len(set.Chapters)}
		for _, task := range set.Tasks() {
			item.TaskCount++
			for _, slug := range accepted {
				if slug == task.Slug() {
					item.AcceptedCount++
				}
			}
		}
		list = append(list, item)
	}
	return writeJSON(w, r, http.StatusOK, list)
}

// setDetails writes a learning set with the state of each problem for the student: accepted, assignable, and its active homework task.
func (s *Server) setDetails(w http.ResponseWriter, r *http.Request, set sets.LearningSet) error {
	ctx := r.Context()
	var slugs []string
	for _, task := range set.Tasks() {
		slugs = append(slugs, task.Slug())
	}
	type libraryRow struct {
		problemID  string
		accepted   bool
		assignable bool
		state      *string
	}
	library := map[string]libraryRow{}
	rows, err := s.db.Query(ctx, "SELECT p.slug,p.id,"+solvedSQL+",("+assignableSQL+"),t.state FROM problems p LEFT JOIN homework_tasks t ON t.problem_id=p.id AND t.state NOT IN ('completed','cancelled') WHERE p.slug=ANY($1)", slugs)
	if err != nil {
		return err
	}
	var slug string
	var row libraryRow
	if _, err := pgx.ForEachRow(rows, []any{&slug, &row.problemID, &row.accepted, &row.assignable, &row.state}, func() error {
		library[slug] = row
		return nil
	}); err != nil {
		return err
	}
	decks, err := s.quizDecks(ctx, s.db)
	if err != nil {
		return err
	}
	type activeTask struct {
		State string `json:"state"`
	}
	type taskRow struct {
		sets.Task
		Slug       string      `json:"slug"`
		Label      string      `json:"label"`
		ProblemID  *string     `json:"problemId"`
		Accepted   bool        `json:"accepted"`
		Assignable bool        `json:"assignable"`
		ActiveTask *activeTask `json:"activeTask"`
	}
	type chapterDetails struct {
		Number  int         `json:"number"`
		Title   string      `json:"title"`
		Summary string      `json:"summary"`
		Links   []sets.Link `json:"links"`
		Quiz    quizDeck    `json:"quiz"`
		Tasks   []taskRow   `json:"tasks"`
	}
	details := struct {
		Slug     string           `json:"slug"`
		Title    string           `json:"title"`
		Summary  string           `json:"summary"`
		Topics   []string         `json:"topics"`
		Chapters []chapterDetails `json:"chapters"`
	}{set.Slug, set.Title, set.Summary, set.Topics, []chapterDetails{}}
	for i, chapter := range set.Chapters {
		item := chapterDetails{Number: i + 1, Title: chapter.Title, Summary: chapter.Summary, Links: chapter.Links, Tasks: []taskRow{}}
		for _, deck := range decks {
			if deck.Chapter.Set == set.Slug && deck.Chapter.Number == int64(i+1) {
				item.Quiz = deck
			}
		}
		for _, task := range chapter.Tasks {
			// A problem that is not in the library can be assigned: the assignment adds it.
			entry := taskRow{Task: task, Slug: task.Slug(), Label: task.Label(), Assignable: true}
			if row, ok := library[entry.Slug]; ok {
				entry.ProblemID, entry.Accepted, entry.Assignable = &row.problemID, row.accepted, row.assignable
				if row.state != nil {
					entry.ActiveTask = &activeTask{*row.state}
				}
			}
			item.Tasks = append(item.Tasks, entry)
		}
		details.Chapters = append(details.Chapters, item)
	}
	return writeJSON(w, r, http.StatusOK, details)
}

package notebook

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// timeLayout is the form of a stored time, for example 2026-10-07T12:34:56.789Z.
const timeLayout = "2006-01-02T15:04:05.000Z"

// deck is the quiz deck of one chapter: its cards and its names.
type deck struct {
	slug    string
	title   string
	chapter chapterRef
	cards   []sets.Card
}

// quizDeck is a quiz deck in the form of the API. Best and BestTotal are the result of the best passed attempt.
// NextAt is the end of the one-hour pause after the last attempt, or null.
type quizDeck struct {
	Slug      string     `json:"slug"`
	Title     string     `json:"title"`
	Chapter   chapterRef `json:"chapter"`
	CardCount int        `json:"cardCount"`
	Best      *int64     `json:"best"`
	BestTotal *int64     `json:"bestTotal"`
	NextAt    *string    `json:"nextAt"`
}

// quizCard is a card without its answer and its explanation.
type quizCard struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Prompt   string   `json:"prompt"`
	Language string   `json:"language,omitempty"`
	Code     string   `json:"code,omitempty"`
	Options  []string `json:"options,omitempty"`
}

// decksOf makes one deck for each chapter of the learning sets.
func decksOf(content []sets.LearningSet) []deck {
	var decks []deck
	for _, set := range content {
		for i, chapter := range set.Chapters {
			decks = append(decks, deck{
				slug:    fmt.Sprintf("%s-%d", set.Slug, i+1),
				title:   fmt.Sprintf("%s, глава %d", set.Title, i+1),
				chapter: chapterRef{Set: set.Slug, SetTitle: set.Title, Number: int64(i + 1), Title: chapter.Title},
				cards:   chapter.Cards,
			})
		}
	}
	return decks
}

// hourAfter returns the moment one hour after a stored time, or "" for a text that is not a time.
func hourAfter(stored string) (time.Time, string) {
	moment, err := time.Parse(timeLayout, stored)
	if err != nil {
		return time.Time{}, ""
	}
	moment = moment.Add(time.Hour)
	return moment, moment.Format(timeLayout)
}

// quizDecks returns each deck with the best passed attempt and the end of the pause.
func (s *Server) quizDecks(ctx context.Context, q querier) ([]quizDeck, error) {
	type state struct {
		last            string
		best, bestTotal *int64
	}
	states := map[string]state{}
	rows, err := q.Query(ctx, "SELECT deck,max(created_at),max(CASE WHEN passed=1 THEN correct END),max(CASE WHEN passed=1 THEN total END) FROM quiz_attempts GROUP BY deck")
	if err != nil {
		return nil, err
	}
	var slug string
	var row state
	if _, err := pgx.ForEachRow(rows, []any{&slug, &row.last, &row.best, &row.bestTotal}, func() error {
		states[slug] = row
		return nil
	}); err != nil {
		return nil, err
	}
	now := s.now()
	decks := []quizDeck{}
	for _, d := range s.decks {
		item := quizDeck{Slug: d.slug, Title: d.title, Chapter: d.chapter, CardCount: len(d.cards)}
		if row, ok := states[d.slug]; ok {
			item.Best, item.BestTotal = row.best, row.bestTotal
			if end, text := hourAfter(row.last); end.After(now) {
				item.NextAt = &text
			}
		}
		decks = append(decks, item)
	}
	return decks, nil
}

// normalize is the form of a short answer for the comparison: lower case, with one space between the words.
func normalize(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), " ")
}

// grade checks the answer of one card. shown is the right answer as the student sees it after a passed attempt.
func grade(card sets.Card, answer any) (correct bool, shown string, err error) {
	if card.Type == "short" {
		given, err := text(answer, 200, true)
		if err != nil {
			return false, "", err
		}
		right, _ := card.Answer.(string)
		return normalize(given) == normalize(right), right, nil
	}
	given, err := integer(answer, 0, int64(len(card.Options)-1))
	if err != nil {
		return false, "", err
	}
	right, _ := card.Answer.(int)
	return given == int64(right), card.Options[right], nil
}

var quizPath = regexp.MustCompile(`^/leetcode/api/quiz(?:/([^/]+)(?:/(cards|attempts))?)?$`)

func (s *Server) handleQuiz(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	match := quizPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	return true, s.quiz(w, r, user, match[1], match[2])
}

func (s *Server) quiz(w http.ResponseWriter, r *http.Request, user *SessionUser, slug, action string) error {
	ctx := r.Context()
	var d *deck
	for i := range s.decks {
		if s.decks[i].slug == slug {
			d = &s.decks[i]
		}
	}
	if slug != "" && d == nil {
		return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой колоды нет."}
	}
	switch action {
	case "cards":
		if err := methods(r, http.MethodGet); err != nil {
			return err
		}
		if err := requireRole(user, "parent"); err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, d.cards)
	case "attempts":
		if err := methods(r, http.MethodPost); err != nil {
			return err
		}
		if err := requireRole(user, "student"); err != nil {
			return err
		}
		return s.quizAttempt(w, r, d)
	}
	if err := methods(r, http.MethodGet); err != nil {
		return err
	}
	summaries, err := s.quizDecks(ctx, s.db)
	if err != nil {
		return err
	}
	if d == nil {
		return writeJSON(w, r, http.StatusOK, summaries)
	}
	round := struct {
		Title   string     `json:"title"`
		Chapter chapterRef `json:"chapter"`
		Cards   []quizCard `json:"cards"`
		NextAt  *string    `json:"nextAt"`
	}{Title: d.title, Chapter: d.chapter, Cards: []quizCard{}}
	for _, card := range d.cards {
		round.Cards = append(round.Cards, quizCard{ID: card.ID, Type: card.Type, Prompt: card.Prompt, Language: card.Language, Code: card.Code, Options: card.Options})
	}
	for _, summary := range summaries {
		if summary.Slug == d.slug {
			round.NextAt = summary.NextAt
		}
	}
	return writeJSON(w, r, http.StatusOK, round)
}

// quizAttempt stores one attempt of a deck. An attempt passes with 80 percent of right answers, and the deck then has a pause of one hour.
// A passed attempt completes the active quiz task of the chapter.
func (s *Server) quizAttempt(w http.ResponseWriter, r *http.Request, d *deck) error {
	ctx := r.Context()
	data, err := readObject(r, "answers")
	if err != nil {
		return err
	}
	ids := make([]string, len(d.cards))
	for i, card := range d.cards {
		ids[i] = card.ID
	}
	answers, err := object(data["answers"], ids...)
	if err != nil {
		return err
	}
	type result struct {
		CardID  string  `json:"cardId"`
		Correct bool    `json:"correct"`
		Answer  *string `json:"answer,omitempty"`
		Explain *string `json:"explain,omitempty"`
	}
	results := make([]result, len(d.cards))
	correct := 0
	for i, card := range d.cards {
		right, shown, err := grade(card, answers[card.ID])
		if err != nil {
			return err
		}
		if right {
			correct++
		}
		results[i] = result{CardID: card.ID, Correct: right, Answer: &shown, Explain: &d.cards[i].Explain}
	}
	total := len(d.cards)
	passed := correct*5 >= total*4
	passedValue := 0
	if passed {
		passedValue = 1
	} else {
		// A failed attempt does not show the right answers and the explanations.
		for i := range results {
			results[i].Answer, results[i].Explain = nil, nil
		}
	}
	now := s.now().UTC()
	moment, id := now.Format(timeLayout), newID()
	var nextAt string
	err = s.tx(ctx, func(tx pgx.Tx) error {
		// The write of the deck row makes two attempts of one deck at the same time conflict. The second one runs again and finds the pause.
		if _, err := tx.Exec(ctx, "INSERT INTO quiz_decks (deck,last_attempt_at) VALUES ($1,$2) ON CONFLICT (deck) DO UPDATE SET last_attempt_at=excluded.last_attempt_at", d.slug, moment); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, "INSERT INTO quiz_attempts SELECT $1,$2,$3,$4,$5,$6 WHERE NOT EXISTS (SELECT 1 FROM quiz_attempts WHERE deck=$2 AND created_at>$7)",
			id, d.slug, correct, total, passedValue, moment, now.Add(-time.Hour).Format(timeLayout))
		if err != nil {
			return err
		}
		var last string
		if err := tx.QueryRow(ctx, "SELECT max(created_at) FROM quiz_attempts WHERE deck=$1", d.slug).Scan(&last); err != nil {
			return err
		}
		_, nextAt = hourAfter(last)
		if tag.RowsAffected() == 0 {
			return &HTTPError{Status: http.StatusTooManyRequests, Code: "quiz_paused", Message: "После попытки квиз закрыт на час.", Fields: map[string]string{"nextAt": nextAt}}
		}
		if passed {
			_, err = tx.Exec(ctx, "UPDATE homework_tasks SET state='completed',completed_at=$1,updated_at=$1,version=version+1 WHERE kind='quiz' AND set_slug=$2 AND chapter=$3 AND state IN ('assigned','in_progress')", moment, d.chapter.Set, d.chapter.Number)
		}
		return err
	})
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusCreated, struct {
		Correct int      `json:"correct"`
		Total   int      `json:"total"`
		Passed  bool     `json:"passed"`
		NextAt  string   `json:"nextAt"`
		Results []result `json:"results"`
	}{correct, total, passed, nextAt, results})
}

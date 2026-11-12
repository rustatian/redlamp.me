package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf16"

	"github.com/jackc/pgx/v5"

	"github.com/rustatian/redlamp.me/backend/internal/answers"
	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// judgeOutput is the output of the student code on one hidden test: the answer text, the error message or a timeout.
type judgeOutput struct {
	Value   *string
	Error   *string
	Timeout bool
}

// judgement is the verdict of an acceptance check. Failed is the first failed test, or null for an accepted solution.
type judgement struct {
	Verdict string      `json:"verdict"`
	Failed  *failedTest `json:"failed"`
}

type failedTest struct {
	Index    int    `json:"index"`
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Output   string `json:"output"`
}

// expected is the answer of a reference solution to an own test, or the reason why it has none.
type expected struct {
	Value *string `json:"value,omitempty"`
	Error string  `json:"error,omitempty"`
}

// units returns the length of a text in UTF-16 code units, the length unit of the web client.
func units(value string) int {
	return len(utf16.Encode([]rune(value)))
}

// head returns the first max UTF-16 code units of a text, and true when the text is longer.
func head(value string, max int) (string, bool) {
	encoded := utf16.Encode([]rune(value))
	if len(encoded) <= max {
		return value, false
	}
	return string(utf16.Decode(encoded[:max])), true
}

// cut keeps the first 2000 characters of a long text.
func cut(value string) string {
	if start, long := head(value, 2000); long {
		return start + "…"
	}
	return value
}

// judge compares the outputs with the answers of the reference. The first failed test gives the verdict, so the reference answers only the tests up to it.
func judge(reference sets.Reference, spec sets.TestSpec, inputs []string, outputs []judgeOutput) judgement {
	for i, input := range inputs {
		want, output := reference.Answer(input), outputs[i]
		if output.Value != nil && answers.Same(answers.Spec{Mode: spec.Mode, Order: spec.Order, IgnoreCase: spec.IgnoreCase}, want, *output.Value) {
			continue
		}
		verdict, shown := "timeout", ""
		switch {
		case output.Value != nil:
			verdict, shown = "wrong", *output.Value
		case output.Error != nil:
			verdict, shown = "error", *output.Error
		}
		return judgement{Verdict: verdict, Failed: &failedTest{Index: i, Input: cut(input), Expected: cut(want), Output: cut(shown)}}
	}
	return judgement{Verdict: "accepted"}
}

// expectedOf returns the answer of the reference to an own test.
// The error is "constraints" for an input outside the constraints of the statement and "input" for an input of a wrong shape or type.
// The error is "size" for an answer of more than 100,000 characters, because the runner parses the answer in the time limit of the run.
func expectedOf(reference sets.Reference, input string) (result expected) {
	defer func() {
		if cause := recover(); cause != nil {
			result = expected{Error: "input"}
			if err, ok := cause.(error); ok && errors.Is(err, sets.ErrConstraints) {
				result.Error = "constraints"
			}
		}
	}()
	text := reference.Answer(input)
	if units(text) > 100_000 {
		return expected{Error: "size"}
	}
	return expected{Value: &text}
}

// outputOf checks one output of a judge request: an object with one of the keys value, error and timeout.
func outputOf(value any) (judgeOutput, error) {
	item, err := object(value, "value", "error", "timeout")
	if err != nil {
		return judgeOutput{}, err
	}
	answer, isValue := item["value"].(string)
	failure, isError := item["error"].(string)
	switch {
	case len(item) != 1:
	case item["timeout"] == true:
		return judgeOutput{Timeout: true}, nil
	case isValue && units(answer) <= 1_000_000:
		return judgeOutput{Value: &answer}, nil
	case isError && units(failure) <= 1_000_000:
		return judgeOutput{Error: &failure}, nil
	}
	return judgeOutput{}, invalid("Проверьте поля запроса.")
}

// hasWork reports whether a document has notes, an explanation or code.
func hasWork(doc attemptDocument) bool {
	texts := []string{doc.Notes}
	for _, a := range doc.Approaches {
		texts = append(texts, a.Idea, a.Correctness, a.Mistakes, a.Code)
	}
	for _, text := range texts {
		if strings.TrimSpace(text) != "" {
			return true
		}
	}
	return false
}

var judgePath = regexp.MustCompile(`^/leetcode/api/problems/([^/]+)/(judge/start|judge|solved|tests/expected)$`)

func (s *Server) handleJudge(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	match := judgePath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	return true, s.judgeRoute(w, r, user, match[1], match[2])
}

func (s *Server) judgeRoute(w http.ResponseWriter, r *http.Request, user *SessionUser, problemID, action string) error {
	ctx := r.Context()
	if err := methods(r, http.MethodPost); err != nil {
		return err
	}
	if action != "tests/expected" {
		if err := requireRole(user, "student"); err != nil {
			return err
		}
	}
	problem, err := getProblem(ctx, s.db, problemID)
	if err != nil {
		return err
	}
	if action == "tests/expected" {
		return s.expectedAnswers(w, r, problem.Slug)
	}
	generate := s.generators[problem.Slug]
	now := isoNow()
	var attemptID string
	var version int64
	result := judgement{Verdict: "accepted"}
	solved := action == "solved"
	switch {
	case solved:
		if generate != nil {
			return &HTTPError{Status: http.StatusConflict, Code: "judged", Message: "Эту задачу проверяют скрытые тесты. Нажмите «Отправить»."}
		}
		data, err := readObject(r, "attemptId")
		if err != nil {
			return err
		}
		if attemptID, err = uuid(data["attemptId"]); err != nil {
			return err
		}
		attempt, err := getAttempt(ctx, s.db, problem.ID, attemptID)
		if err != nil {
			return err
		}
		version = attempt.Version
		var doc attemptDocument
		if err := json.Unmarshal(attempt.Document, &doc); err != nil {
			return err
		}
		if !hasWork(doc) {
			return &HTTPError{Status: http.StatusBadRequest, Code: "empty_attempt", Message: "Добавьте объяснение, заметки или код перед отправкой."}
		}
	case generate == nil:
		return &HTTPError{Status: http.StatusNotFound, Code: "no_judge", Message: "У этой задачи нет скрытых тестов."}
	case action == "judge/start":
		seed := uint32(s.random(1 << 32))
		if _, err := s.db.Exec(ctx, "INSERT INTO judge_runs (problem_id,seed,created_at) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING", problem.ID, int64(seed), now); err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, struct {
			Seed   uint32   `json:"seed"`
			Inputs []string `json:"inputs"`
		}{seed, generate(sets.Mulberry32(seed))})
	default:
		// 16 outputs of 1,000,000 ASCII characters fit in 16 MiB.
		body, err := readJSON(r, 16_777_216)
		if err != nil {
			return err
		}
		data, err := object(body, "seed", "attemptId", "version", "outputs")
		if err != nil {
			return err
		}
		seed, err := integer(data["seed"], 0, 1<<32-1)
		if err != nil {
			return err
		}
		inputs := generate(sets.Mulberry32(uint32(seed)))
		if attemptID, err = uuid(data["attemptId"]); err != nil {
			return err
		}
		if version, err = integer(data["version"], 1, maxSafeInteger); err != nil {
			return err
		}
		raw, ok := data["outputs"].([]any)
		if !ok || len(raw) != len(inputs) {
			return invalid("Проверьте поля запроса.")
		}
		// A wrong verdict shows the answer of the failed test, so each seed gives one verdict.
		var open bool
		err = s.tx(ctx, func(tx pgx.Tx) error {
			tag, err := tx.Exec(ctx, "DELETE FROM judge_runs WHERE problem_id=$1 AND seed=$2", problem.ID, seed)
			open = err == nil && tag.RowsAffected() == 1
			return err
		})
		if err != nil {
			return err
		}
		if !open {
			return &HTTPError{Status: http.StatusConflict, Code: "stale_seed", Message: "Эта проверка устарела. Нажмите «Отправить» ещё раз."}
		}
		outputs := make([]judgeOutput, len(raw))
		for i, item := range raw {
			if outputs[i], err = outputOf(item); err != nil {
				return err
			}
		}
		// A task with a generator has a reference and a test spec.
		task, _ := s.setTask(problem.Slug)
		result = judge(s.references[problem.Slug], task.Tests.TestSpec, inputs, outputs)
		if result.Verdict != "accepted" {
			return writeJSON(w, r, http.StatusOK, result)
		}
	}
	if err := s.accept(ctx, problem.ID, attemptID, version, now, solved); err != nil {
		return err
	}
	if !solved {
		return writeJSON(w, r, http.StatusOK, result)
	}
	attempt, err := getAttempt(ctx, s.db, problem.ID, attemptID)
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, attempt)
}

// accept makes the draft the accepted attempt of the problem and changes the active homework task of the problem, in one transaction.
// The acceptance check completes the task. The solved route submits the task for the review of the parent.
func (s *Server) accept(ctx context.Context, problemID, attemptID string, version int64, now string, solved bool) error {
	states, state := "'assigned','in_progress','changes_requested','submitted'", "completed"
	if solved {
		states, state = "'assigned','in_progress','changes_requested'", "submitted"
	}
	submissionID := newID()
	return s.tx(ctx, func(tx pgx.Tx) error {
		// The understanding is "with help" when the student asked the tutor about the problem.
		tag, err := tx.Exec(ctx, "UPDATE attempts SET state='saved',acceptance='accepted',understanding=CASE WHEN EXISTS(SELECT 1 FROM tutor_messages WHERE context_kind='problem' AND context_ref=$2 AND author='student') THEN 'with_help' ELSE 'independent' END,version=version+1,saved_at=$3,updated_at=$3 WHERE id=$1 AND problem_id=$2 AND state='draft' AND version=$4",
			attemptID, problemID, now, version)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			if _, err := getAttempt(ctx, tx, problemID, attemptID); err != nil {
				return err
			}
			return changed(tag)
		}
		// The write of the problem row makes the acceptance conflict with an archive, a slug change or a new homework task of the problem.
		if _, err := tx.Exec(ctx, "UPDATE problems SET understanding=(SELECT understanding FROM attempts WHERE id=$1),progress_version=progress_version+1,updated_at=$2 WHERE id=$3", attemptID, now, problemID); err != nil {
			return err
		}
		tag, err = tx.Exec(ctx, "INSERT INTO submissions (id,task_id,attempt_id,task_version,created_at) SELECT $1,id,$2,version+1,$3 FROM homework_tasks WHERE problem_id=$4 AND state IN ("+states+")", submissionID, attemptID, now, problemID)
		if err != nil || tag.RowsAffected() == 0 {
			return err
		}
		if !solved {
			if _, err := tx.Exec(ctx, "INSERT INTO feedback (id,problem_id,task_id,submission_id,author,kind,body,created_at) SELECT $1,$2,task_id,id,'student','completed','Принято проверкой.',$3 FROM submissions WHERE id=$4", newID(), problemID, now, submissionID); err != nil {
				return err
			}
		}
		_, err = tx.Exec(ctx, "UPDATE homework_tasks SET state=$1,current_submission_id=$2,version=version+1,updated_at=$3 WHERE id=(SELECT task_id FROM submissions WHERE id=$2)", state, submissionID, now)
		return err
	})
}

// expectedAnswers answers the inputs of own tests with the reference of the problem. A problem without a reference gets null for each input.
func (s *Server) expectedAnswers(w http.ResponseWriter, r *http.Request, slug string) error {
	// 20 inputs of 10,000 characters at two bytes each fit in 400 KB of JSON.
	body, err := readJSON(r, 409_600)
	if err != nil {
		return err
	}
	data, err := object(body, "inputs")
	if err != nil {
		return err
	}
	inputs, ok := data["inputs"].([]any)
	if !ok || len(inputs) > 20 {
		return invalid("Проверьте поля запроса.")
	}
	// The inputs can be the hidden inputs of an open seed. A slug can move to another problem, so the answers make the seeds of all problems stale.
	err = s.tx(r.Context(), func(tx pgx.Tx) error {
		_, err := tx.Exec(r.Context(), "DELETE FROM judge_runs")
		return err
	})
	if err != nil {
		return err
	}
	reference, hasReference := s.references[slug]
	results := make([]*expected, len(inputs))
	for i, item := range inputs {
		input, err := text(item, 10_000, true)
		if err != nil {
			return err
		}
		if hasReference {
			answer := expectedOf(reference, input)
			results[i] = &answer
		}
	}
	return writeJSON(w, r, http.StatusOK, map[string][]*expected{"expected": results})
}

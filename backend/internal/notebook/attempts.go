package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

// maxDocumentBytes is the size limit of the document of an attempt.
const maxDocumentBytes = 256_000

const attemptColumns = "id,problem_id,state,acceptance,understanding,version,created_at,updated_at,saved_at"

// Attempt is an attempt in the form of the API. A list has no Document.
type Attempt struct {
	ID            string  `json:"id"`
	ProblemID     string  `json:"problemId"`
	State         string  `json:"state"`
	Acceptance    string  `json:"acceptance"`
	Understanding string  `json:"understanding"`
	Version       int64   `json:"version"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
	SavedAt       *string `json:"savedAt"`
	// Document is the stored JSON text, so a read returns the stored bytes.
	Document json.RawMessage `json:"document,omitempty"`
}

// approach is one approach of an attempt document. The field order is the key order of the web client.
type approach struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	Idea            string `json:"idea"`
	Correctness     string `json:"correctness"`
	TimeComplexity  string `json:"timeComplexity"`
	SpaceComplexity string `json:"spaceComplexity"`
	EdgeCases       string `json:"edgeCases"`
	Mistakes        string `json:"mistakes"`
	Language        string `json:"language"`
	Code            string `json:"code"`
}

type attemptDocument struct {
	Notes      string     `json:"notes"`
	Approaches []approach `json:"approaches"`
}

func scanAttempt(row pgx.Row, withDocument bool) (Attempt, error) {
	var a Attempt
	fields := []any{&a.ID, &a.ProblemID, &a.State, &a.Acceptance, &a.Understanding, &a.Version, &a.CreatedAt, &a.UpdatedAt, &a.SavedAt}
	var document string
	if withDocument {
		fields = append(fields, &document)
	}
	if err := row.Scan(fields...); err != nil {
		return a, err
	}
	if withDocument {
		a.Document = json.RawMessage(document)
	}
	return a, nil
}

func getAttempt(ctx context.Context, q querier, problemID, id string) (Attempt, error) {
	if _, err := uuid(id); err != nil {
		return Attempt{}, err
	}
	if _, err := uuid(problemID); err != nil {
		return Attempt{}, err
	}
	a, err := scanAttempt(q.QueryRow(ctx, "SELECT "+attemptColumns+",document FROM attempts WHERE id=$1 AND problem_id=$2", id, problemID), true)
	if errors.Is(err, pgx.ErrNoRows) {
		return a, &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой попытки нет."}
	}
	return a, err
}

var approachKeys = []string{"id", "label", "idea", "correctness", "timeComplexity", "spaceComplexity", "edgeCases", "mistakes", "language", "code"}

// documentValue checks the decoded JSON of an attempt document: the notes and one to eight approaches with different IDs.
// It returns the checked value, so the stored text is a serialization of exactly this value.
func documentValue(value any) (attemptDocument, error) {
	var doc attemptDocument
	fields, err := object(value, "notes", "approaches")
	if err != nil {
		return doc, err
	}
	if doc.Notes, err = text(fields["notes"], maxDocumentBytes, false); err != nil {
		return doc, err
	}
	items, ok := fields["approaches"].([]any)
	if !ok || len(items) < 1 || len(items) > 8 {
		return doc, invalid("Укажите от одного до восьми подходов.")
	}
	ids := map[string]bool{}
	for _, item := range items {
		entry, err := object(item, approachKeys...)
		if err != nil {
			return doc, err
		}
		var a approach
		if a.ID, err = uuid(entry["id"]); err != nil {
			return doc, err
		}
		if ids[a.ID] {
			return doc, invalid("У каждого подхода должен быть свой идентификатор.")
		}
		ids[a.ID] = true
		for i, field := range []*string{&a.Label, &a.Idea, &a.Correctness, &a.TimeComplexity, &a.SpaceComplexity, &a.EdgeCases, &a.Mistakes, &a.Language, &a.Code} {
			limit := maxDocumentBytes
			switch approachKeys[i+1] {
			case "label":
				limit = 200
			case "language":
				limit = 40
			}
			if *field, err = text(entry[approachKeys[i+1]], limit, false); err != nil {
				return doc, err
			}
		}
		doc.Approaches = append(doc.Approaches, a)
	}
	return doc, nil
}

// attemptSearch is the search text of a document: the notes and the explanations of each approach, without the code.
func attemptSearch(doc attemptDocument) string {
	parts := []string{doc.Notes}
	for _, a := range doc.Approaches {
		parts = append(parts, a.Label, a.Idea, a.Correctness, a.TimeComplexity, a.SpaceComplexity, a.EdgeCases, a.Mistakes)
	}
	return searchText(strings.Join(parts, " "))
}

var attemptsPath = regexp.MustCompile(`^/leetcode/api/problems/([^/]+)/attempts(?:/([^/]+)(/save)?)?$`)

func (s *Server) handleAttempts(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	match := attemptsPath.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return false, nil
	}
	return true, s.attempts(w, r, user, match[1], match[2], match[3] != "")
}

func (s *Server) attempts(w http.ResponseWriter, r *http.Request, user *SessionUser, problemID, id string, save bool) error {
	ctx := r.Context()
	if _, err := uuid(problemID); err != nil {
		return err
	}
	allowed := []string{http.MethodGet, http.MethodPost}
	if id != "" {
		if _, err := uuid(id); err != nil {
			return err
		}
		allowed = []string{http.MethodGet, http.MethodPut}
		if save {
			allowed = []string{http.MethodPost}
		}
	}
	if err := methods(r, allowed...); err != nil {
		return err
	}
	if r.Method == http.MethodGet || id == "" {
		if _, err := getProblem(ctx, s.db, problemID); err != nil {
			return err
		}
	}
	if r.Method == http.MethodGet {
		if id != "" {
			attempt, err := getAttempt(ctx, s.db, problemID, id)
			if err != nil {
				return err
			}
			return writeJSON(w, r, http.StatusOK, attempt)
		}
		limit, offset, err := pagination(r.URL.Query())
		if err != nil {
			return err
		}
		rows, err := s.db.Query(ctx, "SELECT "+attemptColumns+" FROM attempts WHERE problem_id=$1 ORDER BY state ASC,created_at DESC,id LIMIT $2 OFFSET $3", problemID, limit+1, offset)
		if err != nil {
			return err
		}
		items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Attempt, error) { return scanAttempt(row, false) })
		if err != nil {
			return err
		}
		return writeJSON(w, r, http.StatusOK, page(items, limit, offset))
	}
	if err := requireRole(user, "student"); err != nil {
		return err
	}
	now := isoNow()
	switch {
	case id == "":
		return s.startAttempt(w, r, problemID, now)
	case save:
		body, err := readJSON(r, 32_768)
		if err != nil {
			return err
		}
		data, err := object(body, "version")
		if err != nil {
			return err
		}
		version, err := integer(data["version"], 1, maxSafeInteger)
		if err != nil {
			return err
		}
		// A draft holds a copied or a default understanding, so the save keeps the progress of the problem.
		return s.changeDraft(w, r, problemID, id, "Эта запись изменилась. Обновите страницу и повторите.",
			"UPDATE attempts SET state='saved',version=version+1,saved_at=$1,updated_at=$1 WHERE id=$2 AND problem_id=$3 AND state='draft' AND version=$4", now, id, problemID, version)
	}
	// Only the judge and the solved route accept an attempt, so a draft update changes only the document.
	body, err := readJSON(r, 272_384)
	if err != nil {
		return err
	}
	data, err := object(body, "version", "document")
	if err != nil {
		return err
	}
	doc, err := documentValue(data["document"])
	if err != nil {
		return err
	}
	// The stored text is the checked value in the normal form: no white space, and the keys in the order of the web client.
	stored, err := stringify(doc)
	if err != nil {
		return err
	}
	if len(stored) > maxDocumentBytes {
		return &HTTPError{Status: http.StatusRequestEntityTooLarge, Code: "too_large", Message: "Размер попытки не должен превышать 256 000 байт."}
	}
	version, err := integer(data["version"], 1, maxSafeInteger)
	if err != nil {
		return err
	}
	return s.changeDraft(w, r, problemID, id, "Этот черновик изменился. Сравните свою работу с сохранённой копией.",
		"UPDATE attempts SET document=$1,search_text=$2,version=version+1,updated_at=$3 WHERE id=$4 AND problem_id=$5 AND state='draft' AND version=$6", string(stored), attemptSearch(doc), now, id, problemID, version)
}

// changeDraft runs an update that has the draft state and the version in its condition, and writes the attempt.
// When the update changes no row, the answer is 404 for an attempt that does not exist and the conflict error for a stale version.
func (s *Server) changeDraft(w http.ResponseWriter, r *http.Request, problemID, id, conflict, sql string, args ...any) error {
	ctx := r.Context()
	var attempt Attempt
	err := s.tx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		if attempt, err = getAttempt(ctx, tx, problemID, id); err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return &HTTPError{Status: http.StatusConflict, Code: "conflict", Message: conflict}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, attempt)
}

// startAttempt makes the draft of a problem: an empty document, or a copy of a saved attempt with new approach IDs.
// The unique index one_draft_per_problem allows one draft for each problem.
func (s *Server) startAttempt(w http.ResponseWriter, r *http.Request, problemID, now string) error {
	ctx := r.Context()
	body, err := readJSON(r, 32_768)
	if err != nil {
		return err
	}
	data, err := object(body, "copyAttemptId")
	if err != nil {
		return err
	}
	doc := attemptDocument{Approaches: []approach{{ID: newID()}}}
	acceptance, understanding := "not_submitted", "needs_practice"
	if copyID := data["copyAttemptId"]; copyID != nil && copyID != "" && copyID != false && copyID != 0.0 {
		sourceID, err := uuid(copyID)
		if err != nil {
			return err
		}
		source, err := getAttempt(ctx, s.db, problemID, sourceID)
		if err != nil {
			return err
		}
		if source.State != "saved" {
			return &HTTPError{Status: http.StatusConflict, Code: "conflict", Message: "Копировать можно только сохранённую попытку."}
		}
		if err := json.Unmarshal(source.Document, &doc); err != nil {
			return err
		}
		for i := range doc.Approaches {
			doc.Approaches[i].ID = newID()
		}
		acceptance, understanding = source.Acceptance, source.Understanding
	}
	document, err := stringify(doc)
	if err != nil {
		return err
	}
	id := newID()
	var attempt Attempt
	err = s.tx(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'draft',$3,$4,$5,$6,$7,$7) ON CONFLICT DO NOTHING",
			id, problemID, string(document), attemptSearch(doc), acceptance, understanding, now)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			var draft string
			if err := tx.QueryRow(ctx, "SELECT id FROM attempts WHERE problem_id=$1 AND state='draft'", problemID).Scan(&draft); err != nil {
				return err
			}
			return &HTTPError{Status: http.StatusConflict, Code: "draft_exists", Message: "Откройте текущий черновик, прежде чем начинать новую попытку.", Fields: map[string]string{"attemptId": draft}}
		}
		attempt, err = getAttempt(ctx, tx, problemID, id)
		return err
	})
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusCreated, attempt)
}

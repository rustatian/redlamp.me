package notebook

import (
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// attemptsApproach is one approach of an attempt document. The fields have the key order of test/attempts.test.ts.
type attemptsApproach struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	Idea            string `json:"idea"`
	Correctness     string `json:"correctness"`
	TimeComplexity  string `json:"timeComplexity"`
	SpaceComplexity string `json:"spaceComplexity"`
	EdgeCases       string `json:"edgeCases"`
	Mistakes        string `json:"mistakes"`
	Language        string `json:"language"`
	Code            any    `json:"code"`
}

type attemptsDocument struct {
	Notes      string             `json:"notes"`
	Approaches []attemptsApproach `json:"approaches"`
}

// attemptsFirst is the approach of test/attempts.test.ts.
var attemptsFirst = attemptsApproach{ID: "11111111-1111-4111-8111-111111111111", Label: "Counting", Idea: "Count the needed values.", Language: "Python", Code: "  return count\n"}

// attemptsBase is the document of test/attempts.test.ts.
var attemptsBase = attemptsDocument{Notes: "", Approaches: []attemptsApproach{attemptsFirst}}

// attemptsStored is JSON.stringify(document): the text of the draft in the database.
const attemptsStored = `{"notes":"","approaches":[{"id":"11111111-1111-4111-8111-111111111111","label":"Counting","idea":"Count the needed values.","correctness":"","timeComplexity":"","spaceComplexity":"","edgeCases":"","mistakes":"","language":"Python","code":"  return count\n"}]}`

// attemptsFixture is the state that beforeEach of test/attempts.test.ts makes.
type attemptsFixture struct {
	pool      *pgxpool.Pool
	s         *Server
	cookie    string
	problemID string
	attemptID string
	path      string
}

// attemptsSetup signs in the student, adds the problem two-sum and stores one draft of it.
func attemptsSetup(t *testing.T) attemptsFixture {
	t.Helper()
	pool := database(t)
	s := newServer(t, pool)
	cookie := loginAs(t, s, "student")
	problemID := addProblem(t, s, cookie, map[string]any{"url": "https://leetcode.com/problems/two-sum/", "number": 1, "title": "Two Sum", "difficulty": "easy", "topics": []string{}, "summary": ""})
	attemptID := newID()
	exec(t, pool, "INSERT INTO attempts (id,problem_id,state,document,search_text,acceptance,understanding,created_at,updated_at) VALUES ($1,$2,'draft',$3,'','accepted','with_help','now','now')", attemptID, problemID, attemptsStored)
	return attemptsFixture{pool: pool, s: s, cookie: cookie, problemID: problemID, attemptID: attemptID, path: "/problems/" + problemID + "/attempts/" + attemptID}
}

func TestAttemptsDraftVersion(t *testing.T) {
	tests := []struct {
		name          string
		state         string
		storedVersion int
		sentVersion   int
		wantStatus    int
		wantVersion   int64
		wantNotes     string
	}{
		{name: "fresh draft update", state: "draft", storedVersion: 1, sentVersion: 1, wantStatus: 200, wantVersion: 2, wantNotes: "New note"},
		{name: "stale draft update", state: "draft", storedVersion: 2, sentVersion: 1, wantStatus: 409, wantVersion: 2, wantNotes: ""},
		{name: "saved attempt update", state: "saved", storedVersion: 1, sentVersion: 1, wantStatus: 409, wantVersion: 1, wantNotes: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := attemptsSetup(t)
			exec(t, f.pool, "UPDATE attempts SET state = $1, version = $2", tt.state, tt.storedVersion)
			body := map[string]any{"document": attemptsDocument{Notes: "New note", Approaches: []attemptsApproach{attemptsFirst}}, "version": tt.sentVersion}
			if w := request(t, f.s, f.path, http.MethodPut, body, f.cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			var version int64
			var text string
			if err := f.pool.QueryRow(context.Background(), "SELECT version, document FROM attempts WHERE id = $1", f.attemptID).Scan(&version, &text); err != nil {
				t.Fatal(err)
			}
			if version != tt.wantVersion {
				t.Errorf("version = %d, want %d", version, tt.wantVersion)
			}
			var document attemptsDocument
			if err := json.Unmarshal([]byte(text), &document); err != nil {
				t.Fatal(err)
			}
			if document.Notes != tt.wantNotes {
				t.Errorf("notes = %q, want %q", document.Notes, tt.wantNotes)
			}
		})
	}
}

func TestAttemptsDraftGrades(t *testing.T) {
	tests := []struct {
		name       string
		extra      map[string]any
		wantStatus int
	}{
		{name: "draft update keeps the stored acceptance and understanding", extra: nil, wantStatus: 200},
		{name: "draft update with an acceptance and an understanding", extra: map[string]any{"acceptance": "accepted", "understanding": "independent"}, wantStatus: 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := attemptsSetup(t)
			exec(t, f.pool, "UPDATE attempts SET acceptance = 'not_submitted', understanding = 'needs_practice'")
			body := map[string]any{"document": attemptsBase, "version": 1}
			maps.Copy(body, tt.extra)
			if w := request(t, f.s, f.path, http.MethodPut, body, f.cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
			var acceptance, understanding string
			if err := f.pool.QueryRow(context.Background(), "SELECT acceptance, understanding FROM attempts WHERE id = $1", f.attemptID).Scan(&acceptance, &understanding); err != nil {
				t.Fatal(err)
			}
			if acceptance != "not_submitted" || understanding != "needs_practice" {
				t.Errorf("acceptance = %s, understanding = %s; want not_submitted and needs_practice", acceptance, understanding)
			}
		})
	}
}

func TestAttemptsDraftApproaches(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		duplicate  bool
		wantStatus int
	}{
		{name: "zero approaches", count: 0, wantStatus: 400},
		{name: "eight approaches", count: 8, wantStatus: 200},
		{name: "nine approaches", count: 9, wantStatus: 400},
		{name: "duplicate approach IDs", count: 2, duplicate: true, wantStatus: 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := attemptsSetup(t)
			approaches := make([]attemptsApproach, 0, tt.count)
			for range tt.count {
				approach := attemptsFirst
				if !tt.duplicate {
					approach.ID = newID()
				}
				approaches = append(approaches, approach)
			}
			body := map[string]any{"document": attemptsDocument{Notes: "", Approaches: approaches}, "version": 1}
			if w := request(t, f.s, f.path, http.MethodPut, body, f.cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

// The document limit is 256,000 bytes of JSON text.
func TestAttemptsDocumentSize(t *testing.T) {
	tests := []struct {
		name       string
		extra      int
		multibyte  bool
		wantStatus int
	}{
		{name: "exact document boundary", extra: 0, multibyte: false, wantStatus: 200},
		{name: "one byte over document boundary", extra: 1, multibyte: false, wantStatus: 413},
		{name: "multibyte overflow", extra: 1, multibyte: true, wantStatus: 413},
	}
	base, err := json.Marshal(attemptsBase)
	if err != nil {
		t.Fatal(err)
	}
	remaining := 256_000 - len(base)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := attemptsSetup(t)
			notes := strings.Repeat("x", remaining+tt.extra)
			if tt.multibyte {
				// "я" is two bytes in UTF-8.
				notes = strings.Repeat("я", (remaining+1)/2+tt.extra)
			}
			body := map[string]any{"document": attemptsDocument{Notes: notes, Approaches: []attemptsApproach{attemptsFirst}}, "version": 1}
			if w := request(t, f.s, f.path, http.MethodPut, body, f.cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestAttemptsDraftEditChecks(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		unknownProblem bool
		code           any
		wantStatus     int
	}{
		{name: "non-string approach code", role: "student", unknownProblem: false, code: 42, wantStatus: 400},
		{name: "parent draft edit", role: "parent", unknownProblem: false, code: attemptsFirst.Code, wantStatus: 403},
		{name: "draft edit under an unknown problem", role: "student", unknownProblem: true, code: attemptsFirst.Code, wantStatus: 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := attemptsSetup(t)
			target := f.path
			if tt.unknownProblem {
				target = "/problems/" + newID() + "/attempts/" + f.attemptID
			}
			cookie := f.cookie
			if tt.role == "parent" {
				cookie = loginAs(t, f.s, "parent")
			}
			approach := attemptsFirst
			approach.Code = tt.code
			body := map[string]any{"document": attemptsDocument{Notes: "", Approaches: []attemptsApproach{approach}}, "version": 1}
			if w := request(t, f.s, target, http.MethodPut, body, cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestAttemptsSaveKeepsProblemProgress(t *testing.T) {
	f := attemptsSetup(t)
	exec(t, f.pool, "UPDATE problems SET understanding = 'needs_practice', next_review_date = '2026-10-20'")
	w := request(t, f.s, f.path+"/save", http.MethodPost, map[string]any{"version": 1}, f.cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("save: status %d, body %s", w.Code, w.Body.String())
	}
	saved := decode[struct {
		State    string          `json:"state"`
		Version  int             `json:"version"`
		Document json.RawMessage `json:"document"`
	}](t, w)
	if saved.State != "saved" || saved.Version != 2 {
		t.Errorf("state = %q, version = %d; want saved and 2", saved.State, saved.Version)
	}
	var got, want any
	if err := json.Unmarshal(saved.Document, &got); err != nil {
		t.Fatalf("document %s: %v", saved.Document, err)
	}
	if err := json.Unmarshal([]byte(attemptsStored), &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("document = %s, want %s", saved.Document, attemptsStored)
	}
	if string(saved.Document) != attemptsStored {
		t.Errorf("document text = %s, want the stored text %s", saved.Document, attemptsStored)
	}
	progress := decode[map[string]any](t, request(t, f.s, "/problems/"+f.problemID, http.MethodGet, nil, f.cookie))
	for key, value := range map[string]any{"solved": true, "understanding": "needs_practice", "nextReviewDate": "2026-10-20", "version": float64(1), "progressVersion": float64(1)} {
		if progress[key] != value {
			t.Errorf("problem %s = %v, want %v", key, progress[key], value)
		}
	}
	if w := request(t, f.s, f.path+"/save", http.MethodPost, map[string]any{"version": 1}, f.cookie); w.Code != http.StatusConflict {
		t.Errorf("second save: status %d, want 409", w.Code)
	}
}

func TestAttemptsCopyOfSavedHistory(t *testing.T) {
	f := attemptsSetup(t)
	request(t, f.s, f.path+"/save", http.MethodPost, map[string]any{"version": 1}, f.cookie)
	copied := request(t, f.s, "/problems/"+f.problemID+"/attempts", http.MethodPost, map[string]any{"copyAttemptId": f.attemptID}, f.cookie)
	if copied.Code != http.StatusCreated {
		t.Fatalf("copy: status %d, body %s", copied.Code, copied.Body.String())
	}
	record := decode[struct {
		ID       string           `json:"id"`
		Document attemptsDocument `json:"document"`
	}](t, copied)
	if record.ID == f.attemptID {
		t.Errorf("copy ID = %s, want a new ID", record.ID)
	}
	if len(record.Document.Approaches) == 0 {
		t.Fatalf("copy body %s, want the approaches", copied.Body.String())
	}
	if first := record.Document.Approaches[0]; first.ID == attemptsFirst.ID || first.Code != "  return count\n" {
		t.Errorf("copied approach ID = %s, code = %q; want a new ID and the code %q", first.ID, first.Code, "  return count\n")
	}
	again := request(t, f.s, "/problems/"+f.problemID+"/attempts", http.MethodPost, map[string]any{"copyAttemptId": f.attemptID}, f.cookie)
	if again.Code != http.StatusConflict {
		t.Fatalf("second copy: status %d, want 409", again.Code)
	}
	existing := decode[struct {
		Error struct {
			Fields struct {
				AttemptID string `json:"attemptId"`
			} `json:"fields"`
		} `json:"error"`
	}](t, again)
	if existing.Error.Fields.AttemptID != record.ID {
		t.Errorf("error fields attemptId = %q, want the draft %s", existing.Error.Fields.AttemptID, record.ID)
	}
}

func TestAttemptsCopySearch(t *testing.T) {
	tests := []struct {
		name  string
		q     string
		found bool
	}{
		{name: "copied draft hides JSON keys from search", q: "timecomplexity", found: false},
		{name: "copied draft hides code from search", q: "return", found: false},
		{name: "copied draft finds the key idea", q: "needed values", found: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := attemptsSetup(t)
			if w := request(t, f.s, f.path+"/save", http.MethodPost, map[string]any{"version": 1}, f.cookie); w.Code != http.StatusOK {
				t.Fatalf("save: status %d, body %s", w.Code, w.Body.String())
			}
			if w := request(t, f.s, "/problems/"+f.problemID+"/attempts", http.MethodPost, map[string]any{"copyAttemptId": f.attemptID}, f.cookie); w.Code != http.StatusCreated {
				t.Fatalf("copy: status %d, body %s", w.Code, w.Body.String())
			}
			w := request(t, f.s, "/problems?"+url.Values{"q": {tt.q}}.Encode(), http.MethodGet, nil, f.cookie)
			got := decode[struct {
				Items []struct {
					ID string `json:"id"`
				} `json:"items"`
			}](t, w)
			if tt.found && (len(got.Items) != 1 || got.Items[0].ID != f.problemID) {
				t.Errorf("body = %s, want the one item %s", w.Body.String(), f.problemID)
			}
			if !tt.found && (got.Items == nil || len(got.Items) != 0) {
				t.Errorf("body = %s, want no items", w.Body.String())
			}
		})
	}
}

func TestAttemptsMissingTargets(t *testing.T) {
	tests := []struct {
		name       string
		suffix     string
		body       map[string]any
		wantStatus int
	}{
		{name: "save of a missing attempt", suffix: "/attempts/" + newID() + "/save", body: map[string]any{"version": 1}, wantStatus: 404},
		{name: "new draft while a draft exists", suffix: "/attempts", body: map[string]any{}, wantStatus: 409},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := attemptsSetup(t)
			if w := request(t, f.s, "/problems/"+f.problemID+tt.suffix, http.MethodPost, tt.body, f.cookie); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

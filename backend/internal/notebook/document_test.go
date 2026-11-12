package notebook

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

const (
	documentApproach = `{"id":"11111111-1111-4111-8111-111111111111","label":"Counting","idea":"","correctness":"","timeComplexity":"","spaceComplexity":"","edgeCases":"","mistakes":"","language":"Python","code":""}`
	// documentBase is the length of the stored text of a document with empty notes and documentApproach.
	documentBase = len(`{"notes":"","approaches":[` + documentApproach + `]}`)
)

// The Worker checks the parsed document and stores JSON.stringify of that value (src/leetcode/attempts.ts).
// So the stored text holds only checked values, in the normal form, and the limit of 256000 bytes counts that text.
func TestDraftUpdateStoresTheCheckedDocument(t *testing.T) {
	tests := []struct {
		name       string
		document   string
		wantStatus int
		wantCode   string
		wantStored string
	}{
		{
			name:       "a document in the form of the web client is stored as it is",
			document:   `{"notes":"plan","approaches":[` + documentApproach + `]}`,
			wantStatus: 200,
			wantStored: `{"notes":"plan","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "a duplicate key of the document keeps the last value",
			document:   `{"notes":"first","notes":"last","approaches":[` + documentApproach + `]}`,
			wantStatus: 200,
			wantStored: `{"notes":"last","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "a duplicate key of an approach keeps the last value",
			document:   `{"notes":"","approaches":[` + strings.Replace(documentApproach, `"label":"Counting"`, `"label":"First","label":"Counting"`, 1) + `]}`,
			wantStatus: 200,
			wantStored: `{"notes":"","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "a duplicate key with a value that is not text is rejected when it is the last one",
			document:   `{"notes":"text","notes":1,"approaches":[` + documentApproach + `]}`,
			wantStatus: 400,
			wantCode:   "invalid_input",
		},
		{
			name:       "an unknown key of the document is rejected",
			document:   `{"notes":"","approaches":[` + documentApproach + `],"extra":""}`,
			wantStatus: 400,
			wantCode:   "invalid_input",
		},
		{
			name:       "an unknown key of an approach is rejected",
			document:   `{"notes":"","approaches":[` + strings.Replace(documentApproach, `"code":""`, `"code":"","extra":""`, 1) + `]}`,
			wantStatus: 400,
			wantCode:   "invalid_input",
		},
		{
			name:       "an escaped character is stored as the character",
			document:   `{"notes":"café сумма","approaches":[` + documentApproach + `]}`,
			wantStatus: 200,
			wantStored: `{"notes":"café сумма","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "white space between the tokens is not stored",
			document:   "{\n  \"notes\" : \"a b\",\n  \"approaches\" : [ " + documentApproach + " ]\n}",
			wantStatus: 200,
			wantStored: `{"notes":"a b","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "keys in another order are stored in the order of the web client",
			document:   `{"approaches":[{"code":"","language":"Python","mistakes":"","edgeCases":"","spaceComplexity":"","timeComplexity":"","correctness":"","idea":"","label":"Counting","id":"11111111-1111-4111-8111-111111111111"}],"notes":""}`,
			wantStatus: 200,
			wantStored: `{"notes":"","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "the characters of HTML and a line feed are stored as JSON.stringify writes them",
			document:   `{"notes":"a<b>&c\n\"q\"","approaches":[` + documentApproach + `]}`,
			wantStatus: 200,
			wantStored: `{"notes":"a<b>&c\n\"q\"","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "escapes that make the request larger than the limit count as the characters",
			document:   `{"notes":"` + strings.Repeat(`é`, 44_000) + `","approaches":[` + documentApproach + `]}`,
			wantStatus: 200,
			wantStored: `{"notes":"` + strings.Repeat("é", 44_000) + `","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "a stored text of 256000 bytes is accepted",
			document:   `{"notes":"` + strings.Repeat("x", 256_000-documentBase) + `","approaches":[` + documentApproach + `]}`,
			wantStatus: 200,
			wantStored: `{"notes":"` + strings.Repeat("x", 256_000-documentBase) + `","approaches":[` + documentApproach + `]}`,
		},
		{
			name:       "a stored text of 256001 bytes is too large",
			document:   `{"notes":"` + strings.Repeat("x", 256_001-documentBase) + `","approaches":[` + documentApproach + `]}`,
			wantStatus: 413,
			wantCode:   "too_large",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServer(t, pool)
			cookie := loginAs(t, s, "student")
			problemID := addProblem(t, s, cookie, raceProblem)
			draft := request(t, s, "/problems/"+problemID+"/attempts", http.MethodPost, `{}`, cookie)
			if draft.Code != http.StatusCreated {
				t.Fatalf("new draft: status %d, body %s", draft.Code, draft.Body.String())
			}
			attemptID := decode[Attempt](t, draft).ID
			before := queryValue[string](t, pool, "SELECT document FROM attempts WHERE id=$1", attemptID)

			w := request(t, s, "/problems/"+problemID+"/attempts/"+attemptID, http.MethodPut, `{"version":1,"document":`+tt.document+`}`, cookie)
			if w.Code != tt.wantStatus || errorCode(w) != tt.wantCode {
				t.Fatalf("status = %d, code = %q; want %d, %q", w.Code, errorCode(w), tt.wantStatus, tt.wantCode)
			}
			stored := queryValue[string](t, pool, "SELECT document FROM attempts WHERE id=$1", attemptID)
			if tt.wantStatus != 200 {
				if stored != before {
					t.Errorf("a rejected update changed the stored document")
				}
				return
			}
			if stored != tt.wantStored {
				t.Errorf("stored document = %.300s, want %.300s", stored, tt.wantStored)
			}
			if got := string(decode[struct {
				Document json.RawMessage `json:"document"`
			}](t, w).Document); got != tt.wantStored {
				t.Errorf("document of the response = %.300s, want the stored text", got)
			}
			again := request(t, s, "/problems/"+problemID+"/attempts/"+attemptID, http.MethodGet, nil, cookie)
			if got := string(decode[struct {
				Document json.RawMessage `json:"document"`
			}](t, again).Document); got != tt.wantStored {
				t.Errorf("document of a later read = %.300s, want the stored text", got)
			}
		})
	}
}

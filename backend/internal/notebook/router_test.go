package notebook

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const publicOrigin = "https://redlamp.me"

// Behind CloudFront the request host is the host of the function URL.
const functionURL = "https://abc123.lambda-url.eu-central-1.on.aws"

func newServer(t *testing.T) *Server {
	t.Helper()
	s, err := New(nil, publicOrigin)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		origin  string
		wantErr bool
	}{
		{name: "accepts the public origin", origin: "https://redlamp.me"},
		{name: "accepts a local origin with a port", origin: "http://localhost:8787"},
		{name: "rejects an empty origin", origin: "", wantErr: true},
		{name: "rejects a host without a scheme", origin: "redlamp.me", wantErr: true},
		{name: "rejects a slash after the host", origin: "https://redlamp.me/", wantErr: true},
		{name: "rejects a path", origin: "https://redlamp.me/leetcode", wantErr: true},
		{name: "rejects upper case", origin: "https://Redlamp.me", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(nil, tt.origin)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(%q) error = %v, want an error: %v", tt.origin, err, tt.wantErr)
			}
		})
	}
}

// The expected values are the cases of test/routes.test.ts and test/auth.test.ts that need no session,
// and the origin rule of the design: a write must come from PUBLIC_ORIGIN.
func TestRoutesWithoutSession(t *testing.T) {
	const forbidden = `{"error":{"code":"forbidden","message":"Отправляйте изменения со страницы тетради."}}`
	const unauthorized = `{"error":{"code":"unauthorized","message":"Войдите, чтобы продолжить."}}`
	const notAllowed = `{"error":{"code":"method_not_allowed","message":"Этот метод недоступен."}}`
	tests := []struct {
		name         string
		method       string
		path         string
		origin       string
		wantStatus   int
		wantType     string
		wantLocation string
		wantAllow    string
		wantBody     string
		outside      bool
	}{
		{name: "private root", method: "GET", path: "/leetcode", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode"},
		{name: "private slash", method: "GET", path: "/leetcode/", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode"},
		{name: "private HTML alias", method: "GET", path: "/leetcode/index.html", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode"},
		{name: "root HTML alias", method: "GET", path: "/leetcode.html", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode"},
		{name: "private page with query", method: "GET", path: "/leetcode/problems?due=2026-10-07", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode%2Fproblems%3Fdue%3D2026-10-07"},
		{name: "private page slash alias", method: "GET", path: "/leetcode/problems/", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode%2Fproblems"},
		{name: "private unknown page", method: "GET", path: "/leetcode/missing", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode%2Fmissing"},
		{name: "private catalog problem page", method: "GET", path: "/leetcode/sets/two-pointers/problems/merge-sorted-array", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode%2Fsets%2Ftwo-pointers%2Fproblems%2Fmerge-sorted-array"},
		{name: "anonymous tutor page", method: "GET", path: "/leetcode/tutor", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode%2Ftutor"},
		{name: "private root with HEAD", method: "HEAD", path: "/leetcode", wantStatus: 302, wantLocation: "/leetcode/login?return=%2Fleetcode"},
		{name: "anonymous API", method: "GET", path: "/leetcode/api/problems", wantStatus: 401, wantType: "application/json", wantBody: unauthorized},
		{name: "anonymous API with HEAD", method: "HEAD", path: "/leetcode/api/problems", wantStatus: 401, wantType: "application/json"},
		{name: "login shell", method: "GET", path: "/leetcode/login", wantStatus: 200, wantType: "text/html; charset=utf-8"},
		{name: "login shell with HEAD", method: "HEAD", path: "/leetcode/login", wantStatus: 200, wantType: "text/html; charset=utf-8"},
		{name: "cross-origin write", method: "POST", path: "/leetcode/api/session", origin: "https://other.example", wantStatus: 403, wantType: "application/json", wantBody: forbidden},
		{name: "write with the null origin", method: "POST", path: "/leetcode/api/session", origin: "null", wantStatus: 403, wantType: "application/json", wantBody: forbidden},
		{name: "write without an origin", method: "POST", path: "/leetcode/api/session", wantStatus: 403, wantType: "application/json", wantBody: forbidden},
		{name: "write with the origin of the function URL", method: "POST", path: "/leetcode/api/session", origin: functionURL, wantStatus: 403, wantType: "application/json", wantBody: forbidden},
		{name: "cross-origin write to the login page", method: "POST", path: "/leetcode/login", origin: "https://other.example", wantStatus: 403, wantType: "application/json", wantBody: forbidden},
		{name: "write from the public origin passes the origin check", method: "POST", path: "/leetcode/api/problems", origin: publicOrigin, wantStatus: 401, wantType: "application/json", wantBody: unauthorized},
		{name: "write to the login page", method: "POST", path: "/leetcode/login", origin: publicOrigin, wantStatus: 405, wantType: "application/json", wantAllow: "GET, HEAD", wantBody: notAllowed},
		{name: "write to a page", method: "DELETE", path: "/leetcode/problems", origin: publicOrigin, wantStatus: 405, wantType: "application/json", wantAllow: "GET, HEAD", wantBody: notAllowed},
		{name: "asset path", method: "GET", path: "/leetcode/assets/main.js", wantStatus: 404, outside: true},
		{name: "public site path", method: "GET", path: "/", wantStatus: 404, outside: true},
		{name: "path with the notebook name as a prefix", method: "GET", path: "/leetcodex", wantStatus: 404, outside: true},
	}
	s := newServer(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, functionURL+tt.path, nil)
			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}
			w := httptest.NewRecorder()
			s.ServeHTTP(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			h := w.Header()
			if got := h.Get("Location"); got != tt.wantLocation {
				t.Errorf("Location = %q, want %q", got, tt.wantLocation)
			}
			if got := h.Get("Allow"); got != tt.wantAllow {
				t.Errorf("Allow = %q, want %q", got, tt.wantAllow)
			}
			if tt.outside {
				if got := h.Get("X-Robots-Tag"); got != "" {
					t.Errorf("X-Robots-Tag = %q, want no notebook header", got)
				}
				return
			}
			if got := h.Get("Content-Type"); got != tt.wantType {
				t.Errorf("Content-Type = %q, want %q", got, tt.wantType)
			}
			for name, want := range map[string]string{"Cache-Control": "private, no-store", "X-Content-Type-Options": "nosniff", "X-Robots-Tag": "noindex, nofollow", "Referrer-Policy": "no-referrer"} {
				if got := h.Get(name); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
			body := w.Body.String()
			if tt.method == "HEAD" && body != "" {
				t.Errorf("HEAD body = %q, want none", body)
			}
			if tt.wantBody != "" && body != tt.wantBody {
				t.Errorf("body = %s, want %s", body, tt.wantBody)
			}
			if tt.name == "login shell" && (!strings.Contains(body, `<div id="root">`) || strings.Contains(body, `id="session"`)) {
				t.Errorf("login shell body = %s, want the shell without a session", body)
			}
		})
	}
}

// The expected values are the signed-in cases of test/routes.test.ts and the paths of NOTEBOOK_PAGE in shared/leetcode.ts.
func TestRoutesWithSession(t *testing.T) {
	const noPage = "Такой страницы тетради нет."
	const noRequest = "Такого запроса к тетради нет."
	const session = `<script id="session" type="application/json">{"role":"student","username":"student-test"}</script>`
	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantLocation string
		wantMessage  string
	}{
		{name: "signed-in HTML alias with query", path: "/leetcode/problems.html?q=x", wantStatus: 302, wantLocation: "/leetcode/problems?q=x"},
		{name: "signed-in slash alias", path: "/leetcode/", wantStatus: 302, wantLocation: "/leetcode"},
		{name: "signed-in index alias of a page", path: "/leetcode/sets/index.html", wantStatus: 302, wantLocation: "/leetcode/sets"},
		{name: "signed-in unknown page", path: "/leetcode/missing", wantStatus: 404, wantMessage: noPage},
		{name: "signed-in task page", path: "/leetcode/tasks/11111111-1111-4111-8111-111111111111", wantStatus: 404, wantMessage: noPage},
		{name: "signed-in chapter page with number zero", path: "/leetcode/sets/permutations-and-combinations/0", wantStatus: 404, wantMessage: noPage},
		{name: "signed-in chapter page with three digits", path: "/leetcode/sets/permutations-and-combinations/100", wantStatus: 404, wantMessage: noPage},
		{name: "signed-in lessons page", path: "/leetcode/lessons", wantStatus: 404, wantMessage: noPage},
		{name: "signed-in login page with a slash", path: "/leetcode/login/", wantStatus: 404, wantMessage: noPage},
		{name: "signed-in root", path: "/leetcode", wantStatus: 200},
		{name: "signed-in shell", path: "/leetcode/problems", wantStatus: 200},
		{name: "signed-in problem page", path: "/leetcode/problems/11111111-1111-4111-8111-111111111111", wantStatus: 200},
		{name: "signed-in homework review page", path: "/leetcode/homework/11111111-1111-4111-8111-111111111111/review", wantStatus: 200},
		{name: "signed-in sets page", path: "/leetcode/sets", wantStatus: 200},
		{name: "signed-in set page", path: "/leetcode/sets/permutations-and-combinations", wantStatus: 200},
		{name: "signed-in chapter page", path: "/leetcode/sets/permutations-and-combinations/2", wantStatus: 200},
		{name: "signed-in catalog problem page", path: "/leetcode/sets/permutations-and-combinations/problems/subsets", wantStatus: 200},
		{name: "signed-in tutor page", path: "/leetcode/tutor", wantStatus: 200},
		{name: "signed-in quiz page", path: "/leetcode/quiz", wantStatus: 200},
		{name: "signed-in deck page", path: "/leetcode/quiz/permutations-and-combinations", wantStatus: 200},
		{name: "signed-in unknown API", path: "/leetcode/api/missing", wantStatus: 404, wantMessage: noRequest},
		{name: "signed-in lessons API", path: "/leetcode/api/lessons", wantStatus: 404, wantMessage: noRequest},
	}
	s := newServer(t)
	user := &SessionUser{Role: "student", Username: "student-test"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, functionURL+tt.path, nil)
			w := httptest.NewRecorder()
			err := s.serve(w, r, r.URL.Path, user)
			if tt.wantMessage != "" {
				var failure *HTTPError
				if !errors.As(err, &failure) || failure.Status != tt.wantStatus || failure.Code != "not_found" || failure.Message != tt.wantMessage {
					t.Errorf("error = %#v, want %d not_found %q", err, tt.wantStatus, tt.wantMessage)
				}
				return
			}
			if err != nil {
				t.Fatalf("serve: %v", err)
			}
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if got := w.Header().Get("Location"); got != tt.wantLocation {
				t.Errorf("Location = %q, want %q", got, tt.wantLocation)
			}
			if tt.wantStatus == 200 && !strings.Contains(w.Body.String(), session) {
				t.Errorf("body = %s, want the session %s", w.Body.String(), session)
			}
		})
	}
}

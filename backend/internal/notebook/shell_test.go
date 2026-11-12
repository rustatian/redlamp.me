package notebook

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

// The nonce is 16 random bytes in base64: 22 characters and two padding characters.
var cspNonce = regexp.MustCompile(`style-src 'self' 'nonce-([A-Za-z0-9+/]{22}==)';`)

// The expected page and policy are the text of src/leetcode/shell.ts.
func TestRenderShell(t *testing.T) {
	tests := []struct {
		name    string
		user    *SessionUser
		session string
	}{
		{name: "a shell without a user has no session", user: nil, session: ""},
		{name: "a shell with a user has the session JSON", user: &SessionUser{Role: "parent", Username: "parent-test"}, session: `<script id="session" type="application/json">{"role":"parent","username":"parent-test"}</script>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if err := renderShell(w, httptest.NewRequest(http.MethodGet, "/leetcode/login", nil), tt.user); err != nil {
				t.Fatal(err)
			}
			if w.Code != http.StatusOK {
				t.Errorf("status = %d, want 200", w.Code)
			}
			if got := w.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
				t.Errorf("Content-Type = %q", got)
			}
			policy := w.Header().Get("Content-Security-Policy")
			match := cspNonce.FindStringSubmatch(policy)
			if match == nil {
				t.Fatalf("Content-Security-Policy = %q, want a style nonce", policy)
			}
			nonce := match[1]
			if want := "default-src 'self'; script-src 'self'; style-src 'self' 'nonce-" + nonce + "'; connect-src 'self'; media-src 'self' blob:; frame-src 'none'; object-src 'none'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"; policy != want {
				t.Errorf("Content-Security-Policy = %q, want %q", policy, want)
			}
			want := `<!doctype html>
<html lang="ru"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Тетрадь LeetCode | redlamp</title><link rel="icon" href="/leetcode/assets/icon.svg" type="image/svg+xml"><link rel="stylesheet" href="/leetcode/assets/style.css"><meta name="csp-nonce" nonce="` + nonce + `">
` + tt.session + `<script type="module" src="/leetcode/assets/main.js"></script></head>
<body><div id="root"></div><noscript>Включите JavaScript, чтобы пользоваться тетрадью.</noscript></body></html>`
			if got := w.Body.String(); got != want {
				t.Errorf("page =\n%s\nwant\n%s", got, want)
			}
		})
	}
}

func TestEachShellResponseHasANewNonce(t *testing.T) {
	s := newServer(t, nil)
	read := func() string {
		w := httptest.NewRecorder()
		s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, functionURL+"/leetcode/login", nil))
		match := cspNonce.FindStringSubmatch(w.Header().Get("Content-Security-Policy"))
		if match == nil {
			t.Fatalf("Content-Security-Policy = %q, want a style nonce", w.Header().Get("Content-Security-Policy"))
		}
		return match[1]
	}
	if first, second := read(), read(); first == second {
		t.Errorf("two responses have the nonce %q", first)
	}
}

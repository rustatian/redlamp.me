package notebook

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/db"
	"github.com/rustatian/redlamp.me/backend/internal/dbtest"
	"github.com/rustatian/redlamp.me/backend/migrations"
)

// The accounts and the passwords of test/fixtures.ts and test/client.ts.
const (
	studentPassword = "student-test-password-01"
	parentPassword  = "parent-test-password-001"
)

func digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func accountsSecret(parentName, parentDigest, studentName, studentDigest string) string {
	return fmt.Sprintf(`{"parent":{"username":%q,"passwordHash":%q},"student":{"username":%q,"passwordHash":%q}}`, parentName, parentDigest, studentName, studentDigest)
}

var testAccounts = accountsSecret("parent-test", digest(parentPassword), "student-test", digest(studentPassword))

// database returns a pool on a new database with the baseline schema.
func database(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool := dbtest.New(t)
	if _, err := db.Migrate(context.Background(), pool, migrations.FS, false); err != nil {
		t.Fatal(err)
	}
	return pool
}

func serverWith(t *testing.T, pool *pgxpool.Pool, origin, accounts string) *Server {
	t.Helper()
	s, err := New(Config{DB: pool, PublicOrigin: origin, Accounts: accounts, LeetCode: &fakeLeetCode{}})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// newServer returns a Server with the public origin and the test accounts. pool is nil for a test that reads no session.
func newServer(t *testing.T, pool *pgxpool.Pool) *Server {
	t.Helper()
	return serverWith(t, pool, publicOrigin, testAccounts)
}

var viewers atomic.Int64

// call sends a JSON request from the public origin, as the web client does. Each call has a new viewer address.
func call(s *Server, method, path, body, cookie string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, functionURL+path, strings.NewReader(body))
	n := viewers.Add(1)
	r.Header.Set("CloudFront-Viewer-Address", fmt.Sprintf("198.51.%d.%d:46532", n/256%256, n%256))
	r.Header.Set("Origin", publicOrigin)
	r.Header.Set("Content-Type", "application/json")
	if cookie != "" {
		r.Header.Set("Cookie", cookie)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	return w
}

// login signs in and returns the name=value pair of the session cookie.
func login(t *testing.T, s *Server, username, password string) string {
	t.Helper()
	w := call(s, http.MethodPost, "/leetcode/api/session", fmt.Sprintf(`{"username":%q,"password":%q}`, username, password), "")
	if w.Code != http.StatusOK {
		t.Fatalf("login of %s: status %d, body %s", username, w.Code, w.Body.String())
	}
	pair, _, _ := strings.Cut(w.Header().Get("Set-Cookie"), ";")
	return pair
}

const (
	invalidCredentials = `{"error":{"code":"invalid_credentials","message":"Неверное имя пользователя или пароль."}}`
	unauthorized       = `{"error":{"code":"unauthorized","message":"Войдите, чтобы продолжить."}}`
	setupRequired      = `{"error":{"code":"setup_required","message":"Аккаунты тетради не настроены."}}`
	studentSession     = `{"role":"student","username":"student-test"}`
)

// The cases of test/auth.test.ts, with the response bodies of src/leetcode/auth.ts.
func TestSignIn(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		origin      string
		contentType string
		wantStatus  int
		wantBody    string
	}{
		{name: "wrong password", body: `{"username":"student-test","password":"wrong"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 401, wantBody: invalidCredentials},
		{name: "unknown username", body: `{"username":"unknown","password":"wrong"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 401, wantBody: invalidCredentials},
		{name: "password of the other account", body: `{"username":"student-test","password":"parent-test-password-001"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 401, wantBody: invalidCredentials},
		{name: "cross-origin login", body: `{"username":"student-test","password":"student-test-password-01"}`, origin: "https://other.example", contentType: "application/json", wantStatus: 403, wantBody: `{"error":{"code":"forbidden","message":"Отправляйте изменения со страницы тетради."}}`},
		{name: "null origin", body: `{"username":"student-test","password":"student-test-password-01"}`, origin: "null", contentType: "application/json", wantStatus: 403, wantBody: `{"error":{"code":"forbidden","message":"Отправляйте изменения со страницы тетради."}}`},
		{name: "valid student", body: `{"username":"student-test","password":"student-test-password-01"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 200, wantBody: studentSession},
		{name: "valid parent", body: `{"username":"parent-test","password":"parent-test-password-001"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 200, wantBody: `{"role":"parent","username":"parent-test"}`},
		{name: "oversized login request", body: `{"username":"student-test","password":"` + strings.Repeat("x", 4096) + `"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 413, wantBody: `{"error":{"code":"too_large","message":"Файл или запрос слишком большой."}}`},
		{name: "plain text login request", body: `{"username":"student-test","password":"student-test-password-01"}`, origin: publicOrigin, contentType: "text/plain", wantStatus: 400, wantBody: `{"error":{"code":"invalid_input","message":"Отправьте запрос в формате JSON."}}`},
		{name: "empty username", body: `{"username":" ","password":"student-test-password-01"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 400, wantBody: `{"error":{"code":"invalid_input","message":"Проверьте текстовые поля."}}`},
		{name: "password above 256 characters", body: `{"username":"student-test","password":"` + strings.Repeat("x", 257) + `"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 400, wantBody: `{"error":{"code":"invalid_input","message":"Проверьте текстовые поля."}}`},
		{name: "request without a password", body: `{"username":"student-test"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 400, wantBody: `{"error":{"code":"invalid_input","message":"Проверьте текстовые поля."}}`},
		{name: "request with another field", body: `{"username":"student-test","password":"student-test-password-01","role":"parent"}`, origin: publicOrigin, contentType: "application/json", wantStatus: 400, wantBody: `{"error":{"code":"invalid_input","message":"Проверьте поля запроса."}}`},
	}
	sessionCookie := regexp.MustCompile(`^__Host-redlamp-leetcode=([a-f0-9]{64}); Path=/; HttpOnly; SameSite=Lax; Max-Age=604800; Secure$`)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			pool := database(t)
			s := newServer(t, pool)
			r := httptest.NewRequest(http.MethodPost, functionURL+"/leetcode/api/session", strings.NewReader(tt.body))
			r.Header.Set("Origin", tt.origin)
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			before := time.Now().UnixMilli()
			s.ServeHTTP(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if got := w.Body.String(); got != tt.wantBody {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
			var sessions int
			if err := pool.QueryRow(ctx, "SELECT count(*) FROM sessions").Scan(&sessions); err != nil {
				t.Fatal(err)
			}
			if tt.wantStatus != 200 {
				if got := w.Header().Get("Set-Cookie"); got != "" || sessions != 0 {
					t.Errorf("Set-Cookie = %q and %d sessions, want no session", got, sessions)
				}
				return
			}
			match := sessionCookie.FindStringSubmatch(w.Header().Get("Set-Cookie"))
			if match == nil {
				t.Fatalf("Set-Cookie = %q, want the session cookie", w.Header().Get("Set-Cookie"))
			}
			var tokenHash, role string
			var createdAt, expiresAt int64
			if err := pool.QueryRow(ctx, "SELECT token_hash, role, created_at, expires_at FROM sessions").Scan(&tokenHash, &role, &createdAt, &expiresAt); err != nil {
				t.Fatal(err)
			}
			if tokenHash != digest(match[1]) {
				t.Errorf("token_hash = %s, want the SHA-256 of the cookie value", tokenHash)
			}
			if !strings.Contains(tt.wantBody, `"role":"`+role+`"`) {
				t.Errorf("role = %s, want the role of %s", role, tt.wantBody)
			}
			if createdAt < before || createdAt > time.Now().UnixMilli() || expiresAt-createdAt != 604_800_000 {
				t.Errorf("created_at = %d, expires_at = %d; want the current time in milliseconds and 7 days", createdAt, expiresAt)
			}
		})
	}
}

// The cases "blocks private data" of test/auth.test.ts, and the rule of the task: a changed username or digest ends the sessions of that account.
func TestSessionRead(t *testing.T) {
	zeros := strings.Repeat("0", 64)
	tests := []struct {
		name       string
		signIn     bool
		change     string
		accounts   string
		wantStatus int
		wantBody   string
	}{
		{name: "signed-in student", signIn: true, accounts: testAccounts, wantStatus: 200, wantBody: studentSession},
		{name: "absent token", signIn: false, accounts: testAccounts, wantStatus: 401, wantBody: unauthorized},
		{name: "expired token", signIn: true, change: "UPDATE sessions SET expires_at = 0", accounts: testAccounts, wantStatus: 401, wantBody: unauthorized},
		{name: "deleted session", signIn: true, change: "DELETE FROM sessions", accounts: testAccounts, wantStatus: 401, wantBody: unauthorized},
		{name: "rotated credentials", signIn: true, accounts: accountsSecret("parent-test", digest(parentPassword), "student-test", zeros), wantStatus: 401, wantBody: unauthorized},
		{name: "changed username", signIn: true, accounts: accountsSecret("parent-test", digest(parentPassword), "student-new", digest(studentPassword)), wantStatus: 401, wantBody: unauthorized},
		{name: "rotated credentials of the other account", signIn: true, accounts: accountsSecret("parent-new", zeros, "student-test", digest(studentPassword)), wantStatus: 200, wantBody: studentSession},
		{name: "logout", signIn: true, change: "logout", accounts: testAccounts, wantStatus: 401, wantBody: unauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServer(t, pool)
			cookie := ""
			if tt.signIn {
				cookie = login(t, s, "student-test", studentPassword)
			}
			switch tt.change {
			case "":
			case "logout":
				if w := call(s, http.MethodDelete, "/leetcode/api/session", "", cookie); w.Code != http.StatusNoContent {
					t.Fatalf("logout status = %d, want 204", w.Code)
				}
			default:
				if _, err := pool.Exec(context.Background(), tt.change); err != nil {
					t.Fatal(err)
				}
			}
			w := call(serverWith(t, pool, publicOrigin, tt.accounts), http.MethodGet, "/leetcode/api/session", "", cookie)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if got := w.Body.String(); got != tt.wantBody {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
		})
	}
}

// The cookie rules of token() in src/leetcode/auth.ts: the first cookie with the name counts, and its value is 64 lower-case hex digits.
func TestSessionCookieForms(t *testing.T) {
	tests := []struct {
		name       string
		header     func(name, value string) string
		wantStatus int
	}{
		{name: "the session cookie alone", header: func(name, value string) string { return name + "=" + value }, wantStatus: 200},
		{name: "the session cookie between other cookies", header: func(name, value string) string { return "theme=dark;  " + name + "=" + value + "; other=1" }, wantStatus: 200},
		{name: "a value in upper case", header: func(name, value string) string { return name + "=" + strings.ToUpper(value) }, wantStatus: 401},
		{name: "a value of 63 digits", header: func(name, value string) string { return name + "=" + value[1:] }, wantStatus: 401},
		{name: "a first cookie with the name and a bad value", header: func(name, value string) string { return name + "=bad; " + name + "=" + value }, wantStatus: 401},
		{name: "the value under the local cookie name", header: func(_, value string) string { return "redlamp-leetcode-dev=" + value }, wantStatus: 401},
	}
	pool := database(t)
	s := newServer(t, pool)
	name, value, _ := strings.Cut(login(t, s, "student-test", studentPassword), "=")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if w := call(s, http.MethodGet, "/leetcode/api/session", "", tt.header(name, value)); w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

// The cases "fails closed" of test/auth.test.ts, and the checks of accounts() in src/leetcode/auth.ts.
func TestAccountsFailClosed(t *testing.T) {
	good := digest(studentPassword)
	tests := []struct {
		name       string
		secret     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "missing secret", secret: "", path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "malformed secret", secret: "{", path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "missing role", secret: `{"parent":{}}`, path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "secret that is an array", secret: `[]`, path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "a third role", secret: `{"parent":{"username":"parent-test","passwordHash":"` + good + `"},"student":{"username":"student-test","passwordHash":"` + good + `"},"tutor":{}}`, path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "another field in an account", secret: `{"parent":{"username":"parent-test","passwordHash":"` + good + `","password":"x"},"student":{"username":"student-test","passwordHash":"` + good + `"}}`, path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "one username for both accounts", secret: accountsSecret("same-name", good, "same-name", good), path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "username of two characters", secret: accountsSecret("pa", good, "student-test", good), path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "username with a space", secret: accountsSecret("parent test", good, "student-test", good), path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "digest in upper case", secret: accountsSecret("parent-test", strings.ToUpper(good), "student-test", good), path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "digest of 63 characters", secret: accountsSecret("parent-test", good[1:], "student-test", good), path: "/leetcode/api/session", wantStatus: 503, wantBody: setupRequired},
		{name: "page request without accounts", secret: "", path: "/leetcode/problems", wantStatus: 503, wantBody: setupRequired},
		{name: "login page without accounts", secret: "", path: "/leetcode/login", wantStatus: 200},
		{name: "valid accounts", secret: testAccounts, path: "/leetcode/api/session", wantStatus: 401, wantBody: unauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := call(serverWith(t, nil, publicOrigin, tt.secret), http.MethodGet, tt.path, "", "")
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if got := w.Body.String(); tt.wantBody != "" && got != tt.wantBody {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
			if got := w.Header().Get("Cache-Control"); got != "private, no-store" {
				t.Errorf("Cache-Control = %q, want private, no-store", got)
			}
		})
	}
}

func TestSignInWithoutAccountsFailsClosed(t *testing.T) {
	w := call(serverWith(t, nil, publicOrigin, ""), http.MethodPost, "/leetcode/api/session", `{"username":"student-test","password":"student-test-password-01"}`, "")
	if w.Code != 503 || w.Body.String() != setupRequired {
		t.Errorf("status = %d, body = %s; want 503 %s", w.Code, w.Body.String(), setupRequired)
	}
}

// The rule of the login limit: 10 attempts in 60 seconds for each source address.
func TestLimiter(t *testing.T) {
	start := time.Unix(1_794_787_200, 0)
	tests := []struct {
		name     string
		previous int
		key      string
		after    time.Duration
		want     bool
	}{
		{name: "the first attempt passes", previous: 0, key: "a", after: 0, want: true},
		{name: "the tenth attempt in the period passes", previous: 9, key: "a", after: 59 * time.Second, want: true},
		{name: "the eleventh attempt in the period is limited", previous: 10, key: "a", after: 59 * time.Second, want: false},
		{name: "an attempt after the period passes", previous: 10, key: "a", after: 60 * time.Second, want: true},
		{name: "another source has its own count", previous: 10, key: "b", after: 0, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &limiter{}
			for range tt.previous {
				l.allow("a", start)
			}
			if got := l.allow(tt.key, start.Add(tt.after)); got != tt.want {
				t.Errorf("allow = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimiterDropsOldSources(t *testing.T) {
	start := time.Unix(1_794_787_200, 0)
	l := &limiter{}
	for i := range 100 {
		l.allow(fmt.Sprint("source-", i), start)
	}
	l.allow("late", start.Add(59*time.Second))
	l.allow("new", start.Add(60*time.Second))
	if got := len(l.windows); got != 2 {
		t.Errorf("the limiter holds %d sources, want the 2 sources of the last 60 seconds", got)
	}
}

// The case "limits the eleventh attempt from one source" of test/auth.test.ts. The source is the viewer address without its port.
func TestLoginLimit(t *testing.T) {
	const limited = `{"error":{"code":"rate_limited","message":"Слишком много попыток входа. Попробуйте снова через минуту."}}`
	tests := []struct {
		name    string
		address func(attempt int) string
		other   string
	}{
		{name: "IPv4 viewer with a new port for each attempt", address: func(attempt int) string { return fmt.Sprintf("192.0.2.10:%d", 40000+attempt) }, other: "192.0.2.11:40000"},
		{name: "viewer address with colons and a new port for each attempt", address: func(attempt int) string { return fmt.Sprintf("2001:db8::10:%d", 40000+attempt) }, other: "2001:db8::11:40000"},
		{name: "request without the viewer header", address: func(int) string { return "" }, other: "192.0.2.12:40000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newServer(t, nil)
			attempt := func(address string) *httptest.ResponseRecorder {
				r := httptest.NewRequest(http.MethodPost, functionURL+"/leetcode/api/session", strings.NewReader(`{"username":"unknown","password":"wrong"}`))
				r.Header.Set("Origin", publicOrigin)
				r.Header.Set("Content-Type", "application/json")
				if address != "" {
					r.Header.Set("CloudFront-Viewer-Address", address)
				}
				w := httptest.NewRecorder()
				s.ServeHTTP(w, r)
				return w
			}
			for i := range 10 {
				if w := attempt(tt.address(i)); w.Code != 401 {
					t.Fatalf("attempt %d: status %d, want 401", i+1, w.Code)
				}
			}
			if w := attempt(tt.address(10)); w.Code != 429 || w.Body.String() != limited {
				t.Errorf("attempt 11: status %d, body %s; want 429 %s", w.Code, w.Body.String(), limited)
			}
			if w := attempt(tt.other); w.Code != 401 {
				t.Errorf("attempt from another source: status %d, want 401", w.Code)
			}
		})
	}
}

// The cookie rule of src/leetcode/auth.ts: a local HTTP origin gets the cookie redlamp-leetcode-dev without Secure.
func TestSessionCookieOfTheOrigin(t *testing.T) {
	const public = `^__Host-redlamp-leetcode=[a-f0-9]{64}; Path=/; HttpOnly; SameSite=Lax; Max-Age=604800; Secure$`
	const local = `^redlamp-leetcode-dev=[a-f0-9]{64}; Path=/; HttpOnly; SameSite=Lax; Max-Age=604800$`
	tests := []struct {
		name   string
		origin string
		want   string
	}{
		{name: "the public origin", origin: "https://redlamp.me", want: public},
		{name: "localhost over HTTP", origin: "http://localhost:8787", want: local},
		{name: "127.0.0.1 over HTTP", origin: "http://127.0.0.1:8787", want: local},
		{name: "localhost over HTTPS", origin: "https://localhost:8787", want: public},
		{name: "another host over HTTP", origin: "http://notebook.test", want: public},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serverWith(t, database(t), tt.origin, testAccounts)
			r := httptest.NewRequest(http.MethodPost, functionURL+"/leetcode/api/session", strings.NewReader(`{"username":"student-test","password":"student-test-password-01"}`))
			r.Header.Set("Origin", tt.origin)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.ServeHTTP(w, r)
			cookie := w.Header().Get("Set-Cookie")
			if !regexp.MustCompile(tt.want).MatchString(cookie) {
				t.Fatalf("Set-Cookie = %q, want %s", cookie, tt.want)
			}
			pair, _, _ := strings.Cut(cookie, ";")
			if w := call(s, http.MethodGet, "/leetcode/api/session", "", pair); w.Code != 200 || w.Body.String() != studentSession {
				t.Errorf("session read with the cookie: status %d, body %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name         string
		withCookie   bool
		wantSessions int
	}{
		{name: "logout ends the session of the cookie", withCookie: true, wantSessions: 1},
		{name: "logout without a cookie ends no session", withCookie: false, wantSessions: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := newServer(t, pool)
			cookie := login(t, s, "student-test", studentPassword)
			parent := login(t, s, "parent-test", parentPassword)
			if !tt.withCookie {
				cookie = ""
			}
			w := call(s, http.MethodDelete, "/leetcode/api/session", "", cookie)
			if w.Code != http.StatusNoContent || w.Body.Len() != 0 {
				t.Errorf("status = %d, body = %q; want 204 without a body", w.Code, w.Body.String())
			}
			if got, want := w.Header().Get("Set-Cookie"), "__Host-redlamp-leetcode=; Path=/; HttpOnly; SameSite=Lax; Max-Age=0; Secure"; got != want {
				t.Errorf("Set-Cookie = %q, want %q", got, want)
			}
			var sessions int
			if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM sessions").Scan(&sessions); err != nil {
				t.Fatal(err)
			}
			if sessions != tt.wantSessions {
				t.Errorf("sessions = %d, want %d", sessions, tt.wantSessions)
			}
			if w := call(s, http.MethodGet, "/leetcode/api/session", "", parent); w.Code != 200 {
				t.Errorf("the session of the parent: status %d, want 200", w.Code)
			}
		})
	}
}

func TestSessionMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{name: "PUT is not a session method", method: http.MethodPut},
		{name: "HEAD is not a session method", method: http.MethodHead},
	}
	s := newServer(t, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := call(s, tt.method, "/leetcode/api/session", "", "")
			if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "GET, POST, DELETE" {
				t.Errorf("status = %d, Allow = %q; want 405 and GET, POST, DELETE", w.Code, w.Header().Get("Allow"))
			}
		})
	}
}

// A sign-in deletes the sessions that expired, as the DELETE of the batch in src/leetcode/auth.ts.
func TestSignInDeletesExpiredSessions(t *testing.T) {
	ctx := context.Background()
	pool := database(t)
	s := newServer(t, pool)
	now := time.Now().UnixMilli()
	if _, err := pool.Exec(ctx, "INSERT INTO sessions VALUES ('expired', 'student', 'f', 0, $1), ('open', 'parent', 'f', 0, $2)", now-1, now+60_000); err != nil {
		t.Fatal(err)
	}
	login(t, s, "student-test", studentPassword)
	var hashes []string
	if err := pool.QueryRow(ctx, "SELECT array_agg(token_hash ORDER BY token_hash) FROM sessions WHERE role = 'parent' OR token_hash = 'expired'").Scan(&hashes); err != nil {
		t.Fatal(err)
	}
	if len(hashes) != 1 || hashes[0] != "open" {
		t.Errorf("old sessions = %v, want only the open session", hashes)
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		role    string
		wantErr bool
	}{
		{name: "the parent does a parent action", user: "parent", role: "parent"},
		{name: "the student does a student action", user: "student", role: "student"},
		{name: "the student does a parent action", user: "student", role: "parent", wantErr: true},
		{name: "the parent does a student action", user: "parent", role: "student", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requireRole(&SessionUser{Role: tt.user, Username: tt.user + "-test"}, tt.role)
			if !tt.wantErr {
				if err != nil {
					t.Errorf("requireRole = %v, want no error", err)
				}
				return
			}
			w := httptest.NewRecorder()
			writeError(w, httptest.NewRequest(http.MethodPost, "/leetcode/api/homework", nil), err)
			if want := `{"error":{"code":"forbidden","message":"Это действие недоступно для этого аккаунта."}}`; w.Code != 403 || w.Body.String() != want {
				t.Errorf("status = %d, body = %s; want 403 %s", w.Code, w.Body.String(), want)
			}
		})
	}
}

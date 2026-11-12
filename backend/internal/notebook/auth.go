package notebook

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/rustatian/redlamp.me/backend/internal/db"
)

// account is one entry of the LEETCODE_ACCOUNTS secret. passwordHash is the SHA-256 (hex) of the password.
type account struct {
	username     string
	passwordHash string
}

var roles = []string{"parent", "student"}

var (
	usernamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{3,40}$`)
	digestPattern   = regexp.MustCompile(`^[a-f0-9]{64}$`)
)

// parseAccounts reads the secret: a JSON object with the accounts of the parent and the student.
func parseAccounts(secret string) (map[string]account, error) {
	var value any
	if err := json.Unmarshal([]byte(secret), &value); err != nil {
		return nil, err
	}
	fields, err := object(value, roles...)
	if err != nil {
		return nil, err
	}
	accounts := map[string]account{}
	for _, role := range roles {
		entry, err := object(fields[role], "username", "passwordHash")
		if err != nil {
			return nil, err
		}
		username, _ := entry["username"].(string)
		passwordHash, _ := entry["passwordHash"].(string)
		if !usernamePattern.MatchString(username) || !digestPattern.MatchString(passwordHash) {
			return nil, fmt.Errorf("the account of the %s is not valid", role)
		}
		accounts[role] = account{username: username, passwordHash: passwordHash}
	}
	if accounts["parent"].username == accounts["student"].username {
		return nil, errors.New("the two accounts have one username")
	}
	return accounts, nil
}

// configured returns the accounts, or the 503 error when the secret is not valid.
func (s *Server) configured() (map[string]account, error) {
	if s.accounts == nil {
		return nil, &HTTPError{Status: http.StatusServiceUnavailable, Code: "setup_required", Message: "Аккаунты тетради не настроены."}
	}
	return s.accounts, nil
}

func hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

// fingerprint ties a session to the username and the password digest of its account.
func fingerprint(role string, a account) string {
	// A list of strings always encodes.
	data, _ := json.Marshal([]string{role, a.username, a.passwordHash})
	return hash(string(data))
}

const sessionSeconds = 604_800

var tokenPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

// token returns the session token of the request, or "".
func (s *Server) token(r *http.Request) string {
	for part := range strings.SplitSeq(r.Header.Get("Cookie"), ";") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(part), s.cookieName+"="); ok {
			if tokenPattern.MatchString(value) {
				return value
			}
			return ""
		}
	}
	return ""
}

func (s *Server) cookie(value string, seconds int) string {
	cookie := fmt.Sprintf("%s=%s; Path=/; HttpOnly; SameSite=Lax; Max-Age=%d", s.cookieName, value, seconds)
	if s.cookieName == localCookie {
		return cookie
	}
	return cookie + "; Secure"
}

// readSession returns the user of the session cookie, or nil for a request without a valid session.
// A session ends when the username or the password digest of its account changes.
func (s *Server) readSession(r *http.Request) (*SessionUser, error) {
	accounts, err := s.configured()
	if err != nil {
		return nil, err
	}
	value := s.token(r)
	if value == "" {
		return nil, nil
	}
	var role, stored string
	err = s.db.QueryRow(r.Context(), "SELECT role, credential_fingerprint FROM sessions WHERE token_hash = $1 AND expires_at > $2", hash(value), time.Now().UnixMilli()).Scan(&role, &stored)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if stored != fingerprint(role, accounts[role]) {
		return nil, nil
	}
	return &SessionUser{Role: role, Username: accounts[role].username}, nil
}

func requireRole(user *SessionUser, role string) error {
	if user.Role != role {
		return &HTTPError{Status: http.StatusForbidden, Code: "forbidden", Message: "Это действие недоступно для этого аккаунта."}
	}
	return nil
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) error {
	if err := methods(r, http.MethodGet, http.MethodPost, http.MethodDelete); err != nil {
		return err
	}
	if r.Method == http.MethodPost {
		return s.signIn(w, r)
	}
	user, err := s.readSession(r)
	if err != nil {
		return err
	}
	if r.Method == http.MethodDelete {
		if value := s.token(r); value != "" {
			if _, err := s.db.Exec(r.Context(), "DELETE FROM sessions WHERE token_hash = $1", hash(value)); err != nil {
				return err
			}
		}
		w.Header().Set("Set-Cookie", s.cookie("", 0))
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	if user == nil {
		return &HTTPError{Status: http.StatusUnauthorized, Code: "unauthorized", Message: "Войдите, чтобы продолжить."}
	}
	return writeJSON(w, r, http.StatusOK, user)
}

func (s *Server) signIn(w http.ResponseWriter, r *http.Request) error {
	if !s.logins.allow(source(r), time.Now()) {
		return &HTTPError{Status: http.StatusTooManyRequests, Code: "rate_limited", Message: "Слишком много попыток входа. Попробуйте снова через минуту."}
	}
	body, err := readJSON(r, 4096)
	if err != nil {
		return err
	}
	data, err := object(body, "username", "password")
	if err != nil {
		return err
	}
	username, err := text(data["username"], 40, true)
	if err != nil {
		return err
	}
	password, err := text(data["password"], 256, false)
	if err != nil {
		return err
	}
	accounts, err := s.configured()
	if err != nil {
		return err
	}
	role, expected := "", strings.Repeat("0", 64)
	for _, candidate := range roles {
		if accounts[candidate].username == username {
			role, expected = candidate, accounts[candidate].passwordHash
		}
	}
	// The comparison takes the same time for each digest, also for an unknown username.
	matches := subtle.ConstantTimeCompare([]byte(hash(password)), []byte(expected)) == 1
	if role == "" || !matches {
		return &HTTPError{Status: http.StatusUnauthorized, Code: "invalid_credentials", Message: "Неверное имя пользователя или пароль."}
	}
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return err
	}
	value := hex.EncodeToString(random)
	now := time.Now().UnixMilli()
	err = db.Tx(r.Context(), s.db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(r.Context(), "DELETE FROM sessions WHERE expires_at <= $1", now); err != nil {
			return err
		}
		_, err := tx.Exec(r.Context(), "INSERT INTO sessions VALUES ($1, $2, $3, $4, $5)", hash(value), role, fingerprint(role, accounts[role]), now, now+sessionSeconds*1000)
		return err
	})
	if err != nil {
		return err
	}
	w.Header().Set("Set-Cookie", s.cookie(value, sessionSeconds))
	return writeJSON(w, r, http.StatusOK, SessionUser{Role: role, Username: username})
}

// source returns the address of the viewer without the source port. CloudFront sends the header in the form 198.51.100.10:46532.
// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/adding-cloudfront-headers.html
func source(r *http.Request) string {
	address := r.Header.Get("CloudFront-Viewer-Address")
	if address == "" {
		return "local"
	}
	if i := strings.LastIndexByte(address, ':'); i >= 0 {
		return address[:i]
	}
	return address
}

const (
	loginAttempts = 10
	loginPeriod   = 60
)

// limiter counts the sign-in attempts of each source address in periods of 60 seconds.
// shortcut: the count is in the memory of one function instance, so the limit applies for each instance. Move the count to a database table if the function cannot reserve concurrency.
type limiter struct {
	mu      sync.Mutex
	windows map[string]window
	swept   int64
}

// window is the count of one source. start is in seconds since the Unix epoch.
type window struct {
	start int64
	count int
}

// allow counts one attempt of key at the time now. It returns false above 10 attempts in the period.
func (l *limiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	seconds := now.Unix()
	if l.windows == nil {
		l.windows = map[string]window{}
	}
	if seconds-l.swept >= loginPeriod {
		for source, w := range l.windows {
			if seconds-w.start >= loginPeriod {
				delete(l.windows, source)
			}
		}
		l.swept = seconds
	}
	w, ok := l.windows[key]
	if !ok || seconds-w.start >= loginPeriod {
		w = window{start: seconds}
	}
	w.count++
	l.windows[key] = w
	return w.count <= loginAttempts
}

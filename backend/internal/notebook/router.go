// Package notebook serves the learning notebook: the pages and the API under /leetcode.
package notebook

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// Server is the http.Handler of the notebook paths.
type Server struct {
	db         *pgxpool.Pool
	origin     string
	cookieName string
	// accounts is nil when the secret is not valid. Each request that needs an account then gets the 503 setup_required error.
	accounts map[string]account
	logins   limiter
	leetcode LeetCode
	sets     []sets.LearningSet
	// decks are the quiz decks of the learning sets, one for each chapter.
	decks      []deck
	references map[string]sets.Reference
	generators map[string]sets.Generator
	// now is the clock of the rules that depend on the current day.
	now func() time.Time
	// random returns a random number from 0 to count-1. It chooses the recall days.
	random func(count int) int
	// tutor is nil when the server has no API key. A question to the tutor then gets the 503 setup_required error.
	tutor *anthropic.Client
	// thinkingEvery is the longest time without an event in the answer stream of the tutor.
	thinkingEvery time.Duration
	// transcriber is nil when the server has no speech service. A dictation then gets the 503 dictation_unavailable error.
	transcriber Transcriber
	// assets answers the paths outside the notebook and the notebook assets. It is nil behind CloudFront, which sends these paths to S3.
	assets http.Handler
}

// Config is the setup of a Server.
type Config struct {
	DB *pgxpool.Pool
	// PublicOrigin is the Origin header that the browser sends from the site, for example https://redlamp.me.
	// Behind CloudFront the request host is the host of the function URL, so the origin check of a write uses this value.
	PublicOrigin string
	// Accounts is the LEETCODE_ACCOUNTS secret: the JSON of the two accounts.
	Accounts string
	// LeetCode gives the data of the problems of leetcode.com.
	LeetCode LeetCode
	// Sets is the content of the learning sets.
	Sets []sets.LearningSet
	// References are the reference solutions of the learning set tasks, by problem slug.
	References map[string]sets.Reference
	// Generators make the hidden tests of the learning set tasks, by problem slug.
	Generators map[string]sets.Generator
	// Anthropic is the client of the Anthropic API for the tutor, or nil when the server has no API key.
	Anthropic *anthropic.Client
	// Transcriber is the speech service of the dictation, or nil when the server has none.
	Transcriber Transcriber
	// Assets serves the built site for a local server: the public pages and the notebook assets. With nil these paths get 404.
	Assets http.Handler
}

const (
	siteCookie = "__Host-redlamp-leetcode"
	// localCookie is the session cookie of a local HTTP origin. A browser refuses a __Host- cookie without Secure.
	localCookie = "redlamp-leetcode-dev"
)

// New returns a Server. It fails when the public origin does not have the form https://host.
func New(cfg Config) (*Server, error) {
	u, err := url.Parse(cfg.PublicOrigin)
	if err != nil || u.Scheme == "" || u.Host == "" || cfg.PublicOrigin != strings.ToLower(u.Scheme+"://"+u.Host) {
		return nil, fmt.Errorf("the public origin %q does not have the form https://host", cfg.PublicOrigin)
	}
	s := &Server{db: cfg.DB, origin: cfg.PublicOrigin, cookieName: siteCookie, leetcode: cfg.LeetCode, sets: cfg.Sets, decks: decksOf(cfg.Sets), references: cfg.References, generators: cfg.Generators, now: time.Now, random: rand.IntN, tutor: cfg.Anthropic, thinkingEvery: 10 * time.Second, transcriber: cfg.Transcriber, assets: cfg.Assets}
	if u.Scheme == "http" && (u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1") {
		s.cookieName = localCookie
	}
	if s.accounts, err = parseAccounts(cfg.Accounts); err != nil {
		log.Printf("the accounts secret is not valid, so sign-in is not available: %v", err)
	}
	return s, nil
}

func isNotebookPath(path string) bool {
	return path == "/leetcode" || path == "/leetcode.html" || strings.HasPrefix(path, "/leetcode/")
}

// notebookPage matches the paths of the web client pages. It equals NOTEBOOK_PAGE of shared/leetcode.ts.
var notebookPage = regexp.MustCompile(`^/leetcode(?:/(?:problems|homework)(?:/[a-f0-9-]{36}(?:/review)?)?|/sets(?:/[a-z0-9-]+(?:/(?:[1-9][0-9]?|problems/[a-z0-9-]+))?)?|/quiz(?:/[a-z0-9-]+)?)?$`)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// CloudFront sends the asset paths and the public site to S3.
	if !isNotebookPath(path) || strings.HasPrefix(path, "/leetcode/assets/") {
		if s.assets != nil {
			s.assets.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}
	h := w.Header()
	h.Set("Cache-Control", "private, no-store")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Robots-Tag", "noindex, nofollow")
	h.Set("Referrer-Policy", "no-referrer")
	if err := s.route(w, r, path); err != nil {
		writeError(w, r, err)
	}
}

func (s *Server) route(w http.ResponseWriter, r *http.Request, path string) error {
	read := r.Method == http.MethodGet || r.Method == http.MethodHead
	if !read && r.Header.Get("Origin") != s.origin {
		return &HTTPError{Status: http.StatusForbidden, Code: "forbidden", Message: "Отправляйте изменения со страницы тетради."}
	}
	if path == "/leetcode/api/session" {
		return s.handleSession(w, r)
	}
	if path == "/leetcode/login" {
		if err := methods(r, http.MethodGet, http.MethodHead); err != nil {
			return err
		}
		return renderShell(w, r, nil)
	}
	user, err := s.readSession(r)
	if err != nil {
		return err
	}
	return s.serve(w, r, path, user)
}

// serve answers an API request or a page request of user. user is nil for a request without a session.
func (s *Server) serve(w http.ResponseWriter, r *http.Request, path string, user *SessionUser) error {
	if strings.HasPrefix(path, "/leetcode/api/") {
		if user == nil {
			return &HTTPError{Status: http.StatusUnauthorized, Code: "unauthorized", Message: "Войдите, чтобы продолжить."}
		}
		// The first handler that knows the path answers the request.
		for _, handle := range []func(http.ResponseWriter, *http.Request, *SessionUser) (bool, error){s.handleProblems, s.handleAttempts, s.handleHomework, s.handleReviews, s.handleSets, s.handleQuiz, s.handleJudge, s.handleTests, s.handleStatement, s.handleRuns, s.handleTutor, s.handleDictation} {
			if handled, err := handle(w, r, user); handled || err != nil {
				return err
			}
		}
		return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такого запроса к тетради нет."}
	}
	if err := methods(r, http.MethodGet, http.MethodHead); err != nil {
		return err
	}
	canonical := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(path, ".html"), "/index"), "/")
	search := ""
	if r.URL.RawQuery != "" {
		search = "?" + r.URL.RawQuery
	}
	if user == nil {
		return redirect(w, "/leetcode/login?return="+url.QueryEscape(canonical+search))
	}
	if !notebookPage.MatchString(canonical) {
		return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такой страницы тетради нет."}
	}
	if canonical != path {
		return redirect(w, canonical+search)
	}
	return renderShell(w, r, user)
}

func redirect(w http.ResponseWriter, location string) error {
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
	return nil
}

// Package notebook serves the learning notebook: the pages and the API under /leetcode.
package notebook

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Server is the http.Handler of the notebook paths.
type Server struct {
	db     *pgxpool.Pool
	origin string
}

// New returns a Server. publicOrigin is the Origin header that the browser sends from the site, for example https://redlamp.me.
// Behind CloudFront the request host is the host of the function URL, so the origin check of a write uses this value.
func New(pool *pgxpool.Pool, publicOrigin string) (*Server, error) {
	u, err := url.Parse(publicOrigin)
	if err != nil || u.Scheme == "" || u.Host == "" || publicOrigin != strings.ToLower(u.Scheme+"://"+u.Host) {
		return nil, fmt.Errorf("the public origin %q does not have the form https://host", publicOrigin)
	}
	return &Server{db: pool, origin: publicOrigin}, nil
}

func isNotebookPath(path string) bool {
	return path == "/leetcode" || path == "/leetcode.html" || strings.HasPrefix(path, "/leetcode/")
}

// notebookPage matches the paths of the web client pages. It equals NOTEBOOK_PAGE of shared/leetcode.ts.
var notebookPage = regexp.MustCompile(`^/leetcode(?:/(?:problems|homework)(?:/[a-f0-9-]{36}(?:/review)?)?|/sets(?:/[a-z0-9-]+(?:/(?:[1-9][0-9]?|problems/[a-z0-9-]+))?)?|/quiz(?:/[a-z0-9-]+)?|/tutor)?$`)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// CloudFront sends the asset paths and the public site to S3.
	if !isNotebookPath(path) || strings.HasPrefix(path, "/leetcode/assets/") {
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
	if path == "/leetcode/login" {
		if err := methods(r, http.MethodGet, http.MethodHead); err != nil {
			return err
		}
		return renderShell(w, r, nil)
	}
	// The server has no sign-in, so each request is anonymous.
	return s.serve(w, r, path, nil)
}

// serve answers an API request or a page request of user. user is nil for a request without a session.
func (s *Server) serve(w http.ResponseWriter, r *http.Request, path string, user *SessionUser) error {
	if strings.HasPrefix(path, "/leetcode/api/") {
		if user == nil {
			return &HTTPError{Status: http.StatusUnauthorized, Code: "unauthorized", Message: "Войдите, чтобы продолжить."}
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

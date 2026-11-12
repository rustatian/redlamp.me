package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// The public addresses of the site keep working: a page has an address with and without the suffix .html.
func TestSite(t *testing.T) {
	dir := t.TempDir()
	for name, text := range map[string]string{
		"index.html":                        "home",
		"examples.html":                     "examples",
		"mathtools/primefactorization.html": "factors",
		"leetcode/assets/main.js":           "script",
		"leetcode/assets/pyodide/lock.json": "lock",
	} {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "home page", path: "/", wantStatus: 200, wantBody: "home"},
		{name: "page without the suffix", path: "/examples", wantStatus: 200, wantBody: "examples"},
		{name: "page with the suffix", path: "/examples.html", wantStatus: 200, wantBody: "examples"},
		{name: "page in a folder without the suffix", path: "/mathtools/primefactorization", wantStatus: 200, wantBody: "factors"},
		{name: "notebook asset", path: "/leetcode/assets/main.js", wantStatus: 200, wantBody: "script"},
		{name: "asset with a dot in a folder", path: "/leetcode/assets/pyodide/lock.json", wantStatus: 200, wantBody: "lock"},
		{name: "unknown page", path: "/missing", wantStatus: 404, wantBody: "404 page not found\n"},
	}
	handler := site(dir)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if w.Code != tt.wantStatus || w.Body.String() != tt.wantBody {
				t.Errorf("status %d with the body %q, want %d with %q", w.Code, w.Body.String(), tt.wantStatus, tt.wantBody)
			}
		})
	}
}

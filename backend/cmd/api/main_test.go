package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
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

// A secret comes from its variable. The variable wins over the Parameter Store name, so a local server needs no AWS access.
// Without the variable, the secret is the decrypted SecureString parameter. A parameter that the function cannot read gives an empty secret.
func TestSecret(t *testing.T) {
	// Test credentials and empty AWS files keep the profiles of the developer and the instance metadata service out of the test.
	t.Setenv("AWS_ACCESS_KEY_ID", "test-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_CONFIG_FILE", os.DevNull)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", os.DevNull)
	t.Setenv("AWS_PROFILE", "")
	tests := []struct {
		name      string
		value     string
		parameter string
		status    int
		response  string
		want      string
		// requests is the number of requests to Parameter Store.
		requests int32
	}{
		{name: "value of the variable", value: "local-secret", want: "local-secret"},
		{name: "variable and parameter name", value: "local-secret", parameter: "/redlamp/test", want: "local-secret"},
		{name: "value of the parameter", parameter: "/redlamp/test", status: http.StatusOK, response: `{"Parameter":{"Value":"stored-secret"}}`, want: "stored-secret", requests: 1},
		{name: "parameter that does not exist", parameter: "/redlamp/missing", status: http.StatusBadRequest, response: `{"__type":"ParameterNotFound"}`, want: "", requests: 1},
		{name: "role without access", parameter: "/redlamp/test", status: http.StatusBadRequest, response: `{"__type":"AccessDeniedException"}`, want: "", requests: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The fake of Parameter Store answers GetParameter of the AWS JSON 1.1 protocol.
			var requests atomic.Int32
			fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				var input struct {
					Name           string
					WithDecryption bool
				}
				err := json.NewDecoder(r.Body).Decode(&input)
				if target := r.Header.Get("X-Amz-Target"); err != nil || target != "AmazonSSM.GetParameter" || input.Name != tt.parameter || !input.WithDecryption {
					t.Errorf("request %s %+v (%v), want AmazonSSM.GetParameter of %q with decryption", target, input, err, tt.parameter)
				}
				w.Header().Set("Content-Type", "application/x-amz-json-1.1")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.response))
			}))
			t.Cleanup(fake.Close)
			t.Setenv("AWS_ENDPOINT_URL_SSM", fake.URL)
			t.Setenv("REDLAMP_TEST_SECRET", tt.value)
			t.Setenv("REDLAMP_TEST_SECRET_PARAMETER", tt.parameter)
			if got := secret(t.Context(), "REDLAMP_TEST_SECRET"); got != tt.want {
				t.Errorf("secret = %q, want %q", got, tt.want)
			}
			if got := requests.Load(); got != tt.requests {
				t.Errorf("%d requests to Parameter Store, want %d", got, tt.requests)
			}
		})
	}
}

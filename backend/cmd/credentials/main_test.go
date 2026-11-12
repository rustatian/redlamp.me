package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/db"
	"github.com/rustatian/redlamp.me/backend/internal/dbtest"
	"github.com/rustatian/redlamp.me/backend/internal/notebook"
	"github.com/rustatian/redlamp.me/backend/migrations"
)

// The username rule: two different names of 3 to 40 letters, digits, dots, dashes or underscores.
func TestGenerateChecksTheUsernames(t *testing.T) {
	tests := []struct {
		name      string
		usernames []string
		wantErr   bool
	}{
		{name: "accepts two usernames", usernames: []string{"parent.name", "student_name-1"}},
		{name: "removes the spaces around a username", usernames: []string{" parent-name ", "student-name"}},
		{name: "rejects one username", usernames: []string{"parent-name"}, wantErr: true},
		{name: "rejects three usernames", usernames: []string{"parent-name", "student-name", "other-name"}, wantErr: true},
		{name: "rejects the same username two times", usernames: []string{"same-name", "same-name"}, wantErr: true},
		{name: "rejects a username of two characters", usernames: []string{"pa", "student-name"}, wantErr: true},
		{name: "rejects a username of 41 characters", usernames: []string{strings.Repeat("p", 41), "student-name"}, wantErr: true},
		{name: "rejects a space in a username", usernames: []string{"parent name", "student-name"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, passwords, err := generate(tt.usernames)
			if tt.wantErr {
				if err == nil {
					t.Errorf("generate(%q) = %s, want an error", tt.usernames, secret)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var accounts map[string]map[string]string
			if err := json.Unmarshal(secret, &accounts); err != nil {
				t.Fatal(err)
			}
			for i, role := range []string{"parent", "student"} {
				if got, want := accounts[role]["username"], strings.TrimSpace(tt.usernames[i]); got != want {
					t.Errorf("%s username = %q, want %q", role, got, want)
				}
				// 18 random bytes in base64url are 24 characters.
				if !regexp.MustCompile(`^[A-Za-z0-9_-]{24}$`).MatchString(passwords[i]) {
					t.Errorf("%s password = %q, want 24 base64url characters", role, passwords[i])
				}
				sum := sha256.Sum256([]byte(passwords[i]))
				if got := accounts[role]["passwordHash"]; got != hex.EncodeToString(sum[:]) {
					t.Errorf("%s passwordHash = %q, want the SHA-256 of the password", role, got)
				}
			}
			if passwords[0] == passwords[1] {
				t.Error("the two accounts have one password")
			}
			if want := `{"parent":{"username":"` + accounts["parent"]["username"] + `","passwordHash":"` + accounts["parent"]["passwordHash"] + `"},"student":{"username":"` + accounts["student"]["username"] + `","passwordHash":"` + accounts["student"]["passwordHash"] + `"}}`; string(secret) != want {
				t.Errorf("secret = %s, want the form %s", secret, want)
			}
		})
	}
}

func TestGeneratedAccountsSignIn(t *testing.T) {
	tests := []struct {
		name     string
		role     int
		username string
	}{
		{name: "the parent signs in with the new password", role: 0, username: "parent-name"},
		{name: "the student signs in with the new password", role: 1, username: "student-name"},
	}
	secret, passwords, err := generate([]string{"parent-name", "student-name"})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := dbtest.New(t)
			if _, err := db.Migrate(context.Background(), pool, migrations.FS, false); err != nil {
				t.Fatal(err)
			}
			s, err := notebook.New(notebook.Config{DB: pool, PublicOrigin: "https://redlamp.me", Accounts: string(secret)})
			if err != nil {
				t.Fatal(err)
			}
			r := httptest.NewRequest(http.MethodPost, "/leetcode/api/session", strings.NewReader(fmt.Sprintf(`{"username":%q,"password":%q}`, tt.username, passwords[tt.role])))
			r.Header.Set("Origin", "https://redlamp.me")
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.ServeHTTP(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("status = %d, body = %s; want 200", w.Code, w.Body.String())
			}
		})
	}
}

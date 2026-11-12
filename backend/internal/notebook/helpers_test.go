package notebook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// request sends an API request from the public origin, as request() of test/client.ts.
// path is the path after /leetcode/api. body is nil for no body, a string for the exact body text, or a value that becomes JSON.
func request(t *testing.T, s *Server, path, method string, body any, cookie string) *httptest.ResponseRecorder {
	t.Helper()
	text := ""
	switch value := body.(type) {
	case nil:
	case string:
		text = value
	default:
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		text = string(data)
	}
	return call(s, method, "/leetcode/api"+path, text, cookie)
}

// decode reads the JSON body of a response into a new value of the type T.
func decode[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(w.Body.Bytes(), &value); err != nil {
		t.Fatalf("status %d, body %q: %v", w.Code, w.Body.String(), err)
	}
	return value
}

// errorCode returns the error code of an error response, or "" for a response without one.
func errorCode(w *httptest.ResponseRecorder) string {
	var value struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &value)
	return value.Error.Code
}

// loginAs signs in the test account of the role "student" or "parent" and returns the session cookie.
func loginAs(t *testing.T, s *Server, role string) string {
	t.Helper()
	if role == "parent" {
		return login(t, s, "parent-test", parentPassword)
	}
	return login(t, s, "student-test", studentPassword)
}

// exec runs one SQL statement with PostgreSQL placeholders ($1, $2, ...) on the test database.
func exec(t *testing.T, pool *pgxpool.Pool, sql string, args ...any) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), sql, args...); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
}

// queryValue returns the first column of the first row of a query.
func queryValue[T any](t *testing.T, pool *pgxpool.Pool, sql string, args ...any) T {
	t.Helper()
	var value T
	if err := pool.QueryRow(context.Background(), sql, args...).Scan(&value); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	return value
}

// leetcodeCall is one request to the fake LeetCode.
type leetcodeCall struct {
	titleSlug string
	fields    string
}

// fakeLeetCode answers each question from reply and records the calls. A nil reply fails the request.
type fakeLeetCode struct {
	reply func(titleSlug, fields string) (json.RawMessage, error)
	calls []leetcodeCall
}

func (f *fakeLeetCode) Question(_ context.Context, titleSlug, fields string) (json.RawMessage, error) {
	f.calls = append(f.calls, leetcodeCall{titleSlug: titleSlug, fields: fields})
	if f.reply == nil {
		return nil, errors.New("the test expects no LeetCode request")
	}
	return f.reply(titleSlug, fields)
}

// newServerWithLeetCode returns a Server with the test accounts and a LeetCode client.
func newServerWithLeetCode(t *testing.T, pool *pgxpool.Pool, leetcode LeetCode) *Server {
	t.Helper()
	s, err := New(Config{DB: pool, PublicOrigin: publicOrigin, Accounts: testAccounts, LeetCode: leetcode})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// newServerWithSets returns a Server with the test accounts and the content of learning sets.
func newServerWithSets(t *testing.T, pool *pgxpool.Pool, sets []LearningSet) *Server {
	t.Helper()
	s, err := New(Config{DB: pool, PublicOrigin: publicOrigin, Accounts: testAccounts, LeetCode: &fakeLeetCode{}, Sets: sets})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// addProblem creates a problem through the API as the student and returns its ID.
func addProblem(t *testing.T, s *Server, cookie string, body any) string {
	t.Helper()
	w := request(t, s, "/problems", http.MethodPost, body, cookie)
	if w.Code != http.StatusCreated {
		t.Fatalf("add problem: status %d, body %s", w.Code, w.Body.String())
	}
	return decode[struct {
		ID string `json:"id"`
	}](t, w).ID
}

// race runs requests that are in progress at the same time. hold is a write that takes a row lock in an open transaction.
// Each request must wait for that lock. When all of them wait, the transaction rolls back and the requests continue.
func race(t *testing.T, pool *pgxpool.Pool, hold string, args []any, requests ...func() *httptest.ResponseRecorder) []*httptest.ResponseRecorder {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, hold, args...); err != nil {
		t.Fatal(err)
	}
	results := make([]*httptest.ResponseRecorder, len(requests))
	var wg sync.WaitGroup
	for i, send := range requests {
		wg.Go(func() { results[i] = send() })
	}
	waiting := 0
	for deadline := time.Now().Add(10 * time.Second); waiting != len(requests) && time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		// The sessions of the test pool have the schema name as their application_name.
		if err := pool.QueryRow(ctx, "SELECT count(*) FROM pg_stat_activity WHERE application_name = current_setting('application_name') AND wait_event_type = 'Lock'").Scan(&waiting); err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	if waiting != len(requests) {
		t.Fatalf("%d of %d requests waited for the lock of the other transaction", waiting, len(requests))
	}
	return results
}

package importer

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/db"
	"github.com/rustatian/redlamp.me/backend/internal/dbtest"
	"github.com/rustatian/redlamp.me/backend/migrations"
)

// fixture is a made-up export: one file in the format of wrangler d1 execute --json for each table.
const fixture = "testdata/export"

// d1File is the output of wrangler d1 execute --json for one SELECT.
type d1File []struct {
	Results []map[string]any `json:"results"`
	Success bool             `json:"success"`
	Meta    map[string]any   `json:"meta"`
}

// database returns a pool on a new schema with all migrations.
func database(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool := dbtest.New(t)
	if _, err := db.Migrate(context.Background(), pool, migrations.FS, false); err != nil {
		t.Fatal(err)
	}
	return pool
}

// copyFixture returns a new folder with a copy of the fixture.
func copyFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS(fixture)); err != nil {
		t.Fatal(err)
	}
	return dir
}

// edit decodes <dir>/<table>.json, lets change edit it, and writes it back. UseNumber keeps each number as written.
func edit(t *testing.T, dir, table string, change func(f d1File)) {
	t.Helper()
	path := filepath.Join(dir, table+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var f d1File
	if err := dec.Decode(&f); err != nil {
		t.Fatal(err)
	}
	change(f)
	writeFile(t, path, f)
}

func writeFile(t *testing.T, path string, f d1File) {
	t.Helper()
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func count(t *testing.T, pool *pgxpool.Pool, table string) int64 {
	t.Helper()
	var n int64
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM "+pgx.Identifier{table}.Sanitize()).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// names reports whether text has name as a whole word, so "homework" does not match "homework_tasks".
func names(text, name string) bool {
	return regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`).MatchString(text)
}

func TestRunImportsTheFixture(t *testing.T) {
	// The row counts of the fixture files, in the order of the rule: parents before children.
	counts := []struct {
		name string
		rows int64
	}{
		{name: "problems", rows: 3},
		{name: "homework", rows: 3},
		{name: "attempts", rows: 3},
		{name: "homework_tasks", rows: 4},
		{name: "submissions", rows: 2},
		{name: "feedback", rows: 3},
		{name: "reviews", rows: 2},
		{name: "discussion_reads", rows: 2},
		{name: "homework_sets", rows: 2},
		{name: "problem_sets", rows: 3},
		{name: "problem_tests", rows: 2},
		{name: "test_runs", rows: 2},
		{name: "postponements", rows: 2},
		{name: "tutor_messages", rows: 3},
		{name: "card_answers", rows: 3},
		{name: "quiz_attempts", rows: 2},
		{name: "earned_badges", rows: 2},
		{name: "visits", rows: 3},
	}
	// The values of the fixture files, read back byte for byte.
	values := []struct {
		name  string
		query string
		want  pgtype.Text
	}{
		{name: "a task with a current submission gets it in the second step", query: "SELECT current_submission_id FROM homework_tasks WHERE id = 'task-1'", want: pgtype.Text{String: "submission-2", Valid: true}},
		{name: "a task without a current submission keeps NULL", query: "SELECT current_submission_id FROM homework_tasks WHERE id = 'task-2'", want: pgtype.Text{}},
		{name: "a chapter task keeps its set and chapter", query: "SELECT kind || ' ' || set_slug || ' ' || chapter FROM homework_tasks WHERE id = 'task-3'", want: pgtype.Text{String: "chapter two-pointers 3", Valid: true}},
		{name: "a text with Russian characters", query: "SELECT title FROM homework WHERE id = 'homework-1'", want: pgtype.Text{String: "Неделя 41: два указателя", Valid: true}},
		{name: "a text with quotes, a backslash and a line break", query: "SELECT body FROM feedback WHERE id = 'feedback-1'", want: pgtype.Text{String: "Use \"left\" and \"right\".\nThe path is C:\\notes\\new.", Valid: true}},
		{name: "an attempt document keeps its JSON text", query: "SELECT document FROM attempts WHERE id = 'attempt-2'", want: pgtype.Text{String: `{"version": 1, "blocks": [{"type": "code", "text": "for i in range(n):\n    pass"}]}`, Valid: true}},
		{name: "a whole number of 10 digits", query: "SELECT number::text FROM problems WHERE id = 'problem-2'", want: pgtype.Text{String: "2147483648", Valid: true}},
		{name: "a NULL number", query: "SELECT number::text FROM problems WHERE id = 'problem-3'", want: pgtype.Text{}},
		{name: "an empty text stays empty", query: "SELECT quote FROM tutor_messages WHERE id = 'message-2'", want: pgtype.Text{String: "", Valid: true}},
	}
	ctx := context.Background()
	pool := database(t)
	var out bytes.Buffer
	if err := Run(ctx, pool, fixture, &out); err != nil {
		t.Fatalf("Run: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != len(counts) {
		t.Errorf("output has %d lines, want %d:\n%s", len(lines), len(counts), out.String())
	}
	var tables []string
	for _, tt := range counts {
		tables = append(tables, tt.name)
		t.Run(tt.name, func(t *testing.T) {
			if got := count(t, pool, tt.name); got != tt.rows {
				t.Errorf("rows = %d, want %d", got, tt.rows)
			}
			n := 0
			for _, line := range lines {
				if names(line, tt.name) {
					n++
				}
			}
			if n != 1 {
				t.Errorf("%d output lines name the table, want 1:\n%s", n, out.String())
			}
		})
	}
	if !slices.Equal(Tables, tables) {
		t.Errorf("Tables = %v, want %v", Tables, tables)
	}
	for _, tt := range values {
		t.Run(tt.name, func(t *testing.T) {
			var got pgtype.Text
			if err := pool.QueryRow(ctx, tt.query).Scan(&got); err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRunCommitsAtMost3000RowsAtATime(t *testing.T) {
	dir := copyFixture(t)
	visits := d1File{{Success: true, Meta: map[string]any{"duration": 0}}}
	first := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range 6500 {
		visits[0].Results = append(visits[0].Results, map[string]any{"day": first.AddDate(0, 0, i).Format(time.DateOnly)})
	}
	writeFile(t, filepath.Join(dir, "visits.json"), visits)

	ctx := context.Background()
	pool := database(t)
	if err := Run(ctx, pool, dir, io.Discard); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// The rows that one transaction inserted share one xmin.
	var rows, batches, largest int64
	err := pool.QueryRow(ctx, "SELECT sum(n)::bigint, count(*), max(n) FROM (SELECT count(*) AS n FROM visits GROUP BY xmin) AS batches").Scan(&rows, &batches, &largest)
	if err != nil {
		t.Fatal(err)
	}
	if rows != 6500 {
		t.Errorf("rows = %d, want 6500", rows)
	}
	if largest > 3000 {
		t.Errorf("largest batch = %d rows, want 3000 or fewer", largest)
	}
	if batches < 3 {
		t.Errorf("batches = %d, want 3 or more", batches)
	}
}

func TestRunRejectsBadInputAndInsertsNothing(t *testing.T) {
	// Each case changes the file of table, or runs sql before Run. A nil change removes the file.
	// The changes are in tables near the end of the order, so Run must check all files before it inserts a row.
	tests := []struct {
		name   string
		table  string
		change func(f d1File)
		sql    string
	}{
		{name: "a file is missing", table: "visits"},
		{name: "a file has success false", table: "earned_badges", change: func(f d1File) { f[0].Success = false }},
		{name: "a row has a column that the table does not have", table: "quiz_attempts", change: func(f d1File) { f[0].Results[0]["score"] = json.Number("2") }},
		{name: "a row lacks a nullable column of the table", table: "tutor_messages", change: func(f d1File) { delete(f[0].Results[0], "quote") }},
		{name: "a bigint column has a string", table: "quiz_attempts", change: func(f d1File) { f[0].Results[0]["total"] = "5" }},
		{name: "a bigint column has a number with a fraction", table: "test_runs", change: func(f d1File) { f[0].Results[0]["passed"] = json.Number("2.5") }},
		{name: "a text column has a number", table: "visits", change: func(f d1File) { f[0].Results[0]["day"] = json.Number("20261005") }},
		{name: "a table has a row before the import", sql: "INSERT INTO visits (day) VALUES ('2020-01-01')"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			pool := database(t)
			dir := copyFixture(t)
			switch {
			case tt.sql != "":
				if _, err := pool.Exec(ctx, tt.sql); err != nil {
					t.Fatal(err)
				}
			case tt.change == nil:
				if err := os.Remove(filepath.Join(dir, tt.table+".json")); err != nil {
					t.Fatal(err)
				}
			default:
				edit(t, dir, tt.table, tt.change)
			}
			if err := Run(ctx, pool, dir, io.Discard); err == nil {
				t.Error("Run: no error, want one")
			}
			if n := count(t, pool, "problems"); n != 0 {
				t.Errorf("problems has %d rows after the failed Run, want 0", n)
			}
		})
	}
}

func TestVerifyFindsDifferences(t *testing.T) {
	// Each case runs one statement after a good import. tables are the tables that differ from the files.
	tests := []struct {
		name   string
		sql    string
		tables []string
	}{
		{name: "no change"},
		{name: "a row is missing", sql: "DELETE FROM visits WHERE day = '2026-10-06'", tables: []string{"visits"}},
		{name: "a row is extra", sql: "INSERT INTO earned_badges (id) VALUES ('streak-30')", tables: []string{"earned_badges"}},
		{name: "one character of a text differs", sql: "UPDATE homework SET title = 'Неделя 41: два указатели' WHERE id = 'homework-1'", tables: []string{"homework"}},
		{name: "a NULL in place of an empty string", sql: "UPDATE tutor_messages SET quote = NULL WHERE id = 'message-2'", tables: []string{"tutor_messages"}},
		{name: "an empty string in place of a NULL", sql: "UPDATE tutor_messages SET quote = '' WHERE id = 'message-1'", tables: []string{"tutor_messages"}},
		{name: "a whole number differs", sql: "UPDATE problems SET number = 2147483649 WHERE id = 'problem-2'", tables: []string{"problems"}},
		{name: "a current submission differs", sql: "UPDATE homework_tasks SET current_submission_id = 'submission-1' WHERE id = 'task-1'", tables: []string{"homework_tasks"}},
		{name: "two tables differ", sql: "WITH gone AS (DELETE FROM visits WHERE day = '2026-10-06') INSERT INTO earned_badges (id) VALUES ('streak-30')", tables: []string{"visits", "earned_badges"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			pool := database(t)
			if err := Run(ctx, pool, fixture, io.Discard); err != nil {
				t.Fatalf("Run: %v", err)
			}
			if tt.sql != "" {
				if _, err := pool.Exec(ctx, tt.sql); err != nil {
					t.Fatal(err)
				}
			}
			err := Verify(ctx, pool, fixture, io.Discard)
			if len(tt.tables) == 0 {
				if err != nil {
					t.Errorf("Verify: %v, want no error", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Verify: no error, want one that names %v", tt.tables)
			}
			for _, table := range tt.tables {
				if !names(err.Error(), table) {
					t.Errorf("Verify error %q does not name %s", err, table)
				}
			}
		})
	}
}

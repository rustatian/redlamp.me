package db

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/dbtest"
	"github.com/rustatian/redlamp.me/backend/migrations"
)

func baseline(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool := dbtest.New(t)
	if _, err := Migrate(context.Background(), pool, migrations.FS, false); err != nil {
		t.Fatal(err)
	}
	return pool
}

// The expected tables and columns are the D1 schema after migrations/0017_homework_content.sql, without tutor_account.
func TestBaselineTables(t *testing.T) {
	tests := []struct {
		name    string
		columns string
	}{
		{name: "attempts", columns: "id,problem_id,state,document,search_text,acceptance,understanding,version,created_at,updated_at,saved_at"},
		{name: "card_answers", columns: "id,card_id,correct,interval_days,answered_on,next_date,created_at"},
		{name: "discussion_reads", columns: "role,problem_id,read_at"},
		{name: "earned_badges", columns: "id"},
		{name: "feedback", columns: "id,problem_id,task_id,submission_id,author,kind,body,created_at"},
		{name: "homework", columns: "id,title,instructions,due_date,version,created_at,updated_at"},
		{name: "homework_sets", columns: "homework_id,set_slug"},
		{name: "homework_tasks", columns: "id,homework_id,problem_id,state,current_submission_id,version,created_at,updated_at,kind,set_slug,chapter,completed_at"},
		{name: "judge_runs", columns: "problem_id,seed,created_at"},
		{name: "postponements", columns: "homework_id,from_date,to_date,created_at"},
		{name: "problem_sets", columns: "problem_id,set_slug"},
		{name: "problem_tests", columns: "id,problem_id,input,created_at"},
		{name: "problems", columns: "id,slug,number,title,difficulty,topics,summary,search_text,understanding,next_review_date,archived_at,version,created_at,updated_at,progress_version"},
		{name: "quiz_attempts", columns: "id,deck,correct,total,passed,created_at"},
		{name: "reviews", columns: "id,problem_id,result,note,reviewed_on,next_review_date,created_at"},
		{name: "sessions", columns: "token_hash,role,credential_fingerprint,created_at,expires_at"},
		{name: "submissions", columns: "id,task_id,attempt_id,task_version,created_at"},
		{name: "test_runs", columns: "id,problem_id,attempt_id,approach_id,code_hash,passed,total,run_on,created_at"},
		{name: "tutor_messages", columns: "id,context_kind,context_ref,author,body,quote,created_at"},
		{name: "visits", columns: "day"},
	}
	ctx := context.Background()
	pool := baseline(t)
	want := []string{"schema_migrations"}
	for _, tt := range tests {
		want = append(want, tt.name)
		t.Run(tt.name, func(t *testing.T) {
			var columns string
			err := pool.QueryRow(ctx, "SELECT coalesce(string_agg(column_name, ',' ORDER BY ordinal_position), '') FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = $1", tt.name).Scan(&columns)
			if err != nil {
				t.Fatal(err)
			}
			if columns != tt.columns {
				t.Errorf("columns = %s, want %s", columns, tt.columns)
			}
		})
	}
	slices.Sort(want)
	var tables []string
	if err := pool.QueryRow(ctx, "SELECT array_agg(table_name::text ORDER BY table_name) FROM information_schema.tables WHERE table_schema = current_schema()").Scan(&tables); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(tables, want) {
		t.Errorf("tables = %v, want %v", tables, want)
	}
}

// The expected indexes are the named indexes of the D1 schema after migrations/0017_homework_content.sql.
func TestBaselineIndexes(t *testing.T) {
	tests := []struct {
		name    string
		table   string
		columns string
		unique  bool
	}{
		{name: "sessions_expiry", table: "sessions", columns: "expires_at"},
		{name: "one_draft_per_problem", table: "attempts", columns: "problem_id", unique: true},
		{name: "attempts_problem_state", table: "attempts", columns: "problem_id,state,saved_at"},
		{name: "one_active_task", table: "homework_tasks", columns: "problem_id", unique: true},
		{name: "one_active_content", table: "homework_tasks", columns: "kind,set_slug,chapter", unique: true},
		{name: "homework_tasks_homework", table: "homework_tasks", columns: "homework_id"},
		{name: "feedback_problem", table: "feedback", columns: "problem_id,created_at"},
		{name: "reviews_problem_date", table: "reviews", columns: "problem_id,reviewed_on"},
		{name: "tutor_messages_thread", table: "tutor_messages", columns: "context_kind,context_ref,created_at"},
		{name: "problem_tests_problem", table: "problem_tests", columns: "problem_id,created_at"},
		{name: "test_runs_problem", table: "test_runs", columns: "problem_id,run_on"},
	}
	ctx := context.Background()
	pool := baseline(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var table, columns string
			var unique bool
			err := pool.QueryRow(ctx, `SELECT t.relname, i.indisunique, (SELECT string_agg(a.attname, ',' ORDER BY k.n)
					FROM unnest(i.indkey::int2[]) WITH ORDINALITY AS k (attnum, n) JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = k.attnum)
				FROM pg_index i JOIN pg_class t ON t.oid = i.indrelid WHERE i.indexrelid = to_regclass($1)`, tt.name).Scan(&table, &unique, &columns)
			if err != nil {
				t.Fatal(err)
			}
			if table != tt.table || columns != tt.columns || unique != tt.unique {
				t.Errorf("index = %s(%s) unique %v, want %s(%s) unique %v", table, columns, unique, tt.table, tt.columns, tt.unique)
			}
		})
	}
}

// The cases run in order on one database: an accepted statement stays for the cases after it.
func TestBaselineRules(t *testing.T) {
	ctx := context.Background()
	pool := baseline(t)
	for _, sql := range []string{
		"INSERT INTO problems (id, slug, title, difficulty, topics, summary, search_text, created_at, updated_at) VALUES ('p1', 'two-sum', 'Two Sum', 'easy', '[]', '', 'two sum', 't0', 't0')",
		"INSERT INTO attempts (id, problem_id, state, document, search_text, acceptance, understanding, created_at, updated_at) VALUES ('a1', 'p1', 'draft', '{}', '', 'not_submitted', 'needs_practice', 't0', 't0')",
		"INSERT INTO homework (id, title, instructions, created_at, updated_at) VALUES ('h1', 'Week 41', '', 't0', 't0')",
		"INSERT INTO homework_tasks (id, homework_id, problem_id, state, created_at, updated_at) VALUES ('t1', 'h1', 'p1', 'assigned', 't0', 't0')",
		"INSERT INTO homework_tasks (id, homework_id, state, created_at, updated_at, kind, set_slug, chapter) VALUES ('t2', 'h1', 'assigned', 't0', 't0', 'chapter', 'two-pointers', 1)",
		"INSERT INTO submissions VALUES ('s1', 't1', 'a1', 2, 't0')",
		"INSERT INTO card_answers VALUES ('c1', 'card', 1, 1, '2026-10-07', '2026-10-08', 't0')",
		"INSERT INTO postponements VALUES ('h1', '2026-10-07', '2026-10-08', 't0')",
	} {
		if _, err := pool.Exec(ctx, sql); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	const attempt = "INSERT INTO attempts (id, problem_id, state, document, search_text, acceptance, understanding, created_at, updated_at) VALUES "
	const problemTask = "INSERT INTO homework_tasks (id, homework_id, problem_id, state, created_at, updated_at) VALUES "
	const contentTask = "INSERT INTO homework_tasks (id, homework_id, problem_id, state, created_at, updated_at, kind, set_slug, chapter) VALUES "
	tests := []struct {
		name     string
		sql      string
		wantCode string
	}{
		{name: "one_draft_per_problem rejects a second draft of a problem", sql: attempt + "('a2', 'p1', 'draft', '{}', '', 'not_submitted', 'needs_practice', 't0', 't0')", wantCode: "23505"},
		{name: "one_draft_per_problem accepts a saved attempt next to the draft", sql: attempt + "('a3', 'p1', 'saved', '{}', '', 'accepted', 'with_help', 't0', 't0')"},
		{name: "one_active_task rejects a second active homework task of a problem", sql: problemTask + "('t3', 'h1', 'p1', 'changes_requested', 't0', 't0')", wantCode: "23505"},
		{name: "one_active_task accepts a completed homework task next to the active one", sql: problemTask + "('t4', 'h1', 'p1', 'completed', 't0', 't0')"},
		{name: "one_active_content rejects a second active homework task of a chapter", sql: contentTask + "('t5', 'h1', NULL, 'in_progress', 't0', 't0', 'chapter', 'two-pointers', 1)", wantCode: "23505"},
		{name: "one_active_content accepts a quiz task of the same chapter", sql: contentTask + "('t6', 'h1', NULL, 'assigned', 't0', 't0', 'quiz', 'two-pointers', 1)"},
		{name: "one_active_content accepts a completed homework task of the chapter", sql: contentTask + "('t7', 'h1', NULL, 'completed', 't0', 't0', 'chapter', 'two-pointers', 1)"},
		{name: "a homework task rejects a current submission that does not exist", sql: "UPDATE homework_tasks SET current_submission_id = 'missing' WHERE id = 't1'", wantCode: "23503"},
		{name: "a homework task accepts its submission as the current one", sql: "UPDATE homework_tasks SET current_submission_id = 's1' WHERE id = 't1'"},
		{name: "a submission rejects a homework task that does not exist", sql: "INSERT INTO submissions VALUES ('s2', 'missing', 'a1', 1, 't0')", wantCode: "23503"},
		{name: "a homework task version has one submission", sql: "INSERT INTO submissions VALUES ('s3', 't1', 'a1', 2, 't0')", wantCode: "23505"},
		{name: "a chapter task has no problem", sql: contentTask + "('t8', 'h1', 'p1', 'completed', 't0', 't0', 'chapter', 'heap', 1)", wantCode: "23514"},
		{name: "a chapter task has no review state", sql: contentTask + "('t9', 'h1', NULL, 'submitted', 't0', 't0', 'chapter', 'heap', 1)", wantCode: "23514"},
		{name: "a chapter number is above zero", sql: contentTask + "('t10', 'h1', NULL, 'completed', 't0', 't0', 'chapter', 'heap', 0)", wantCode: "23514"},
		{name: "a problem task has a problem", sql: problemTask + "('t11', 'h1', NULL, 'completed', 't0', 't0')", wantCode: "23514"},
		{name: "an attempt document is JSON", sql: attempt + "('a4', 'p1', 'saved', '{', '', 'accepted', 'with_help', 't0', 't0')", wantCode: "23514"},
		{name: "an attempt document of 256000 bytes is stored", sql: attempt + "('a5', 'p1', 'saved', '\"' || repeat('я', 127999) || '\"', '', 'accepted', 'with_help', 't0', 't0')"},
		{name: "an attempt document of 256002 bytes in 128002 characters is too large", sql: attempt + "('a6', 'p1', 'saved', '\"' || repeat('я', 128000) || '\"', '', 'accepted', 'with_help', 't0', 't0')", wantCode: "23514"},
		{name: "problem topics are JSON", sql: "INSERT INTO problems (id, slug, title, difficulty, topics, summary, search_text, created_at, updated_at) VALUES ('p2', 'subsets', 'Subsets', 'medium', 'Array', '', '', 't0', 't0')", wantCode: "23514"},
		{name: "a problem slug is unique", sql: "INSERT INTO problems (id, slug, title, difficulty, topics, summary, search_text, created_at, updated_at) VALUES ('p3', 'two-sum', 'Two Sum', 'easy', '[]', '', '', 't0', 't0')", wantCode: "23505"},
		{name: "the level badge is not stored", sql: "INSERT INTO earned_badges VALUES ('level-five')", wantCode: "23514"},
		{name: "a quiz attempt has no more correct answers than cards", sql: "INSERT INTO quiz_attempts VALUES ('q1', 'deck', 3, 2, 1, 't0')", wantCode: "23514"},
		{name: "a card has one answer for each day", sql: "INSERT INTO card_answers VALUES ('c2', 'card', 0, 1, '2026-10-07', '2026-10-08', 't0')", wantCode: "23505"},
		{name: "a card interval is one of the five steps", sql: "INSERT INTO card_answers VALUES ('c3', 'card', 1, 2, '2026-10-09', '2026-10-11', 't0')", wantCode: "23514"},
		{name: "a homework has one postponement", sql: "INSERT INTO postponements VALUES ('h1', '2026-10-08', '2026-10-09', 't0')", wantCode: "23505"},
		{name: "a judge seed uses 32 bits without a sign", sql: "INSERT INTO judge_runs VALUES ('p1', 4294967295, 't0')"},
		{name: "a session time is in milliseconds", sql: "INSERT INTO sessions VALUES ('hash', 'student', 'fingerprint', 1794787200000, 1795392000000)"},
		{name: "the ChatGPT connection has no table", sql: "SELECT 1 FROM tutor_account", wantCode: "42P01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, tt.sql)
			if got := sqlState(err); got != tt.wantCode {
				t.Errorf("SQLSTATE %q (%v), want %q\n%s", got, err, tt.wantCode, strings.TrimSpace(tt.sql))
			}
		})
	}
}

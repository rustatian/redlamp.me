package db

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rustatian/redlamp.me/backend/internal/dbtest"
)

func sqlFiles(files map[string]string) fstest.MapFS {
	fsys := fstest.MapFS{}
	for name, sql := range files {
		fsys[name] = &fstest.MapFile{Data: []byte(sql)}
	}
	return fsys
}

func TestMigrate(t *testing.T) {
	ab := map[string]string{"0002_b.sql": "CREATE TABLE b (id int)", "0001_a.sql": "CREATE TABLE a (id int);\n"}
	tests := []struct {
		name         string
		files        map[string]string
		runs         int
		wantApplied  []string
		wantErr      string
		wantVersions []string
	}{
		{name: "applies files in name order", files: ab, runs: 1, wantApplied: []string{"0001_a", "0002_b"}, wantVersions: []string{"0001_a", "0002_b"}},
		{name: "second run applies nothing", files: ab, runs: 2, wantApplied: nil, wantVersions: []string{"0001_a", "0002_b"}},
		{name: "a failed file keeps earlier versions", files: map[string]string{"0001_a.sql": "CREATE TABLE a (id int)", "0002_bad.sql": "CREATE TABLE"}, runs: 1, wantApplied: []string{"0001_a"}, wantErr: "0002_bad.sql", wantVersions: []string{"0001_a"}},
		{name: "rejects a file with two statements", files: map[string]string{"0001_two.sql": "CREATE TABLE a (id int); CREATE TABLE b (id int);"}, runs: 1, wantApplied: nil, wantErr: "0001_two.sql: one statement per file", wantVersions: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			pool := dbtest.New(t)
			var applied []string
			var err error
			for range tt.runs {
				applied, err = Migrate(ctx, pool, sqlFiles(tt.files), false)
			}
			if tt.wantErr == "" && err != nil {
				t.Fatalf("Migrate: %v", err)
			}
			if tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Fatalf("Migrate error = %v, want it to contain %q", err, tt.wantErr)
			}
			if !slices.Equal(applied, tt.wantApplied) {
				t.Errorf("applied = %v, want %v", applied, tt.wantApplied)
			}
			rows, err := pool.Query(ctx, "SELECT version FROM schema_migrations ORDER BY version")
			if err != nil {
				t.Fatal(err)
			}
			var versions []string
			for rows.Next() {
				var v string
				if err := rows.Scan(&v); err != nil {
					t.Fatal(err)
				}
				versions = append(versions, v)
			}
			if !slices.Equal(versions, tt.wantVersions) {
				t.Errorf("schema_migrations = %v, want %v", versions, tt.wantVersions)
			}
		})
	}
}

func TestMigrateIndexFiles(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		wantUnique    bool
		wantDuplicate string
	}{
		{name: "builds an index from an ASYNC statement", file: "CREATE INDEX ASYNC a_index ON a (id);", wantUnique: false, wantDuplicate: ""},
		{name: "builds a unique partial index from an ASYNC statement after a comment", file: "-- One row with the kind x for each id.\nCREATE UNIQUE INDEX ASYNC a_index ON a (id) WHERE kind = 'x';\n", wantUnique: true, wantDuplicate: "23505"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			pool := dbtest.New(t)
			files := sqlFiles(map[string]string{"0001_a.sql": "CREATE TABLE a (id int, kind text)", "0002_a_index.sql": tt.file})
			applied, err := Migrate(ctx, pool, files, false)
			if err != nil {
				t.Fatalf("Migrate: %v", err)
			}
			if want := []string{"0001_a", "0002_a_index"}; !slices.Equal(applied, want) {
				t.Errorf("applied = %v, want %v", applied, want)
			}
			var unique bool
			if err := pool.QueryRow(ctx, "SELECT indisunique FROM pg_index WHERE indexrelid = 'a_index'::regclass").Scan(&unique); err != nil {
				t.Fatal(err)
			}
			if unique != tt.wantUnique {
				t.Errorf("unique = %v, want %v", unique, tt.wantUnique)
			}
			if _, err := pool.Exec(ctx, "INSERT INTO a VALUES (1, 'y'), (1, 'y')"); err != nil {
				t.Errorf("two rows outside the index predicate: %v", err)
			}
			_, err = pool.Exec(ctx, "INSERT INTO a VALUES (1, 'x'), (1, 'x')")
			if got := sqlState(err); got != tt.wantDuplicate {
				t.Errorf("two rows with the kind x: SQLSTATE %q (%v), want %q", got, err, tt.wantDuplicate)
			}
		})
	}
}

// sqlState returns the SQLSTATE of err, or "" for no error.
func sqlState(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

func TestSetupAppRoleRejectsBadARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
	}{
		{name: "rejects a quote", arn: "arn:aws:iam::123456789012:role/a'b"},
		{name: "rejects a user ARN", arn: "arn:aws:iam::123456789012:user/x"},
		{name: "rejects a short account", arn: "arn:aws:iam::123:role/x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetupAppRole(context.Background(), nil, tt.arn)
			if err == nil || !strings.Contains(err.Error(), "invalid IAM role ARN") {
				t.Errorf("SetupAppRole(%q) error = %v, want an invalid ARN error", tt.arn, err)
			}
		})
	}
}

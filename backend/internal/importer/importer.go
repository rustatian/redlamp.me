// Package importer moves the notebook data from the D1 database of the Worker into the new database. It runs one time, at the cutover.
package importer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rustatian/redlamp.me/backend/internal/db"
)

// Tables are the tables that move, parents before children.
// The tables sessions and judge_runs get no rows: a user signs in again, and an open seed of the hidden tests is valid for one run.
// The table quiz_decks gets no rows: a quiz attempt writes the row of its deck.
var Tables = []string{
	"problems", "homework", "attempts", "homework_tasks", "submissions", "feedback", "reviews", "discussion_reads", "homework_sets",
	"problem_sets", "problem_tests", "test_runs", "postponements", "tutor_messages", "card_answers", "quiz_attempts", "earned_badges", "visits",
}

// One transaction of Aurora DSQL changes 3,000 rows and 10 MiB at most: https://docs.aws.amazon.com/aurora-dsql/latest/userguide/CHAP_quotas.html
// A statement of PostgreSQL has 65,535 parameters at most: https://www.postgresql.org/docs/16/limits.html
const (
	maxRows       = 3000
	maxBytes      = 8 << 20
	maxParameters = 65535
)

type column struct {
	name    string
	integer bool
}

// table has the rows of one file. Each row has its values in the order of the columns of the database: a string, an int64 or nil.
type table struct {
	name    string
	columns []column
	rows    [][]any
}

func (t table) names() string {
	names := make([]string, len(t.columns))
	for i, c := range t.columns {
		names[i] = pgx.Identifier{c.name}.Sanitize()
	}
	return strings.Join(names, ",")
}

// load reads the file of each table and checks each row against the columns of the database.
// A file is the output of "wrangler d1 execute --json" for one SELECT of the whole table.
func load(ctx context.Context, pool *pgxpool.Pool, dir string) ([]table, error) {
	tables := make([]table, 0, len(Tables))
	for _, name := range Tables {
		t := table{name: name}
		rows, err := pool.Query(ctx, "SELECT column_name,data_type FROM information_schema.columns WHERE table_schema=current_schema() AND table_name=$1 ORDER BY ordinal_position", name)
		if err != nil {
			return nil, err
		}
		var kind string
		var c column
		if _, err := pgx.ForEachRow(rows, []any{&c.name, &kind}, func() error {
			if kind != "text" && kind != "bigint" {
				return fmt.Errorf("%s.%s has the type %s, and the import knows text and bigint", name, c.name, kind)
			}
			c.integer = kind == "bigint"
			t.columns = append(t.columns, c)
			return nil
		}); err != nil {
			return nil, err
		}
		if len(t.columns) == 0 {
			return nil, fmt.Errorf("the database has no table %s", name)
		}
		data, err := os.ReadFile(filepath.Join(dir, name+".json"))
		if err != nil {
			return nil, err
		}
		var export []struct {
			Results []map[string]any `json:"results"`
			Success bool             `json:"success"`
		}
		// A number stays text, so a large whole number does not lose digits.
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if err := decoder.Decode(&export); err != nil {
			return nil, fmt.Errorf("%s.json: %w", name, err)
		}
		if len(export) != 1 || !export[0].Success {
			return nil, fmt.Errorf("%s.json is not the result of one statement that passed", name)
		}
		for i, source := range export[0].Results {
			if len(source) != len(t.columns) {
				return nil, fmt.Errorf("%s.json row %d has %d columns, and the table has %d", name, i+1, len(source), len(t.columns))
			}
			row := make([]any, len(t.columns))
			for j, c := range t.columns {
				value, ok := source[c.name]
				if !ok {
					return nil, fmt.Errorf("%s.json row %d has no column %s", name, i+1, c.name)
				}
				switch v := value.(type) {
				case nil:
				case string:
					if c.integer {
						return nil, fmt.Errorf("%s.json row %d: %s is a text, and the column is a whole number", name, i+1, c.name)
					}
					row[j] = v
				case json.Number:
					number, err := v.Int64()
					if err != nil || !c.integer {
						return nil, fmt.Errorf("%s.json row %d: %s is the number %s, and the column does not take it", name, i+1, c.name, v)
					}
					row[j] = number
				default:
					return nil, fmt.Errorf("%s.json row %d: %s is not a text, a whole number or null", name, i+1, c.name)
				}
			}
			t.rows = append(t.rows, row)
		}
		tables = append(tables, t)
	}
	return tables, nil
}

// insert stores the rows of a table in batches. Each batch is one statement in one transaction.
func insert(ctx context.Context, pool *pgxpool.Pool, t table) error {
	limit := min(maxRows, maxParameters/len(t.columns))
	for start := 0; start < len(t.rows); {
		var values []string
		var args []any
		size := 0
		for end := start; end < len(t.rows) && end-start < limit; end++ {
			row := t.rows[end]
			for _, value := range row {
				if text, ok := value.(string); ok {
					size += len(text)
				}
			}
			if end > start && size > maxBytes {
				break
			}
			marks := make([]string, len(row))
			for i, value := range row {
				// A task row gets its submission in the second step, because the submission row points to the task row.
				if t.name == "homework_tasks" && t.columns[i].name == "current_submission_id" {
					value = nil
				}
				args = append(args, value)
				marks[i] = fmt.Sprintf("$%d", len(args))
			}
			values = append(values, "("+strings.Join(marks, ",")+")")
		}
		statement := "INSERT INTO " + t.name + " (" + t.names() + ") VALUES " + strings.Join(values, ",")
		if err := db.Tx(ctx, pool, func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, statement, args...)
			return err
		}); err != nil {
			return fmt.Errorf("%s rows %d to %d: %w", t.name, start+1, start+len(values), err)
		}
		start += len(values)
	}
	return nil
}

// link gives each homework task its current submission, after the submissions are in the database.
func link(ctx context.Context, pool *pgxpool.Pool, tasks table) error {
	id := slices.IndexFunc(tasks.columns, func(c column) bool { return c.name == "id" })
	submission := slices.IndexFunc(tasks.columns, func(c column) bool { return c.name == "current_submission_id" })
	var linked [][]any
	for _, row := range tasks.rows {
		if row[submission] != nil {
			linked = append(linked, row)
		}
	}
	for batch := range slices.Chunk(linked, maxRows) {
		if err := db.Tx(ctx, pool, func(tx pgx.Tx) error {
			for _, row := range batch {
				if _, err := tx.Exec(ctx, "UPDATE homework_tasks SET current_submission_id=$1 WHERE id=$2", row[submission], row[id]); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return fmt.Errorf("homework_tasks.current_submission_id: %w", err)
		}
	}
	return nil
}

// checksum is the SHA-256 of the rows in a fixed order: each row is one line of JSON, and the lines are in byte order.
func checksum(rows [][]any) (string, error) {
	lines := make([]string, len(rows))
	for i, row := range rows {
		line, err := json.Marshal(row)
		if err != nil {
			return "", err
		}
		lines[i] = string(line)
	}
	slices.Sort(lines)
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return hex.EncodeToString(sum[:]), nil
}

// verify compares the rows of each table in the database with the rows of its file: the row count and the checksum.
func verify(ctx context.Context, pool *pgxpool.Pool, tables []table, out io.Writer) error {
	var different []string
	for _, t := range tables {
		rows, err := pool.Query(ctx, "SELECT "+t.names()+" FROM "+t.name)
		if err != nil {
			return err
		}
		stored, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) ([]any, error) { return row.Values() })
		if err != nil {
			return err
		}
		want, err := checksum(t.rows)
		if err != nil {
			return err
		}
		got, err := checksum(stored)
		if err != nil {
			return err
		}
		state := "equal"
		if len(stored) != len(t.rows) || got != want {
			state = "DIFFERENT"
			different = append(different, t.name)
		}
		if _, err := fmt.Fprintf(out, "%-16s file: %6d rows %s  database: %6d rows %s  %s\n", t.name, len(t.rows), want[:12], len(stored), got[:12], state); err != nil {
			return err
		}
	}
	if len(different) > 0 {
		return fmt.Errorf("the database differs from the files in: %s", strings.Join(different, ", "))
	}
	return nil
}

// Verify compares the database with the files of dir: for each table the row count and the checksum over the ordered rows.
// It writes one line for each table and returns an error when a table differs.
func Verify(ctx context.Context, pool *pgxpool.Pool, dir string, out io.Writer) error {
	tables, err := load(ctx, pool, dir)
	if err != nil {
		return err
	}
	return verify(ctx, pool, tables, out)
}

// Run reads one file for each table from dir, stores the rows in the empty tables of the database and then compares the database with the files.
// It checks all files before the first write.
func Run(ctx context.Context, pool *pgxpool.Pool, dir string, out io.Writer) error {
	tables, err := load(ctx, pool, dir)
	if err != nil {
		return err
	}
	for _, t := range tables {
		var used bool
		if err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM "+t.name+")").Scan(&used); err != nil {
			return err
		}
		if used {
			return fmt.Errorf("the table %s has rows: the import needs empty tables", t.name)
		}
	}
	for _, t := range tables {
		if err := insert(ctx, pool, t); err != nil {
			return err
		}
	}
	tasks := slices.IndexFunc(tables, func(t table) bool { return t.name == "homework_tasks" })
	if err := link(ctx, pool, tables[tasks]); err != nil {
		return err
	}
	return verify(ctx, pool, tables, out)
}

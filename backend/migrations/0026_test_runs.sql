-- The test runs of the student. code_hash is the SHA-256 (hex) of the code of the run. run_on is the local date of the student.
CREATE TABLE test_runs (
	id text PRIMARY KEY,
	problem_id text NOT NULL REFERENCES problems (id),
	attempt_id text NOT NULL REFERENCES attempts (id),
	approach_id text NOT NULL,
	code_hash text NOT NULL,
	passed bigint NOT NULL,
	total bigint NOT NULL,
	run_on text NOT NULL,
	created_at text NOT NULL
);

-- The open seeds of the hidden tests. The judge accepts a seed one time. A seed is an unsigned 32-bit number.
CREATE TABLE judge_runs (
	problem_id text NOT NULL REFERENCES problems (id),
	seed bigint NOT NULL,
	created_at text NOT NULL,
	PRIMARY KEY (problem_id, seed)
);

-- A submission of a contest problem. error is the message of an error verdict, else NULL.
CREATE TABLE contest_submissions (
	id text PRIMARY KEY,
	contest_id text NOT NULL REFERENCES contests (id),
	position bigint NOT NULL,
	verdict text NOT NULL CHECK (verdict IN ('accepted', 'wrong', 'timeout', 'error')),
	failed_test bigint,
	error text,
	code text NOT NULL,
	created_at text NOT NULL
);

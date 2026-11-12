-- The own tests of a problem. Both roles add them. An own test keeps only its input: the reference solution gives the correct answer.
CREATE TABLE problem_tests (
	id text PRIMARY KEY,
	problem_id text NOT NULL REFERENCES problems (id),
	input text NOT NULL,
	created_at text NOT NULL
);

CREATE TABLE reviews (
	id text PRIMARY KEY,
	problem_id text NOT NULL REFERENCES problems (id),
	result text NOT NULL CHECK (result IN ('needs_practice', 'with_help', 'independent')),
	note text NOT NULL,
	reviewed_on text NOT NULL,
	next_review_date text,
	created_at text NOT NULL
);

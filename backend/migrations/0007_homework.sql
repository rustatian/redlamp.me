-- Homework is a titled collection of homework tasks.
CREATE TABLE homework (
	id text PRIMARY KEY,
	title text NOT NULL,
	instructions text NOT NULL,
	due_date text,
	version bigint NOT NULL DEFAULT 1,
	created_at text NOT NULL,
	updated_at text NOT NULL
);

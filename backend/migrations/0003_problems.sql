-- A timestamp column is text in the ISO 8601 form of the API, for example 2026-10-07T12:34:56.789Z.
-- A date column is text in the form 2026-10-07. topics is a JSON array of names.
CREATE TABLE problems (
	id text PRIMARY KEY,
	slug text NOT NULL UNIQUE,
	number bigint,
	title text NOT NULL,
	difficulty text NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
	topics text NOT NULL CHECK (topics IS JSON),
	summary text NOT NULL,
	search_text text NOT NULL,
	understanding text NOT NULL DEFAULT 'needs_practice' CHECK (understanding IN ('needs_practice', 'with_help', 'independent')),
	next_review_date text,
	archived_at text,
	version bigint NOT NULL DEFAULT 1,
	created_at text NOT NULL,
	updated_at text NOT NULL,
	progress_version bigint NOT NULL DEFAULT 1
);

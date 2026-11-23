-- A contest is a timed set of problems. ends_at is started_at plus 90 minutes.
CREATE TABLE contests (
	id text PRIMARY KEY,
	started_at text NOT NULL,
	ends_at text NOT NULL
);

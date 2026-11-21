-- The problems of a contest. The title, the rating, the difficulty and the points are a copy of the bank entry at the start.
CREATE TABLE contest_problems (
	contest_id text NOT NULL REFERENCES contests (id),
	position bigint NOT NULL CHECK (position BETWEEN 1 AND 4),
	slug text NOT NULL,
	title text NOT NULL,
	rating bigint NOT NULL,
	topic text NOT NULL,
	difficulty text NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
	points bigint NOT NULL,
	code text NOT NULL DEFAULT '',
	PRIMARY KEY (contest_id, position)
);

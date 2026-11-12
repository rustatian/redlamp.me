-- One quiz attempt answers all cards of one deck. The deck is repository content, so it has no foreign key.
CREATE TABLE quiz_attempts (
	id text PRIMARY KEY,
	deck text NOT NULL,
	correct bigint NOT NULL,
	total bigint NOT NULL CHECK (total > 0 AND correct BETWEEN 0 AND total),
	passed bigint NOT NULL CHECK (passed IN (0, 1)),
	created_at text NOT NULL
);

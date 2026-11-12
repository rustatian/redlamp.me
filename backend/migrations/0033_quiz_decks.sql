-- One row for each quiz deck that has an attempt. A quiz attempt writes the row of its deck, so two attempts of one deck
-- at the same time conflict, and the one-hour pause holds. D1 has no such table: it runs one request at a time.
CREATE TABLE quiz_decks (
	deck text PRIMARY KEY,
	last_attempt_at text NOT NULL
);

-- The guard row of one active contest. Each start writes it, so two starts at the same time conflict, as quiz_decks does.
CREATE TABLE contest_starts (
	id bigint PRIMARY KEY CHECK (id = 1),
	started_at text NOT NULL
);

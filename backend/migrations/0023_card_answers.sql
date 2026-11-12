-- Quiz answers. A card is repository content, so card_id has no foreign key.
-- The newest row of a card is its schedule. interval_days is the interval that gave next_date.
-- A card has one row for each day.
CREATE TABLE card_answers (
	id text PRIMARY KEY,
	card_id text NOT NULL,
	correct bigint NOT NULL CHECK (correct IN (0, 1)),
	interval_days bigint NOT NULL CHECK (interval_days IN (1, 3, 7, 14, 30)),
	answered_on text NOT NULL,
	next_date text NOT NULL,
	created_at text NOT NULL,
	UNIQUE (card_id, answered_on)
);

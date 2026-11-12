-- The awards that the student has. The level badge follows the current XP, so it has no row.
CREATE TABLE earned_badges (
	id text PRIMARY KEY CHECK (id <> 'level-five')
);

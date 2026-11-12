-- Links from homework to learning sets. A set is repository content, so set_slug has no foreign key.
CREATE TABLE homework_sets (
	homework_id text NOT NULL REFERENCES homework (id),
	set_slug text NOT NULL,
	PRIMARY KEY (homework_id, set_slug)
);

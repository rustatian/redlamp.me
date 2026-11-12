-- Links from problems to learning sets. A set is repository content, so set_slug has no foreign key.
CREATE TABLE problem_sets (
	problem_id text NOT NULL REFERENCES problems (id),
	set_slug text NOT NULL,
	PRIMARY KEY (problem_id, set_slug)
);

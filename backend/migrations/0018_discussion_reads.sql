CREATE TABLE discussion_reads (
	role text NOT NULL CHECK (role IN ('parent', 'student')),
	problem_id text NOT NULL REFERENCES problems (id),
	read_at text NOT NULL,
	PRIMARY KEY (role, problem_id)
);

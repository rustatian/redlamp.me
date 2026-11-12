CREATE TABLE feedback (
	id text PRIMARY KEY,
	problem_id text NOT NULL REFERENCES problems (id),
	task_id text REFERENCES homework_tasks (id),
	submission_id text REFERENCES submissions (id),
	author text NOT NULL CHECK (author IN ('parent', 'student')),
	kind text NOT NULL CHECK (kind IN ('reply', 'changes_requested', 'completed')),
	body text NOT NULL,
	created_at text NOT NULL
);

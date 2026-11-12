-- AI tutor chat. A thread is one context: a problem, a learning set or a lesson.
-- context_ref is a problem ID, a set slug or a lesson ID. A set is repository content, so the column has no foreign key.
CREATE TABLE tutor_messages (
	id text PRIMARY KEY,
	context_kind text NOT NULL CHECK (context_kind IN ('problem', 'set', 'lesson')),
	context_ref text NOT NULL,
	author text NOT NULL CHECK (author IN ('student', 'assistant')),
	body text NOT NULL,
	quote text,
	created_at text NOT NULL
);

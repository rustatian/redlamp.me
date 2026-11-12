-- A homework task assigns one problem, chapter or quiz deck. A chapter or a quiz deck is repository content, so set_slug has no foreign key.
-- The foreign key of current_submission_id is in a later file, because submissions refers to this table.
CREATE TABLE homework_tasks (
	id text PRIMARY KEY,
	homework_id text NOT NULL REFERENCES homework (id),
	problem_id text REFERENCES problems (id),
	state text NOT NULL CHECK (state IN ('assigned', 'in_progress', 'submitted', 'changes_requested', 'completed', 'cancelled')),
	current_submission_id text,
	version bigint NOT NULL DEFAULT 1,
	created_at text NOT NULL,
	updated_at text NOT NULL,
	kind text NOT NULL DEFAULT 'problem' CHECK (kind IN ('problem', 'chapter', 'quiz')),
	set_slug text,
	chapter bigint,
	completed_at text,
	CHECK ((kind = 'problem' AND problem_id IS NOT NULL AND set_slug IS NULL AND chapter IS NULL)
		OR (kind IN ('chapter', 'quiz') AND problem_id IS NULL AND set_slug IS NOT NULL AND chapter IS NOT NULL AND chapter > 0 AND current_submission_id IS NULL AND state IN ('assigned', 'in_progress', 'completed', 'cancelled')))
);

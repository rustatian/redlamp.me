-- Each chapter and each quiz deck has at most one active homework task.
CREATE UNIQUE INDEX ASYNC one_active_content ON homework_tasks (kind, set_slug, chapter) WHERE state IN ('assigned', 'in_progress');

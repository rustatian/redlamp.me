-- Each problem has at most one active homework task.
CREATE UNIQUE INDEX ASYNC one_active_task ON homework_tasks (problem_id) WHERE state IN ('assigned', 'in_progress', 'submitted', 'changes_requested');

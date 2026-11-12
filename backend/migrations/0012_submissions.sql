-- A submission is a saved attempt submitted for one version of a homework task.
CREATE TABLE submissions (
	id text PRIMARY KEY,
	task_id text NOT NULL REFERENCES homework_tasks (id),
	attempt_id text NOT NULL REFERENCES attempts (id),
	task_version bigint NOT NULL,
	created_at text NOT NULL,
	UNIQUE (task_id, task_version)
);

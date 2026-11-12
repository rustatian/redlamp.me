-- Each problem has at most one draft.
CREATE UNIQUE INDEX ASYNC one_draft_per_problem ON attempts (problem_id) WHERE state = 'draft';

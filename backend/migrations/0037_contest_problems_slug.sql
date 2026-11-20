-- A bank problem is in at most one contest.
CREATE UNIQUE INDEX ASYNC contest_problems_slug ON contest_problems (slug);

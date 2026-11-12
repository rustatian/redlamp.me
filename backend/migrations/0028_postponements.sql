-- The student moves the due date of a homework by one day. The primary key allows one postponement for each homework.
CREATE TABLE postponements (
	homework_id text PRIMARY KEY REFERENCES homework (id),
	from_date text NOT NULL,
	to_date text NOT NULL,
	created_at text NOT NULL
);

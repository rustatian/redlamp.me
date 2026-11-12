-- Aurora DSQL adds a foreign key to an existing table only with NOT VALID. New and changed rows are still checked.
-- https://docs.aws.amazon.com/aurora-dsql/latest/userguide/alter-table-syntax-support.html
ALTER TABLE homework_tasks ADD CONSTRAINT homework_tasks_current_submission_id_fkey FOREIGN KEY (current_submission_id) REFERENCES submissions (id) NOT VALID;

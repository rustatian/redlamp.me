-- A sign-in session. created_at and expires_at are milliseconds since the Unix epoch.
CREATE TABLE sessions (
	token_hash text PRIMARY KEY,
	role text NOT NULL CHECK (role IN ('parent', 'student')),
	credential_fingerprint text NOT NULL,
	created_at bigint NOT NULL,
	expires_at bigint NOT NULL
);

-- document is JSON text of at most 256000 bytes. The column is text, so a read returns the stored bytes.
CREATE TABLE attempts (
	id text PRIMARY KEY,
	problem_id text NOT NULL REFERENCES problems (id),
	state text NOT NULL CHECK (state IN ('draft', 'saved')),
	document text NOT NULL CHECK (document IS JSON AND octet_length(document) <= 256000),
	search_text text NOT NULL,
	acceptance text NOT NULL CHECK (acceptance IN ('not_submitted', 'not_accepted', 'accepted')),
	understanding text NOT NULL CHECK (understanding IN ('needs_practice', 'with_help', 'independent')),
	version bigint NOT NULL DEFAULT 1,
	created_at text NOT NULL,
	updated_at text NOT NULL,
	saved_at text
);

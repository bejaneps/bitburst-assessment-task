CREATE TABLE IF NOT EXISTS bitburst."objects" (
	o_id INTEGER PRIMARY KEY NOT NULL,
	online BOOLEAN NOT NULL DEFAULT TRUE,
	last_seen TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

-- create an index on o_id column for faster reads
CREATE INDEX o_id_idx ON bitburst."objects" (o_id);
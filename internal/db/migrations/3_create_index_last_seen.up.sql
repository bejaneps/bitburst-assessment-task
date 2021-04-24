-- create an index on last_seen column for faster deletes
CREATE INDEX IF NOT EXISTS last_seen_idx ON bitburst."objects" (last_seen);
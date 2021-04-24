-- create an index on last_seen column for faster deletes
DELETE INDEX IF EXISTS last_seen_idx ON bitburst."objects" (last_seen);
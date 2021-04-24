-- name: InsertObjectsOrUpdate :many
INSERT INTO bitburst."objects" ( o_id, online )
VALUES ( UNNEST($1::INT[]), UNNEST($2::BOOLEAN[]) ) ON CONFLICT ( o_id ) DO
UPDATE
	SET last_seen = CURRENT_TIMESTAMP,
		online = EXCLUDED.online
	RETURNING o_id;

-- name: DeleteNotSeenObjects :many
DELETE
FROM
	bitburst."objects"
WHERE
	last_seen < CURRENT_TIMESTAMP - INTERVAL '30 seconds' RETURNING o_id;
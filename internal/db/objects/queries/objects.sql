-- name: InsertObjectsOrUpdate :many
INSERT INTO bitburst."objects" ( o_id )
VALUES ( UNNEST($1::INT[]) ) ON CONFLICT ( o_id ) DO
UPDATE
	SET last_seen = CURRENT_TIMESTAMP,
		online = true
RETURNING o_id;

-- name: UpdateObjects :many
UPDATE bitburst."objects"
	SET online = false
WHERE
	o_id=ANY($1::INT[])
RETURNING o_id;

-- name: DeleteNotSeenObjects :many
DELETE
FROM
	bitburst."objects"
WHERE
	last_seen < CURRENT_TIMESTAMP - INTERVAL '30 seconds' RETURNING o_id;
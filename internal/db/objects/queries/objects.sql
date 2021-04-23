-- name: InsertObjectOrUpdate :one
INSERT INTO bitburst."objects" ( o_id )
VALUES
	( $1 ) ON CONFLICT ( o_id ) DO
UPDATE
	SET last_seen = CURRENT_TIMESTAMP,
		online = TRUE
	RETURNING o_id;

-- name: InsertObjectsOrUpdate :many
INSERT INTO bitburst."objects" ( o_id ) 
VALUES ( UNNEST(@o_ids::INT[]) ) ON CONFLICT ( o_id ) DO
UPDATE
	SET last_seen = CURRENT_TIMESTAMP,
		online = TRUE
	RETURNING o_id;

-- name: UpdateObject :one
UPDATE bitburst."objects"
	SET online = FALSE
	WHERE o_id = $1
	RETURNING o_id;

-- name: DeleteNotSeenObjects :many
DELETE
FROM
	bitburst."objects"
WHERE
	last_seen < CURRENT_TIMESTAMP - INTERVAL '30 seconds' RETURNING o_id;
-- name: CreateFeedFollows :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
returning *;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows where user_id = $1;

-- name: DeleteFeedFollows :exec
delete from feed_follows where id = $1 and user_id = $2;

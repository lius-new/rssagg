-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
SELECT * FROM feeds;

--  NULLS FIRST 表示nulls的排序(nulls排在非nulls前面)

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY  last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetchd :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;

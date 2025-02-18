-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name as user_name FROM feeds f
INNER JOIN users u ON f.user_id = u.id;

-- name: FindFeedFromURL :one
SELECT id, name, url FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT * from feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;

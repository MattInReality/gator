-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT u.name as user_name, f.name as feed_name, iff.*
FROM inserted_feed_follow iff
INNER JOIN users u ON iff.user_id = u.id
INNER JOIN feeds f ON iff.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
SELECT u.name as user_name, f.name as feed_name, ff.*
FROM feed_follows ff
INNER JOIN users u ON ff.user_id = u.id
INNER JOIN feeds f ON ff.feed_id = f.id
WHERE u.name = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows USING feeds
WHERE url = $1 AND feed_follows.user_id = $2;



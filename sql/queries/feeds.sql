-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE name = $1;

-- name: ResetFeeds :exec
TRUNCATE feeds;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, feeds.user_id, users.name
FROM feeds
INNER JOIN users
  ON feeds.user_id = users.id;

-- name: CreateFeedFollow :exec
INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW());



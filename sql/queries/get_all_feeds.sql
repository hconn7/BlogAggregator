-- name: GetFeedFollowsByUser :many
SELECT id, feed_id, user_id, created_at, updated_at
FROM feed_follows
WHERE user_id = $1;


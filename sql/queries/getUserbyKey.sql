-- name: GetUserByApiKey :one
SELECT id, created_at, updated_at, name, apiKey
FROM users
WHERE apiKey = $1;

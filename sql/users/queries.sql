-- name: CreateUser :one
INSERT INTO users (username, email, pw_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, pw_hash, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT id, username, email, pw_hash, created_at, updated_at
FROM users
WHERE username = $1;

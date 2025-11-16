-- name: CreateUser :one
INSERT INTO users (username, email, pw_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, pw_hash, created_at, updated_at;

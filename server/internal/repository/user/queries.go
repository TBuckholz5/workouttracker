package user

const createUser = `INSERT INTO users (username, email, pw_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, pw_hash, created_at, updated_at
`

const getUserByUsername = `SELECT id, username, email, pw_hash, created_at, updated_at
FROM users
WHERE username = $1
`

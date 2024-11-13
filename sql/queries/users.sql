-- name: SetUser :one
INSERT INTO users(id, created_at, updated_at, username, firstname, lastname, email, password, avatar_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CheckUserExists :one
SELECT EXISTS(SELECT true FROM users WHERE id = $1);


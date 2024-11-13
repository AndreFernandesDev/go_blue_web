-- name: SetAccount :one
INSERT INTO accounts(id, user_id, created_at, updated_at, provider, provider_id, access_token, refresh_token, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE provider_id = $1;


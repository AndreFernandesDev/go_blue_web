-- +goose Up
CREATE TABLE accounts (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	provider TEXT NOT NULL,
	provider_id TEXT NOT NULL UNIQUE,
	access_token TEXT NOT NULL,
	refresh_token TEXT NOT NULL,
	expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE accounts;

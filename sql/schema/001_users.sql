-- +goose Up
CREATE TABLE users (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	username TEXT NOT NULL,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL,
	avatar_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;

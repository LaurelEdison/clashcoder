-- +goose Up

CREATE TABLE users (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	name TEXT NOT NULL,	
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	last_login_at TIMESTAMP,
	CONSTRAINT email_format CHECK (email ~* '^[^@]+@[^@]+\.[^@]+$'),
	role TEXT NOT NULL DEFAULT 'user'

);

-- +goose Down

DROP TABLE users;

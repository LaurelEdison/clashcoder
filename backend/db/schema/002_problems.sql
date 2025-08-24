-- +goose Up

CREATE TABLE problems (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	difficulty TEXT NOT NULL check (difficulty in ('easy', 'medium', 'hard')),
	time_limit INT NOT NULL DEFAULT 2000,
	memory_limit_mb INT NOT NULL DEFAULT 256
);

-- +goose Down

DROP TABLE problems;

-- +goose Up

CREATE TABLE problems (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	difficulty TEXT NOT NULL check (difficulty in ('easy', 'medium', 'hard')),
	starter_code TEXT NOT NULL,
	time_limit INT NOT NULL DEFAULT 2000,
	memory_limit_mb INT NOT NULL DEFAULT 256
);

CREATE TABLE problem_tests(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
	test_code TEXT NOT NULL
);

-- +goose Down

DROP TABLE problem_tests;
DROP TABLE problems;

-- +goose Up

CREATE TABLE submissions (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
	code TEXT NOT NULL,
	language TEXT NOT NULL, --like 'go', 'cpp', 'python'
	status TEXT NOT NULL CHECK (
		status in (
			'running',
			'pending',
			'accepted',
			'wrong_answer',
			'runtime_error',
			'time_limit_exceed',
			'memory_limit_exceed'
		)
	),
	runtime_ms INT,
	memory_kb INT,
	output TEXT
);

-- +goose Down

DROP TABLE submissions;

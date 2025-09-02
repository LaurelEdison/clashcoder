-- +goose Up

CREATE TABLE lobbies(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	invite_code TEXT UNIQUE NOT NULL,
	max_users INT NOT NULL DEFAULT 2,
	started_at TIMESTAMP,
	ended_at TIMESTAMP,
	status TEXT NOT NULL DEFAULT 'waiting' CHECK (
		status in (
			'waiting',
			'in_progress',
			'finished'
		)
	),
	ready_state BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down

DROP TABLE lobbies;

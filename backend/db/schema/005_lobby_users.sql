-- +goose Up

CREATE TABLE lobby_users(
	lobby_id UUID NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
	role TEXT NOT NULL DEFAULT 'player' CHECK(role in ('player', 'host', 'spectator')), --Spectator optional for future
	PRIMARY KEY (lobby_id, user_id)
);

CREATE UNIQUE INDEX one_host_per_lobby
ON lobby_users(lobby_id)
WHERE role = 'host'

-- +goose Down

DROP TABLE lobby_users;

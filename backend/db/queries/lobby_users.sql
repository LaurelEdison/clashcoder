-- name: CreateLobbyUser :one
INSERT INTO lobby_users(lobby_id, user_id, joined_at, role)
VALUES($1,$2,$3,$4)
RETURNING *;

-- name: GetLobbyUsersByLobbyID :many
SELECT * FROM lobby_users WHERE lobby_id = $1;

-- name: GetHostFromLobbyID :one
SELECT * FROM lobby_users WHERE lobby_id = $1 AND role = 'host';

-- name: RemoveLobbyUserFromLobby :exec
DELETE FROM lobby_users WHERE lobby_id = $1 AND user_id = $2;

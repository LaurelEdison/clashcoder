-- name: CreateLobby :one
INSERT INTO lobbies(id, created_at, updated_at, invite_code, max_users, started_at, ended_at, status,ready_state)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;
-- name: GetLobbyById :one
SELECT * FROM lobbies WHERE id = $1;


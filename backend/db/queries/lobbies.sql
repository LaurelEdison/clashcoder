-- name: CreateLobby :one
INSERT INTO lobbies(id, created_at, updated_at, invite_code, max_users, started_at, ended_at, status,ready_state)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;
-- name: GetLobbyById :one
SELECT * FROM lobbies WHERE id = $1;
-- name: UpdateLobbyStatus :exec
UPDATE lobbies
SET status = $2
WHERE id = $1;
-- name: SelectProblem :exec
UPDATE lobbies
SET problem_id = $2
WHERE id = $1;
-- name: UpdateLobbyStartEnd :exec
UPDATE lobbies
SET started_at = $2,
ended_at = $3
WHERE id = $1;
-- name: StartLobby :exec
UPDATE lobbies
SET problem_id = $2,
status = $3,
started_at = $4,
ended_at = $5,
updated_at = $6
WHERE id = $1;


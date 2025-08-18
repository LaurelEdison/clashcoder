-- name: CreateProblem :one

INSERT INTO problems (id, created_at, updated_at, title, description, difficulty, time_limit, memory_limit_mb)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetProblemByID :one

SELECT * FROM problems
WHERE (id = $1);

-- name: GetAllProblems :many

SELECT * FROM problems;

-- name: GetRandomProblem :one
SELECT * FROM problems ORDER BY RANDOM() LIMIT 1;


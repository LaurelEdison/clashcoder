-- name: CreateProblemTests :one
INSERT INTO problem_tests(id, created_at, updated_at, problem_id, test_code)
VALUES($1,$2,$3,$4,$5)
RETURNING *;
-- name: GetProblemTestsByProblemID :many
SELECT * FROM problem_tests WHERE problem_id = $1;
-- name: GetLatestProblemTestByProblemID :one
SELECT * FROM problem_tests 
WHERE problem_id = $1
ORDER BY created_at DESC LIMIT 1;
-- name: GetProblemTestByID :one
SELECT * FROM problem_tests WHERE id= $1;

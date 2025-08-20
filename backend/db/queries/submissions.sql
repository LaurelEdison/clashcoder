-- name: CreateSubmission :one

INSERT INTO submissions (id, created_at, user_id, problem_id, code, language, status, runtime_ms, memory_kb, output)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

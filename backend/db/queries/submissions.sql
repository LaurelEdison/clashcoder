-- name: CreateSubmission :one

INSERT INTO submissions (id, created_at, user_id, problem_id, code, language, status, runtime_ms, memory_kb, output)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetSubmissionByID :one
SELECT * FROM submissions 
WHERE (id = $1);

-- name: GetLatestSubmissionByUserID :one
SELECT * FROM submissions 
WHERE user_id = $1 AND problem_id = $2
ORDER BY created_at DESC LIMIT 1;

-- name: GetAllSubmissionByUserID :many
SELECT * FROM submissions 
WHERE user_id = $1 AND problem_id = $2
ORDER BY created_at;

-- name: SelectPendingSubmission :one

UPDATE submissions
SET status = 'running'
WHERE id = (
	SELECT id from submissions
	WHERE status = 'pending'
	ORDER BY created_at
	LIMIT 1
	FOR UPDATE SKIP LOCKED
)
	RETURNING *;

-- name: UpdateSubmissionResult :exec
UPDATE submissions
SET status = $2,
output = $3
WHERE id = $1;

-- name: UpdateSubmissionStatus :exec
UPDATE submissions
SET status = $2
WHERE id = $1;

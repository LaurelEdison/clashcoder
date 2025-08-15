-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, password_hash, last_login_at)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

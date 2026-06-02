-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password, balance)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;
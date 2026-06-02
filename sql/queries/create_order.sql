-- name: CreateOrder :one
INSERT INTO orders (order_id, user_id, side, type, price, quantity, remaining_quantity, status, created_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;
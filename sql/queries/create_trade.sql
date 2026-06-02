-- name: CreateTrade :one
INSERT INTO trades (trade_id, maker_order_id, taker_order_id, maker_user_id, taker_user_id, price, quantity, created_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) 
RETURNING *;
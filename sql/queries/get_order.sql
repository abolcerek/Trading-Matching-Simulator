-- name: GetOrder :one
SELECT * FROM orders
WHERE order_id = $1 LIMIT 1;
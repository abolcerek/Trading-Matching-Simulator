-- name: UpdateOrder :exec
UPDATE orders
SET remaining_quantity = $1, status = $2
WHERE order_id = $3;
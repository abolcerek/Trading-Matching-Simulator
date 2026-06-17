-- name: Replay :many
SELECT * FROM orders 
WHERE status IN ('open', 'partially_filled')
ORDER BY sequence_num;
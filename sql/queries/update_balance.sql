-- name: UpdateBalance :exec
UPDATE users
SET balance = balance + $1
WHERE id = $2;
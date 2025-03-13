-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY created_at;

-- name: SelectTransfer :one
SELECT * FROM transfers
where id = $1;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;

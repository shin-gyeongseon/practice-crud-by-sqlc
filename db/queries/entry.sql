-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY created_at;

-- name: SelectEntry :one
SELECT * FROM entries
where id = $1;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;
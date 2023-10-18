-- name: CreateAccount :one
INSERT INTO accounts (
  owner, 
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY created_at;

-- name: SelectAccount :one
SELECT * FROM accounts
WHERE id = $1;

-- name: SelectAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: UpdateAccount :one
UPDATE accounts
  set balance = $2
WHERE id = $1
RETURNING *;
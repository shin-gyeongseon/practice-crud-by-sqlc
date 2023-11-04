// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: entry.sql

package tutorial

import (
	"context"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID int64
	Amount    int64
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEntry = `-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteEntry, id)
	return err
}

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
ORDER BY created_at
`

func (q *Queries) ListEntries(ctx context.Context) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectEntry = `-- name: SelectEntry :one
SELECT id, account_id, amount, created_at FROM entries
where id = $1
`

func (q *Queries) SelectEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, selectEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

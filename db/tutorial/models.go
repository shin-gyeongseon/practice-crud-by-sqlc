// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package tutorial

import (
	"database/sql"
)

type Account struct {
	ID        int64
	Owner     interface{}
	Balance   int64
	Currency  string
	CreatedAt interface{}
}

type Entry struct {
	ID        int64
	AccountID sql.NullInt64
	// can be negative our?
	Amount    int64
	CreatedAt interface{}
}

type Transfer struct {
	ID            int64
	FromAccountID sql.NullInt64
	ToAccountID   sql.NullInt64
	// must be positive
	Amount    int64
	CreatedAt interface{}
}

type User struct {
	Username          string
	HashedPassword    string
	FullName          interface{}
	Email             string
	PasswordChangedAt interface{}
	CreatedAt         interface{}
}

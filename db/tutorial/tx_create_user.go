package tutorial

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
}

type CreateUserTxResponse struct {
	User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResponse, error) {
	var result CreateUserTxResponse

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User = user
		return nil
	})

	return result, err
}

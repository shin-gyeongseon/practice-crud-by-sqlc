package tutorial

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	fromAccountID := 1
	toAccountID := 2
	amount := 10

	beforeFromAcocunt, err := testSqlStore.SelectAccount(context.Background(), int64(fromAccountID))
	require.NoError(t,err)
	beforeToAccount, err := testSqlStore.SelectAccount(context.Background(), int64(toAccountID))
	require.NoError(t,err)


	result, err := testSqlStore.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: int64(fromAccountID),
		ToAccountID:   int64(toAccountID),
		Amount:        int64(amount),
	})

	require.NoError(t, err)

	require.Equal(t, int64(fromAccountID), result.FromAccount.ID)
	require.Equal(t, int64(toAccountID), result.ToAccount.ID)

	require.Equal(t, beforeFromAcocunt.Balance - int64(amount), result.FromAccount.Balance)
	require.Equal(t, beforeToAccount.Balance + int64(amount), result.ToAccount.Balance)
}

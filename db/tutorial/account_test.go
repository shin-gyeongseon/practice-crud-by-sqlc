package tutorial

import (
	"context"
	"go-practice/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney().Int64,
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    arg.Owner,
		Balance:  arg.Balance,
		Currency: arg.Currency,
	})

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestListAccount(t *testing.T) {
	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  10,
		Offset: 1 * 5,
	})

	for _, v := range accounts {
		t.Log(v)
	}

	require.NoError(t, err)
}

func TestDelete(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	newaccount := createRandomAccount(t)
	account, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      newaccount.ID,
		Balance: newaccount.Balance,
	})

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, newaccount.Balance, account.Balance)
}

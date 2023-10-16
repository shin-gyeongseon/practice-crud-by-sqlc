package tutorial

import (
	"context"
	"go-practice/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount() CreateAccountParams {
	return CreateAccountParams{
		Owner:    util.RnadomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func TestCreateAccount(t *testing.T) Account {
	arg := createRandomAccount()
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	
	return account
}

func TestListAccount(t *testing.T) {
	accounts, err := testQueries.ListAccounts(context.Background())

	for _, v := range accounts {
		t.Log(v)
	}

	require.NoError(t, err)
	require.Greater(t, len(accounts), 0)
}

func TestDelete(t *testing.T) {
	account := TestCreateAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	// get account 가 있어야지 조회를 해서 비교할 수 있네 
	// err := testQueries.DeleteAccount(context.Background(), account.ID)
	// require.Error(t, err)
}

func TestUpdate(t *testing.T) {
	arg := UpdateAccountParams{
		ID:      2,
		Owner:   "updated-owner",
		Balance: 10,
	}
	account, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
}

package tutorial

import (
	"context"
	"go-practice/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)

	return user
}

func TestSelectUser(t *testing.T) {
	user := CreateRandomUser(t)
	selectUser, err := testQueries.SelectUser(context.Background(), user.Username)
	
	require.NoError(t, err)
	require.Equal(t, user.Username, selectUser.Username)
	require.Equal(t, user.Email, selectUser.Email)
	require.Equal(t, user.FullName, selectUser.FullName)
	require.Equal(t, user.HashedPassword, selectUser.HashedPassword)

	require.NotZero(t, selectUser.CreatedAt)
	require.NotZero(t, selectUser.PasswordChangedAt)
}
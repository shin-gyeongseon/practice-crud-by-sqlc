package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	for i := 0; i < 5; i++ {
		password := RandomString(6)
		
		hashedPassword, err := HashPassword(password)
		fmt.Printf("hashedPassword: %v\n", hashedPassword)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)
		
		err = CheckPassword(password, hashedPassword)
		require.NoError(t, err)
		
		wrongPassword := RandomString(6)
		err = CheckPassword(wrongPassword, hashedPassword)
		require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	}
}

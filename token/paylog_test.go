package token

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		tokenID, err := uuid.NewRandom()
		require.NoError(t, err)
	
		fmt.Printf("tokenID: %v\n", tokenID.String())
	}
}

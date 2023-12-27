package api

import (
	"context"
	mock_tutorial "go-practice/db/mock"
	"go-practice/db/tutorial"
	"go-practice/util"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

const (
	USD = "USD"
	EUR = "EUR"
)

func TestTransfer(t *testing.T) {

	testcase := []struct {
		name      string
		body      transferRequest
		buildStub func(store *mock_tutorial.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "succeed",
			body: transferRequest{
				FromAccountID: 1,
				ToAccountID: 2,
				Amount: 10,
				Currency: USD,
			},
			buildStub: func(store *mock_tutorial.MockStore) {
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Any()).
				Times(1).
				Return(tutorial.TransferTxResult{
					FromAccount: randomAccount(),
				})
			},
		},
	}
	

}

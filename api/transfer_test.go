package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	mock_tutorial "go-practice/db/mock"
	"go-practice/db/tutorial"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	USD = "USD"
	EUR = "EUR"
)

func TestTransfer(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
	var amount int64 = 10

	testcase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_tutorial.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        USD,
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				store.EXPECT().SelectAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().SelectAccount(gomock.Any(), account2.ID).Times(1).Return(account2, nil)

				// store.EXPECT().
				// 	TransferTx(context.Background(), gomock.Any()).
				// 	Times(1).
				// 	Return(tutorial.TransferTxResult{
				// 		FromAccount: account1,
				// 		ToAccount:   account2,
				// 		// from, to entity, transfer 없는 상태로 합시다. 크게 상관없으니까.
				// 	}, nil)

				arg := tutorial.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "bad_request",
			body: gin.H{
				"from_account_id": account1.Currency,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        USD,
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				store.EXPECT().SelectAccount(gomock.Any(), account1.ID).Times(0)
				store.EXPECT().SelectAccount(gomock.Any(), account2.ID).Times(0)

				arg := tutorial.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "validation internal server error",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        USD,
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				store.EXPECT().SelectAccount(gomock.Any(), account1.ID).Times(1).Return(account1, errors.New("temporary error"))
				store.EXPECT().SelectAccount(gomock.Any(), account2.ID).Times(0)

				arg := tutorial.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "validation not find row error",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        USD,
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				store.EXPECT().SelectAccount(gomock.Any(), account1.ID).Times(1).Return(account1, sql.ErrNoRows)
				store.EXPECT().SelectAccount(gomock.Any(), account2.ID).Times(0)

				arg := tutorial.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testcase {
		tc := testcase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_tutorial.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store, TestGlobalConfig)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}

}

package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mock_tutorial "go-practice/db/mock"
	"go-practice/db/tutorial"
	"go-practice/util"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      tutorial.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(tutorial.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserParamMatcher(arg tutorial.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		arg:      arg,
		password: password,
	}
}

func TestCreateUserAPI(t *testing.T) {
	randomUser, password := randomUser(t)

	testcase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_tutorial.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "CREATED",
			body: gin.H{
				"user_name": randomUser.Username,
				"full_name": randomUser.FullName,
				"email":     randomUser.Email,
				"password":  password,
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				arg := tutorial.CreateUserParams{
					Username: randomUser.Username,
					FullName: randomUser.FullName,
					Email:    randomUser.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParamMatcher(arg, password)).
					Times(1).
					Return(randomUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BADREQUEST",
			body: gin.H{
				"user_name": randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     "this is not a email format",
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, 400, http.StatusBadRequest)
			},
		},
		{
			name: "INTERNAL ERROR",
			body: gin.H{
				"user_name": randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(tutorial.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DUPLICATE USERNAME",
			body: gin.H{
				"user_name": randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mock_tutorial.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(tutorial.User{}, tutorial.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomUser(t *testing.T) (user tutorial.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = tutorial.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

package gapi

import (
	"context"
	"errors"
	mock_tutorial "go-practice/db/mock"
	"go-practice/db/tutorial"
	"go-practice/pb"
	"go-practice/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func randomUserFn() (user tutorial.User, password string) {
	password = util.RandomString(10)
	hasedPassword, _ := util.HashPassword(password)

	user = tutorial.User{
		Username:          util.RandomString(6),
		HashedPassword:    hasedPassword,
		FullName:          util.RandomString(10),
		Email:             util.RandomEmail(),
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}
	return
}

func TestRpcCreateUser(t *testing.T) {
	randomUser, pwd := randomUserFn()

	testcases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStub     func(store *mock_tutorial.MockStore)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "create",
			req: &pb.CreateUserRequest{
				Username: randomUser.Username,
				FullName: randomUser.FullName,
				Email:    randomUser.Email,
				Password: pwd,
			},
			buildStub: func(store *mock_tutorial.MockStore) {
				resultResponse := &pb.CreateUserResponse{
					User: &pb.User{
						Username:          randomUser.Username,
						FullName:          randomUser.FullName,
						Email:             randomUser.Email,
						PasswordChangedAt: timestamppb.New(randomUser.PasswordChangedAt),
						CreatedAt:         timestamppb.New(randomUser.CreatedAt),
					},
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(resultResponse, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res.GetUser()
				require.Equal(t, randomUser.FullName, createdUser.Username)
				require.Equal(t, randomUser.Username, createdUser.Username)
				require.Equal(t, randomUser.Email, createdUser.Email)
			},
		},
		{
			name: "duplicated ids",
			req: &pb.CreateUserRequest{
				Username: randomUser.Username,
				FullName: randomUser.FullName,
				Email:    randomUser.Email,
				Password: pwd,
			},
			buildStub: func(store *mock_tutorial.MockStore) {
				responseError := tutorial.ErrUniqueViolation

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, responseError)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			name: "internal error",
			req: &pb.CreateUserRequest{
				Username: randomUser.Username,
				FullName: randomUser.FullName,
				Email:    randomUser.Email,
				Password: pwd,
			},
			buildStub: func(store *mock_tutorial.MockStore) {
				store.EXPECT().
				CreateUserTx(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, errors.New("Internal any error"))
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)	
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for _, v := range testcases {
		/*
			파라미터를 설정하고 파라미터를 던지고
			결과를 통해서 확인할 수 있지않을까 ?
		*/
		storeController := gomock.NewController(t)
		defer storeController.Finish()
		mockStore := mock_tutorial.NewMockStore(storeController)

		v.buildStub(mockStore)

		server := NewServer(mockStore, util.Config{}) // config 생성은 필요 없을 것 같습니다.

		cur, err := server.UnimplementedSimpleBankServer.CreateUser(context.Background(), v.req) // 여기서 PRC를 실행하는거야 그래서 여기서 error 가 나오면 안된다는거고
		v.checkResponse(t, cur, err)
	}

}

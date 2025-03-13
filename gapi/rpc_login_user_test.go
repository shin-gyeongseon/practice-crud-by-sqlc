package gapi

import (
	mock_tutorial "go-practice/db/mock"
	"go-practice/pb"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGrpcLoginUser(t *testing.T) {
	// random user ?! 안
	testcases := []struct {
		name          string
		requestParam  *pb.LoginUserRequest
		buildStub     func(store *mock_tutorial.MockStore, loginRequest *pb.LoginUserRequest) (pb.LoginUserResponse, error)
		checkResponse func(t *testing.T, loginResponse *pb.LoginUserResponse, err error)
	}{}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			mockStore := mock_tutorial.NewMockStore(controller)

			loginUserResponse, err := testcase.buildStub(mockStore, testcase.requestParam)

			testcase.checkResponse(t, &loginUserResponse, err)
		})
	}
}

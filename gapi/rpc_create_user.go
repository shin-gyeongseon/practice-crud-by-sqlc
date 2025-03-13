package gapi

import (
	"context"
	"go-practice/db/tutorial"
	"go-practice/pb"
	"go-practice/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := tutorial.CreateUserTxParams{
		CreateUserParams: tutorial.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
	}

	// create user transaction
	createUserResponse, err2 := server.store.CreateUserTx(ctx, arg)
	if err2 != nil {
		if tutorial.ErrorCode(err) == tutorial.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, err2.Error())
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err2)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(createUserResponse.User),
	}

	return rsp, nil
}

package gapi

import (
	"context"
	"go-practice/db/tutorial"
	"go-practice/pb"
	"go-practice/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.ValidUser(ctx, req.UserName)
	if err != nil {
		return nil, err
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		err = status.Errorf(codes.Unauthenticated, "invalid password please check again")
		return nil, err
	}

	token, payload, err := server.tokenMaker.CreateToken(req.UserName, server.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	loginSession := tutorial.CreateSessionParams{
		ID:           payload.ID.String(),
		ExpiresAt:    payload.ExpiredAt,
		RefreshToken: token,
		Username:     req.UserName,
		// UserAgent:    ctx.Request.UserAgent(), // gRPC에서는 이거 일단 사용 안하는걸로 합시다.
		// ClientIp:     ctx.ClientIP(),
		IsBlocked: false,
	}
	session, err := server.store.CreateSession(ctx, loginSession)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	loginResponse := &pb.LoginUserResponse{
		User:               convertUser(user),
		SessionId:          session.ID,
		Token:              token,
		AccessTokenExpired: timestamppb.New(payload.ExpiredAt),
	}

	return loginResponse, nil
}

func (server *Server) ValidUser(ctx context.Context, userName string) (tutorial.User, error) {
	user, err := server.store.SelectUser(ctx, userName)
	if err != nil {
		if tutorial.ErrorCode(err) == tutorial.ErrUniqueViolation.Code {
			return tutorial.User{}, status.Errorf(codes.AlreadyExists, "Already exist user name")
		}

		return tutorial.User{}, status.Errorf(codes.Internal, "valid user internal server errror")
	}

	return user, nil
}

package gapi

import (
	"go-practice/db/tutorial"
	"go-practice/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user tutorial.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

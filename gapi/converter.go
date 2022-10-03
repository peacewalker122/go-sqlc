package gapi

import (
	db "sqlc/db/sqlc"
	"sqlc/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convert(user db.User) *pb.User{
	return &pb.User{
		Username:       user.Username,
		Fullname:       user.FullName,
		Email:          user.Email,
		PasswordChange: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:      timestamppb.New(user.CreatedAt),
	}
}
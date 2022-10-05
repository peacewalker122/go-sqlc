package gapi

import (
	db "github.com/peacewalker122/go-sqlc/db/sqlc"
	"github.com/peacewalker122/go-sqlc/pb"

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
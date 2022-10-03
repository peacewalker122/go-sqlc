package gapi

import (
	"context"
	db "sqlc/db/sqlc"
	"sqlc/pb"
	"sqlc/util"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.ResponseUser, error) {

	pass, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hashing password due: ", err)
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: pass,
		FullName:       req.Fullname,
		Email:          req.Email,
	}

	account, err := s.store.CreateUser(c, arg)
	if err != nil {
		if PqErr, ok := err.(*pq.Error); ok {
			switch PqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exist ", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "cannot Create User due: ", err)
	}
	resp := &pb.ResponseUser{
		User: convert(account),
	}
	return resp, nil
}

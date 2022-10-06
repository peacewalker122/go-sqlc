package gapi

import (
	"context"

	db "github.com/peacewalker122/go-sqlc/db/sqlc"
	"github.com/peacewalker122/go-sqlc/pb"
	"github.com/peacewalker122/go-sqlc/util"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.ResponseUser, error) {
	validation := validateUserReq(req)
	if validation != nil{
		return nil,InvalidArgument(validation)
	}

	pass, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hashing password due: %v", err)
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
				return nil, status.Errorf(codes.AlreadyExists, "username already exist: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "cannot Create User due: %v", err)
	}
	resp := &pb.ResponseUser{
		User: convert(account),
	}
	return resp, nil
}

func validateUserReq(req *pb.CreateUserRequest) (violation []*errdetails.BadRequest_FieldViolation) {
	if err := validateUsername(req.Username); err != nil {
		violation = append(violation, fieldValidation("username", err))
	}
	if err := validatePassword(req.Password); err != nil {
		violation = append(violation, fieldValidation("password", err))
	}
	if err := validateFullname(req.Fullname); err != nil {
		violation = append(violation, fieldValidation("fullname", err))
	}
	if err := validateEmail(req.Email); err != nil {
		violation = append(violation, fieldValidation("email", err))
	}
	return violation
}

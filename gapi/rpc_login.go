package gapi

import (
	"context"
	"database/sql"

	db "github.com/peacewalker122/go-sqlc/db/sqlc"
	"github.com/peacewalker122/go-sqlc/pb"
	"github.com/peacewalker122/go-sqlc/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) Login(c context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := s.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "No Such Username")
		}
		return nil, status.Errorf(codes.Internal, "Can't Get The User Due: %v", err)
	}

	err = util.CheckPass(req.Password, res.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Password Not Same")
	}

	token, AccessPayload, err := s.TokenMaker.CreateToken(res.Username, s.config.Duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot Create Acces Token")
	}

	RefreshToken, RefreshPayload, err := s.TokenMaker.CreateToken(res.Username, s.config.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot Create Refresh Payload")
	}
	
	mtdt := s.extractMetadata(c)
	session, err := s.store.CreateSession(c, db.CreateSessionParams{
		ID:           RefreshPayload.ID,
		Username:     res.Username,
		RefreshToken: RefreshToken,
		UserAgent:    mtdt.HostName,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    RefreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot Create Session")
	}
	rsp := pb.LoginResponse{
		User:                  convert(res),
		SessionId:             session.ID.String(),
		AccessToken:           token,
		RefreshToken:          RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(AccessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(RefreshPayload.ExpiredAt),
	}
	return &rsp, nil
}

package gapi

import (
	"fmt"
	db "sqlc/db/sqlc"
	"sqlc/token"
	"sqlc/util"
	"sqlc/pb"
)

type server struct {
	pb.SimpleBankServer
	config     util.Config
	store      db.Store
	TokenMaker token.Maker
}

func Newserver(c util.Config, store db.Store) (*server, error) {
	Newtoken, err := token.NewJWTmaker(c.SymmectricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token %v", err.Error())
	}
	server := &server{
		config:     c,
		store:      store,
		TokenMaker: Newtoken,
	}
	return server, nil
}

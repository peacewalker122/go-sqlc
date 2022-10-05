package gapi

import (
	"fmt"

	db "github.com/peacewalker122/go-sqlc/db/sqlc"
	"github.com/peacewalker122/go-sqlc/pb"
	"github.com/peacewalker122/go-sqlc/token"
	"github.com/peacewalker122/go-sqlc/util"
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

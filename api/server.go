package api

import (
	"fmt"
	db "sqlc/db/sqlc"
	"sqlc/token"
	"sqlc/util"

	"github.com/gin-gonic/gin"
)

type server struct {
	config   util.Config
	router   *gin.Engine
	store    db.Store
	AccToken token.Maker
}

func Newserver(c util.Config, store db.Store) (*server, error) {
	Newtoken, err := token.NewJWTmaker(c.SymmectricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token %v", err.Error())
	}
	server := &server{
		config:   c,
		store:    store,
		AccToken: Newtoken,
	}
	server.routerhandle()
	return server, nil
}

func (server *server) routerhandle() {
	router := gin.Default()

	router.POST("/user/login",server.serverLogin)
	router.POST("/user", server.createUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getaccountid)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.Transfertx)
	server.router = router
}

func (s *server) Runserver(target string) error {
	return s.router.Run(target)
}

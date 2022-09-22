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
	TokenMaker token.Maker
}

func Newserver(c util.Config, store db.Store) (*server, error) {
	Newtoken, err := token.NewJWTmaker(c.SymmectricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token %v", err.Error())
	}
	server := &server{
		config:   c,
		store:    store,
		TokenMaker: Newtoken,
	}
	server.routerhandle()
	return server, nil
}

func (server *server) routerhandle() {
	router := gin.Default()
	
	authRouter := router.Group("/").Use(authMiddleware(server.TokenMaker))

	authRouter.POST("/user/login",server.serverLogin)
	authRouter.POST("/user", server.createUser)

	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getaccountid)
	authRouter.GET("/accounts", server.listAccount)
	
	authRouter.POST("/transfers", server.Transfertx)
	server.router = router
}

func (s *server) Runserver(target string) error {
	return s.router.Run(target)
}

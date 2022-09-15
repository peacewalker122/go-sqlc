package api

import (
	db "sqlc/db/sqlc"

	"github.com/gin-gonic/gin"
)

type server struct {
	router *gin.Engine
	store  db.Store
}

func Newserver(store db.Store) *server {
	server := &server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getaccountid)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.Transfertx)
	server.router = router
	return server
}

func (s *server) Runserver(target string) error {
	return s.router.Run(target)
}

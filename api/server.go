package api

import (
	"fmt"
	db "sqlc/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type server struct {
	router *gin.Engine
	store db.Store
}

func Newserver(store db.Store) *server {
	server := &server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getaccountid)
	router.GET("/accounts", server.listAccount)
	server.router = router
	return server
}

func (s *server) Runserver(target string) error {
	return s.router.Run(target)
}

func errorvalidator(err error) gin.H {
	ermsg := []string{}
	for _, e := range err.(validator.ValidationErrors) {
		errmsg := fmt.Sprintf("error happen in %s, due %s", e.Field(), e.Error())
		ermsg = append(ermsg, errmsg)
	}
	r := gin.H{
		"errors": ermsg,
	}
	return r
}

func errorhandle(err error) gin.H {
	r := gin.H{
		"errors": err.Error(),
	}
	return r
}

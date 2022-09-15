package api

import (
	"database/sql"
	"fmt"
	"net/http"
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
	router.POST("/transfers", server.Transfertx)
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

func(s *server) transferValidator(c *gin.Context, accountID int64, currency string) bool{
	account,err := s.store.GetAccount(c,accountID)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(http.StatusBadRequest, errorhandle(err))
			return false
		}
		c.JSON(http.StatusInternalServerError,errorhandle(err))
		return false
	}
	if account.Currency != currency{
		err := fmt.Errorf("Different Currency, expected %v",account.Currency)
		c.JSON(http.StatusBadRequest,errorhandle(err))
		return false
	}
	return true
}
package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "sqlc/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func errorvalidator(err error) gin.H {
	ermsg := []string{}
	for _, e := range err.(validator.ValidationErrors) {
		errmsg := fmt.Sprintf("error happen in %s, due %s, expected %s", e.Field(), e.Value(), e.Param())
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

func (s *server) transferValidator(c *gin.Context, FaccountID, TAccountID int64, currency string) (db.Account,db.Account, bool) {
	account, err := s.store.GetAccount(c, FaccountID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorhandle(err))
			return account,db.Account{},false
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
		return account,db.Account{},false
	}

	account2, err := s.store.GetAccount(c, TAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorhandle(err))
			return db.Account{},account2,false
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
		return db.Account{},account2,false
	}

	if account.Currency != account2.Currency {
		err := fmt.Errorf("Different Currency, expected %v. can't process %v into %v", account.Currency, account.Currency, account2.Currency)
		c.JSON(http.StatusBadRequest, errorhandle(err))
		return account,account2,false
	}
	return account,account2,true
}

// func returns(s string) gin.H {
// 	r := fmt.Errorf("errors in: ", s)
// 	return gin.H{
// 		"error": r,
// 	}
// }

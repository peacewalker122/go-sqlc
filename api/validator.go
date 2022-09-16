package api

import (
	"database/sql"
	"fmt"
	"net/http"

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

func (s *server) transferValidator(c *gin.Context, FaccountID, TAccountID int64, currency string) bool {
	account, err := s.store.GetAccount(c, FaccountID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, errorhandle(err))
			return false
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
		return false
	}

	account2, err := s.store.GetAccount(c, TAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, errorhandle(err))
			return false
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
		return false
	}

	if account.Currency != account2.Currency {
		err := fmt.Errorf("Different Currency, expected %v. can't process %v into %v", account.Currency, account.Currency, account2.Currency)
		c.JSON(http.StatusBadRequest, errorhandle(err))
		return false
	}
	return true
}

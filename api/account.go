package api

import (
	"database/sql"
	"net/http"
	db "sqlc/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createAccountParam struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR IDR"`
}

func (s *server) createAccount(c *gin.Context) {
	var req createAccountParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		r := errorvalidator(err)
		c.JSON(http.StatusBadRequest, r)
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(c, arg)
	if err != nil {
		r := errorvalidator(err)
		c.JSON(http.StatusBadRequest, r)
		return
	}
	c.JSON(http.StatusOK, account)
}

type getaccountidParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *server) getaccountid(c *gin.Context) {
	var res getaccountidParam
	err := c.ShouldBindUri(&res)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorhandle(err))
		return
	}

	getid, err := s.store.GetAccount(c, res.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, errorhandle(err))
			return
		}
		c.JSON(http.StatusBadRequest, errorhandle(err))
		return
	}

	c.JSON(http.StatusOK, getid)
}

type listAccountRequest struct {
	pageID   int32 `form:"pageid" binding:"required,min=1"`
	pageSize int32 `form:"pagesize" binding:"required,min=5,max=10"`
}

func (server *server) listAccount(ctx *gin.Context) {
    var req listAccountRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorhandle(err))
        return
    }

    arg := db.ListAccountsParams{
        Limit:  req.pageSize,
        Offset: (req.pageID - 1) * req.pageSize,
    }

    accounts, err := server.store.ListAccounts(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorhandle(err))
        return
    }

    ctx.JSON(http.StatusOK, accounts)
}


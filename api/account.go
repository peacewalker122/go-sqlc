package api

import (
	"database/sql"
	"log"
	"net/http"
	db "sqlc/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lib/pq"
)

type createAccountParam struct {
	Owner    string `json:"owner" binding:"required,alphanum"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR IDR GBP"`
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

		if PqErr,ok := err.(*pq.Error); ok{
			switch PqErr.Code.Name() {
			case "foreign_key_violation","unique_violation":
				c.JSON(http.StatusForbidden,errorhandle(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
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
			c.JSON(http.StatusNotFound, errorhandle(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
		return
	}

	c.JSON(http.StatusOK, getid)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

func (server *server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindWith(&req, binding.Query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorvalidator(err))
		return
	}

	log.Println(req.PageSize, " and ", req.PageID)

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorhandle(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

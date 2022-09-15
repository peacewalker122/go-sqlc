package api

import (
	"net/http"
	db "sqlc/db/sqlc"

	"github.com/gin-gonic/gin"
)

type TransferParam struct {
	FromAccountID int64  `json:"from_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR IDR GBP"`
}

func (s *server) Transfertx(c *gin.Context) {
	var res TransferParam
	err := c.ShouldBindJSON(&res)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorvalidator(err))
		return
	}

	if !s.transferValidator(c, res.FromAccountID, res.Currency) {
		return
	}

	if !s.transferValidator(c, res.ToAccountID, res.Currency) {
		return
	}

	arg := db.TransferctxParams{
		FromAccountID: res.FromAccountID,
		ToAccountID:   res.ToAccountID,
		Amount:        res.Amount,
	}

	getid, err := s.store.TransferCtx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorvalidator(err))
	}

	c.JSON(http.StatusOK, getid)
}

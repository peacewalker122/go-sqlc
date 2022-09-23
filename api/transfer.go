package api

import (
	"fmt"
	"net/http"
	db "sqlc/db/sqlc"
	"sqlc/token"

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
	
	Param,_,ok := s.transferValidator(c, res.FromAccountID, res.ToAccountID, res.Currency)
	if !ok {
		return
	}

	authPayload := c.MustGet(authPayload).(*token.Payload)
	if authPayload.Username != Param.Owner {
		err := fmt.Errorf("Unauthorized username for this owner")
		c.JSON(http.StatusUnauthorized,errorhandle(err))
		return
	}

	arg := db.TransferctxParams{
		FromAccountID: res.FromAccountID,
		ToAccountID:   res.ToAccountID,
		Amount:        res.Amount,
	}

	getid, err := s.store.TransferCtx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandle(err))
	}

	c.JSON(http.StatusOK, getid)
}

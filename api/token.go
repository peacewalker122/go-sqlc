package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AccesTokenParam struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccesTokenResp struct {
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (s *server) serverAccesToken(c *gin.Context) {
	var req AccesTokenParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorvalidator(err))
		return
	}

	refreshPayload, err := s.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorhandle(err))
		return
	}

	session, err := s.store.GetSession(c, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorhandle(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		c.JSON(http.StatusUnauthorized, errorhandle(err))
		return
	}

	if session.ID != refreshPayload.ID {
		err := fmt.Errorf("incorrect session user")
		c.JSON(http.StatusUnauthorized, errorhandle(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch session token")
		c.JSON(http.StatusUnauthorized, errorhandle(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		c.JSON(http.StatusUnauthorized, errorhandle(err))
		return
	}

	Accesstoken, AccessPayload, err := s.TokenMaker.CreateToken(refreshPayload.Username, s.config.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandle(err))
	}

	rsp := AccesTokenResp{
		RefreshToken:          Accesstoken,
		RefreshTokenExpiresAt: AccessPayload.ExpiredAt,
	}
	
	c.JSON(http.StatusOK, rsp)
}

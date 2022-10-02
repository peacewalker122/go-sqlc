package api

import (
	"database/sql"
	"net/http"
	db "sqlc/db/sqlc"
	"sqlc/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUserParam struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"full_name" binding:"required,alpha"`
	Email    string `json:"email" binding:"required,email"`
}

type userResp struct {
	Username       string    `json:"username"`
	Fullname       string    `json:"full_name"`
	Email          string    `json:"email"`
	PasswordChange time.Time `json:"password_change"`
	CreatedAt      time.Time `json:"created_at"`
}

func newUserResp(u db.User) userResp {
	return userResp{
		Username:       u.Username,
		Fullname:       u.FullName,
		Email:          u.Email,
		PasswordChange: u.PasswordChangedAt,
		CreatedAt:      u.CreatedAt,
	}
}

func (s *server) createUser(c *gin.Context) {
	var req createUserParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		r := errorvalidator(err)
		c.JSON(http.StatusBadRequest, r)
		return
	}

	pass, err := util.HashPassword(req.Password)
	if err != nil {
		r := errorhandle(err)
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: pass,
		FullName:       req.Fullname,
		Email:          req.Email,
	}

	account, err := s.store.CreateUser(c, arg)
	if err != nil {
		if PqErr, ok := err.(*pq.Error); ok {
			switch PqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorhandle(err))
				return
			}
		}
		r := errorhandle(err)
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	resp := newUserResp(account)

	c.JSON(http.StatusOK, resp)
}

type loginParam struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResp struct {
	SessionID             uuid.UUID `json:"session_id"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  userResp  `json:"user"`
	AccesToken            string    `json:"acc_token"`
	AccesTokenExpiresAt   time.Time `json:"acces_token_expire_sat"`
}

func (s *server) serverLogin(c *gin.Context) {
	var req loginParam
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorvalidator(err))
		return
	}

	res, err := s.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorhandle(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorhandle(err))
		return
	}

	err = util.CheckPass(req.Password, res.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorhandle(err))
		return
	}

	token, AccessPayload, err := s.TokenMaker.CreateToken(res.Username, s.config.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandle(err))
	}

	RefreshToken, RefreshPayload, err := s.TokenMaker.CreateToken(res.Username, s.config.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandle(err))
	}

	session, err := s.store.CreateSession(c, db.CreateSessionParams{
		ID:           RefreshPayload.ID,
		Username:     res.Username,
		RefreshToken: RefreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    RefreshPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandle(err))
	}

	rsp := loginResp{
		SessionID:             session.ID,
		RefreshToken:          RefreshToken,
		RefreshTokenExpiresAt: RefreshPayload.ExpiredAt,
		User:                  newUserResp(res),
		AccesToken:            token,
		AccesTokenExpiresAt:   AccessPayload.ExpiredAt,
	}
	c.JSON(http.StatusOK, rsp)
}

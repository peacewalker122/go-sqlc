package api

import (
	"database/sql"
	"net/http"
	db "sqlc/db/sqlc"
	"sqlc/util"
	"time"

	"github.com/gin-gonic/gin"
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
	Token string   `json:"acc_token"`
	User  userResp `json:"user"`
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
	}

	err = util.CheckPass(req.Password, res.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorhandle(err))
		return
	}

	token, err := s.AccToken.CreateToken(res.Username, s.config.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorhandle(err))
	}

	rsp := loginResp{
		Token: token,
		User:  newUserResp(res),
	}
	c.JSON(http.StatusOK, rsp)
}

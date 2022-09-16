package api

import (
	"net/http"
	db "sqlc/db/sqlc"
	"sqlc/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserParam struct {
	Username string `json:"username" binding:"required,alpha"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"full_name" binding:"required,alpha"`
	Email    string `json:"email" binding:"required,email"`
}

type userResp struct {
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	Email string `json:"email"`
	PasswordChange time.Time `json:"password_change"`
	CreatedAt time.Time `json:"created_at"`
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
	resp := userResp{
		Username:       account.Username,
		Fullname:       account.FullName,
		Email:          account.Email,
		PasswordChange: account.PasswordChangedAt,
		CreatedAt:      account.CreatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

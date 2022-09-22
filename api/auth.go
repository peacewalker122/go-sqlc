package api

import (
	"errors"
	"fmt"
	"net/http"
	"sqlc/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const(
	authHeaderkey = "authentication"
	authTypeBearer = "bearer"
	authPayload = "authorization_payload"
)

func authMiddleware(token token.Maker) gin.HandlerFunc {
	return func (ctx *gin.Context)  {

		//to get the header
		authorizationHeader := ctx.GetHeader(authHeaderkey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is empty")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorhandle(err))
			return
		}

		// to split the strings to slice
		authorizationHeaderFields := strings.Fields(authorizationHeader)
		if len(authorizationHeaderFields) < 2 {
			err := errors.New("unknown authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorhandle(err))
			return
		}

		// authentication to check the type the authentication
		authType := strings.ToLower(authorizationHeaderFields[0])
		if authType != authTypeBearer {
			err := fmt.Errorf("invalid authorization %v type", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorhandle(err))
			return
		}

		// declaring the token variable
		authToken := authorizationHeaderFields[1]

		// to verify the token
		payload,err := token.VerifyToken(authToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorhandle(err))
			return
		}

		ctx.Set(authPayload,payload)
		ctx.Next()
	}
}

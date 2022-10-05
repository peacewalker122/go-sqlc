package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/peacewalker122/go-sqlc/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Addauthorization(
	t *testing.T,
	req *http.Request,
	tokenMaker token.Maker,
	authType string,
	username string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	AuthHeader := fmt.Sprintf("%s %s", authType, token)
	assert.NoError(t, err)
	req.Header.Set(authHeaderkey, AuthHeader)
}

func TestAuth(t *testing.T) {

	TestCases := []struct {
		name             string
		setupAuth        func(t *testing.T, request *http.Request, token token.Maker)
		responseRecorder func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, token token.Maker) {
				Addauthorization(t, request, token, authTypeBearer, "test", time.Minute)
			},
			responseRecorder: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuth",
			setupAuth: func(t *testing.T, request *http.Request, token token.Maker) {
				Addauthorization(t, request, token, "awsLock", "test", time.Minute)
			},
			responseRecorder: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnauthorizedAuth",
			setupAuth: func(t *testing.T, request *http.Request, token token.Maker) {
				//Addauthorization(t,request,token,"awsLock","test",time.Minute)
			},
			responseRecorder: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthentication",
			setupAuth: func(t *testing.T, request *http.Request, token token.Maker) {
				Addauthorization(t, request, token, "", "test", time.Minute)
			},
			responseRecorder: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, token token.Maker) {
				Addauthorization(t, request, token, authTypeBearer, "test", -time.Minute)
			},
			responseRecorder: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range TestCases {
		tc := TestCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := NewTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.TokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.TokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.responseRecorder(t, recorder)
		})
	}
}

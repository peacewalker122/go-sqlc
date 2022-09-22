package api

import (
	"os"
	db "sqlc/db/sqlc"
	"sqlc/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T,store db.Store) *server{
	config := util.Config{
		SymmectricKey: util.Randomstring(32),
		Duration: time.Minute,
	}
	server,err := Newserver(config,store)
	require.NoError(t,err)

	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

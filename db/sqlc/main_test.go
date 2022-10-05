package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/peacewalker122/go-sqlc/util"

	_ "github.com/lib/pq"
)

var Testqueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can't definef config ", err.Error())
	}
	TestDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	Testqueries = New(TestDB)
	os.Exit(m.Run())
}

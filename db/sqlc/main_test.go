package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
)

const (
	DBdriver = "postgres"
	DBsource = "postgresql://postgres:test123@localhost:5432/simple_bank"
)

var Testqueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	TestDB, err = sql.Open(DBdriver, DBsource)
	if err != nil {
		log.Fatal(err)
	}
	Testqueries = New(TestDB)
	os.Exit(m.Run())
}

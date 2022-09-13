package main

import (
	"database/sql"
	"log"
	"sqlc/api"
	db "sqlc/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	DBdriver = "postgres"
	DBsource = "postgresql://postgres:test123@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(DBdriver, DBsource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.Newstore(conn)
	server := api.Newserver(store)

	err = server.Runserver("localhost:8080")
	if err != nil {
		log.Fatal("can't establish connection due ", err.Error())
	}

}

package main

import (
	"database/sql"
	"log"
	"sqlc/api"
	db "sqlc/db/sqlc"

	"sqlc/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't definef config ", err.Error())
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.Newstore(conn)
	server := api.Newserver(store)

	err = server.Runserver(config.ServerAddress)
	if err != nil {
		log.Fatal("can't establish connection due ", err.Error())
	}

}

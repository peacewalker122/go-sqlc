package main

import (
	"database/sql"
	"log"
	"net"
	"sqlc/api"
	db "sqlc/db/sqlc"
	"sqlc/gapi"
	"sqlc/pb"

	"sqlc/util"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	GinServer(config,store)
}

func gRPCServer(config util.Config, store db.Store){
	server, err := gapi.Newserver(config, store)
	if err != nil {
		log.Fatal("can't establish server due ", err.Error())
	}
	Grpcserver := grpc.NewServer()
	pb.RegisterSimpleBankServer(Grpcserver,server)
	
	reflection.Register(Grpcserver)

	listener,err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Panic("can't listen to the server")
	}

	log.Printf("Establish gRPC connection at %s",listener.Addr().String())
	err = Grpcserver.Serve(listener)
	if err != nil {
		log.Fatal("can't establish server")
	}
}

func GinServer(config util.Config, store db.Store){
	server, err := api.Newserver(config, store)
	if err != nil {
		log.Fatal("can't establish server due ", err.Error())
	}

	err = server.Runserver(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("can't establish connection due ", err.Error())
	}
}

package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/peacewalker122/go-sqlc/api"
	db "github.com/peacewalker122/go-sqlc/db/sqlc"
	"github.com/peacewalker122/go-sqlc/gapi"
	"github.com/peacewalker122/go-sqlc/pb"

	"github.com/peacewalker122/go-sqlc/util"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't define config: ", err.Error())
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.Newstore(conn)
	go GatewayServer(config, store)
	gRPCServer(config, store)
}

func gRPCServer(config util.Config, store db.Store) {
	server, err := gapi.Newserver(config, store)
	if err != nil {
		log.Fatal("can't establish server due ", err.Error())
	}
	Grpcserver := grpc.NewServer()
	pb.RegisterSimpleBankServer(Grpcserver, server)

	reflection.Register(Grpcserver)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Panic("can't listen to the server")
	}

	log.Printf("Establish gRPC connection at %s", listener.Addr().String())
	err = Grpcserver.Serve(listener)
	if err != nil {
		log.Fatal("can't establish server")
	}
}

func GatewayServer(config util.Config, store db.Store) {
	server, err := gapi.Newserver(config, store)
	if err != nil {
		log.Fatal("can't establish server due ", err.Error())
	}
	gRPCmux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, gRPCmux, server)
	if err != nil {
		log.Fatal("can't register server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", gRPCmux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Panic("can't listen to the server")
	}

	log.Printf("Establish HTTP connection at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("can't establish server")
	}
}

func GinServer(config util.Config, store db.Store) {
	server, err := api.Newserver(config, store)
	if err != nil {
		log.Fatal("can't establish server due ", err.Error())
	}

	err = server.Runserver(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("can't establish connection due ", err.Error())
	}
}

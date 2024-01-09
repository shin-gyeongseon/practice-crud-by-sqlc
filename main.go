package main

import (
	"database/sql"
	"go-practice/api"
	"go-practice/db/tutorial"
	"go-practice/gapi"
	"go-practice/pb"
	"go-practice/util"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't not connect DB", err)
	}

	store := tutorial.NewStore(conn)
	// runGinServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store tutorial.Store) {
	server := gapi.NewServer(store, config)

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer) // 이거는 뭐하는거지 ? 

	listener, err := net.Listen("tcp", config.GrpcADDRESS)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store tutorial.Store) {
	server := api.NewServer(store, config)
	err := server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connot start server: ", err)
	}
}

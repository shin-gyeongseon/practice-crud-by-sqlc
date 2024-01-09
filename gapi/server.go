package gapi

import (
	"go-practice/db/tutorial"
	"go-practice/pb"
	"go-practice/token"
	"go-practice/util"
	"log"
)

type ServerIf interface {
}

// Server servers gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	store      tutorial.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(store tutorial.Store, config util.Config) *Server {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmenticKey)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server
}

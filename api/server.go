package api

import (
	"go-practice/db/tutorial"
	"go-practice/token"
	"go-practice/util"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ServerIf interface {
}

// Server servers HTTP requests for our banking service
type Server struct {
	store  tutorial.Store
	router *gin.Engine
	maker  token.Maker
	config util.Config
}

func NewServer(store tutorial.Store, config util.Config) *Server {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmenticKey)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	server := &Server{
		store: store,
		maker: tokenMaker,
		config: config,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// account
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)

	// transfer
	router.POST("/transfer", server.CreateTransfer)

	// user
	router.POST("/user", server.CreateUser)

	// login
	router.POST("/login", server.LoginUser)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

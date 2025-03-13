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
	store      tutorial.Store
	router     *gin.Engine
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
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// user
	router.POST("/user", server.CreateUser)
	router.POST("/login", server.LoginUser)

	// account
	authRoutes := router.Group("/").Use(authMiddleware(tokenMaker))
	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.GetAccount)

	// transfer
	authRoutes.POST("/transfer", server.CreateTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

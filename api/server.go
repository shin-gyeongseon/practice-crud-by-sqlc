package api

import (
	"go-practice/db/tutorial"

	"github.com/gin-gonic/gin"
)

// Server servers HTTP requests for our banking service
type Server struct {
	*tutorial.Store
	router *gin.Engine
}

func NewServer (store *tutorial.Store) *Server {
	server := &Server{Store: store}
	router := gin.Default()

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
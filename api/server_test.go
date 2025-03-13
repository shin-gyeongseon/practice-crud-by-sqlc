package api

import (
	"go-practice/db/tutorial"
	"go-practice/util"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, db tutorial.Store) *Server{
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	server := NewServer(db, config)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

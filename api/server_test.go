package api

import (
	"go-practice/util"
	"log"
	"testing"
)

var TestGlobalConfig util.Config

func TestMain(m *testing.M) {
	t, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalln(err)
		return
	}

	TestGlobalConfig = t
}

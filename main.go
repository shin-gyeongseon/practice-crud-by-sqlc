package main

import (
	"database/sql"
	"go-practice/api"
	"go-practice/db/tutorial"
	"go-practice/util"
	"log"

	_ "github.com/lib/pq"
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
	server := api.NewServer(store, config)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connot start server: ", err)
	}
}

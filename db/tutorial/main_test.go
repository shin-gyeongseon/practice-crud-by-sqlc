package tutorial

import (
	"database/sql"
	"go-practice/token"
	"go-practice/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB
var testSqlStore Store
var testTokenMaer token.Maker
var testConfig util.Config

func TestMain(m *testing.M) {
	var err error

	// test query
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't not connect DB", err)
	}

	testQueries = New(testDB)
	testSqlStore = NewStore(testDB)

	// config , token maker 
	config, err2 := util.LoadConfig("../")
	if err2 != nil {
		log.Fatalln(err2)
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmenticKey)
	if err != nil {
		log.Fatalln(err)
	}
	 
	testConfig = config
	testTokenMaer = tokenMaker

	os.Exit(m.Run())
}

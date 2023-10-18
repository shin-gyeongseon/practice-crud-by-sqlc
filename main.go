package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-practice/db/tutorial"
	"log"

	_ "github.com/lib/pq"
)

// connection and create account , quote sqlc get star
func run() error {
	ctx := context.Background()

	db, err := sql.Open("postgres", "user=root password=secret dbname=simple_bank sslmode=disable")
	if err != nil {
		return err
	}

	queries := tutorial.New(db)

	account, err := queries.CreateAccount(ctx, tutorial.CreateAccountParams{
		Owner:    "sample owner",
		Balance:  30,
		Currency: "sample Currency",
	})

	if err != nil {
		return err
	}
	log.Println("create account ", account)

	accounts, err := queries.ListAccounts(ctx)
	if err != nil {
		return err
	}
	log.Println(accounts)

	return nil
}

type OperateFunc func(int, int) int

func getOperator(op string) OperateFunc {
	if op == "+" {
		return func(a, b int) int {
			return a + b
		}
	} else if op == "*" {
		return func(a, b int) int {
			return a * b
		}
	} else {
		return nil
	}
}

func main() {
	// if err := run(); err != nil {
	// 	log.Fatal(err)
	// }

	var op OperateFunc

	op = getOperator("+")
	result := op(10, 20)
	fmt.Println(result)
}

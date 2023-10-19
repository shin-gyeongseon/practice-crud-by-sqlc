postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres test server
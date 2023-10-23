postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

test:
	go test -v -cover -short ./...

winsqlc:
	docker run --rm -v C:\shiftone-projects\go-practice\db:/src -w /src sqlc/sqlc:latest generate

winmock
	mockgen -destination db\mock\store.go -source .\db\tutorial\store.go

server:
	go run main.go

.PHONY:postgres test winsqlc winmock server
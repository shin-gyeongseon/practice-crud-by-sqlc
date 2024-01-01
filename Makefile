DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

test:
	go test -v -cover -short ./...

winsqlc:
	docker run --rm -v C:\shiftone-projects\go-practice\db:/src -w /src sqlc/sqlc:latest generate

macsqlc:
	sqlc generate

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

mock:
	mockgen -source db/tutorial/store.go -destination db/mock/store.go
	
server:
	go run main.go

.PHONY:postgres test winsqlc macsqlc migrateup migrateup1 migrateup2 migrateup3 migratedown migratedown1 migratedown2 migratedown3 new_migration winmock server
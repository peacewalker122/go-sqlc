postgres:
	postgres12 -p 5432:5432 -e POSTGRES_USER=golang -e POSTGRES_PASSWORD=test123 -d postgres:12-alpine
createdb:
	createdb --username=postgres --owner=postgres simple_bank
dropdb:
	dropdb simple_bank 
migrateup:
	migrate -path db/migration -database "postgresql://postgres:test123@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://postgres:test123@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v ./...
server:
	go run main.go
open:
	psql -U postgres -d simple_bank

postgres:
	postgres12 -p 5432:5432 -e POSTGRES_USER=golang -e POSTGRES_PASSWORD=test123 -d postgres:12-alpine

createdb:
	createdb --username=postgres --owner=postgres simple_bank

dropdb:
	dropdb simple_bank 

migrateup:
	migrate -path db/migration -database "postgresql://postgres:test123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:test123@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:test123@localhost:5432/simple_bank?sslmode=disable" -verbose down
	
migratedown1:
	migrate -path db/migration -database "postgresql://postgres:test123@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v ./... -cover

server:
	go run main.go
	
mock:
	 mockgen -package mockdb -destination db/mock/store.go sqlc/db/sqlc Store

grpc:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl
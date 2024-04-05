build:
	@go build -o bin/server cmd/server/main.go

run: build
	@./bin/server

test:
	@go test -v -cover ./...

# can pass migration file name as name=file_name
migrate_init:
	@goose -dir sql/schema postgres ${DATABASE_URL} create $(name) sql

migrate_up:
	@goose -dir sql/schema postgres ${DATABASE_URL} up

migrate_down:
	@goose -dir sql/schema postgres ${DATABASE_URL} down
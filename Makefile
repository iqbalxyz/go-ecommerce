ifneq (,$(wildcard .env))
    include .env
    export
endif

MIGRATIONS_PATH ?= ./migrations
SERVER_PATH ?= ./cmd/api
BIN_PATH ?= ./bin/api

.PHONY: run build test tidy \
        migrate-up migrate-down migrate-reset migrate-version migrate-create

# Server
run:
	go run $(SERVER_PATH)

build:
	go build -o $(BIN_PATH) $(SERVER_PATH)

test:
	go test ./... -v

tidy:
	go mod tidy

# Migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-reset:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-create:
	@powershell -Command "if ('$(NAME)' -eq '') { Write-Host 'Error: NAME is required. Usage: make migrate-create NAME=migrations_name' -ForegroundColor Red; exit 1 }"
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
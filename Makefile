# Load .env if exist
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Default variables can be override via environment
DB_URL ?= $(env_DB_URL)
MIGRATIONS_PATH ?= ./migrations
SERVER_MAIN ?= main.go

.PHONY: migrate-up migrate-down migrate-reset migrate-version migrate-create run

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

run:
	go run $(SERVER_MAIN)
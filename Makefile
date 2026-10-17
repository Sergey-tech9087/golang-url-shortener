-include .env
export

export PROJECT_ROOT=$(shell pwd)
COMPOSE      := docker compose
COMPOSE_PG   := $(COMPOSE) --profile postgres

# --- PostgreSQL ---
env-up:
	@$(COMPOSE_PG) up -d url-shortener-postgres

env-down:
	@$(COMPOSE_PG) down url-shortener-postgres

env-cleanup: env-down
	@rm -rf out/pgdata

# --- Миграции ---
migrate-create:
	@docker compose run --rm url-shortener-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@docker compose run --rm url-shortener-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@url-shortener-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"${action}"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

# --- Приложение ---
app-up-memory:
	STORAGE_TYPE=memory $(COMPOSE) up --build -d url-shortener-app

app-up-postgres:
	STORAGE_TYPE=postgres $(COMPOSE_PG) up --build -d url-shortener-app

app-down:
	$(COMPOSE) down

url-shortener-run:
	@go run ./cmd/url-shortener-app/main.go
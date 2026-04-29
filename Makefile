.PHONY: run test docker-up docker-down migrate-up help

run:
	go run cmd/server/main.go

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

docker-logs:
	docker compose -f deployments/docker-compose.yml logs -f

migrate-up:
	@echo "Running database migrations..."
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f internal/repository/postgres/migrations/001_init.sql
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f internal/repository/postgres/migrations/002_add_indexes.sql

test:
	go test ./... -v

tidy:
	go mod tidy

help:
	@echo "Available commands:"
	@echo "  make run         - Start the server"
	@echo "  make docker-up   - Start all dependencies (PostgreSQL, InfluxDB, Redis, EMQX)"
	@echo "  make docker-down - Stop all dependencies"
	@echo "  make docker-logs - View dependency logs"
	@echo "  make migrate-up  - Run database migrations"
	@echo "  make test        - Run tests"
	@echo "  make tidy        - Tidy go modules"

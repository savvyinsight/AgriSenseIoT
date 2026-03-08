.PHONY: run test docker-up docker-down migrate-up help

run:
	go run cmd/server/main.go

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

docker-logs:
	docker compose -f deployments/docker-compose.yml logs -f

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
	@echo "  make test        - Run tests"
	@echo "  make tidy        - Tidy go modules"

.PHONY: help build run stop clean test dev logs ps restart db-shell db-backup db-restore

# Default target
help:
	@echo "URL Shortener - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start all services in development mode"
	@echo "  make logs         - View logs from all services"
	@echo "  make ps           - Show running containers"
	@echo "  make restart      - Restart all services"
	@echo "  make stop         - Stop all services"
	@echo "  make clean        - Stop and remove all containers, volumes, and images"
	@echo ""
	@echo "Build & Run:"
	@echo "  make build        - Build the application container"
	@echo "  make run          - Start all services"
	@echo "  make run-local    - Run application locally (without Docker)"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make bench        - Run benchmarks"
	@echo ""
	@echo "Database:"
	@echo "  make db-shell     - Open PostgreSQL shell"
	@echo "  make db-backup    - Backup database to backup.sql"
	@echo "  make db-restore   - Restore database from backup.sql"
	@echo "  make db-logs      - View database logs"
	@echo ""

# Build the application
build:
	@echo "Building application container..."
	docker-compose build app

# Start all services
run:
	@echo "Starting all services..."
	docker-compose up -d
	@echo "Services started! Access at http://localhost:8080"

# Run application locally (without Docker)
run-local:
	@echo "Starting PostgreSQL only..."
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Starting application locally..."
	go run main.go

# Stop all services
stop:
	@echo "Stopping all services..."
	docker-compose down

# Clean up everything
clean:
	@echo "Cleaning up containers, volumes, and images..."
	docker-compose down -v --rmi all
	@echo "Removing binary..."
	@rm -f url-shortener
	@echo "Cleanup complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test ./internal/utils -bench=. -benchmem

# Start services in development mode (with logs)
dev:
	@echo "Starting services in development mode..."
	docker-compose up --build

# View logs
logs:
	docker-compose logs -f

# View database logs only
db-logs:
	docker-compose logs -f postgres

# Show running containers
ps:
	docker-compose ps

# Restart all services
restart:
	@echo "Restarting all services..."
	docker-compose restart

# Open PostgreSQL shell
db-shell:
	docker exec -it url-shortener-db psql -U urlshortener -d urlshortener

# Backup database
db-backup:
	@echo "Backing up database to backup.sql..."
	docker exec url-shortener-db pg_dump -U urlshortener urlshortener > backup.sql
	@echo "Backup complete: backup.sql"

# Restore database
db-restore:
	@echo "Restoring database from backup.sql..."
	@cat backup.sql | docker exec -i url-shortener-db psql -U urlshortener urlshortener
	@echo "Restore complete!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod verify

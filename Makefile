# Makefile for Internal Transfer System

.PHONY: help setup run test clean docker-up docker-down deps fmt lint

# Default target
help:
	@echo "Available commands:"
	@echo "  setup       - Setup the project (start DB, install deps)"
	@echo "  run         - Run the application"
	@echo "  test        - Run API tests"
	@echo "  deps        - Install Go dependencies"
	@echo "  fmt         - Format Go code"
	@echo "  lint        - Run linting (requires golangci-lint)"
	@echo "  docker-up   - Start PostgreSQL container"
	@echo "  docker-down - Stop PostgreSQL container"
	@echo "  clean       - Clean up containers and dependencies"

# Setup the project
setup: docker-up deps
	@echo "Setup complete! You can now run 'make run' to start the server."

# Start PostgreSQL container
docker-up:
	@echo "Starting PostgreSQL container..."
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5

# Stop PostgreSQL container
docker-down:
	@echo "Stopping PostgreSQL container..."
	docker-compose down

# Install Go dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy

# Run the application
run:
	@echo "Starting Internal Transfer System..."
	go run cmd/main.go

# Run API tests
test:
	@echo "Running API tests..."
	@if [ ! -f test_api.sh ]; then echo "test_api.sh not found!"; exit 1; fi
	@chmod +x test_api.sh
	./test_api.sh

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run linting (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean up
clean: docker-down
	@echo "Cleaning up..."
	go clean
	go mod tidy

# Build the application
build:
	@echo "Building application..."
	go build -o bin/internal-transfer-system cmd/main.go

# Run in development mode with auto-reload (requires air)
dev:
	@echo "Starting in development mode..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running normally..."; \
		make run; \
	fi 
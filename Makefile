.PHONY: help install dev build test clean docker-up docker-down migrate seed

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing Go dependencies..."
	@cd services/auth-service && go mod download
	@cd services/org-service && go mod download
	@cd services/api-gateway && go mod download
	@echo "Installing Python dependencies..."
	@cd services/ai-service && pip install -r requirements.txt
	@echo "Dependencies installed!"

dev: ## Run all services in development mode
	@echo "Starting all services..."
	docker-compose up

build: ## Build all services
	@echo "Building all services..."
	@cd services/auth-service && go build -o ../../bin/auth-service ./cmd/main.go
	@cd services/org-service && go build -o ../../bin/org-service ./cmd/main.go
	@cd services/api-gateway && go build -o ../../bin/api-gateway ./cmd/main.go
	@echo "Build complete!"

test: ## Run tests for all services
	@echo "Running tests..."
	@cd services/auth-service && go test ./...
	@cd services/org-service && go test ./...
	@cd services/api-gateway && go test ./...
	@echo "Tests complete!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf services/*/bin/
	@echo "Clean complete!"

docker-up: ## Start all Docker services
	docker-compose up -d

docker-down: ## Stop all Docker services
	docker-compose down

docker-clean: ## Remove all Docker containers and volumes
	docker-compose down -v

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run scripts/migrate/main.go
	@echo "Migrations complete!"

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@go run scripts/seed/main.go
	@echo "Seeding complete!"

proto: ## Generate protobuf files
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		shared/proto/*.proto
	@echo "Protobuf generation complete!"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/main.go -o docs/swagger --pd
	@echo "Swagger docs generated!"

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...
	@echo "Linting complete!"

format: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@go mod tidy
	@echo "Formatting complete!"

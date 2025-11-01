.PHONY: help build run test clean migrate setup install docker-build docker-run

# Application variables
APP_NAME=restaurant-backend
BINARY_NAME=bin/server
MAIN_PATH=cmd/server/main.go

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Set up the project (install dependencies)
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully!"

install: setup ## Install dependencies (alias for setup)

build: ## Build the application
	@echo "Building application..."
	@mkdir -p bin
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

run: ## Run the application
	@echo "Starting application..."
	go run $(MAIN_PATH)

run-dev: ## Run the application in development mode
	@echo "Starting application in development mode..."
	ENVIRONMENT=development LOG_LEVEL=debug go run $(MAIN_PATH)

migrate: ## Run database migrations (schema only, no bootstrap)
	@echo "Running database migrations..."
	go run $(MAIN_PATH) --migrate

bootstrap: ## Bootstrap platform organization and admin user
	@echo "Bootstrapping platform organization and admin user..."
	go run $(MAIN_PATH) --bootstrap

setup-db: migrate bootstrap ## Set up database: run migrations and bootstrap
	@echo "Database setup complete!"

test: ## Run tests
	@echo "Running tests..."
	go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

check-env: ## Check if .env file exists
	@if [ ! -f .env ]; then \
		echo "Warning: .env file not found. Creating from .env.example..."; \
		cp .env.example .env; \
		echo "Please edit .env file with your configuration."; \
	fi

dev: check-env migrate run-dev ## Development workflow: check env, migrate, and run

# Database commands
db-migrate-up: ## Run database migrations (using migrate tool if available)
	@echo "Running database migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up; \
	else \
		echo "migrate tool not found. Using built-in migration..."; \
		go run $(MAIN_PATH) --migrate; \
	fi

db-migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down; \
	else \
		echo "migrate tool not found. Using built-in migration..."; \
		go run $(MAIN_PATH) --migrate-down; \
	fi

db-migrate-status: ## Show migration status
	@echo "Checking migration status..."
	go run $(MAIN_PATH) --migrate-status

# Quick start
start: check-env build run ## Quick start: check env, build, and run


.PHONY: build run test clean migrate docker-build docker-run help

# Variables
BINARY_NAME=lookie
BUILD_DIR=bin
VERSION ?= 1.0.0
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_HASH := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitHash=$(GIT_HASH)"

# Default target
all: build

## Build commands
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/lookie/main.go

build-linux: ## Build for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux cmd/lookie/main.go

## Development commands
run: ## Run the application
	go run cmd/lookie/main.go

run-server: ## Run server only mode
	go run cmd/lookie/main.go -server-only

run-scheduler: ## Run scheduler only mode
	go run cmd/lookie/main.go -scheduler-only

dev: ## Run with development configuration
	@echo "Starting development server..."
	go run cmd/lookie/main.go

## Database commands
migrate: ## Run database migrations
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

## Testing commands
test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -cover ./...

test-race: ## Run tests with race detection
	go test -v -race ./...

## Code quality
fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	golangci-lint run

## Docker commands
docker-build: ## Build Docker image
	docker build -t lookie:$(VERSION) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env lookie:$(VERSION)

docker-compose-up: ## Start with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop docker-compose
	docker-compose down

## Deployment commands
deploy-build: build-linux ## Build for deployment
	@echo "Built deployment binary: $(BUILD_DIR)/$(BINARY_NAME)-linux"

## Utility commands
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

deps: ## Download dependencies
	go mod download
	go mod tidy

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Hash: $(GIT_HASH)"

setup: ## Initial project setup
	@echo "Setting up project..."
	cp config.example.yaml config.yaml
	cp .env.example .env
	mkdir -p data
	@echo "Please edit config.yaml and .env with your settings"

health-check: ## Check if the service is running
	@curl -f http://localhost:8080/health || echo "Service not running or unhealthy"

## Help
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
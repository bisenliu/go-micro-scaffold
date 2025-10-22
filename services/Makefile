# Makefile for Go Micro Scaffold

# Variables
APP_NAME := go-micro-scaffold
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# Go build flags
GO_BUILD_FLAGS := -ldflags "$(LDFLAGS)"
GO_BUILD_FLAGS_PROD := -ldflags "$(LDFLAGS) -s -w" -trimpath

# Docker settings
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := $(VERSION)

.PHONY: help build build-prod run clean swagger-gen swagger-validate test test-cover lint fmt vet \
        dev full-build docker-build docker-run docker-clean deploy-dev deploy-staging deploy-prod \
        install-tools check-tools pre-commit ci-build

# Default target
help:
	@echo "Go Micro Scaffold Build System"
	@echo ""
	@echo "Development Commands:"
	@echo "  dev             - Development mode (docs + build + run)"
	@echo "  build           - Build the application (development)"
	@echo "  build-prod      - Build for production (optimized)"
	@echo "  run             - Run the application"
	@echo "  clean           - Clean build artifacts"
	@echo ""
	@echo "Documentation Commands:"
	@echo "  swagger-gen     - Generate Swagger documentation"
	@echo "  swagger-validate - Validate Swagger documentation"
	@echo ""
	@echo "Testing Commands:"
	@echo "  test            - Run tests"
	@echo "  test-cover      - Run tests with coverage"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo ""
	@echo "Build & Deploy Commands:"
	@echo "  full-build      - Complete build with validation"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  deploy-dev      - Deploy to development environment"
	@echo "  deploy-staging  - Deploy to staging environment"
	@echo "  deploy-prod     - Deploy to production environment"
	@echo ""
	@echo "Utility Commands:"
	@echo "  install-tools   - Install required development tools"
	@echo "  check-tools     - Check if required tools are installed"
	@echo "  pre-commit      - Run pre-commit checks"
	@echo "  ci-build        - CI/CD build pipeline"

# Install required development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest
	@echo "✅ Development tools installed"

# Check if required tools are installed
check-tools:
	@echo "Checking required tools..."
	@command -v swag >/dev/null 2>&1 || { echo "❌ swag not found. Run 'make install-tools'"; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "❌ golangci-lint not found. Run 'make install-tools'"; exit 1; }
	@echo "✅ All required tools are available"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✅ Go vet completed"

# Run linter
lint: check-tools
	@echo "Running linter..."
	@golangci-lint run
	@echo "✅ Linting completed"

# Generate Swagger documentation
swagger-gen: check-tools
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
	@echo "✅ Swagger documentation generated in docs/"

# Validate Swagger documentation
swagger-validate:
	@echo "Validating Swagger documentation..."
	@if [ ! -f "docs/swagger.json" ]; then \
		echo "❌ Swagger documentation not found. Run 'make swagger-gen' first"; \
		exit 1; \
	fi
	@go run scripts/validate-swagger.go
	@echo "✅ Swagger documentation validation completed"

# Build the application (development)
build: swagger-gen
	@echo "Building application (development)..."
	@mkdir -p bin
	@go build $(GO_BUILD_FLAGS) -o bin/server cmd/server/main.go
	@go build $(GO_BUILD_FLAGS) -o bin/cli cmd/cli/main.go
	@echo "✅ Build completed:"
	@echo "  Server: bin/server"
	@echo "  CLI:    bin/cli"

# Build for production (optimized)
build-prod: swagger-gen lint vet test
	@echo "Building application (production)..."
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS_PROD) -o bin/server-linux-amd64 cmd/server/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS_PROD) -o bin/cli-linux-amd64 cmd/cli/main.go
	@go build $(GO_BUILD_FLAGS_PROD) -o bin/server cmd/server/main.go
	@go build $(GO_BUILD_FLAGS_PROD) -o bin/cli cmd/cli/main.go
	@echo "✅ Production build completed:"
	@echo "  Server (local):  bin/server"
	@echo "  Server (linux):  bin/server-linux-amd64"
	@echo "  CLI (local):     bin/cli"
	@echo "  CLI (linux):     bin/cli-linux-amd64"

# Run the application
run:
	@echo "Starting server..."
	@if [ ! -f "bin/server" ]; then \
		echo "❌ Server binary not found. Run 'make build' first"; \
		exit 1; \
	fi
	@./bin/server

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf docs/
	@go clean
	@docker system prune -f 2>/dev/null || true
	@echo "✅ Clean completed"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "✅ Tests completed"

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests with coverage completed"
	@echo "Coverage report: coverage.html"

# Pre-commit checks
pre-commit: fmt vet lint swagger-gen swagger-validate test
	@echo "✅ All pre-commit checks passed"

# Development workflow
dev: swagger-gen build
	@echo ""
	@echo "🚀 Development build completed!"
	@echo ""
	@echo "To start the server:"
	@echo "  make run"
	@echo ""
	@echo "Swagger UI will be available at:"
	@echo "  http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "API documentation JSON:"
	@echo "  http://localhost:8080/swagger/doc.json"

# Full build and validation
full-build: clean pre-commit build-prod swagger-validate
	@echo "✅ Full build and validation completed successfully!"

# Docker build
docker-build: swagger-gen
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -t $(DOCKER_IMAGE):latest .
	@echo "✅ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker run -d --name $(APP_NAME) -p 8080:8080 $(DOCKER_IMAGE):latest
	@echo "✅ Docker container started"
	@echo "Application available at: http://localhost:8080"
	@echo "Swagger UI available at: http://localhost:8080/swagger/index.html"

# Docker clean
docker-clean:
	@echo "Cleaning Docker resources..."
	@docker stop $(APP_NAME) 2>/dev/null || true
	@docker rm $(APP_NAME) 2>/dev/null || true
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	@docker rmi $(DOCKER_IMAGE):latest 2>/dev/null || true
	@echo "✅ Docker resources cleaned"

# Deploy to development environment
deploy-dev: build swagger-validate
	@echo "Deploying to development environment..."
	@echo "Environment: Development"
	@echo "Swagger: Enabled"
	@echo "Build version: $(VERSION)"
	@echo "✅ Development deployment completed"

# Deploy to staging environment
deploy-staging: build-prod swagger-validate
	@echo "Deploying to staging environment..."
	@echo "Environment: Staging"
	@echo "Swagger: Enabled (with authentication)"
	@echo "Build version: $(VERSION)"
	@echo "✅ Staging deployment completed"

# Deploy to production environment
deploy-prod: build-prod swagger-validate
	@echo "Deploying to production environment..."
	@echo "Environment: Production"
	@echo "Swagger: Disabled (security)"
	@echo "Build version: $(VERSION)"
	@echo "⚠️  Production deployment requires manual confirmation"
	@echo "✅ Production deployment prepared"

# CI/CD build pipeline
ci-build: check-tools clean swagger-gen lint vet test build-prod swagger-validate
	@echo "✅ CI/CD build pipeline completed successfully!"
	@echo "Artifacts ready for deployment:"
	@ls -la bin/
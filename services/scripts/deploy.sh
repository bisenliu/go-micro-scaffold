#!/bin/bash

# Deployment script for Go Micro Scaffold
# Supports multiple environments with different Swagger configurations

set -e

# Configuration
APP_NAME="go-micro-scaffold"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help function
show_help() {
    cat << EOF
Go Micro Scaffold Deployment Script

Usage: $0 [ENVIRONMENT] [OPTIONS]

Environments:
  dev         Deploy to development environment (Swagger enabled)
  staging     Deploy to staging environment (Swagger with auth)
  prod        Deploy to production environment (Swagger disabled)

Options:
  -h, --help          Show this help message
  -v, --version       Show version information
  -c, --config FILE   Use custom configuration file
  -d, --dry-run       Show what would be deployed without executing
  -f, --force         Force deployment without confirmation
  --skip-tests        Skip running tests before deployment
  --skip-docs         Skip generating Swagger documentation
  --docker            Deploy using Docker containers

Examples:
  $0 dev                    # Deploy to development
  $0 staging --docker       # Deploy to staging using Docker
  $0 prod --dry-run         # Show production deployment plan
  $0 dev -c custom.yaml     # Deploy with custom config

Environment Variables:
  VERSION                   Override version (default: git describe)
  CONFIG_FILE              Configuration file path
  DOCKER_REGISTRY          Docker registry URL
  DEPLOY_TARGET            Deployment target (local, k8s, docker-compose)

EOF
}

# Parse command line arguments
ENVIRONMENT=""
CONFIG_FILE=""
DRY_RUN=false
FORCE=false
SKIP_TESTS=false
SKIP_DOCS=false
USE_DOCKER=false

while [[ $# -gt 0 ]]; do
    case $1 in
        dev|staging|prod)
            ENVIRONMENT="$1"
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--version)
            echo "Version: $VERSION"
            echo "Build Time: $BUILD_TIME"
            echo "Git Commit: $GIT_COMMIT"
            exit 0
            ;;
        -c|--config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --skip-docs)
            SKIP_DOCS=true
            shift
            ;;
        --docker)
            USE_DOCKER=true
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate environment
if [[ -z "$ENVIRONMENT" ]]; then
    log_error "Environment is required"
    show_help
    exit 1
fi

if [[ ! "$ENVIRONMENT" =~ ^(dev|staging|prod)$ ]]; then
    log_error "Invalid environment: $ENVIRONMENT"
    log_error "Valid environments: dev, staging, prod"
    exit 1
fi

# Set environment-specific configurations
setup_environment() {
    case $ENVIRONMENT in
        dev)
            SWAGGER_ENABLED=true
            SWAGGER_AUTH=false
            LOG_LEVEL="debug"
            PORT=8080
            ;;
        staging)
            SWAGGER_ENABLED=true
            SWAGGER_AUTH=true
            LOG_LEVEL="info"
            PORT=8080
            ;;
        prod)
            SWAGGER_ENABLED=false
            SWAGGER_AUTH=false
            LOG_LEVEL="warn"
            PORT=8080
            ;;
    esac
    
    log_info "Environment: $ENVIRONMENT"
    log_info "Swagger Enabled: $SWAGGER_ENABLED"
    log_info "Swagger Auth: $SWAGGER_AUTH"
    log_info "Log Level: $LOG_LEVEL"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if we're in the right directory
    if [[ ! -f "cmd/server/main.go" ]]; then
        log_error "Not in services directory or main.go not found"
        exit 1
    fi
    
    # Check required tools
    local missing_tools=()
    
    if ! command -v go &> /dev/null; then
        missing_tools+=("go")
    fi
    
    if [[ "$SKIP_DOCS" != true ]] && ! command -v swag &> /dev/null; then
        missing_tools+=("swag")
    fi
    
    if [[ "$USE_DOCKER" == true ]] && ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        log_error "Run 'make install-tools' to install missing tools"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Generate configuration file
generate_config() {
    local config_file="${CONFIG_FILE:-configs/app.yaml}"
    
    if [[ "$DRY_RUN" == true ]]; then
        log_info "Would generate configuration: $config_file"
        return
    fi
    
    log_info "Generating configuration for $ENVIRONMENT environment..."
    
    # Create environment-specific config
    cat > "$config_file" << EOF
# Generated configuration for $ENVIRONMENT environment
# Generated at: $BUILD_TIME
# Version: $VERSION

server:
  port: $PORT
  mode: $ENVIRONMENT

swagger:
  enabled: $SWAGGER_ENABLED
  title: "Go Micro Scaffold API"
  description: "微服务脚手架 API 文档 ($ENVIRONMENT)"
  version: "$VERSION"
  host: "localhost:$PORT"
  base_path: "/api/v1"

logging:
  level: "$LOG_LEVEL"
  format: "json"

database:
  # Database configuration should be set via environment variables
  # or external configuration management

redis:
  # Redis configuration should be set via environment variables
  # or external configuration management
EOF
    
    log_success "Configuration generated: $config_file"
}

# Generate Swagger documentation
generate_docs() {
    if [[ "$SKIP_DOCS" == true ]]; then
        log_warning "Skipping Swagger documentation generation"
        return
    fi
    
    if [[ "$DRY_RUN" == true ]]; then
        log_info "Would generate Swagger documentation"
        return
    fi
    
    log_info "Generating Swagger documentation..."
    
    if ! swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal; then
        log_error "Failed to generate Swagger documentation"
        exit 1
    fi
    
    # Validate generated documentation
    if [[ -f "docs/swagger.json" && -f "docs/swagger.yaml" && -f "docs/docs.go" ]]; then
        log_success "Swagger documentation generated successfully"
    else
        log_error "Swagger documentation generation incomplete"
        exit 1
    fi
}

# Run tests
run_tests() {
    if [[ "$SKIP_TESTS" == true ]]; then
        log_warning "Skipping tests"
        return
    fi
    
    if [[ "$DRY_RUN" == true ]]; then
        log_info "Would run tests"
        return
    fi
    
    log_info "Running tests..."
    
    if ! go test ./...; then
        log_error "Tests failed"
        exit 1
    fi
    
    log_success "All tests passed"
}

# Build application
build_app() {
    if [[ "$DRY_RUN" == true ]]; then
        log_info "Would build application for $ENVIRONMENT"
        return
    fi
    
    log_info "Building application for $ENVIRONMENT..."
    
    # Set build flags
    local ldflags="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"
    
    if [[ "$ENVIRONMENT" == "prod" ]]; then
        ldflags="$ldflags -s -w"
    fi
    
    # Create bin directory
    mkdir -p bin
    
    # Build server
    if ! go build -ldflags "$ldflags" -o bin/server cmd/server/main.go; then
        log_error "Failed to build server"
        exit 1
    fi
    
    # Build CLI
    if ! go build -ldflags "$ldflags" -o bin/cli cmd/cli/main.go; then
        log_error "Failed to build CLI"
        exit 1
    fi
    
    log_success "Application built successfully"
}

# Build Docker image
build_docker() {
    if [[ "$USE_DOCKER" != true ]]; then
        return
    fi
    
    if [[ "$DRY_RUN" == true ]]; then
        log_info "Would build Docker image: $APP_NAME:$VERSION"
        return
    fi
    
    log_info "Building Docker image..."
    
    if ! docker build -t "$APP_NAME:$VERSION" -t "$APP_NAME:latest" .; then
        log_error "Failed to build Docker image"
        exit 1
    fi
    
    log_success "Docker image built: $APP_NAME:$VERSION"
}

# Deploy application
deploy_app() {
    if [[ "$DRY_RUN" == true ]]; then
        log_info "Would deploy to $ENVIRONMENT environment"
        show_deployment_summary
        return
    fi
    
    # Production deployment requires confirmation
    if [[ "$ENVIRONMENT" == "prod" && "$FORCE" != true ]]; then
        log_warning "Production deployment requires confirmation"
        echo -n "Are you sure you want to deploy to production? (yes/no): "
        read -r confirmation
        
        if [[ "$confirmation" != "yes" ]]; then
            log_info "Deployment cancelled"
            exit 0
        fi
    fi
    
    log_info "Deploying to $ENVIRONMENT environment..."
    
    if [[ "$USE_DOCKER" == true ]]; then
        deploy_docker
    else
        deploy_binary
    fi
    
    log_success "Deployment completed successfully"
    show_deployment_summary
}

# Deploy using Docker
deploy_docker() {
    log_info "Deploying using Docker..."
    
    # Stop existing container
    docker stop "$APP_NAME" 2>/dev/null || true
    docker rm "$APP_NAME" 2>/dev/null || true
    
    # Run new container
    docker run -d \
        --name "$APP_NAME" \
        -p "$PORT:8080" \
        -e GO_ENV="$ENVIRONMENT" \
        "$APP_NAME:$VERSION"
    
    log_success "Docker container deployed"
}

# Deploy binary
deploy_binary() {
    log_info "Deploying binary..."
    
    # In a real deployment, this would copy files to target servers
    # For now, we just show what would happen
    log_info "Binary deployment would:"
    log_info "  - Copy bin/server to target server"
    log_info "  - Copy configuration files"
    log_info "  - Restart service"
    log_info "  - Verify health checks"
}

# Show deployment summary
show_deployment_summary() {
    echo ""
    echo "======================================"
    echo "       DEPLOYMENT SUMMARY"
    echo "======================================"
    echo "Environment:     $ENVIRONMENT"
    echo "Version:         $VERSION"
    echo "Build Time:      $BUILD_TIME"
    echo "Git Commit:      $GIT_COMMIT"
    echo "Swagger Enabled: $SWAGGER_ENABLED"
    echo "Docker:          $USE_DOCKER"
    echo "Port:            $PORT"
    echo ""
    echo "Application URL: http://localhost:$PORT"
    
    if [[ "$SWAGGER_ENABLED" == true ]]; then
        echo "Swagger UI:      http://localhost:$PORT/swagger/index.html"
        echo "API Docs JSON:   http://localhost:$PORT/swagger/doc.json"
    fi
    
    echo "======================================"
}

# Cleanup function
cleanup() {
    if [[ "$DRY_RUN" != true ]]; then
        log_info "Cleaning up temporary files..."
        # Add cleanup logic here if needed
    fi
}

# Trap cleanup on exit
trap cleanup EXIT

# Main deployment flow
main() {
    log_info "Starting deployment process..."
    log_info "Version: $VERSION"
    log_info "Environment: $ENVIRONMENT"
    
    setup_environment
    check_prerequisites
    generate_config
    generate_docs
    run_tests
    build_app
    build_docker
    deploy_app
    
    log_success "Deployment process completed!"
}

# Run main function
main "$@"
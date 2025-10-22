# Multi-stage Dockerfile for Go Micro Scaffold with Swagger Integration

# Build arguments
ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.18

# Stage 1: Build dependencies and tools
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS tools

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl

# Install Swagger CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Stage 2: Build application
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT=unknown

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Copy Swagger tool from tools stage
COPY --from=tools /go/bin/swag /usr/local/bin/swag

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -s -w" \
    -trimpath \
    -o bin/server \
    cmd/server/main.go

# Build CLI tool
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -s -w" \
    -trimpath \
    -o bin/cli \
    cmd/cli/main.go

# Stage 3: Development image (with source code and tools)
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS development

# Install development tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl \
    bash \
    make

# Install Go development tools
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Expose ports
EXPOSE 8080 2345

# Default command for development
CMD ["air", "-c", ".air.toml"]

# Stage 4: Production image (minimal)
FROM alpine:${ALPINE_VERSION} AS production

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/cli /app/cli

# Copy configuration files
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/docs /app/docs

# Create logs directory
RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command
CMD ["/app/server"]

# Stage 5: Final stage selector
FROM ${TARGET_STAGE:-production} AS final
# Build stage
FROM golang:1.24-alpine AS builder

# Install git and make
RUN apk add --no-cache git make

WORKDIR /app

# Copy workspace files
COPY go.work go.work.sum ./
COPY common/go.mod common/go.sum ./common/
COPY services/go.mod services/go.sum ./services/

# Download dependencies
RUN go mod download

# Copy source code
COPY common ./common
COPY services ./services

# Build the application
RUN go build -o /app/server services/cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (including netcat for wait-for script)
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

# Copy binary from builder
COPY --from=builder /app/server .

# Copy configs
COPY services/configs ./configs

# Expose port (adjust if needed, default 8080 based on main.go comments)
EXPOSE 8080

# Run the application
CMD ["./server"]

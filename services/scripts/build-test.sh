#!/bin/bash

# Build test script for Swagger integration
echo "=== Building Go Micro Scaffold with Swagger Integration ==="

# Change to services directory
cd "$(dirname "$0")/.."

# Clean any previous builds
echo "Cleaning previous builds..."
go clean

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Verify all imports can be resolved
echo "Verifying imports..."
go mod verify

# Build the application
echo "Building application..."
if go build -o bin/server cmd/server/main.go; then
    echo "✅ Build successful!"
    
    # Check if swagger docs exist
    if [ -f "docs/swagger.json" ] && [ -f "docs/swagger.yaml" ] && [ -f "docs/docs.go" ]; then
        echo "✅ Swagger documentation files found!"
    else
        echo "❌ Swagger documentation files missing!"
        exit 1
    fi
    
    # Validate swagger docs
    echo "Validating Swagger documentation..."
    if go run scripts/validate-swagger.go; then
        echo "✅ Swagger documentation validation passed!"
    else
        echo "❌ Swagger documentation validation failed!"
        exit 1
    fi
    
    echo ""
    echo "=== Build Summary ==="
    echo "✅ Application compiled successfully"
    echo "✅ Swagger documentation generated"
    echo "✅ All validations passed"
    echo ""
    echo "To start the server:"
    echo "  ./bin/server"
    echo ""
    echo "Swagger UI will be available at:"
    echo "  http://localhost:8080/swagger/index.html"
    
else
    echo "❌ Build failed!"
    exit 1
fi
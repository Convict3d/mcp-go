# MCP Go Client Library Makefile

.PHONY: help build test lint clean examples doc doc-serve validate check deps

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build all packages"
	@echo "  test      - Run all tests"
	@echo "  lint      - Run linting tools"
	@echo "  clean     - Clean build artifacts"
	@echo "  examples  - Build all examples"
	@echo "  doc       - Show package documentation"
	@echo "  doc-serve - Serve documentation locally"
	@echo "  validate  - Validate package for public release"
	@echo "  check     - Run all checks (test + lint)"
	@echo "  deps      - Install/update dependencies"

# Build all packages
build:
	@echo "Building all packages..."
	go build ./...

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linting
lint:
	@echo "Running linting..."
	go vet ./...
	go fmt ./...
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean ./...
	rm -f coverage.out coverage.html

# Build examples
examples:
	@echo "Building examples..."
	go build ./examples/...

# Generate documentation
doc:
	@echo "Generating documentation..."
	@echo "=== Root Package ==="
	@go doc .
	@echo ""
	@echo "=== Client Package ==="
	@go doc ./client
	@echo ""
	@echo "=== Types Package ==="
	@go doc ./types
	@echo ""
	@echo "=== HTTP Transport ==="
	@go doc ./transport/http
	@echo ""
	@echo "=== Stdio Transport ==="
	@go doc ./transport/stdio
	@echo ""
	@echo "To serve documentation locally:"
	@echo "  go install golang.org/x/pkgsite/cmd/pkgsite@latest"
	@echo "  pkgsite -http=:8080"
	@echo "  Open: http://localhost:8080/github.com/Convict3d/mcp-go"

# Serve documentation locally
doc-serve:
	@echo "Starting documentation server..."
	@echo "Visit: http://localhost:8080/github.com/Convict3d/mcp-go"
	pkgsite -http=:8080

# Validate package for public release
validate:
	@echo "Validating package for public release..."
	@echo "Checking go.mod..."
	@go mod tidy
	@go mod verify
	@echo "Checking documentation..."
	@go doc . > /dev/null || (echo "Error: Missing package documentation" && exit 1)
	@go doc ./client > /dev/null || (echo "Error: Missing client package documentation" && exit 1)
	@go doc ./types > /dev/null || (echo "Error: Missing types package documentation" && exit 1)
	@echo "Checking examples..."
	@for example in examples/*.go; do \
		if [ -f "$$example" ]; then \
			go build -o /tmp/example "$$example" || exit 1; \
		fi; \
	done
	@echo "Running tests..."
	@go test ./...
	@echo "Package validation successful! Ready for pkg.go.dev"

# Run all checks
check: test lint
	@echo "All checks passed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run example
run-example:
	@echo "To run the basic example:"
	@echo "  go run ./examples/basic/main.go"

# Release preparation
prepare-release: clean test lint doc
	@echo "Release preparation complete!"

# Development setup
dev-setup: deps
	@echo "Development setup complete!"
	@echo "Available commands:"
	@echo "  make build    - Build the library"
	@echo "  make test     - Run tests"
	@echo "  make examples - Build examples"

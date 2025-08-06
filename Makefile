# MCP Go Client Library Makefile

.PHONY: help build test lint clean examples doc

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build all packages"
	@echo "  test      - Run all tests"
	@echo "  lint      - Run linting tools"
	@echo "  clean     - Clean build artifacts"
	@echo "  examples  - Build all examples"
	@echo "  doc       - Generate documentation"
	@echo "  check     - Run all checks (test + lint)"

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
	go doc -all .
	@if command -v godoc > /dev/null; then \
		echo "Run 'godoc -http=:6060' to serve documentation"; \
	fi

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

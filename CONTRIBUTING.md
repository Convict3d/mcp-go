# Contributing to MCP Go Client Library

Thank you for your interest in contributing to the MCP Go Client Library! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Running Tests](#running-tests)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker to report bugs
- Before creating an issue, please check if a similar issue already exists
- Include as much detail as possible:
  - Go version
  - Operating system
  - Steps to reproduce
  - Expected vs actual behavior
  - Error messages or logs

### Suggesting Enhancements

- Use the GitHub issue tracker to suggest new features
- Explain the use case and why the enhancement would be valuable
- Consider if the enhancement aligns with the project's goals

### Contributing Code

1. Fork the repository
2. Create a feature branch from `master`
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Clone and Setup

```bash
git clone https://github.com/Convict3d/mcp-go.git
cd mcp-go
go mod download
```

### Install Development Tools

```bash
# Install linting tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install documentation tools
go install golang.org/x/pkgsite/cmd/pkgsite@latest
```

## Making Changes

### Project Structure

```
├── client/          # High-level MCP client implementation
├── types/           # MCP protocol type definitions
├── transport/       # Transport layer implementations
│   ├── http/        # HTTP transport
│   └── stdio/       # Standard I/O transport
├── examples/        # Usage examples
├── .github/         # GitHub Actions workflows
└── docs/            # Documentation
```

### Branch Naming

Use descriptive branch names:
- `feature/add-new-transport`
- `fix/client-timeout-issue`
- `docs/improve-examples`

### Commit Messages

Follow conventional commit format:
- `feat: add support for new MCP protocol version`
- `fix: resolve client timeout issue`
- `docs: update API documentation`
- `test: add tests for stdio transport`

## Running Tests

### Basic Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

```bash
# Run linting
make lint

# Or run golangci-lint directly
golangci-lint run
```

### Full Validation

```bash
# Run complete validation
make validate
```

## Submitting Changes

### Pull Request Process

1. **Update Documentation**: Ensure any new features are documented
2. **Add Tests**: All new code should include appropriate tests
3. **Update Examples**: Add or update examples if relevant
4. **Run Tests**: Ensure all tests pass locally
5. **Update Changelog**: Add an entry describing your changes

### Pull Request Description

Include in your PR description:
- What changes were made and why
- How to test the changes
- Any breaking changes
- Links to related issues

### Review Process

- All PRs require at least one review
- Address review feedback promptly
- Keep PRs focused and reasonably sized
- Rebase on master before merging

## Code Style

### Go Style Guidelines

Follow standard Go conventions:
- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types

### Specific Guidelines

1. **Error Handling**: Always handle errors appropriately
   ```go
   result, err := someOperation()
   if err != nil {
       return fmt.Errorf("operation failed: %w", err)
   }
   ```

2. **Context Usage**: Use context for cancellation and timeouts
   ```go
   func (c *Client) SomeOperation(ctx context.Context) error {
       // Implementation
   }
   ```

3. **Interface Design**: Keep interfaces small and focused
   ```go
   type Transport interface {
       Send(request interface{}) (interface{}, error)
       Close() error
   }
   ```

4. **Option Pattern**: Use functional options for configuration
   ```go
   func WithTimeout(timeout time.Duration) Option {
       return func(c *Config) {
           c.Timeout = timeout
       }
   }
   ```

### Linting Rules

The project uses `golangci-lint` with specific rules. See `.golangci.yml` for the complete configuration.

## Documentation

### Package Documentation

- Add package-level documentation in `doc.go` files
- Document all exported functions, types, and constants
- Include usage examples in documentation comments

### Examples

- Add runnable examples to the `examples/` directory
- Ensure examples are well-commented and demonstrate best practices
- Update examples when adding new features

### README Updates

When adding significant features:
- Update the main README.md
- Add usage examples
- Update the feature list

## Testing Guidelines

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    // Setup
    client := setupTestClient()
    
    // Execute
    result, err := client.DoSomething()
    
    // Verify
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

### Test Coverage

- Aim for high test coverage (>80%)
- Test both success and error cases
- Include edge cases and boundary conditions
- Use table-driven tests for multiple scenarios

### Integration Tests

- Add integration tests for complete workflows
- Test with real MCP servers when possible
- Use test doubles for external dependencies

## Questions?

If you have questions about contributing:
- Check existing issues and documentation
- Open a discussion on GitHub
- Contact the maintainers

Thank you for contributing to the MCP Go Client Library!

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.9.0] - 2025-08-06

### Added
- Initial public release of MCP Go client library
- High-level client with option pattern configuration
- HTTP transport implementation with JSON-RPC support
- Stdio transport implementation for process-based servers
- Complete MCP protocol type definitions (version 2025-06-18)
- Support for all MCP operations:
  - Tools (list, call with arguments)
  - Resources (list, read, templates)
  - Prompts (list, get with arguments)
- Content type support (text, image, audio, resource links)
- Comprehensive test coverage with 31+ test functions
- Professional documentation with examples
- GitHub Actions CI/CD pipeline
- Linting configuration with golangci-lint

### Features
- Clean API design with Go option pattern
- Type-safe interfaces for all MCP operations
- Production-ready error handling and timeouts
- Minimal dependencies (jsonrpc, sync primitives)
- Complete protocol compliance with MCP 2025-06-18
- Thread-safe client implementation
- Configurable transports (HTTP, stdio)
- Server capability detection
- Custom headers and authentication support

### Documentation
- Comprehensive README with examples
- Package-level documentation for all modules
- Working examples in examples/ directory
- Contributing guidelines
- MIT license

### Technical
- Go 1.21+ compatibility
- Semantic versioning
- Professional package structure for pkg.go.dev
- Full test coverage with integration tests

[Unreleased]: https://github.com/Convict3d/mcp-go/compare/v0.9.0...HEAD
[0.9.0]: https://github.com/Convict3d/mcp-go/releases/tag/v0.9.0

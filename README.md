# MCP Go Client Library

A professional Go client library for the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/). This library provides a clean, type-safe implementation for connecting Go applications to MCP servers.

[![Go Reference](https://pkg.go.dev/badge/github.com/Convict3d/mcp-go.svg)](https://pkg.go.dev/github.com/Convict3d/mcp-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Convict3d/mcp-go)](https://goreportcard.com/report/github.com/Convict3d/mcp-go)

## Overview

The Model Context Protocol enables seamless integration between LLM applications and external data sources and tools. This library implements the MCP specification with a focus on:

- **Clean API Design** - Intuitive Go option pattern for configuration
- **Type Safety** - Strongly typed interfaces for all MCP operations
- **Production Ready** - Comprehensive error handling and timeouts
- **Minimal Dependencies** - Lightweight with essential dependencies only
- **Complete Protocol Support** - Tools, resources, prompts, and all content types

## Installation

```bash
go get github.com/Convict3d/mcp-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/Convict3d/mcp-go/client"
    "github.com/Convict3d/mcp-go/types"
)

func main() {
    // Create client with option pattern
    c := client.NewClient("http://localhost:9831/mcp",
        client.WithClientInfo("my-app", "1.0.0"),
        client.WithTimeout(30*time.Second),
        client.WithSSESupport(),
    )

    // Initialize connection
    err := c.Initialize(types.LatestProtocolVersion)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // List available tools
    tools, err := c.ListTools()
    if err != nil {
        log.Fatal(err)
    }

    for _, tool := range tools {
        fmt.Printf("Tool: %s - %s\n", tool.Name, tool.Description)
    }

    // Call a tool
    result, err := c.CallTool("search", map[string]interface{}{
        "query": "golang mcp",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Extract text content from result
    textContents := result.GetTextContent()
    for _, content := range textContents {
        fmt.Printf("Result: %s\n", content.Text)
    }
}
```

## Client Configuration

Configure the client using the option pattern:

```go
c := client.NewClient("http://localhost:9831/mcp",
    // Set client identification
    client.WithClientInfo("my-app", "2.1.0"),
    
    // Configure timeouts
    client.WithTimeout(45*time.Second),
    
    // Add custom headers
    client.WithHeader("Authorization", "Bearer token"),
    client.WithCustomHeaders(map[string]string{
        "X-API-Key": "secret",
        "User-Agent": "MyApp/1.0",
    }),
    
    // Enable Server-Sent Events
    client.WithSSESupport(),
)
```

### Custom Transport

You can also provide your own transport for advanced configurations:

```go
// HTTP Transport - Client connects to running MCP server
httpTransport := transport.NewHTTPTransport("http://localhost:9831/mcp",
    transport.WithTimeout(60*time.Second),
    transport.WithHeader("X-Custom-Auth", "token123"),
    transport.WithSSESupport(),
)

// Stdio Transport - Client launches MCP server as subprocess  
stdioTransport, err := transport.NewStdioTransport("mcp-filesystem-server",
    transport.WithArgs("--root", "/home/user/documents"),
    transport.WithWorkingDir("/tmp"),
)

// Use the transport with the client
c := client.NewClient(
    client.WithTransport(stdioTransport),
    client.WithClientInfo("my-app", "1.0.0"),
)
```

### Stdio Transport Use Cases

The stdio transport supports different scenarios:

```go
// 1. Client launches MCP server (most common)
transport, err := transport.NewStdioTransport("mcp-server-executable",
    transport.WithArgs("--config", "config.json"),
)

// 2. Use existing streams (for custom scenarios)  
transport, err := transport.NewStdioTransportFromStreams(stdin, stdout, stderr)

// 3. Current process IS an MCP server (rare)
transport, err := transport.NewStdioTransportFromOS()
```

### Configuration Options

- `WithClientInfo(name, version)` - Set client identification
- `WithTimeout(duration)` - Configure request timeout  
- `WithHeader(key, value)` - Add single custom header
- `WithCustomHeaders(map)` - Add multiple custom headers
- `WithTransport(transport)` - Use custom transport
- `WithSSESupport()` - Enable Server-Sent Events support

### Simple Client

For quick prototyping:

```go
c := client.NewSimpleClient("http://localhost:9831/mcp")
```

## Features

### Tools

Execute tools on the MCP server:

```go
// List available tools
tools, err := c.ListTools()
if err != nil {
    log.Fatal(err)
}

// Call a tool with arguments
result, err := c.CallTool("calculator", map[string]interface{}{
    "operation": "multiply",
    "a": 15,
    "b": 3,
})
if err != nil {
    log.Fatal(err)
}

// Access different content types from results
textContent := result.GetTextContent()
imageContent := result.GetImageContent()
audioContent := result.GetAudioContent()
```

### Resources

Access server resources:

```go
// List available resources
resources, err := c.ListResources()
if err != nil {
    log.Fatal(err)
}

// Read a specific resource
content, err := c.ReadResource("file://path/to/document.txt")
if err != nil {
    log.Fatal(err)
}

// Process the resource content
for _, item := range content.Contents {
    switch content := item.(type) {
    case *types.TextContent:
        fmt.Println("Text:", content.Text)
    case *types.ImageContent:
        fmt.Println("Image data length:", len(content.Data))
    }
}
```

### Prompts

Work with server prompts:

```go
// List available prompts
prompts, err := c.ListPrompts()
if err != nil {
    log.Fatal(err)
}

// Get a prompt with arguments
prompt, err := c.GetPrompt("code-review", map[string]string{
    "language": "go",
    "file": "main.go",
})
if err != nil {
    log.Fatal(err)
}

// Use the prompt messages
for _, message := range prompt.Messages {
    fmt.Printf("Role: %s\n", message.Role)
    // Process message content...
}
```

## Content Types

The library supports all MCP content types with full type safety:

```go
// Text content
textContent := &types.TextContent{
    Type: "text",
    Text: "Hello, world!",
}

// Image content
imageContent := &types.ImageContent{
    Type: "image",
    Data: base64ImageData,
    MimeType: "image/png",
}

// Audio content
audioContent := &types.AudioContent{
    Type: "audio", 
    Data: base64AudioData,
    MimeType: "audio/wav",
}

// Resource links
resourceLink := &types.ResourceLinkContent{
    Type: "resource",
    URI: "file://document.pdf",
}

// Embedded resources
resource := &types.ResourceContent{
    Type: "resource",
    Resource: types.EmbeddedResource{
        URI: "data://example",
        Text: "embedded content",
        MimeType: "text/plain",
    },
}
```

## Project Structure

Clean, focused package organization:

```
mcp-go/
├── client/          # High-level MCP client implementation  
├── transport/       # HTTP transport layer
├── types/          # Complete MCP protocol types
├── examples/       # Working examples
├── go.mod          # Module definition
└── README.md       # This file
```

### Package Details

- **`client/`** - Main client with option pattern configuration
- **`transport/`** - HTTP transport implementation with JSON-RPC support
- **`types/`** - Complete type definitions for MCP protocol messages
- **`examples/`** - Real-world usage examples and patterns

## Error Handling

Comprehensive error handling throughout:

```go
// Check server capabilities before operations
if !c.HasTools() {
    log.Println("Server does not support tools")
    return
}

// Handle specific errors
tools, err := c.ListTools()
if err != nil {
    log.Printf("Failed to list tools: %v", err)
    return
}

// Graceful resource cleanup
defer func() {
    if err := c.Close(); err != nil {
        log.Printf("Error closing client: %v", err)
    }
}()
```

## Server Capabilities

Check what features the server supports:

```go
// After initialization, check capabilities
fmt.Printf("Server: %s v%s\n", 
    c.GetServerInfo().Name, 
    c.GetServerInfo().Version)

if c.HasTools() {
    fmt.Println("✓ Tools supported")
}

if c.HasResources() {
    fmt.Println("✓ Resources supported") 
}

if c.HasPrompts() {
    fmt.Println("✓ Prompts supported")
}

// Access detailed capabilities
capabilities := c.GetCapabilities()
// Inspect specific capability details...
```

## Examples

Run the included examples:

```bash
# Basic client usage
cd examples/basic
go run main.go

# Real-world integration
cd examples/playwright-http-real  
go run main.go
```

## Advanced Usage

## Advanced Usage

### Custom Transport with Advanced Configuration

Create highly customized transports and use them with the client:

```go
// Create a transport with advanced configuration
transport := transport.NewHTTPTransport("http://localhost:9831/mcp",
    transport.WithTimeout(60*time.Second),
    transport.WithHeader("X-Custom-Auth", "token123"),
    transport.WithHeader("X-Request-ID", generateRequestID()),
    transport.WithCustomHeaders(map[string]string{
        "X-Client-Version": "2.1.0",
        "X-Platform": "golang",
    }),
    transport.WithSSESupport(),
)

// Use the custom transport with the client
c := client.NewClient("http://localhost:9831/mcp",
    client.WithTransport(transport),
    client.WithClientInfo("advanced-client", "2.1.0"),
)

// The client will use your custom transport configuration
err := c.Initialize(types.LatestProtocolVersion)
```

### Direct Transport Access

For maximum control, use the transport layer directly:

```go
import "github.com/Convict3d/mcp-go/transport"

// Create transport with options
t := transport.NewHTTPTransport("http://localhost:9831/mcp",
    transport.WithTimeout(45*time.Second),
    transport.WithHeader("Authorization", "Bearer token"),
    transport.WithCustomHeaders(map[string]string{
        "X-API-Key": "secret",
    }),
    transport.WithSSESupport(),
)

// Make direct JSON-RPC calls
var result map[string]interface{}
err := t.Call(context.Background(), &result, "tools/list")
```

### Legacy Transport Config

For backward compatibility, you can still use the config struct:

```go
config := transport.Config{
    ServerURL:     "http://localhost:9831/mcp",
    Timeout:       30 * time.Second,
    CustomHeaders: map[string]string{
        "Accept": "application/json",
    },
}

t := transport.NewHTTPTransportWithConfig(config)
```

### Custom Context

Use custom contexts for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

c := client.NewClient("http://localhost:9831/mcp",
    client.WithContext(ctx),
)
```

## Protocol Compliance

This library implements the complete MCP specification including:

- JSON-RPC 2.0 transport layer
- Initialization and capability negotiation
- Resources (listing, reading, templates)
- Tools (listing, calling with arguments)  
- Prompts (listing, getting with arguments)
- All content types (text, image, audio, resource links, embedded resources)
- Progress notifications and error handling
- Session management

## Requirements

- Go 1.24 or later
- Minimal external dependencies (see go.mod)

## Dependencies

The library uses these well-maintained dependencies:

- `github.com/creachadair/jrpc2` - JSON-RPC 2.0 implementation
- `github.com/ybbus/jsonrpc/v3` - Additional JSON-RPC utilities
- `golang.org/x/sync` - Go synchronization primitives

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch  
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Related

- [Model Context Protocol](https://modelcontextprotocol.io/) - Official specification
- [MCP Servers](https://github.com/modelcontextprotocol/servers) - Server implementations
- [Claude Desktop](https://claude.ai/desktop) - Desktop app with MCP support

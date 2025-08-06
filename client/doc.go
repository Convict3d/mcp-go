/*
Package client provides a high-level MCP client implementation.

This package offers a simple, intuitive API for connecting to MCP servers and performing
all protocol operations. The client handles connection management, request/response
processing, and error handling automatically.

# Basic Usage

Create and initialize a client:

	c := client.NewClient(
		client.WithClientInfo("my-app", "1.0.0"),
		client.WithTimeout(30*time.Second),
	)

	err := c.Initialize(types.LatestProtocolVersion)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

# Configuration Options

The client supports various configuration options:

	c := client.NewClient(
		client.WithClientInfo("app-name", "1.0.0"),    // Set client identification
		client.WithTimeout(30*time.Second),            // Set operation timeout
		client.WithTransport(customTransport),         // Use custom transport
	)

# Supported Operations

The client supports all MCP protocol operations:

  - ListTools() - List available tools
  - CallTool(name, args) - Call a specific tool
  - ListResources() - List available resources
  - ReadResource(uri) - Read resource content
  - ListPrompts() - List available prompts
  - GetPrompt(name, args) - Get prompt with arguments

# Error Handling

All client methods return appropriate Go errors:

	tools, err := c.ListTools()
	if err != nil {
		// Handle the error appropriately
		return fmt.Errorf("failed to list tools: %w", err)
	}

# Transport Layer

The client works with different transport implementations:

  - HTTP transport for REST-like connections
  - Stdio transport for process-based servers

Custom transports can be implemented by satisfying the transport.Transport interface.

# Thread Safety

The Client type is safe for concurrent use, but individual operations should be
performed from a single goroutine.
*/
package client

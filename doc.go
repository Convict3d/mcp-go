/*
Package mcp provides a comprehensive Go client library for the Model Context Protocol (MCP).

The Model Context Protocol enables seamless integration between LLM applications and external
data sources and tools. This library implements the MCP specification with a focus on clean
API design, type safety, and production readiness.

# Quick Start

Basic usage example:

	package main

	import (
		"log"
		"time"

		"github.com/Convict3d/mcp-go/client"
		"github.com/Convict3d/mcp-go/types"
	)

	func main() {
		// Create client with configuration options
		c := client.NewClient(
			client.WithClientInfo("my-app", "1.0.0"),
			client.WithTimeout(30*time.Second),
		)

		// Initialize MCP connection
		err := c.Initialize(types.LatestProtocolVersion)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		// Use the client to interact with MCP server
		tools, err := c.ListTools()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Found %d tools", len(tools))
	}

# Package Structure

The library is organized into the following packages:

  - client: High-level MCP client implementation
  - types: MCP protocol type definitions and constants
  - transport/http: HTTP transport implementation
  - transport/stdio: Standard I/O transport implementation

# Protocol Support

This library supports MCP protocol version 2025-06-18 and provides:

  - Tools: List and call external tools
  - Resources: Read and manage external resources
  - Prompts: Get and use prompt templates
  - Content: Handle text, image, and audio content types
  - Logging: Protocol-level logging capabilities
  - Sampling: LLM sampling and model preferences

# Configuration

The client supports flexible configuration through the option pattern:

	c := client.NewClient(
		client.WithClientInfo("app-name", "1.0.0"),
		client.WithTimeout(30*time.Second),
		client.WithTransport(customTransport),
	)

# Error Handling

All client methods return Go errors that should be checked:

	tools, err := client.ListTools()
	if err != nil {
		// Handle error appropriately
		log.Printf("Failed to list tools: %v", err)
		return
	}

# Thread Safety

The client is designed to be safe for concurrent use from multiple goroutines.
However, individual operations should be used from a single goroutine.

For more detailed examples and documentation, see the individual package documentation
and the examples directory.
*/
package mcp

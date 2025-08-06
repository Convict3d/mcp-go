/*
Package http provides HTTP transport implementation for MCP clients.

This package implements the transport.Transport interface using HTTP as the
underlying communication protocol. It supports both standard HTTP requests
and Server-Sent Events (SSE) for real-time communication.

# Basic Usage

Create an HTTP transport:

	transport := http.NewTransport("http://localhost:9831/mcp")

	// Use with client
	c := client.NewClient(client.WithTransport(transport))

# Configuration Options

The HTTP transport supports various configuration options:

	transport := http.NewTransport("http://localhost:9831/mcp",
		http.WithTimeout(30*time.Second),
		http.WithUserAgent("my-app/1.0"),
		http.WithCustomHeaders(map[string]string{
			"Authorization": "Bearer token",
		}),
	)

# Protocol Support

The HTTP transport supports:

  - Standard HTTP requests for MCP operations
  - Server-Sent Events (SSE) for real-time notifications
  - Custom headers for authentication and metadata
  - Configurable timeouts and retry logic

# Error Handling

The transport returns appropriate HTTP-related errors:

  - Network connectivity issues
  - HTTP status code errors
  - JSON parsing errors
  - Timeout errors

# Security

When using HTTP transport:

  - Always use HTTPS in production
  - Implement proper authentication headers
  - Validate SSL certificates
  - Use appropriate timeout values

# Thread Safety

The HTTP transport is safe for concurrent use from multiple goroutines.
*/
package http

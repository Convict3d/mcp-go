/*
Package stdio provides standard I/O transport implementation for MCP clients.

This package implements the transport.Transport interface using standard input/output
streams for communication. This is commonly used for MCP servers that run as
separate processes and communicate via stdin/stdout.

# Basic Usage

Create a stdio transport:

	transport := stdio.NewTransport()

	// Use with client
	c := client.NewClient(client.WithTransport(transport))

# Process Management

The stdio transport can manage external processes:

	transport := stdio.NewTransport(
		stdio.WithCommand("python", "my_mcp_server.py"),
		stdio.WithWorkDir("/path/to/server"),
		stdio.WithEnv(map[string]string{
			"DEBUG": "1",
			"CONFIG_PATH": "/etc/myapp",
		}),
	)

# Configuration Options

Available configuration options:

  - WithCommand() - Set the command to execute
  - WithWorkDir() - Set working directory for the process
  - WithEnv() - Set environment variables
  - WithStdin() - Use custom stdin reader
  - WithStdout() - Use custom stdout writer

# Process Lifecycle

The transport handles the complete process lifecycle:

 1. Start the external process
 2. Set up stdin/stdout communication channels
 3. Handle process termination and cleanup
 4. Manage process errors and restarts

# Communication Protocol

Communication follows the JSON-RPC protocol over stdio:

  - Requests are sent to the process stdin
  - Responses are read from the process stdout
  - Each message is newline-delimited JSON

# Error Handling

The transport handles various error conditions:

  - Process startup failures
  - Communication errors
  - Process crashes and unexpected termination
  - JSON parsing errors

# Thread Safety

The stdio transport is safe for concurrent use from multiple goroutines,
with proper synchronization for stdin/stdout access.

# Best Practices

When using stdio transport:

  - Ensure the target process follows MCP stdio protocol
  - Handle process termination gracefully
  - Monitor process health and restart if necessary
  - Use appropriate timeouts for operations
*/
package stdio

// Package transport provides stdio transport implementation for MCP
package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ybbus/jsonrpc/v3"
)

// StdioTransport implements MCP over stdio (standard input/output)
type StdioTransport struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	scanner   *bufio.Scanner
	mu        sync.RWMutex
	sessionID string
	nextID    int64
	closed    bool

	// For handling concurrent requests
	pendingRequests map[int64]chan *jsonrpc.RPCResponse
	requestsMu      sync.RWMutex

	// For handling notifications and server requests
	notificationHandler func(method string, params interface{})
	requestHandler      func(method string, params interface{}) (interface{}, error)

	// Background reader control
	stopReader chan struct{}
	readerDone chan struct{}
}

// StdioConfig holds configuration for stdio transport
type StdioConfig struct {
	Command    string   // Command to execute
	Args       []string // Command arguments
	WorkingDir string   // Working directory for the command
	Env        []string // Environment variables
}

// StdioOption defines a function that configures the stdio transport
type StdioOption func(*StdioConfig)

// WithCommand sets the command to execute
func WithCommand(command string) StdioOption {
	return func(c *StdioConfig) {
		c.Command = command
	}
}

// WithArgs sets the command arguments
func WithArgs(args ...string) StdioOption {
	return func(c *StdioConfig) {
		c.Args = args
	}
}

// WithWorkingDir sets the working directory
func WithWorkingDir(dir string) StdioOption {
	return func(c *StdioConfig) {
		c.WorkingDir = dir
	}
}

// WithEnv sets environment variables
func WithEnv(env []string) StdioOption {
	return func(c *StdioConfig) {
		c.Env = env
	}
}

// defaultStdioConfig returns a default stdio configuration
func defaultStdioConfig() *StdioConfig {
	return &StdioConfig{
		Command: "",
		Args:    []string{},
		Env:     os.Environ(),
	}
}

// NewStdioTransport creates a new stdio transport with options
func NewStdioTransport(command string, opts ...StdioOption) (*StdioTransport, error) {
	config := defaultStdioConfig()
	config.Command = command

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	return NewStdioTransportWithConfig(*config)
}

// NewStdioTransportFromStreams creates a stdio transport using existing streams
// This is useful when your program IS the MCP server and wants to communicate
// over its own stdin/stdout, or when you have custom streams
func NewStdioTransportFromStreams(stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser) (*StdioTransport, error) {
	transport := &StdioTransport{
		cmd:             nil, // No subprocess when using existing streams
		stdin:           stdin,
		stdout:          stdout,
		stderr:          stderr,
		scanner:         bufio.NewScanner(stdout),
		nextID:          1,
		pendingRequests: make(map[int64]chan *jsonrpc.RPCResponse),
		stopReader:      make(chan struct{}),
		readerDone:      make(chan struct{}),
	}

	// Start reading stderr in background (if provided)
	if stderr != nil {
		go transport.readStderr()
	}

	// Start reading stdout in background
	go transport.readMessages()

	return transport, nil
}

// NewStdioTransportFromOS creates a stdio transport using the current process's stdin/stdout
// This is useful when your Go program IS an MCP server
func NewStdioTransportFromOS() (*StdioTransport, error) {
	// Note: We don't close os.Stdin/Stdout in this case since they're owned by the OS
	return NewStdioTransportFromStreams(
		&nopCloser{os.Stdin},  // Wrap to prevent closing
		&nopCloser{os.Stdout}, // Wrap to prevent closing
		&nopCloser{os.Stderr}, // Wrap to prevent closing
	)
}

// nopCloser wraps a ReadWriteCloser but makes Close() a no-op
type nopCloser struct {
	io.ReadWriter
}

func (nc *nopCloser) Close() error {
	return nil // Don't actually close os.Stdin/Stdout/Stderr
}

// NewStdioTransportWithConfig creates a new stdio transport with config
func NewStdioTransportWithConfig(config StdioConfig) (*StdioTransport, error) {
	if config.Command == "" {
		return nil, fmt.Errorf("command is required for stdio transport")
	}

	// Create the command
	cmd := exec.Command(config.Command, config.Args...)
	if config.WorkingDir != "" {
		cmd.Dir = config.WorkingDir
	}
	if len(config.Env) > 0 {
		cmd.Env = config.Env
	}

	// Set up pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		stderr.Close()
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	transport := &StdioTransport{
		cmd:             cmd,
		stdin:           stdin,
		stdout:          stdout,
		stderr:          stderr,
		scanner:         bufio.NewScanner(stdout),
		nextID:          1,
		pendingRequests: make(map[int64]chan *jsonrpc.RPCResponse),
		stopReader:      make(chan struct{}),
		readerDone:      make(chan struct{}),
	}

	// Start reading stderr in background
	go transport.readStderr()

	// Start reading stdout in background
	go transport.readMessages()

	return transport, nil
}

// readMessages reads and processes JSON-RPC messages from stdout
func (t *StdioTransport) readMessages() {
	defer close(t.readerDone)

	for {
		select {
		case <-t.stopReader:
			return
		default:
			if !t.scanner.Scan() {
				if err := t.scanner.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "MCP stdio scanner error: %v\n", err)
				}
				return
			}

			line := t.scanner.Text()
			if line == "" {
				continue
			}

			t.processMessage(line)
		}
	}
}

// processMessage processes a single JSON-RPC message
func (t *StdioTransport) processMessage(line string) {
	var message map[string]interface{}
	if err := json.Unmarshal([]byte(line), &message); err != nil {
		fmt.Fprintf(os.Stderr, "MCP invalid JSON received: %v\n", err)
		return
	}

	// Check if this is a response (has id and result/error)
	if id, hasID := message["id"]; hasID {
		if _, hasResult := message["result"]; hasResult {
			t.handleResponse(message, id)
			return
		}
		if _, hasError := message["error"]; hasError {
			t.handleResponse(message, id)
			return
		}

		// This is a request from server
		t.handleServerRequest(message, id)
		return
	}

	// This is a notification
	t.handleNotification(message)
}

// handleResponse handles JSON-RPC responses
func (t *StdioTransport) handleResponse(message map[string]interface{}, id interface{}) {
	idFloat, ok := id.(float64)
	if !ok {
		return
	}

	requestID := int64(idFloat)

	t.requestsMu.Lock()
	responseChan, exists := t.pendingRequests[requestID]
	if exists {
		delete(t.pendingRequests, requestID)
	}
	t.requestsMu.Unlock()

	if exists {
		response := &jsonrpc.RPCResponse{}

		// Handle result
		if result, hasResult := message["result"]; hasResult {
			response.Result = result
		}

		// Handle error
		if errorData, hasError := message["error"]; hasError {
			if errorMap, ok := errorData.(map[string]interface{}); ok {
				rpcError := &jsonrpc.RPCError{}

				if code, ok := errorMap["code"].(float64); ok {
					rpcError.Code = int(code)
				}
				if msg, ok := errorMap["message"].(string); ok {
					rpcError.Message = msg
				}
				if data, ok := errorMap["data"]; ok {
					rpcError.Data = data
				}

				response.Error = rpcError
			}
		}

		select {
		case responseChan <- response:
		case <-time.After(1 * time.Second):
			// Channel might be closed or blocking
		}
	}
}

// handleServerRequest handles requests from the server
func (t *StdioTransport) handleServerRequest(message map[string]interface{}, id interface{}) {
	method, ok := message["method"].(string)
	if !ok {
		t.sendErrorResponse(id, -32600, "Invalid request: missing method")
		return
	}

	params := message["params"]

	if t.requestHandler != nil {
		result, err := t.requestHandler(method, params)
		if err != nil {
			t.sendErrorResponse(id, -32000, err.Error())
		} else {
			t.sendSuccessResponse(id, result)
		}
	} else {
		t.sendErrorResponse(id, -32601, "Method not found")
	}
}

// handleNotification handles notifications from the server
func (t *StdioTransport) handleNotification(message map[string]interface{}) {
	method, ok := message["method"].(string)
	if !ok {
		return
	}

	params := message["params"]

	if t.notificationHandler != nil {
		t.notificationHandler(method, params)
	}
}

// sendErrorResponse sends an error response to the server
func (t *StdioTransport) sendErrorResponse(id interface{}, code int, message string) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}

	t.sendMessage(response)
}

// sendSuccessResponse sends a success response to the server
func (t *StdioTransport) sendSuccessResponse(id interface{}, result interface{}) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}

	t.sendMessage(response)
}

// sendMessage sends a JSON-RPC message to the server
func (t *StdioTransport) sendMessage(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("transport is closed")
	}

	_, err = t.stdin.Write(append(data, '\n'))
	return err
}

// SetNotificationHandler sets the handler for server notifications
func (t *StdioTransport) SetNotificationHandler(handler func(method string, params interface{})) {
	t.notificationHandler = handler
}

// SetRequestHandler sets the handler for server requests
func (t *StdioTransport) SetRequestHandler(handler func(method string, params interface{}) (interface{}, error)) {
	t.requestHandler = handler
}

// readStderr reads and logs stderr output
func (t *StdioTransport) readStderr() {
	scanner := bufio.NewScanner(t.stderr)
	for scanner.Scan() {
		// Log stderr to help with debugging
		// In production, you might want to use a proper logger
		fmt.Fprintf(os.Stderr, "MCP Server stderr: %s\n", scanner.Text())
	}
}

// Call makes a JSON-RPC call over stdio
func (t *StdioTransport) Call(ctx context.Context, result interface{}, method string, params ...interface{}) error {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return fmt.Errorf("transport is closed")
	}
	t.mu.RUnlock()

	// Build the JSON-RPC request
	id := atomic.AddInt64(&t.nextID, 1)
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
	}

	if len(params) > 0 {
		request["params"] = params[0]
	}

	// Create response channel
	responseChan := make(chan *jsonrpc.RPCResponse, 1)

	t.requestsMu.Lock()
	t.pendingRequests[id] = responseChan
	t.requestsMu.Unlock()

	// Cleanup on function exit
	defer func() {
		t.requestsMu.Lock()
		delete(t.pendingRequests, id)
		t.requestsMu.Unlock()
		close(responseChan)
	}()

	// Send the request
	if err := t.sendMessage(request); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for response or context cancellation
	select {
	case response := <-responseChan:
		if response.Error != nil {
			return response.Error
		}

		if result != nil && response.Result != nil {
			resultBytes, err := json.Marshal(response.Result)
			if err != nil {
				return fmt.Errorf("failed to marshal result: %w", err)
			}

			if err := json.Unmarshal(resultBytes, result); err != nil {
				return fmt.Errorf("failed to unmarshal result: %w", err)
			}
		}

		return nil

	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(30 * time.Second):
		return fmt.Errorf("request timeout")
	}
}

// CallRaw makes a JSON-RPC call and returns the raw response
func (t *StdioTransport) CallRaw(ctx context.Context, method string, params interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := t.Call(ctx, &result, method, params)
	return result, err
}

// GetSessionID returns the current session ID
func (t *StdioTransport) GetSessionID() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.sessionID
}

// Close closes the stdio transport
func (t *StdioTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	// Stop the background message reader
	close(t.stopReader)

	// Close pending requests with error
	t.requestsMu.Lock()
	for id, ch := range t.pendingRequests {
		select {
		case ch <- &jsonrpc.RPCResponse{
			Error: &jsonrpc.RPCError{
				Code:    -32000,
				Message: "Transport closed",
			},
		}:
		default:
		}
		delete(t.pendingRequests, id)
	}
	t.requestsMu.Unlock()

	// Close pipes
	if t.stdin != nil {
		t.stdin.Close()
	}
	if t.stdout != nil {
		t.stdout.Close()
	}
	if t.stderr != nil {
		t.stderr.Close()
	}

	// Wait for reader to finish
	select {
	case <-t.readerDone:
	case <-time.After(1 * time.Second):
	}

	// Wait for the command to finish (only if we started a subprocess)
	if t.cmd != nil && t.cmd.Process != nil {
		// Give the process a moment to exit gracefully
		done := make(chan error, 1)
		go func() {
			done <- t.cmd.Wait()
		}()

		select {
		case err := <-done:
			return err
		case <-time.After(5 * time.Second):
			// Force kill if it doesn't exit gracefully
			if t.cmd.Process != nil {
				t.cmd.Process.Kill()
				return t.cmd.Wait()
			}
		}
	}

	return nil
}

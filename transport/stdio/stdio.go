// Package stdio implements MCP over stdio transport
package stdio

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

// Transport implements MCP over stdio (standard input/output)
type Transport struct {
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

// Config holds configuration for stdio transport
type Config struct {
	Command    string        // Command to execute
	Args       []string      // Command arguments
	WorkingDir string        // Working directory for the command
	Env        []string      // Environment variables
	Timeout    time.Duration // Request timeout
}

// Option defines a function that configures the stdio transport
type Option func(*Config)

// WithCommand sets the command to execute
func WithCommand(command string) Option {
	return func(c *Config) {
		c.Command = command
	}
}

// WithArgs sets the command arguments
func WithArgs(args ...string) Option {
	return func(c *Config) {
		c.Args = args
	}
}

// WithWorkingDir sets the working directory
func WithWorkingDir(dir string) Option {
	return func(c *Config) {
		c.WorkingDir = dir
	}
}

// WithEnv sets environment variables
func WithEnv(env []string) Option {
	return func(c *Config) {
		c.Env = env
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// defaultConfig returns a default stdio configuration
func defaultConfig() *Config {
	return &Config{
		Command: "",
		Args:    []string{},
		Env:     os.Environ(),
		Timeout: 30 * time.Second,
	}
}

// NewTransport creates a new stdio transport with options
func NewTransport(command string, args []string, opts ...Option) (*Transport, error) {
	config := defaultConfig()
	config.Command = command
	config.Args = args

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	return NewTransportWithConfig(*config)
}

// NewTransportFromStreams creates a stdio transport using existing streams
// This is useful when your program IS the MCP server and wants to communicate
// over its own stdin/stdout, or when you have custom streams
func NewTransportFromStreams(stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser) (*Transport, error) {
	transport := &Transport{
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

// NewTransportFromOS creates a stdio transport using the current process's stdin/stdout
// This is useful when your Go program IS an MCP server
func NewTransportFromOS() (*Transport, error) {
	// Note: We don't close os.Stdin/Stdout in this case since they're owned by the OS
	return NewTransportFromStreams(
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

// NewTransportWithConfig creates a new stdio transport with config
func NewTransportWithConfig(config Config) (*Transport, error) {
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

	transport := &Transport{
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
func (t *Transport) readMessages() {
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
func (t *Transport) processMessage(line string) {
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
func (t *Transport) handleResponse(message map[string]interface{}, id interface{}) {
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
func (t *Transport) handleServerRequest(message map[string]interface{}, id interface{}) {
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
func (t *Transport) handleNotification(message map[string]interface{}) {
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
func (t *Transport) sendErrorResponse(id interface{}, code int, message string) {
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
func (t *Transport) sendSuccessResponse(id interface{}, result interface{}) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}

	t.sendMessage(response)
}

// sendMessage sends a JSON-RPC message to the server
func (t *Transport) sendMessage(message interface{}) error {
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

// generateRequestID generates a unique request ID
func (t *Transport) generateRequestID() int64 {
	return atomic.AddInt64(&t.nextID, 1)
}

// SetNotificationHandler sets the handler for server notifications
func (t *Transport) SetNotificationHandler(handler func(method string, params interface{})) {
	t.notificationHandler = handler
}

// SetRequestHandler sets the handler for server requests
func (t *Transport) SetRequestHandler(handler func(method string, params interface{}) (interface{}, error)) {
	t.requestHandler = handler
}

// readStderr reads and logs stderr output
func (t *Transport) readStderr() {
	scanner := bufio.NewScanner(t.stderr)
	for scanner.Scan() {
		// Log stderr to help with debugging
		// In production, you might want to use a proper logger
		fmt.Fprintf(os.Stderr, "MCP Server stderr: %s\n", scanner.Text())
	}
}

// Call makes a JSON-RPC call over stdio
func (t *Transport) Call(ctx context.Context, result interface{}, method string, params ...interface{}) error {
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
func (t *Transport) CallRaw(ctx context.Context, method string, params interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := t.Call(ctx, &result, method, params)
	return result, err
}

// GetSessionID returns the current session ID
func (t *Transport) GetSessionID() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.sessionID
}

// Close closes the stdio transport
func (t *Transport) Close() error {
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

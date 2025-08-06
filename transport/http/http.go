// Package http implements MCP over HTTP transport
package http

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ybbus/jsonrpc/v3"
)

// SessionAwareHTTPClient is an HTTP client that handles MCP session management
type SessionAwareHTTPClient struct {
	client    *http.Client
	sessionID string
}

// NewSessionAwareHTTPClient creates a new session-aware HTTP client
func NewSessionAwareHTTPClient(timeout time.Duration) *SessionAwareHTTPClient {
	return &SessionAwareHTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do performs an HTTP request with session management
func (s *SessionAwareHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Add session ID if we have one
	if s.sessionID != "" {
		req.Header.Set("Mcp-Session-Id", s.sessionID)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Extract session ID from response
	if sessionID := resp.Header.Get("Mcp-Session-Id"); sessionID != "" {
		s.sessionID = sessionID
	}

	// Handle Server-Sent Events format
	if strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
		body, err := s.parseSSEResponse(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body = body
	}

	return resp, nil
}

// GetSessionID returns the current session ID
func (s *SessionAwareHTTPClient) GetSessionID() string {
	return s.sessionID
}

// parseSSEResponse parses Server-Sent Events format and extracts JSON data
func (s *SessionAwareHTTPClient) parseSSEResponse(body io.ReadCloser) (io.ReadCloser, error) {
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var jsonData strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data: ") {
			jsonLine := strings.TrimPrefix(line, "data: ")
			if jsonLine != "" && jsonLine != "[DONE]" {
				jsonData.WriteString(jsonLine)
			}
		}
	}

	return io.NopCloser(strings.NewReader(jsonData.String())), nil
}

// Config represents transport configuration
type Config struct {
	ServerURL     string
	Timeout       time.Duration
	CustomHeaders map[string]string
}

// Option defines a function that configures the transport
type Option func(*Config)

// WithTimeout sets the transport timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithCustomHeaders sets custom HTTP headers
func WithCustomHeaders(headers map[string]string) Option {
	return func(c *Config) {
		c.CustomHeaders = headers
	}
}

// WithHeader adds a single custom HTTP header
func WithHeader(key, value string) Option {
	return func(c *Config) {
		if c.CustomHeaders == nil {
			c.CustomHeaders = make(map[string]string)
		}
		c.CustomHeaders[key] = value
	}
}

// WithSSESupport is a convenience option that adds SSE headers
func WithSSESupport() Option {
	return WithHeader("Accept", "application/json, text/event-stream")
}

// defaultConfig returns a default transport configuration
func defaultConfig() *Config {
	return &Config{
		Timeout: 30 * time.Second,
		CustomHeaders: map[string]string{
			"Accept": "application/json, text/event-stream",
		},
	}
}

// HTTPTransport implements MCP over HTTP
type HTTPTransport struct {
	config Config
	client jsonrpc.RPCClient
	http   *SessionAwareHTTPClient
}

// NewHTTPTransport creates a new HTTP transport with options
func NewHTTPTransport(serverURL string, opts ...Option) *HTTPTransport {
	config := defaultConfig()
	config.ServerURL = serverURL

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	httpClient := NewSessionAwareHTTPClient(config.Timeout)

	// Create JSON-RPC client with session-aware HTTP client
	rpcClient := jsonrpc.NewClientWithOpts(config.ServerURL, &jsonrpc.RPCClientOpts{
		HTTPClient: httpClient,
		CustomHeaders: func() map[string]string {
			headers := map[string]string{
				"Accept": "application/json, text/event-stream",
			}
			// Add custom headers
			for k, v := range config.CustomHeaders {
				headers[k] = v
			}
			return headers
		}(),
	})

	return &HTTPTransport{
		config: *config,
		client: rpcClient,
		http:   httpClient,
	}
}

// NewHTTPTransportWithConfig creates a new HTTP transport with config (legacy)
// Deprecated: Use NewHTTPTransport with options instead
func NewHTTPTransportWithConfig(config Config) *HTTPTransport {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := NewSessionAwareHTTPClient(config.Timeout)

	// Create JSON-RPC client with session-aware HTTP client
	rpcClient := jsonrpc.NewClientWithOpts(config.ServerURL, &jsonrpc.RPCClientOpts{
		HTTPClient: httpClient,
		CustomHeaders: func() map[string]string {
			headers := map[string]string{
				"Accept": "application/json, text/event-stream",
			}
			// Add custom headers
			for k, v := range config.CustomHeaders {
				headers[k] = v
			}
			return headers
		}(),
	})

	return &HTTPTransport{
		config: config,
		client: rpcClient,
		http:   httpClient,
	}
}

// Call makes a JSON-RPC call
func (t *HTTPTransport) Call(ctx context.Context, result interface{}, method string, params ...interface{}) error {
	if len(params) == 0 {
		return t.client.CallFor(ctx, result, method)
	}
	return t.client.CallFor(ctx, result, method, params[0])
}

// CallRaw makes a JSON-RPC call and returns the raw response
func (t *HTTPTransport) CallRaw(ctx context.Context, method string, params interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := t.Call(ctx, &result, method, params)
	return result, err
}

// GetSessionID returns the current session ID
func (t *HTTPTransport) GetSessionID() string {
	return t.http.GetSessionID()
}

// Close closes the transport (no-op for HTTP)
func (t *HTTPTransport) Close() error {
	return nil
}

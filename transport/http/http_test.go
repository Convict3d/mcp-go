package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewSessionAwareHTTPClient(t *testing.T) {
	timeout := 45 * time.Second
	client := NewSessionAwareHTTPClient(timeout)

	if client == nil {
		t.Fatal("NewSessionAwareHTTPClient returned nil")
	}

	if client.client.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.client.Timeout)
	}

	// Initially should have no session ID
	sessionID := client.GetSessionID()
	if sessionID != "" {
		t.Errorf("Expected empty session ID initially, got: %s", sessionID)
	}
}

func TestSessionAwareHTTPClient_SessionManagement(t *testing.T) {
	// Create a test server that returns a session ID
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if session ID was sent
		if sessionID := r.Header.Get("Mcp-Session-Id"); sessionID != "" {
			w.Header().Set("Mcp-Session-Id", sessionID)
		} else {
			// Return a new session ID
			w.Header().Set("Mcp-Session-Id", "test-session-123")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "ok"}`))
	}))
	defer server.Close()

	client := NewSessionAwareHTTPClient(30 * time.Second)

	// Make first request - should get session ID
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	resp.Body.Close()

	// Check that session ID was stored
	sessionID := client.GetSessionID()
	if sessionID != "test-session-123" {
		t.Errorf("Expected session ID 'test-session-123', got '%s'", sessionID)
	}

	// Make second request - should send session ID
	req2, _ := http.NewRequest("GET", server.URL, nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	resp2.Body.Close()
}

func TestSessionAwareHTTPClient_SSEParsing(t *testing.T) {
	// Create a test server that returns SSE format
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data: {\"result\": \"test\"}\n\ndata: [DONE]\n\n"))
	}))
	defer server.Close()

	client := NewSessionAwareHTTPClient(30 * time.Second)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the parsed response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := `{"result": "test"}`
	if string(body) != expected {
		t.Errorf("Expected parsed body '%s', got '%s'", expected, string(body))
	}
}

func TestParseSSEResponse(t *testing.T) {
	client := &SessionAwareHTTPClient{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single data line",
			input:    "data: {\"test\": \"value\"}\n\n",
			expected: `{"test": "value"}`,
		},
		{
			name:     "multiple data lines",
			input:    "data: {\"line1\": \"value1\"}\n\ndata: {\"line2\": \"value2\"}\n\n",
			expected: `{"line1": "value1"}{"line2": "value2"}`,
		},
		{
			name:     "with DONE marker",
			input:    "data: {\"test\": \"value\"}\n\ndata: [DONE]\n\n",
			expected: `{"test": "value"}`,
		},
		{
			name:     "empty data lines",
			input:    "data: \n\ndata: {\"test\": \"value\"}\n\ndata: \n\n",
			expected: `{"test": "value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := io.NopCloser(strings.NewReader(tt.input))
			result, err := client.parseSSEResponse(reader)
			if err != nil {
				t.Fatalf("parseSSEResponse failed: %v", err)
			}

			body, err := io.ReadAll(result)
			if err != nil {
				t.Fatalf("Failed to read result: %v", err)
			}

			if string(body) != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, string(body))
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", config.Timeout)
	}

	if config.CustomHeaders == nil {
		t.Error("Expected default custom headers to be initialized")
	}

	if accept, exists := config.CustomHeaders["Accept"]; !exists || accept != "application/json, text/event-stream" {
		t.Error("Expected default Accept header for SSE")
	}
}

func TestHTTPTransportOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		checkFn  func(*Config) bool
		expected string
	}{
		{
			name:   "WithTimeout",
			option: WithTimeout(60 * time.Second),
			checkFn: func(c *Config) bool {
				return c.Timeout == 60*time.Second
			},
			expected: "timeout should be set",
		},
		{
			name:   "WithHeader",
			option: WithHeader("X-Custom", "value"),
			checkFn: func(c *Config) bool {
				return c.CustomHeaders["X-Custom"] == "value"
			},
			expected: "custom header should be set",
		},
		{
			name: "WithCustomHeaders",
			option: WithCustomHeaders(map[string]string{
				"X-Test1": "value1",
				"X-Test2": "value2",
			}),
			checkFn: func(c *Config) bool {
				return c.CustomHeaders["X-Test1"] == "value1" && c.CustomHeaders["X-Test2"] == "value2"
			},
			expected: "custom headers should be set",
		},
		{
			name:   "WithSSESupport",
			option: WithSSESupport(),
			checkFn: func(c *Config) bool {
				return c.CustomHeaders["Accept"] == "application/json, text/event-stream"
			},
			expected: "SSE headers should be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := defaultConfig()
			tt.option(config)

			if !tt.checkFn(config) {
				t.Error(tt.expected)
			}
		})
	}
}

func TestNewHTTPTransport(t *testing.T) {
	tests := []struct {
		name      string
		serverURL string
		opts      []Option
		wantErr   bool
	}{
		{
			name:      "basic transport creation",
			serverURL: "http://localhost:9831/mcp",
			opts:      nil,
			wantErr:   false,
		},
		{
			name:      "transport with options",
			serverURL: "http://localhost:9831/mcp",
			opts: []Option{
				WithTimeout(45 * time.Second),
				WithHeader("X-Test", "test"),
				WithSSESupport(),
			},
			wantErr: false,
		},
		{
			name:      "transport with custom headers",
			serverURL: "http://localhost:9831/mcp",
			opts: []Option{
				WithCustomHeaders(map[string]string{
					"Authorization": "Bearer token",
					"X-API-Key":     "secret",
				}),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := NewHTTPTransport(tt.serverURL, tt.opts...)

			if transport == nil {
				t.Fatal("NewHTTPTransport returned nil")
			}

			if transport.config.ServerURL != tt.serverURL {
				t.Errorf("Expected ServerURL %s, got %s", tt.serverURL, transport.config.ServerURL)
			}

			// Clean up
			transport.Close()
		})
	}
}

func TestNewHTTPTransportWithConfig(t *testing.T) {
	config := Config{
		ServerURL: "http://localhost:9831/mcp",
		Timeout:   45 * time.Second,
		CustomHeaders: map[string]string{
			"X-Test": "value",
		},
	}

	transport := NewHTTPTransportWithConfig(config)

	if transport == nil {
		t.Fatal("NewHTTPTransportWithConfig returned nil")
	}

	if transport.config.ServerURL != config.ServerURL {
		t.Errorf("Expected ServerURL %s, got %s", config.ServerURL, transport.config.ServerURL)
	}

	if transport.config.Timeout != config.Timeout {
		t.Errorf("Expected Timeout %v, got %v", config.Timeout, transport.config.Timeout)
	}

	transport.Close()
}

func TestHTTPTransportClose(t *testing.T) {
	transport := NewHTTPTransport("http://localhost:9831/mcp")

	err := transport.Close()
	if err != nil {
		t.Errorf("Close should not return an error, got: %v", err)
	}
}

func TestHTTPTransportGetSessionID(t *testing.T) {
	transport := NewHTTPTransport("http://localhost:9831/mcp")
	defer transport.Close()

	sessionID := transport.GetSessionID()
	// Should return empty string initially
	if sessionID != "" {
		t.Errorf("Expected empty session ID initially, got: %s", sessionID)
	}
}

func TestHTTPTransportCall(t *testing.T) {
	// Create a test server that responds to JSON-RPC calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc": "2.0", "result": {"test": "value"}, "id": 1}`))
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	defer transport.Close()

	var result map[string]interface{}
	err := transport.Call(context.Background(), &result, "test/method")
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	if result["test"] != "value" {
		t.Errorf("Expected result['test'] = 'value', got '%v'", result["test"])
	}
}

func TestHTTPTransportCallRaw(t *testing.T) {
	// Create a test server that responds to JSON-RPC calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc": "2.0", "result": {"raw": "data"}, "id": 1}`))
	}))
	defer server.Close()

	transport := NewHTTPTransport(server.URL)
	defer transport.Close()

	result, err := transport.CallRaw(context.Background(), "test/method", nil)
	if err != nil {
		t.Fatalf("CallRaw failed: %v", err)
	}

	if result["raw"] != "data" {
		t.Errorf("Expected result['raw'] = 'data', got '%v'", result["raw"])
	}
}

package stdio

import (
	"os"
	"testing"
	"time"
)

func TestNewTransport(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "basic echo command",
			command: "echo",
			args:    []string{"hello"},
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "with timeout option",
			command: "echo",
			args:    []string{"hello"},
			opts:    []Option{WithTimeout(10 * time.Second)},
			wantErr: false,
		},
		{
			name:    "with custom env",
			command: "echo",
			args:    []string{"hello"},
			opts:    []Option{WithEnv([]string{"TEST=value"})},
			wantErr: false,
		},
		{
			name:    "with working directory",
			command: "echo",
			args:    []string{"hello"},
			opts:    []Option{WithWorkingDir("/tmp")},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := NewTransport(tt.command, tt.args, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && transport == nil {
				t.Error("NewTransport() returned nil transport")
			}

			if transport != nil {
				transport.Close()
			}
		})
	}
}

func TestTransportOptions(t *testing.T) {
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
			name:   "WithEnv",
			option: WithEnv([]string{"TEST=value", "DEBUG=true"}),
			checkFn: func(c *Config) bool {
				return len(c.Env) == 2 && c.Env[0] == "TEST=value" && c.Env[1] == "DEBUG=true"
			},
			expected: "env should be set",
		},
		{
			name:   "WithWorkingDir",
			option: WithWorkingDir("/tmp"),
			checkFn: func(c *Config) bool {
				return c.WorkingDir == "/tmp"
			},
			expected: "working directory should be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{}
			tt.option(config)

			if !tt.checkFn(config) {
				t.Error(tt.expected)
			}
		})
	}
}

func TestTransportClose(t *testing.T) {
	// Use a simple command that exits quickly
	transport, err := NewTransport("echo", []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}

	// Give it a moment to start
	time.Sleep(50 * time.Millisecond)

	err = transport.Close()
	// For echo commands, broken pipe errors are expected when the process exits quickly
	if err != nil && err.Error() != "signal: broken pipe" {
		t.Errorf("Close() returned unexpected error: %v", err)
	}

	// Calling Close again should not cause issues
	err = transport.Close()
	if err != nil {
		t.Errorf("Second Close() should not return error: %v", err)
	}
}

func TestTransportGetSessionID(t *testing.T) {
	transport, err := NewTransport("echo", []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}
	defer transport.Close()

	sessionID := transport.GetSessionID()
	// Should return empty string for stdio transport
	if sessionID != "" {
		t.Errorf("Expected empty session ID for stdio transport, got: %s", sessionID)
	}
}

func TestTransportWithInvalidCommand(t *testing.T) {
	// Try to create transport with non-existent command
	transport, err := NewTransport("/non/existent/command", []string{})
	if err == nil {
		t.Error("Expected error for non-existent command")
		if transport != nil {
			transport.Close()
		}
	}
}

func TestTransportRequestIDGeneration(t *testing.T) {
	transport, err := NewTransport("echo", []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}
	defer transport.Close()

	// Test that request IDs are generated uniquely
	id1 := transport.generateRequestID()
	id2 := transport.generateRequestID()

	if id1 == id2 {
		t.Error("Request IDs should be unique")
	}

	if id1 <= 0 || id2 <= 0 {
		t.Error("Request IDs should be positive integers")
	}
}

func TestTransportMessageSending(t *testing.T) {
	// Test with a command that can handle JSON input/output
	// We'll use 'cat' which will echo back what we send
	if _, err := os.Stat("/bin/cat"); os.IsNotExist(err) {
		t.Skip("cat command not available")
	}

	transport, err := NewTransport("cat", []string{})
	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}
	defer transport.Close()

	// Give the process a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test sending a message
	message := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "test/method",
		"id":      1,
	}

	err = transport.sendMessage(message)
	if err != nil {
		t.Errorf("sendMessage() failed: %v", err)
	}
}

func TestTransportNotificationHandling(t *testing.T) {
	transport, err := NewTransport("echo", []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}
	defer transport.Close()

	// Test notification handling
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/message",
		"params": map[string]interface{}{
			"level": "info",
			"text":  "Test notification",
		},
	}

	// This should not panic or error
	transport.handleNotification(notification)
}

func TestTransportRequestHandling(t *testing.T) {
	transport, err := NewTransport("echo", []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}
	defer transport.Close()

	// Test request handling
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "ping",
		"id":      1,
	}

	// This should not panic or error
	transport.handleServerRequest(request, 1)
}

func TestTransportConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "empty command",
			command: "",
			args:    []string{},
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "valid command",
			command: "echo",
			args:    []string{"test"},
			opts:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := NewTransport(tt.command, tt.args, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
			}

			if transport != nil {
				transport.Close()
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", config.Timeout)
	}

	if config.Env == nil {
		t.Error("Expected default env to be initialized")
	}

	if config.Command != "" {
		t.Errorf("Expected empty default command, got %s", config.Command)
	}
}

func TestNewTransportFromOS(t *testing.T) {
	transport, err := NewTransportFromOS()
	if err != nil {
		t.Fatalf("NewTransportFromOS() failed: %v", err)
	}
	defer transport.Close()

	if transport == nil {
		t.Error("NewTransportFromOS() returned nil transport")
	}

	// Should have no command since it uses OS streams
	if transport.cmd != nil {
		t.Error("Expected no subprocess when using OS streams")
	}
}

func TestTransportSetHandlers(t *testing.T) {
	transport, err := NewTransport("echo", []string{"test"})
	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}
	defer transport.Close()

	// Test setting notification handler
	notificationCalled := false
	transport.SetNotificationHandler(func(method string, params interface{}) {
		notificationCalled = true
	})

	// Test setting request handler
	requestCalled := false
	transport.SetRequestHandler(func(method string, params interface{}) (interface{}, error) {
		requestCalled = true
		return "response", nil
	})

	// Test notification
	transport.handleNotification(map[string]interface{}{
		"method": "test",
		"params": nil,
	})

	if !notificationCalled {
		t.Error("Notification handler was not called")
	}

	// Test request
	transport.handleServerRequest(map[string]interface{}{
		"method": "test",
		"params": nil,
	}, 1)

	if !requestCalled {
		t.Error("Request handler was not called")
	}
}

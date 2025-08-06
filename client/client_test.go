package client

import (
	"context"
	"testing"
	"time"

	"github.com/Convict3d/mcp-go/types"
)

// MockTransport implements the Transport interface for testing
type MockTransport struct {
	callFunc    func(ctx context.Context, result interface{}, method string, params ...interface{}) error
	callRawFunc func(ctx context.Context, method string, params interface{}) (map[string]interface{}, error)
	sessionID   string
	closed      bool
}

func (m *MockTransport) Call(ctx context.Context, result interface{}, method string, params ...interface{}) error {
	if m.callFunc != nil {
		return m.callFunc(ctx, result, method, params...)
	}
	return nil
}

func (m *MockTransport) CallRaw(ctx context.Context, method string, params interface{}) (map[string]interface{}, error) {
	if m.callRawFunc != nil {
		return m.callRawFunc(ctx, method, params)
	}
	return make(map[string]interface{}), nil
}

func (m *MockTransport) GetSessionID() string {
	return m.sessionID
}

func (m *MockTransport) Close() error {
	m.closed = true
	return nil
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "basic client creation",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "client with options",
			opts: []Option{
				WithClientInfo("test-client", "1.0.0"),
				WithTimeout(45 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "client with custom transport",
			opts: []Option{
				WithTransport(&MockTransport{}),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.opts...)

			if client == nil {
				t.Fatal("NewClient returned nil")
			}

			if client.config == nil {
				t.Fatal("Client config is nil")
			}

			// Clean up
			if client.transport != nil {
				client.Close()
			}
		})
	}
}

func TestNewSimpleClient(t *testing.T) {
	client := NewSimpleClient()

	if client == nil {
		t.Fatal("NewSimpleClient returned nil")
	}

	// Clean up
	if client.transport != nil {
		client.Close()
	}
}

func TestClientOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		checkFn  func(*Config) bool
		expected string
	}{
		{
			name:   "WithClientInfo",
			option: WithClientInfo("test-app", "2.0.0"),
			checkFn: func(c *Config) bool {
				return c.ClientName == "test-app" && c.ClientVersion == "2.0.0"
			},
			expected: "client info should be set",
		},
		{
			name:   "WithTimeout",
			option: WithTimeout(60 * time.Second),
			checkFn: func(c *Config) bool {
				return c.Timeout == 60*time.Second
			},
			expected: "timeout should be set",
		},
		{
			name:   "WithTransport",
			option: WithTransport(&MockTransport{}),
			checkFn: func(c *Config) bool {
				return c.Transport != nil
			},
			expected: "transport should be set",
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

func TestClientInitialize(t *testing.T) {
	mockTransport := &MockTransport{
		callFunc: func(ctx context.Context, result interface{}, method string, params ...interface{}) error {
			if method == "initialize" {
				// Simulate successful initialization
				if initResult, ok := result.(*types.InitializeResult); ok {
					initResult.ServerInfo = types.Implementation{
						Name:    "test-server",
						Version: "1.0.0",
					}
					initResult.Capabilities = types.ServerCapabilities{
						Tools: &types.ToolsCapability{
							ListChanged: true,
						},
					}
				}
			}
			return nil
		},
	}

	client := NewClient(
		WithTransport(mockTransport),
		WithClientInfo("test-client", "1.0.0"),
	)
	defer client.Close()

	err := client.Initialize(types.LatestProtocolVersion)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Check server info was set
	serverInfo := client.GetServerInfo()
	if serverInfo == nil {
		t.Fatal("Server info not set after initialization")
	}

	if serverInfo.Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", serverInfo.Name)
	}

	// Check capabilities were set
	capabilities := client.GetCapabilities()
	if capabilities == nil {
		t.Fatal("Capabilities not set after initialization")
	}

	if !client.HasTools() {
		t.Error("Client should have tools capability")
	}
}

func TestClientCapabilityChecks(t *testing.T) {
	tests := []struct {
		name         string
		capabilities *types.ServerCapabilities
		hasTools     bool
		hasResources bool
		hasPrompts   bool
	}{
		{
			name:         "no capabilities",
			capabilities: nil,
			hasTools:     false,
			hasResources: false,
			hasPrompts:   false,
		},
		{
			name: "tools only",
			capabilities: &types.ServerCapabilities{
				Tools: &types.ToolsCapability{},
			},
			hasTools:     true,
			hasResources: false,
			hasPrompts:   false,
		},
		{
			name: "all capabilities",
			capabilities: &types.ServerCapabilities{
				Tools:     &types.ToolsCapability{},
				Resources: &types.ResourcesCapability{},
				Prompts:   &types.PromptsCapability{},
			},
			hasTools:     true,
			hasResources: true,
			hasPrompts:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(WithTransport(&MockTransport{}))
			defer client.Close()

			client.capabilities = tt.capabilities

			if client.HasTools() != tt.hasTools {
				t.Errorf("HasTools() = %v, want %v", client.HasTools(), tt.hasTools)
			}

			if client.HasResources() != tt.hasResources {
				t.Errorf("HasResources() = %v, want %v", client.HasResources(), tt.hasResources)
			}

			if client.HasPrompts() != tt.hasPrompts {
				t.Errorf("HasPrompts() = %v, want %v", client.HasPrompts(), tt.hasPrompts)
			}
		})
	}
}

func TestClientListTools(t *testing.T) {
	mockTools := []types.Tool{
		{
			BaseMetadata: types.BaseMetadata{
				Name:        "calculator",
				Description: "A simple calculator",
			},
			InputSchema: types.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"operation": map[string]interface{}{"type": "string"},
					"a":         map[string]interface{}{"type": "number"},
					"b":         map[string]interface{}{"type": "number"},
				},
			},
		},
	}

	mockTransport := &MockTransport{
		callFunc: func(ctx context.Context, result interface{}, method string, params ...interface{}) error {
			if method == "tools/list" {
				if listResult, ok := result.(*types.ListToolsResult); ok {
					listResult.Tools = mockTools
				}
			}
			return nil
		},
	}

	client := NewClient(WithTransport(mockTransport))
	defer client.Close()

	// Set capabilities to indicate tools are supported
	client.capabilities = &types.ServerCapabilities{
		Tools: &types.ToolsCapability{},
	}

	tools, err := client.ListTools()
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	if len(tools) != 1 {
		t.Fatalf("Expected 1 tool, got %d", len(tools))
	}

	if tools[0].Name != "calculator" {
		t.Errorf("Expected tool name 'calculator', got '%s'", tools[0].Name)
	}
}

func TestClientCallTool(t *testing.T) {
	mockResult := &types.CallToolResult{
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": "Result: 15",
			},
		},
	}

	mockTransport := &MockTransport{
		callFunc: func(ctx context.Context, result interface{}, method string, params ...interface{}) error {
			if method == "tools/call" {
				if callResult, ok := result.(*types.CallToolResult); ok {
					*callResult = *mockResult
				}
			}
			return nil
		},
	}

	client := NewClient(WithTransport(mockTransport))
	defer client.Close()

	// Set capabilities to indicate tools are supported
	client.capabilities = &types.ServerCapabilities{
		Tools: &types.ToolsCapability{},
	}

	result, err := client.CallTool("calculator", map[string]interface{}{
		"operation": "multiply",
		"a":         3,
		"b":         5,
	})
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if result == nil {
		t.Fatal("CallTool returned nil result")
	}

	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(result.Content))
	}
}

func TestClientListResources(t *testing.T) {
	mockResources := []types.Resource{
		{
			BaseMetadata: types.BaseMetadata{
				Name:        "test-file",
				Description: "A test file resource",
			},
			URI:      "file://test.txt",
			MimeType: "text/plain",
		},
	}

	mockTransport := &MockTransport{
		callFunc: func(ctx context.Context, result interface{}, method string, params ...interface{}) error {
			if method == "resources/list" {
				if listResult, ok := result.(*types.ListResourcesResult); ok {
					listResult.Resources = mockResources
				}
			}
			return nil
		},
	}

	client := NewClient(WithTransport(mockTransport))
	defer client.Close()

	// Set capabilities to indicate resources are supported
	client.capabilities = &types.ServerCapabilities{
		Resources: &types.ResourcesCapability{},
	}

	resources, err := client.ListResources()
	if err != nil {
		t.Fatalf("ListResources failed: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	if resources[0].Name != "test-file" {
		t.Errorf("Expected resource name 'test-file', got '%s'", resources[0].Name)
	}
}

func TestClientClose(t *testing.T) {
	mockTransport := &MockTransport{}
	client := NewClient(WithTransport(mockTransport))

	err := client.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if !mockTransport.closed {
		t.Error("Transport was not closed")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()

	if config.ClientName != "go-mcp-client" {
		t.Errorf("Expected default client name 'go-mcp-client', got '%s'", config.ClientName)
	}

	if config.ClientVersion != "1.0.0" {
		t.Errorf("Expected default client version '1.0.0', got '%s'", config.ClientVersion)
	}

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

func TestClientWithoutTransport(t *testing.T) {
	// Test creating a client without a transport (should handle gracefully)
	client := NewClient(WithClientInfo("test", "1.0.0"))

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.transport != nil {
		t.Error("Expected nil transport when none provided")
	}

	// Should not panic when calling methods without transport
	serverInfo := client.GetServerInfo()
	if serverInfo != nil {
		t.Error("Expected nil server info when not initialized")
	}

	capabilities := client.GetCapabilities()
	if capabilities != nil {
		t.Error("Expected nil capabilities when not initialized")
	}
}

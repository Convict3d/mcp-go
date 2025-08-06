// Package client provides a high-level MCP client implementation
package client

import (
	"context"
	"time"

	"github.com/Convict3d/mcp-go/transport"
	"github.com/Convict3d/mcp-go/types"
)

// Client represents a high-level MCP client
type Client struct {
	transport    transport.Transport
	ctx          context.Context
	config       *Config
	serverInfo   *types.Implementation
	capabilities *types.ServerCapabilities
}

// Config holds configuration for the MCP client
type Config struct {
	ServerURL     string
	ClientName    string
	ClientVersion string
	Timeout       time.Duration
	CustomHeaders map[string]string
	Transport     transport.Transport // Custom transport
}

// Option defines a function that configures the client
type Option func(*Config)

// WithClientInfo sets the client name and version
func WithClientInfo(name, version string) Option {
	return func(c *Config) {
		c.ClientName = name
		c.ClientVersion = version
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithTransport sets a custom transport
func WithTransport(t transport.Transport) Option {
	return func(c *Config) {
		c.Transport = t
	}
}

// WithContext sets a custom context (advanced usage)
func WithContext(ctx context.Context) Option {
	return func(c *Config) {
		// Store context in config for later use in NewClient
		// We'll handle this in the client creation
	}
}

// defaultConfig returns a default client configuration
func defaultConfig() *Config {
	return &Config{
		ClientName:    "go-mcp-client",
		ClientVersion: "1.0.0",
		Timeout:       30 * time.Second,
		CustomHeaders: map[string]string{
			"Accept": "application/json, text/event-stream",
		},
	}
}

// NewClient creates a new high-level MCP client with options
func NewClient(opts ...Option) *Client {
	config := defaultConfig()

	// Apply all options
	for _, opt := range opts {
		opt(config)
	}

	return &Client{
		transport: config.Transport,
		ctx:       context.Background(),
		config:    config,
	}
}

// NewSimpleClient creates a client with minimal configuration for quick setup
func NewSimpleClient() *Client {
	return NewClient()
}

// Initialize connects to the MCP server and performs the handshake
func (c *Client) Initialize(protocolVersion string) error {
	params := struct {
		ProtocolVersion string                   `json:"protocolVersion"`
		Capabilities    types.ClientCapabilities `json:"capabilities"`
		ClientInfo      types.Implementation     `json:"clientInfo"`
	}{
		ProtocolVersion: protocolVersion,
		Capabilities:    types.ClientCapabilities{},
		ClientInfo: types.Implementation{
			Name:    c.config.ClientName,
			Version: c.config.ClientVersion,
		},
	}

	var result types.InitializeResult
	err := c.transport.Call(c.ctx, &result, "initialize", params)
	if err != nil {
		return err
	}

	c.serverInfo = &result.ServerInfo
	c.capabilities = &result.Capabilities

	// Note: Some MCP servers don't require the initialized notification
	// If needed, uncomment the following lines:
	// var empty interface{}
	// return c.transport.Call(c.ctx, &empty, "initialized", struct{}{})

	return nil
}

// GetServerInfo returns information about the connected server
func (c *Client) GetServerInfo() *types.Implementation {
	return c.serverInfo
}

// GetCapabilities returns the server's capabilities
func (c *Client) GetCapabilities() *types.ServerCapabilities {
	return c.capabilities
}

// GetSessionID returns the current session ID
func (c *Client) GetSessionID() string {
	return c.transport.GetSessionID()
}

// HasTools returns true if the server supports tools
func (c *Client) HasTools() bool {
	return c.capabilities != nil && c.capabilities.Tools != nil
}

// HasResources returns true if the server supports resources
func (c *Client) HasResources() bool {
	return c.capabilities != nil && c.capabilities.Resources != nil
}

// HasPrompts returns true if the server supports prompts
func (c *Client) HasPrompts() bool {
	return c.capabilities != nil && c.capabilities.Prompts != nil
}

// ListTools retrieves the list of available tools from the server
func (c *Client) ListTools() ([]types.Tool, error) {
	if !c.HasTools() {
		return nil, nil
	}

	var result types.ListToolsResult
	err := c.transport.Call(c.ctx, &result, "tools/list")
	if err != nil {
		return nil, err
	}

	return result.Tools, nil
}

// CallTool executes a tool with the given arguments
func (c *Client) CallTool(name string, arguments map[string]interface{}) (*types.CallToolResult, error) {
	if !c.HasTools() {
		return nil, nil
	}

	params := struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments,omitempty"`
		Meta      types.Meta             `json:"_meta,omitempty"`
	}{
		Name:      name,
		Arguments: arguments,
	}

	var result types.CallToolResult
	err := c.transport.Call(c.ctx, &result, "tools/call", params)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListResources retrieves the list of available resources from the server
func (c *Client) ListResources() ([]types.Resource, error) {
	if !c.HasResources() {
		return nil, nil
	}

	var result types.ListResourcesResult
	err := c.transport.Call(c.ctx, &result, "resources/list")
	if err != nil {
		return nil, err
	}

	return result.Resources, nil
}

// ReadResource reads the content of a specific resource
func (c *Client) ReadResource(uri string) (*types.ReadResourceResult, error) {
	if !c.HasResources() {
		return nil, nil
	}

	params := struct {
		URI  string     `json:"uri"`
		Meta types.Meta `json:"_meta,omitempty"`
	}{
		URI: uri,
	}

	var result types.ReadResourceResult
	err := c.transport.Call(c.ctx, &result, "resources/read", params)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListPrompts retrieves the list of available prompts from the server
func (c *Client) ListPrompts() ([]types.Prompt, error) {
	if !c.HasPrompts() {
		return nil, nil
	}

	var result types.ListPromptsResult
	err := c.transport.Call(c.ctx, &result, "prompts/list")
	if err != nil {
		return nil, err
	}

	return result.Prompts, nil
}

// GetPrompt retrieves a specific prompt with arguments
func (c *Client) GetPrompt(name string, arguments map[string]string) (*types.GetPromptResult, error) {
	if !c.HasPrompts() {
		return nil, nil
	}

	params := struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments,omitempty"`
		Meta      types.Meta        `json:"_meta,omitempty"`
	}{
		Name:      name,
		Arguments: arguments,
	}

	var result types.GetPromptResult
	err := c.transport.Call(c.ctx, &result, "prompts/get", params)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Close closes the client and cleans up resources
func (c *Client) Close() error {
	return c.transport.Close()
}

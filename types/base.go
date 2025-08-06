// Package types contains all MCP protocol type definitions
package types

import (
	"encoding/json"
)

// Protocol constants
const (
	LatestProtocolVersion = "2025-06-18"
	JSONRPCVersion        = "2.0"
)

// RequestID represents a uniquely identifying ID for a request in JSON-RPC
type RequestID interface{}

// ProgressToken is used to associate progress notifications with the original request
type ProgressToken interface{}

// Cursor represents an opaque token used for pagination
type Cursor string

// Meta provides additional metadata for MCP interactions
type Meta map[string]interface{}

// Base JSON-RPC message types

// JSONRPCMessage represents any valid JSON-RPC object
type JSONRPCMessage interface {
	GetJSONRPCVersion() string
}

// Request represents a base request structure
type Request struct {
	Method string        `json:"method"`
	Params RequestParams `json:"params,omitempty"`
}

// RequestParams contains parameters for requests
type RequestParams struct {
	Meta   Meta                   `json:"_meta,omitempty"`
	Fields map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for RequestParams
func (rp RequestParams) MarshalJSON() ([]byte, error) {
	// Start with the fields map
	result := make(map[string]interface{})
	for k, v := range rp.Fields {
		result[k] = v
	}

	// Add Meta if present
	if rp.Meta != nil {
		result["_meta"] = rp.Meta
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for RequestParams
func (rp *RequestParams) UnmarshalJSON(data []byte) error {
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Extract Meta if present
	if meta, ok := temp["_meta"]; ok {
		if metaMap, ok := meta.(map[string]interface{}); ok {
			rp.Meta = metaMap
		}
		delete(temp, "_meta")
	}

	// Store remaining fields
	rp.Fields = temp
	return nil
}

// Response represents a base response structure
type Response struct {
	Meta   Meta                   `json:"_meta,omitempty"`
	Result map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for Response
func (r Response) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	for k, v := range r.Result {
		result[k] = v
	}

	if r.Meta != nil {
		result["_meta"] = r.Meta
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for Response
func (r *Response) UnmarshalJSON(data []byte) error {
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if meta, ok := temp["_meta"]; ok {
		if metaMap, ok := meta.(map[string]interface{}); ok {
			r.Meta = metaMap
		}
		delete(temp, "_meta")
	}

	r.Result = temp
	return nil
}

// Notification represents a base notification structure
type Notification struct {
	Method string             `json:"method"`
	Params NotificationParams `json:"params,omitempty"`
}

// NotificationParams contains parameters for notifications
type NotificationParams struct {
	Meta   Meta                   `json:"_meta,omitempty"`
	Fields map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for NotificationParams
func (np NotificationParams) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	for k, v := range np.Fields {
		result[k] = v
	}

	if np.Meta != nil {
		result["_meta"] = np.Meta
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for NotificationParams
func (np *NotificationParams) UnmarshalJSON(data []byte) error {
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if meta, ok := temp["_meta"]; ok {
		if metaMap, ok := meta.(map[string]interface{}); ok {
			np.Meta = metaMap
		}
		delete(temp, "_meta")
	}

	np.Fields = temp
	return nil
}

// Base metadata structure
type BaseMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Role represents a role in a conversation
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// LoggingLevel represents different logging levels
type LoggingLevel string

const (
	LoggingLevelDebug     LoggingLevel = "debug"
	LoggingLevelInfo      LoggingLevel = "info"
	LoggingLevelNotice    LoggingLevel = "notice"
	LoggingLevelWarning   LoggingLevel = "warning"
	LoggingLevelError     LoggingLevel = "error"
	LoggingLevelCritical  LoggingLevel = "critical"
	LoggingLevelAlert     LoggingLevel = "alert"
	LoggingLevelEmergency LoggingLevel = "emergency"
)

// Content types
type ContentType = string

const (
	ContentTypeText         ContentType = "text"
	ContentTypeImage        ContentType = "image"
	ContentTypeAudio        ContentType = "audio"
	ContentTypeResourceLink ContentType = "resource_link"
	ContentTypeResource     ContentType = "resource"
)

// ContentBlock represents different types of content that can be sent
type ContentBlock interface {
	ContentType() string
}

// TextContent represents text content
type TextContent struct {
	Type        string       `json:"type"`
	Text        string       `json:"text"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// ContentType returns the content type for TextContent
func (tc TextContent) ContentType() string {
	return ContentTypeText
}

// ImageContent represents image content
type ImageContent struct {
	Type        string       `json:"type"`
	Data        string       `json:"data"` // base64 encoded
	MimeType    string       `json:"mimeType"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// ContentType returns the content type for ImageContent
func (ic ImageContent) ContentType() string {
	return ContentTypeImage
}

// AudioContent represents audio content
type AudioContent struct {
	Type        string       `json:"type"`
	Data        string       `json:"data"` // base64 encoded
	MimeType    string       `json:"mimeType"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// ContentType returns the content type for AudioContent
func (ac AudioContent) ContentType() string {
	return ContentTypeAudio
}

// ResourceLinkContent represents a link to a resource
type ResourceLinkContent struct {
	Type        string       `json:"type"`
	URI         string       `json:"uri"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	MimeType    string       `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// ContentType returns the content type for ResourceLinkContent
func (rlc ResourceLinkContent) ContentType() string {
	return ContentTypeResourceLink
}

// ResourceContent represents an embedded resource
type ResourceContent struct {
	Type        string       `json:"type"`
	Resource    *Resource    `json:"resource"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// ContentType returns the content type for ResourceContent
func (rc ResourceContent) ContentType() string {
	return ContentTypeResource
}

// Capabilities

// ClientCapabilities represents what the client supports
type ClientCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Sampling     *SamplingCapability    `json:"sampling,omitempty"`
	Roots        *RootsCapability       `json:"roots,omitempty"`
}

// ServerCapabilities represents what the server supports
type ServerCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Logging      *LoggingCapability     `json:"logging,omitempty"`
	Prompts      *PromptsCapability     `json:"prompts,omitempty"`
	Resources    *ResourcesCapability   `json:"resources,omitempty"`
	Tools        *ToolsCapability       `json:"tools,omitempty"`
}

// Individual capability structures
type SamplingCapability struct{}
type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}
type LoggingCapability struct{}
type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}
type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Initialize types

// Implementation represents information about the client or server implementation
type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeRequest is sent from client to server to initialize the connection
type InitializeRequest struct {
	Method string `json:"method"`
	Params struct {
		ProtocolVersion string             `json:"protocolVersion"`
		Capabilities    ClientCapabilities `json:"capabilities"`
		ClientInfo      Implementation     `json:"clientInfo"`
	} `json:"params"`
}

// InitializeResult is the server's response to initialization
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      Implementation     `json:"serverInfo"`
	Instructions    string             `json:"instructions,omitempty"`
}

// InitializedNotification is sent from client to server after initialization
type InitializedNotification struct {
	Method string `json:"method"`
}

// Ping types

// PingRequest is sent to test liveness/heartbeat
type PingRequest struct {
	Method string `json:"method"`
}

// PingResult is the response to a ping
type PingResult struct{}

// Progress types

// ProgressNotification provides updates on long-running operations
type ProgressNotification struct {
	Method string `json:"method"`
	Params struct {
		ProgressToken ProgressToken `json:"progressToken"`
		Progress      int           `json:"progress"`
		Total         *int          `json:"total,omitempty"`
	} `json:"params"`
}

// Pagination types

// PaginatedRequest represents a request that supports pagination
type PaginatedRequest struct {
	Cursor *Cursor `json:"cursor,omitempty"`
}

// PaginatedResult represents a paginated response
type PaginatedResult struct {
	NextCursor *Cursor `json:"nextCursor,omitempty"`
}

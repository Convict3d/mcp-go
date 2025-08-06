package types

import (
	"encoding/json"
	"testing"
)

func TestRequestParams_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		params   RequestParams
		expected string
	}{
		{
			name: "empty params",
			params: RequestParams{
				Fields: map[string]interface{}{},
			},
			expected: `{}`,
		},
		{
			name: "with fields only",
			params: RequestParams{
				Fields: map[string]interface{}{
					"key1": "value1",
					"key2": 42,
				},
			},
			expected: `{"key1":"value1","key2":42}`,
		},
		{
			name: "with meta only",
			params: RequestParams{
				Meta: Meta{
					"author": "test",
				},
				Fields: map[string]interface{}{},
			},
			expected: `{"_meta":{"author":"test"}}`,
		},
		{
			name: "with both fields and meta",
			params: RequestParams{
				Meta: Meta{
					"version": "1.0",
				},
				Fields: map[string]interface{}{
					"data": "test",
				},
			},
			expected: `{"_meta":{"version":"1.0"},"data":"test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// Parse both JSON strings to compare content regardless of order
			var expectedMap, resultMap map[string]interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedMap); err != nil {
				t.Fatalf("Failed to parse expected JSON: %v", err)
			}
			if err := json.Unmarshal(result, &resultMap); err != nil {
				t.Fatalf("Failed to parse result JSON: %v", err)
			}

			if !equalMaps(expectedMap, resultMap) {
				t.Errorf("Expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestRequestParams_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RequestParams
		wantErr  bool
	}{
		{
			name:  "empty object",
			input: `{}`,
			expected: RequestParams{
				Fields: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:  "with fields only",
			input: `{"key1":"value1","key2":42}`,
			expected: RequestParams{
				Fields: map[string]interface{}{
					"key1": "value1",
					"key2": float64(42), // JSON numbers unmarshal as float64
				},
			},
			wantErr: false,
		},
		{
			name:  "with meta only",
			input: `{"_meta":{"author":"test"}}`,
			expected: RequestParams{
				Meta: Meta{
					"author": "test",
				},
				Fields: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:  "with both fields and meta",
			input: `{"_meta":{"version":"1.0"},"data":"test"}`,
			expected: RequestParams{
				Meta: Meta{
					"version": "1.0",
				},
				Fields: map[string]interface{}{
					"data": "test",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params RequestParams
			err := json.Unmarshal([]byte(tt.input), &params)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !equalMaps(tt.expected.Meta, params.Meta) {
				t.Errorf("Meta mismatch. Expected %v, got %v", tt.expected.Meta, params.Meta)
			}

			if !equalMaps(tt.expected.Fields, params.Fields) {
				t.Errorf("Fields mismatch. Expected %v, got %v", tt.expected.Fields, params.Fields)
			}
		})
	}
}

func TestResponse_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		response Response
		expected string
	}{
		{
			name: "empty response",
			response: Response{
				Result: map[string]interface{}{},
			},
			expected: `{}`,
		},
		{
			name: "with result only",
			response: Response{
				Result: map[string]interface{}{
					"status": "success",
					"count":  10,
				},
			},
			expected: `{"status":"success","count":10}`,
		},
		{
			name: "with meta only",
			response: Response{
				Meta: Meta{
					"timestamp": "2023-01-01",
				},
				Result: map[string]interface{}{},
			},
			expected: `{"_meta":{"timestamp":"2023-01-01"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var expectedMap, resultMap map[string]interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedMap); err != nil {
				t.Fatalf("Failed to parse expected JSON: %v", err)
			}
			if err := json.Unmarshal(result, &resultMap); err != nil {
				t.Fatalf("Failed to parse result JSON: %v", err)
			}

			if !equalMaps(expectedMap, resultMap) {
				t.Errorf("Expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestResponse_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Response
		wantErr  bool
	}{
		{
			name:  "empty object",
			input: `{}`,
			expected: Response{
				Result: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:  "with result fields",
			input: `{"status":"success","count":10}`,
			expected: Response{
				Result: map[string]interface{}{
					"status": "success",
					"count":  float64(10),
				},
			},
			wantErr: false,
		},
		{
			name:  "with meta",
			input: `{"_meta":{"timestamp":"2023-01-01"},"status":"ok"}`,
			expected: Response{
				Meta: Meta{
					"timestamp": "2023-01-01",
				},
				Result: map[string]interface{}{
					"status": "ok",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response Response
			err := json.Unmarshal([]byte(tt.input), &response)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !equalMaps(tt.expected.Meta, response.Meta) {
				t.Errorf("Meta mismatch. Expected %v, got %v", tt.expected.Meta, response.Meta)
			}

			if !equalMaps(tt.expected.Result, response.Result) {
				t.Errorf("Result mismatch. Expected %v, got %v", tt.expected.Result, response.Result)
			}
		})
	}
}

func TestNotificationParams_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		params   NotificationParams
		expected string
	}{
		{
			name: "empty params",
			params: NotificationParams{
				Fields: map[string]interface{}{},
			},
			expected: `{}`,
		},
		{
			name: "with fields and meta",
			params: NotificationParams{
				Meta: Meta{
					"source": "server",
				},
				Fields: map[string]interface{}{
					"message": "Hello",
					"level":   "info",
				},
			},
			expected: `{"_meta":{"source":"server"},"message":"Hello","level":"info"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var expectedMap, resultMap map[string]interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedMap); err != nil {
				t.Fatalf("Failed to parse expected JSON: %v", err)
			}
			if err := json.Unmarshal(result, &resultMap); err != nil {
				t.Fatalf("Failed to parse result JSON: %v", err)
			}

			if !equalMaps(expectedMap, resultMap) {
				t.Errorf("Expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestContentBlocks(t *testing.T) {
	tests := []struct {
		name         string
		content      ContentBlock
		expectedType string
	}{
		{
			name: "text content",
			content: TextContent{
				Type: "text",
				Text: "Hello world",
			},
			expectedType: ContentTypeText,
		},
		{
			name: "image content",
			content: ImageContent{
				Type:     "image",
				Data:     "base64data",
				MimeType: "image/png",
			},
			expectedType: ContentTypeImage,
		},
		{
			name: "audio content",
			content: AudioContent{
				Type:     "audio",
				Data:     "base64audio",
				MimeType: "audio/mp3",
			},
			expectedType: ContentTypeAudio,
		},
		{
			name: "resource link content",
			content: ResourceLinkContent{
				Type: "resource_link",
				URI:  "https://example.com/resource",
				Name: "Example Resource",
			},
			expectedType: ContentTypeResourceLink,
		},
		{
			name: "resource content",
			content: ResourceContent{
				Type: "resource",
				Resource: &Resource{
					BaseMetadata: BaseMetadata{
						Name: "test",
					},
					URI: "file://test.txt",
				},
			},
			expectedType: ContentTypeResource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.content.ContentType() != tt.expectedType {
				t.Errorf("Expected content type %s, got %s", tt.expectedType, tt.content.ContentType())
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test protocol constants
	if LatestProtocolVersion != "2025-06-18" {
		t.Errorf("Expected protocol version 2025-06-18, got %s", LatestProtocolVersion)
	}

	if JSONRPCVersion != "2.0" {
		t.Errorf("Expected JSON-RPC version 2.0, got %s", JSONRPCVersion)
	}

	// Test role constants
	if RoleUser != "user" {
		t.Errorf("Expected user role to be 'user', got %s", RoleUser)
	}

	if RoleAssistant != "assistant" {
		t.Errorf("Expected assistant role to be 'assistant', got %s", RoleAssistant)
	}

	// Test logging level constants
	levels := []LoggingLevel{
		LoggingLevelDebug,
		LoggingLevelInfo,
		LoggingLevelNotice,
		LoggingLevelWarning,
		LoggingLevelError,
		LoggingLevelCritical,
		LoggingLevelAlert,
		LoggingLevelEmergency,
	}

	expectedLevels := []string{
		"debug", "info", "notice", "warning",
		"error", "critical", "alert", "emergency",
	}

	for i, level := range levels {
		if string(level) != expectedLevels[i] {
			t.Errorf("Expected logging level %s, got %s", expectedLevels[i], string(level))
		}
	}

	// Test content type constants
	contentTypes := []ContentType{
		ContentTypeText,
		ContentTypeImage,
		ContentTypeAudio,
		ContentTypeResourceLink,
		ContentTypeResource,
	}

	expectedContentTypes := []string{
		"text", "image", "audio", "resource_link", "resource",
	}

	for i, ct := range contentTypes {
		if ct != expectedContentTypes[i] {
			t.Errorf("Expected content type %s, got %s", expectedContentTypes[i], ct)
		}
	}
}

func TestJSONSerialization(t *testing.T) {
	// Test complex structure serialization
	init := InitializeRequest{
		Method: "initialize",
		Params: struct {
			ProtocolVersion string             `json:"protocolVersion"`
			Capabilities    ClientCapabilities `json:"capabilities"`
			ClientInfo      Implementation     `json:"clientInfo"`
		}{
			ProtocolVersion: LatestProtocolVersion,
			Capabilities: ClientCapabilities{
				Sampling: &SamplingCapability{},
				Roots: &RootsCapability{
					ListChanged: true,
				},
				Experimental: map[string]interface{}{
					"feature1": true,
				},
			},
			ClientInfo: Implementation{
				Name:    "test-client",
				Version: "1.0.0",
			},
		},
	}

	data, err := json.Marshal(init)
	if err != nil {
		t.Fatalf("Failed to marshal InitializeRequest: %v", err)
	}

	var unmarshaled InitializeRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal InitializeRequest: %v", err)
	}

	if unmarshaled.Method != "initialize" {
		t.Errorf("Expected method 'initialize', got %s", unmarshaled.Method)
	}

	if unmarshaled.Params.ProtocolVersion != LatestProtocolVersion {
		t.Errorf("Expected protocol version %s, got %s", LatestProtocolVersion, unmarshaled.Params.ProtocolVersion)
	}

	if unmarshaled.Params.ClientInfo.Name != "test-client" {
		t.Errorf("Expected client name 'test-client', got %s", unmarshaled.Params.ClientInfo.Name)
	}
}

// Helper function to compare maps (handles nil cases)
func equalMaps(a, b map[string]interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return len(a) == 0 && len(b) == 0
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !equalValues(v, bv) {
			return false
		}
	}
	return true
}

// Helper function to compare values (handles different numeric types)
func equalValues(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle nested maps
	if mapA, okA := a.(map[string]interface{}); okA {
		if mapB, okB := b.(map[string]interface{}); okB {
			return equalMaps(mapA, mapB)
		}
		return false
	}

	return a == b
}

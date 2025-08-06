package types

import (
	"encoding/json"
	"testing"
)

func TestTool_JSONSerialization(t *testing.T) {
	tool := Tool{
		BaseMetadata: BaseMetadata{
			Name:        "test_tool",
			Description: "A test tool",
		},
		Description: "Detailed description",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"param1": map[string]interface{}{
					"type":        "string",
					"description": "First parameter",
				},
				"param2": map[string]interface{}{
					"type": "number",
				},
			},
			Required: []string{"param1"},
		},
		OutputSchema: &ToolOutputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"result": map[string]interface{}{
					"type": "string",
				},
			},
		},
		Annotations: &ToolAnnotations{
			Title:           "Test Tool",
			ReadOnlyHint:    true,
			DestructiveHint: false,
			IdempotentHint:  true,
			OpenWorldHint:   false,
		},
		Meta: Meta{
			"version": "1.0",
		},
	}

	// Test marshaling
	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("Failed to marshal Tool: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Tool
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Tool: %v", err)
	}

	// Verify core fields
	if unmarshaled.Name != "test_tool" {
		t.Errorf("Expected name 'test_tool', got %s", unmarshaled.Name)
	}

	if unmarshaled.Description != "Detailed description" {
		t.Errorf("Expected description 'Detailed description', got %s", unmarshaled.Description)
	}

	if unmarshaled.InputSchema.Type != "object" {
		t.Errorf("Expected input schema type 'object', got %s", unmarshaled.InputSchema.Type)
	}

	if len(unmarshaled.InputSchema.Required) != 1 || unmarshaled.InputSchema.Required[0] != "param1" {
		t.Errorf("Expected required fields ['param1'], got %v", unmarshaled.InputSchema.Required)
	}

	if unmarshaled.OutputSchema == nil {
		t.Error("Expected output schema to be present")
	} else if unmarshaled.OutputSchema.Type != "object" {
		t.Errorf("Expected output schema type 'object', got %s", unmarshaled.OutputSchema.Type)
	}

	if unmarshaled.Annotations == nil {
		t.Error("Expected annotations to be present")
	} else {
		if unmarshaled.Annotations.Title != "Test Tool" {
			t.Errorf("Expected title 'Test Tool', got %s", unmarshaled.Annotations.Title)
		}
		if !unmarshaled.Annotations.ReadOnlyHint {
			t.Error("Expected ReadOnlyHint to be true")
		}
		if unmarshaled.Annotations.DestructiveHint {
			t.Error("Expected DestructiveHint to be false")
		}
		if !unmarshaled.Annotations.IdempotentHint {
			t.Error("Expected IdempotentHint to be true")
		}
		if unmarshaled.Annotations.OpenWorldHint {
			t.Error("Expected OpenWorldHint to be false")
		}
	}

	if unmarshaled.Meta == nil {
		t.Error("Expected meta to be present")
	} else if unmarshaled.Meta["version"] != "1.0" {
		t.Errorf("Expected meta version '1.0', got %v", unmarshaled.Meta["version"])
	}
}

func TestToolInputSchema(t *testing.T) {
	schema := ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name parameter",
				"minLength":   1,
			},
			"age": map[string]interface{}{
				"type":    "integer",
				"minimum": 0,
				"maximum": 150,
			},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		Required: []string{"name"},
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal ToolInputSchema: %v", err)
	}

	var unmarshaled ToolInputSchema
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ToolInputSchema: %v", err)
	}

	if unmarshaled.Type != "object" {
		t.Errorf("Expected type 'object', got %s", unmarshaled.Type)
	}

	if len(unmarshaled.Properties) != 3 {
		t.Errorf("Expected 3 properties, got %d", len(unmarshaled.Properties))
	}

	if len(unmarshaled.Required) != 1 || unmarshaled.Required[0] != "name" {
		t.Errorf("Expected required ['name'], got %v", unmarshaled.Required)
	}

	// Check specific property
	nameProperty, exists := unmarshaled.Properties["name"]
	if !exists {
		t.Error("Expected 'name' property to exist")
	} else {
		nameProp := nameProperty.(map[string]interface{})
		if nameProp["type"] != "string" {
			t.Errorf("Expected name type 'string', got %v", nameProp["type"])
		}
		if nameProp["description"] != "Name parameter" {
			t.Errorf("Expected name description 'Name parameter', got %v", nameProp["description"])
		}
	}
}

func TestToolOutputSchema(t *testing.T) {
	schema := ToolOutputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{"success", "error"},
			},
			"data": map[string]interface{}{
				"type": "object",
			},
		},
		Required: []string{"status"},
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal ToolOutputSchema: %v", err)
	}

	var unmarshaled ToolOutputSchema
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ToolOutputSchema: %v", err)
	}

	if unmarshaled.Type != "object" {
		t.Errorf("Expected type 'object', got %s", unmarshaled.Type)
	}

	if len(unmarshaled.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(unmarshaled.Properties))
	}

	if len(unmarshaled.Required) != 1 || unmarshaled.Required[0] != "status" {
		t.Errorf("Expected required ['status'], got %v", unmarshaled.Required)
	}
}

func TestToolAnnotations(t *testing.T) {
	annotations := ToolAnnotations{
		Title:           "My Tool",
		ReadOnlyHint:    true,
		DestructiveHint: false,
		IdempotentHint:  true,
		OpenWorldHint:   false,
	}

	data, err := json.Marshal(annotations)
	if err != nil {
		t.Fatalf("Failed to marshal ToolAnnotations: %v", err)
	}

	var unmarshaled ToolAnnotations
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ToolAnnotations: %v", err)
	}

	if unmarshaled.Title != "My Tool" {
		t.Errorf("Expected title 'My Tool', got %s", unmarshaled.Title)
	}

	if !unmarshaled.ReadOnlyHint {
		t.Error("Expected ReadOnlyHint to be true")
	}

	if unmarshaled.DestructiveHint {
		t.Error("Expected DestructiveHint to be false")
	}

	if !unmarshaled.IdempotentHint {
		t.Error("Expected IdempotentHint to be true")
	}

	if unmarshaled.OpenWorldHint {
		t.Error("Expected OpenWorldHint to be false")
	}
}

func TestListToolsRequest(t *testing.T) {
	request := ListToolsRequest{
		Method: "tools/list",
		Params: struct {
			Cursor Cursor `json:"cursor,omitempty"`
		}{
			Cursor: "page_token_123",
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal ListToolsRequest: %v", err)
	}

	var unmarshaled ListToolsRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ListToolsRequest: %v", err)
	}

	if unmarshaled.Method != "tools/list" {
		t.Errorf("Expected method 'tools/list', got %s", unmarshaled.Method)
	}

	if unmarshaled.Params.Cursor != "page_token_123" {
		t.Errorf("Expected cursor 'page_token_123', got %s", unmarshaled.Params.Cursor)
	}
}

func TestListToolsResult(t *testing.T) {
	cursor := Cursor("next_page_token")
	result := ListToolsResult{
		Tools: []Tool{
			{
				BaseMetadata: BaseMetadata{
					Name:        "tool1",
					Description: "First tool",
				},
				InputSchema: ToolInputSchema{
					Type: "object",
				},
			},
			{
				BaseMetadata: BaseMetadata{
					Name:        "tool2",
					Description: "Second tool",
				},
				InputSchema: ToolInputSchema{
					Type: "object",
				},
			},
		},
		NextCursor: &cursor,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ListToolsResult: %v", err)
	}

	var unmarshaled ListToolsResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ListToolsResult: %v", err)
	}

	if len(unmarshaled.Tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(unmarshaled.Tools))
	}

	if unmarshaled.Tools[0].Name != "tool1" {
		t.Errorf("Expected first tool name 'tool1', got %s", unmarshaled.Tools[0].Name)
	}

	if unmarshaled.Tools[1].Name != "tool2" {
		t.Errorf("Expected second tool name 'tool2', got %s", unmarshaled.Tools[1].Name)
	}

	if unmarshaled.NextCursor == nil {
		t.Error("Expected NextCursor to be present")
	} else if *unmarshaled.NextCursor != "next_page_token" {
		t.Errorf("Expected NextCursor 'next_page_token', got %s", *unmarshaled.NextCursor)
	}
}

func TestToolMinimal(t *testing.T) {
	// Test tool with minimal required fields
	tool := Tool{
		BaseMetadata: BaseMetadata{
			Name: "minimal_tool",
		},
		InputSchema: ToolInputSchema{
			Type: "object",
		},
	}

	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("Failed to marshal minimal Tool: %v", err)
	}

	var unmarshaled Tool
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal minimal Tool: %v", err)
	}

	if unmarshaled.Name != "minimal_tool" {
		t.Errorf("Expected name 'minimal_tool', got %s", unmarshaled.Name)
	}

	if unmarshaled.InputSchema.Type != "object" {
		t.Errorf("Expected input schema type 'object', got %s", unmarshaled.InputSchema.Type)
	}

	// Optional fields should be nil/empty
	if unmarshaled.OutputSchema != nil {
		t.Error("Expected OutputSchema to be nil for minimal tool")
	}

	if unmarshaled.Annotations != nil {
		t.Error("Expected Annotations to be nil for minimal tool")
	}

	if unmarshaled.Description != "" {
		t.Errorf("Expected empty description, got %s", unmarshaled.Description)
	}
}

func TestToolComplexSchema(t *testing.T) {
	// Test tool with complex nested schema
	tool := Tool{
		BaseMetadata: BaseMetadata{
			Name: "complex_tool",
		},
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"config": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"nested_field": map[string]interface{}{
							"type":        "string",
							"description": "A nested field",
						},
						"array_field": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"item_name": map[string]interface{}{
										"type": "string",
									},
								},
							},
						},
					},
					"required": []string{"nested_field"},
				},
			},
			Required: []string{"config"},
		},
	}

	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("Failed to marshal complex Tool: %v", err)
	}

	var unmarshaled Tool
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal complex Tool: %v", err)
	}

	// Verify the complex schema structure is preserved
	configProp, exists := unmarshaled.InputSchema.Properties["config"]
	if !exists {
		t.Error("Expected 'config' property to exist")
	} else {
		configMap := configProp.(map[string]interface{})
		if configMap["type"] != "object" {
			t.Errorf("Expected config type 'object', got %v", configMap["type"])
		}

		// Check nested properties
		if props, ok := configMap["properties"].(map[string]interface{}); ok {
			if nestedField, exists := props["nested_field"]; exists {
				nestedMap := nestedField.(map[string]interface{})
				if nestedMap["type"] != "string" {
					t.Errorf("Expected nested_field type 'string', got %v", nestedMap["type"])
				}
			} else {
				t.Error("Expected nested_field to exist in config properties")
			}
		} else {
			t.Error("Expected config to have properties")
		}
	}
}

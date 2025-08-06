package types

import (
	"encoding/json"
	"testing"
)

func TestPrompt_JSONSerialization(t *testing.T) {
	prompt := Prompt{
		BaseMetadata: BaseMetadata{
			Name:        "test_prompt",
			Description: "A test prompt",
		},
		Description: "Detailed description of the prompt",
		Arguments: []PromptArgument{
			{
				BaseMetadata: BaseMetadata{
					Name:        "user_name",
					Description: "The user's name",
				},
				Description: "The name of the user to greet",
				Required:    true,
			},
			{
				BaseMetadata: BaseMetadata{
					Name:        "greeting_type",
					Description: "Type of greeting",
				},
				Description: "The type of greeting to use",
				Required:    false,
			},
		},
		Meta: Meta{
			"category": "greeting",
			"version":  "1.0",
		},
	}

	// Test marshaling
	data, err := json.Marshal(prompt)
	if err != nil {
		t.Fatalf("Failed to marshal Prompt: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Prompt
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Prompt: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != "test_prompt" {
		t.Errorf("Expected name 'test_prompt', got %s", unmarshaled.Name)
	}

	if unmarshaled.Description != "Detailed description of the prompt" {
		t.Errorf("Expected description 'Detailed description of the prompt', got %s", unmarshaled.Description)
	}

	if len(unmarshaled.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(unmarshaled.Arguments))
	}

	// Check first argument
	if unmarshaled.Arguments[0].Name != "user_name" {
		t.Errorf("Expected first argument name 'user_name', got %s", unmarshaled.Arguments[0].Name)
	}

	if !unmarshaled.Arguments[0].Required {
		t.Error("Expected first argument to be required")
	}

	// Check second argument
	if unmarshaled.Arguments[1].Name != "greeting_type" {
		t.Errorf("Expected second argument name 'greeting_type', got %s", unmarshaled.Arguments[1].Name)
	}

	if unmarshaled.Arguments[1].Required {
		t.Error("Expected second argument to not be required")
	}

	// Check meta
	if unmarshaled.Meta == nil {
		t.Error("Expected meta to be present")
	} else {
		if unmarshaled.Meta["category"] != "greeting" {
			t.Errorf("Expected category 'greeting', got %v", unmarshaled.Meta["category"])
		}
		if unmarshaled.Meta["version"] != "1.0" {
			t.Errorf("Expected version '1.0', got %v", unmarshaled.Meta["version"])
		}
	}
}

func TestPromptArgument_JSONSerialization(t *testing.T) {
	arg := PromptArgument{
		BaseMetadata: BaseMetadata{
			Name:        "input_text",
			Description: "Text to process",
		},
		Description: "The input text that will be processed",
		Required:    true,
	}

	data, err := json.Marshal(arg)
	if err != nil {
		t.Fatalf("Failed to marshal PromptArgument: %v", err)
	}

	var unmarshaled PromptArgument
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal PromptArgument: %v", err)
	}

	if unmarshaled.Name != "input_text" {
		t.Errorf("Expected name 'input_text', got %s", unmarshaled.Name)
	}

	if unmarshaled.Description != "The input text that will be processed" {
		t.Errorf("Expected description 'The input text that will be processed', got %s", unmarshaled.Description)
	}

	if !unmarshaled.Required {
		t.Error("Expected argument to be required")
	}
}

func TestPromptMessage_JSONSerialization(t *testing.T) {
	message := PromptMessage{
		Role: RoleUser,
		Content: TextContent{
			Type: "text",
			Text: "Hello, how can I help you?",
		},
	}

	data, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal PromptMessage: %v", err)
	}

	// Test that marshaling works and verify structure with map
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if temp["role"] != string(RoleUser) {
		t.Errorf("Expected role %s, got %v", RoleUser, temp["role"])
	}

	if content, ok := temp["content"].(map[string]interface{}); ok {
		if content["type"] != "text" {
			t.Errorf("Expected content type 'text', got %v", content["type"])
		}
		if content["text"] != "Hello, how can I help you?" {
			t.Errorf("Expected content text 'Hello, how can I help you?', got %v", content["text"])
		}
	} else {
		t.Error("Expected content to be an object")
	}
}

func TestListPromptsRequest(t *testing.T) {
	request := ListPromptsRequest{
		Method: "prompts/list",
		Params: struct {
			Cursor Cursor `json:"cursor,omitempty"`
		}{
			Cursor: "page_token_456",
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal ListPromptsRequest: %v", err)
	}

	var unmarshaled ListPromptsRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ListPromptsRequest: %v", err)
	}

	if unmarshaled.Method != "prompts/list" {
		t.Errorf("Expected method 'prompts/list', got %s", unmarshaled.Method)
	}

	if unmarshaled.Params.Cursor != "page_token_456" {
		t.Errorf("Expected cursor 'page_token_456', got %s", unmarshaled.Params.Cursor)
	}
}

func TestListPromptsResult(t *testing.T) {
	cursor := Cursor("next_page_token_prompts")
	result := ListPromptsResult{
		Prompts: []Prompt{
			{
				BaseMetadata: BaseMetadata{
					Name:        "prompt1",
					Description: "First prompt",
				},
				Arguments: []PromptArgument{
					{
						BaseMetadata: BaseMetadata{
							Name: "arg1",
						},
						Required: true,
					},
				},
			},
			{
				BaseMetadata: BaseMetadata{
					Name:        "prompt2",
					Description: "Second prompt",
				},
				Arguments: []PromptArgument{},
			},
		},
		NextCursor: &cursor,
		Meta: Meta{
			"total": 25,
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ListPromptsResult: %v", err)
	}

	var unmarshaled ListPromptsResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ListPromptsResult: %v", err)
	}

	if len(unmarshaled.Prompts) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(unmarshaled.Prompts))
	}

	if unmarshaled.Prompts[0].Name != "prompt1" {
		t.Errorf("Expected first prompt name 'prompt1', got %s", unmarshaled.Prompts[0].Name)
	}

	if unmarshaled.Prompts[1].Name != "prompt2" {
		t.Errorf("Expected second prompt name 'prompt2', got %s", unmarshaled.Prompts[1].Name)
	}

	if unmarshaled.NextCursor == nil {
		t.Error("Expected NextCursor to be present")
	} else if *unmarshaled.NextCursor != "next_page_token_prompts" {
		t.Errorf("Expected NextCursor 'next_page_token_prompts', got %s", *unmarshaled.NextCursor)
	}

	if unmarshaled.Meta == nil {
		t.Error("Expected meta to be present")
	} else if unmarshaled.Meta["total"] != float64(25) {
		t.Errorf("Expected total 25, got %v", unmarshaled.Meta["total"])
	}
}

func TestGetPromptRequest(t *testing.T) {
	request := GetPromptRequest{
		Method: "prompts/get",
		Params: struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments,omitempty"`
			Meta      Meta              `json:"_meta,omitempty"`
		}{
			Name: "greeting_prompt",
			Arguments: map[string]string{
				"user_name":     "Alice",
				"greeting_type": "formal",
			},
			Meta: Meta{
				"source": "client",
			},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal GetPromptRequest: %v", err)
	}

	var unmarshaled GetPromptRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal GetPromptRequest: %v", err)
	}

	if unmarshaled.Method != "prompts/get" {
		t.Errorf("Expected method 'prompts/get', got %s", unmarshaled.Method)
	}

	if unmarshaled.Params.Name != "greeting_prompt" {
		t.Errorf("Expected name 'greeting_prompt', got %s", unmarshaled.Params.Name)
	}

	if len(unmarshaled.Params.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(unmarshaled.Params.Arguments))
	}

	if unmarshaled.Params.Arguments["user_name"] != "Alice" {
		t.Errorf("Expected user_name 'Alice', got %s", unmarshaled.Params.Arguments["user_name"])
	}

	if unmarshaled.Params.Arguments["greeting_type"] != "formal" {
		t.Errorf("Expected greeting_type 'formal', got %s", unmarshaled.Params.Arguments["greeting_type"])
	}

	if unmarshaled.Params.Meta == nil {
		t.Error("Expected meta to be present")
	} else if unmarshaled.Params.Meta["source"] != "client" {
		t.Errorf("Expected source 'client', got %v", unmarshaled.Params.Meta["source"])
	}
}

func TestPromptMinimal(t *testing.T) {
	// Test prompt with minimal required fields
	prompt := Prompt{
		BaseMetadata: BaseMetadata{
			Name: "minimal_prompt",
		},
	}

	data, err := json.Marshal(prompt)
	if err != nil {
		t.Fatalf("Failed to marshal minimal Prompt: %v", err)
	}

	var unmarshaled Prompt
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal minimal Prompt: %v", err)
	}

	if unmarshaled.Name != "minimal_prompt" {
		t.Errorf("Expected name 'minimal_prompt', got %s", unmarshaled.Name)
	}

	// Optional fields should be empty/nil
	if unmarshaled.Description != "" {
		t.Errorf("Expected empty description, got %s", unmarshaled.Description)
	}

	if len(unmarshaled.Arguments) != 0 {
		t.Errorf("Expected empty arguments, got %d", len(unmarshaled.Arguments))
	}

	if unmarshaled.Meta != nil {
		t.Errorf("Expected nil meta, got %v", unmarshaled.Meta)
	}
}

func TestPromptArgumentMinimal(t *testing.T) {
	// Test argument with minimal required fields
	arg := PromptArgument{
		BaseMetadata: BaseMetadata{
			Name: "minimal_arg",
		},
	}

	data, err := json.Marshal(arg)
	if err != nil {
		t.Fatalf("Failed to marshal minimal PromptArgument: %v", err)
	}

	var unmarshaled PromptArgument
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal minimal PromptArgument: %v", err)
	}

	if unmarshaled.Name != "minimal_arg" {
		t.Errorf("Expected name 'minimal_arg', got %s", unmarshaled.Name)
	}

	// Optional fields should be empty/false
	if unmarshaled.Description != "" {
		t.Errorf("Expected empty description, got %s", unmarshaled.Description)
	}

	if unmarshaled.Required {
		t.Error("Expected Required to be false by default")
	}
}

package types

import (
	"encoding/json"
	"testing"
)

func TestSamplingMessage_JSONSerialization(t *testing.T) {
	message := SamplingMessage{
		Role: RoleUser,
		Content: TextContent{
			Type: "text",
			Text: "What is the weather like today?",
		},
	}

	data, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal SamplingMessage: %v", err)
	}

	// Test that marshaling works - unmarshaling ContentBlock interfaces is complex
	// so we just verify the JSON structure contains expected fields
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
		if content["text"] != "What is the weather like today?" {
			t.Errorf("Expected content text 'What is the weather like today?', got %v", content["text"])
		}
	} else {
		t.Error("Expected content to be an object")
	}
}

func TestCreateMessageRequest_JSONSerialization(t *testing.T) {
	temperature := 0.7
	costPriority := 0.8
	speedPriority := 0.5
	intelligencePriority := 0.9

	request := CreateMessageRequest{
		Method: "sampling/createMessage",
		Params: struct {
			Messages         []SamplingMessage      `json:"messages"`
			ModelPreferences *ModelPreferences      `json:"modelPreferences,omitempty"`
			SystemPrompt     string                 `json:"systemPrompt,omitempty"`
			IncludeContext   string                 `json:"includeContext,omitempty"`
			Temperature      *float64               `json:"temperature,omitempty"`
			MaxTokens        int                    `json:"maxTokens"`
			StopSequences    []string               `json:"stopSequences,omitempty"`
			Metadata         map[string]interface{} `json:"metadata,omitempty"`
		}{
			Messages: []SamplingMessage{
				{
					Role: RoleUser,
					Content: TextContent{
						Type: "text",
						Text: "Hello, how are you?",
					},
				},
				{
					Role: RoleAssistant,
					Content: TextContent{
						Type: "text",
						Text: "I'm doing well, thank you!",
					},
				},
			},
			ModelPreferences: &ModelPreferences{
				Hints: []ModelHint{
					{Name: "gpt-4"},
					{Name: "claude-3"},
				},
				CostPriority:         &costPriority,
				SpeedPriority:        &speedPriority,
				IntelligencePriority: &intelligencePriority,
			},
			SystemPrompt:   "You are a helpful assistant.",
			IncludeContext: "thisServer",
			Temperature:    &temperature,
			MaxTokens:      1000,
			StopSequences:  []string{"\n\n", "END"},
			Metadata: map[string]interface{}{
				"requestId": "req_123",
				"priority":  "high",
			},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal CreateMessageRequest: %v", err)
	}

	// Test that marshaling works and verify structure with map
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if temp["method"] != "sampling/createMessage" {
		t.Errorf("Expected method 'sampling/createMessage', got %v", temp["method"])
	}

	params, ok := temp["params"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected params to be an object")
	}

	messages, ok := params["messages"].([]interface{})
	if !ok {
		t.Fatal("Expected messages to be an array")
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	if params["systemPrompt"] != "You are a helpful assistant." {
		t.Errorf("Expected SystemPrompt 'You are a helpful assistant.', got %v", params["systemPrompt"])
	}

	if params["includeContext"] != "thisServer" {
		t.Errorf("Expected IncludeContext 'thisServer', got %v", params["includeContext"])
	}

	if params["temperature"] != 0.7 {
		t.Errorf("Expected Temperature 0.7, got %v", params["temperature"])
	}

	if params["maxTokens"] != float64(1000) {
		t.Errorf("Expected MaxTokens 1000, got %v", params["maxTokens"])
	}
}

func TestCreateMessageResult_JSONSerialization(t *testing.T) {
	result := CreateMessageResult{
		SamplingMessage: SamplingMessage{
			Role: RoleAssistant,
			Content: TextContent{
				Type: "text",
				Text: "Here is my response to your question.",
			},
		},
		Model:      "gpt-4-turbo",
		StopReason: "endTurn",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal CreateMessageResult: %v", err)
	}

	// Test that marshaling works and verify structure with map
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if temp["role"] != string(RoleAssistant) {
		t.Errorf("Expected role %s, got %v", RoleAssistant, temp["role"])
	}

	if temp["model"] != "gpt-4-turbo" {
		t.Errorf("Expected model 'gpt-4-turbo', got %v", temp["model"])
	}

	if temp["stopReason"] != "endTurn" {
		t.Errorf("Expected StopReason 'endTurn', got %v", temp["stopReason"])
	}
}

func TestModelPreferences_JSONSerialization(t *testing.T) {
	costPriority := 0.3
	speedPriority := 0.7
	intelligencePriority := 0.9

	prefs := ModelPreferences{
		Hints: []ModelHint{
			{Name: "claude-3-sonnet"},
			{Name: "gpt-4"},
			{Name: "gemini-pro"},
		},
		CostPriority:         &costPriority,
		SpeedPriority:        &speedPriority,
		IntelligencePriority: &intelligencePriority,
	}

	data, err := json.Marshal(prefs)
	if err != nil {
		t.Fatalf("Failed to marshal ModelPreferences: %v", err)
	}

	var unmarshaled ModelPreferences
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ModelPreferences: %v", err)
	}

	if len(unmarshaled.Hints) != 3 {
		t.Errorf("Expected 3 hints, got %d", len(unmarshaled.Hints))
	}

	if unmarshaled.Hints[0].Name != "claude-3-sonnet" {
		t.Errorf("Expected first hint 'claude-3-sonnet', got %s", unmarshaled.Hints[0].Name)
	}

	if unmarshaled.CostPriority == nil {
		t.Error("Expected CostPriority to be present")
	} else if *unmarshaled.CostPriority != 0.3 {
		t.Errorf("Expected CostPriority 0.3, got %f", *unmarshaled.CostPriority)
	}

	if unmarshaled.SpeedPriority == nil {
		t.Error("Expected SpeedPriority to be present")
	} else if *unmarshaled.SpeedPriority != 0.7 {
		t.Errorf("Expected SpeedPriority 0.7, got %f", *unmarshaled.SpeedPriority)
	}

	if unmarshaled.IntelligencePriority == nil {
		t.Error("Expected IntelligencePriority to be present")
	} else if *unmarshaled.IntelligencePriority != 0.9 {
		t.Errorf("Expected IntelligencePriority 0.9, got %f", *unmarshaled.IntelligencePriority)
	}
}

func TestModelHint_JSONSerialization(t *testing.T) {
	hint := ModelHint{
		Name: "gpt-4-vision",
	}

	data, err := json.Marshal(hint)
	if err != nil {
		t.Fatalf("Failed to marshal ModelHint: %v", err)
	}

	var unmarshaled ModelHint
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ModelHint: %v", err)
	}

	if unmarshaled.Name != "gpt-4-vision" {
		t.Errorf("Expected name 'gpt-4-vision', got %s", unmarshaled.Name)
	}
}

func TestCreateMessageRequestMinimal(t *testing.T) {
	// Test with minimal required fields
	request := CreateMessageRequest{
		Method: "sampling/createMessage",
		Params: struct {
			Messages         []SamplingMessage      `json:"messages"`
			ModelPreferences *ModelPreferences      `json:"modelPreferences,omitempty"`
			SystemPrompt     string                 `json:"systemPrompt,omitempty"`
			IncludeContext   string                 `json:"includeContext,omitempty"`
			Temperature      *float64               `json:"temperature,omitempty"`
			MaxTokens        int                    `json:"maxTokens"`
			StopSequences    []string               `json:"stopSequences,omitempty"`
			Metadata         map[string]interface{} `json:"metadata,omitempty"`
		}{
			Messages: []SamplingMessage{
				{
					Role: RoleUser,
					Content: TextContent{
						Type: "text",
						Text: "Simple question",
					},
				},
			},
			MaxTokens: 100,
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal minimal CreateMessageRequest: %v", err)
	}

	// Test that marshaling works and verify structure with map
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	params, ok := temp["params"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected params to be an object")
	}

	messages, ok := params["messages"].([]interface{})
	if !ok {
		t.Fatal("Expected messages to be an array")
	}

	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if params["maxTokens"] != float64(100) {
		t.Errorf("Expected MaxTokens 100, got %v", params["maxTokens"])
	}

	// Optional fields should not be present in minimal JSON
	if _, hasModelPrefs := params["modelPreferences"]; hasModelPrefs {
		t.Error("Expected ModelPreferences to be omitted")
	}

	if params["systemPrompt"] != "" && params["systemPrompt"] != nil {
		t.Errorf("Expected empty SystemPrompt, got %v", params["systemPrompt"])
	}
}

func TestModelPreferencesMinimal(t *testing.T) {
	// Test with minimal fields
	prefs := ModelPreferences{
		Hints: []ModelHint{
			{Name: "any-model"},
		},
	}

	data, err := json.Marshal(prefs)
	if err != nil {
		t.Fatalf("Failed to marshal minimal ModelPreferences: %v", err)
	}

	var unmarshaled ModelPreferences
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal minimal ModelPreferences: %v", err)
	}

	if len(unmarshaled.Hints) != 1 {
		t.Errorf("Expected 1 hint, got %d", len(unmarshaled.Hints))
	}

	// Optional priority fields should be nil
	if unmarshaled.CostPriority != nil {
		t.Errorf("Expected CostPriority to be nil, got %v", unmarshaled.CostPriority)
	}

	if unmarshaled.SpeedPriority != nil {
		t.Errorf("Expected SpeedPriority to be nil, got %v", unmarshaled.SpeedPriority)
	}

	if unmarshaled.IntelligencePriority != nil {
		t.Errorf("Expected IntelligencePriority to be nil, got %v", unmarshaled.IntelligencePriority)
	}
}

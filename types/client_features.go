// Package types contains MCP protocol client feature definitions
package types

// Sampling and client feature types

// SamplingMessage describes a message issued to or received from an LLM API
type SamplingMessage struct {
	Role    Role         `json:"role"`
	Content ContentBlock `json:"content"`
}

// CreateMessageRequest is a request from the server to sample an LLM via the client
type CreateMessageRequest struct {
	Method string `json:"method"`
	Params struct {
		Messages         []SamplingMessage      `json:"messages"`
		ModelPreferences *ModelPreferences      `json:"modelPreferences,omitempty"`
		SystemPrompt     string                 `json:"systemPrompt,omitempty"`
		IncludeContext   string                 `json:"includeContext,omitempty"` // "none" | "thisServer" | "allServers"
		Temperature      *float64               `json:"temperature,omitempty"`
		MaxTokens        int                    `json:"maxTokens"`
		StopSequences    []string               `json:"stopSequences,omitempty"`
		Metadata         map[string]interface{} `json:"metadata,omitempty"`
	} `json:"params"`
}

// CreateMessageResult is the client's response to a sampling/create_message request
type CreateMessageResult struct {
	SamplingMessage
	Model      string `json:"model"`
	StopReason string `json:"stopReason,omitempty"` // "endTurn" | "stopSequence" | "maxTokens" | string
}

// ModelPreferences represents the server's preferences for model selection
type ModelPreferences struct {
	Hints                []ModelHint `json:"hints,omitempty"`
	CostPriority         *float64    `json:"costPriority,omitempty"`
	SpeedPriority        *float64    `json:"speedPriority,omitempty"`
	IntelligencePriority *float64    `json:"intelligencePriority,omitempty"`
}

// ModelHint provides hints for model selection
type ModelHint struct {
	Name string `json:"name,omitempty"`
}

// Roots-related types

// Root represents a root directory or file that the server can operate on
type Root struct {
	URI  string `json:"uri"`
	Name string `json:"name,omitempty"`
	Meta Meta   `json:"_meta,omitempty"`
}

// ListRootsRequest is sent from the server to request a list of root URIs from the client
type ListRootsRequest struct {
	Method string `json:"method"`
}

// ListRootsResult is the client's response to a roots/list request
type ListRootsResult struct {
	Roots []Root `json:"roots"`
}

// RootsListChangedNotification informs that the list of roots has changed
type RootsListChangedNotification struct {
	Method string `json:"method"`
}

// Elicitation types

// ElicitRequest is a request from the server to elicit additional information from the user
type ElicitRequest struct {
	Method string `json:"method"`
	Params struct {
		Message         string              `json:"message"`
		RequestedSchema ElicitRequestSchema `json:"requestedSchema"`
	} `json:"params"`
}

// ElicitRequestSchema represents a restricted subset of JSON Schema
type ElicitRequestSchema struct {
	Type       string                               `json:"type"`
	Properties map[string]PrimitiveSchemaDefinition `json:"properties"`
	Required   []string                             `json:"required,omitempty"`
}

// PrimitiveSchemaDefinition represents restricted schema definitions
type PrimitiveSchemaDefinition interface{}

// StringSchema represents a string schema definition
type StringSchema struct {
	Type        string `json:"type"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	MinLength   *int   `json:"minLength,omitempty"`
	MaxLength   *int   `json:"maxLength,omitempty"`
	Format      string `json:"format,omitempty"` // "email" | "uri" | "date" | "date-time"
}

// NumberSchema represents a number schema definition
type NumberSchema struct {
	Type        string   `json:"type"` // "number" | "integer"
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Minimum     *float64 `json:"minimum,omitempty"`
	Maximum     *float64 `json:"maximum,omitempty"`
}

// BooleanSchema represents a boolean schema definition
type BooleanSchema struct {
	Type        string `json:"type"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     *bool  `json:"default,omitempty"`
}

// EnumSchema represents an enum schema definition
type EnumSchema struct {
	Type        string   `json:"type"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum"`
	EnumNames   []string `json:"enumNames,omitempty"`
}

// ElicitResult is the client's response to an elicitation request
type ElicitResult struct {
	Action  string                 `json:"action"` // "accept" | "decline" | "cancel"
	Content map[string]interface{} `json:"content,omitempty"`
}

// Logging types

// SetLevelRequest is a request from the client to the server to enable or adjust logging
type SetLevelRequest struct {
	Method string `json:"method"`
	Params struct {
		Level LoggingLevel `json:"level"`
	} `json:"params"`
}

// LoggingMessageNotification is a notification of a log message from server to client
type LoggingMessageNotification struct {
	Method string `json:"method"`
	Params struct {
		Level  LoggingLevel `json:"level"`
		Logger string       `json:"logger,omitempty"`
		Data   interface{}  `json:"data"`
	} `json:"params"`
}

// Completion types

// CompleteRequest is a request from the client to the server for completion options
type CompleteRequest struct {
	Method string `json:"method"`
	Params struct {
		Ref      CompletionReference `json:"ref"`
		Argument struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"argument"`
		Context *CompletionContext `json:"context,omitempty"`
	} `json:"params"`
}

// CompletionReference represents a reference to a resource or prompt
type CompletionReference interface{}

// ResourceTemplateReference represents a reference to a resource template
type ResourceTemplateReference struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

// PromptReference represents a reference to a prompt
type PromptReference struct {
	BaseMetadata
	Type string `json:"type"`
}

// CompletionContext provides additional context for completions
type CompletionContext struct {
	Arguments map[string]string `json:"arguments,omitempty"`
}

// CompleteResult is the server's response to a completion/complete request
type CompleteResult struct {
	Completion struct {
		Values  []string `json:"values"`
		Total   *int     `json:"total,omitempty"`
		HasMore *bool    `json:"hasMore,omitempty"`
	} `json:"completion"`
}

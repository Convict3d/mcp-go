// Package types contains MCP protocol prompt definitions
package types

// Prompt represents a prompt or prompt template that the server offers
type Prompt struct {
	BaseMetadata
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
	Meta        Meta             `json:"_meta,omitempty"`
}

// PromptArgument describes an argument that a prompt can accept
type PromptArgument struct {
	BaseMetadata
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptMessage describes a message returned as part of a prompt
type PromptMessage struct {
	Role    Role         `json:"role"`
	Content ContentBlock `json:"content"`
}

// Prompt request/response types

// ListPromptsRequest is sent from the client to request a list of prompts
type ListPromptsRequest struct {
	Method string `json:"method"`
	Params struct {
		Cursor Cursor `json:"cursor,omitempty"`
	} `json:"params,omitempty"`
}

// ListPromptsResult is the server's response to a prompts/list request
type ListPromptsResult struct {
	Prompts    []Prompt `json:"prompts"`
	NextCursor *Cursor  `json:"nextCursor,omitempty"`
	Meta       Meta     `json:"_meta,omitempty"`
}

// GetPromptRequest is used by the client to get a prompt provided by the server
type GetPromptRequest struct {
	Method string `json:"method"`
	Params struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments,omitempty"`
		Meta      Meta              `json:"_meta,omitempty"`
	} `json:"params"`
}

// GetPromptResult is the server's response to a prompts/get request
type GetPromptResult struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
	Meta        Meta            `json:"_meta,omitempty"`
}

// PromptsListChangedNotification informs that the list of prompts has changed
type PromptsListChangedNotification struct {
	Method string `json:"method"`
	Params struct {
		Meta Meta `json:"_meta,omitempty"`
	} `json:"params,omitempty"`
}

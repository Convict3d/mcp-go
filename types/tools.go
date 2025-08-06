// Package types contains MCP protocol tool definitions
package types

// Tool represents a tool the client can call
type Tool struct {
	BaseMetadata
	Description  string            `json:"description,omitempty"`
	InputSchema  ToolInputSchema   `json:"inputSchema"`
	OutputSchema *ToolOutputSchema `json:"outputSchema,omitempty"`
	Annotations  *ToolAnnotations  `json:"annotations,omitempty"`
	Meta         Meta              `json:"_meta,omitempty"`
}

// ToolInputSchema defines the expected parameters for the tool
type ToolInputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// ToolOutputSchema defines the structure of the tool's output
type ToolOutputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// ToolAnnotations provide additional properties describing a Tool to clients
type ToolAnnotations struct {
	Title           string `json:"title,omitempty"`
	ReadOnlyHint    bool   `json:"readOnlyHint,omitempty"`
	DestructiveHint bool   `json:"destructiveHint,omitempty"`
	IdempotentHint  bool   `json:"idempotentHint,omitempty"`
	OpenWorldHint   bool   `json:"openWorldHint,omitempty"`
}

// Tool request/response types

// ListToolsRequest is sent from the client to request a list of tools
type ListToolsRequest struct {
	Method string `json:"method"`
	Params struct {
		Cursor Cursor `json:"cursor,omitempty"`
	} `json:"params,omitempty"`
}

// ListToolsResult is the server's response to a tools/list request
type ListToolsResult struct {
	Tools      []Tool  `json:"tools"`
	NextCursor *Cursor `json:"nextCursor,omitempty"`
	Meta       Meta    `json:"_meta,omitempty"`
}

// CallToolRequest is sent from the client to invoke a tool
type CallToolRequest struct {
	Method string `json:"method"`
	Params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments,omitempty"`
		Meta      Meta                   `json:"_meta,omitempty"`
	} `json:"params"`
}

// CallToolResult is the server's response to a tools/call request
type CallToolResult struct {
	Content []interface{} `json:"content"` // Using interface{} for flexible JSON unmarshaling
	IsError bool          `json:"isError,omitempty"`
	Meta    Meta          `json:"_meta,omitempty"`
}

// GetTextContent extracts text content from the result as properly typed TextContent structs
func (ctr *CallToolResult) GetTextContent() []TextContent {
	var texts []TextContent
	for _, content := range ctr.Content {
		if contentMap, ok := content.(map[string]interface{}); ok {
			if contentType, ok := contentMap["type"].(string); ok && contentType == "text" {
				if text, ok := contentMap["text"].(string); ok {
					textContent := TextContent{
						Type: contentType,
						Text: text,
					}
					// Parse annotations if present
					if annotations, ok := contentMap["annotations"].(map[string]interface{}); ok {
						textContent.Annotations = parseAnnotations(annotations)
					}
					texts = append(texts, textContent)
				}
			}
		}
	}
	return texts
}

// GetTextStrings is a convenience method that returns just the text strings
func (ctr *CallToolResult) GetTextStrings() []string {
	var texts []string
	textContents := ctr.GetTextContent()
	for _, tc := range textContents {
		texts = append(texts, tc.Text)
	}
	return texts
}

// GetImageContent extracts image content from the result
func (ctr *CallToolResult) GetImageContent() []ImageContent {
	var images []ImageContent
	for _, content := range ctr.Content {
		if contentMap, ok := content.(map[string]interface{}); ok {
			if contentType, ok := contentMap["type"].(string); ok && contentType == "image" {
				image := ImageContent{Type: contentType}
				if data, ok := contentMap["data"].(string); ok {
					image.Data = data
				}
				if mimeType, ok := contentMap["mimeType"].(string); ok {
					image.MimeType = mimeType
				}
				// Parse annotations if present
				if annotations, ok := contentMap["annotations"].(map[string]interface{}); ok {
					image.Annotations = parseAnnotations(annotations)
				}
				images = append(images, image)
			}
		}
	}
	return images
}

// GetAudioContent extracts audio content from the result
func (ctr *CallToolResult) GetAudioContent() []AudioContent {
	var audios []AudioContent
	for _, content := range ctr.Content {
		if contentMap, ok := content.(map[string]interface{}); ok {
			if contentType, ok := contentMap["type"].(string); ok && contentType == "audio" {
				audio := AudioContent{Type: contentType}
				if data, ok := contentMap["data"].(string); ok {
					audio.Data = data
				}
				if mimeType, ok := contentMap["mimeType"].(string); ok {
					audio.MimeType = mimeType
				}
				// Parse annotations if present
				if annotations, ok := contentMap["annotations"].(map[string]interface{}); ok {
					audio.Annotations = parseAnnotations(annotations)
				}
				audios = append(audios, audio)
			}
		}
	}
	return audios
}

// GetResourceLinkContent extracts resource link content from the result
func (ctr *CallToolResult) GetResourceLinkContent() []ResourceLinkContent {
	var resourceLinks []ResourceLinkContent
	for _, content := range ctr.Content {
		if contentMap, ok := content.(map[string]interface{}); ok {
			if contentType, ok := contentMap["type"].(string); ok && contentType == "resource_link" {
				resourceLink := ResourceLinkContent{Type: contentType}
				if uri, ok := contentMap["uri"].(string); ok {
					resourceLink.URI = uri
				}
				if name, ok := contentMap["name"].(string); ok {
					resourceLink.Name = name
				}
				if description, ok := contentMap["description"].(string); ok {
					resourceLink.Description = description
				}
				if mimeType, ok := contentMap["mimeType"].(string); ok {
					resourceLink.MimeType = mimeType
				}
				// Parse annotations if present
				if annotations, ok := contentMap["annotations"].(map[string]interface{}); ok {
					resourceLink.Annotations = parseAnnotations(annotations)
				}
				resourceLinks = append(resourceLinks, resourceLink)
			}
		}
	}
	return resourceLinks
}

// GetResourceContent extracts embedded resource content from the result
func (ctr *CallToolResult) GetResourceContent() []ResourceContent {
	var resources []ResourceContent
	for _, content := range ctr.Content {
		if contentMap, ok := content.(map[string]interface{}); ok {
			if contentType, ok := contentMap["type"].(string); ok && contentType == "resource" {
				resourceContent := ResourceContent{Type: contentType}
				if resourceData, ok := contentMap["resource"].(map[string]interface{}); ok {
					resource := parseResource(resourceData)
					resourceContent.Resource = &resource
				}
				// Parse annotations if present
				if annotations, ok := contentMap["annotations"].(map[string]interface{}); ok {
					resourceContent.Annotations = parseAnnotations(annotations)
				}
				resources = append(resources, resourceContent)
			}
		}
	}
	return resources
}

// GetAllContent returns all content items as their proper types
func (ctr *CallToolResult) GetAllContent() []ContentBlock {
	var contents []ContentBlock

	// Add text content
	for _, tc := range ctr.GetTextContent() {
		contents = append(contents, tc)
	}

	// Add image content
	for _, ic := range ctr.GetImageContent() {
		contents = append(contents, ic)
	}

	// Add audio content
	for _, ac := range ctr.GetAudioContent() {
		contents = append(contents, ac)
	}

	// Add resource link content
	for _, rlc := range ctr.GetResourceLinkContent() {
		contents = append(contents, rlc)
	}

	// Add resource content
	for _, rc := range ctr.GetResourceContent() {
		contents = append(contents, rc)
	}

	return contents
}

// GetContentType returns the type of the first content item
func (ctr *CallToolResult) GetContentType() string {
	if len(ctr.Content) == 0 {
		return ""
	}
	if contentMap, ok := ctr.Content[0].(map[string]interface{}); ok {
		if contentType, ok := contentMap["type"].(string); ok {
			return contentType
		}
	}
	return ""
}

// parseAnnotations converts a map to Annotations struct
func parseAnnotations(annotationMap map[string]interface{}) *Annotations {
	annotations := &Annotations{}

	if audience, ok := annotationMap["audience"].([]interface{}); ok {
		for _, role := range audience {
			if roleStr, ok := role.(string); ok {
				annotations.Audience = append(annotations.Audience, Role(roleStr))
			}
		}
	}

	if priority, ok := annotationMap["priority"].(float64); ok {
		priorityInt := int(priority)
		annotations.Priority = &priorityInt
	}

	return annotations
}

// parseResource converts a map to Resource struct
func parseResource(resourceMap map[string]interface{}) Resource {
	resource := Resource{}

	if uri, ok := resourceMap["uri"].(string); ok {
		resource.URI = uri
	}
	if name, ok := resourceMap["name"].(string); ok {
		resource.BaseMetadata.Name = name
	}
	if description, ok := resourceMap["description"].(string); ok {
		resource.Description = description
	}
	if mimeType, ok := resourceMap["mimeType"].(string); ok {
		resource.MimeType = mimeType
	}
	if size, ok := resourceMap["size"].(float64); ok {
		sizeInt := int(size)
		resource.Size = &sizeInt
	}

	// Parse annotations if present
	if annotations, ok := resourceMap["annotations"].(map[string]interface{}); ok {
		resource.Annotations = parseAnnotations(annotations)
	}

	return resource
}

// ToolsListChangedNotification informs that the list of tools has changed
type ToolsListChangedNotification struct {
	Method string `json:"method"`
	Params struct {
		Meta Meta `json:"_meta,omitempty"`
	} `json:"params,omitempty"`
}

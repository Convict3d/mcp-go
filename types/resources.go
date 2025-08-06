// Package types contains MCP protocol resource definitions
package types

// Resource represents a known resource that the server is capable of reading
type Resource struct {
	BaseMetadata
	URI         string       `json:"uri"`
	Description string       `json:"description,omitempty"`
	MimeType    string       `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
	Size        *int         `json:"size,omitempty"`
	Meta        Meta         `json:"_meta,omitempty"`
}

// ResourceTemplate represents a template description for resources available on the server
type ResourceTemplate struct {
	BaseMetadata
	URITemplate string       `json:"uriTemplate"`
	Description string       `json:"description,omitempty"`
	MimeType    string       `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
	Meta        Meta         `json:"_meta,omitempty"`
}

// ResourceContents represents the contents of a specific resource or sub-resource
type ResourceContents struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Meta     Meta   `json:"_meta,omitempty"`
}

// TextResourceContents represents text resource contents
type TextResourceContents struct {
	ResourceContents
	Text string `json:"text"`
}

// BlobResourceContents represents binary resource contents
type BlobResourceContents struct {
	ResourceContents
	Blob string `json:"blob"`
}

// ResourceLink represents a resource that can be included in prompts or tool results
type ResourceLink struct {
	Resource
	Type string `json:"type"`
}

// Annotations provide additional metadata
type Annotations struct {
	Audience []Role `json:"audience,omitempty"`
	Priority *int   `json:"priority,omitempty"`
}

// Resource request/response types

// ListResourcesRequest is sent from the client to request a list of resources
type ListResourcesRequest struct {
	Method string `json:"method"`
	Params struct {
		Cursor Cursor `json:"cursor,omitempty"`
	} `json:"params,omitempty"`
}

// ListResourcesResult is the server's response to a resources/list request
type ListResourcesResult struct {
	Resources  []Resource `json:"resources"`
	NextCursor *Cursor    `json:"nextCursor,omitempty"`
	Meta       Meta       `json:"_meta,omitempty"`
}

// ListResourceTemplatesRequest is sent from the client to request resource templates
type ListResourceTemplatesRequest struct {
	Method string `json:"method"`
	Params struct {
		Cursor Cursor `json:"cursor,omitempty"`
	} `json:"params,omitempty"`
}

// ListResourceTemplatesResult is the server's response to a resources/templates/list request
type ListResourceTemplatesResult struct {
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates"`
	NextCursor        *Cursor            `json:"nextCursor,omitempty"`
	Meta              Meta               `json:"_meta,omitempty"`
}

// ReadResourceRequest is sent from the client to read a specific resource
type ReadResourceRequest struct {
	Method string `json:"method"`
	Params struct {
		URI  string `json:"uri"`
		Meta Meta   `json:"_meta,omitempty"`
	} `json:"params"`
}

// ReadResourceResult is the server's response to a resources/read request
type ReadResourceResult struct {
	Contents []ResourceContents `json:"contents"`
	Meta     Meta               `json:"_meta,omitempty"`
}

// SubscribeToResourceRequest is sent from the client to subscribe to resource changes
type SubscribeToResourceRequest struct {
	Method string `json:"method"`
	Params struct {
		URI  string `json:"uri"`
		Meta Meta   `json:"_meta,omitempty"`
	} `json:"params"`
}

// UnsubscribeFromResourceRequest is sent from the client to unsubscribe from resource changes
type UnsubscribeFromResourceRequest struct {
	Method string `json:"method"`
	Params struct {
		URI  string `json:"uri"`
		Meta Meta   `json:"_meta,omitempty"`
	} `json:"params"`
}

// ResourceUpdatedNotification informs that a subscribed resource has changed
type ResourceUpdatedNotification struct {
	Method string `json:"method"`
	Params struct {
		URI  string `json:"uri"`
		Meta Meta   `json:"_meta,omitempty"`
	} `json:"params"`
}

// ResourcesListChangedNotification informs that the list of resources has changed
type ResourcesListChangedNotification struct {
	Method string `json:"method"`
	Params struct {
		Meta Meta `json:"_meta,omitempty"`
	} `json:"params,omitempty"`
}

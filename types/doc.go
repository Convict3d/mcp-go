/*
Package types contains all MCP protocol type definitions.

This package provides comprehensive type definitions for the Model Context Protocol (MCP)
specification version 2025-06-18. All types are designed for JSON serialization and
include proper validation and documentation.

# Core Protocol Types

The package defines all fundamental MCP types:

  - Request and Response types for JSON-RPC communication
  - Content types for text, image, and audio data
  - Tool definitions and schemas
  - Resource definitions and templates
  - Prompt definitions and arguments

# Content Types

MCP supports various content types for rich data exchange:

	// Text content
	text := &types.TextContent{
		Type: types.ContentTypeText,
		Text: "Hello, world!",
	}

	// Image content
	image := &types.ImageContent{
		Type:     types.ContentTypeImage,
		Data:     base64ImageData,
		MimeType: "image/png",
	}

# Protocol Constants

Important protocol constants are defined:

	const LatestProtocolVersion = "2025-06-18"

	const (
		ContentTypeText  = "text"
		ContentTypeImage = "image"
		ContentTypeAudio = "audio"
	)

# Tool Definitions

Tools are defined with comprehensive schemas:

	tool := types.Tool{
		BaseMetadata: types.BaseMetadata{
			Name:        "calculator",
			Description: "Basic calculator operations",
		},
		InputSchema: types.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"enum": []string{"add", "subtract", "multiply", "divide"},
				},
				"a": map[string]interface{}{"type": "number"},
				"b": map[string]interface{}{"type": "number"},
			},
			Required: []string{"operation", "a", "b"},
		},
	}

# Resource Management

Resources represent external data sources:

	resource := types.Resource{
		BaseMetadata: types.BaseMetadata{
			Name:        "config",
			Description: "Application configuration",
		},
		URI:      "file:///app/config.json",
		MimeType: "application/json",
	}

# JSON Serialization

All types support proper JSON marshaling and unmarshaling:

	data, err := json.Marshal(tool)
	if err != nil {
		return err
	}

	var unmarshaled types.Tool
	err = json.Unmarshal(data, &unmarshaled)

# Validation

Types include built-in validation where appropriate. Required fields are enforced
through Go's type system and JSON tags.

# Protocol Compliance

All types are designed to be fully compliant with the MCP specification and
include proper field names, types, and validation rules.
*/
package types

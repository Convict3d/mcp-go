package types

import (
	"encoding/json"
	"testing"
)

func TestResource_JSONSerialization(t *testing.T) {
	size := 1024
	priority := 1
	resource := Resource{
		BaseMetadata: BaseMetadata{
			Name:        "test_resource",
			Description: "A test resource",
		},
		URI:         "file:///path/to/resource.txt",
		Description: "Detailed description",
		MimeType:    "text/plain",
		Annotations: &Annotations{
			Audience: []Role{RoleUser, RoleAssistant},
			Priority: &priority,
		},
		Size: &size,
		Meta: Meta{
			"lastModified": "2023-01-01T00:00:00Z",
		},
	}

	// Test marshaling
	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal Resource: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Resource: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != "test_resource" {
		t.Errorf("Expected name 'test_resource', got %s", unmarshaled.Name)
	}

	if unmarshaled.URI != "file:///path/to/resource.txt" {
		t.Errorf("Expected URI 'file:///path/to/resource.txt', got %s", unmarshaled.URI)
	}

	if unmarshaled.Description != "Detailed description" {
		t.Errorf("Expected description 'Detailed description', got %s", unmarshaled.Description)
	}

	if unmarshaled.MimeType != "text/plain" {
		t.Errorf("Expected mime type 'text/plain', got %s", unmarshaled.MimeType)
	}

	if unmarshaled.Size == nil {
		t.Error("Expected size to be present")
	} else if *unmarshaled.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", *unmarshaled.Size)
	}

	if unmarshaled.Annotations == nil {
		t.Error("Expected annotations to be present")
	} else {
		if len(unmarshaled.Annotations.Audience) != 2 {
			t.Errorf("Expected 2 audience roles, got %d", len(unmarshaled.Annotations.Audience))
		}
		if unmarshaled.Annotations.Priority == nil {
			t.Error("Expected priority to be present")
		} else if *unmarshaled.Annotations.Priority != 1 {
			t.Errorf("Expected priority 1, got %d", *unmarshaled.Annotations.Priority)
		}
	}

	if unmarshaled.Meta == nil {
		t.Error("Expected meta to be present")
	} else if unmarshaled.Meta["lastModified"] != "2023-01-01T00:00:00Z" {
		t.Errorf("Expected lastModified '2023-01-01T00:00:00Z', got %v", unmarshaled.Meta["lastModified"])
	}
}

func TestResourceTemplate_JSONSerialization(t *testing.T) {
	priority := 1
	template := ResourceTemplate{
		BaseMetadata: BaseMetadata{
			Name:        "file_template",
			Description: "Template for file resources",
		},
		URITemplate: "file:///documents/{filename}",
		Description: "Access files in the documents directory",
		MimeType:    "application/octet-stream",
		Annotations: &Annotations{
			Audience: []Role{RoleUser},
			Priority: &priority,
		},
		Meta: Meta{
			"pattern": "*.txt",
		},
	}

	data, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceTemplate: %v", err)
	}

	var unmarshaled ResourceTemplate
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ResourceTemplate: %v", err)
	}

	if unmarshaled.Name != "file_template" {
		t.Errorf("Expected name 'file_template', got %s", unmarshaled.Name)
	}

	if unmarshaled.URITemplate != "file:///documents/{filename}" {
		t.Errorf("Expected URI template 'file:///documents/{filename}', got %s", unmarshaled.URITemplate)
	}

	if unmarshaled.Description != "Access files in the documents directory" {
		t.Errorf("Expected description 'Access files in the documents directory', got %s", unmarshaled.Description)
	}

	if unmarshaled.MimeType != "application/octet-stream" {
		t.Errorf("Expected mime type 'application/octet-stream', got %s", unmarshaled.MimeType)
	}
}

func TestTextResourceContents(t *testing.T) {
	contents := TextResourceContents{
		ResourceContents: ResourceContents{
			URI:      "file:///test.txt",
			MimeType: "text/plain",
			Meta: Meta{
				"encoding": "utf-8",
			},
		},
		Text: "Hello, world!\nThis is test content.",
	}

	data, err := json.Marshal(contents)
	if err != nil {
		t.Fatalf("Failed to marshal TextResourceContents: %v", err)
	}

	var unmarshaled TextResourceContents
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal TextResourceContents: %v", err)
	}

	if unmarshaled.URI != "file:///test.txt" {
		t.Errorf("Expected URI 'file:///test.txt', got %s", unmarshaled.URI)
	}

	if unmarshaled.MimeType != "text/plain" {
		t.Errorf("Expected mime type 'text/plain', got %s", unmarshaled.MimeType)
	}

	if unmarshaled.Text != "Hello, world!\nThis is test content." {
		t.Errorf("Expected text 'Hello, world!\nThis is test content.', got %s", unmarshaled.Text)
	}

	if unmarshaled.Meta == nil {
		t.Error("Expected meta to be present")
	} else if unmarshaled.Meta["encoding"] != "utf-8" {
		t.Errorf("Expected encoding 'utf-8', got %v", unmarshaled.Meta["encoding"])
	}
}

func TestBlobResourceContents(t *testing.T) {
	contents := BlobResourceContents{
		ResourceContents: ResourceContents{
			URI:      "file:///image.png",
			MimeType: "image/png",
			Meta: Meta{
				"size": 2048,
			},
		},
		Blob: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
	}

	data, err := json.Marshal(contents)
	if err != nil {
		t.Fatalf("Failed to marshal BlobResourceContents: %v", err)
	}

	var unmarshaled BlobResourceContents
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal BlobResourceContents: %v", err)
	}

	if unmarshaled.URI != "file:///image.png" {
		t.Errorf("Expected URI 'file:///image.png', got %s", unmarshaled.URI)
	}

	if unmarshaled.MimeType != "image/png" {
		t.Errorf("Expected mime type 'image/png', got %s", unmarshaled.MimeType)
	}

	expectedBlob := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="
	if unmarshaled.Blob != expectedBlob {
		t.Errorf("Expected blob data to match, got different content")
	}

	if unmarshaled.Meta == nil {
		t.Error("Expected meta to be present")
	}
}

func TestResourceLink(t *testing.T) {
	link := ResourceLink{
		Resource: Resource{
			BaseMetadata: BaseMetadata{
				Name:        "linked_resource",
				Description: "A linked resource",
			},
			URI:      "https://example.com/api/resource/123",
			MimeType: "application/json",
		},
		Type: "resource_link",
	}

	data, err := json.Marshal(link)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceLink: %v", err)
	}

	var unmarshaled ResourceLink
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ResourceLink: %v", err)
	}

	if unmarshaled.Name != "linked_resource" {
		t.Errorf("Expected name 'linked_resource', got %s", unmarshaled.Name)
	}

	if unmarshaled.URI != "https://example.com/api/resource/123" {
		t.Errorf("Expected URI 'https://example.com/api/resource/123', got %s", unmarshaled.URI)
	}

	if unmarshaled.Type != "resource_link" {
		t.Errorf("Expected type 'resource_link', got %s", unmarshaled.Type)
	}
}

func TestResourceMinimal(t *testing.T) {
	// Test resource with minimal required fields
	resource := Resource{
		BaseMetadata: BaseMetadata{
			Name: "minimal_resource",
		},
		URI: "file:///minimal.txt",
	}

	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal minimal Resource: %v", err)
	}

	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal minimal Resource: %v", err)
	}

	if unmarshaled.Name != "minimal_resource" {
		t.Errorf("Expected name 'minimal_resource', got %s", unmarshaled.Name)
	}

	if unmarshaled.URI != "file:///minimal.txt" {
		t.Errorf("Expected URI 'file:///minimal.txt', got %s", unmarshaled.URI)
	}

	// Optional fields should be empty/nil
	if unmarshaled.Description != "" {
		t.Errorf("Expected empty description, got %s", unmarshaled.Description)
	}

	if unmarshaled.MimeType != "" {
		t.Errorf("Expected empty mime type, got %s", unmarshaled.MimeType)
	}

	if unmarshaled.Size != nil {
		t.Errorf("Expected nil size, got %v", unmarshaled.Size)
	}

	if unmarshaled.Annotations != nil {
		t.Errorf("Expected nil annotations, got %v", unmarshaled.Annotations)
	}
}

func TestResourceTemplateMinimal(t *testing.T) {
	// Test resource template with minimal required fields
	template := ResourceTemplate{
		BaseMetadata: BaseMetadata{
			Name: "minimal_template",
		},
		URITemplate: "file:///{path}",
	}

	data, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal minimal ResourceTemplate: %v", err)
	}

	var unmarshaled ResourceTemplate
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal minimal ResourceTemplate: %v", err)
	}

	if unmarshaled.Name != "minimal_template" {
		t.Errorf("Expected name 'minimal_template', got %s", unmarshaled.Name)
	}

	if unmarshaled.URITemplate != "file:///{path}" {
		t.Errorf("Expected URI template 'file:///{path}', got %s", unmarshaled.URITemplate)
	}
}

func TestResourceContents(t *testing.T) {
	// Test base ResourceContents struct
	contents := ResourceContents{
		URI:      "https://api.example.com/data",
		MimeType: "application/json",
		Meta: Meta{
			"version": "v1.0",
			"cached":  true,
		},
	}

	data, err := json.Marshal(contents)
	if err != nil {
		t.Fatalf("Failed to marshal ResourceContents: %v", err)
	}

	var unmarshaled ResourceContents
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ResourceContents: %v", err)
	}

	if unmarshaled.URI != "https://api.example.com/data" {
		t.Errorf("Expected URI 'https://api.example.com/data', got %s", unmarshaled.URI)
	}

	if unmarshaled.MimeType != "application/json" {
		t.Errorf("Expected mime type 'application/json', got %s", unmarshaled.MimeType)
	}

	if unmarshaled.Meta == nil {
		t.Error("Expected meta to be present")
	} else {
		if unmarshaled.Meta["version"] != "v1.0" {
			t.Errorf("Expected version 'v1.0', got %v", unmarshaled.Meta["version"])
		}
		if unmarshaled.Meta["cached"] != true {
			t.Errorf("Expected cached true, got %v", unmarshaled.Meta["cached"])
		}
	}
}

func TestResourceWithComplexAnnotations(t *testing.T) {
	priority := 2
	resource := Resource{
		BaseMetadata: BaseMetadata{
			Name: "annotated_resource",
		},
		URI: "file:///annotated.txt",
		Annotations: &Annotations{
			Audience: []Role{RoleUser, RoleAssistant},
			Priority: &priority,
		},
	}

	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal Resource with annotations: %v", err)
	}

	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Resource with annotations: %v", err)
	}

	if unmarshaled.Annotations == nil {
		t.Error("Expected annotations to be present")
	} else {
		if len(unmarshaled.Annotations.Audience) != 2 {
			t.Errorf("Expected 2 audience roles, got %d", len(unmarshaled.Annotations.Audience))
		}

		if unmarshaled.Annotations.Audience[0] != RoleUser {
			t.Errorf("Expected first audience role to be %s, got %s", RoleUser, unmarshaled.Annotations.Audience[0])
		}

		if unmarshaled.Annotations.Audience[1] != RoleAssistant {
			t.Errorf("Expected second audience role to be %s, got %s", RoleAssistant, unmarshaled.Annotations.Audience[1])
		}

		if unmarshaled.Annotations.Priority == nil {
			t.Error("Expected priority to be present")
		} else if *unmarshaled.Annotations.Priority != 2 {
			t.Errorf("Expected priority 2, got %d", *unmarshaled.Annotations.Priority)
		}
	}
}

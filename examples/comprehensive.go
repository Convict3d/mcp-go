// Package main demonstrates comprehensive MCP client usage.
//
// This example shows how to use the MCP Go client library to connect to an MCP server
// and perform various operations including listing tools, resources, and prompts.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Convict3d/mcp-go/client"
	"github.com/Convict3d/mcp-go/types"
)

func main() {
	// Create a new MCP client with configuration options
	c := client.NewClient(
		client.WithClientInfo("mcp-example", "1.0.0"),
		client.WithTimeout(30*time.Second),
	)

	// Initialize the MCP connection
	fmt.Println("Initializing MCP client...")
	err := c.Initialize(types.LatestProtocolVersion)
	if err != nil {
		log.Fatalf("Failed to initialize MCP client: %v", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	fmt.Printf("Successfully connected to MCP server (protocol version: %s)\n\n",
		types.LatestProtocolVersion)

	// Demonstrate listing tools
	demonstrateTools(c)

	// Demonstrate listing resources
	demonstrateResources(c)

	// Demonstrate listing prompts
	demonstratePrompts(c)

	fmt.Println("=== MCP Client Demo Complete ===")
}

// demonstrateTools shows how to list and work with tools
func demonstrateTools(c *client.Client) {
	fmt.Println("=== Tools Demo ===")

	tools, err := c.ListTools()
	if err != nil {
		log.Printf("Failed to list tools: %v", err)
		return
	}

	if len(tools) == 0 {
		fmt.Println("No tools available from the server")
		return
	}

	fmt.Printf("Found %d tools:\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("%d. %s\n", i+1, tool.Name)
		if tool.Description != "" {
			fmt.Printf("   Description: %s\n", tool.Description)
		}

		// Show input schema type
		fmt.Printf("   Input Schema: %s\n", tool.InputSchema.Type)

		// Show required parameters if any
		if len(tool.InputSchema.Required) > 0 {
			fmt.Printf("   Required Parameters: %v\n", tool.InputSchema.Required)
		}

		// Show annotations if present
		if tool.Annotations != nil {
			fmt.Printf("   Title: %s\n", tool.Annotations.Title)
			if tool.Annotations.ReadOnlyHint {
				fmt.Printf("   (Read-only operation)\n")
			}
		}
		fmt.Println()
	}
}

// demonstrateResources shows how to list and work with resources
func demonstrateResources(c *client.Client) {
	fmt.Println("=== Resources Demo ===")

	resources, err := c.ListResources()
	if err != nil {
		log.Printf("Failed to list resources: %v", err)
		return
	}

	if len(resources) == 0 {
		fmt.Println("No resources available from the server")
		return
	}

	fmt.Printf("Found %d resources:\n", len(resources))
	for i, resource := range resources {
		fmt.Printf("%d. %s\n", i+1, resource.Name)
		if resource.Description != "" {
			fmt.Printf("   Description: %s\n", resource.Description)
		}
		fmt.Printf("   URI: %s\n", resource.URI)
		if resource.MimeType != "" {
			fmt.Printf("   MIME Type: %s\n", resource.MimeType)
		}
		fmt.Println()
	}
}

// demonstratePrompts shows how to list and work with prompts
func demonstratePrompts(c *client.Client) {
	fmt.Println("=== Prompts Demo ===")

	prompts, err := c.ListPrompts()
	if err != nil {
		log.Printf("Failed to list prompts: %v", err)
		return
	}

	if len(prompts) == 0 {
		fmt.Println("No prompts available from the server")
		return
	}

	fmt.Printf("Found %d prompts:\n", len(prompts))
	for i, prompt := range prompts {
		fmt.Printf("%d. %s\n", i+1, prompt.Name)
		if prompt.Description != "" {
			fmt.Printf("   Description: %s\n", prompt.Description)
		}

		// Show prompt arguments if any
		if len(prompt.Arguments) > 0 {
			fmt.Printf("   Arguments:\n")
			for _, arg := range prompt.Arguments {
				fmt.Printf("     - %s", arg.Name)
				if arg.Description != "" {
					fmt.Printf(": %s", arg.Description)
				}
				if arg.Required {
					fmt.Printf(" (required)")
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
}

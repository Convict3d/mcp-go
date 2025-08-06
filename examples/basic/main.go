package main

import (
	"fmt"
	"log"
	"time"

	"github.com/convict3d/mcp-go/client"
	"github.com/convict3d/mcp-go/transport/http"
	"github.com/convict3d/mcp-go/types"
)

func main() {
	fmt.Println("🚀 MCP Go Client Library - Professional Example")
	fmt.Println("==============================================")

	transport := http.NewHTTPTransport("http://localhost:9831/mcp",
		http.WithTimeout(30*time.Second),
		http.WithHeader("Accept", "application/json, text/event-stream"),
	)

	// Create the MCP client using option pattern
	c := client.NewClient(
		client.WithTransport(transport),
		client.WithClientInfo("professional-example", "1.0.0"),
	)

	// Initialize connection
	fmt.Println("\n📡 Connecting to MCP server...")
	err := c.Initialize(types.LatestProtocolVersion)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Get server info
	serverInfo := c.GetServerInfo()
	if serverInfo != nil {
		fmt.Printf("✅ Connected to: %s v%s\n", serverInfo.Name, serverInfo.Version)
	}

	// Get server capabilities
	fmt.Printf("🔧 Server capabilities:\n")
	fmt.Printf("   Tools: %v\n", c.HasTools())
	fmt.Printf("   Resources: %v\n", c.HasResources())
	fmt.Printf("   Prompts: %v\n", c.HasPrompts())
	fmt.Printf("   Session ID: %s\n", c.GetSessionID())

	// List and call tools if available
	if c.HasTools() {
		fmt.Println("\n🛠️  Listing tools...")
		tools, err := c.ListTools()
		if err != nil {
			log.Printf("Failed to list tools: %v", err)
		} else {
			fmt.Printf("Found %d tools:\n", len(tools))
			for i, tool := range tools {
				fmt.Printf("   %d. %s - %s (%+v)\n", i+1, tool.Name, tool.Description, tool.OutputSchema)
			}

			// Try calling a tool if available
			if len(tools) > 0 {
				fmt.Printf("\n🎯 Calling tool: %s\n", "browser_navigate")

				// Example tool call - adjust arguments based on your server
				result, err := c.CallTool("browser_navigate", map[string]interface{}{
					"url": "https://facebook.com",
				})

				if err != nil {
					log.Printf("Tool call failed: %v", err)
				} else {
					fmt.Printf("✅ Tool executed successfully!\n")
					if result.IsError {
						fmt.Printf("⚠️  Tool returned an error\n")
					}
					if len(result.Content) > 0 {
						contentType := result.GetContentType()
						fmt.Printf("📄 Content type: %s\n", contentType)

						switch contentType {
						case "text":
							textContents := result.GetTextContent()
							if len(textContents) > 0 {
								fmt.Printf("📄 Tool output: %s\n", textContents[0].Text)
							}
						case "image":
							imageContents := result.GetImageContent()
							if len(imageContents) > 0 {
								fmt.Printf("📄 Image output: %s (%s)\n", imageContents[0].MimeType, "data available")
							}
						default:
							fmt.Printf("📄 Tool output available (%d items)\n", len(result.Content))
						}
					}
				}

				result, erro := c.CallTool("browser_take_screenshot", map[string]interface{}{
					"format": "png",
				})
				if erro != nil {
					log.Printf("Screenshot tool call failed: %v", erro)
				} else {
					fmt.Printf("✅ Screenshot taken successfully!\n")
					if result.IsError {
						fmt.Printf("⚠️  Screenshot tool returned an error\n")
					}
					if len(result.Content) > 0 {
						allContent := result.GetAllContent()

						for _, content := range allContent {
							fmt.Printf("📄 Content type: %s\n", content.ContentType())
							switch content.ContentType() {
							case "text":
								textContents := result.GetTextContent()
								if len(textContents) > 0 {
									fmt.Printf("📄 Tool output: %s\n", textContents[0].Text)
								}
							case "image":
								imageContents := result.GetImageContent()
								if len(imageContents) > 0 {
									for _, img := range imageContents {
										fmt.Printf("📸 Screenshot captured (image data available, type: %s)\n", img.MimeType)
									}
								}
							default:
								fmt.Printf("📄 Tool output available (%d items)\n", len(result.Content))
							}
						}
					}
				}

				result, err = c.CallTool("browser_snapshot", map[string]interface{}{})
				if err != nil {
					log.Printf("Snapshot tool call failed: %v", err)
				} else {
					fmt.Printf("✅ Snapshot taken successfully!\n")
					if result.IsError {
						fmt.Printf("⚠️  Snapshot tool returned an error\n")
					}
					if len(result.Content) > 0 {
						allContent := result.GetAllContent()

						for _, content := range allContent {
							fmt.Printf("📄 Content type: %s\n", content.ContentType())
							switch content.ContentType() {
							case "text":
								textContents := result.GetTextContent()
								if len(textContents) > 0 {
									fmt.Printf("📄 Tool output: %s\n", textContents[0].Text)
								}
							case "image":
								imageContents := result.GetImageContent()
								if len(imageContents) > 0 {
									for _, img := range imageContents {
										fmt.Printf("📸 Snapshot captured (image data available, type: %s)\n", img.MimeType)
									}
								}
							default:
								fmt.Printf("📄 Tool output available (%d items)\n", len(result.Content))
							}
						}
					}
				}

				result, err = c.CallTool("browser_click", map[string]interface{}{
					"element": "link",
					"ref":     "e107",
				})
				if err != nil {
					log.Printf("Click tool call failed: %v", err)
				} else {
					fmt.Printf("✅ Click action performed successfully!\n")
					if result.IsError {
						fmt.Printf("⚠️  Click tool returned an error\n")
					}
					if len(result.Content) > 0 {
						allContent := result.GetAllContent()

						for _, content := range allContent {
							fmt.Printf("📄 Content type: %s\n", content.ContentType())
							switch content.ContentType() {
							case "text":
								textContents := result.GetTextContent()
								if len(textContents) > 0 {
									fmt.Printf("📄 Tool output: %s\n", textContents[0].Text)
								}
							case "image":
								imageContents := result.GetImageContent()
								if len(imageContents) > 0 {
									for _, img := range imageContents {
										fmt.Printf("📸 Click action captured (image data available, type: %s)\n", img.MimeType)
									}
								}
							default:
								fmt.Printf("📄 Tool output available (%d items)\n", len(result.Content))
							}
						}
					}
				}
			}
		}
	}

	// List resources if available
	if c.HasResources() {
		fmt.Println("\n📁 Listing resources...")
		resources, err := c.ListResources()
		if err != nil {
			log.Printf("Failed to list resources: %v", err)
		} else {
			fmt.Printf("Found %d resources:\n", len(resources))
			for i, resource := range resources {
				fmt.Printf("   %d. %s (%s)\n", i+1, resource.Name, resource.URI)
			}
		}
	}

	// List prompts if available
	if c.HasPrompts() {
		fmt.Println("\n💬 Listing prompts...")
		prompts, err := c.ListPrompts()
		if err != nil {
			log.Printf("Failed to list prompts: %v", err)
		} else {
			fmt.Printf("Found %d prompts:\n", len(prompts))
			for i, prompt := range prompts {
				fmt.Printf("   %d. %s - %s\n", i+1, prompt.Name, prompt.Description)
			}
		}
	}

	// Clean up
	err = c.Close()
	if err != nil {
		log.Printf("Failed to close client: %v", err)
	}

	fmt.Println("\n🏁 Example completed successfully!")
	fmt.Println("Your professional MCP Go client library is working! 🎉")
}

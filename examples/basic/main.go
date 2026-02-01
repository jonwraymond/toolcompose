// Example: Basic toolset composition
//
// Demonstrates creating a toolset from tools and exporting to a protocol format.
package main

import (
	"fmt"
	"log"

	"github.com/jonwraymond/toolcompose/set"
	adapt "github.com/jonwraymond/toolfoundation/adapter"
)

func main() {
	// Create sample tools
	tools := []*adapt.CanonicalTool{
		{
			Name:        "search",
			Namespace:   "github",
			Description: "Search GitHub repositories",
			InputSchema: &adapt.JSONSchema{
				Type: "object",
				Properties: map[string]*adapt.JSONSchema{
					"query": {Type: "string", Description: "Search query"},
				},
				Required: []string{"query"},
			},
		},
		{
			Name:        "create-issue",
			Namespace:   "github",
			Description: "Create a GitHub issue",
			InputSchema: &adapt.JSONSchema{
				Type: "object",
				Properties: map[string]*adapt.JSONSchema{
					"title": {Type: "string", Description: "Issue title"},
					"body":  {Type: "string", Description: "Issue body"},
				},
				Required: []string{"title"},
			},
		},
		{
			Name:        "list-files",
			Namespace:   "filesystem",
			Description: "List files in a directory",
			InputSchema: &adapt.JSONSchema{
				Type: "object",
				Properties: map[string]*adapt.JSONSchema{
					"path": {Type: "string", Description: "Directory path"},
				},
			},
		},
	}

	// Build toolset from tools
	ts, err := set.NewBuilder("my-tools").
		FromTools(tools).
		Build()
	if err != nil {
		log.Fatalf("Failed to build toolset: %v", err)
	}

	fmt.Printf("Toolset '%s' contains %d tools:\n", ts.Name(), ts.Count())
	for _, id := range ts.IDs() {
		tool, _ := ts.Get(id)
		fmt.Printf("  - %s: %s\n", id, tool.Description)
	}

	// Export to protocol format using adapter
	fmt.Println("\nExporting to MCP format...")
	adapter := adapt.NewMCPAdapter()
	exposure := set.NewExposure(ts, adapter)

	exported, err := exposure.Export()
	if err != nil {
		log.Fatalf("Failed to export: %v", err)
	}

	fmt.Printf("Exported %d tools to MCP format\n", len(exported))
}

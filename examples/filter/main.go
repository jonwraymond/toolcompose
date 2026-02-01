// Example: Toolset filtering
//
// Demonstrates using filters and policies to create focused toolsets.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jonwraymond/toolcompose/set"
	adapt "github.com/jonwraymond/toolfoundation/adapter"
)

func main() {
	// Create a diverse set of tools
	tools := []*adapt.CanonicalTool{
		{Name: "search", Namespace: "github", Tags: []string{"read", "public"}, Description: "Search repos"},
		{Name: "create-issue", Namespace: "github", Tags: []string{"write", "issues"}, Description: "Create issue"},
		{Name: "delete-repo", Namespace: "github", Tags: []string{"write", "admin", "dangerous"}, Description: "Delete repo"},
		{Name: "list-files", Namespace: "filesystem", Tags: []string{"read"}, Description: "List files"},
		{Name: "write-file", Namespace: "filesystem", Tags: []string{"write"}, Description: "Write file"},
		{Name: "delete-file", Namespace: "filesystem", Tags: []string{"write", "dangerous"}, Description: "Delete file"},
	}

	fmt.Println("=== Filter by Namespace ===")
	githubOnly, err := set.NewBuilder("github-tools").
		FromTools(tools).
		WithNamespace("github").
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printToolset(githubOnly)

	fmt.Println("\n=== Filter by Tags (read-only) ===")
	readOnly, err := set.NewBuilder("read-only").
		FromTools(tools).
		WithFilter(set.TagsAll("read")).
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printToolset(readOnly)

	fmt.Println("\n=== Custom Filter (no dangerous) ===")
	safe, err := set.NewBuilder("safe-tools").
		FromTools(tools).
		WithFilter(func(t *adapt.CanonicalTool) bool {
			for _, tag := range t.Tags {
				if tag == "dangerous" {
					return false
				}
			}
			return true
		}).
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printToolset(safe)

	fmt.Println("\n=== Policy: Allow specific IDs only ===")
	restricted, err := set.NewBuilder("restricted").
		FromTools(tools).
		WithFilter(set.AllowIDs("github:search", "github:create-issue", "filesystem:list-files")).
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printToolset(restricted)

	fmt.Println("\n=== Policy: Deny tags ===")
	noDangerous, err := set.NewBuilder("no-dangerous").
		FromTools(tools).
		WithPolicy(set.DenyTags("dangerous")).
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printToolset(noDangerous)

	fmt.Println("\n=== Combined: Namespace + Tags + Custom Filter ===")
	combined, err := set.NewBuilder("github-write").
		FromTools(tools).
		WithNamespace("github").
		WithFilter(set.TagsAll("write")).
		WithFilter(func(t *adapt.CanonicalTool) bool {
			// Exclude tools with "admin" tag
			for _, tag := range t.Tags {
				if tag == "admin" {
					return false
				}
			}
			return true
		}).
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printToolset(combined)
}

func printToolset(ts *set.Toolset) {
	fmt.Printf("Toolset '%s' (%d tools):\n", ts.Name(), ts.Count())
	for _, id := range ts.IDs() {
		tool, _ := ts.Get(id)
		tags := strings.Join(tool.Tags, ", ")
		if tags == "" {
			tags = "(no tags)"
		}
		fmt.Printf("  - %s [%s]\n", id, tags)
	}
}

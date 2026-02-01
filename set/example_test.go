package set_test

import (
	"fmt"

	"github.com/jonwraymond/toolcompose/set"
	"github.com/jonwraymond/toolfoundation/adapter"
)

func ExampleNewBuilder() {
	tools := []*adapter.CanonicalTool{
		{Name: "create_issue", Namespace: "github", Tags: []string{"write", "issues"}},
		{Name: "list_repos", Namespace: "github", Tags: []string{"read", "repos"}},
		{Name: "get_user", Namespace: "github", Tags: []string{"read", "users"}},
	}

	ts, err := set.NewBuilder("readonly").
		FromTools(tools).
		WithTags([]string{"read"}).
		Build()
	if err != nil {
		panic(err)
	}

	fmt.Println("Count:", ts.Count())
	fmt.Println("IDs:", ts.IDs())
	// Output:
	// Count: 2
	// IDs: [github:get_user github:list_repos]
}

func ExampleNewBuilder_withNamespace() {
	tools := []*adapter.CanonicalTool{
		{Name: "create_issue", Namespace: "github"},
		{Name: "list_files", Namespace: "filesystem"},
		{Name: "list_repos", Namespace: "github"},
	}

	ts, _ := set.NewBuilder("github-only").
		FromTools(tools).
		WithNamespace("github").
		Build()

	fmt.Println("IDs:", ts.IDs())
	// Output:
	// IDs: [github:create_issue github:list_repos]
}

func ExampleNewBuilder_withPolicy() {
	tools := []*adapter.CanonicalTool{
		{Name: "create_issue", Namespace: "github", Tags: []string{"write"}},
		{Name: "delete_repo", Namespace: "github", Tags: []string{"write", "dangerous"}},
		{Name: "list_repos", Namespace: "github", Tags: []string{"read"}},
	}

	ts, _ := set.NewBuilder("safe-tools").
		FromTools(tools).
		WithPolicy(set.DenyTags("dangerous")).
		Build()

	fmt.Println("IDs:", ts.IDs())
	// Output:
	// IDs: [github:create_issue github:list_repos]
}

func ExampleToolset_Filter() {
	ts := set.New("all-tools")
	ts.Add(&adapter.CanonicalTool{Name: "tool1", Namespace: "ns1", Tags: []string{"fast"}})
	ts.Add(&adapter.CanonicalTool{Name: "tool2", Namespace: "ns1", Tags: []string{"slow"}})
	ts.Add(&adapter.CanonicalTool{Name: "tool3", Namespace: "ns2", Tags: []string{"fast"}})

	filtered := ts.Filter(set.TagsAny("fast"))

	fmt.Println("Original count:", ts.Count())
	fmt.Println("Filtered count:", filtered.Count())
	fmt.Println("Filtered IDs:", filtered.IDs())
	// Output:
	// Original count: 3
	// Filtered count: 2
	// Filtered IDs: [ns1:tool1 ns2:tool3]
}

func ExampleAllowScopes() {
	tools := []*adapter.CanonicalTool{
		{Name: "read_file", Namespace: "fs", RequiredScopes: []string{"read"}},
		{Name: "write_file", Namespace: "fs", RequiredScopes: []string{"read", "write"}},
		{Name: "delete_file", Namespace: "fs", RequiredScopes: []string{"read", "write", "delete"}},
	}

	// User only has read scope
	ts, _ := set.NewBuilder("user-tools").
		FromTools(tools).
		WithPolicy(set.AllowScopes("read")).
		Build()

	fmt.Println("Tools with read scope:", ts.IDs())
	// Output:
	// Tools with read scope: [fs:read_file]
}

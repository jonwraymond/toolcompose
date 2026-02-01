// Package set provides composable tool collection building.
//
// Toolset enables curated, filtered, and access-controlled tool surfaces
// from multiple sources. It is pure data composition with no I/O, execution,
// or network dependencies.
//
// # Core Concepts
//
//   - Toolset: Thread-safe collection of canonical tools
//   - Builder: Fluent API for constructing toolsets from sources
//   - FilterFunc: Predicates for filtering tools (AND-composed)
//   - Policy: Access control decisions (applied after filters)
//   - Exposure: Export to MCP/OpenAI/Anthropic via adapter
//
// # Basic Usage
//
//	ts, err := set.NewBuilder("github-readonly").
//	    FromRegistry(myRegistry).
//	    WithNamespace("github").
//	    WithTags([]string{"read"}).
//	    WithPolicy(set.DenyTags("deprecated")).
//	    Build()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Export to OpenAI format
//	exp := set.NewExposure(ts, adapter.NewOpenAIAdapter())
//	tools, warnings, errs := exp.ExportWithWarnings()
//
// # Filter Composition
//
// Filters are AND-composed: a tool must pass ALL filters to be included.
// Use custom FilterFunc for OR logic:
//
//	orFilter := func(t *adapter.CanonicalTool) bool {
//	    return set.NamespaceFilter("ns1")(t) || set.NamespaceFilter("ns2")(t)
//	}
//	ts, _ := set.NewBuilder("custom").
//	    FromTools(tools).
//	    WithFilter(orFilter).
//	    Build()
//
// # Built-in Filters
//
//   - NamespaceFilter(ns...): Match tools in any listed namespace (OR)
//   - TagsAny(tags...): Match tools with ANY of the tags (OR)
//   - TagsAll(tags...): Match tools with ALL tags (AND)
//   - TagsNone(tags...): Match tools with NONE of the tags
//   - CategoryFilter(categories...): Match tools in any category (OR)
//   - AllowIDs(ids...): Include only listed tool IDs
//   - DenyIDs(ids...): Exclude listed tool IDs
//
// # Policies
//
// Policies provide access control after filtering:
//
//   - AllowAll(): Allow all non-nil tools
//   - DenyAll(): Deny all tools
//   - AllowNamespaces(ns...): Allow only listed namespaces
//   - DenyTags(tags...): Deny tools with any forbidden tag
//   - AllowScopes(scopes...): Allow tools requiring only allowed scopes
//
// # Thread Safety
//
//   - Toolset: Safe for concurrent use (RWMutex protected)
//   - Builder: NOT safe for concurrent use (single-goroutine construction)
//   - Policy: Safe after construction (immutable closures)
//   - FilterFunc: Must be safe if Toolset is shared across goroutines
//
// # Integration
//
// The set package integrates with:
//
//   - [github.com/jonwraymond/toolfoundation/adapter] for protocol conversion
//   - [github.com/jonwraymond/tooldiscovery/index] as a Registry source
//
// # Error Handling
//
// The package defines sentinel errors for programmatic handling:
//
//   - ErrNoSource: Builder.Build() called without FromTools or FromRegistry
//   - ErrNilAdapter: Exposure created with nil adapter
//   - ErrNilToolset: Operation requires non-nil Toolset
//
// Use errors.Is() for reliable error checking:
//
//	if errors.Is(err, set.ErrNoSource) {
//	    // handle missing source
//	}
package set

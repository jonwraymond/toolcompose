package set

import "github.com/jonwraymond/toolfoundation/adapter"

// Registry provides tools for the builder.
//
// Contract:
// - Concurrency: Tools may be called concurrently; implementations must be safe or document otherwise.
// - Ownership: returned slice is caller-owned; tools are shared and read-only.
// - Determinism: ordering should be deterministic for identical registry state.
// - Nil handling: returning nil is treated as empty.
type Registry interface {
	Tools() []*adapter.CanonicalTool
}

// Builder constructs Toolsets with filtering.
//
// Thread-safety: Builder is NOT safe for concurrent use. Create and
// configure Builders in a single goroutine, then call Build() once.
// The resulting Toolset IS thread-safe.
//
// Execution flow:
//  1. Source tools gathered from FromTools or FromRegistry
//  2. Filters applied in order (AND-composed)
//  3. Policy applied last (if set)
//  4. Resulting tools added to new Toolset
type Builder struct {
	name      string
	source    []*adapter.CanonicalTool
	sourceSet bool // tracks whether FromTools was called (even with nil)
	registry  Registry
	filters   []FilterFunc
	policy    Policy
}

// NewBuilder creates a new Builder with the given toolset name.
func NewBuilder(name string) *Builder {
	return &Builder{name: name}
}

// FromTools sets tools as the source.
func (b *Builder) FromTools(tools []*adapter.CanonicalTool) *Builder {
	b.source = tools
	b.sourceSet = true
	return b
}

// FromRegistry sets a registry as the source.
func (b *Builder) FromRegistry(r Registry) *Builder {
	b.registry = r
	return b
}

// WithNamespace filters to a single namespace.
func (b *Builder) WithNamespace(ns string) *Builder {
	b.filters = append(b.filters, NamespaceFilter(ns))
	return b
}

// WithNamespaces filters to multiple namespaces.
func (b *Builder) WithNamespaces(ns []string) *Builder {
	b.filters = append(b.filters, NamespaceFilter(ns...))
	return b
}

// WithTags filters to tools with ALL specified tags.
func (b *Builder) WithTags(tags []string) *Builder {
	b.filters = append(b.filters, TagsAll(tags...))
	return b
}

// WithCategories filters to tools with ANY category.
func (b *Builder) WithCategories(categories []string) *Builder {
	b.filters = append(b.filters, CategoryFilter(categories...))
	return b
}

// WithTools includes only listed tool IDs.
func (b *Builder) WithTools(ids []string) *Builder {
	b.filters = append(b.filters, AllowIDs(ids...))
	return b
}

// ExcludeTools excludes listed tool IDs.
func (b *Builder) ExcludeTools(ids []string) *Builder {
	b.filters = append(b.filters, DenyIDs(ids...))
	return b
}

// WithFilter adds a custom filter.
func (b *Builder) WithFilter(fn FilterFunc) *Builder {
	b.filters = append(b.filters, fn)
	return b
}

// WithPolicy sets the access control policy (applied after filters).
func (b *Builder) WithPolicy(p Policy) *Builder {
	b.policy = p
	return b
}

// Build creates the Toolset.
func (b *Builder) Build() (*Toolset, error) {
	// Gather source tools
	var tools []*adapter.CanonicalTool
	if b.registry != nil {
		tools = b.registry.Tools()
	} else if b.source != nil || b.sourceSet {
		tools = b.source
	} else {
		return nil, ErrNoSource
	}

	// Apply filters (AND composition)
	for _, filter := range b.filters {
		var filtered []*adapter.CanonicalTool
		for _, t := range tools {
			if filter(t) {
				filtered = append(filtered, t)
			}
		}
		tools = filtered
	}

	// Apply policy (last)
	if b.policy != nil {
		var allowed []*adapter.CanonicalTool
		for _, t := range tools {
			if b.policy.Allow(t) {
				allowed = append(allowed, t)
			}
		}
		tools = allowed
	}

	// Build toolset
	ts := New(b.name)
	for _, t := range tools {
		ts.Add(t)
	}
	return ts, nil
}

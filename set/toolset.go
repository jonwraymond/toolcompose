package set

import (
	"sort"
	"sync"

	"github.com/jonwraymond/toolfoundation/adapter"
)

// FilterFunc is a predicate for filtering tools.
//
// Contract:
//   - Must not call methods on the Toolset being filtered (would deadlock)
//   - Must not mutate the tool (treat as read-only)
//   - Must be safe for concurrent calls if Toolset is shared
//   - Must return false for nil tools
type FilterFunc func(*adapter.CanonicalTool) bool

// Toolset is a thread-safe collection of canonical tools.
//
// All methods are safe for concurrent use. Tools are stored by reference;
// callers should not mutate tools after adding them to a Toolset.
// Results from IDs() and Tools() are deterministically sorted by ID.
type Toolset struct {
	name  string
	mu    sync.RWMutex
	tools map[string]*adapter.CanonicalTool // keyed by ID()
}

// New creates a new Toolset with the given name.
func New(name string) *Toolset {
	return &Toolset{
		name:  name,
		tools: make(map[string]*adapter.CanonicalTool),
	}
}

// Name returns the toolset's name.
func (ts *Toolset) Name() string { return ts.name }

// Add adds a tool. Nil tools are silently ignored.
func (ts *Toolset) Add(tool *adapter.CanonicalTool) {
	if tool == nil {
		return
	}
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tools[tool.ID()] = tool
}

// Get retrieves a tool by ID. Returns (nil, false) if not found.
func (ts *Toolset) Get(id string) (*adapter.CanonicalTool, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tool, ok := ts.tools[id]
	return tool, ok
}

// Remove removes a tool by ID. Returns true if found and removed.
func (ts *Toolset) Remove(id string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if _, ok := ts.tools[id]; ok {
		delete(ts.tools, id)
		return true
	}
	return false
}

// Count returns the number of tools.
func (ts *Toolset) Count() int {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return len(ts.tools)
}

// IDs returns tool IDs sorted lexicographically.
func (ts *Toolset) IDs() []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	ids := make([]string, 0, len(ts.tools))
	for id := range ts.tools {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Tools returns all tools sorted lexicographically by ID.
func (ts *Toolset) Tools() []*adapter.CanonicalTool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tools := make([]*adapter.CanonicalTool, 0, len(ts.tools))
	for _, t := range ts.tools {
		tools = append(tools, t)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].ID() < tools[j].ID()
	})
	return tools
}

// Filter returns a new Toolset with tools matching fn.
// The original Toolset is not modified.
func (ts *Toolset) Filter(fn FilterFunc) *Toolset {
	ts.mu.RLock()
	// Snapshot matching tools while holding lock
	var matches []*adapter.CanonicalTool
	for _, t := range ts.tools {
		if fn(t) {
			matches = append(matches, t)
		}
	}
	ts.mu.RUnlock()

	// Build new toolset from snapshot (no lock needed)
	filtered := New(ts.name + "-filtered")
	for _, t := range matches {
		filtered.tools[t.ID()] = t
	}
	return filtered
}

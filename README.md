# toolcompose

[![CI](https://github.com/jonwraymond/toolcompose/actions/workflows/ci.yml/badge.svg)](https://github.com/jonwraymond/toolcompose/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jonwraymond/toolcompose.svg)](https://pkg.go.dev/github.com/jonwraymond/toolcompose)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonwraymond/toolcompose)](https://goreportcard.com/report/github.com/jonwraymond/toolcompose)

Toolset composition, filtering, and skill orchestration for the
[ApertureStack](https://github.com/jonwraymond) AI tool ecosystem.

## Installation

```bash
go get github.com/jonwraymond/toolcompose@latest
```

## Packages

| Package | Description |
|---------|-------------|
| [`set`](https://pkg.go.dev/github.com/jonwraymond/toolcompose/set) | Toolset composition, filtering, and protocol exposure |
| [`skill`](https://pkg.go.dev/github.com/jonwraymond/toolcompose/skill) | Declarative skill planning and execution |

## Quick Start

### Toolset Composition

```go
import (
    "github.com/jonwraymond/toolcompose/set"
    "github.com/jonwraymond/toolfoundation/adapter"
)

// Create toolset from registry with filters and policy
ts, err := set.NewBuilder("my-toolset").
    FromRegistry(registry).
    WithNamespace("github").
    WithFilter(set.TagsAll("read")).
    WithPolicy(set.DenyTags("dangerous")).
    Build()

// Export to protocol format
exposure := set.NewExposure(ts, adapter.NewMCPAdapter())
tools, err := exposure.Export()
```

### Skill Execution

```go
import "github.com/jonwraymond/toolcompose/skill"

// Define a skill
s := skill.Skill{
    Name: "create-issue",
    Steps: []skill.Step{
        {ID: "search", ToolID: "github:search", Inputs: map[string]any{"q": "bug"}},
        {ID: "create", ToolID: "github:create-issue", Inputs: map[string]any{"title": "Fix bug"}},
    },
}

// Plan and execute
planner := skill.NewPlanner()
plan, err := planner.Plan(s)
results, err := skill.Execute(ctx, plan, myRunner)
```

## Features

- **Toolset Composition**: Build filtered toolsets from registries or tool slices
- **Filter Chaining**: Compose filters with AND semantics
- **Built-in Filters**: Namespace, tags, scopes, and custom predicates
- **Policy System**: Allow/deny lists with DenyByDefault or AllowByDefault
- **Protocol Exposure**: Export to MCP, OpenAI, or Anthropic formats
- **Feature Warnings**: Detect schema feature loss during conversion
- **Skill Planning**: Validate and plan multi-step tool sequences
- **Guard System**: Max steps, allowed tool IDs, and custom guards
- **Deterministic Execution**: Reproducible results for testing
- **Thread-Safe**: Toolsets are safe for concurrent reads

## Thread Safety

| Type | Thread-Safe | Notes |
|------|-------------|-------|
| `Toolset` | Yes | Safe for concurrent reads |
| `Builder` | No | Use in single goroutine |
| `Exposure` | Yes | Safe after construction |
| `Planner` | Yes | Stateless, safe to share |
| `Plan` | Yes | Immutable after creation |

## Documentation

- [Package documentation](https://pkg.go.dev/github.com/jonwraymond/toolcompose)
- [Design notes](./docs/design-notes.md)
- [Examples](./examples/)

## Examples

```bash
go run ./examples/basic
go run ./examples/filter
go run ./examples/skill
```

## Related Packages

toolcompose is part of the ApertureStack ecosystem:

```
toolfoundation     Core types, adapters, model definitions
     │
     v
tooldiscovery      Search, indexing, semantic discovery
     │
     v
toolexec           Execution, runtime isolation, backends
     │
     v
toolcompose        Composition, filtering, skill orchestration (this package)
```

## License

MIT License - see [LICENSE](./LICENSE)

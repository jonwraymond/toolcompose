# toolcompose

Composition layer for building filtered tool collections and declarative skills.

## Packages

| Package | Purpose |
|---------|---------|
| `set` | Toolset composition, filtering, and exposure |
| `skill` | Declarative skill planning and execution |

## Installation

```bash
go get github.com/jonwraymond/toolcompose@latest
```

## Quick Start: Toolset

```go
import (
  "fmt"
  "log"

  "github.com/jonwraymond/toolcompose/set"
  "github.com/jonwraymond/toolfoundation/adapter"
)

tools := []*adapter.CanonicalTool{
  {
    Namespace: "github",
    Name:      "create_issue",
    Tags:      []string{"issues"},
    InputSchema: &adapter.JSONSchema{Type: "object"},
  },
  {
    Namespace: "github",
    Name:      "list_issues",
    Tags:      []string{"issues"},
    InputSchema: &adapter.JSONSchema{Type: "object"},
  },
}

ts, err := set.NewBuilder("github-issues").
  FromTools(tools).
  WithNamespace("github").
  WithTags([]string{"issues"}).
  WithPolicy(set.DenyTags("danger")).
  Build()
if err != nil {
  log.Fatal(err)
}

fmt.Println(ts.IDs())
```

## Quick Start: Skill

```go
import (
  "context"
  "log"

  "github.com/jonwraymond/toolcompose/skill"
  "github.com/jonwraymond/toolexec/run"
)

// Define a skill
sk := skill.Skill{
  Name: "triage-issue",
  Steps: []skill.Step{
    {ID: "create", ToolID: "github:create_issue", Inputs: map[string]any{"title": "Bug"}},
    {ID: "label", ToolID: "github:add_labels", Inputs: map[string]any{"labels": []string{"bug"}}},
  },
}

plan, err := skill.NewPlanner().Plan(sk)
if err != nil {
  log.Fatal(err)
}

runner := run.NewRunner()

// Adapt toolexec runner to skill.Runner
type runAdapter struct{ exec run.Runner }
func (r runAdapter) Run(ctx context.Context, step skill.Step) (any, error) {
  res, err := r.exec.Run(ctx, step.ToolID, step.Inputs)
  if err != nil {
    return nil, err
  }
  return res.Output, nil
}

ctx := context.Background()
results, err := skill.Execute(ctx, plan, runAdapter{exec: runner})
if err != nil {
  log.Fatal(err)
}
_ = results
```

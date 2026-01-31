# toolcompose User Journey

## 1. Installation

```bash
go get github.com/jonwraymond/toolcompose@latest
```

## 2. Build a Toolset

```go
import (
  "fmt"
  "log"

  "github.com/jonwraymond/toolcompose/set"
  "github.com/jonwraymond/toolfoundation/adapter"
)

tools := []*adapter.CanonicalTool{
  {Namespace: "github", Name: "create_issue", Tags: []string{"issues"}, InputSchema: &adapter.JSONSchema{Type: "object"}},
  {Namespace: "github", Name: "add_labels", Tags: []string{"issues"}, InputSchema: &adapter.JSONSchema{Type: "object"}},
  {Namespace: "slack", Name: "post_message", Tags: []string{"chat"}, InputSchema: &adapter.JSONSchema{Type: "object"}},
}

ts, err := set.NewBuilder("issue-workflow").
  FromTools(tools).
  WithNamespace("github").
  WithTags([]string{"issues"}).
  WithPolicy(set.DenyTags("danger")).
  Build()
if err != nil {
  log.Fatal(err)
}

fmt.Println("Toolset IDs:", ts.IDs())
```

## 3. Filter an Existing Toolset

```go
filtered := ts.Filter(set.TagsAny("issues"))
fmt.Println(filtered.IDs())
```

## 4. Export Toolsets

```go
import (
  "log"

  "github.com/jonwraymond/toolfoundation/adapter/mcp"
)

exposure := set.NewExposure(ts, mcp.NewAdapter())
tools, warnings, errs := exposure.ExportWithWarnings()

if len(errs) > 0 {
  log.Printf("conversion errors: %v", errs)
}
_ = warnings
_ = tools
```

## 5. Define a Skill

```go
import "github.com/jonwraymond/toolcompose/skill"

sk := skill.Skill{
  Name: "triage-issue",
  Steps: []skill.Step{
    {ID: "create", ToolID: "github:create_issue", Inputs: map[string]any{"title": "Bug"}},
    {ID: "label", ToolID: "github:add_labels", Inputs: map[string]any{"labels": []string{"bug"}}},
  },
}
```

## 6. Plan + Guard

```go
planner := skill.NewPlanner()
plan, err := planner.Plan(sk)
if err != nil {
  log.Fatal(err)
}

guard := skill.MaxStepsGuard(5)
if err := guard.Validate(sk); err != nil {
  log.Fatal(err)
}
```

## 7. Execute via toolexec

```go
import (
  "context"

  "github.com/jonwraymond/toolcompose/skill"
  "github.com/jonwraymond/toolexec/run"
)

type runAdapter struct{ exec run.Runner }
func (r runAdapter) Run(ctx context.Context, step skill.Step) (any, error) {
  res, err := r.exec.Run(ctx, step.ToolID, step.Inputs)
  if err != nil {
    return nil, err
  }
  return res.Output, nil
}

runner := run.NewRunner()
results, err := skill.Execute(context.Background(), plan, runAdapter{exec: runner})
if err != nil {
  log.Fatal(err)
}
_ = results
```

## Next Steps

- Combine with `toolcompose/set` to scope which tools a skill can access.
- Use `toolcompose/set` exposures to export protocol-specific tool lists.

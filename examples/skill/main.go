// Example: Skill planning and execution
//
// Demonstrates declarative skill definitions and deterministic execution.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jonwraymond/toolcompose/skill"
)

// MockRunner simulates tool execution for demonstration.
type MockRunner struct {
	results map[string]any
}

func (r *MockRunner) Run(ctx context.Context, step skill.Step) (any, error) {
	fmt.Printf("  Executing step '%s' (tool: %s)\n", step.ID, step.ToolID)
	fmt.Printf("    Inputs: %v\n", step.Inputs)

	// Simulate tool execution
	if result, ok := r.results[step.ID]; ok {
		fmt.Printf("    Result: %v\n", result)
		return result, nil
	}
	return nil, nil
}

func main() {
	// Define a multi-step skill
	s := skill.Skill{
		Name: "create-issue-from-search",
		Steps: []skill.Step{
			{
				ID:     "1-search",
				ToolID: "github:search",
				Inputs: map[string]any{
					"query": "is:issue label:bug",
					"limit": 10,
				},
			},
			{
				ID:     "2-analyze",
				ToolID: "ai:summarize",
				Inputs: map[string]any{
					"text": "{{1-search.result}}",
				},
			},
			{
				ID:     "3-create",
				ToolID: "github:create-issue",
				Inputs: map[string]any{
					"title": "Bug Summary",
					"body":  "{{2-analyze.result}}",
				},
			},
		},
	}

	fmt.Println("=== Skill Validation ===")
	if err := s.Validate(); err != nil {
		log.Fatalf("Invalid skill: %v", err)
	}
	fmt.Printf("Skill '%s' is valid with %d steps\n", s.Name, len(s.Steps))

	fmt.Println("\n=== Planning ===")
	planner := skill.NewPlanner()
	plan, err := planner.Plan(s)
	if err != nil {
		log.Fatalf("Planning failed: %v", err)
	}
	fmt.Printf("Plan created with %d steps (sorted by ID):\n", len(plan.Steps))
	for i, step := range plan.Steps {
		fmt.Printf("  %d. %s -> %s\n", i+1, step.ID, step.ToolID)
	}

	fmt.Println("\n=== Execution ===")
	runner := &MockRunner{
		results: map[string]any{
			"1-search":  []string{"issue-123", "issue-456"},
			"2-analyze": "Found 2 critical bugs related to auth",
			"3-create":  "issue-789",
		},
	}

	ctx := context.Background()
	results, err := skill.Execute(ctx, plan, runner)
	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	fmt.Println("\n=== Results ===")
	for _, r := range results {
		status := "OK"
		if r.Err != nil {
			status = fmt.Sprintf("ERROR: %v", r.Err)
		}
		fmt.Printf("  Step '%s': %s (value: %v)\n", r.StepID, status, r.Value)
	}
}

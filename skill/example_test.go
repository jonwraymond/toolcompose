package skill_test

import (
	"context"
	"fmt"

	"github.com/jonwraymond/toolcompose/skill"
)

// mockRunner implements skill.Runner for examples.
type mockRunner struct {
	results map[string]any
}

func (r *mockRunner) Run(ctx context.Context, step skill.Step) (any, error) {
	if result, ok := r.results[step.ID]; ok {
		return result, nil
	}
	return fmt.Sprintf("executed %s", step.ToolID), nil
}

func ExampleNewPlanner() {
	s := skill.Skill{
		Name: "greet-workflow",
		Steps: []skill.Step{
			{ID: "step-b", ToolID: "ns:greet", Inputs: map[string]any{"name": "World"}},
			{ID: "step-a", ToolID: "ns:log", Inputs: map[string]any{"msg": "starting"}},
		},
	}

	planner := skill.NewPlanner()
	plan, err := planner.Plan(s)
	if err != nil {
		panic(err)
	}

	// Steps are sorted by ID
	fmt.Println("Plan name:", plan.Name)
	fmt.Println("First step:", plan.Steps[0].ID)
	fmt.Println("Second step:", plan.Steps[1].ID)
	// Output:
	// Plan name: greet-workflow
	// First step: step-a
	// Second step: step-b
}

func ExampleExecute() {
	plan := skill.Plan{
		Name: "demo",
		Steps: []skill.Step{
			{ID: "1", ToolID: "ns:tool1"},
			{ID: "2", ToolID: "ns:tool2"},
		},
	}

	runner := &mockRunner{
		results: map[string]any{
			"1": "result-1",
			"2": "result-2",
		},
	}

	results, err := skill.Execute(context.Background(), plan, runner)
	if err != nil {
		panic(err)
	}

	fmt.Println("Results count:", len(results))
	fmt.Println("Step 1 result:", results[0].Value)
	fmt.Println("Step 2 result:", results[1].Value)
	// Output:
	// Results count: 2
	// Step 1 result: result-1
	// Step 2 result: result-2
}

func ExampleMaxStepsGuard() {
	guard := skill.MaxStepsGuard(2)

	smallSkill := skill.Skill{
		Name:  "small",
		Steps: []skill.Step{{ID: "1", ToolID: "ns:tool"}},
	}
	largeSkill := skill.Skill{
		Name: "large",
		Steps: []skill.Step{
			{ID: "1", ToolID: "ns:tool"},
			{ID: "2", ToolID: "ns:tool"},
			{ID: "3", ToolID: "ns:tool"},
		},
	}

	fmt.Println("Small skill valid:", guard.Validate(smallSkill) == nil)
	fmt.Println("Large skill valid:", guard.Validate(largeSkill) == nil)
	// Output:
	// Small skill valid: true
	// Large skill valid: false
}

func ExampleAllowedToolIDsGuard() {
	guard := skill.AllowedToolIDsGuard([]string{"safe:read", "safe:list"})

	safeSkill := skill.Skill{
		Name:  "safe",
		Steps: []skill.Step{{ID: "1", ToolID: "safe:read"}},
	}
	unsafeSkill := skill.Skill{
		Name:  "unsafe",
		Steps: []skill.Step{{ID: "1", ToolID: "dangerous:delete"}},
	}

	fmt.Println("Safe skill valid:", guard.Validate(safeSkill) == nil)
	fmt.Println("Unsafe skill valid:", guard.Validate(unsafeSkill) == nil)
	// Output:
	// Safe skill valid: true
	// Unsafe skill valid: false
}

func ExampleSkill_Validate() {
	validSkill := skill.Skill{
		Name:  "valid",
		Steps: []skill.Step{{ID: "1", ToolID: "ns:tool"}},
	}

	invalidSkill := skill.Skill{
		Name:  "", // Empty name
		Steps: []skill.Step{{ID: "1", ToolID: "ns:tool"}},
	}

	fmt.Println("Valid skill:", validSkill.Validate() == nil)
	fmt.Println("Invalid skill:", invalidSkill.Validate() == nil)
	// Output:
	// Valid skill: true
	// Invalid skill: false
}

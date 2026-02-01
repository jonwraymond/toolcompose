// Package skill provides skill composition primitives built on tool execution.
//
// A skill is a reusable workflow composed of tool calls with pre/post
// conditions and validation guards. This library defines the orchestration
// model; actual tool execution is delegated to a Runner implementation.
//
// # Core Concepts
//
//   - Skill: Declarative workflow definition with named steps
//   - Step: Individual tool invocation with ID, ToolID, and Inputs
//   - Plan: Compiled, deterministically-ordered execution plan
//   - Planner: Compiles Skills into Plans (validates and sorts by step ID)
//   - Runner: Executes individual steps (user-provided implementation)
//   - Guard: Validates skills before execution (max steps, allowed tools)
//
// # Basic Usage
//
//	// Define a skill
//	s := skill.Skill{
//	    Name: "create-and-label-issue",
//	    Steps: []skill.Step{
//	        {ID: "1", ToolID: "github:create_issue", Inputs: map[string]any{"title": "Bug"}},
//	        {ID: "2", ToolID: "github:add_labels", Inputs: map[string]any{"labels": []string{"bug"}}},
//	    },
//	}
//
//	// Create a plan (validates and sorts steps)
//	planner := skill.NewPlanner()
//	plan, err := planner.Plan(s)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Execute the plan
//	results, err := skill.Execute(ctx, plan, myRunner)
//
// # Guards
//
// Guards validate skills before execution:
//
//	// Limit maximum steps
//	guard := skill.MaxStepsGuard(10)
//	if err := guard.Validate(mySkill); err != nil {
//	    // Handle validation failure
//	}
//
//	// Restrict to allowed tools
//	guard := skill.AllowedToolIDsGuard([]string{"github:create_issue", "github:list_repos"})
//
// # Runner Interface
//
// Implement Runner to execute steps:
//
//	type myRunner struct {
//	    exec *toolexec.Exec
//	}
//
//	func (r *myRunner) Run(ctx context.Context, step skill.Step) (any, error) {
//	    result, err := r.exec.RunTool(ctx, step.ToolID, step.Inputs)
//	    if err != nil {
//	        return nil, err
//	    }
//	    return result.Value, nil
//	}
//
// # Execution Model
//
// Execute runs steps sequentially in plan order:
//   - Steps are executed one at a time
//   - Execution stops on first error (fail-fast)
//   - Partial results are returned on error
//   - Context cancellation is respected
//
// # Determinism
//
// Plans are deterministic: steps are sorted by ID (lexicographically).
// This ensures reproducible execution order regardless of definition order.
//
// # Thread Safety
//
//   - Skill, Step, Plan: Immutable after creation (safe for concurrent read)
//   - Planner: Safe for concurrent use
//   - Runner: Must be safe for concurrent use (per contract)
//   - Guard: Safe after construction (immutable closures)
//
// # Error Handling
//
// Sentinel errors for programmatic handling:
//
//   - ErrInvalidRunner: Execute called with nil runner
//   - ErrInvalidSkillName: Skill has empty name
//   - ErrInvalidStepID: Step has empty ID
//   - ErrInvalidToolID: Step has empty ToolID
//   - ErrNoSteps: Skill has no steps
//   - ErrMaxStepsExceeded: MaxStepsGuard limit exceeded
//   - ErrToolNotAllowed: AllowedToolIDsGuard rejected tool
//
// Use errors.Is() for reliable error checking:
//
//	if errors.Is(err, skill.ErrNoSteps) {
//	    // Handle empty skill
//	}
package skill

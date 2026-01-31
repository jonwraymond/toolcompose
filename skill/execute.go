package skill

import "context"

// Runner executes a single step.
//
// Contract:
// - Concurrency: implementations must be safe for concurrent use.
// - Context: must honor cancellation/deadlines and return ctx.Err() when canceled.
// - Errors: execution failures should be returned directly; callers may wrap.
type Runner interface {
	Run(ctx context.Context, step Step) (any, error)
}

// StepResult captures execution results.
type StepResult struct {
	StepID string
	Value  any
	Err    error
}

// Execute runs the plan in order using the provided runner.
func Execute(ctx context.Context, plan Plan, runner Runner) ([]StepResult, error) {
	if runner == nil {
		return nil, ErrInvalidRunner
	}

	results := make([]StepResult, 0, len(plan.Steps))
	for _, step := range plan.Steps {
		val, err := runner.Run(ctx, step)
		res := StepResult{StepID: step.ID, Value: val, Err: err}
		results = append(results, res)
		if err != nil {
			return results, err
		}
	}
	return results, nil
}

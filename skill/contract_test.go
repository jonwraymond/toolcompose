package skill

import (
	"context"
	"testing"
)

func TestRunnerContract_ExecuteUsesRunner(t *testing.T) {
	runner := &mockRunner{}
	plan := Plan{Steps: []Step{{ID: "one"}}}
	_, err := Execute(context.Background(), plan, runner)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected runner to be called once, got %d", len(runner.calls))
	}
}

func TestGuardContract_Basic(t *testing.T) {
	guard := MaxStepsGuard(1)
	err := guard.Validate(Skill{Steps: []Step{{ID: "one"}, {ID: "two"}}})
	if err == nil {
		t.Fatalf("expected guard validation error")
	}
}

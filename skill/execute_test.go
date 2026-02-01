package skill

import (
	"context"
	"errors"
	"testing"
)

type ctxKey string

const testCtxKey ctxKey = "key"

type mockRunner struct {
	calls []string
	fail  bool
}

func (m *mockRunner) Run(ctx context.Context, step Step) (any, error) {
	m.calls = append(m.calls, step.ID)
	if m.fail {
		return nil, errors.New("failed")
	}
	return step.ID + "-ok", nil
}

func TestExecute_OrderPreserved(t *testing.T) {
	planner := NewPlanner()
	skill := Skill{
		Name: "workflow",
		Steps: []Step{
			{ID: "b", ToolID: "t2"},
			{ID: "a", ToolID: "t1"},
		},
	}
	plan, err := planner.Plan(skill)
	if err != nil {
		t.Fatalf("plan failed: %v", err)
	}

	runner := &mockRunner{}
	results, err := Execute(context.Background(), plan, runner)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	want := []string{"a", "b"}
	if len(runner.calls) != len(want) {
		t.Fatalf("calls = %d, want %d", len(runner.calls), len(want))
	}
	for i := range want {
		if runner.calls[i] != want[i] {
			t.Fatalf("call[%d] = %q, want %q", i, runner.calls[i], want[i])
		}
	}
	if len(results) != 2 {
		t.Fatalf("results length = %d, want 2", len(results))
	}
}

func TestExecute_ErrorPropagation(t *testing.T) {
	plan := Plan{Name: "workflow", Steps: []Step{{ID: "a", ToolID: "t1"}}}
	runner := &mockRunner{fail: true}

	_, err := Execute(context.Background(), plan, runner)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestExecute_ContextPropagation(t *testing.T) {
	plan := Plan{Name: "workflow", Steps: []Step{{ID: "a", ToolID: "t1"}}}

	runner := &mockRunner{}
	ctx := context.WithValue(context.Background(), testCtxKey, "value")
	_, err := Execute(ctx, plan, runner)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
}

func TestExecute_NilRunner(t *testing.T) {
	plan := Plan{Name: "workflow", Steps: []Step{{ID: "a", ToolID: "t1"}}}

	_, err := Execute(context.Background(), plan, nil)
	if !errors.Is(err, ErrInvalidRunner) {
		t.Fatalf("expected ErrInvalidRunner, got: %v", err)
	}
}

func TestExecute_EmptyPlan(t *testing.T) {
	plan := Plan{Name: "workflow", Steps: []Step{}}
	runner := &mockRunner{}

	results, err := Execute(context.Background(), plan, runner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestExecute_PartialResultsOnError(t *testing.T) {
	// Runner that fails on second step
	runner := &partialFailRunner{failOnStep: "b"}
	plan := Plan{
		Name: "workflow",
		Steps: []Step{
			{ID: "a", ToolID: "t1"},
			{ID: "b", ToolID: "t2"},
			{ID: "c", ToolID: "t3"},
		},
	}

	results, err := Execute(context.Background(), plan, runner)
	if err == nil {
		t.Fatalf("expected error")
	}
	// Should have partial results (a succeeded, b failed)
	if len(results) != 2 {
		t.Fatalf("expected 2 partial results, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("first result should not have error")
	}
	if results[1].Err == nil {
		t.Fatalf("second result should have error")
	}
}

type partialFailRunner struct {
	failOnStep string
}

func (r *partialFailRunner) Run(ctx context.Context, step Step) (any, error) {
	if step.ID == r.failOnStep {
		return nil, errors.New("step failed")
	}
	return step.ID + "-ok", nil
}

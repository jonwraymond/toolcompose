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

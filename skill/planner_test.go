package skill

import "testing"

func TestPlanner_DeterministicOrdering(t *testing.T) {
	skill := Skill{
		Name: "workflow",
		Steps: []Step{
			{ID: "b", ToolID: "t2"},
			{ID: "a", ToolID: "t1"},
			{ID: "c", ToolID: "t3"},
		},
	}

	plan, err := NewPlanner().Plan(skill)
	if err != nil {
		t.Fatalf("plan failed: %v", err)
	}

	want := []string{"a", "b", "c"}
	if len(plan.Steps) != len(want) {
		t.Fatalf("plan steps = %d, want %d", len(plan.Steps), len(want))
	}
	for i := range want {
		if plan.Steps[i].ID != want[i] {
			t.Fatalf("step[%d] = %q, want %q", i, plan.Steps[i].ID, want[i])
		}
	}
}

func TestPlanner_DetectsMissingSteps(t *testing.T) {
	skill := Skill{Name: "workflow"}
	_, err := NewPlanner().Plan(skill)
	if err == nil {
		t.Fatalf("expected error for empty steps")
	}
}

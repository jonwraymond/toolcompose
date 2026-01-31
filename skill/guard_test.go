package skill

import "testing"

func TestGuard_MaxSteps(t *testing.T) {
	guard := MaxStepsGuard(2)

	skill := Skill{
		Name: "test",
		Steps: []Step{{ID: "a", ToolID: "t1"}, {ID: "b", ToolID: "t2"}, {ID: "c", ToolID: "t3"}},
	}

	if err := guard.Validate(skill); err == nil {
		t.Fatalf("expected error for exceeding max steps")
	}
}

func TestGuard_AllowedToolIDs(t *testing.T) {
	guard := AllowedToolIDsGuard([]string{"t1", "t2"})

	skill := Skill{
		Name: "test",
		Steps: []Step{{ID: "a", ToolID: "t1"}, {ID: "b", ToolID: "t3"}},
	}

	if err := guard.Validate(skill); err == nil {
		t.Fatalf("expected error for disallowed tool id")
	}
}

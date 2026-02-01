package skill

import "testing"

func TestGuard_MaxSteps(t *testing.T) {
	guard := MaxStepsGuard(2)

	skill := Skill{
		Name:  "test",
		Steps: []Step{{ID: "a", ToolID: "t1"}, {ID: "b", ToolID: "t2"}, {ID: "c", ToolID: "t3"}},
	}

	if err := guard.Validate(skill); err == nil {
		t.Fatalf("expected error for exceeding max steps")
	}
}

func TestGuard_AllowedToolIDs(t *testing.T) {
	guard := AllowedToolIDsGuard([]string{"t1", "t2"})

	skill := Skill{
		Name:  "test",
		Steps: []Step{{ID: "a", ToolID: "t1"}, {ID: "b", ToolID: "t3"}},
	}

	if err := guard.Validate(skill); err == nil {
		t.Fatalf("expected error for disallowed tool id")
	}
}

func TestGuard_MaxSteps_ZeroOrNegative(t *testing.T) {
	tests := []struct {
		name string
		max  int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guard := MaxStepsGuard(tt.max)
			skill := Skill{
				Name:  "test",
				Steps: []Step{{ID: "a", ToolID: "t1"}},
			}
			// Should allow all when max <= 0
			if err := guard.Validate(skill); err != nil {
				t.Fatalf("unexpected error for max=%d: %v", tt.max, err)
			}
		})
	}
}

func TestGuard_MaxSteps_ExactLimit(t *testing.T) {
	guard := MaxStepsGuard(2)

	// Exactly at limit should pass
	skill := Skill{
		Name:  "test",
		Steps: []Step{{ID: "a", ToolID: "t1"}, {ID: "b", ToolID: "t2"}},
	}
	if err := guard.Validate(skill); err != nil {
		t.Fatalf("unexpected error at exact limit: %v", err)
	}
}

func TestGuard_AllowedToolIDs_AllAllowed(t *testing.T) {
	guard := AllowedToolIDsGuard([]string{"t1", "t2", "t3"})

	skill := Skill{
		Name:  "test",
		Steps: []Step{{ID: "a", ToolID: "t1"}, {ID: "b", ToolID: "t2"}},
	}

	if err := guard.Validate(skill); err != nil {
		t.Fatalf("unexpected error when all tools allowed: %v", err)
	}
}

func TestGuard_AllowedToolIDs_EmptyAllowlist(t *testing.T) {
	guard := AllowedToolIDsGuard([]string{})

	skill := Skill{
		Name:  "test",
		Steps: []Step{{ID: "a", ToolID: "t1"}},
	}

	if err := guard.Validate(skill); err == nil {
		t.Fatalf("expected error with empty allowlist")
	}
}

func TestGuard_AllowedToolIDs_EmptySteps(t *testing.T) {
	guard := AllowedToolIDsGuard([]string{"t1"})

	skill := Skill{
		Name:  "test",
		Steps: []Step{},
	}

	// Empty steps should pass (nothing to validate)
	if err := guard.Validate(skill); err != nil {
		t.Fatalf("unexpected error with empty steps: %v", err)
	}
}

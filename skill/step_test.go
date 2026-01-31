package skill

import "testing"

func TestStep_ValidatesToolID(t *testing.T) {
	step := Step{ID: "s1", ToolID: ""}
	if err := step.Validate(); err == nil {
		t.Fatalf("expected validation error for empty tool id")
	}
}

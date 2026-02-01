package skill

import (
	"errors"
	"testing"
)

func TestStep_ValidatesToolID(t *testing.T) {
	step := Step{ID: "s1", ToolID: ""}
	if err := step.Validate(); err == nil {
		t.Fatalf("expected validation error for empty tool id")
	}
}

func TestStep_ValidatesStepID(t *testing.T) {
	step := Step{ID: "", ToolID: "ns:tool"}
	err := step.Validate()
	if err == nil {
		t.Fatalf("expected validation error for empty step id")
	}
	if !errors.Is(err, ErrInvalidStepID) {
		t.Fatalf("got %v, want ErrInvalidStepID", err)
	}
}

func TestStep_Validate_Valid(t *testing.T) {
	step := Step{ID: "s1", ToolID: "ns:tool", Inputs: map[string]any{"key": "value"}}
	if err := step.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestStep_Validate_BothEmpty(t *testing.T) {
	step := Step{ID: "", ToolID: ""}
	err := step.Validate()
	if err == nil {
		t.Fatalf("expected validation error")
	}
	// First check is ID, so should get ErrInvalidStepID
	if !errors.Is(err, ErrInvalidStepID) {
		t.Fatalf("got %v, want ErrInvalidStepID", err)
	}
}

func TestStep_Validate_NilInputs(t *testing.T) {
	step := Step{ID: "s1", ToolID: "ns:tool", Inputs: nil}
	if err := step.Validate(); err != nil {
		t.Fatalf("unexpected validation error for nil inputs: %v", err)
	}
}

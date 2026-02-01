package skill

import (
	"encoding/json"
	"testing"
)

func TestSkill_SerializationDeterministic(t *testing.T) {
	skill := Skill{
		Name: "summarize",
		Steps: []Step{
			{ID: "search", ToolID: "mcp:search", Inputs: map[string]any{"q": "foo"}},
			{ID: "summarize", ToolID: "mcp:summarize", Inputs: map[string]any{"text": "bar"}},
		},
	}

	b1, err := json.Marshal(skill)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	b2, err := json.Marshal(skill)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if string(b1) != string(b2) {
		t.Fatalf("serialization not deterministic:\n%s\n%s", string(b1), string(b2))
	}
}

func TestSkill_Validate(t *testing.T) {
	skill := Skill{
		Name:  "",
		Steps: []Step{{ID: "s1", ToolID: "mcp:search"}},
	}

	if err := skill.Validate(); err == nil {
		t.Fatalf("expected validation error for empty name")
	}
}

func TestSkill_Validate_Valid(t *testing.T) {
	skill := Skill{
		Name:  "valid-skill",
		Steps: []Step{{ID: "s1", ToolID: "mcp:search"}},
	}

	if err := skill.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestSkill_Validate_InvalidStep(t *testing.T) {
	tests := []struct {
		name    string
		skill   Skill
		wantErr error
	}{
		{
			name: "empty step ID",
			skill: Skill{
				Name:  "skill",
				Steps: []Step{{ID: "", ToolID: "mcp:search"}},
			},
			wantErr: ErrInvalidStepID,
		},
		{
			name: "empty tool ID",
			skill: Skill{
				Name:  "skill",
				Steps: []Step{{ID: "s1", ToolID: ""}},
			},
			wantErr: ErrInvalidToolID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.skill.Validate()
			if err == nil {
				t.Fatalf("expected error")
			}
			if err != tt.wantErr {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSkill_Validate_EmptySteps(t *testing.T) {
	skill := Skill{
		Name:  "valid-skill",
		Steps: []Step{},
	}

	// Validation should pass - no steps is valid for Skill.Validate
	// (Planner.Plan will reject empty steps)
	if err := skill.Validate(); err != nil {
		t.Fatalf("unexpected error for empty steps: %v", err)
	}
}

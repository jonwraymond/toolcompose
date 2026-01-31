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
		Name: "",
		Steps: []Step{{ID: "s1", ToolID: "mcp:search"}},
	}

	if err := skill.Validate(); err == nil {
		t.Fatalf("expected validation error for empty name")
	}
}

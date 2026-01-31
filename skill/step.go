package skill

import "errors"

var ErrInvalidStepID = errors.New("skill: step id is required")
var ErrInvalidToolID = errors.New("skill: tool id is required")

// Step references a tool and its bindings.
type Step struct {
	ID     string
	ToolID string
	Inputs map[string]any
}

// Validate validates a step.
func (s Step) Validate() error {
	if s.ID == "" {
		return ErrInvalidStepID
	}
	if s.ToolID == "" {
		return ErrInvalidToolID
	}
	return nil
}

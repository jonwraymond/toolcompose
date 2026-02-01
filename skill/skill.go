package skill

// Skill represents a declarative workflow.
type Skill struct {
	Name  string
	Steps []Step
}

// Validate validates the skill definition.
func (s Skill) Validate() error {
	if s.Name == "" {
		return ErrInvalidSkillName
	}
	for _, step := range s.Steps {
		if err := step.Validate(); err != nil {
			return err
		}
	}
	return nil
}

package skill

import "sort"

// Plan is the compiled, deterministic execution plan.
type Plan struct {
	Name  string
	Steps []Step
}

// Planner produces deterministic plans from skills.
type Planner struct{}

// NewPlanner creates a Planner.
func NewPlanner() *Planner {
	return &Planner{}
}

// Plan validates and orders steps deterministically by ID.
func (p *Planner) Plan(skill Skill) (Plan, error) {
	if err := skill.Validate(); err != nil {
		return Plan{}, err
	}
	if len(skill.Steps) == 0 {
		return Plan{}, ErrNoSteps
	}

	steps := make([]Step, len(skill.Steps))
	copy(steps, skill.Steps)
	sort.Slice(steps, func(i, j int) bool { return steps[i].ID < steps[j].ID })

	return Plan{Name: skill.Name, Steps: steps}, nil
}

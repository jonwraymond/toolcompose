package skill

// Guard validates a skill or step.
//
// Contract:
//   - Concurrency: implementations must be safe for concurrent use.
//   - Errors: validation failures must return non-nil error; must not panic.
//   - Determinism: same skill should yield same result.
type Guard interface {
	Validate(skill Skill) error
}

type guardFunc func(skill Skill) error

func (g guardFunc) Validate(skill Skill) error { return g(skill) }

// MaxStepsGuard enforces a maximum number of steps.
func MaxStepsGuard(max int) Guard {
	return guardFunc(func(skill Skill) error {
		if max <= 0 {
			return nil
		}
		if len(skill.Steps) > max {
			return ErrMaxStepsExceeded
		}
		return nil
	})
}

// AllowedToolIDsGuard restricts steps to a set of tool IDs.
func AllowedToolIDsGuard(allowed []string) Guard {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, id := range allowed {
		allowedSet[id] = struct{}{}
	}

	return guardFunc(func(skill Skill) error {
		for _, step := range skill.Steps {
			if _, ok := allowedSet[step.ToolID]; !ok {
				return ErrToolNotAllowed
			}
		}
		return nil
	})
}

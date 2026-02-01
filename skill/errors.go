package skill

import "errors"

// Sentinel errors for the skill package.
var (
	// ErrInvalidRunner is returned when Execute is called with a nil runner.
	ErrInvalidRunner = errors.New("skill: runner is required")

	// ErrInvalidSkillName is returned when a Skill has an empty name.
	ErrInvalidSkillName = errors.New("skill: skill name is required")

	// ErrInvalidStepID is returned when a Step has an empty ID.
	ErrInvalidStepID = errors.New("skill: step id is required")

	// ErrInvalidToolID is returned when a Step has an empty ToolID.
	ErrInvalidToolID = errors.New("skill: tool id is required")

	// ErrNoSteps is returned when planning a Skill with no steps.
	ErrNoSteps = errors.New("skill: skill has no steps")

	// ErrMaxStepsExceeded is returned by MaxStepsGuard when limit is exceeded.
	ErrMaxStepsExceeded = errors.New("skill: max steps exceeded")

	// ErrToolNotAllowed is returned by AllowedToolIDsGuard for unlisted tools.
	ErrToolNotAllowed = errors.New("skill: tool id not allowed")
)

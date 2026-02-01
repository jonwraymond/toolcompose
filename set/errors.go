package set

import "errors"

// Sentinel errors for the set package.
var (
	// ErrNoSource is returned when Build() is called without a source.
	ErrNoSource = errors.New("set: no source provided (call FromTools or FromRegistry)")

	// ErrNilAdapter is returned when an Exposure is created with a nil adapter.
	ErrNilAdapter = errors.New("set: adapter is nil")

	// ErrNilToolset is returned when an operation requires a non-nil Toolset.
	ErrNilToolset = errors.New("set: toolset is nil")
)

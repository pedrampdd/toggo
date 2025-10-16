package toggo

import "errors"

var (
	// ErrFlagNotFound is returned when a requested flag doesn't exist in the store
	ErrFlagNotFound = errors.New("flag not found")

	// ErrInvalidOperator is returned when an unsupported operator is encountered
	ErrInvalidOperator = errors.New("invalid operator")

	// ErrInvalidRollout is returned when rollout percentage is not between 0 and 100
	ErrInvalidRollout = errors.New("rollout must be between 0 and 100")

	// ErrInvalidCondition is returned when a condition is malformed
	ErrInvalidCondition = errors.New("invalid condition")

	// ErrRolloutKeyMissing is returned when the specified rollout key is not in context
	ErrRolloutKeyMissing = errors.New("rollout key missing from context")
)

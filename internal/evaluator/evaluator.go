package evaluator

import "github.com/pedrampdd/toggo"

// Evaluator defines the interface for condition evaluation
type Evaluator interface {
	// Evaluate checks if a condition matches the given context
	Evaluate(condition toggo.Condition, ctx toggo.Context) (bool, error)

	// EvaluateAll checks if all conditions match the given context
	EvaluateAll(conditions []toggo.Condition, ctx toggo.Context) (bool, error)
}

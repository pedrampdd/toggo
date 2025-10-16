package evaluator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pedrampdd/toggo"
)

// StandardEvaluator is the default implementation of the Evaluator interface
type StandardEvaluator struct{}

// NewStandard creates a new standard evaluator
func NewStandard() *StandardEvaluator {
	return &StandardEvaluator{}
}

// Evaluate checks if a single condition matches the context
func (e *StandardEvaluator) Evaluate(condition toggo.Condition, ctx toggo.Context) (bool, error) {
	if err := condition.Validate(); err != nil {
		return false, err
	}

	value, exists := ctx.Get(condition.Attribute)
	if !exists {
		// If attribute doesn't exist in context, condition fails
		return e.applyNegate(false, condition.Negate), nil
	}

	result, err := e.evaluateOperator(condition.Operator, value, condition.Value)
	if err != nil {
		return false, err
	}

	return e.applyNegate(result, condition.Negate), nil
}

// EvaluateAll checks if all conditions match (AND logic)
func (e *StandardEvaluator) EvaluateAll(conditions []toggo.Condition, ctx toggo.Context) (bool, error) {
	for _, cond := range conditions {
		match, err := e.Evaluate(cond, ctx)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

// applyNegate applies negation to the result if negate is true
func (e *StandardEvaluator) applyNegate(result, negate bool) bool {
	if negate {
		return !result
	}
	return result
}

// evaluateOperator performs the actual comparison based on operator
func (e *StandardEvaluator) evaluateOperator(op toggo.Operator, ctxValue, condValue interface{}) (bool, error) {
	switch op {
	case toggo.OperatorEqual:
		return e.evaluateEqual(ctxValue, condValue), nil
	case toggo.OperatorNotEqual:
		return !e.evaluateEqual(ctxValue, condValue), nil
	case toggo.OperatorIn:
		return e.evaluateIn(ctxValue, condValue), nil
	case toggo.OperatorNotIn:
		return !e.evaluateIn(ctxValue, condValue), nil
	case toggo.OperatorGreaterThan:
		return e.evaluateGreaterThan(ctxValue, condValue, false), nil
	case toggo.OperatorGreaterThanOrEqual:
		return e.evaluateGreaterThan(ctxValue, condValue, true), nil
	case toggo.OperatorLessThan:
		return e.evaluateLessThan(ctxValue, condValue, false), nil
	case toggo.OperatorLessThanOrEqual:
		return e.evaluateLessThan(ctxValue, condValue, true), nil
	case toggo.OperatorContains:
		return e.evaluateContains(ctxValue, condValue), nil
	case toggo.OperatorStartsWith:
		return e.evaluateStartsWith(ctxValue, condValue), nil
	case toggo.OperatorEndsWith:
		return e.evaluateEndsWith(ctxValue, condValue), nil
	case toggo.OperatorRegex:
		return e.evaluateRegex(ctxValue, condValue)
	default:
		return false, toggo.ErrInvalidOperator
	}
}

// evaluateEqual checks equality
func (e *StandardEvaluator) evaluateEqual(ctxValue, condValue interface{}) bool {
	return fmt.Sprint(ctxValue) == fmt.Sprint(condValue)
}

// evaluateIn checks if value is in a list
func (e *StandardEvaluator) evaluateIn(ctxValue, condValue interface{}) bool {
	ctxStr := fmt.Sprint(ctxValue)

	// Handle slice of interfaces
	switch v := condValue.(type) {
	case []interface{}:
		for _, item := range v {
			if fmt.Sprint(item) == ctxStr {
				return true
			}
		}
	case []string:
		for _, item := range v {
			if item == ctxStr {
				return true
			}
		}
	default:
		// If it's not a slice, treat as single value comparison
		return e.evaluateEqual(ctxValue, condValue)
	}

	return false
}

// evaluateGreaterThan checks if context value is greater than condition value
func (e *StandardEvaluator) evaluateGreaterThan(ctxValue, condValue interface{}, orEqual bool) bool {
	ctxNum, err1 := e.toFloat64(ctxValue)
	condNum, err2 := e.toFloat64(condValue)

	if err1 != nil || err2 != nil {
		// Fallback to string comparison
		ctxStr := fmt.Sprint(ctxValue)
		condStr := fmt.Sprint(condValue)
		if orEqual {
			return ctxStr >= condStr
		}
		return ctxStr > condStr
	}

	if orEqual {
		return ctxNum >= condNum
	}
	return ctxNum > condNum
}

// evaluateLessThan checks if context value is less than condition value
func (e *StandardEvaluator) evaluateLessThan(ctxValue, condValue interface{}, orEqual bool) bool {
	ctxNum, err1 := e.toFloat64(ctxValue)
	condNum, err2 := e.toFloat64(condValue)

	if err1 != nil || err2 != nil {
		// Fallback to string comparison
		ctxStr := fmt.Sprint(ctxValue)
		condStr := fmt.Sprint(condValue)
		if orEqual {
			return ctxStr <= condStr
		}
		return ctxStr < condStr
	}

	if orEqual {
		return ctxNum <= condNum
	}
	return ctxNum < condNum
}

// evaluateContains checks if context string contains condition string
func (e *StandardEvaluator) evaluateContains(ctxValue, condValue interface{}) bool {
	ctxStr := fmt.Sprint(ctxValue)
	condStr := fmt.Sprint(condValue)
	return strings.Contains(ctxStr, condStr)
}

// evaluateStartsWith checks if context string starts with condition string
func (e *StandardEvaluator) evaluateStartsWith(ctxValue, condValue interface{}) bool {
	ctxStr := fmt.Sprint(ctxValue)
	condStr := fmt.Sprint(condValue)
	return strings.HasPrefix(ctxStr, condStr)
}

// evaluateEndsWith checks if context string ends with condition string
func (e *StandardEvaluator) evaluateEndsWith(ctxValue, condValue interface{}) bool {
	ctxStr := fmt.Sprint(ctxValue)
	condStr := fmt.Sprint(condValue)
	return strings.HasSuffix(ctxStr, condStr)
}

// evaluateRegex checks if context string matches regex pattern
func (e *StandardEvaluator) evaluateRegex(ctxValue, condValue interface{}) (bool, error) {
	ctxStr := fmt.Sprint(ctxValue)
	pattern := fmt.Sprint(condValue)

	matched, err := regexp.MatchString(pattern, ctxStr)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// toFloat64 converts interface{} to float64
func (e *StandardEvaluator) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

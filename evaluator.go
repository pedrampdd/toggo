package toggo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// conditionEvaluator handles the evaluation of conditions against contexts
type conditionEvaluator struct{}

// newConditionEvaluator creates a new condition evaluator
func newConditionEvaluator() *conditionEvaluator {
	return &conditionEvaluator{}
}

// evaluate checks if a single condition matches the context
func (e *conditionEvaluator) evaluate(condition Condition, ctx Context) (bool, error) {
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

// evaluateAll checks if all conditions match (AND logic)
func (e *conditionEvaluator) evaluateAll(conditions []Condition, ctx Context) (bool, error) {
	for _, cond := range conditions {
		match, err := e.evaluate(cond, ctx)
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
func (e *conditionEvaluator) applyNegate(result, negate bool) bool {
	if negate {
		return !result
	}
	return result
}

// evaluateOperator performs the actual comparison based on operator
func (e *conditionEvaluator) evaluateOperator(op Operator, ctxValue, condValue interface{}) (bool, error) {
	switch op {
	case OperatorEqual:
		return e.evaluateEqual(ctxValue, condValue), nil
	case OperatorNotEqual:
		return !e.evaluateEqual(ctxValue, condValue), nil
	case OperatorIn:
		return e.evaluateIn(ctxValue, condValue), nil
	case OperatorNotIn:
		return !e.evaluateIn(ctxValue, condValue), nil
	case OperatorGreaterThan:
		return e.evaluateGreaterThan(ctxValue, condValue, false), nil
	case OperatorGreaterThanOrEqual:
		return e.evaluateGreaterThan(ctxValue, condValue, true), nil
	case OperatorLessThan:
		return e.evaluateLessThan(ctxValue, condValue, false), nil
	case OperatorLessThanOrEqual:
		return e.evaluateLessThan(ctxValue, condValue, true), nil
	case OperatorContains:
		return e.evaluateContains(ctxValue, condValue), nil
	case OperatorStartsWith:
		return e.evaluateStartsWith(ctxValue, condValue), nil
	case OperatorEndsWith:
		return e.evaluateEndsWith(ctxValue, condValue), nil
	case OperatorRegex:
		return e.evaluateRegex(ctxValue, condValue)
	default:
		return false, ErrInvalidOperator
	}
}

// evaluateEqual checks equality
func (e *conditionEvaluator) evaluateEqual(ctxValue, condValue interface{}) bool {
	return fmt.Sprint(ctxValue) == fmt.Sprint(condValue)
}

// evaluateIn checks if value is in a list
func (e *conditionEvaluator) evaluateIn(ctxValue, condValue interface{}) bool {
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
func (e *conditionEvaluator) evaluateGreaterThan(ctxValue, condValue interface{}, orEqual bool) bool {
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
func (e *conditionEvaluator) evaluateLessThan(ctxValue, condValue interface{}, orEqual bool) bool {
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
func (e *conditionEvaluator) evaluateContains(ctxValue, condValue interface{}) bool {
	ctxStr := fmt.Sprint(ctxValue)
	condStr := fmt.Sprint(condValue)
	return strings.Contains(ctxStr, condStr)
}

// evaluateStartsWith checks if context string starts with condition string
func (e *conditionEvaluator) evaluateStartsWith(ctxValue, condValue interface{}) bool {
	ctxStr := fmt.Sprint(ctxValue)
	condStr := fmt.Sprint(condValue)
	return strings.HasPrefix(ctxStr, condStr)
}

// evaluateEndsWith checks if context string ends with condition string
func (e *conditionEvaluator) evaluateEndsWith(ctxValue, condValue interface{}) bool {
	ctxStr := fmt.Sprint(ctxValue)
	condStr := fmt.Sprint(condValue)
	return strings.HasSuffix(ctxStr, condStr)
}

// evaluateRegex checks if context string matches regex pattern
func (e *conditionEvaluator) evaluateRegex(ctxValue, condValue interface{}) (bool, error) {
	ctxStr := fmt.Sprint(ctxValue)
	pattern := fmt.Sprint(condValue)

	matched, err := regexp.MatchString(pattern, ctxStr)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// toFloat64 converts interface{} to float64
func (e *conditionEvaluator) toFloat64(value interface{}) (float64, error) {
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

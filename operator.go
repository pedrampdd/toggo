package toggo

// Operator represents a comparison operator for condition evaluation
type Operator string

const (
	// OperatorEqual checks if attribute equals value
	OperatorEqual Operator = "=="

	// OperatorNotEqual checks if attribute does not equal value
	OperatorNotEqual Operator = "!="

	// OperatorIn checks if attribute is in a list of values
	OperatorIn Operator = "in"

	// OperatorNotIn checks if attribute is not in a list of values
	OperatorNotIn Operator = "not_in"

	// OperatorGreaterThan checks if attribute is greater than value
	OperatorGreaterThan Operator = ">"

	// OperatorGreaterThanOrEqual checks if attribute is greater than or equal to value
	OperatorGreaterThanOrEqual Operator = ">="

	// OperatorLessThan checks if attribute is less than value
	OperatorLessThan Operator = "<"

	// OperatorLessThanOrEqual checks if attribute is less than or equal to value
	OperatorLessThanOrEqual Operator = "<="

	// OperatorContains checks if attribute string contains value
	OperatorContains Operator = "contains"

	// OperatorStartsWith checks if attribute string starts with value
	OperatorStartsWith Operator = "starts_with"

	// OperatorEndsWith checks if attribute string ends with value
	OperatorEndsWith Operator = "ends_with"

	// OperatorRegex checks if attribute matches regex pattern
	OperatorRegex Operator = "regex"
)

// IsValid checks if the operator is supported
func (o Operator) IsValid() bool {
	switch o {
	case OperatorEqual, OperatorNotEqual, OperatorIn, OperatorNotIn,
		OperatorGreaterThan, OperatorGreaterThanOrEqual,
		OperatorLessThan, OperatorLessThanOrEqual,
		OperatorContains, OperatorStartsWith, OperatorEndsWith,
		OperatorRegex:
		return true
	}
	return false
}

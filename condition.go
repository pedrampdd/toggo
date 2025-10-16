package toggo

// Condition represents a single evaluation condition
type Condition struct {
	// Attribute is the key to lookup in the context
	Attribute string `json:"attribute" yaml:"attribute"`

	// Operator is the comparison operator to use
	Operator Operator `json:"operator" yaml:"operator"`

	// Value is the value to compare against (can be string, number, array, etc.)
	Value interface{} `json:"value" yaml:"value"`

	// Negate inverts the condition result if true
	Negate bool `json:"negate,omitempty" yaml:"negate,omitempty"`
}

// Validate checks if the condition is properly formed
func (c *Condition) Validate() error {
	if c.Attribute == "" {
		return ErrInvalidCondition
	}
	if !c.Operator.IsValid() {
		return ErrInvalidOperator
	}
	return nil
}

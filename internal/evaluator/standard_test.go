package evaluator

import (
	"testing"

	"github.com/pedrampdd/toggo"
)

func TestStandardEvaluator_Equal(t *testing.T) {
	eval := NewStandard()

	tests := []struct {
		name      string
		condition toggo.Condition
		ctx       toggo.Context
		expected  bool
	}{
		{
			name: "equal strings",
			condition: toggo.Condition{
				Attribute: "country",
				Operator:  toggo.OperatorEqual,
				Value:     "US",
			},
			ctx:      toggo.Context{"country": "US"},
			expected: true,
		},
		{
			name: "not equal strings",
			condition: toggo.Condition{
				Attribute: "country",
				Operator:  toggo.OperatorEqual,
				Value:     "US",
			},
			ctx:      toggo.Context{"country": "CA"},
			expected: false,
		},
		{
			name: "equal numbers",
			condition: toggo.Condition{
				Attribute: "age",
				Operator:  toggo.OperatorEqual,
				Value:     25,
			},
			ctx:      toggo.Context{"age": 25},
			expected: true,
		},
		{
			name: "missing attribute",
			condition: toggo.Condition{
				Attribute: "country",
				Operator:  toggo.OperatorEqual,
				Value:     "US",
			},
			ctx:      toggo.Context{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.Evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStandardEvaluator_In(t *testing.T) {
	eval := NewStandard()

	tests := []struct {
		name      string
		condition toggo.Condition
		ctx       toggo.Context
		expected  bool
	}{
		{
			name: "value in list",
			condition: toggo.Condition{
				Attribute: "country",
				Operator:  toggo.OperatorIn,
				Value:     []interface{}{"US", "CA", "UK"},
			},
			ctx:      toggo.Context{"country": "US"},
			expected: true,
		},
		{
			name: "value not in list",
			condition: toggo.Condition{
				Attribute: "country",
				Operator:  toggo.OperatorIn,
				Value:     []interface{}{"US", "CA", "UK"},
			},
			ctx:      toggo.Context{"country": "DE"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.Evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStandardEvaluator_Comparison(t *testing.T) {
	eval := NewStandard()

	tests := []struct {
		name      string
		condition toggo.Condition
		ctx       toggo.Context
		expected  bool
	}{
		{
			name: "greater than - true",
			condition: toggo.Condition{
				Attribute: "age",
				Operator:  toggo.OperatorGreaterThan,
				Value:     18,
			},
			ctx:      toggo.Context{"age": 25},
			expected: true,
		},
		{
			name: "greater than - false",
			condition: toggo.Condition{
				Attribute: "age",
				Operator:  toggo.OperatorGreaterThan,
				Value:     30,
			},
			ctx:      toggo.Context{"age": 25},
			expected: false,
		},
		{
			name: "less than - true",
			condition: toggo.Condition{
				Attribute: "age",
				Operator:  toggo.OperatorLessThan,
				Value:     30,
			},
			ctx:      toggo.Context{"age": 25},
			expected: true,
		},
		{
			name: "greater than or equal - equal",
			condition: toggo.Condition{
				Attribute: "age",
				Operator:  toggo.OperatorGreaterThanOrEqual,
				Value:     25,
			},
			ctx:      toggo.Context{"age": 25},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.Evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStandardEvaluator_StringOperations(t *testing.T) {
	eval := NewStandard()

	tests := []struct {
		name      string
		condition toggo.Condition
		ctx       toggo.Context
		expected  bool
	}{
		{
			name: "contains - true",
			condition: toggo.Condition{
				Attribute: "email",
				Operator:  toggo.OperatorContains,
				Value:     "@gmail.com",
			},
			ctx:      toggo.Context{"email": "user@gmail.com"},
			expected: true,
		},
		{
			name: "starts with - true",
			condition: toggo.Condition{
				Attribute: "name",
				Operator:  toggo.OperatorStartsWith,
				Value:     "John",
			},
			ctx:      toggo.Context{"name": "John Doe"},
			expected: true,
		},
		{
			name: "ends with - false",
			condition: toggo.Condition{
				Attribute: "filename",
				Operator:  toggo.OperatorEndsWith,
				Value:     ".jpg",
			},
			ctx:      toggo.Context{"filename": "image.png"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.Evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStandardEvaluator_Negate(t *testing.T) {
	eval := NewStandard()

	condition := toggo.Condition{
		Attribute: "country",
		Operator:  toggo.OperatorEqual,
		Value:     "US",
		Negate:    true,
	}

	ctx := toggo.Context{"country": "CA"}
	result, err := eval.Evaluate(condition, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// CA != US, so result is false, but negated should be true
	if !result {
		t.Error("expected true after negation")
	}
}

func TestStandardEvaluator_EvaluateAll(t *testing.T) {
	eval := NewStandard()

	conditions := []toggo.Condition{
		{
			Attribute: "country",
			Operator:  toggo.OperatorIn,
			Value:     []interface{}{"US", "CA"},
		},
		{
			Attribute: "plan",
			Operator:  toggo.OperatorEqual,
			Value:     "premium",
		},
	}

	tests := []struct {
		name     string
		ctx      toggo.Context
		expected bool
	}{
		{
			name:     "all match",
			ctx:      toggo.Context{"country": "US", "plan": "premium"},
			expected: true,
		},
		{
			name:     "first fails",
			ctx:      toggo.Context{"country": "DE", "plan": "premium"},
			expected: false,
		},
		{
			name:     "second fails",
			ctx:      toggo.Context{"country": "US", "plan": "basic"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.EvaluateAll(conditions, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

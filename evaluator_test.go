package toggo

import (
	"testing"
)

func TestConditionEvaluator_Equal(t *testing.T) {
	eval := newConditionEvaluator()

	tests := []struct {
		name      string
		condition Condition
		ctx       Context
		expected  bool
	}{
		{
			name: "equal strings",
			condition: Condition{
				Attribute: "country",
				Operator:  OperatorEqual,
				Value:     "US",
			},
			ctx:      Context{"country": "US"},
			expected: true,
		},
		{
			name: "not equal strings",
			condition: Condition{
				Attribute: "country",
				Operator:  OperatorEqual,
				Value:     "US",
			},
			ctx:      Context{"country": "CA"},
			expected: false,
		},
		{
			name: "equal numbers",
			condition: Condition{
				Attribute: "age",
				Operator:  OperatorEqual,
				Value:     25,
			},
			ctx:      Context{"age": 25},
			expected: true,
		},
		{
			name: "missing attribute",
			condition: Condition{
				Attribute: "country",
				Operator:  OperatorEqual,
				Value:     "US",
			},
			ctx:      Context{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConditionEvaluator_In(t *testing.T) {
	eval := newConditionEvaluator()

	tests := []struct {
		name      string
		condition Condition
		ctx       Context
		expected  bool
	}{
		{
			name: "value in list",
			condition: Condition{
				Attribute: "country",
				Operator:  OperatorIn,
				Value:     []interface{}{"US", "CA", "UK"},
			},
			ctx:      Context{"country": "US"},
			expected: true,
		},
		{
			name: "value not in list",
			condition: Condition{
				Attribute: "country",
				Operator:  OperatorIn,
				Value:     []interface{}{"US", "CA", "UK"},
			},
			ctx:      Context{"country": "DE"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConditionEvaluator_Comparison(t *testing.T) {
	eval := newConditionEvaluator()

	tests := []struct {
		name      string
		condition Condition
		ctx       Context
		expected  bool
	}{
		{
			name: "greater than - true",
			condition: Condition{
				Attribute: "age",
				Operator:  OperatorGreaterThan,
				Value:     18,
			},
			ctx:      Context{"age": 25},
			expected: true,
		},
		{
			name: "greater than - false",
			condition: Condition{
				Attribute: "age",
				Operator:  OperatorGreaterThan,
				Value:     30,
			},
			ctx:      Context{"age": 25},
			expected: false,
		},
		{
			name: "less than - true",
			condition: Condition{
				Attribute: "age",
				Operator:  OperatorLessThan,
				Value:     30,
			},
			ctx:      Context{"age": 25},
			expected: true,
		},
		{
			name: "greater than or equal - equal",
			condition: Condition{
				Attribute: "age",
				Operator:  OperatorGreaterThanOrEqual,
				Value:     25,
			},
			ctx:      Context{"age": 25},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConditionEvaluator_StringOperations(t *testing.T) {
	eval := newConditionEvaluator()

	tests := []struct {
		name      string
		condition Condition
		ctx       Context
		expected  bool
	}{
		{
			name: "contains - true",
			condition: Condition{
				Attribute: "email",
				Operator:  OperatorContains,
				Value:     "@gmail.com",
			},
			ctx:      Context{"email": "user@gmail.com"},
			expected: true,
		},
		{
			name: "starts with - true",
			condition: Condition{
				Attribute: "name",
				Operator:  OperatorStartsWith,
				Value:     "John",
			},
			ctx:      Context{"name": "John Doe"},
			expected: true,
		},
		{
			name: "ends with - false",
			condition: Condition{
				Attribute: "filename",
				Operator:  OperatorEndsWith,
				Value:     ".jpg",
			},
			ctx:      Context{"filename": "image.png"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.evaluate(tt.condition, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConditionEvaluator_Negate(t *testing.T) {
	eval := newConditionEvaluator()

	condition := Condition{
		Attribute: "country",
		Operator:  OperatorEqual,
		Value:     "US",
		Negate:    true,
	}

	ctx := Context{"country": "CA"}
	result, err := eval.evaluate(condition, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// CA != US, so result is false, but negated should be true
	if !result {
		t.Error("expected true after negation")
	}
}

func TestConditionEvaluator_EvaluateAll(t *testing.T) {
	eval := newConditionEvaluator()

	conditions := []Condition{
		{
			Attribute: "country",
			Operator:  OperatorIn,
			Value:     []interface{}{"US", "CA"},
		},
		{
			Attribute: "plan",
			Operator:  OperatorEqual,
			Value:     "premium",
		},
	}

	tests := []struct {
		name     string
		ctx      Context
		expected bool
	}{
		{
			name:     "all match",
			ctx:      Context{"country": "US", "plan": "premium"},
			expected: true,
		},
		{
			name:     "first fails",
			ctx:      Context{"country": "DE", "plan": "premium"},
			expected: false,
		},
		{
			name:     "second fails",
			ctx:      Context{"country": "US", "plan": "basic"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eval.evaluateAll(conditions, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

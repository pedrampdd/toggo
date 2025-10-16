package main

import (
	"fmt"

	"github.com/pedrampdd/toggo"
)

func main() {
	store := toggo.NewStore()

	// Define a flag with conditions
	premiumFlag := &toggo.Flag{
		Name:    "premium_features",
		Enabled: true,
		Rollout: 100,
		Conditions: []toggo.Condition{
			{
				Attribute: "country",
				Operator:  toggo.OperatorIn,
				Value:     []interface{}{"US", "CA", "UK"},
			},
			{
				Attribute: "plan",
				Operator:  toggo.OperatorEqual,
				Value:     "premium",
			},
		},
	}

	store.AddFlag(premiumFlag)

	// Test different users
	testCases := []struct {
		name    string
		context toggo.Context
	}{
		{
			name: "Premium US user",
			context: toggo.Context{
				"user_id": "user_1",
				"country": "US",
				"plan":    "premium",
			},
		},
		{
			name: "Basic US user",
			context: toggo.Context{
				"user_id": "user_2",
				"country": "US",
				"plan":    "basic",
			},
		},
		{
			name: "Premium German user",
			context: toggo.Context{
				"user_id": "user_3",
				"country": "DE",
				"plan":    "premium",
			},
		},
		{
			name: "Premium UK user",
			context: toggo.Context{
				"user_id": "user_4",
				"country": "UK",
				"plan":    "premium",
			},
		},
	}

	fmt.Println("Testing conditional flag:")
	for _, tc := range testCases {
		enabled := store.IsEnabled("premium_features", tc.context)
		status := "✗ disabled"
		if enabled {
			status = "✓ enabled"
		}
		fmt.Printf("  %s: %s\n", tc.name, status)
	}

	// Example with numeric conditions
	ageFlag := &toggo.Flag{
		Name:    "age_restricted_content",
		Enabled: true,
		Rollout: 100,
		Conditions: []toggo.Condition{
			{
				Attribute: "age",
				Operator:  toggo.OperatorGreaterThanOrEqual,
				Value:     18,
			},
		},
	}

	store.AddFlag(ageFlag)

	fmt.Println("\nTesting age restriction:")
	ages := []int{15, 18, 25}
	for _, age := range ages {
		ctx := toggo.Context{
			"user_id": "user",
			"age":     age,
		}
		enabled := store.IsEnabled("age_restricted_content", ctx)
		status := "✗ disabled"
		if enabled {
			status = "✓ enabled"
		}
		fmt.Printf("  Age %d: %s\n", age, status)
	}
}

package main

import (
	"fmt"

	"github.com/pedram/toggo"
)

func main() {
	// Create a new feature flag store
	store := toggo.NewStore()

	// Define a feature flag with conditions
	premiumFeature := &toggo.Flag{
		Name:    "premium_dashboard",
		Enabled: true,
		Rollout: 100,
		Conditions: []toggo.Condition{
			{
				Attribute: "plan",
				Operator:  toggo.OperatorEqual,
				Value:     "premium",
			},
		},
	}

	store.AddFlag(premiumFeature)

	// Define a geographic-targeted feature
	geoFeature := &toggo.Flag{
		Name:    "new_checkout",
		Enabled: true,
		Rollout: 50,
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

	store.AddFlag(geoFeature)

	// Test different user scenarios
	testUsers := []struct {
		name    string
		context toggo.Context
	}{
		{
			name: "Premium US User",
			context: toggo.Context{
				"user_id": "user_1",
				"plan":    "premium",
				"country": "US",
			},
		},
		{
			name: "Basic US User",
			context: toggo.Context{
				"user_id": "user_2",
				"plan":    "basic",
				"country": "US",
			},
		},
		{
			name: "Premium DE User",
			context: toggo.Context{
				"user_id": "user_3",
				"plan":    "premium",
				"country": "DE",
			},
		},
	}

	fmt.Println("Feature Flag Evaluation:")
	fmt.Println("========================")

	for _, user := range testUsers {
		fmt.Printf("%s:\n", user.name)

		premiumEnabled := store.IsEnabled("premium_dashboard", user.context)
		fmt.Printf("  Premium Dashboard: %v\n", premiumEnabled)

		checkoutEnabled := store.IsEnabled("new_checkout", user.context)
		fmt.Printf("  New Checkout: %v\n", checkoutEnabled)

		fmt.Println()
	}

	// Example with numeric comparison
	ageGatedFeature := &toggo.Flag{
		Name:    "adult_content",
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

	store.AddFlag(ageGatedFeature)

	// Test age-gated feature
	fmt.Println("Age-Gated Feature:")
	fmt.Println("==================")

	ages := []int{15, 18, 25}
	for _, age := range ages {
		ctx := toggo.Context{
			"user_id": fmt.Sprintf("user_age_%d", age),
			"age":     age,
		}

		enabled := store.IsEnabled("adult_content", ctx)
		fmt.Printf("User (age %d): %v\n", age, enabled)
	}
}

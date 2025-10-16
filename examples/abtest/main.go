package main

import (
	"fmt"

	"github.com/pedrampdd/toggo"
)

func main() {
	// Create a new feature flag store
	store := toggo.NewStore()

	// Define an A/B test with multiple variants
	abTest := &toggo.Flag{
		Name:           "pricing_test",
		Enabled:        true,
		RolloutKey:     "user_id",
		DefaultVariant: "control",
		Variants: []toggo.Variant{
			{
				Name:   "control",
				Weight: 33, // 33% of users
			},
			{
				Name:   "price_low",
				Weight: 33, // 33% of users
			},
			{
				Name:   "price_high",
				Weight: 34, // 34% of users
			},
		},
	}

	if err := store.AddFlag(abTest); err != nil {
		fmt.Printf("Error adding flag: %v\n", err)
		return
	}

	// Simulate users and show which variant they get
	fmt.Println("A/B Test Results:")
	fmt.Println("=================")

	variantCounts := make(map[string]int)

	for i := 1; i <= 100; i++ {
		ctx := toggo.Context{
			"user_id": fmt.Sprintf("user_%d", i),
		}

		variant, enabled := store.GetVariant("pricing_test", ctx)
		if enabled {
			variantCounts[variant]++

			// Show first 10 assignments
			if i <= 10 {
				fmt.Printf("User %d -> %s\n", i, variant)
			}
		}
	}

	// Show distribution
	fmt.Println("\nVariant Distribution (100 users):")
	for variant, count := range variantCounts {
		fmt.Printf("%s: %d%%\n", variant, count)
	}

	// Example: Applying different pricing based on variant
	ctx := toggo.Context{"user_id": "user_42"}
	variant, _ := store.GetVariant("pricing_test", ctx)

	var price float64
	switch variant {
	case "control":
		price = 99.99
	case "price_low":
		price = 79.99
	case "price_high":
		price = 119.99
	default:
		price = 99.99
	}

	fmt.Printf("\nUser 42 sees variant '%s' with price: $%.2f\n", variant, price)
}

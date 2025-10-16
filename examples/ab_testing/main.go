package main

import (
	"fmt"

	"github.com/pedrampdd/toggo"
)

func main() {
	store := toggo.NewStore()

	// Define an A/B test with multiple variants
	pricingTest := &toggo.Flag{
		Name:           "pricing_experiment",
		Enabled:        true,
		RolloutKey:     "user_id",
		DefaultVariant: "control",
		Variants: []toggo.Variant{
			{
				Name:   "control",
				Weight: 34,
			},
			{
				Name:   "price_low",
				Weight: 33,
			},
			{
				Name:   "price_high",
				Weight: 33,
			},
		},
	}

	store.AddFlag(pricingTest)

	// Simulate 100 users and count variant distribution
	variantCounts := make(map[string]int)

	fmt.Println("Simulating A/B test with 100 users:")
	for i := 1; i <= 100; i++ {
		ctx := toggo.Context{
			"user_id": fmt.Sprintf("user_%d", i),
		}

		variant, enabled := store.GetVariant("pricing_experiment", ctx)
		if enabled {
			variantCounts[variant]++

			// Show first 10 assignments
			if i <= 10 {
				fmt.Printf("  User %d assigned to: %s\n", i, variant)
			}
		}
	}

	// Show distribution
	fmt.Println("\nVariant distribution:")
	for variant, count := range variantCounts {
		fmt.Printf("  %s: %d users (%d%%)\n", variant, count, count)
	}

	// Example: Using variant to show different prices
	ctx := toggo.Context{"user_id": "demo_user"}
	variant, _ := store.GetVariant("pricing_experiment", ctx)

	fmt.Printf("\nExample for demo_user (variant: %s):\n", variant)

	switch variant {
	case "control":
		fmt.Println("  Showing original price: $99/month")
	case "price_low":
		fmt.Println("  Showing discounted price: $79/month")
	case "price_high":
		fmt.Println("  Showing premium price: $129/month")
	default:
		fmt.Println("  Showing default price: $99/month")
	}
}

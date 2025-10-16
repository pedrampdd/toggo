package main

import (
	"fmt"

	"github.com/pedram/toggo"
)

func main() {
	// Create a new feature flag store
	store := toggo.NewStore()

	// Define a simple feature flag
	simpleFlag := &toggo.Flag{
		Name:    "new_ui",
		Enabled: true,
		Rollout: 100, // 100% rollout
	}

	// Add the flag to the store
	if err := store.AddFlag(simpleFlag); err != nil {
		panic(err)
	}

	// Create a context with user attributes
	ctx := toggo.Context{
		"user_id": "user_123",
		"country": "US",
	}

	// Check if the flag is enabled
	if store.IsEnabled("new_ui", ctx) {
		fmt.Println("✓ New UI is enabled for this user")
	} else {
		fmt.Println("✗ New UI is disabled for this user")
	}

	// Example with 50% rollout
	rolloutFlag := &toggo.Flag{
		Name:       "beta_feature",
		Enabled:    true,
		Rollout:    50,
		RolloutKey: "user_id",
	}

	store.AddFlag(rolloutFlag)

	// Test with different users
	fmt.Println("\nTesting 50% rollout with different users:")
	for i := 1; i <= 5; i++ {
		userCtx := toggo.Context{
			"user_id": fmt.Sprintf("user_%d", i),
		}

		enabled := store.IsEnabled("beta_feature", userCtx)
		status := "disabled"
		if enabled {
			status = "enabled"
		}

		fmt.Printf("  User %d: %s\n", i, status)
	}
}

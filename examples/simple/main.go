package main

import (
	"fmt"

	"github.com/pedrampdd/toggo"
)

func main() {
	// Create a new feature flag store
	store := toggo.NewStore()

	// Define a simple feature flag
	flag := &toggo.Flag{
		Name:    "dark_mode",
		Enabled: true,
		Rollout: 100, // 100% rollout
	}

	// Add the flag to the store
	if err := store.AddFlag(flag); err != nil {
		fmt.Printf("Error adding flag: %v\n", err)
		return
	}

	// Create a context with user information
	ctx := toggo.Context{
		"user_id": "user_12345",
	}

	// Check if the feature is enabled
	if store.IsEnabled("dark_mode", ctx) {
		fmt.Println("Dark mode is enabled!")
	} else {
		fmt.Println("Dark mode is disabled")
	}

	// Example with 50% rollout
	rolloutFlag := &toggo.Flag{
		Name:       "new_ui",
		Enabled:    true,
		Rollout:    50, // 50% of users will see this
		RolloutKey: "user_id",
	}

	store.AddFlag(rolloutFlag)

	// Test with multiple users
	fmt.Println("\nTesting 50% rollout:")
	for i := 1; i <= 10; i++ {
		userCtx := toggo.Context{
			"user_id": fmt.Sprintf("user_%d", i),
		}

		enabled := store.IsEnabled("new_ui", userCtx)
		fmt.Printf("User %d: new_ui = %v\n", i, enabled)
	}
}

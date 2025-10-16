package main

import (
	"fmt"
	"log"

	"github.com/pedram/toggo"
	"github.com/pedram/toggo/loader"
)

func main() {
	// Create a new store
	store := toggo.NewStore()

	// Load flags from YAML file
	yamlLoader := loader.NewYAMLFile("../../testdata/flags.yaml")
	if err := yamlLoader.LoadIntoStore(store); err != nil {
		log.Fatalf("Failed to load flags: %v", err)
	}

	fmt.Printf("Loaded %d flags from YAML\n\n", store.Size())

	// List all loaded flags
	fmt.Println("Available flags:")
	for _, name := range store.ListFlags() {
		flag, _ := store.GetFlag(name)
		status := "disabled"
		if flag.Enabled {
			status = "enabled"
		}
		fmt.Printf("  - %s (%s, %d%% rollout)\n", name, status, flag.Rollout)
	}

	// Test a loaded flag
	ctx := toggo.Context{
		"user_id":     "test_user",
		"country":     "US",
		"plan":        "premium",
		"beta_tester": true,
	}

	fmt.Println("\nTesting flags with context:", ctx)

	if store.IsEnabled("new_checkout", ctx) {
		fmt.Println("  ✓ new_checkout is enabled")
	} else {
		fmt.Println("  ✗ new_checkout is disabled")
	}

	if store.IsEnabled("dark_mode", ctx) {
		fmt.Println("  ✓ dark_mode is enabled")
	} else {
		fmt.Println("  ✗ dark_mode is disabled")
	}

	// Test A/B test variant
	variant, enabled := store.GetVariant("pricing_experiment", ctx)
	if enabled {
		fmt.Printf("  ✓ pricing_experiment variant: %s\n", variant)
	} else {
		fmt.Printf("  ✗ pricing_experiment is disabled (default: %s)\n", variant)
	}
}

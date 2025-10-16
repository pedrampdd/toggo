package main

import (
	"fmt"
	"strings"

	"github.com/pedrampdd/toggo"
	"github.com/pedrampdd/toggo/loader"
)

func main() {
	// Example 1: Load from JSON string
	jsonConfig := `{
		"flags": [
			{
				"name": "feature_a",
				"enabled": true,
				"rollout": 100
			},
			{
				"name": "feature_b",
				"enabled": true,
				"rollout": 50,
				"rollout_key": "user_id",
				"conditions": [
					{
						"attribute": "country",
						"operator": "in",
						"value": ["US", "CA"]
					}
				]
			}
		]
	}`

	store := toggo.NewStore()

	jsonLoader := loader.NewJSONReader(strings.NewReader(jsonConfig))
	if err := jsonLoader.LoadIntoStore(store); err != nil {
		fmt.Printf("Error loading JSON: %v\n", err)
		return
	}

	fmt.Println("Loaded flags from JSON:")
	for _, name := range store.ListFlags() {
		fmt.Printf("  - %s\n", name)
	}

	// Example 2: Load from YAML string
	yamlConfig := `
flags:
  - name: dark_mode
    enabled: true
    rollout: 100
  - name: beta_features
    enabled: true
    rollout: 25
    rollout_key: user_id
    conditions:
      - attribute: beta_tester
        operator: "=="
        value: true
  - name: pricing_experiment
    enabled: true
    rollout_key: user_id
    default_variant: control
    variants:
      - name: control
        weight: 50
      - name: variant_a
        weight: 50
`

	storeYAML := toggo.NewStore()

	yamlLoader := loader.NewYAMLReader(strings.NewReader(yamlConfig))
	if err := yamlLoader.LoadIntoStore(storeYAML); err != nil {
		fmt.Printf("Error loading YAML: %v\n", err)
		return
	}

	fmt.Println("\nLoaded flags from YAML:")
	for _, name := range storeYAML.ListFlags() {
		fmt.Printf("  - %s\n", name)
	}

	// Example 3: Test the loaded flags
	fmt.Println("\nTesting loaded flags:")

	ctx := toggo.Context{
		"user_id":     "user_123",
		"country":     "US",
		"beta_tester": true,
	}

	// Test simple flag
	if storeYAML.IsEnabled("dark_mode", ctx) {
		fmt.Println("✓ Dark mode is enabled")
	}

	// Test conditional flag
	if storeYAML.IsEnabled("beta_features", ctx) {
		fmt.Println("✓ Beta features are enabled for beta tester")
	}

	// Test A/B test variant
	variant, _ := storeYAML.GetVariant("pricing_experiment", ctx)
	fmt.Printf("✓ User assigned to variant: %s\n", variant)

	// Example 4: Loading from file (uncomment to use)
	// fileLoader := loader.NewYAMLFile("flags.yaml")
	// if err := fileLoader.LoadIntoStore(store); err != nil {
	//     fmt.Printf("Error loading from file: %v\n", err)
	//     return
	// }

	fmt.Println("\nTo load from a file, use:")
	fmt.Println(`  loader := loader.NewYAMLFile("flags.yaml")`)
	fmt.Println(`  loader.LoadIntoStore(store)`)
}

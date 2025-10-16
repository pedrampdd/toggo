// Package toggo provides a flexible and performant feature flag and A/B testing SDK for Go.
//
// Toggo enables you to manage feature rollouts, conduct A/B tests, and control feature access
// with fine-grained targeting conditions. It supports:
//   - Simple on/off feature flags
//   - Percentage-based rollouts with deterministic hashing
//   - Complex conditional targeting (country, plan, custom attributes)
//   - A/B testing with multiple variants
//   - Thread-safe operations
//   - Configuration loading from JSON/YAML files
//
// # Basic Usage
//
// Create a store and add a feature flag:
//
//	store := toggo.NewStore()
//	flag := &toggo.Flag{
//		Name:    "new_checkout",
//		Enabled: true,
//		Rollout: 50,
//	}
//	store.AddFlag(flag)
//
// Check if a feature is enabled:
//
//	ctx := toggo.Context{
//		"user_id": "12345",
//		"country": "US",
//	}
//	if store.IsEnabled("new_checkout", ctx) {
//		// Use new checkout flow
//	}
//
// # A/B Testing
//
// For A/B testing with variants:
//
//	flag := &toggo.Flag{
//		Name:           "pricing_test",
//		Enabled:        true,
//		DefaultVariant: "control",
//		Variants: []toggo.Variant{
//			{Name: "control", Weight: 50},
//			{Name: "variant_a", Weight: 50},
//		},
//	}
//	store.AddFlag(flag)
//
//	variant, _ := store.GetVariant("pricing_test", ctx)
//	switch variant {
//	case "control":
//		// Original pricing
//	case "variant_a":
//		// New pricing
//	}
//
// # Conditional Targeting
//
// Add conditions to target specific users:
//
//	flag := &toggo.Flag{
//		Name:    "premium_feature",
//		Enabled: true,
//		Rollout: 100,
//		Conditions: []toggo.Condition{
//			{
//				Attribute: "plan",
//				Operator:  toggo.OperatorEqual,
//				Value:     "premium",
//			},
//			{
//				Attribute: "country",
//				Operator:  toggo.OperatorIn,
//				Value:     []interface{}{"US", "CA", "UK"},
//			},
//		},
//	}
//
// # Loading from Configuration Files
//
// Load flags from JSON or YAML:
//
//	loader := loader.NewYAMLFile("flags.yaml")
//	loader.LoadIntoStore(store)
//
// # Operators
//
// Toggo supports various comparison operators:
//   - == (equal)
//   - != (not equal)
//   - in (value in list)
//   - not_in (value not in list)
//   - > (greater than)
//   - >= (greater than or equal)
//   - < (less than)
//   - <= (less than or equal)
//   - contains (string contains)
//   - starts_with (string starts with)
//   - ends_with (string ends with)
//   - regex (regular expression match)
package toggo

const (
	// Version is the current version of the toggo SDK
	Version = "1.0.0"
)

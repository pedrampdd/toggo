package main

import (
	"fmt"
	"time"

	"github.com/pedrampdd/toggo"
)

func main() {
	fmt.Println("=== Switchback Testing Examples ===")
	fmt.Println()

	// Example 1: Driver Rebate Switchback Test
	driverRebateExample()

	fmt.Println("\n" + separator(60) + "\n")

	// Example 2: UI Feature Switchback Test
	uiFeatureExample()

	fmt.Println("\n" + separator(60) + "\n")

	// Example 3: Pricing Switchback with Daily Swap
	pricingWithDailySwapExample()
}

// Example 1: Testing different rebate strategies for drivers
func driverRebateExample() {
	fmt.Println("Example 1: Driver Rebate Switchback Test")
	fmt.Println("Testing two rebate strategies that switch every 30 minutes")
	fmt.Println()

	// Configure switchback with 30-minute intervals
	store := toggo.NewStore(
		toggo.WithSwitchback(
			toggo.WithStartTime(time.Now().Add(-24*time.Hour)),
			toggo.WithIntervalMinutes(30),
			toggo.WithDailySwap(true),
		),
	)

	// Define the rebate experiment flag with two rebate types
	rebateFlag := &toggo.Flag{
		Name:           "driver_rebate_experiment",
		Enabled:        true,
		DefaultVariant: "standard_rebate",
		Variants: []toggo.Variant{
			{Name: "standard_rebate"}, // 10% cashback
			{Name: "premium_rebate"},  // 15% cashback
		},
	}

	store.AddFlag(rebateFlag)

	// Simulate checking the rebate for a driver
	driverCtx := toggo.Context{
		"driver_id": "DRV-12345",
		"city":      "san_francisco",
	}

	// Get current rebate variant
	rebateType, enabled := store.GetVariant("driver_rebate_experiment", driverCtx)

	if enabled {
		fmt.Printf("Driver DRV-12345 should receive: %s\n", rebateType)

		// Apply the appropriate rebate logic
		switch rebateType {
		case "standard_rebate":
			fmt.Println("  → Applying 10%% cashback on completed rides")
		case "premium_rebate":
			fmt.Println("  → Applying 15%% cashback on completed rides")
		}
	}

	// Show switchback timing info
	if info := toggo.GetSwitchbackInfo(store); info != nil {
		fmt.Printf("\nSwitchback Status:\n")
		fmt.Printf("  Current interval: %d\n", info.CurrentInterval)
		fmt.Printf("  Next switch in: %v\n", info.TimeUntilSwitch.Round(time.Second))
		fmt.Printf("  Interval duration: %v\n", info.IntervalDuration)
	}

	fmt.Println("\nNote: ALL drivers see the same rebate type at the same time.")
	fmt.Println("This allows comparing performance metrics between time periods.")
}

// Example 2: Testing UI features across all users
func uiFeatureExample() {
	fmt.Println("Example 2: UI Feature Switchback Test")
	fmt.Println("Testing different checkout flows that switch every 15 minutes")
	fmt.Println()

	store := toggo.NewStore(
		toggo.WithSwitchback(
			toggo.WithIntervalMinutes(15),
		),
	)

	checkoutFlag := &toggo.Flag{
		Name:           "checkout_flow_test",
		Enabled:        true,
		DefaultVariant: "classic_checkout",
		Variants: []toggo.Variant{
			{Name: "classic_checkout", Weight: 33},
			{Name: "express_checkout", Weight: 33},
			{Name: "onepage_checkout", Weight: 34},
		},
	}

	store.AddFlag(checkoutFlag)

	// Check what flow users should see
	userCtx := toggo.Context{"user_id": "user_789"}
	flow, _ := store.GetVariant("checkout_flow_test", userCtx)

	fmt.Printf("Current checkout flow for all users: %s\n", flow)

	// Show how the flow cycles through variants
	fmt.Println("\nSwitchback Schedule (15-minute intervals):")
	fmt.Println("  Interval 0: classic_checkout")
	fmt.Println("  Interval 1: express_checkout")
	fmt.Println("  Interval 2: onepage_checkout")
	fmt.Println("  Interval 3: classic_checkout (repeats)")
	fmt.Println("  ...")
}

// Example 3: Pricing test with daily pattern swap
func pricingWithDailySwapExample() {
	fmt.Println("Example 3: Pricing Test with Daily Swap")
	fmt.Println("Switch between pricing tiers every 2 hours")
	fmt.Println("Pattern reverses each day to control for time-of-day effects")
	fmt.Println()

	// Start at beginning of today for clean day boundaries
	startOfToday := time.Now().Truncate(24 * time.Hour)

	store := toggo.NewStore(
		toggo.WithSwitchback(
			toggo.WithIntervalMinutes(120), // 2-hour intervals
			toggo.WithStartTime(startOfToday),
			toggo.WithDailySwap(true), // Enable daily pattern reversal
		),
	)

	pricingFlag := &toggo.Flag{
		Name:           "pricing_tier_test",
		Enabled:        true,
		DefaultVariant: "tier_standard",
		Variants: []toggo.Variant{
			{Name: "tier_economy", Weight: 50},
			{Name: "tier_standard", Weight: 50},
		},
	}

	store.AddFlag(pricingFlag)

	// Show the pattern
	fmt.Println("Day 0 Pattern:")
	fmt.Println("  00:00-02:00 → tier_economy")
	fmt.Println("  02:00-04:00 → tier_standard")
	fmt.Println("  04:00-06:00 → tier_economy")
	fmt.Println("  06:00-08:00 → tier_standard")
	fmt.Println("  ... (pattern continues)")

	fmt.Println("\nDay 1 Pattern (SWAPPED):")
	fmt.Println("  00:00-02:00 → tier_standard")
	fmt.Println("  02:00-04:00 → tier_economy")
	fmt.Println("  04:00-06:00 → tier_standard")
	fmt.Println("  06:00-08:00 → tier_economy")
	fmt.Println("  ... (pattern continues)")

	fmt.Println("\nDay 2 Pattern (back to Day 0):")
	fmt.Println("  00:00-02:00 → tier_economy")
	fmt.Println("  ... (same as Day 0)")

	// Get current pricing tier
	ctx := toggo.Context{"user_id": "customer_456"}
	tier, _ := store.GetVariant("pricing_tier_test", ctx)

	if info := toggo.GetSwitchbackInfo(store); info != nil {
		fmt.Printf("\nCurrent Status:\n")
		fmt.Printf("  Active pricing tier: %s\n", tier)
		fmt.Printf("  Day: %d\n", info.CurrentDay)
		fmt.Printf("  Interval: %d\n", info.CurrentInterval)
	}

	fmt.Println("\nBenefit: Daily swap controls for time-of-day effects")
	fmt.Println("(e.g., morning vs evening traffic patterns)")
}

func separator(length int) string {
	s := ""
	for i := 0; i < length; i++ {
		s += "-"
	}
	return s
}

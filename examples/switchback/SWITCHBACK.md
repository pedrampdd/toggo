# Switchback Testing in Toggo

## What is Switchback Testing?

Switchback testing is a time-based experimentation method where **all users experience the same variant at the same time**, and the variant switches at regular intervals. Unlike traditional A/B testing where users are randomly assigned to groups, switchback testing rotates everyone through different treatments based on time.

## When to Use Switchback Testing

Switchback testing is ideal for:

1. **Marketplace Experiments**: Testing driver incentives, pricing strategies, or supply/demand interventions where individual randomization isn't feasible
2. **System-Wide Changes**: Testing infrastructure changes, algorithm updates, or platform-wide features
3. **Network Effects**: Scenarios where user experience depends on what others are experiencing
4. **Controlling Time Effects**: Eliminating time-of-day bias by alternating patterns daily

## Key Differences from Standard A/B Testing

| Aspect | Standard A/B Testing | Switchback Testing |
|--------|---------------------|-------------------|
| Assignment | Per user (hashed) | All users at once |
| Consistency | User stays in same group | All users rotate together |
| Use Case | Individual features | System-wide changes |
| Time Factor | Controlled via random assignment | Explicitly alternated |

## Basic Usage

### Simple Switchback

```go
import "github.com/pedrampdd/toggo"

// Create store with switchback strategy
store := toggo.NewStore(
    toggo.WithSwitchback(
        toggo.WithIntervalMinutes(30), // Switch every 30 minutes
    ),
)

// Define your variants
flag := &toggo.Flag{
    Name:           "pricing_strategy",
    Enabled:        true,
    DefaultVariant: "standard",
    Variants: []toggo.Variant{
        {Name: "standard", Weight: 50},
        {Name: "dynamic", Weight: 50},
    },
}

store.AddFlag(flag)

// Everyone gets the same variant at the same time
ctx := toggo.Context{"user_id": "any_user"}
variant, _ := store.GetVariant("pricing_strategy", ctx)

// variant will be "standard" or "dynamic" based on current time interval
```

### With Daily Pattern Swap

To control for time-of-day effects, enable daily swap:

```go
store := toggo.NewStore(
    toggo.WithSwitchback(
        toggo.WithIntervalMinutes(30),
        toggo.WithDailySwap(true), // Reverse pattern each day
    ),
)
```

**Day 0 Schedule:**
- 00:00-00:30 → Variant A
- 00:30-01:00 → Variant B
- 01:00-01:30 → Variant A
- ... continues

**Day 1 Schedule (reversed):**
- 00:00-00:30 → Variant B
- 00:30-01:00 → Variant A
- 01:00-01:30 → Variant B
- ... continues

This ensures that each variant gets equal exposure to different times of day.

### Custom Start Time

Set a specific start time for clean boundaries:

```go
import "time"

startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

store := toggo.NewStore(
    toggo.WithSwitchback(
        toggo.WithIntervalMinutes(60),
        toggo.WithStartTime(startTime),
    ),
)
```

## Getting Timing Information

```go
// Get detailed switchback timing info
if info := toggo.GetSwitchbackInfo(store); info != nil {
    fmt.Printf("Current interval: %d\n", info.CurrentInterval)
    fmt.Printf("Current day: %d\n", info.CurrentDay)
    fmt.Printf("Time until next switch: %v\n", info.TimeUntilSwitch)
    fmt.Printf("Interval duration: %v\n", info.IntervalDuration)
}
```

## Real-World Examples

### Example 1: Driver Rebate Program

You want to test two rebate strategies for drivers:

```go
store := toggo.NewStore(
    toggo.WithSwitchback(
        toggo.WithIntervalMinutes(30),
        toggo.WithDailySwap(true),
    ),
)

rebateFlag := &toggo.Flag{
    Name:           "driver_rebate_test",
    Enabled:        true,
    DefaultVariant: "standard_rebate",
    Variants: []toggo.Variant{
        {Name: "standard_rebate", Weight: 50},  // 10% cashback
        {Name: "premium_rebate", Weight: 50},   // 15% cashback
    },
}

store.AddFlag(rebateFlag)

// In your driver payment logic
ctx := toggo.Context{"driver_id": driverID}
rebateType, _ := store.GetVariant("driver_rebate_test", ctx)

switch rebateType {
case "standard_rebate":
    cashback = rideEarnings * 0.10
case "premium_rebate":
    cashback = rideEarnings * 0.15
}

applyRebate(driverID, cashback)
```

**Why switchback?** All drivers see the same rebate at the same time, allowing you to measure aggregate marketplace effects (supply, acceptance rates) without confounding variables.

### Example 2: Dynamic Pricing Algorithm

Testing a new pricing algorithm:

```go
store := toggo.NewStore(
    toggo.WithSwitchback(
        toggo.WithIntervalMinutes(120), // 2-hour intervals
        toggo.WithDailySwap(true),
    ),
)

pricingFlag := &toggo.Flag{
    Name:           "pricing_algorithm",
    Enabled:        true,
    DefaultVariant: "legacy_pricing",
    Variants: []toggo.Variant{
        {Name: "legacy_pricing", Weight: 50},
        {Name: "ml_based_pricing", Weight: 50},
    },
}

store.AddFlag(pricingFlag)

// When calculating ride price
ctx := toggo.Context{"ride_id": rideID}
algorithm, _ := store.GetVariant("pricing_algorithm", ctx)

var price float64
switch algorithm {
case "legacy_pricing":
    price = calculateLegacyPrice(distance, duration)
case "ml_based_pricing":
    price = calculateMLPrice(distance, duration, demand, supply)
}
```

### Example 3: Three-Way Test

Testing three different notification strategies:

```go
store := toggo.NewStore(
    toggo.WithSwitchback(
        toggo.WithIntervalMinutes(20), // 20-minute intervals
    ),
)

notificationFlag := &toggo.Flag{
    Name:           "notification_strategy",
    Enabled:        true,
    DefaultVariant: "standard",
    Variants: []toggo.Variant{
        {Name: "standard", Weight: 33},   // Original notifications
        {Name: "frequent", Weight: 33},    // More frequent
        {Name: "smart_timing", Weight: 34}, // ML-optimized timing
    },
}

store.AddFlag(notificationFlag)

// Variants rotate: standard → frequent → smart_timing → standard → ...
ctx := toggo.Context{"user_id": userID}
strategy, _ := store.GetVariant("notification_strategy", ctx)
```

## Best Practices

### 1. Choose Appropriate Interval Length

- **Short intervals (15-30 min)**: Good for high-traffic scenarios, gets results faster
- **Long intervals (2-4 hours)**: Better for capturing longer-term effects
- Consider your metric collection frequency

### 2. Enable Daily Swap for Time-Sensitive Metrics

If your metrics vary by time of day (e.g., morning vs. evening traffic), enable daily swap:

```go
toggo.WithDailySwap(true)
```

### 3. Log Timing Information

Include interval/day information in your metrics for analysis:

```go
info := toggo.GetSwitchbackInfo(store)
logMetric("conversion", value, map[string]interface{}{
    "interval": info.CurrentInterval,
    "day":      info.CurrentDay,
    "variant":  variant,
})
```

### 4. Account for Carryover Effects

Some interventions have carryover effects. Consider:
- Adding "washout" periods between switches
- Using longer intervals
- Analyzing lag effects in your data

### 5. Monitor Switches

Set up monitoring around switch times:

```go
if info := toggo.GetSwitchbackInfo(store); info != nil {
    if info.TimeUntilSwitch < 5*time.Minute {
        // Alert: approaching switch time
        // Useful for monitoring metric transitions
    }
}
```

## Statistical Considerations

### Sample Size

Switchback tests require careful sample size calculation:
- Each interval is one observation
- Need sufficient intervals for statistical power
- Consider using longer test periods

### Analysis

When analyzing results:
1. Group data by interval
2. Compare metrics across variants
3. Account for time trends
4. Use paired comparisons (especially with daily swap)

### Variance

Switchback tests may have higher variance than traditional A/B tests because:
- Fewer independent units (intervals vs. users)
- Time-based correlation between observations

Compensate by:
- Running tests longer
- Using daily swap to balance time effects
- Collecting more intervals

## Configuration Options

### WithIntervalMinutes(minutes int)

Sets how long each variant runs before switching.

```go
toggo.WithIntervalMinutes(30) // 30-minute intervals
```

### WithStartTime(t time.Time)

Sets the reference start time for calculating intervals and days.

```go
startTime := time.Now().Truncate(24 * time.Hour) // Start of today
toggo.WithStartTime(startTime)
```

### WithDailySwap(enabled bool)

Enables daily pattern reversal to control for time-of-day effects.

```go
toggo.WithDailySwap(true)
```

## Advanced: Testing Setup

For testing your switchback implementation:

```go
func TestYourSwitchbackLogic(t *testing.T) {
    store := toggo.NewStore(
        toggo.WithSwitchback(
            toggo.WithIntervalMinutes(30),
            toggo.WithStartTime(fixedStartTime),
        ),
    )
    
    // Access the strategy to control time in tests
    strategy := store.GetRolloutStrategy().(*toggo.SwitchbackRolloutStrategy)
    strategy.SetTimeFunc(func() time.Time {
        return testTime // Your controlled time
    })
    
    // Now test your logic with controlled time
}
```

## Troubleshooting

### Problem: Variants switching too frequently
**Solution**: Increase interval duration with `WithIntervalMinutes()`

### Problem: Need to align with business hours
**Solution**: Set `WithStartTime()` to align intervals with your timezone

### Problem: Time-of-day effects in results
**Solution**: Enable `WithDailySwap(true)` to balance patterns

### Problem: Need to know when next switch happens
**Solution**: Use `GetSwitchbackInfo()` to get timing details

## Further Reading

- [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture details
- [examples/switchback/main.go](examples/switchback/main.go) - Complete working examples
- [switchback_test.go](switchback_test.go) - Test cases showing various scenarios


# Switchback Testing Feature - Implementation Summary

## Overview

I've implemented a comprehensive **Switchback Testing** feature for your Toggo feature flag library. This enables time-based A/B testing where all users experience the same variant at the same time, with automatic switching at regular intervals.

## What Was Implemented

### 1. Core Implementation (`switchback.go`)

**SwitchbackRolloutStrategy** - A new rollout strategy that implements time-based variant switching:

- **Time-based intervals**: Switch between variants at configurable intervals (e.g., every 30 minutes)
- **Daily pattern swap**: Optional reversal of variant order on alternating days to control for time-of-day effects
- **Multiple variants support**: Works with 2+ variants, cycling through them in order
- **Timing information**: Get current interval, day number, and time until next switch

**Configuration Options**:
```go
toggo.WithSwitchback(
    toggo.WithIntervalMinutes(30),      // Switch every 30 minutes
    toggo.WithStartTime(time.Now()),    // Reference start time
    toggo.WithDailySwap(true),          // Enable daily pattern reversal
)
```

### 2. Integration with Existing System

The switchback strategy integrates seamlessly with your existing Flag/Variant system:

```go
// Works with existing Flag structure
flag := &toggo.Flag{
    Name:           "driver_rebate",
    Enabled:        true,
    DefaultVariant: "standard",
    Variants: []toggo.Variant{
        {Name: "standard_rebate", Weight: 50},
        {Name: "premium_rebate", Weight: 50},
    },
}

// Use existing GetVariant API
variant, enabled := store.GetVariant("driver_rebate", ctx)
```

### 3. Helper Functions

- `GetSwitchbackInfo(store)` - Get timing information from a store
- `GetInfo()` - Get detailed timing data (interval, day, time until switch)
- `GetCurrentInterval()` - Which time interval we're in
- `GetCurrentDay()` - Which day since start
- `GetTimeUntilNextSwitch()` - Time remaining in current interval

### 4. Comprehensive Tests (`switchback_test.go`)

- 11 test functions covering all scenarios
- Tests for 2-variant and 3-variant cases
- Daily swap behavior verification
- Integration tests with full Store
- Custom interval durations
- All tests passing ✅

### 5. Real-World Examples (`examples/switchback/main.go`)

Three complete examples demonstrating:

1. **Driver Rebate Test**: Your use case - switching between rebate strategies
2. **UI Feature Test**: Testing different checkout flows
3. **Pricing with Daily Swap**: Demonstrating time-of-day effect control

### 6. Documentation

- **README.md**: Updated with switchback section and examples
- **SWITCHBACK.md**: Comprehensive guide with:
  - What switchback testing is
  - When to use it
  - Real-world examples
  - Best practices
  - Statistical considerations
  - Troubleshooting guide

## How It Works

### Basic Concept

```
Traditional A/B Test:          Switchback Test:
User A → Variant 1 (always)   Time 0-30min  → All users see Variant 1
User B → Variant 2 (always)   Time 30-60min → All users see Variant 2
User C → Variant 1 (always)   Time 60-90min → All users see Variant 1
```

### Daily Swap (Optional)

Controls for time-of-day effects:

```
Day 0:                         Day 1 (swapped):
00:00-00:30 → Variant A       00:00-00:30 → Variant B
00:30-01:00 → Variant B       00:30-01:00 → Variant A
01:00-01:30 → Variant A       01:00-01:30 → Variant B
```

## Your Driver Rebate Use Case

Here's exactly how you'd use it for your rebate testing:

```go
package main

import (
    "github.com/pedrampdd/toggo"
    "time"
)

func main() {
    // Create store with 30-minute switchback intervals
    store := toggo.NewStore(
        toggo.WithSwitchback(
            toggo.WithIntervalMinutes(30),
            toggo.WithDailySwap(true), // Balance time-of-day effects
        ),
    )

    // Define your rebate variants
    rebateFlag := &toggo.Flag{
        Name:           "driver_rebate_experiment",
        Enabled:        true,
        DefaultVariant: "standard_rebate",
        Variants: []toggo.Variant{
            {Name: "standard_rebate", Weight: 50},
            {Name: "premium_rebate", Weight: 50},
        },
    }

    store.AddFlag(rebateFlag)

    // In your driver payment processing
    processDriverPayment := func(driverID string, rideEarnings float64) {
        ctx := toggo.Context{"driver_id": driverID}
        
        // Get current rebate type (same for all drivers at this time)
        rebateType, _ := store.GetVariant("driver_rebate_experiment", ctx)
        
        var rebateAmount float64
        switch rebateType {
        case "standard_rebate":
            rebateAmount = rideEarnings * 0.10 // 10% cashback
        case "premium_rebate":
            rebateAmount = rideEarnings * 0.15 // 15% cashback
        }
        
        applyRebateToDriver(driverID, rebateAmount)
        
        // Log for analysis (include timing info)
        if info := toggo.GetSwitchbackInfo(store); info != nil {
            logMetric("rebate_applied", map[string]interface{}{
                "driver_id":        driverID,
                "rebate_type":      rebateType,
                "rebate_amount":    rebateAmount,
                "interval":         info.CurrentInterval,
                "day":              info.CurrentDay,
            })
        }
    }
}
```

## Key Benefits for Your Use Case

1. **All drivers get same rebate at same time** - Eliminates driver-to-driver comparison issues
2. **Measure marketplace effects** - See how different rebates affect supply, acceptance rates, etc.
3. **Control for time effects** - Daily swap ensures each rebate gets equal exposure to peak/off-peak times
4. **Easy to implement** - Works with your existing Flag system
5. **Get timing info** - Know which group is active and when it switches

## Files Added/Modified

### New Files:
- `switchback.go` - Core implementation
- `switchback_test.go` - Comprehensive tests
- `examples/switchback/main.go` - Working examples
- `SWITCHBACK.md` - Complete documentation guide
- `FEATURE_SUMMARY.md` - This file

### Modified Files:
- `store.go` - Added `GetRolloutStrategy()` method for advanced usage
- `README.md` - Added switchback feature documentation

## Testing

All tests pass:
```bash
$ go test -v ./... -run TestSwitchback
# 11 tests, all passing
PASS
```

Example runs successfully:
```bash
$ cd examples/switchback && go run main.go
# Shows all three examples working correctly
```

## Generalization

The implementation is fully generalized and supports:

✅ **Any number of variants** (2, 3, 4, etc.)
✅ **Any interval duration** (minutes, hours)
✅ **Any use case** (rebates, pricing, features, algorithms)
✅ **Custom variant names** (no default naming like "segment_a")
✅ **Integration with existing conditions** (can still use Flag conditions)
✅ **Time zone control** (via WithStartTime)
✅ **Testing support** (injectable time provider)

## Next Steps / Usage

1. **Import the library**: The feature is ready to use
2. **Configure your experiment**: Choose interval duration and whether to use daily swap
3. **Define your variants**: Use meaningful names for your specific use case
4. **Integrate into your code**: Use `GetVariant()` to get current active variant
5. **Log with timing info**: Include interval/day numbers in your metrics
6. **Analyze results**: Compare metrics across intervals/variants

## Example Metrics Analysis

When analyzing your results, group by interval:

```sql
-- Example SQL for analyzing rebate experiment
SELECT 
    rebate_type,
    AVG(acceptance_rate) as avg_acceptance,
    AVG(rides_per_hour) as avg_supply,
    COUNT(*) as num_intervals
FROM metrics
WHERE experiment = 'driver_rebate_experiment'
GROUP BY rebate_type
```

## Questions?

Check the comprehensive documentation in `SWITCHBACK.md` for:
- Detailed usage examples
- Best practices
- Statistical considerations
- Troubleshooting guide
- Advanced configuration

---

**The feature is production-ready and fully tested!** 🚀


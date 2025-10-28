package toggo

import (
	"fmt"
	"testing"
	"time"
)

func TestSwitchbackRolloutStrategy_GetCurrentInterval(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		intervalMinutes  int
		currentTime      time.Time
		expectedInterval int
	}{
		{
			name:             "start of first interval",
			intervalMinutes:  30,
			currentTime:      startTime,
			expectedInterval: 0,
		},
		{
			name:             "middle of first interval",
			intervalMinutes:  30,
			currentTime:      startTime.Add(15 * time.Minute),
			expectedInterval: 0,
		},
		{
			name:             "start of second interval",
			intervalMinutes:  30,
			currentTime:      startTime.Add(30 * time.Minute),
			expectedInterval: 1,
		},
		{
			name:             "middle of third interval",
			intervalMinutes:  30,
			currentTime:      startTime.Add(75 * time.Minute),
			expectedInterval: 2,
		},
		{
			name:             "2 hour intervals",
			intervalMinutes:  120,
			currentTime:      startTime.Add(5 * time.Hour),
			expectedInterval: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := NewSwitchbackRolloutStrategy(
				WithIntervalMinutes(tt.intervalMinutes),
				WithStartTime(startTime),
			)
			strategy.timeProvider = func() time.Time { return tt.currentTime }

			interval := strategy.GetCurrentInterval()
			if interval != tt.expectedInterval {
				t.Errorf("GetCurrentInterval() = %v, want %v", interval, tt.expectedInterval)
			}
		})
	}
}

func TestSwitchbackRolloutStrategy_GetCurrentDay(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		currentTime time.Time
		expectedDay int
	}{
		{
			name:        "day 0",
			currentTime: startTime.Add(12 * time.Hour),
			expectedDay: 0,
		},
		{
			name:        "day 1",
			currentTime: startTime.Add(25 * time.Hour),
			expectedDay: 1,
		},
		{
			name:        "day 2",
			currentTime: startTime.Add(48 * time.Hour),
			expectedDay: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := NewSwitchbackRolloutStrategy(WithStartTime(startTime))
			strategy.timeProvider = func() time.Time { return tt.currentTime }

			day := strategy.GetCurrentDay()
			if day != tt.expectedDay {
				t.Errorf("GetCurrentDay() = %v, want %v", day, tt.expectedDay)
			}
		})
	}
}

func TestSwitchbackRolloutStrategy_GetVariant(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	flag := &Flag{
		Name:           "test_flag",
		Enabled:        true,
		DefaultVariant: "default",
		Variants: []Variant{
			{Name: "variant_a", Weight: 50},
			{Name: "variant_b", Weight: 50},
		},
	}

	tests := []struct {
		name            string
		currentTime     time.Time
		swapDaily       bool
		expectedVariant string
	}{
		{
			name:            "interval 0 shows first variant",
			currentTime:     startTime,
			swapDaily:       false,
			expectedVariant: "variant_a",
		},
		{
			name:            "interval 1 shows second variant",
			currentTime:     startTime.Add(30 * time.Minute),
			swapDaily:       false,
			expectedVariant: "variant_b",
		},
		{
			name:            "interval 2 cycles back to first",
			currentTime:     startTime.Add(60 * time.Minute),
			swapDaily:       false,
			expectedVariant: "variant_a",
		},
		{
			name:            "interval 3 shows second variant",
			currentTime:     startTime.Add(90 * time.Minute),
			swapDaily:       false,
			expectedVariant: "variant_b",
		},
		{
			name:            "day 1 interval 0 with swap shows second variant",
			currentTime:     startTime.Add(24 * time.Hour),
			swapDaily:       true,
			expectedVariant: "variant_b",
		},
		{
			name:            "day 1 interval 1 with swap shows first variant",
			currentTime:     startTime.Add(24*time.Hour + 30*time.Minute),
			swapDaily:       true,
			expectedVariant: "variant_a",
		},
		{
			name:            "day 2 interval 0 with swap back to first variant",
			currentTime:     startTime.Add(48 * time.Hour),
			swapDaily:       true,
			expectedVariant: "variant_a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := NewSwitchbackRolloutStrategy(
				WithIntervalMinutes(30),
				WithStartTime(startTime),
				WithDailySwap(tt.swapDaily),
			)
			strategy.timeProvider = func() time.Time { return tt.currentTime }

			ctx := Context{"user_id": "test_user"}
			variant, err := strategy.GetVariant(flag, ctx)
			if err != nil {
				t.Errorf("GetVariant() error = %v", err)
			}
			if variant != tt.expectedVariant {
				t.Errorf("GetVariant() = %v, want %v", variant, tt.expectedVariant)
			}
		})
	}
}

func TestSwitchbackRolloutStrategy_ThreeVariants(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	flag := &Flag{
		Name:           "test_flag",
		Enabled:        true,
		DefaultVariant: "default",
		Variants: []Variant{
			{Name: "variant_a", Weight: 33},
			{Name: "variant_b", Weight: 33},
			{Name: "variant_c", Weight: 34},
		},
	}

	tests := []struct {
		interval        int
		expectedVariant string
	}{
		{0, "variant_a"},
		{1, "variant_b"},
		{2, "variant_c"},
		{3, "variant_a"}, // cycles back
		{4, "variant_b"},
		{5, "variant_c"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("interval_%d", tt.interval), func(t *testing.T) {
			strategy := NewSwitchbackRolloutStrategy(
				WithIntervalMinutes(30),
				WithStartTime(startTime),
			)
			currentTime := startTime.Add(time.Duration(tt.interval*30) * time.Minute)
			strategy.timeProvider = func() time.Time { return currentTime }

			ctx := Context{"user_id": "test_user"}
			variant, err := strategy.GetVariant(flag, ctx)
			if err != nil {
				t.Errorf("GetVariant() error = %v", err)
			}
			if variant != tt.expectedVariant {
				t.Errorf("Interval %d: GetVariant() = %v, want %v", tt.interval, variant, tt.expectedVariant)
			}
		})
	}
}

func TestSwitchbackRolloutStrategy_ShouldRollout(t *testing.T) {
	strategy := NewSwitchbackRolloutStrategy()

	flag := &Flag{Name: "test", Enabled: true}
	ctx := Context{"user_id": "123"}

	shouldRollout, err := strategy.ShouldRollout(flag, ctx)
	if err != nil {
		t.Errorf("ShouldRollout() error = %v", err)
	}
	if !shouldRollout {
		t.Error("ShouldRollout() should always return true for switchback")
	}
}

func TestSwitchbackRolloutStrategy_GetTimeUntilNextSwitch(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentTime := startTime.Add(25 * time.Minute) // 25 minutes into 30-minute interval

	strategy := NewSwitchbackRolloutStrategy(
		WithIntervalMinutes(30),
		WithStartTime(startTime),
	)
	strategy.timeProvider = func() time.Time { return currentTime }

	timeUntilSwitch := strategy.GetTimeUntilNextSwitch()
	expected := 5 * time.Minute

	if timeUntilSwitch != expected {
		t.Errorf("GetTimeUntilNextSwitch() = %v, want %v", timeUntilSwitch, expected)
	}
}

func TestSwitchbackRolloutStrategy_GetInfo(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currentTime := startTime.Add(2*time.Hour + 45*time.Minute)

	strategy := NewSwitchbackRolloutStrategy(
		WithIntervalMinutes(30),
		WithStartTime(startTime),
	)
	strategy.timeProvider = func() time.Time { return currentTime }

	info := strategy.GetInfo()

	// After 2h 45m with 30-minute intervals = 5 complete intervals + 15 minutes
	expectedInterval := 5
	if info.CurrentInterval != expectedInterval {
		t.Errorf("Info.CurrentInterval = %v, want %v", info.CurrentInterval, expectedInterval)
	}

	expectedDay := 0
	if info.CurrentDay != expectedDay {
		t.Errorf("Info.CurrentDay = %v, want %v", info.CurrentDay, expectedDay)
	}

	expectedTimeUntilSwitch := 15 * time.Minute
	if info.TimeUntilSwitch != expectedTimeUntilSwitch {
		t.Errorf("Info.TimeUntilSwitch = %v, want %v", info.TimeUntilSwitch, expectedTimeUntilSwitch)
	}

	expectedDuration := 30 * time.Minute
	if info.IntervalDuration != expectedDuration {
		t.Errorf("Info.IntervalDuration = %v, want %v", info.IntervalDuration, expectedDuration)
	}
}

func TestSwitchbackRolloutStrategy_NoVariants(t *testing.T) {
	strategy := NewSwitchbackRolloutStrategy()

	flag := &Flag{
		Name:           "test",
		Enabled:        true,
		DefaultVariant: "default",
		Variants:       []Variant{},
	}

	ctx := Context{"user_id": "123"}
	variant, err := strategy.GetVariant(flag, ctx)
	if err != nil {
		t.Errorf("GetVariant() error = %v", err)
	}
	if variant != "default" {
		t.Errorf("GetVariant() = %v, want default", variant)
	}
}

func TestWithSwitchback(t *testing.T) {
	store := NewStore(
		WithSwitchback(
			WithIntervalMinutes(15),
		),
	)

	// Verify strategy was set
	_, ok := store.rolloutStrategy.(*SwitchbackRolloutStrategy)
	if !ok {
		t.Error("Store should have SwitchbackRolloutStrategy")
	}
}

func TestGetSwitchbackInfo(t *testing.T) {
	t.Run("with switchback strategy", func(t *testing.T) {
		store := NewStore(WithSwitchback())
		info := GetSwitchbackInfo(store)
		if info == nil {
			t.Error("GetSwitchbackInfo() should not return nil for switchback store")
		}
	})

	t.Run("without switchback strategy", func(t *testing.T) {
		store := NewStore() // default strategy
		info := GetSwitchbackInfo(store)
		if info != nil {
			t.Error("GetSwitchbackInfo() should return nil for non-switchback store")
		}
	})
}

func TestSwitchbackIntegration(t *testing.T) {
	// Integration test with full Store
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	store := NewStore(
		WithSwitchback(
			WithIntervalMinutes(60),
			WithStartTime(startTime),
			WithDailySwap(true),
		),
	)

	flag := &Flag{
		Name:           "rebate_test",
		Enabled:        true,
		DefaultVariant: "standard",
		Variants: []Variant{
			{Name: "standard", Weight: 50},
			{Name: "premium", Weight: 50},
		},
	}

	err := store.AddFlag(flag)
	if err != nil {
		t.Fatalf("AddFlag() error = %v", err)
	}

	// Get the strategy to control time
	strategy := store.rolloutStrategy.(*SwitchbackRolloutStrategy)

	ctx := Context{"driver_id": "DRV-001"}

	// Test at hour 0 (day 0, interval 0) - should be standard
	strategy.timeProvider = func() time.Time { return startTime }
	variant, enabled := store.GetVariant("rebate_test", ctx)
	if !enabled {
		t.Error("Flag should be enabled")
	}
	if variant != "standard" {
		t.Errorf("Hour 0: variant = %v, want standard", variant)
	}

	// Test at hour 1 (day 0, interval 1) - should be premium
	strategy.timeProvider = func() time.Time { return startTime.Add(1 * time.Hour) }
	variant, enabled = store.GetVariant("rebate_test", ctx)
	if !enabled {
		t.Error("Flag should be enabled")
	}
	if variant != "premium" {
		t.Errorf("Hour 1: variant = %v, want premium", variant)
	}

	// Test at hour 24 (day 1, interval 0) - should be premium (swapped)
	strategy.timeProvider = func() time.Time { return startTime.Add(24 * time.Hour) }
	variant, enabled = store.GetVariant("rebate_test", ctx)
	if !enabled {
		t.Error("Flag should be enabled")
	}
	if variant != "premium" {
		t.Errorf("Day 1 Hour 0: variant = %v, want premium (swapped)", variant)
	}

	// Test at hour 25 (day 1, interval 1) - should be standard (swapped)
	strategy.timeProvider = func() time.Time { return startTime.Add(25 * time.Hour) }
	variant, enabled = store.GetVariant("rebate_test", ctx)
	if !enabled {
		t.Error("Flag should be enabled")
	}
	if variant != "standard" {
		t.Errorf("Day 1 Hour 1: variant = %v, want standard (swapped)", variant)
	}

	// Verify all users get the same variant at the same time
	ctx2 := Context{"driver_id": "DRV-002"}
	variant2, _ := store.GetVariant("rebate_test", ctx2)
	if variant2 != variant {
		t.Errorf("Different users should get same variant: %v != %v", variant, variant2)
	}
}

func TestSwitchbackRolloutStrategy_String(t *testing.T) {
	strategy := NewSwitchbackRolloutStrategy()
	str := strategy.String()
	if str == "" {
		t.Error("String() should not be empty")
	}
	if len(str) < 10 {
		t.Error("String() should provide meaningful description")
	}
}

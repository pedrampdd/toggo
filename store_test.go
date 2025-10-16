package toggo

import (
	"testing"
)

func TestStore_AddFlag(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:    "test_flag",
		Enabled: true,
		Rollout: 100,
	}

	err := store.AddFlag(flag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify flag was added
	retrieved, err := store.GetFlag("test_flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrieved.Name != flag.Name {
		t.Errorf("expected %s, got %s", flag.Name, retrieved.Name)
	}
}

func TestStore_IsEnabled_Simple(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:    "simple_flag",
		Enabled: true,
		Rollout: 100,
	}

	store.AddFlag(flag)

	ctx := Context{"user_id": "123"}

	if !store.IsEnabled("simple_flag", ctx) {
		t.Error("expected flag to be enabled")
	}
}

func TestStore_IsEnabled_Disabled(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:    "disabled_flag",
		Enabled: false,
		Rollout: 100,
	}

	store.AddFlag(flag)

	ctx := Context{"user_id": "123"}

	if store.IsEnabled("disabled_flag", ctx) {
		t.Error("expected flag to be disabled")
	}
}

func TestStore_IsEnabled_WithConditions(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:    "conditional_flag",
		Enabled: true,
		Rollout: 100,
		Conditions: []Condition{
			{
				Attribute: "country",
				Operator:  OperatorIn,
				Value:     []interface{}{"US", "CA"},
			},
			{
				Attribute: "plan",
				Operator:  OperatorEqual,
				Value:     "premium",
			},
		},
	}

	store.AddFlag(flag)

	tests := []struct {
		name     string
		ctx      Context
		expected bool
	}{
		{
			name:     "all conditions match",
			ctx:      Context{"user_id": "123", "country": "US", "plan": "premium"},
			expected: true,
		},
		{
			name:     "country doesn't match",
			ctx:      Context{"user_id": "123", "country": "DE", "plan": "premium"},
			expected: false,
		},
		{
			name:     "plan doesn't match",
			ctx:      Context{"user_id": "123", "country": "US", "plan": "basic"},
			expected: false,
		},
		{
			name:     "missing attribute",
			ctx:      Context{"user_id": "123"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := store.IsEnabled("conditional_flag", tt.ctx)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStore_IsEnabled_Rollout(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:       "rollout_flag",
		Enabled:    true,
		Rollout:    50,
		RolloutKey: "user_id",
	}

	store.AddFlag(flag)

	// Test with multiple users - some should be enabled, some not
	enabledCount := 0
	for i := 0; i < 100; i++ {
		ctx := Context{"user_id": i}
		if store.IsEnabled("rollout_flag", ctx) {
			enabledCount++
		}
	}

	// With 50% rollout, we expect roughly 50 users enabled
	// Allow some variance (30-70%)
	if enabledCount < 30 || enabledCount > 70 {
		t.Errorf("expected roughly 50%% rollout, got %d/100", enabledCount)
	}

	// Test determinism - same user should always get same result
	ctx := Context{"user_id": "test_user_123"}
	result1 := store.IsEnabled("rollout_flag", ctx)
	result2 := store.IsEnabled("rollout_flag", ctx)
	result3 := store.IsEnabled("rollout_flag", ctx)

	if result1 != result2 || result2 != result3 {
		t.Error("rollout is not deterministic for same user")
	}
}

func TestStore_GetVariant(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:           "ab_test",
		Enabled:        true,
		DefaultVariant: "control",
		Variants: []Variant{
			{Name: "control", Weight: 50},
			{Name: "variation_a", Weight: 25},
			{Name: "variation_b", Weight: 25},
		},
	}

	store.AddFlag(flag)

	// Test with multiple users
	variantCounts := make(map[string]int)
	for i := 0; i < 100; i++ {
		ctx := Context{"user_id": i}
		variant, enabled := store.GetVariant("ab_test", ctx)
		if enabled {
			variantCounts[variant]++
		}
	}

	// Check that we have reasonable distribution
	if len(variantCounts) == 0 {
		t.Error("no variants assigned")
	}

	// Test determinism
	ctx := Context{"user_id": "test_user"}
	variant1, _ := store.GetVariant("ab_test", ctx)
	variant2, _ := store.GetVariant("ab_test", ctx)

	if variant1 != variant2 {
		t.Error("variant assignment is not deterministic")
	}
}

func TestStore_GetVariant_Disabled(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:           "disabled_ab",
		Enabled:        false,
		DefaultVariant: "control",
		Variants: []Variant{
			{Name: "control", Weight: 50},
			{Name: "variation", Weight: 50},
		},
	}

	store.AddFlag(flag)

	ctx := Context{"user_id": "123"}
	variant, enabled := store.GetVariant("disabled_ab", ctx)

	if enabled {
		t.Error("expected variant to be disabled")
	}

	if variant != "control" {
		t.Errorf("expected default variant 'control', got %s", variant)
	}
}

func TestStore_RemoveFlag(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:    "temp_flag",
		Enabled: true,
		Rollout: 100,
	}

	store.AddFlag(flag)

	// Verify it exists
	if _, err := store.GetFlag("temp_flag"); err != nil {
		t.Fatal("flag should exist")
	}

	// Remove it
	store.RemoveFlag("temp_flag")

	// Verify it's gone
	if _, err := store.GetFlag("temp_flag"); err == nil {
		t.Error("flag should not exist after removal")
	}
}

func TestStore_ListFlags(t *testing.T) {
	store := NewStore()

	flags := []*Flag{
		{Name: "flag1", Enabled: true, Rollout: 100},
		{Name: "flag2", Enabled: true, Rollout: 100},
		{Name: "flag3", Enabled: true, Rollout: 100},
	}

	for _, flag := range flags {
		store.AddFlag(flag)
	}

	names := store.ListFlags()

	if len(names) != 3 {
		t.Errorf("expected 3 flags, got %d", len(names))
	}
}

func TestStore_Clear(t *testing.T) {
	store := NewStore()

	store.AddFlag(&Flag{Name: "flag1", Enabled: true, Rollout: 100})
	store.AddFlag(&Flag{Name: "flag2", Enabled: true, Rollout: 100})

	if store.Size() != 2 {
		t.Errorf("expected size 2, got %d", store.Size())
	}

	store.Clear()

	if store.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", store.Size())
	}
}

func TestStore_ThreadSafety(t *testing.T) {
	store := NewStore()

	flag := &Flag{
		Name:    "concurrent_flag",
		Enabled: true,
		Rollout: 100,
	}

	store.AddFlag(flag)

	// Run concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			ctx := Context{"user_id": id}
			for j := 0; j < 100; j++ {
				store.IsEnabled("concurrent_flag", ctx)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

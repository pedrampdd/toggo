package toggo

import (
	"sync"
)

// Store manages feature flags and provides thread-safe evaluation
type Store struct {
	mu              sync.RWMutex
	flags           map[string]*Flag
	evaluator       *conditionEvaluator
	rolloutStrategy RolloutStrategy
}

// StoreOption is a functional option for configuring the Store
type StoreOption func(*Store)

// NewStore creates a new feature flag store
func NewStore(opts ...StoreOption) *Store {
	store := &Store{
		flags:           make(map[string]*Flag),
		evaluator:       newConditionEvaluator(),
		rolloutStrategy: NewDefaultRolloutStrategy(nil),
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

// AddFlag adds or updates a flag in the store
func (s *Store) AddFlag(flag *Flag) error {
	if err := flag.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.flags[flag.Name] = flag
	return nil
}

// AddFlags adds multiple flags to the store
func (s *Store) AddFlags(flags []*Flag) error {
	for _, flag := range flags {
		if err := s.AddFlag(flag); err != nil {
			return err
		}
	}
	return nil
}

// RemoveFlag removes a flag from the store
func (s *Store) RemoveFlag(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.flags, name)
}

// GetFlag retrieves a flag by name
func (s *Store) GetFlag(name string) (*Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	flag, ok := s.flags[name]
	if !ok {
		return nil, ErrFlagNotFound
	}

	return flag, nil
}

// ListFlags returns all flag names
func (s *Store) ListFlags() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.flags))
	for name := range s.flags {
		names = append(names, name)
	}

	return names
}

// IsEnabled checks if a feature flag is enabled for the given context
// This is the primary method for simple on/off feature flags
func (s *Store) IsEnabled(name string, ctx Context) bool {
	result, _ := s.IsEnabledWithError(name, ctx)
	return result
}

// IsEnabledWithError checks if a feature flag is enabled and returns any error
func (s *Store) IsEnabledWithError(name string, ctx Context) (bool, error) {
	flag, err := s.GetFlag(name)
	if err != nil {
		return false, err
	}

	// If flag is disabled, return false immediately
	if !flag.Enabled {
		return false, nil
	}

	// If flag has variants, IsEnabled should return false
	// User should use GetVariant instead
	if flag.HasVariants() {
		return false, nil
	}

	// Evaluate all conditions
	match, err := s.evaluator.evaluateAll(flag.Conditions, ctx)
	if err != nil {
		return false, err
	}

	// If conditions don't match, return false
	if !match {
		return false, nil
	}

	// Apply rollout strategy
	shouldRollout, err := s.rolloutStrategy.ShouldRollout(flag, ctx)
	if err != nil {
		return false, err
	}

	return shouldRollout, nil
}

// GetVariant returns the variant for A/B testing
// Returns the variant name and whether the flag is enabled
func (s *Store) GetVariant(name string, ctx Context) (string, bool) {
	variant, enabled, _ := s.GetVariantWithError(name, ctx)
	return variant, enabled
}

// GetVariantWithError returns the variant with detailed error information
func (s *Store) GetVariantWithError(name string, ctx Context) (string, bool, error) {
	flag, err := s.GetFlag(name)
	if err != nil {
		return "", false, err
	}

	// If flag is disabled, return default variant
	if !flag.Enabled {
		return flag.DefaultVariant, false, nil
	}

	// Evaluate global flag conditions
	match, err := s.evaluator.evaluateAll(flag.Conditions, ctx)
	if err != nil {
		return "", false, err
	}

	// If global conditions don't match, return default variant
	if !match {
		return flag.DefaultVariant, false, nil
	}

	// If no variants configured, this is a simple on/off flag
	if !flag.HasVariants() {
		// Apply rollout
		shouldRollout, err := s.rolloutStrategy.ShouldRollout(flag, ctx)
		if err != nil {
			return "", false, err
		}
		if shouldRollout {
			return "on", true, nil
		}
		return "off", false, nil
	}

	// Get variant based on rollout strategy
	variantName, err := s.rolloutStrategy.GetVariant(flag, ctx)
	if err != nil {
		return "", false, err
	}

	// Find the variant and check its conditions
	for _, variant := range flag.Variants {
		if variant.Name == variantName {
			// Evaluate variant-specific conditions if any
			if len(variant.Conditions) > 0 {
				match, err := s.evaluator.evaluateAll(variant.Conditions, ctx)
				if err != nil {
					return "", false, err
				}
				if !match {
					return flag.DefaultVariant, false, nil
				}
			}
			return variant.Name, true, nil
		}
	}

	return flag.DefaultVariant, false, nil
}

// Clear removes all flags from the store
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.flags = make(map[string]*Flag)
}

// Size returns the number of flags in the store
func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.flags)
}

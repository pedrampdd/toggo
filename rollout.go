package toggo

import (
	"fmt"

	"github.com/pedram/toggo/internal/hash"
)

// RolloutStrategy defines how rollout decisions are made
type RolloutStrategy interface {
	// ShouldRollout determines if a flag should be enabled based on rollout percentage
	ShouldRollout(flag *Flag, ctx Context) (bool, error)

	// GetVariant determines which variant to return for A/B testing
	GetVariant(flag *Flag, ctx Context) (string, error)
}

// DefaultRolloutStrategy implements standard percentage-based rollout
type DefaultRolloutStrategy struct {
	hasher hash.Hasher
}

// NewDefaultRolloutStrategy creates a new default rollout strategy
func NewDefaultRolloutStrategy(hasher hash.Hasher) *DefaultRolloutStrategy {
	if hasher == nil {
		hasher = hash.NewFNV()
	}
	return &DefaultRolloutStrategy{
		hasher: hasher,
	}
}

// ShouldRollout determines if the flag should be enabled based on rollout percentage
func (r *DefaultRolloutStrategy) ShouldRollout(flag *Flag, ctx Context) (bool, error) {
	// If rollout is 100, always return true
	if flag.Rollout >= 100 {
		return true, nil
	}

	// If rollout is 0, always return false
	if flag.Rollout <= 0 {
		return false, nil
	}

	// Get the rollout key value from context
	rolloutKey := flag.GetRolloutKey()
	keyValue, exists := ctx.Get(rolloutKey)
	if !exists {
		// If rollout key is missing, we can't make a consistent decision
		// Return false to be conservative
		return false, nil
	}

	// Create deterministic hash key
	hashKey := fmt.Sprintf("%s:%s", flag.Name, fmt.Sprint(keyValue))
	hashValue := r.hasher.Hash(hashKey)

	// Check if hash falls within rollout percentage
	return hashValue < flag.Rollout, nil
}

// GetVariant determines which variant to return based on weights
func (r *DefaultRolloutStrategy) GetVariant(flag *Flag, ctx Context) (string, error) {
	if !flag.HasVariants() {
		return flag.DefaultVariant, nil
	}

	// Get the rollout key value from context
	rolloutKey := flag.GetRolloutKey()
	keyValue, exists := ctx.Get(rolloutKey)
	if !exists {
		return flag.DefaultVariant, nil
	}

	// Create deterministic hash key for variant selection
	hashKey := fmt.Sprintf("%s:variant:%s", flag.Name, fmt.Sprint(keyValue))
	hashValue := r.hasher.Hash(hashKey)

	// Find the variant based on cumulative weights
	cumulative := 0
	for _, variant := range flag.Variants {
		cumulative += variant.Weight
		if hashValue < cumulative {
			return variant.Name, nil
		}
	}

	// If no variant matched (shouldn't happen with proper config), return default
	return flag.DefaultVariant, nil
}

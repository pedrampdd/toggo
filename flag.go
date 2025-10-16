package toggo

// Flag represents a feature flag configuration
type Flag struct {
	// Name is the unique identifier for this flag
	Name string `json:"name" yaml:"name"`

	// Enabled controls whether this flag is active
	Enabled bool `json:"enabled" yaml:"enabled"`

	// Rollout is the percentage (0-100) of users who should see this flag
	// when all conditions are met
	Rollout int `json:"rollout,omitempty" yaml:"rollout,omitempty"`

	// RolloutKey specifies which context attribute to use for rollout hashing
	// Defaults to "user_id" if not specified
	RolloutKey string `json:"rollout_key,omitempty" yaml:"rollout_key,omitempty"`

	// Conditions are the rules that must ALL be satisfied for the flag to be enabled
	Conditions []Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`

	// Variants enables A/B testing with multiple variations
	// If set, IsEnabled returns false and GetVariant should be used instead
	Variants []Variant `json:"variants,omitempty" yaml:"variants,omitempty"`

	// DefaultVariant is returned when no variant matches
	DefaultVariant string `json:"default_variant,omitempty" yaml:"default_variant,omitempty"`
}

// Variant represents an A/B test variant
type Variant struct {
	// Name is the variant identifier
	Name string `json:"name" yaml:"name"`

	// Weight is the percentage (0-100) of traffic allocated to this variant
	Weight int `json:"weight" yaml:"weight"`

	// Conditions are additional conditions specific to this variant
	Conditions []Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
}

// Validate checks if the flag configuration is valid
func (f *Flag) Validate() error {
	if f.Name == "" {
		return ErrInvalidCondition
	}

	if f.Rollout < 0 || f.Rollout > 100 {
		return ErrInvalidRollout
	}

	for _, cond := range f.Conditions {
		if err := cond.Validate(); err != nil {
			return err
		}
	}

	// Validate variants
	totalWeight := 0
	for _, variant := range f.Variants {
		if variant.Weight < 0 || variant.Weight > 100 {
			return ErrInvalidRollout
		}
		totalWeight += variant.Weight
		for _, cond := range variant.Conditions {
			if err := cond.Validate(); err != nil {
				return err
			}
		}
	}

	if len(f.Variants) > 0 && totalWeight > 100 {
		return ErrInvalidRollout
	}

	return nil
}

// HasVariants returns true if this flag has A/B test variants configured
func (f *Flag) HasVariants() bool {
	return len(f.Variants) > 0
}

// GetRolloutKey returns the key to use for rollout hashing
func (f *Flag) GetRolloutKey() string {
	if f.RolloutKey != "" {
		return f.RolloutKey
	}
	return "user_id" // default
}

package toggo

import (
	"fmt"
	"time"

	"github.com/pedrampdd/toggo/internal/hash"
)

// SwitchbackRolloutStrategy implements time-based switchback testing
// In switchback tests, all users see the same variant at the same time,
// and the variant switches at regular intervals
type SwitchbackRolloutStrategy struct {
	baseStrategy    *DefaultRolloutStrategy
	intervalMinutes int
	startTime       time.Time
	swapDaily       bool
	timeProvider    func() time.Time
}

// SwitchbackOption configures a switchback strategy
type SwitchbackOption func(*SwitchbackRolloutStrategy)

// WithIntervalMinutes sets the duration of each switchback interval in minutes
func WithIntervalMinutes(minutes int) SwitchbackOption {
	return func(s *SwitchbackRolloutStrategy) {
		s.intervalMinutes = minutes
	}
}

// WithStartTime sets the reference start time for calculating intervals
func WithStartTime(t time.Time) SwitchbackOption {
	return func(s *SwitchbackRolloutStrategy) {
		s.startTime = t
	}
}

// WithDailySwap enables swapping the variant order on alternating days
// Day 0: variants in order, Day 1: variants in reverse order, etc.
func WithDailySwap(enabled bool) SwitchbackOption {
	return func(s *SwitchbackRolloutStrategy) {
		s.swapDaily = enabled
	}
}

// NewSwitchbackRolloutStrategy creates a new switchback rollout strategy
func NewSwitchbackRolloutStrategy(opts ...SwitchbackOption) *SwitchbackRolloutStrategy {
	s := &SwitchbackRolloutStrategy{
		baseStrategy:    NewDefaultRolloutStrategy(hash.NewFNV()),
		intervalMinutes: 30, // default 30 minutes
		startTime:       time.Now().Truncate(24 * time.Hour),
		swapDaily:       false,
		timeProvider:    time.Now,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// GetCurrentInterval returns which time interval we're currently in
func (s *SwitchbackRolloutStrategy) GetCurrentInterval() int {
	now := s.timeProvider()
	elapsed := now.Sub(s.startTime)
	intervalDuration := time.Duration(s.intervalMinutes) * time.Minute
	return int(elapsed / intervalDuration)
}

// GetCurrentDay returns which day number we're in since start time
func (s *SwitchbackRolloutStrategy) GetCurrentDay() int {
	now := s.timeProvider()
	elapsed := now.Sub(s.startTime)
	return int(elapsed / (24 * time.Hour))
}

// GetTimeUntilNextSwitch returns how much time until the next interval switch
func (s *SwitchbackRolloutStrategy) GetTimeUntilNextSwitch() time.Duration {
	now := s.timeProvider()
	elapsed := now.Sub(s.startTime)
	intervalDuration := time.Duration(s.intervalMinutes) * time.Minute
	currentInterval := int(elapsed / intervalDuration)
	nextSwitchTime := s.startTime.Add(time.Duration(currentInterval+1) * intervalDuration)
	return nextSwitchTime.Sub(now)
}

// ShouldRollout always returns true for switchback tests since all users participate
func (s *SwitchbackRolloutStrategy) ShouldRollout(flag *Flag, ctx Context) (bool, error) {
	return true, nil
}

// GetVariant returns the current variant based on time interval
// All users get the same variant at the same time
func (s *SwitchbackRolloutStrategy) GetVariant(flag *Flag, ctx Context) (string, error) {
	if !flag.HasVariants() {
		return flag.DefaultVariant, nil
	}

	intervalNum := s.GetCurrentInterval()
	dayNum := s.GetCurrentDay()

	// Calculate which variant index to use
	numVariants := len(flag.Variants)
	if numVariants == 0 {
		return flag.DefaultVariant, nil
	}

	// Determine base index from interval
	variantIndex := intervalNum % numVariants

	// If daily swap is enabled and we're on an odd day, reverse the order
	if s.swapDaily && dayNum%2 == 1 {
		variantIndex = (numVariants - 1) - variantIndex
	}

	return flag.Variants[variantIndex].Name, nil
}

// GetSwitchbackInfo returns detailed information about current switchback state
type SwitchbackInfo struct {
	CurrentInterval  int
	CurrentDay       int
	TimeUntilSwitch  time.Duration
	IntervalDuration time.Duration
}

// GetInfo returns detailed switchback timing information
func (s *SwitchbackRolloutStrategy) GetInfo() SwitchbackInfo {
	return SwitchbackInfo{
		CurrentInterval:  s.GetCurrentInterval(),
		CurrentDay:       s.GetCurrentDay(),
		TimeUntilSwitch:  s.GetTimeUntilNextSwitch(),
		IntervalDuration: time.Duration(s.intervalMinutes) * time.Minute,
	}
}

// WithSwitchback is a StoreOption that configures switchback testing
func WithSwitchback(opts ...SwitchbackOption) StoreOption {
	return func(store *Store) {
		store.rolloutStrategy = NewSwitchbackRolloutStrategy(opts...)
	}
}

// GetSwitchbackInfo is a convenience method to get switchback info from a store
// Returns nil if the store is not using switchback strategy
func GetSwitchbackInfo(store *Store) *SwitchbackInfo {
	if strategy, ok := store.rolloutStrategy.(*SwitchbackRolloutStrategy); ok {
		info := strategy.GetInfo()
		return &info
	}
	return nil
}

// String returns a human-readable description of the switchback state
func (s *SwitchbackRolloutStrategy) String() string {
	info := s.GetInfo()
	return fmt.Sprintf(
		"Switchback: Interval %d, Day %d, Next switch in %v",
		info.CurrentInterval,
		info.CurrentDay,
		info.TimeUntilSwitch.Round(time.Second),
	)
}

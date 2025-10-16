package loader

import (
	"github.com/pedram/toggo"
)

// Loader defines the interface for loading feature flags from various sources
type Loader interface {
	// Load reads flags from a source and returns them
	Load() ([]*toggo.Flag, error)
}

// Config represents the structure of a feature flags configuration file
type Config struct {
	Flags []*toggo.Flag `json:"flags" yaml:"flags"`
}

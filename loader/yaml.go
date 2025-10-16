package loader

import (
	"io"
	"os"

	"github.com/pedrampdd/toggo"
	"gopkg.in/yaml.v3"
)

// YAMLLoader loads feature flags from YAML files or readers
type YAMLLoader struct {
	source interface{} // can be string (file path) or io.Reader
}

// NewYAMLFile creates a loader that reads from a YAML file
func NewYAMLFile(filepath string) *YAMLLoader {
	return &YAMLLoader{source: filepath}
}

// NewYAMLReader creates a loader that reads from an io.Reader
func NewYAMLReader(reader io.Reader) *YAMLLoader {
	return &YAMLLoader{source: reader}
}

// Load reads and parses the YAML configuration
func (l *YAMLLoader) Load() ([]*toggo.Flag, error) {
	var reader io.Reader

	switch src := l.source.(type) {
	case string:
		file, err := os.Open(src)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	case io.Reader:
		reader = src
	}

	var config Config
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	// Validate all flags
	for _, flag := range config.Flags {
		if err := flag.Validate(); err != nil {
			return nil, err
		}
	}

	return config.Flags, nil
}

// LoadIntoStore is a convenience method that loads flags directly into a store
func (l *YAMLLoader) LoadIntoStore(store *toggo.Store) error {
	flags, err := l.Load()
	if err != nil {
		return err
	}
	return store.AddFlags(flags)
}

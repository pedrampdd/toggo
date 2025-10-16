package loader

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pedram/toggo"
)

// JSONLoader loads feature flags from JSON files or readers
type JSONLoader struct {
	source interface{} // can be string (file path) or io.Reader
}

// NewJSONFile creates a loader that reads from a JSON file
func NewJSONFile(filepath string) *JSONLoader {
	return &JSONLoader{source: filepath}
}

// NewJSONReader creates a loader that reads from an io.Reader
func NewJSONReader(reader io.Reader) *JSONLoader {
	return &JSONLoader{source: reader}
}

// Load reads and parses the JSON configuration
func (l *JSONLoader) Load() ([]*toggo.Flag, error) {
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
	decoder := json.NewDecoder(reader)
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
func (l *JSONLoader) LoadIntoStore(store *toggo.Store) error {
	flags, err := l.Load()
	if err != nil {
		return err
	}
	return store.AddFlags(flags)
}

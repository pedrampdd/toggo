# Toggo 🚀

A flexible, performant, and production-ready feature flag and A/B testing SDK for Go.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Features

- ✨ **Simple on/off feature flags** - Control features with boolean flags
- 📊 **Percentage-based rollouts** - Gradually roll out features with deterministic hashing
- 🎯 **Conditional targeting** - Target users based on attributes (country, plan, custom fields)
- 🧪 **A/B testing** - Run experiments with multiple variants
- 🔒 **Thread-safe** - Safe for concurrent access
- 📝 **JSON/YAML configuration** - Load flags from configuration files
- 🎨 **Flexible operators** - Support for ==, !=, in, >, <, contains, regex, and more
- 🚀 **Zero dependencies** (except yaml parser)
- 📦 **Clean API** - Simple and intuitive interface

## Installation

```bash
go get github.com/pedrampdd/toggo
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/pedrampdd/toggo"
)

func main() {
    // Create a feature flag store
    store := toggo.NewStore()

    // Define a feature flag
    flag := &toggo.Flag{
        Name:    "new_checkout",
        Enabled: true,
        Rollout: 50, // 50% of users
    }
    
    store.AddFlag(flag)

    // Check if enabled for a user
    ctx := toggo.Context{
        "user_id": "12345",
        "country": "US",
    }

    if store.IsEnabled("new_checkout", ctx) {
        // Show new checkout flow
    }
}
```

## Core Concepts

### Context

A `Context` is a map of user attributes used for flag evaluation:

```go
ctx := toggo.Context{
    "user_id": "12345",
    "country": "US",
    "plan":    "premium",
    "age":     25,
}
```

### Flags

Flags control feature availability:

```go
flag := &toggo.Flag{
    Name:       "feature_name",
    Enabled:    true,
    Rollout:    100,        // 0-100 percentage
    RolloutKey: "user_id",  // Context key for hashing (default: "user_id")
    Conditions: []toggo.Condition{
        // Optional targeting conditions
    },
}
```

### Conditions

Target specific users with conditions:

```go
condition := toggo.Condition{
    Attribute: "country",
    Operator:  toggo.OperatorIn,
    Value:     []interface{}{"US", "CA", "UK"},
}
```

### Supported Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `==` | Equal | `plan == "premium"` |
| `!=` | Not equal | `country != "US"` |
| `in` | In list | `country in ["US", "CA"]` |
| `not_in` | Not in list | `country not_in ["DE", "FR"]` |
| `>` | Greater than | `age > 18` |
| `>=` | Greater than or equal | `age >= 21` |
| `<` | Less than | `age < 65` |
| `<=` | Less than or equal | `age <= 25` |
| `contains` | String contains | `email contains "@company.com"` |
| `starts_with` | String starts with | `name starts_with "John"` |
| `ends_with` | String ends with | `file ends_with ".pdf"` |
| `regex` | Regex match | `email regex ".*@example\\.com"` |

## Usage Examples

### Simple Feature Flag

```go
store := toggo.NewStore()

flag := &toggo.Flag{
    Name:    "dark_mode",
    Enabled: true,
    Rollout: 100,
}

store.AddFlag(flag)

ctx := toggo.Context{"user_id": "123"}
if store.IsEnabled("dark_mode", ctx) {
    // Enable dark mode
}
```

### Percentage Rollout

```go
flag := &toggo.Flag{
    Name:       "new_ui",
    Enabled:    true,
    Rollout:    25, // 25% of users
    RolloutKey: "user_id",
}

store.AddFlag(flag)

// Same user always gets same result (deterministic)
ctx := toggo.Context{"user_id": "user_42"}
enabled := store.IsEnabled("new_ui", ctx) // Consistent for this user
```

### Conditional Targeting

```go
flag := &toggo.Flag{
    Name:    "premium_feature",
    Enabled: true,
    Rollout: 100,
    Conditions: []toggo.Condition{
        {
            Attribute: "plan",
            Operator:  toggo.OperatorEqual,
            Value:     "premium",
        },
        {
            Attribute: "country",
            Operator:  toggo.OperatorIn,
            Value:     []interface{}{"US", "CA", "UK"},
        },
    },
}

store.AddFlag(flag)

ctx := toggo.Context{
    "user_id": "123",
    "plan":    "premium",
    "country": "US",
}

// Enabled only if ALL conditions match
if store.IsEnabled("premium_feature", ctx) {
    // Show premium feature
}
```

### A/B Testing

```go
flag := &toggo.Flag{
    Name:           "pricing_test",
    Enabled:        true,
    RolloutKey:     "user_id",
    DefaultVariant: "control",
    Variants: []toggo.Variant{
        {Name: "control", Weight: 33},
        {Name: "price_low", Weight: 33},
        {Name: "price_high", Weight: 34},
    },
}

store.AddFlag(flag)

ctx := toggo.Context{"user_id": "user_42"}
variant, _ := store.GetVariant("pricing_test", ctx)

switch variant {
case "control":
    price = 99.99
case "price_low":
    price = 79.99
case "price_high":
    price = 119.99
}
```

### Loading from Configuration Files

#### JSON

```json
{
  "flags": [
    {
      "name": "new_checkout",
      "enabled": true,
      "rollout": 50,
      "conditions": [
        {
          "attribute": "country",
          "operator": "in",
          "value": ["US", "CA"]
        }
      ]
    }
  ]
}
```

```go
import "github.com/pedrampdd/toggo/loader"

store := toggo.NewStore()
l := loader.NewJSONFile("flags.json")
l.LoadIntoStore(store)
```

#### YAML

```yaml
flags:
  - name: dark_mode
    enabled: true
    rollout: 100
  - name: beta_features
    enabled: true
    rollout: 25
    conditions:
      - attribute: beta_tester
        operator: "=="
        value: true
```

```go
import "github.com/pedrampdd/toggo/loader"

store := toggo.NewStore()
l := loader.NewYAMLFile("flags.yaml")
l.LoadIntoStore(store)
```

## API Reference

### Store

#### `NewStore(opts ...StoreOption) *Store`

Creates a new feature flag store.

#### `AddFlag(flag *Flag) error`

Adds or updates a flag in the store. Returns error if validation fails.

#### `IsEnabled(name string, ctx Context) bool`

Checks if a feature flag is enabled for the given context. Returns `false` if flag not found or conditions don't match.

#### `GetVariant(name string, ctx Context) (string, bool)`

Returns the variant name for A/B testing. Second return value indicates if flag is enabled.

#### `GetFlag(name string) (*Flag, error)`

Retrieves a flag by name. Returns `ErrFlagNotFound` if not found.

#### `ListFlags() []string`

Returns all flag names.

#### `RemoveFlag(name string)`

Removes a flag from the store.

#### `Clear()`

Removes all flags from the store.

#### `Size() int`

Returns the number of flags in the store.

### Flag

```go
type Flag struct {
    Name           string
    Enabled        bool
    Rollout        int        // 0-100
    RolloutKey     string     // Default: "user_id"
    Conditions     []Condition
    Variants       []Variant
    DefaultVariant string
}
```

### Condition

```go
type Condition struct {
    Attribute string
    Operator  Operator
    Value     interface{}
    Negate    bool
}
```

### Variant

```go
type Variant struct {
    Name       string
    Weight     int           // 0-100
    Conditions []Condition
}
```

## Project Structure

```
toggo/
├── toggo.go              # Main package file with documentation
├── context.go            # Context type and methods
├── flag.go              # Flag and Variant types
├── condition.go         # Condition type
├── operator.go          # Operator constants
├── store.go             # Store implementation
├── rollout.go           # Rollout strategy
├── errors.go            # Error definitions
├── internal/            # Internal implementation details
│   ├── evaluator/      # Condition evaluation logic
│   └── hash/           # Hashing for rollouts
├── loader/             # Configuration loaders
│   ├── json.go
│   └── yaml.go
├── examples/           # Usage examples
│   ├── simple/
│   ├── abtest/
│   ├── conditional/
│   └── config_loader/
└── testdata/          # Test fixtures
```

## Testing

Run all tests:

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

Run specific package:

```bash
go test ./internal/evaluator
```

## Examples

Explore the `examples/` directory for complete working examples:

- **simple** - Basic feature flag usage
- **abtest** - A/B testing with variants
- **conditional** - Conditional targeting
- **config_loader** - Loading flags from JSON/YAML

Run an example:

```bash
cd examples/simple
go run main.go
```

## Best Practices

1. **Use deterministic rollout keys** - Always use stable user identifiers (user_id, session_id) for rollout keys to ensure consistent experience.

2. **Validate flags** - Flags are validated when added to the store. Handle errors appropriately.

3. **Keep conditions simple** - Complex condition trees can impact performance. Consider splitting into multiple flags.

4. **Use variants for A/B tests** - Don't use multiple flags for variants of the same experiment.

5. **Load from config files** - Store flag definitions in version-controlled YAML/JSON files for easier management.

6. **Test coverage** - Always test both enabled and disabled states of features.

## Performance

- **Thread-safe** - Uses `sync.RWMutex` for concurrent reads
- **Fast evaluation** - O(1) flag lookup, O(n) condition evaluation where n is number of conditions
- **Deterministic hashing** - FNV-1a hash for consistent, fast rollout decisions
- **Zero allocations** - Designed to minimize allocations in hot paths

## Roadmap

- [ ] Remote flag management integration
- [ ] Metrics and analytics hooks
- [ ] Flag scheduling (enable/disable at specific times)
- [ ] User segments for reusable targeting
- [ ] Admin UI for flag management
- [ ] WebSocket/SSE for real-time flag updates

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Authors

Built with ❤️ for the Go community

---

**Questions?** Open an issue or start a discussion!

# Toggo Architecture

This document describes the architecture and design decisions behind Toggo.

## Design Principles

1. **Simplicity** - Easy to use API with sensible defaults
2. **Performance** - Thread-safe with minimal allocations
3. **Flexibility** - Extensible through interfaces and options
4. **Type Safety** - Strong typing with clear error handling
5. **Best Practices** - Follows Go idioms and conventions

## Package Structure

```
toggo/
├── toggo.go              # Package documentation and exports
├── context.go            # Context type for evaluation
├── flag.go              # Flag and Variant definitions
├── condition.go         # Condition type
├── operator.go          # Operator constants
├── store.go             # Main Store implementation
├── evaluator.go         # Condition evaluation logic
├── rollout.go           # Rollout strategy implementation
├── errors.go            # Error definitions
├── internal/            # Internal implementation details
│   └── hash/           # Hashing for deterministic rollouts
│       ├── hasher.go   # Hasher interface
│       └── fnv.go      # FNV-1a implementation
└── loader/             # Configuration loaders
    ├── loader.go       # Loader interface
    ├── json.go         # JSON loader
    └── yaml.go         # YAML loader
```

## Core Components

### 1. Context

The `Context` is a simple `map[string]interface{}` that holds arbitrary user attributes:

```go
type Context map[string]interface{}
```

**Design Decisions:**
- Uses `interface{}` for maximum flexibility
- Provides helper methods for type-safe access
- No struct definition allows dynamic attributes
- Users can add any custom fields

### 2. Store

The `Store` is the main entry point for flag evaluation:

```go
type Store struct {
    mu              sync.RWMutex
    flags           map[string]*Flag
    evaluator       *conditionEvaluator
    rolloutStrategy RolloutStrategy
}
```

**Design Decisions:**
- Uses `sync.RWMutex` for thread-safe concurrent access
- Read-heavy workload optimized with RLock
- Functional options pattern for configuration
- Composition over inheritance for flexibility

**Thread Safety:**
- Readers use `RLock()` for concurrent reads
- Writers use `Lock()` for exclusive access
- No lock held during evaluation (flags are immutable once added)

### 3. Flag

Flags represent feature configurations:

```go
type Flag struct {
    Name           string
    Enabled        bool
    Rollout        int        // 0-100
    RolloutKey     string
    Conditions     []Condition
    Variants       []Variant
    DefaultVariant string
}
```

**Design Decisions:**
- Validation on add prevents runtime errors
- Supports both simple flags and A/B tests
- Conditions use AND logic (all must match)
- Variants enable multi-variate testing

### 4. Condition Evaluation

The `conditionEvaluator` handles all condition matching logic:

```go
type conditionEvaluator struct{}
```

**Design Decisions:**
- Unexported to keep implementation flexible
- Stateless for thread-safety
- Type coercion for numeric comparisons
- String fallback for non-numeric types
- Short-circuit evaluation for performance

**Operator Support:**
- Equality: `==`, `!=`
- Membership: `in`, `not_in`
- Comparison: `>`, `>=`, `<`, `<=`
- String: `contains`, `starts_with`, `ends_with`
- Pattern: `regex`

### 5. Rollout Strategy

Determines which users see a feature based on percentage:

```go
type RolloutStrategy interface {
    ShouldRollout(flag *Flag, ctx Context) (bool, error)
    GetVariant(flag *Flag, ctx Context) (string, error)
}
```

**Design Decisions:**
- Interface allows custom strategies
- Default uses FNV-1a hash for determinism
- Hash key: `"flagname:userid"` ensures consistency
- Modulo 100 gives percentage buckets

**Deterministic Hashing:**
```
hash = FNV1a(flagName + ":" + userId) % 100
enabled = hash < rolloutPercentage
```

This ensures:
- Same user always gets same result
- Uniform distribution across users
- Independent decisions per flag
- No coordination needed

### 6. Configuration Loaders

Support for external configuration:

```go
type Loader interface {
    Load() ([]*Flag, error)
}
```

**Design Decisions:**
- JSON and YAML support out of the box
- Can load from files or io.Reader
- Validates on load
- Convenience method `LoadIntoStore()`

## Data Flow

### IsEnabled Flow

```
User Request
    ↓
Store.IsEnabled(name, ctx)
    ↓
1. Get flag (with RLock)
    ↓
2. Check flag.Enabled
    ↓
3. Evaluate all conditions (AND logic)
    ↓
4. Apply rollout strategy
    ↓
Return boolean
```

### GetVariant Flow

```
User Request
    ↓
Store.GetVariant(name, ctx)
    ↓
1. Get flag (with RLock)
    ↓
2. Check flag.Enabled
    ↓
3. Evaluate global conditions
    ↓
4. Select variant based on hash
    ↓
5. Evaluate variant conditions
    ↓
Return variant name
```

## Performance Characteristics

### Time Complexity
- Flag lookup: **O(1)** - map lookup
- Condition evaluation: **O(n)** - n = number of conditions
- Variant selection: **O(m)** - m = number of variants
- Overall: **O(n + m)** per evaluation

### Space Complexity
- Store: **O(f)** - f = number of flags
- Context: **O(a)** - a = number of attributes
- No allocation during evaluation

### Concurrency
- Read operations: Fully concurrent
- Write operations: Exclusive lock
- No contention for read-heavy workloads

## Design Patterns Used

### 1. Functional Options Pattern
```go
store := NewStore(
    WithRolloutStrategy(customStrategy),
)
```

**Benefits:**
- Optional configuration
- Backward compatible
- Self-documenting
- Extensible

### 2. Strategy Pattern
```go
type RolloutStrategy interface {
    ShouldRollout(flag *Flag, ctx Context) (bool, error)
}
```

**Benefits:**
- Pluggable algorithms
- Easy testing
- Runtime selection

### 3. Builder Pattern (Fluent Configuration)
```go
flag := &Flag{
    Name: "feature",
    Enabled: true,
    Conditions: []Condition{...},
}
```

### 4. Interface Segregation
- Small, focused interfaces
- Easy to implement
- Minimal dependencies

## Internal Package

The `internal/` directory follows Go's internal package convention:

```
internal/
└── hash/
    ├── hasher.go  # Interface
    └── fnv.go     # Implementation
```

**Why internal?**
- Implementation details hidden
- Can change without breaking users
- Enforces encapsulation
- Only exported toggo types are public API

## Error Handling

Errors are predefined and exported:

```go
var (
    ErrFlagNotFound     = errors.New("flag not found")
    ErrInvalidOperator  = errors.New("invalid operator")
    ErrInvalidRollout   = errors.New("rollout must be between 0 and 100")
    // ...
)
```

**Error Strategy:**
- Fail fast on configuration errors
- Fail safe on runtime errors (return false)
- Provide detailed error messages
- Both error and non-error variants of methods

## Testing Strategy

### Unit Tests
- Test each component in isolation
- Table-driven tests
- Edge case coverage
- Mock dependencies

### Integration Tests
- Load from config files
- End-to-end flag evaluation
- Thread safety tests

### Test Coverage
- Aim for >90% coverage
- Focus on business logic
- Test error paths

## Extension Points

Users can extend Toggo by:

1. **Custom Rollout Strategies**
   ```go
   type MyStrategy struct{}
   func (s *MyStrategy) ShouldRollout(...) (bool, error) {}
   ```

2. **Custom Loaders**
   ```go
   type DatabaseLoader struct{}
   func (l *DatabaseLoader) Load() ([]*Flag, error) {}
   ```

3. **Custom Operators** (future)
   - Would require evaluator interface exposure

## Scalability Considerations

### Horizontal Scaling
- Stateless evaluation
- No coordination needed
- Load flags independently per instance

### Vertical Scaling
- Thread-safe for multiple goroutines
- Read-optimized with RWMutex
- Minimal memory footprint

### Large Scale
- For 1000+ flags: consider caching
- For complex rules: consider rule engine
- For real-time updates: add watch mechanism

## Security Considerations

1. **Input Validation**
   - Flags validated on add
   - Conditions validated on evaluation

2. **Resource Limits**
   - No recursion (prevents stack overflow)
   - Bounded iterations

3. **Regex Safety**
   - Use Go's regexp package (safe)
   - Consider timeout for complex patterns

## Future Enhancements

### Potential Features
- [ ] Remote configuration sync
- [ ] Real-time updates via WebSocket
- [ ] Metrics and analytics
- [ ] Admin UI
- [ ] User segments
- [ ] Scheduled rollouts
- [ ] Dependency between flags
- [ ] OR logic for conditions
- [ ] Custom evaluation context per flag

### Backward Compatibility
- Follow semantic versioning
- Deprecate before removing
- Maintain stable API surface

## Conclusion

Toggo is designed to be:
- **Simple** - Easy to understand and use
- **Fast** - Optimized for read-heavy workloads
- **Safe** - Thread-safe and type-safe
- **Flexible** - Extensible through interfaces
- **Maintainable** - Clean architecture and good tests

The architecture follows Go best practices and provides a solid foundation for feature flag management and A/B testing.


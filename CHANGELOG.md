# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-16

### Added
- Initial release of Toggo feature flag and A/B testing SDK
- Core feature flag functionality with `Store` type
- Context-based evaluation with flexible attributes
- Support for multiple operators: `==`, `!=`, `in`, `not_in`, `>`, `>=`, `<`, `<=`, `contains`, `starts_with`, `ends_with`, `regex`
- Percentage-based rollout with deterministic hashing (FNV-1a)
- A/B testing with multiple variants and weight distribution
- Conditional targeting with AND logic
- Thread-safe concurrent access with RWMutex
- JSON configuration loader
- YAML configuration loader
- Comprehensive test coverage (>95%)
- Full documentation and examples
- Four example applications (simple, abtest, conditional, config_loader)

### Features
- **Simple Flags**: Boolean on/off feature flags
- **Rollout**: Gradual percentage-based rollouts (0-100%)
- **Conditions**: Target users based on attributes
- **Variants**: Multi-variate A/B testing
- **Thread-Safe**: Safe for concurrent use
- **Deterministic**: Same user always gets same result
- **Flexible**: Dynamic context attributes
- **Validated**: Configuration validated on load

### Documentation
- README.md with quick start and examples
- ARCHITECTURE.md with design decisions
- CONTRIBUTING.md with contribution guidelines
- Inline documentation for all exported types
- Complete example applications

### Performance
- O(1) flag lookup
- O(n) condition evaluation
- Minimal allocations
- Optimized for read-heavy workloads

[1.0.0]: https://github.com/pedram/toggo/releases/tag/v1.0.0


# Contributing to Toggo

Thank you for your interest in contributing to Toggo! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/toggo.git`
3. Create a feature branch: `git checkout -b feature/my-new-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit your changes: `git commit -am 'Add new feature'`
7. Push to the branch: `git push origin feature/my-new-feature`
8. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Building

```bash
go build ./...
```

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests with race detector:
```bash
go test -race ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Add comments for exported functions, types, and constants
- Write descriptive commit messages

## Testing Guidelines

- Write tests for all new features
- Maintain or improve code coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate
- Ensure all tests pass before submitting PR

Example test structure:
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Feature(tt.input)
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

## Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Update documentation for API changes
- Add tests for new functionality
- Ensure all tests pass
- Update CHANGELOG.md (if applicable)
- Reference related issues in the PR description

## Commit Message Format

Use clear, descriptive commit messages:

```
<type>: <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

Example:
```
feat: add regex operator support

Add support for regex pattern matching in conditions.
This allows users to match attributes against regular expressions.

Closes #123
```

## Adding New Features

When adding new features:

1. Check if an issue exists, or create one
2. Discuss the feature in the issue before implementing
3. Write tests first (TDD approach recommended)
4. Implement the feature
5. Update documentation
6. Add examples if applicable

## Project Structure

```
toggo/
├── toggo.go              # Main package file
├── context.go            # Context type
├── flag.go              # Flag and Variant types
├── condition.go         # Condition type
├── operator.go          # Operator constants
├── store.go             # Store implementation
├── evaluator.go         # Condition evaluation
├── rollout.go           # Rollout strategy
├── errors.go            # Error definitions
├── internal/            # Internal packages
│   └── hash/           # Hashing implementation
├── loader/             # Configuration loaders
├── examples/           # Usage examples
└── testdata/          # Test fixtures
```

## Internal Packages

Code in `internal/` is not part of the public API and can change without notice. Only use exported types and functions from the main `toggo` package.

## Documentation

- Update README.md for user-facing changes
- Add inline documentation for exported symbols
- Include usage examples
- Update package documentation in toggo.go

## Reporting Issues

When reporting issues, please include:

- Go version (`go version`)
- Operating system
- Minimal code to reproduce the issue
- Expected vs actual behavior
- Stack trace (if applicable)

## Feature Requests

Feature requests are welcome! Please:

- Search existing issues first
- Provide clear use cases
- Explain why the feature would be useful
- Be open to discussion and alternatives

## Code Review Process

1. All PRs require review from maintainers
2. Address review comments promptly
3. Keep discussions respectful and constructive
4. Be patient - reviews may take time

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to:
- Open an issue for questions
- Start a discussion
- Reach out to maintainers

Thank you for contributing to Toggo! 🎉


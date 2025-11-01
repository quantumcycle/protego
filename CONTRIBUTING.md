# Contributing to Protego

Thank you for your interest in contributing to Protego! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and inclusive. We're all here to build great software together.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear, descriptive title
- Steps to reproduce the behavior
- Expected vs actual behavior
- Code samples if applicable
- Go version and OS

### Suggesting Enhancements

Enhancement suggestions are welcome! Please create an issue with:
- A clear description of the feature
- Why this feature would be useful
- Example code showing how it would be used

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the coding standards below
3. **Add tests** for any new functionality
4. **Ensure tests pass** by running `go test ./...`
5. **Run the linter** with `golangci-lint run`
6. **Update documentation** including README.md if needed
7. **Submit a pull request** with a clear description of changes

## Development Setup

### Prerequisites

- Go 1.22 or higher
- golangci-lint (optional, for linting)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/quantumcycle/protego.git
cd protego

# Run tests
go test ./validation/... -v
go test ./playground/... -v

# Run with coverage
go test ./validation/... -cover
go test ./playground/... -cover

# Run linter (if installed)
golangci-lint run
```

## Coding Standards

### General Guidelines

- Follow standard Go conventions and idioms
- Use meaningful variable and function names
- Keep functions small and focused
- Avoid premature optimization

### Testing

- Write BDD-style tests using Gomega
- Test both success and failure cases
- Test edge cases and boundaries
- Aim for >95% code coverage
- Use descriptive test names

Example:
```go
func TestValidatorName(t *testing.T) {
    t.Run("passes when condition met", func(t *testing.T) {
        g := NewWithT(t)
        err := validation.Validate(value, validator)
        g.Expect(err).To(BeNil())
    })

    t.Run("fails when condition not met", func(t *testing.T) {
        g := NewWithT(t)
        err := validation.Validate(badValue, validator)
        g.Expect(err).ToNot(BeNil())
    })
}
```

### Documentation

- Add godoc comments for all exported functions
- Include usage examples in comments
- Update README.md for new features
- Add examples to example_test.go

Example:
```go
// ValidatorName validates that a value meets specific criteria.
// It returns an error if validation fails.
//
// Example:
//
//	err := validation.Validate(value, validation.ValidatorName(param))
func ValidatorName(param Type) Validator[T] {
    return func(v T) error {
        // implementation
    }
}
```

### Package Organization

- **validation/**: Core validators with zero external dependencies
- **playground/**: Validators using go-playground/validator

Keep the validation package dependency-free!

## Adding New Validators

### Core Validators (validation package)

1. Add the validator function to the appropriate file:
   - `string.go` - String validators
   - `numeric.go` - Number validators
   - `collection.go` - Slice/map validators
   - `date.go` - Date/time validators
   - `optional.go` - Pointer/optional validators
   - `helpers.go` - Combinator validators

2. Write comprehensive tests
3. Add examples to `example_test.go`
4. Document in README.md
5. Ensure zero external dependencies

### Playground Validators

1. Add pre-built validators to `playground/playground.go`
2. Use `FromTag[T](tag)` to wrap go-playground validators
3. Group validators logically (network, encoding, etc.)
4. Add tests to `playground_test.go`
5. Document in README.md

## Commit Messages

Use clear, descriptive commit messages:

```
Add FloatValidator for map validation

- Adds FloatValidator to convert float64 validators for use with any type
- Handles float64, float32, int, and int64 conversions
- Includes comprehensive tests
- Updates documentation
```

## Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR

## Questions?

Feel free to open an issue for any questions about contributing!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

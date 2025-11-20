package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// MinLength validates that a string has at least the specified minimum length.
//
// Example:
//
//	validation.Validate(name, validation.MinLength(3))
func MinLength(minimum int) Validator[string] {
	return func(v string) error {
		if len(v) < minimum {
			return NewValidationError(fmt.Sprintf("must be at least %d characters", minimum))
		}
		return nil
	}
}

// MaxLength validates that a string has at most the specified maximum length.
//
// Example:
//
//	validation.Validate(name, validation.MaxLength(50))
func MaxLength(maximum int) Validator[string] {
	return func(v string) error {
		if len(v) > maximum {
			return NewValidationError(fmt.Sprintf("must be at most %d characters", maximum))
		}
		return nil
	}
}

// Length validates that a string length is between minimum and maximum (inclusive).
//
// Example:
//
//	validation.Validate(username, validation.Length(3, 20))
func Length(minimum, maximum int) Validator[string] {
	return func(v string) error {
		length := len(v)
		if length < minimum || length > maximum {
			return NewValidationError(fmt.Sprintf("must be between %d and %d characters", minimum, maximum))
		}
		return nil
	}
}

// IsInt validates that a string represents a valid integer.
//
// Example:
//
//	validation.Validate(idString, validation.IsInt())
func IsInt() Validator[string] {
	return func(v string) error {
		if _, err := strconv.Atoi(v); err != nil {
			return NewValidationError("must be a valid integer")
		}
		return nil
	}
}

// MatchesPattern validates that a string matches the given regular expression pattern.
//
// Example:
//
//	validation.Validate(code, validation.MatchesPattern(`^[A-Z]{3}-\d{4}$`))
func MatchesPattern(pattern string) Validator[string] {
	regex := regexp.MustCompile(pattern)
	return func(v string) error {
		if !regex.MatchString(v) {
			return NewValidationError(fmt.Sprintf("must match pattern %q", pattern))
		}
		return nil
	}
}

// StartsWith validates that a string starts with the specified prefix.
//
// Example:
//
//	validation.Validate(packageName, validation.StartsWith("com."))
func StartsWith(prefix string) Validator[string] {
	return func(v string) error {
		if !strings.HasPrefix(v, prefix) {
			return NewValidationError(fmt.Sprintf("must start with %q", prefix))
		}
		return nil
	}
}

// EndsWith validates that a string ends with the specified suffix.
//
// Example:
//
//	validation.Validate(filename, validation.EndsWith(".json"))
func EndsWith(suffix string) Validator[string] {
	return func(v string) error {
		if !strings.HasSuffix(v, suffix) {
			return NewValidationError(fmt.Sprintf("must end with %q", suffix))
		}
		return nil
	}
}

// Contains validates that a string contains the specified substring.
//
// Example:
//
//	validation.Validate(text, validation.Contains("important"))
func Contains(substring string) Validator[string] {
	return func(v string) error {
		if !strings.Contains(v, substring) {
			return NewValidationError(fmt.Sprintf("must contain %q", substring))
		}
		return nil
	}
}

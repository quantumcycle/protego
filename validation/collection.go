package validation

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// In validates that a value is in the allowed list.
// If caseInsensitive is true, string values are compared case-insensitively.
//
// Example:
//
//	validation.Validate(status, validation.In(false, "ACTIVE", "INACTIVE", "PENDING"))
//	validation.Validate(mode, validation.In(true, "single", "multiple")) // case-insensitive
func In[T comparable](caseInsensitive bool, allowed ...T) Validator[T] {
	return func(v T) error {
		if caseInsensitive {
			// For strings, do case-insensitive comparison
			vs := strings.ToLower(fmt.Sprint(v))
			for _, a := range allowed {
				if strings.ToLower(fmt.Sprint(a)) == vs {
					return nil
				}
			}
		} else if slices.Contains(allowed, v) {
			return nil
		}
		return fmt.Errorf("must be one of: %v", allowed)
	}
}

// InSlice validates that a value is in the allowed slice.
// This is similar to In but takes a slice instead of variadic arguments.
//
// Example:
//
//	allowedCurrencies := []string{"USD", "EUR", "GBP"}
//	validation.Validate(currency, validation.InSlice(true, allowedCurrencies))
func InSlice[T comparable](caseInsensitive bool, allowed []T) Validator[T] {
	return In(caseInsensitive, allowed...)
}

// NotIn validates that a value is NOT in the forbidden list.
//
// Example:
//
//	validation.Validate(username, validation.NotIn(false, "admin", "root", "system"))
func NotIn[T comparable](caseInsensitive bool, forbidden ...T) Validator[T] {
	return func(v T) error {
		if caseInsensitive {
			vs := strings.ToLower(fmt.Sprint(v))
			for _, f := range forbidden {
				if strings.ToLower(fmt.Sprint(f)) == vs {
					return fmt.Errorf("cannot be one of: %v", forbidden)
				}
			}
		} else if slices.Contains(forbidden, v) {
			return fmt.Errorf("cannot be one of: %v", forbidden)
		}
		return nil
	}
}

// Each validates each element in a slice using the provided element validator.
// All errors are collected and returned as a joined error.
//
// Example:
//
//	validation.Validate(emails, validation.Each(validation.IsEmail()))
//	validation.Validate(ages, validation.Each(validation.Range(0, 120)))
func Each[T any](elementValidator Validator[T]) Validator[[]T] {
	return func(values []T) error {
		var errs []error
		for i, v := range values {
			if err := elementValidator(v); err != nil {
				errs = append(errs, fmt.Errorf("index %d: %w", i, err))
			}
		}
		return errors.Join(errs...)
	}
}

// NotEmpty validates that a slice is not empty.
//
// Example:
//
//	validation.Validate(tags, validation.NotEmpty[string]())
func NotEmpty[T any]() Validator[[]T] {
	return func(values []T) error {
		if len(values) == 0 {
			return fmt.Errorf("cannot be empty")
		}
		return nil
	}
}

// MinItems validates that a slice has at least the specified minimum number of items.
//
// Example:
//
//	validation.Validate(tags, validation.MinItems[string](1))
func MinItems[T any](minimum int) Validator[[]T] {
	return func(values []T) error {
		if len(values) < minimum {
			return fmt.Errorf("must have at least %d items", minimum)
		}
		return nil
	}
}

// MaxItems validates that a slice has at most the specified maximum number of items.
//
// Example:
//
//	validation.Validate(tags, validation.MaxItems[string](10))
func MaxItems[T any](maximum int) Validator[[]T] {
	return func(values []T) error {
		if len(values) > maximum {
			return fmt.Errorf("must have at most %d items", maximum)
		}
		return nil
	}
}

// UniqueItems validates that all items in a slice are unique.
//
// Example:
//
//	validation.Validate(ids, validation.UniqueItems[string]())
func UniqueItems[T comparable]() Validator[[]T] {
	return func(values []T) error {
		seen := make(map[T]bool)
		for i, v := range values {
			if seen[v] {
				return fmt.Errorf("duplicate item at index %d: %v", i, v)
			}
			seen[v] = true
		}
		return nil
	}
}

// MapKeyRule represents a validation rule for a specific key in a map.
type MapKeyRule[V any] struct {
	key        string
	required   bool
	validators []Validator[V]
}

// MapKey creates a validation rule for a map key.
// If required is true, the key must exist in the map.
//
// Example:
//
//	validation.MapKey("name", true, validation.Required[string](), validation.MinLength(3))
func MapKey[V any](key string, required bool, validators ...Validator[V]) MapKeyRule[V] {
	return MapKeyRule[V]{
		key:        key,
		required:   required,
		validators: validators,
	}
}

// ValidateStringMap validates a map[string]string with the specified rules.
// If allowExtra is false, any keys not defined in rules will cause an error.
//
// Example:
//
//	err := validation.ValidateStringMap(
//	    config,
//	    true, // allow extra keys
//	    validation.MapKey("host", true, validation.Required[string]()),
//	    validation.MapKey("port", true, validation.IsInt()),
//	)
func ValidateStringMap(m map[string]string, allowExtra bool, rules ...MapKeyRule[string]) error {
	validated := make(map[string]bool)

	for _, rule := range rules {
		validated[rule.key] = true

		value, exists := m[rule.key]
		if !exists && rule.required {
			return fmt.Errorf("key %q is required", rule.key)
		}

		if exists {
			for _, validator := range rule.validators {
				if err := validator(value); err != nil {
					return fmt.Errorf("key %q: %w", rule.key, err)
				}
			}
		}
	}

	// Check for extra keys if not allowed
	if !allowExtra {
		for key := range m {
			if !validated[key] {
				return fmt.Errorf("key %q not expected", key)
			}
		}
	}

	return nil
}

// ValidateAnyMap validates a map[string]any (JSON-style map) with the specified rules.
// If allowExtra is false, any keys not defined in rules will cause an error.
//
// Example:
//
//	err := validation.ValidateAnyMap(
//	    jsonData,
//	    true, // allow extra keys
//	    validation.MapKey("name", true, validation.StringValidator(validation.Required[string]())),
//	    validation.MapKey("age", false, validation.AnyValidator(func(v any) error {
//	        age, ok := v.(float64)
//	        if !ok { return fmt.Errorf("must be a number") }
//	        return validation.Validate(int(age), validation.Range(0, 120))
//	    })),
//	)
func ValidateAnyMap(m map[string]any, allowExtra bool, rules ...MapKeyRule[any]) error {
	validated := make(map[string]bool)

	for _, rule := range rules {
		validated[rule.key] = true

		value, exists := m[rule.key]
		if !exists && rule.required {
			return fmt.Errorf("key %q is required", rule.key)
		}

		if exists {
			for _, validator := range rule.validators {
				if err := validator(value); err != nil {
					return fmt.Errorf("key %q: %w", rule.key, err)
				}
			}
		}
	}

	// Check for extra keys if not allowed
	if !allowExtra {
		for key := range m {
			if !validated[key] {
				return fmt.Errorf("key %q not expected", key)
			}
		}
	}

	return nil
}

// StringValidator converts a string validator to work with any type by first asserting it's a string.
// This is useful for ValidateAnyMap when you know a value should be a string.
//
// Example:
//
//	validation.MapKey("email", true, validation.StringValidator(validation.IsEmail()))
func StringValidator(validator Validator[string]) Validator[any] {
	return func(v any) error {
		str, ok := v.(string)
		if !ok {
			return fmt.Errorf("must be a string")
		}
		return validator(str)
	}
}

// IntValidator converts an int validator to work with any type by first asserting it's a number.
// This is useful for ValidateAnyMap when you know a value should be an int.
//
// Example:
//
//	validation.MapKey("age", true, validation.IntValidator(validation.Range(0, 120)))
func IntValidator(validator Validator[int]) Validator[any] {
	return func(v any) error {
		// JSON numbers are float64
		switch val := v.(type) {
		case int:
			return validator(val)
		case float64:
			return validator(int(val))
		case int64:
			return validator(int(val))
		default:
			return fmt.Errorf("must be a number")
		}
	}
}

// FloatValidator converts a float64 validator to work with any type by first asserting it's a number.
// This is useful for ValidateAnyMap when you know a value should be a float.
//
// Example:
//
//	validation.MapKey("price", true, validation.FloatValidator(validation.Min(0.0)))
func FloatValidator(validator Validator[float64]) Validator[any] {
	return func(v any) error {
		// JSON numbers are float64
		switch val := v.(type) {
		case float64:
			return validator(val)
		case float32:
			return validator(float64(val))
		case int:
			return validator(float64(val))
		case int64:
			return validator(float64(val))
		default:
			return fmt.Errorf("must be a number")
		}
	}
}

// BoolValidator converts a bool validator to work with any type by first asserting it's a boolean.
// This is useful for ValidateAnyMap when you know a value should be a bool.
//
// Example:
//
//	validation.MapKey("active", true, validation.BoolValidator(validation.Custom(func(v bool) error {
//	    if !v { return fmt.Errorf("must be active") }
//	    return nil
//	})))
func BoolValidator(validator Validator[bool]) Validator[any] {
	return func(v any) error {
		val, ok := v.(bool)
		if !ok {
			return fmt.Errorf("must be a boolean")
		}
		return validator(val)
	}
}

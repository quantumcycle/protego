package validation

import "fmt"

// Required validates that a value is not the zero value for its type.
// For strings, this means not empty. For numbers, this means not zero.
// For pointers, this means not nil.
//
// Example:
//
//	validation.Validate(name, validation.Required[string]())
//	validation.Validate(age, validation.Required[int]())
func Required[T comparable]() Validator[T] {
	return func(v T) error {
		var zero T
		if v == zero {
			return fmt.Errorf("required")
		}
		return nil
	}
}

// RequiredIf validates that a value is not the zero value if the condition is true.
// If the condition is false, validation passes regardless of the value.
//
// Example:
//
//	validation.Validate(input.ShippingAddress,
//	    validation.RequiredIf[string](input.RequiresShipping),
//	)
func RequiredIf[T comparable](condition bool) Validator[T] {
	return func(v T) error {
		if !condition {
			return nil
		}
		var zero T
		if v == zero {
			return fmt.Errorf("required")
		}
		return nil
	}
}

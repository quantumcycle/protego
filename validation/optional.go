package validation

import "fmt"

// NilOrNotEmpty validates that a pointer to string is either nil or not empty.
// This is useful for optional fields that, if provided, must not be empty.
//
// Example:
//
//	type UpdateUserInput struct {
//	    Name *string
//	}
//
//	validation.Validate(input.Name, validation.NilOrNotEmpty())
func NilOrNotEmpty() Validator[*string] {
	return func(v *string) error {
		if v == nil {
			return nil // Nil is okay
		}
		if *v == "" {
			return fmt.Errorf("cannot be empty string (must be nil or non-empty)")
		}
		return nil
	}
}

// NilOr wraps a validator to make it work with optional (pointer) fields.
// If the pointer is nil, validation passes. If not nil, the wrapped validator is applied.
//
// Example:
//
//	type UpdateUserInput struct {
//	    Email    *string
//	    Age      *int
//	    Priority *int
//	}
//
//	func (input UpdateUserInput) Validate() error {
//	    return errors.Join(
//	        validation.Validate(input.Email,
//	            validation.NilOr(validation.IsEmail()),
//	        ),
//	        validation.Validate(input.Age,
//	            validation.NilOr(validation.Range(0, 120)),
//	        ),
//	        validation.Validate(input.Priority,
//	            validation.NilOr(validation.In(false, 1, 2, 3)),
//	        ),
//	    )
//	}
func NilOr[T any](validator Validator[T]) Validator[*T] {
	return func(v *T) error {
		if v == nil {
			return nil // Nil is okay
		}
		return validator(*v) // Validate the dereferenced value
	}
}

// NotNil validates that a pointer is not nil.
//
// Example:
//
//	validation.Validate(input.RequiredField, validation.NotNil[string]())
func NotNil[T any]() Validator[*T] {
	return func(v *T) error {
		if v == nil {
			return fmt.Errorf("cannot be nil")
		}
		return nil
	}
}

// OptionalWith validates an optional field with multiple validators if it's not nil.
// This is a convenience function that combines NilOr with multiple validators.
//
// Example:
//
//	validation.Validate(input.Email,
//	    validation.OptionalWith(
//	        validation.IsEmail(),
//	        validation.MinLength(5),
//	    ),
//	)
func OptionalWith[T any](validators ...Validator[T]) Validator[*T] {
	return func(v *T) error {
		if v == nil {
			return nil
		}
		return Validate(*v, validators...)
	}
}

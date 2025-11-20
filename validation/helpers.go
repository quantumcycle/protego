package validation

import (
	"errors"
	"fmt"
)

// WithMessage wraps a validator and replaces its error message with a custom message.
// This is useful when you want to provide more specific or user-friendly error messages.
//
// Example:
//
//	validation.Validate(input.Count,
//	    validation.WithMessage(
//	        validation.Required[int](),
//	        "count cannot be zero",
//	    ),
//	)
func WithMessage[T any](validator Validator[T], message string) Validator[T] {
	return func(v T) error {
		if err := validator(v); err != nil {
			return NewValidationError(message)
		}
		return nil
	}
}

// And combines multiple validators - all must pass for validation to succeed.
// This is the default behavior of Validate, but And can be useful for composition.
//
// Example:
//
//	passwordValidator := validation.And(
//	    validation.MinLength(8),
//	    validation.MatchesPattern(`[A-Z]`), // must have uppercase
//	    validation.MatchesPattern(`[0-9]`), // must have number
//	)
//	validation.Validate(password, passwordValidator)
func And[T any](validators ...Validator[T]) Validator[T] {
	return func(v T) error {
		for _, validator := range validators {
			if err := validator(v); err != nil {
				return err
			}
		}
		return nil
	}
}

// Or combines multiple validators - at least one must pass for validation to succeed.
// If all validators fail, a combined error is returned.
//
// Example:
//
//	// Accept either email or phone number
//	validation.Validate(contact,
//	    validation.Or(
//	        validation.IsEmail(),
//	        validation.MatchesPattern(`^\+\d{10,}$`), // phone pattern
//	    ),
//	)
func Or[T any](validators ...Validator[T]) Validator[T] {
	return func(v T) error {
		var errs []error
		for _, validator := range validators {
			if err := validator(v); err == nil {
				return nil // One passed, we're good
			} else {
				errs = append(errs, err)
			}
		}
		if len(errs) == 1 {
			return errs[0]
		}
		return WrapError(fmt.Errorf("all validators failed: %w", errors.Join(errs...)))
	}
}

// Not inverts a validator - it passes if the validator fails, and vice versa.
//
// Example:
//
//	validation.Validate(username,
//	    validation.Not(validation.In(false, "admin", "root")), // must NOT be admin or root
//	)
func Not[T any](validator Validator[T]) Validator[T] {
	return func(v T) error {
		if err := validator(v); err != nil {
			return nil // Validator failed, so Not passes
		}
		return NewValidationError("validation should have failed but passed")
	}
}

// When applies a validator only if the condition is true.
// This is useful for conditional validation.
//
// Example:
//
//	validation.Validate(input.ShippingAddress,
//	    validation.When(
//	        input.RequiresShipping,
//	        validation.Required[string](),
//	    ),
//	)
func When[T any](condition bool, validator Validator[T]) Validator[T] {
	return func(v T) error {
		if !condition {
			return nil // Condition not met, skip validation
		}
		return validator(v)
	}
}

// Unless applies a validator only if the condition is false.
// This is the opposite of When.
//
// Example:
//
//	validation.Validate(input.Reason,
//	    validation.Unless(
//	        input.IsApproved,
//	        validation.Required[string](), // reason required if not approved
//	    ),
//	)
func Unless[T any](condition bool, validator Validator[T]) Validator[T] {
	return When(!condition, validator)
}

// Custom creates a custom validator from a function.
// This is useful for inline validators or wrapping complex validation logic.
//
// Example:
//
//	validation.Validate(username,
//	    validation.Custom(func(v string) error {
//	        if containsProfanity(v) {
//	            return fmt.Errorf("contains inappropriate content")
//	        }
//	        return nil
//	    }),
//	)
func Custom[T any](fn func(T) error) Validator[T] {
	return fn
}

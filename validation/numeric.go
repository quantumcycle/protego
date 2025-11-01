package validation

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Min validates that a value is greater than or equal to the specified minimum.
//
// Example:
//
//	validation.Validate(age, validation.Min(18))
//	validation.Validate(price, validation.Min(0.0))
func Min[T constraints.Ordered](minimum T) Validator[T] {
	return func(v T) error {
		if v < minimum {
			return fmt.Errorf("must be at least %v", minimum)
		}
		return nil
	}
}

// Max validates that a value is less than or equal to the specified maximum.
//
// Example:
//
//	validation.Validate(age, validation.Max(120))
//	validation.Validate(discount, validation.Max(1.0))
func Max[T constraints.Ordered](maximum T) Validator[T] {
	return func(v T) error {
		if v > maximum {
			return fmt.Errorf("must be at most %v", maximum)
		}
		return nil
	}
}

// Range validates that a value is between minimum and maximum (inclusive).
//
// Example:
//
//	validation.Validate(age, validation.Range(18, 120))
//	validation.Validate(rating, validation.Range(1, 5))
func Range[T constraints.Ordered](minimum, maximum T) Validator[T] {
	return func(v T) error {
		if v < minimum || v > maximum {
			return fmt.Errorf("must be between %v and %v", minimum, maximum)
		}
		return nil
	}
}

// GreaterThan validates that a value is strictly greater than the specified value.
//
// Example:
//
//	validation.Validate(count, validation.GreaterThan(0))
func GreaterThan[T constraints.Ordered](threshold T) Validator[T] {
	return func(v T) error {
		if v <= threshold {
			return fmt.Errorf("must be greater than %v", threshold)
		}
		return nil
	}
}

// LessThan validates that a value is strictly less than the specified value.
//
// Example:
//
//	validation.Validate(percentage, validation.LessThan(100.0))
func LessThan[T constraints.Ordered](threshold T) Validator[T] {
	return func(v T) error {
		if v >= threshold {
			return fmt.Errorf("must be less than %v", threshold)
		}
		return nil
	}
}

// Positive validates that a numeric value is greater than zero.
//
// Example:
//
//	validation.Validate(amount, validation.Positive[int]())
//	validation.Validate(price, validation.Positive[float64]())
func Positive[T constraints.Ordered]() Validator[T] {
	return func(v T) error {
		var zero T
		if v <= zero {
			return fmt.Errorf("must be positive")
		}
		return nil
	}
}

// NonNegative validates that a numeric value is greater than or equal to zero.
//
// Example:
//
//	validation.Validate(count, validation.NonNegative[int]())
func NonNegative[T constraints.Ordered]() Validator[T] {
	return func(v T) error {
		var zero T
		if v < zero {
			return fmt.Errorf("must be non-negative")
		}
		return nil
	}
}

// Negative validates that a numeric value is less than zero.
//
// Example:
//
//	validation.Validate(debit, validation.Negative[float64]())
func Negative[T constraints.Ordered]() Validator[T] {
	return func(v T) error {
		var zero T
		if v >= zero {
			return fmt.Errorf("must be negative")
		}
		return nil
	}
}

// MultipleOf validates that a numeric value is a multiple of the specified divisor.
//
// Example:
//
//	validation.Validate(quantity, validation.MultipleOf(5))
func MultipleOf[T constraints.Integer](divisor T) Validator[T] {
	return func(v T) error {
		if v%divisor != 0 {
			return fmt.Errorf("must be a multiple of %v", divisor)
		}
		return nil
	}
}

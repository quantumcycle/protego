// Package validation provides a type-safe validation framework using Go generics.
//
// This package offers compile-time type safety for validators, ensuring that
// validators can only be applied to compatible types. For example, a MinLength
// validator can only be used on strings, not integers.
//
// Example usage:
//
//	func (input CreateUserInput) Validate() error {
//	    return errors.Join(
//	        validation.Validate(input.Email,
//	            validation.Required[string](),
//	            validation.IsEmail(),
//	        ),
//	        validation.Validate(input.Age,
//	            validation.Required[int](),
//	            validation.Range(18, 120),
//	        ),
//	    )
//	}
package validation

// Validator is a generic validation function that validates a value of type T.
// It returns an error if validation fails, or nil if the value is valid.
type Validator[T any] func(T) error

// Validate applies multiple validators to a value and returns the first error encountered.
// All validators are applied in order, and validation stops at the first failure.
//
// Example:
//
//	err := validation.Validate(name,
//	    validation.Required[string](),
//	    validation.MinLength(3),
//	    validation.MaxLength(50),
//	)
func Validate[T any](value T, validators ...Validator[T]) error {
	for _, validator := range validators {
		if err := validator(value); err != nil {
			return err
		}
	}
	return nil
}

// Validatable is an interface for types that have a Validate() method.
// This is useful for nested struct validation.
type Validatable interface {
	Validate() error
}

// ValidateNested checks if the value implements Validatable and calls its Validate() method.
// If the value doesn't implement Validatable, it returns nil (no validation performed).
//
// Example:
//
//	type Address struct { ... }
//	func (a Address) Validate() error { ... }
//
//	type User struct {
//	    Name    string
//	    Address Address
//	}
//
//	func (u User) Validate() error {
//	    return errors.Join(
//	        validation.Validate(u.Name, validation.Required[string]()),
//	        validation.ValidateNested(u.Address), // Calls Address.Validate()
//	    )
//	}
func ValidateNested[T any](value T) error {
	if v, ok := any(value).(Validatable); ok {
		return v.Validate()
	}
	return nil
}

// Nested returns a validator that calls the Validate() method on nested structs.
// This is a wrapper around ValidateNested for use with the Validate function.
//
// Example:
//
//	validation.Validate(user.Address, validation.Nested[Address]())
func Nested[T Validatable]() Validator[T] {
	return func(v T) error {
		return v.Validate()
	}
}

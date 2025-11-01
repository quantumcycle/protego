package validation_test

import (
	"errors"
	"fmt"

	"github.com/quantumcycle/protego/validation"
)

// Example demonstrates basic validation usage
func Example() {
	type CreateUser struct {
		Email    string
		Username string
		Age      int
	}

	input := CreateUser{
		Email:    "test@example.com",
		Username: "john123",
		Age:      25,
	}

	err := errors.Join(
		validation.Validate(input.Email,
			validation.Required[string](),
			validation.Contains("@"),
		),
		validation.Validate(input.Username,
			validation.Required[string](),
			validation.Length(3, 20),
			validation.MatchesPattern(`^[a-zA-Z0-9]+$`),
		),
		validation.Validate(input.Age,
			validation.Required[int](),
			validation.Range(18, 120),
		),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// ExampleOptionalFields demonstrates validation of optional pointer fields
func Example_optionalFields() {
	type UpdateUser struct {
		Email *string
		Age   *int
	}

	email := "test@example.com"
	age := 25
	input := UpdateUser{
		Email: &email,
		Age:   &age,
	}

	err := errors.Join(
		validation.Validate(input.Email,
			validation.NilOr(validation.Contains("@")),
		),
		validation.Validate(input.Age,
			validation.NilOr(validation.Range(0, 120)),
		),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// ExampleSliceValidation demonstrates validation of slices
func Example_sliceValidation() {
	emails := []string{"test1@example.com", "test2@example.com", "test3@example.com"}

	err := validation.Validate(emails,
		validation.NotEmpty[string](),
		validation.MaxItems[string](10),
		validation.Each(validation.Contains("@")),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// ExampleCustomValidator demonstrates creating custom validators
func Example_customValidator() {
	// Define a custom validator
	IsPositiveEven := func() validation.Validator[int] {
		return func(v int) error {
			if v <= 0 {
				return fmt.Errorf("must be positive")
			}
			if v%2 != 0 {
				return fmt.Errorf("must be even")
			}
			return nil
		}
	}

	err := validation.Validate(4, IsPositiveEven())

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// ExampleConditionalValidation demonstrates conditional validation
func Example_conditionalValidation() {
	type Order struct {
		RequiresShipping bool
		ShippingAddress  string
	}

	input := Order{
		RequiresShipping: false,
		ShippingAddress:  "", // Empty, but that's okay since RequiresShipping is false
	}

	err := validation.Validate(input.ShippingAddress,
		validation.When(
			input.RequiresShipping,
			validation.Required[string](),
		),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// Example_mapValidation demonstrates validating map[string]string
func Example_mapValidation() {
	config := map[string]string{
		"host": "localhost",
		"port": "8080",
	}

	err := validation.ValidateStringMap(
		config,
		true, // allow extra keys
		validation.MapKey("host", true, validation.Required[string]()),
		validation.MapKey("port", true, validation.IsInt()),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// Example_anyMapValidation demonstrates validating map[string]any (JSON-style data)
func Example_anyMapValidation() {
	data := map[string]any{
		"name":   "John Doe",
		"age":    30,
		"active": true,
	}

	err := validation.ValidateAnyMap(
		data,
		true, // allow extra keys
		validation.MapKey("name", true, validation.StringValidator(validation.Required[string]())),
		validation.MapKey("age", true, validation.IntValidator(validation.Range(0, 120))),
		validation.MapKey("active", false, validation.BoolValidator(validation.Custom(func(v bool) error {
			return nil // Any bool value is fine
		}))),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// Example_nestedValidation demonstrates validating nested structs
func Example_nestedValidation() {
	type Address struct {
		Street string
		City   string
	}

	type User struct {
		Name    string
		Address Address
	}

	// Address implements Validatable interface
	validateAddress := func(a Address) error {
		return errors.Join(
			validation.Validate(a.Street, validation.Required[string]()),
			validation.Validate(a.City, validation.Required[string]()),
		)
	}

	user := User{
		Name: "John Doe",
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
		},
	}

	err := errors.Join(
		validation.Validate(user.Name, validation.Required[string]()),
		validateAddress(user.Address),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// Example_composedValidators demonstrates using And, Or, and Not
func Example_composedValidators() {
	// Complex password validator: min 8 chars, has uppercase, has number
	passwordValidator := validation.And(
		validation.MinLength(8),
		validation.MatchesPattern(`[A-Z]`), // has uppercase
		validation.MatchesPattern(`[0-9]`), // has number
	)

	// Email or phone validator
	contactValidator := validation.Or(
		validation.Contains("@"),
		validation.MatchesPattern(`^\+\d{10,}$`),
	)

	password := "SecureP@ss123"
	contact := "test@example.com"

	err := errors.Join(
		validation.Validate(password, passwordValidator),
		validation.Validate(contact, contactValidator),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

// Example_dateValidation demonstrates date and time validation
func Example_dateValidation() {
	// RFC3339 datetime
	timestamp := "2024-12-31T23:59:59Z"

	// ISO8601 date only
	birthDate := "1990-05-15"

	err := errors.Join(
		validation.Validate(timestamp, validation.IsRFC3339DateTime()),
		validation.Validate(birthDate, validation.IsISO8601Date()),
	)

	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed")
	}
	// Output: Validation passed
}

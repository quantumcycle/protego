package validation_test

import (
	"errors"
	"testing"

	"github.com/quantumcycle/protego/validation"
)

func TestValidateStruct(t *testing.T) {
	type User struct {
		Name  string
		Email string `json:"email"`
		Age   int
	}

	t.Run("all valid", func(t *testing.T) {
		u := User{Name: "John", Email: "john@example.com", Age: 30}
		err := validation.ValidateStruct(
			validation.Field(&u, &u.Name, validation.Required[string]()),
			validation.Field(&u, &u.Email, validation.Required[string](), validation.Contains("@")),
			validation.Field(&u, &u.Age, validation.Range(18, 120)),
		)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("field error uses json tag as name", func(t *testing.T) {
		u := User{Name: "John", Email: "not-an-email", Age: 30}
		err := validation.ValidateStruct(
			validation.Field(&u, &u.Email, validation.Contains("@")),
		)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error()[:5] != "email" {
			t.Fatalf("expected error prefixed with 'email', got: %v", err)
		}
	})

	t.Run("field error uses struct field name when no json tag", func(t *testing.T) {
		u := User{Name: "", Email: "john@example.com", Age: 30}
		err := validation.ValidateStruct(
			validation.Field(&u, &u.Name, validation.Required[string]()),
		)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error()[:4] != "Name" {
			t.Fatalf("expected error prefixed with 'Name', got: %v", err)
		}
	})

	t.Run("multiple errors joined", func(t *testing.T) {
		u := User{Name: "", Email: "bad", Age: 5}
		err := validation.ValidateStruct(
			validation.Field(&u, &u.Name, validation.Required[string]()),
			validation.Field(&u, &u.Email, validation.Contains("@")),
			validation.Field(&u, &u.Age, validation.Range(18, 120)),
		)
		if err == nil {
			t.Fatal("expected errors")
		}
		var errs interface{ Unwrap() []error }
		if !errors.As(err, &errs) || len(errs.Unwrap()) != 3 {
			t.Fatalf("expected 3 joined errors, got: %v", err)
		}
	})
}

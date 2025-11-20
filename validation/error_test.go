package validation_test

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/quantumcycle/protego/validation"
)

func TestValidationError(t *testing.T) {
	t.Run("NewValidationError creates a ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.NewValidationError("test error")
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(Equal("test error"))
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("WrapError wraps an existing error", func(t *testing.T) {
		g := NewWithT(t)
		originalErr := errors.New("original error")
		wrappedErr := validation.WrapError(originalErr)
		g.Expect(wrappedErr).ToNot(BeNil())
		g.Expect(wrappedErr.Error()).To(Equal("original error"))
		g.Expect(validation.IsValidationError(wrappedErr)).To(BeTrue())
	})

	t.Run("WrapError preserves ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		originalErr := validation.NewValidationError("validation error")
		wrappedErr := validation.WrapError(originalErr)
		g.Expect(wrappedErr).To(Equal(originalErr))
	})

	t.Run("WrapError returns nil for nil error", func(t *testing.T) {
		g := NewWithT(t)
		wrappedErr := validation.WrapError(nil)
		g.Expect(wrappedErr).To(BeNil())
	})

	t.Run("IsValidationError returns false for non-ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := errors.New("regular error")
		g.Expect(validation.IsValidationError(err)).To(BeFalse())
	})

	t.Run("IsValidationError returns false for nil", func(t *testing.T) {
		g := NewWithT(t)
		g.Expect(validation.IsValidationError(nil)).To(BeFalse())
	})

	t.Run("Error unwrapping works", func(t *testing.T) {
		g := NewWithT(t)
		originalErr := errors.New("original error")
		wrappedErr := validation.WrapError(originalErr)
		g.Expect(errors.Unwrap(wrappedErr)).To(Equal(originalErr))
	})

	t.Run("errors.Is works with validation.Error", func(t *testing.T) {
		g := NewWithT(t)
		err1 := validation.NewValidationError("test error")
		err2 := validation.NewValidationError("another error")
		g.Expect(errors.Is(err1, &validation.Error{})).To(BeTrue())
		g.Expect(errors.Is(err2, &validation.Error{})).To(BeTrue())
	})
}

func TestValidatorsReturnValidationError(t *testing.T) {
	t.Run("Required validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.Required[string]())
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("MinLength validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ab", validation.MinLength(3))
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("Range validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(150, validation.Range(0, 120))
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("In validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("invalid", validation.In(false, "valid1", "valid2"))
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("NotEmpty validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{}, validation.NotEmpty[string]())
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("IsRFC3339DateTime validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("invalid-date", validation.IsRFC3339DateTime())
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("WithMessage validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.WithMessage(validation.Required[string](), "custom message"))
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(Equal("custom message"))
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("Or validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("test", validation.Or(
			validation.MinLength(10),
			validation.MaxLength(2),
		))
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("NotNil validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		var nilStr *string
		err := validation.Validate(nilStr, validation.NotNil[string]())
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})
}

func TestErrorMessagesPreserved(t *testing.T) {
	t.Run("Required error message preserved", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.Required[string]())
		g.Expect(err.Error()).To(Equal("required"))
	})

	t.Run("MinLength error message preserved", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ab", validation.MinLength(3))
		g.Expect(err.Error()).To(Equal("must be at least 3 characters"))
	})

	t.Run("Range error message preserved", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(150, validation.Range(0, 120))
		g.Expect(err.Error()).To(Equal("must be between 0 and 120"))
	})
}

package playground_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/quantumcycle/protego/playground"
	"github.com/quantumcycle/protego/validation"
)

func TestPlaygroundValidatorsReturnValidationError(t *testing.T) {
	t.Run("IsEmail validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-an-email", playground.IsEmail)
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("IsUUID4 validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-uuid", playground.IsUUID4)
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("IsURL validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-url", playground.IsURL)
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("IsIPv4 validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-an-ip", playground.IsIPv4)
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("FromTag validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		customValidator := playground.FromTag[string]("uuid")
		err := validation.Validate("not-a-uuid", customValidator)
		g.Expect(err).ToNot(BeNil())
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})

	t.Run("FromTagWithMessage validator returns ValidationError", func(t *testing.T) {
		g := NewWithT(t)
		customValidator := playground.FromTagWithMessage[string]("uuid", "custom error message")
		err := validation.Validate("not-a-uuid", customValidator)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(Equal("custom error message"))
		g.Expect(validation.IsValidationError(err)).To(BeTrue())
	})
}

func TestPlaygroundValidatorsPass(t *testing.T) {
	t.Run("IsEmail validator passes for valid email", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("test@example.com", playground.IsEmail)
		g.Expect(err).To(BeNil())
	})

	t.Run("IsUUID4 validator passes for valid UUID", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("f47ac10b-58cc-4372-a567-0e02b2c3d479", playground.IsUUID4)
		g.Expect(err).To(BeNil())
	})

	t.Run("IsURL validator passes for valid URL", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("https://example.com", playground.IsURL)
		g.Expect(err).To(BeNil())
	})
}

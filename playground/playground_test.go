package playground_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/quantumcycle/protego/playground"
	"github.com/quantumcycle/protego/validation"
)

func TestFromTag(t *testing.T) {

	t.Run("works with uuid4", func(t *testing.T) {
		g := NewWithT(t)
		IsUUID4 := playground.FromTag[string]("uuid4")
		err := validation.Validate("550e8400-e29b-41d4-a716-446655440000", IsUUID4)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid uuid4", func(t *testing.T) {
		g := NewWithT(t)
		IsUUID4 := playground.FromTag[string]("uuid4")
		err := validation.Validate("not-a-uuid", IsUUID4)
		g.Expect(err).NotTo(BeNil())
	})

	t.Run("works with ipv4", func(t *testing.T) {
		g := NewWithT(t)
		IsIPv4 := playground.FromTag[string]("ipv4")
		err := validation.Validate("192.168.1.1", IsIPv4)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid ipv4", func(t *testing.T) {
		g := NewWithT(t)
		IsIPv4 := playground.FromTag[string]("ipv4")
		err := validation.Validate("999.999.999.999", IsIPv4)
		g.Expect(err).NotTo(BeNil())
	})

	t.Run("works with numeric range", func(t *testing.T) {
		g := NewWithT(t)
		PortValidator := playground.FromTag[int]("min=1,max=65535")
		err := validation.Validate(8080, PortValidator)
		g.Expect(err).To(BeNil())
	})

	t.Run("works with custom message", func(t *testing.T) {
		g := NewWithT(t)
		IsUUID4 := playground.FromTagWithMessage[string]("uuid4", "must be a valid UUID v4")
		err := validation.Validate("invalid", IsUUID4)
		g.Expect(err).To(MatchError("must be a valid UUID v4"))
	})
}

func TestIsEmail(t *testing.T) {

	t.Run("passes with valid email", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("test@example.com", playground.IsEmail)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid email", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-an-email", playground.IsEmail)
		g.Expect(err).NotTo(BeNil())
	})
}

func TestIsURL(t *testing.T) {

	t.Run("passes with valid URL", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("https://example.com", playground.IsURL)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid URL", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not a url", playground.IsURL)
		g.Expect(err).NotTo(BeNil())
	})
}

func TestIsSemver(t *testing.T) {

	t.Run("passes with valid semver", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("1.2.3", playground.IsSemver)
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with prerelease", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("1.2.3-alpha", playground.IsSemver)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid semver", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-version", playground.IsSemver)
		g.Expect(err).NotTo(BeNil())
	})
}

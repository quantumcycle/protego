package validation_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/quantumcycle/protego/validation"
)

func TestRequired(t *testing.T) {

	t.Run("string - passes when not empty", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("test", validation.Required[string]())
		g.Expect(err).To(BeNil())
	})

	t.Run("string - fails when empty", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.Required[string]())
		g.Expect(err).To(MatchError("required"))
	})

	t.Run("int - passes when not zero", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(42, validation.Required[int]())
		g.Expect(err).To(BeNil())
	})

	t.Run("int - fails when zero", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(0, validation.Required[int]())
		g.Expect(err).To(MatchError("required"))
	})
}

func TestMinLength(t *testing.T) {

	t.Run("passes when length equals min", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("abc", validation.MinLength(3))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when length exceeds min", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("abcd", validation.MinLength(3))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when length below min", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ab", validation.MinLength(3))
		g.Expect(err).To(MatchError("must be at least 3 characters"))
	})
}

func TestMaxLength(t *testing.T) {

	t.Run("passes when length equals max", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("abc", validation.MaxLength(3))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when length below max", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ab", validation.MaxLength(3))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when length exceeds max", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("abcd", validation.MaxLength(3))
		g.Expect(err).To(MatchError("must be at most 3 characters"))
	})
}

func TestLength(t *testing.T) {

	t.Run("passes when within range", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("test", validation.Length(3, 5))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when below min", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ab", validation.Length(3, 5))
		g.Expect(err).To(MatchError("must be between 3 and 5 characters"))
	})

	t.Run("fails when above max", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("abcdef", validation.Length(3, 5))
		g.Expect(err).To(MatchError("must be between 3 and 5 characters"))
	})
}

func TestMin(t *testing.T) {

	t.Run("int - passes when equal", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.Min(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("int - passes when greater", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10, validation.Min(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("int - fails when less", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(3, validation.Min(5))
		g.Expect(err).To(MatchError("must be at least 5"))
	})

	t.Run("float - passes when equal", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5.5, validation.Min(5.5))
		g.Expect(err).To(BeNil())
	})

	t.Run("float - passes when greater", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10.5, validation.Min(5.5))
		g.Expect(err).To(BeNil())
	})

	t.Run("float - fails when less", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(3.5, validation.Min(5.5))
		g.Expect(err).To(MatchError("must be at least 5.5"))
	})
}

func TestMax(t *testing.T) {
	t.Run("passes when equal", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.Max(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when less", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(3, validation.Max(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when greater", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10, validation.Max(5))
		g.Expect(err).To(MatchError("must be at most 5"))
	})
}

func TestRange(t *testing.T) {

	t.Run("passes when within range", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.Range(1, 10))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes at min boundary", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(1, validation.Range(1, 10))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes at max boundary", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10, validation.Range(1, 10))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when below min", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(0, validation.Range(1, 10))
		g.Expect(err).To(MatchError("must be between 1 and 10"))
	})

	t.Run("fails when above max", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(11, validation.Range(1, 10))
		g.Expect(err).To(MatchError("must be between 1 and 10"))
	})
}

func TestIn(t *testing.T) {

	t.Run("case sensitive - passes when in list", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ACTIVE", validation.In(false, "ACTIVE", "INACTIVE"))
		g.Expect(err).To(BeNil())
	})

	t.Run("case sensitive - fails when not in list", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("PENDING", validation.In(false, "ACTIVE", "INACTIVE"))
		g.Expect(err).To(MatchError("must be one of: [ACTIVE INACTIVE]"))
	})

	t.Run("case sensitive - fails on case mismatch", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("active", validation.In(false, "ACTIVE", "INACTIVE"))
		g.Expect(err).To(MatchError("must be one of: [ACTIVE INACTIVE]"))
	})

	t.Run("case insensitive - passes with different case", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("active", validation.In(true, "ACTIVE", "INACTIVE"))
		g.Expect(err).To(BeNil())
	})

	t.Run("case insensitive - passes with mixed case", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("AcTiVe", validation.In(true, "ACTIVE", "INACTIVE"))
		g.Expect(err).To(BeNil())
	})
}

func TestInSlice(t *testing.T) {

	allowed := []string{
		"USD", "EUR", "GBP"}

	t.Run("passes when in slice", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("USD", validation.InSlice(false, allowed))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when not in slice", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("JPY", validation.InSlice(false, allowed))
		g.Expect(err).To(MatchError("must be one of: [USD EUR GBP]"))
	})
}

func TestEach(t *testing.T) {

	t.Run("passes when all elements valid", func(t *testing.T) {
		g := NewWithT(t)
		values := []string{
			"test1", "test2", "test3"}
		err := validation.Validate(values, validation.Each(validation.MinLength(3)))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when one element invalid", func(t *testing.T) {
		g := NewWithT(t)
		values := []string{
			"test", "ab", "test3"}
		err := validation.Validate(values, validation.Each(validation.MinLength(3)))
		g.Expect(err).To(MatchError(ContainSubstring("index 1")))
	})

	t.Run("collects all errors", func(t *testing.T) {
		g := NewWithT(t)
		values := []string{
			"ab", "cd", "test"}
		err := validation.Validate(values, validation.Each(validation.MinLength(3)))
		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("index 0"))
		g.Expect(err.Error()).To(ContainSubstring("index 1"))
	})
}

func TestNilOrNotEmpty(t *testing.T) {

	t.Run("passes when nil", func(t *testing.T) {
		g := NewWithT(t)
		var str *string
		err := validation.Validate(str, validation.NilOrNotEmpty())
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when not empty", func(t *testing.T) {
		g := NewWithT(t)
		str := "test"
		err := validation.Validate(&str, validation.NilOrNotEmpty())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when empty string", func(t *testing.T) {
		g := NewWithT(t)
		str := ""
		err := validation.Validate(&str, validation.NilOrNotEmpty())
		g.Expect(err).To(MatchError("cannot be empty string (must be nil or non-empty)"))
	})
}

func TestNilOr(t *testing.T) {

	t.Run("passes when nil", func(t *testing.T) {
		g := NewWithT(t)
		var num *int
		err := validation.Validate(num, validation.NilOr(validation.Range(1, 10)))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when valid value", func(t *testing.T) {
		g := NewWithT(t)
		num := 5
		err := validation.Validate(&num, validation.NilOr(validation.Range(1, 10)))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when invalid value", func(t *testing.T) {
		g := NewWithT(t)
		num := 15
		err := validation.Validate(&num, validation.NilOr(validation.Range(1, 10)))
		g.Expect(err).To(MatchError("must be between 1 and 10"))
	})
}

func TestIsInt(t *testing.T) {

	t.Run("passes with valid integer string", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("123", validation.IsInt())
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with negative integer", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("-123", validation.IsInt())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with non-integer", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("abc", validation.IsInt())
		g.Expect(err).To(MatchError("must be a valid integer"))
	})
}

func TestWithMessage(t *testing.T) {

	t.Run("uses custom message on failure", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(0,
			validation.WithMessage(
				validation.Required[int](),
				"count cannot be zero",
			),
		)
		g.Expect(err).To(MatchError("count cannot be zero"))
	})

	t.Run("passes when validation succeeds", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5,
			validation.WithMessage(
				validation.Required[int](),
				"count cannot be zero",
			),
		)
		g.Expect(err).To(BeNil())
	})
}

func TestAnd(t *testing.T) {

	t.Run("passes when all validators pass", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.And(
			validation.MinLength(3),
			validation.MaxLength(10),
		)
		err := validation.Validate("test", validator)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when one validator fails", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.And(
			validation.MinLength(3),
			validation.MaxLength(10),
		)
		err := validation.Validate("ab", validator)
		g.Expect(err).NotTo(BeNil())
	})
}

func TestOr(t *testing.T) {

	t.Run("passes when first validator passes", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.Or(
			validation.Contains("@"),
			validation.StartsWith("http"),
		)
		err := validation.Validate("test@example.com", validator)
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when second validator passes", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.Or(
			validation.Contains("@"),
			validation.StartsWith("http"),
		)
		err := validation.Validate("https://example.com", validator)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when all validators fail", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.Or(
			validation.Contains("@"),
			validation.StartsWith("http"),
		)
		err := validation.Validate("not-valid", validator)
		g.Expect(err).To(MatchError(ContainSubstring("all validators failed")))
	})

	t.Run("returns single error when only one validator fails", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.Or(
			validation.Contains("@"),
		)
		err := validation.Validate("no-at-sign", validator)
		g.Expect(err).To(MatchError(ContainSubstring("must contain")))
		g.Expect(err).NotTo(MatchError(ContainSubstring("all validators failed")))
	})
}

func TestNestedValidation(t *testing.T) {

	type Address struct {
		Street string
		City   string
	}

	// Address implements Validatable
	validateAddress := func(a Address) error {
		return errors.Join(
			validation.Validate(a.Street, validation.Required[string]()),
			validation.Validate(a.City, validation.Required[string]()),
		)
	}

	t.Run("passes with valid nested struct", func(t *testing.T) {
		g := NewWithT(t)
		address := Address{
			Street: "123 Main St",
			City:   "New York",
		}
		err := validateAddress(address)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid nested struct", func(t *testing.T) {
		g := NewWithT(t)
		address := Address{
			Street: "",
			City:   "New York",
		}
		err := validateAddress(address)
		g.Expect(err).NotTo(BeNil())
	})
}

func TestMultipleValidators(t *testing.T) {

	t.Run("applies all validators in order", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("test@example.com",
			validation.Required[string](),
			validation.MinLength(3),
			validation.Contains("@"),
		)
		g.Expect(err).To(BeNil())
	})

	t.Run("stops at first error", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("",
			validation.Required[string](),
			validation.MinLength(3),  // Won't be checked
			validation.Contains("@"), // Won't be checked
		)
		g.Expect(err).To(MatchError("required"))
	})
}

func TestGreaterThan(t *testing.T) {

	t.Run("passes when greater", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10, validation.GreaterThan(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when equal", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.GreaterThan(5))
		g.Expect(err).To(MatchError(ContainSubstring("must be greater than 5")))
	})

	t.Run("fails when less", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(3, validation.GreaterThan(5))
		g.Expect(err).To(MatchError(ContainSubstring("must be greater than 5")))
	})
}

func TestLessThan(t *testing.T) {

	t.Run("passes when less", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(3, validation.LessThan(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when equal", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.LessThan(5))
		g.Expect(err).To(MatchError(ContainSubstring("must be less than 5")))
	})

	t.Run("fails when greater", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10, validation.LessThan(5))
		g.Expect(err).To(MatchError(ContainSubstring("must be less than 5")))
	})
}

func TestPositive(t *testing.T) {

	t.Run("passes when positive", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.Positive[int]())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when zero", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(0, validation.Positive[int]())
		g.Expect(err).To(MatchError(ContainSubstring("must be positive")))
	})

	t.Run("fails when negative", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(-5, validation.Positive[int]())
		g.Expect(err).To(MatchError(ContainSubstring("must be positive")))
	})
}

func TestNonNegative(t *testing.T) {

	t.Run("passes when positive", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.NonNegative[int]())
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when zero", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(0, validation.NonNegative[int]())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when negative", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(-5, validation.NonNegative[int]())
		g.Expect(err).To(MatchError(ContainSubstring("must be non-negative")))
	})
}

func TestNegative(t *testing.T) {

	t.Run("passes when negative", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(-5, validation.Negative[int]())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when zero", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(0, validation.Negative[int]())
		g.Expect(err).To(MatchError(ContainSubstring("must be negative")))
	})

	t.Run("fails when positive", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(5, validation.Negative[int]())
		g.Expect(err).To(MatchError(ContainSubstring("must be negative")))
	})
}

func TestMultipleOf(t *testing.T) {

	t.Run("passes when multiple", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10, validation.MultipleOf(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when zero", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(0, validation.MultipleOf(5))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when not multiple", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(7, validation.MultipleOf(5))
		g.Expect(err).To(MatchError(ContainSubstring("must be a multiple of 5")))
	})
}

func TestNotIn(t *testing.T) {

	t.Run("passes when not in list", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("other", validation.NotIn(false, "admin", "root"))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when in list", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("admin", validation.NotIn(false, "admin", "root"))
		g.Expect(err).To(MatchError(ContainSubstring("cannot be one of")))
	})

	t.Run("case insensitive - fails with different case", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ADMIN", validation.NotIn(true, "admin", "root"))
		g.Expect(err).To(MatchError(ContainSubstring("cannot be one of")))
	})
}

func TestMinItems(t *testing.T) {

	t.Run("passes when meets minimum", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a", "b"}, validation.MinItems[string](2))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when exceeds minimum", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a", "b", "c"}, validation.MinItems[string](2))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when below minimum", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a"}, validation.MinItems[string](2))
		g.Expect(err).To(MatchError(ContainSubstring("must have at least 2 items")))
	})
}

func TestMaxItems(t *testing.T) {

	t.Run("passes when below maximum", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a"}, validation.MaxItems[string](2))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when equals maximum", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a", "b"}, validation.MaxItems[string](2))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when exceeds maximum", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a", "b", "c"}, validation.MaxItems[string](2))
		g.Expect(err).To(MatchError(ContainSubstring("must have at most 2 items")))
	})
}

func TestUniqueItems(t *testing.T) {

	t.Run("passes when all unique", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a", "b", "c"}, validation.UniqueItems[string]())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when duplicate", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate([]string{
			"a", "b", "a"}, validation.UniqueItems[string]())
		g.Expect(err).To(MatchError(ContainSubstring("duplicate item")))
	})
}

func TestStartsWith(t *testing.T) {

	t.Run("passes when starts with prefix", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("hello world", validation.StartsWith("hello"))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when does not start with prefix", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("world hello", validation.StartsWith("hello"))
		g.Expect(err).To(MatchError(ContainSubstring("must start with")))
	})
}

func TestEndsWith(t *testing.T) {

	t.Run("passes when ends with suffix", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("hello world", validation.EndsWith("world"))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when does not end with suffix", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("world hello", validation.EndsWith("world"))
		g.Expect(err).To(MatchError(ContainSubstring("must end with")))
	})
}

func TestContains(t *testing.T) {

	t.Run("passes when contains substring", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("hello world", validation.Contains("lo wo"))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when does not contain substring", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("hello world", validation.Contains("xyz"))
		g.Expect(err).To(MatchError(ContainSubstring("must contain")))
	})
}

func TestMatchesPattern(t *testing.T) {

	t.Run("passes when matches pattern", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("ABC-1234", validation.MatchesPattern(`^[A-Z]{3}-\d{4}$`))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when does not match pattern", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("abc-1234", validation.MatchesPattern(`^[A-Z]{3}-\d{4}$`))
		g.Expect(err).To(MatchError(ContainSubstring("must match pattern")))
	})
}

func TestIsRFC3339DateTime(t *testing.T) {

	t.Run("passes with valid RFC3339", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-15T10:30:00Z", validation.IsRFC3339DateTime())
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with RFC3339 with timezone", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-15T10:30:00-05:00", validation.IsRFC3339DateTime())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid RFC3339", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-15", validation.IsRFC3339DateTime())
		g.Expect(err).To(MatchError(ContainSubstring("must be a valid RFC3339 date-time")))
	})

	t.Run("fails with invalid format", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-date", validation.IsRFC3339DateTime())
		g.Expect(err).ToNot(BeNil())
	})
}

func TestIsISO8601Date(t *testing.T) {

	t.Run("passes with valid ISO8601 date", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-15", validation.IsISO8601Date())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with datetime instead of date", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-15T10:30:00Z", validation.IsISO8601Date())
		g.Expect(err).To(MatchError(ContainSubstring("must be a valid ISO8601 date")))
	})

	t.Run("fails with invalid date", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-13-45", validation.IsISO8601Date())
		g.Expect(err).ToNot(BeNil())
	})

	t.Run("fails with non-date string", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-date", validation.IsISO8601Date())
		g.Expect(err).To(MatchError(ContainSubstring("must be a valid ISO8601 date")))
	})
}

func TestIsDateFormat(t *testing.T) {

	t.Run("passes with matching format", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-15 10:30:00", validation.IsDateFormat("2006-01-02 15:04:05"))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with non-matching format", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-15", validation.IsDateFormat("2006-01-02 15:04:05"))
		g.Expect(err).ToNot(BeNil())
	})
}

func TestRequiredIf(t *testing.T) {

	t.Run("passes when condition false", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.RequiredIf[string](false))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when condition true and value provided", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("value", validation.RequiredIf[string](true))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when condition true and value empty", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.RequiredIf[string](true))
		g.Expect(err).To(MatchError(ContainSubstring("required")))
	})
}

func TestNotNil(t *testing.T) {

	t.Run("passes when not nil", func(t *testing.T) {
		g := NewWithT(t)
		value := "test"
		err := validation.Validate(&value, validation.NotNil[string]())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when nil", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate((*string)(nil), validation.NotNil[string]())
		g.Expect(err).To(MatchError(ContainSubstring("cannot be nil")))
	})
}

func TestOptionalWith(t *testing.T) {

	t.Run("passes when nil", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate((*string)(nil), validation.OptionalWith(validation.MinLength(5)))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes when valid", func(t *testing.T) {
		g := NewWithT(t)
		value := "hello"
		err := validation.Validate(&value, validation.OptionalWith(validation.MinLength(5)))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when invalid", func(t *testing.T) {
		g := NewWithT(t)
		value := "hi"
		err := validation.Validate(&value, validation.OptionalWith(validation.MinLength(5)))
		g.Expect(err).ToNot(BeNil())
	})
}

func TestNot(t *testing.T) {

	t.Run("passes when validator fails", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("noemail", validation.Not(validation.Contains("@")))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when validator passes", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("test@example.com", validation.Not(validation.Contains("@")))
		g.Expect(err).ToNot(BeNil())
	})
}

func TestWhen(t *testing.T) {

	t.Run("applies validator when condition true", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.When(true, validation.Required[string]()))
		g.Expect(err).ToNot(BeNil())
	})

	t.Run("skips validator when condition false", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.When(false, validation.Required[string]()))
		g.Expect(err).To(BeNil())
	})
}

func TestUnless(t *testing.T) {

	t.Run("skips validator when condition true", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.Unless(true, validation.Required[string]()))
		g.Expect(err).To(BeNil())
	})

	t.Run("applies validator when condition false", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("", validation.Unless(false, validation.Required[string]()))
		g.Expect(err).ToNot(BeNil())
	})
}

func TestCustom(t *testing.T) {

	t.Run("passes with custom validator", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(10, validation.Custom(func(v int) error {
			if v%2 == 0 {
				return nil
			}
			return fmt.Errorf("must be even")
		}))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with custom validator", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate(11, validation.Custom(func(v int) error {
			if v%2 == 0 {
				return nil
			}
			return fmt.Errorf("must be even")
		}))
		g.Expect(err).To(MatchError("must be even"))
	})
}

func TestValidateStringMap(t *testing.T) {

	t.Run("passes with valid map", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]string{
			"name": "John", "age": "30"}
		err := validation.ValidateStringMap(m, true,
			validation.MapKey("name", true, validation.Required[string]()),
			validation.MapKey("age", true, validation.IsInt()),
		)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when required key missing", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]string{
			"age": "30"}
		err := validation.ValidateStringMap(m, true,
			validation.MapKey("name", true, validation.Required[string]()),
		)
		g.Expect(err).To(MatchError(ContainSubstring("is required")))
	})

	t.Run("fails with extra keys when not allowed", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]string{
			"name": "John", "extra": "value"}
		err := validation.ValidateStringMap(m, false,
			validation.MapKey("name", true, validation.Required[string]()),
		)
		g.Expect(err).To(MatchError(ContainSubstring("not expected")))
	})
}

func TestIsFutureDate(t *testing.T) {

	t.Run("passes with future date", func(t *testing.T) {
		g := NewWithT(t)
		futureDate := "2099-12-31T23:59:59Z"
		err := validation.Validate(futureDate, validation.IsFutureDate())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with past date", func(t *testing.T) {
		g := NewWithT(t)
		pastDate := "2020-01-01T00:00:00Z"
		err := validation.Validate(pastDate, validation.IsFutureDate())
		g.Expect(err).To(MatchError(ContainSubstring("must be a future date")))
	})

	t.Run("fails with invalid format", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-date", validation.IsFutureDate())
		g.Expect(err).To(MatchError(ContainSubstring("invalid date format")))
	})
}

func TestIsPastDate(t *testing.T) {

	t.Run("passes with past date", func(t *testing.T) {
		g := NewWithT(t)
		pastDate := "2020-01-01T00:00:00Z"
		err := validation.Validate(pastDate, validation.IsPastDate())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with future date", func(t *testing.T) {
		g := NewWithT(t)
		futureDate := "2099-12-31T23:59:59Z"
		err := validation.Validate(futureDate, validation.IsPastDate())
		g.Expect(err).To(MatchError(ContainSubstring("must be a past date")))
	})

	t.Run("fails with invalid format", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-date", validation.IsPastDate())
		g.Expect(err).To(MatchError(ContainSubstring("invalid date format")))
	})
}

func TestIsDateBefore(t *testing.T) {

	t.Run("passes when date is before", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-01T00:00:00Z", validation.IsDateBefore("2024-12-31T23:59:59Z"))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when date is after", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-12-31T23:59:59Z", validation.IsDateBefore("2024-01-01T00:00:00Z"))
		g.Expect(err).To(MatchError(ContainSubstring("must be before")))
	})

	t.Run("fails when date is equal", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-01T00:00:00Z", validation.IsDateBefore("2024-01-01T00:00:00Z"))
		g.Expect(err).To(MatchError(ContainSubstring("must be before")))
	})

	t.Run("fails with invalid date format", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-date", validation.IsDateBefore("2024-01-01T00:00:00Z"))
		g.Expect(err).To(MatchError(ContainSubstring("invalid date format")))
	})

	t.Run("fails with invalid beforeDate parameter", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-01T00:00:00Z", validation.IsDateBefore("not-a-date"))
		g.Expect(err).To(MatchError(ContainSubstring("invalid before date format")))
	})
}

func TestIsDateAfter(t *testing.T) {

	t.Run("passes when date is after", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-12-31T23:59:59Z", validation.IsDateAfter("2024-01-01T00:00:00Z"))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when date is before", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-01T00:00:00Z", validation.IsDateAfter("2024-12-31T23:59:59Z"))
		g.Expect(err).To(MatchError(ContainSubstring("must be after")))
	})

	t.Run("fails when date is equal", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-01T00:00:00Z", validation.IsDateAfter("2024-01-01T00:00:00Z"))
		g.Expect(err).To(MatchError(ContainSubstring("must be after")))
	})

	t.Run("fails with invalid date format", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("not-a-date", validation.IsDateAfter("2024-01-01T00:00:00Z"))
		g.Expect(err).To(MatchError(ContainSubstring("invalid date format")))
	})

	t.Run("fails with invalid afterDate parameter", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.Validate("2024-01-01T00:00:00Z", validation.IsDateAfter("not-a-date"))
		g.Expect(err).To(MatchError(ContainSubstring("invalid after date format")))
	})
}

func TestValidateAnyMap(t *testing.T) {

	t.Run("passes with valid map", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]any{"name": "John", "age": 30}
		err := validation.ValidateAnyMap(m, true,
			validation.MapKey("name", true, validation.StringValidator(validation.Required[string]())),
			validation.MapKey("age", true, validation.IntValidator(validation.Range(0, 120))),
		)
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with float64 as int", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]any{"age": float64(30)} // JSON numbers are float64
		err := validation.ValidateAnyMap(m, true,
			validation.MapKey("age", true, validation.IntValidator(validation.Range(0, 120))),
		)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when required key missing", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]any{"age": 30}
		err := validation.ValidateAnyMap(m, true,
			validation.MapKey("name", true, validation.StringValidator(validation.Required[string]())),
		)
		g.Expect(err).To(MatchError(ContainSubstring("is required")))
	})

	t.Run("fails with extra keys when not allowed", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]any{"name": "John", "extra": "value"}
		err := validation.ValidateAnyMap(m, false,
			validation.MapKey("name", true, validation.StringValidator(validation.Required[string]())),
		)
		g.Expect(err).To(MatchError(ContainSubstring("not expected")))
	})

	t.Run("allows extra keys when allowed", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]any{"name": "John", "extra": "value"}
		err := validation.ValidateAnyMap(m, true,
			validation.MapKey("name", true, validation.StringValidator(validation.Required[string]())),
		)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with validation error", func(t *testing.T) {
		g := NewWithT(t)
		m := map[string]any{"age": 150}
		err := validation.ValidateAnyMap(m, true,
			validation.MapKey("age", true, validation.IntValidator(validation.Range(0, 120))),
		)
		g.Expect(err).To(MatchError(ContainSubstring("age")))
		g.Expect(err).To(MatchError(ContainSubstring("must be between")))
	})
}

func TestStringValidator(t *testing.T) {

	t.Run("passes with valid string", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.StringValidator(validation.MinLength(3))
		err := validator("hello")
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid string", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.StringValidator(validation.MinLength(3))
		err := validator("hi")
		g.Expect(err).To(MatchError(ContainSubstring("must be at least 3 characters")))
	})

	t.Run("fails when value is not a string", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.StringValidator(validation.MinLength(3))
		err := validator(123)
		g.Expect(err).To(MatchError(ContainSubstring("must be a string")))
	})

	t.Run("fails when value is nil", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.StringValidator(validation.Required[string]())
		err := validator(nil)
		g.Expect(err).To(MatchError(ContainSubstring("must be a string")))
	})
}

func TestIntValidator(t *testing.T) {

	t.Run("passes with int", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.IntValidator(validation.Range(0, 120))
		err := validator(30)
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with float64", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.IntValidator(validation.Range(0, 120))
		err := validator(float64(30))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with int64", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.IntValidator(validation.Range(0, 120))
		err := validator(int64(30))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid int", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.IntValidator(validation.Range(0, 120))
		err := validator(150)
		g.Expect(err).To(MatchError(ContainSubstring("must be between")))
	})

	t.Run("fails when value is not a number", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.IntValidator(validation.Range(0, 120))
		err := validator("not a number")
		g.Expect(err).To(MatchError(ContainSubstring("must be a number")))
	})

	t.Run("fails when value is nil", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.IntValidator(validation.Range(0, 120))
		err := validator(nil)
		g.Expect(err).To(MatchError(ContainSubstring("must be a number")))
	})
}

func TestFloatValidator(t *testing.T) {

	t.Run("passes with float64", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.FloatValidator(validation.Range(0.0, 100.0))
		err := validator(50.5)
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with float32", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.FloatValidator(validation.Range(0.0, 100.0))
		err := validator(float32(50.5))
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with int converted to float", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.FloatValidator(validation.Range(0.0, 100.0))
		err := validator(50)
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with int64 converted to float", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.FloatValidator(validation.Range(0.0, 100.0))
		err := validator(int64(50))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with invalid float", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.FloatValidator(validation.Range(0.0, 100.0))
		err := validator(150.5)
		g.Expect(err).To(MatchError(ContainSubstring("must be between")))
	})

	t.Run("fails when value is not a number", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.FloatValidator(validation.Range(0.0, 100.0))
		err := validator("not a number")
		g.Expect(err).To(MatchError(ContainSubstring("must be a number")))
	})

	t.Run("fails when value is nil", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.FloatValidator(validation.Range(0.0, 100.0))
		err := validator(nil)
		g.Expect(err).To(MatchError(ContainSubstring("must be a number")))
	})
}

func TestBoolValidator(t *testing.T) {

	t.Run("passes with true", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.BoolValidator(validation.Custom(func(v bool) error {
			return nil // Accept any bool
		}))
		err := validator(true)
		g.Expect(err).To(BeNil())
	})

	t.Run("passes with false", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.BoolValidator(validation.Custom(func(v bool) error {
			return nil // Accept any bool
		}))
		err := validator(false)
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with custom validation", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.BoolValidator(validation.Custom(func(v bool) error {
			if !v {
				return fmt.Errorf("must be true")
			}
			return nil
		}))
		err := validator(false)
		g.Expect(err).To(MatchError(ContainSubstring("must be true")))
	})

	t.Run("fails when value is not a boolean", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.BoolValidator(validation.Custom(func(v bool) error {
			return nil
		}))
		err := validator("not a bool")
		g.Expect(err).To(MatchError(ContainSubstring("must be a boolean")))
	})

	t.Run("fails when value is nil", func(t *testing.T) {
		g := NewWithT(t)
		validator := validation.BoolValidator(validation.Custom(func(v bool) error {
			return nil
		}))
		err := validator(nil)
		g.Expect(err).To(MatchError(ContainSubstring("must be a boolean")))
	})
}

// Test types for ValidateNested
type testAddress struct {
	Street string
	City   string
}

type testValidatableAddress struct {
	Street string
	City   string
}

func (a testValidatableAddress) Validate() error {
	return errors.Join(
		validation.Validate(a.Street, validation.Required[string]()),
		validation.Validate(a.City, validation.Required[string]()),
	)
}

func TestValidateNested(t *testing.T) {

	t.Run("validates struct that implements Validatable", func(t *testing.T) {
		g := NewWithT(t)
		addr := testValidatableAddress{Street: "123 Main St", City: "NYC"}
		err := validation.ValidateNested(addr)
		g.Expect(err).To(BeNil())
	})

	t.Run("returns error for invalid Validatable struct", func(t *testing.T) {
		g := NewWithT(t)
		addr := testValidatableAddress{Street: "", City: "NYC"}
		err := validation.ValidateNested(addr)
		g.Expect(err).ToNot(BeNil())
	})

	t.Run("returns nil for non-Validatable struct", func(t *testing.T) {
		g := NewWithT(t)
		addr := testAddress{Street: "", City: ""}
		err := validation.ValidateNested(addr)
		g.Expect(err).To(BeNil()) // Non-Validatable structs pass through
	})

	t.Run("returns nil for primitive type", func(t *testing.T) {
		g := NewWithT(t)
		err := validation.ValidateNested("string")
		g.Expect(err).To(BeNil())
	})
}

func TestNested(t *testing.T) {

	t.Run("validates valid nested struct", func(t *testing.T) {
		g := NewWithT(t)
		addr := testValidatableAddress{Street: "123 Main St", City: "NYC"}
		validator := validation.Nested[testValidatableAddress]()
		err := validator(addr)
		g.Expect(err).To(BeNil())
	})

	t.Run("returns error for invalid nested struct", func(t *testing.T) {
		g := NewWithT(t)
		addr := testValidatableAddress{Street: "", City: "NYC"}
		validator := validation.Nested[testValidatableAddress]()
		err := validator(addr)
		g.Expect(err).ToNot(BeNil())
	})

	t.Run("works with Validate function", func(t *testing.T) {
		g := NewWithT(t)
		addr := testValidatableAddress{Street: "123 Main St", City: "NYC"}
		err := validation.Validate(addr, validation.Nested[testValidatableAddress]())
		g.Expect(err).To(BeNil())
	})
}

func TestIsFutureTime(t *testing.T) {

	t.Run("passes with future time", func(t *testing.T) {
		g := NewWithT(t)
		futureTime := time.Now().Add(24 * time.Hour)
		err := validation.Validate(futureTime, validation.IsFutureTime())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with past time", func(t *testing.T) {
		g := NewWithT(t)
		pastTime := time.Now().Add(-24 * time.Hour)
		err := validation.Validate(pastTime, validation.IsFutureTime())
		g.Expect(err).To(MatchError(ContainSubstring("must be a future time")))
	})

	t.Run("fails with now or past time", func(t *testing.T) {
		g := NewWithT(t)
		// Use a time slightly in the past to avoid timing issues
		almostNow := time.Now().Add(-1 * time.Millisecond)
		err := validation.Validate(almostNow, validation.IsFutureTime())
		g.Expect(err).To(MatchError(ContainSubstring("must be a future time")))
	})
}

func TestIsPastTime(t *testing.T) {

	t.Run("passes with past time", func(t *testing.T) {
		g := NewWithT(t)
		pastTime := time.Now().Add(-24 * time.Hour)
		err := validation.Validate(pastTime, validation.IsPastTime())
		g.Expect(err).To(BeNil())
	})

	t.Run("fails with future time", func(t *testing.T) {
		g := NewWithT(t)
		futureTime := time.Now().Add(24 * time.Hour)
		err := validation.Validate(futureTime, validation.IsPastTime())
		g.Expect(err).To(MatchError(ContainSubstring("must be a past time")))
	})

	t.Run("fails with now or future time", func(t *testing.T) {
		g := NewWithT(t)
		// Use a time slightly in the future to avoid timing issues
		almostNow := time.Now().Add(1 * time.Millisecond)
		err := validation.Validate(almostNow, validation.IsPastTime())
		g.Expect(err).To(MatchError(ContainSubstring("must be a past time")))
	})
}

func TestIsTimeBefore(t *testing.T) {

	t.Run("passes when time is before", func(t *testing.T) {
		g := NewWithT(t)
		before := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		timeToValidate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		err := validation.Validate(timeToValidate, validation.IsTimeBefore(before))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when time is after", func(t *testing.T) {
		g := NewWithT(t)
		before := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		timeToValidate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		err := validation.Validate(timeToValidate, validation.IsTimeBefore(before))
		g.Expect(err).To(MatchError(ContainSubstring("must be before")))
	})

	t.Run("fails when time is equal", func(t *testing.T) {
		g := NewWithT(t)
		sameTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		err := validation.Validate(sameTime, validation.IsTimeBefore(sameTime))
		g.Expect(err).To(MatchError(ContainSubstring("must be before")))
	})
}

func TestIsTimeAfter(t *testing.T) {

	t.Run("passes when time is after", func(t *testing.T) {
		g := NewWithT(t)
		after := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		timeToValidate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		err := validation.Validate(timeToValidate, validation.IsTimeAfter(after))
		g.Expect(err).To(BeNil())
	})

	t.Run("fails when time is before", func(t *testing.T) {
		g := NewWithT(t)
		after := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		timeToValidate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		err := validation.Validate(timeToValidate, validation.IsTimeAfter(after))
		g.Expect(err).To(MatchError(ContainSubstring("must be after")))
	})

	t.Run("fails when time is equal", func(t *testing.T) {
		g := NewWithT(t)
		sameTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		err := validation.Validate(sameTime, validation.IsTimeAfter(sameTime))
		g.Expect(err).To(MatchError(ContainSubstring("must be after")))
	})
}

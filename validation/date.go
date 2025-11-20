package validation

import (
	"fmt"
	"time"
)

// IsRFC3339DateTime validates that a string is a valid RFC3339 date-time.
// RFC3339 is the standard format used by JSON and most APIs (e.g., "2006-01-02T15:04:05Z07:00").
//
// Example:
//
//	validation.Validate(timestamp, validation.IsRFC3339DateTime())
func IsRFC3339DateTime() Validator[string] {
	return func(v string) error {
		_, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return NewValidationError("must be a valid RFC3339 date-time")
		}
		return nil
	}
}

// IsISO8601Date validates that a string is a valid ISO8601 date in YYYY-MM-DD format.
//
// Example:
//
//	validation.Validate(birthDate, validation.IsISO8601Date())
func IsISO8601Date() Validator[string] {
	return func(v string) error {
		_, err := time.Parse("2006-01-02", v)
		if err != nil {
			return NewValidationError("must be a valid ISO8601 date (YYYY-MM-DD)")
		}
		return nil
	}
}

// IsDateFormat validates that a string matches the specified date format.
// The layout uses Go's time format layout (reference time: Mon Jan 2 15:04:05 MST 2006).
//
// Example:
//
//	validation.Validate(date, validation.IsDateFormat("2006-01-02 15:04:05"))
func IsDateFormat(layout string) Validator[string] {
	return func(v string) error {
		_, err := time.Parse(layout, v)
		if err != nil {
			return NewValidationError(fmt.Sprintf("must match date format %q", layout))
		}
		return nil
	}
}

// IsFutureDateFormat validates that a string represents a date in the future.
// The date is parsed using the specified layout.
//
// Example:
//
//	validation.Validate(expiryDate, validation.IsFutureDateFormat("2006-01-02"))
func IsFutureDateFormat(layout string) Validator[string] {
	return func(v string) error {
		t, err := time.Parse(layout, v)
		if err != nil {
			return NewValidationError("invalid date format")
		}
		if !t.After(time.Now()) {
			return NewValidationError("must be a future date")
		}
		return nil
	}
}

// IsFutureDate validates that a string represents a date in the future.
// The date is parsed as RFC3339 format.
//
// Example:
//
//	validation.Validate(expiryDate, validation.IsFutureDate())
func IsFutureDate() Validator[string] {
	return IsFutureDateFormat(time.RFC3339)
}

// IsPastDateFormat validates that a string represents a date in the past.
// The date is parsed using the specified layout.
//
// Example:
//
//	validation.Validate(birthDate, validation.IsPastDateFormat("2006-01-02"))
func IsPastDateFormat(layout string) Validator[string] {
	return func(v string) error {
		t, err := time.Parse(layout, v)
		if err != nil {
			return NewValidationError("invalid date format")
		}
		if !t.Before(time.Now()) {
			return NewValidationError("must be a past date")
		}
		return nil
	}
}

// IsPastDate validates that a string represents a date in the past.
// The date is parsed as RFC3339 format.
//
// Example:
//
//	validation.Validate(birthDate, validation.IsPastDate())
func IsPastDate() Validator[string] {
	return IsPastDateFormat(time.RFC3339)
}

// IsDateBeforeFormat validates that a string represents a date before the specified date.
// Both dates are parsed using the specified layout.
//
// Example:
//
//	validation.Validate(startDate, validation.IsDateBeforeFormat("2024-12-31", "2006-01-02"))
func IsDateBeforeFormat(beforeDate, layout string) Validator[string] {
	return func(v string) error {
		t, err := time.Parse(layout, v)
		if err != nil {
			return NewValidationError("invalid date format")
		}
		before, err := time.Parse(layout, beforeDate)
		if err != nil {
			return NewValidationError("invalid before date format")
		}
		if !t.Before(before) {
			return NewValidationError(fmt.Sprintf("must be before %s", before.Format(layout)))
		}
		return nil
	}
}

// IsDateBefore validates that a string represents a date before the specified date.
// Both dates are parsed as RFC3339 format.
//
// Example:
//
//	validation.Validate(startDate, validation.IsDateBefore("2024-12-31T23:59:59Z"))
func IsDateBefore(beforeDate string) Validator[string] {
	return IsDateBeforeFormat(beforeDate, time.RFC3339)
}

// IsDateAfterFormat validates that a string represents a date after the specified date.
// Both dates are parsed using the specified layout.
//
// Example:
//
//	validation.Validate(endDate, validation.IsDateAfterFormat("2024-01-01", "2006-01-02"))
func IsDateAfterFormat(afterDate, layout string) Validator[string] {
	return func(v string) error {
		t, err := time.Parse(layout, v)
		if err != nil {
			return NewValidationError("invalid date format")
		}
		after, err := time.Parse(layout, afterDate)
		if err != nil {
			return NewValidationError("invalid after date format")
		}
		if !t.After(after) {
			return NewValidationError(fmt.Sprintf("must be after %s", after.Format(layout)))
		}
		return nil
	}
}

// IsDateAfter validates that a string represents a date after the specified date.
// Both dates are parsed as RFC3339 format.
//
// Example:
//
//	validation.Validate(endDate, validation.IsDateAfter("2024-01-01T00:00:00Z"))
func IsDateAfter(afterDate string) Validator[string] {
	return IsDateAfterFormat(afterDate, time.RFC3339)
}

// IsFutureTime validates that a time.Time is in the future.
//
// Example:
//
//	validation.Validate(expiryTime, validation.IsFutureTime())
func IsFutureTime() Validator[time.Time] {
	return func(v time.Time) error {
		if !v.After(time.Now()) {
			return NewValidationError("must be a future time")
		}
		return nil
	}
}

// IsPastTime validates that a time.Time is in the past.
//
// Example:
//
//	validation.Validate(birthTime, validation.IsPastTime())
func IsPastTime() Validator[time.Time] {
	return func(v time.Time) error {
		if !v.Before(time.Now()) {
			return NewValidationError("must be a past time")
		}
		return nil
	}
}

// IsTimeBefore validates that a time.Time is before the specified time.
//
// Example:
//
//	validation.Validate(startTime, validation.IsTimeBefore(endTime))
func IsTimeBefore(before time.Time) Validator[time.Time] {
	return func(v time.Time) error {
		if !v.Before(before) {
			return NewValidationError(fmt.Sprintf("must be before %s", before.Format(time.RFC3339)))
		}
		return nil
	}
}

// IsTimeAfter validates that a time.Time is after the specified time.
//
// Example:
//
//	validation.Validate(endTime, validation.IsTimeAfter(startTime))
func IsTimeAfter(after time.Time) Validator[time.Time] {
	return func(v time.Time) error {
		if !v.After(after) {
			return NewValidationError(fmt.Sprintf("must be after %s", after.Format(time.RFC3339)))
		}
		return nil
	}
}

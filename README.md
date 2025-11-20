# Protego: Type-Safe Golang Validation Package

A type-safe validation framework for Go using generics. This package provides compile-time type safety, ensuring validators can only be applied to compatible types.

## Features

- **Type-safe**: Validators are strongly typed - you can't apply a string validator to an integer
- **Composable**: Chain multiple validators together
- **Reusable**: Define validators once, use them everywhere
- **Flexible**: Wrap go-playground/validator, ozzo-validation, or ANY validation library
- **Zero dependencies**: Core package has no external dependencies
- **Battle-tested**: Optional playground package provides 100+ validators from go-playground
- **Error Detection**: All errors are wrapped in `validation.Error` for easy identification
- **Clean API**: Simple, readable validation code

## Philosophy

This package provides **primitives** for type-safe validation. The core package intentionally has **zero external dependencies** and only implements simple, fast validators.

**For complex validators** (email, UUID, IP addresses, etc.), you can:
- Use the `playground` subpackage (40+ pre-built validators using github.com/go-playground/validator/v10)
- Create your own using `playground.FromTag` (100+ go-playground validators)
- Wrap ANY other validation library you prefer

**You control the dependencies**, not the validation framework.

## Installation

Install the core validation package:

```bash
go get github.com/quantumcycle/protego/validation
```

Optionally, install the playground package for 100+ pre-built validators:

```bash
go get github.com/quantumcycle/protego/playground
```

Then import in your code:

```go
import "github.com/quantumcycle/protego/validation"
import "github.com/quantumcycle/protego/playground"  // Optional
```

## Quick Start

```go
import (
    "github.com/quantumcycle/protego/validation"
    "github.com/quantumcycle/protego/playground"
)

type CreateUserInput struct {
    Email    string
    Username string
    Age      int
}

func (input CreateUserInput) Validate() error {
    return errors.Join(
        validation.Validate(input.Email,
            validation.Required[string](),     // Core validator
            playground.IsEmail,                // Playground validator (go-playground)
        ),
        validation.Validate(input.Username,
            validation.Required[string](),     // Core validator
            validation.Length(3, 20),          // Core validator
            playground.IsAlphanumeric,         // Playground validator (go-playground)
        ),
        validation.Validate(input.Age,
            validation.Required[int](),        // Core validator
            validation.Range(18, 120),         // Core validator
        ),
    )
}
```

**Mix and match**: Use core validators for simple checks, playground validators for complex formats, or wrap your own libraries!

## Type Safety in Action

The compiler prevents invalid validator usage:

```go
// This compiles - MinLength works on strings
validation.Validate(name, validation.MinLength(3))

// This compiles - Range works on ints
validation.Validate(age, validation.Range(0, 120))

// This won't compile - MinLength doesn't work on ints
validation.Validate(age, validation.MinLength(3))

// This won't compile - Range expects same type
validation.Validate(age, validation.Range("0", "120"))
```

## Error Detection

All validation errors in Protego are wrapped in a `validation.Error` type, making it easy to detect and handle Protego-specific errors:

```go
import (
    "errors"
    "github.com/quantumcycle/protego/validation"
    "github.com/quantumcycle/protego/playground"
)

func ProcessUser(input CreateUserInput) error {
    err := input.Validate()
    if err != nil {
        // Check if this is a Protego validation error
        if validation.IsValidationError(err) {
            // Handle validation errors specifically
            return fmt.Errorf("validation failed: %w", err)
        }
        // Handle other types of errors
        return fmt.Errorf("unexpected error: %w", err)
    }
    // Process valid input
    return nil
}
```

### Error Detection Features

- **Type Detection**: Use `validation.IsValidationError(err)` to check if an error came from Protego
- **Error Wrapping**: All validators wrap errors using `validation.Error`, including playground validators
- **Error Unwrapping**: Supports Go's standard `errors.Unwrap()` and `errors.Is()` functions
- **Preserved Messages**: Original error messages remain unchanged for backward compatibility

### Examples

```go
// Detect validation errors from core validators
err := validation.Validate("", validation.Required[string]())
if validation.IsValidationError(err) {
    fmt.Println("Protego validation error:", err.Error()) // Output: required
}

// Detect validation errors from playground validators
err = validation.Validate("invalid-email", playground.IsEmail)
if validation.IsValidationError(err) {
    fmt.Println("Email validation failed:", err.Error())
}

// Use with errors.Join for multiple validations
err = errors.Join(
    validation.Validate("", validation.Required[string]()),
    validation.Validate("ab", validation.MinLength(3)),
)
// Check if any are validation errors
if validation.IsValidationError(err) {
    fmt.Println("Contains validation errors")
}

// Error unwrapping works
originalErr := errors.Unwrap(err)
```

## Available Validators

### Required Validators

```go
validation.Required[T]()                    // Value must not be zero value
validation.RequiredIf[T](condition)         // Required if condition is true
```

### String Validators

```go
validation.MinLength(min)                   // Minimum string length
validation.MaxLength(max)                   // Maximum string length
validation.Length(min, max)                 // String length range
validation.IsInt()                          // String represents integer
validation.MatchesPattern(regex)            // Matches regex pattern
validation.StartsWith(prefix)               // Starts with prefix
validation.EndsWith(suffix)                 // Ends with suffix
validation.Contains(substring)              // Contains substring
```

### Numeric Validators

```go
validation.Min[T](min)                      // Minimum value
validation.Max[T](max)                      // Maximum value
validation.Range[T](min, max)               // Value range
validation.GreaterThan[T](threshold)        // Strictly greater than
validation.LessThan[T](threshold)           // Strictly less than
validation.Positive[T]()                    // Greater than zero
validation.NonNegative[T]()                 // Greater than or equal to zero
validation.Negative[T]()                    // Less than zero
validation.MultipleOf(divisor)              // Multiple of divisor
```

### Collection Validators

```go
validation.In(caseSensitive, allowed...)    // Value in allowed list
validation.InSlice(caseSensitive, allowed)  // Value in allowed slice
validation.NotIn(caseSensitive, forbidden)  // Value not in forbidden list
validation.Each(validator)                  // Validate each element
validation.NotEmpty[T]()                    // Slice not empty
validation.MinItems[T](min)                 // Minimum slice length
validation.MaxItems[T](max)                 // Maximum slice length
validation.UniqueItems[T]()                 // All items unique
```

### Date/Time Validators

```go
validation.IsRFC3339DateTime()              // Valid RFC3339 date-time
validation.IsISO8601Date()                  // Valid ISO8601 date (YYYY-MM-DD)
validation.IsDateFormat(layout)             // Matches date format
validation.IsFutureDate()                   // Date in the future
validation.IsPastDate()                     // Date in the past
validation.IsDateBefore(date)               // Date before specified date
validation.IsDateAfter(date)                // Date after specified date
```

### Optional/Pointer Validators

```go
validation.NilOrNotEmpty()                  // Nil or non-empty string
validation.NilOr(validator)                 // Nil or passes validator
validation.NotNil[T]()                      // Not nil
validation.OptionalWith(validators...)      // Nil or passes all validators
```

### Helper Functions

```go
validation.WithMessage(validator, msg)      // Custom error message
validation.And(validators...)               // All must pass
validation.Or(validators...)                // At least one must pass
validation.Not(validator)                   // Inverts validator
validation.When(condition, validator)       // Apply if condition true
validation.Unless(condition, validator)     // Apply if condition false
validation.Custom(fn)                       // Custom validator function
```

## Playground Package - Pre-Built Validators

The `playground` subpackage provides **40+ pre-built validators** for common use cases:

```go
import "github.com/quantumcycle/protego/playground"

func (input CreateServerInput) Validate() error {
    return errors.Join(
        validation.Validate(input.Email, playground.IsEmail),
        validation.Validate(input.URL, playground.IsURL),
        validation.Validate(input.ID, playground.IsUUID4),
        validation.Validate(input.IP, playground.IsIPv4),
    )
}
```

### Available in playground package:

**String Format**: `IsEmail`, `IsURL`, `IsURI`, `IsAlpha`, `IsAlphanumeric`, `IsNumeric`, `IsLowercase`, `IsUppercase`, `IsASCII`

**Network**: `IsIPv4`, `IsIPv6`, `IsIP`, `IsCIDR`, `IsMAC`, `IsHostname`, `IsFQDN`

**Identifiers**: `IsUUID`, `IsUUID3`, `IsUUID4`, `IsUUID5`, `IsULID`

**Encoding**: `IsBase64`, `IsBase64URL`, `IsHex`, `IsHexColor`, `IsJSON`, `IsJWT`

**Payment**: `IsCreditCard`, `IsBTC`, `IsETH`, `IsIBAN`

**Versions**: `IsSemver`, `IsISBN`, `IsISBN10`, `IsISBN13`

**Other**: `IsE164` (phone), `IsLatitude`, `IsLongitude`, `IsBoolean`

See [playground/playground.go](/pkg/playground/playground.go) for the complete list.

## Extending with go-playground/validator

You can access **100+ validators** from https://github.com/go-playground/validator using the `playground.FromTag` helper:

### Using FromTag

```go
import "github.com/quantumcycle/protego/playground"

// Create reusable validators from go-playground tags
var (
    IsUUID     = playground.FromTag[string]("uuid")
    IsUUID4    = playground.FromTag[string]("uuid4")
    IsIPv4Custom = playground.FromTag[string]("ipv4")
    IsIPv6Custom = playground.FromTag[string]("ipv6")
    IsCIDRCustom = playground.FromTag[string]("cidr")
    IsMACCustom  = playground.FromTag[string]("mac")
)

// Use them like any other validator
func (input CreateServerInput) Validate() error {
    return errors.Join(
        validation.Validate(input.ID, validation.Required[string](), IsUUID4),
        validation.Validate(input.IP, validation.Required[string](), IsIPv4Custom),
    )
}
```

### Inline Usage

```go
import (
    "github.com/quantumcycle/protego/validation"
    "github.com/quantumcycle/protego/playground"
)

// Or use FromTag inline without creating variables
func (input Input) Validate() error {
    return errors.Join(
        validation.Validate(input.SessionID,
            validation.Required[string](),
            playground.FromTag[string]("uuid4"),
        ),
        validation.Validate(input.Port,
            validation.Required[int](),
            playground.FromTag[int]("min=1,max=65535"),
        ),
    )
}
```

### Available go-playground Validators

go-playground/validator provides 100+ validators including:

**Network**:
- `ipv4`, `ipv6`, `ip`, `cidr`, `cidrv4`, `cidrv6`
- `mac`, `hostname`, `fqdn`, `url`, `uri`, `url_encoded`
- `tcp_addr`, `tcp4_addr`, `tcp6_addr`, `udp_addr`, `udp4_addr`, `udp6_addr`

**Identifiers**:
- `uuid`, `uuid3`, `uuid4`, `uuid5`, `ulid`, `ascii`, `printascii`

**Encoding**:
- `base64`, `base64url`, `base64rawurl`, `hexadecimal`, `hexcolor`, `rgb`, `rgba`
- `hsl`, `hsla`, `html`, `html_encoded`, `json`, `jwt`

**Payment & Finance**:
- `credit_card`, `luhn_checksum`, `btc_addr`, `btc_addr_bech32`, `eth_addr`
- `iban`, `bic`, `semver`, `mongodb`, `cron`

**Geographic**:
- `latitude`, `longitude`, `postcode_iso3166_alpha2`, `postcode_iso3166_alpha2_field`

**Misc**:
- `e164` (phone), `isbn`, `isbn10`, `isbn13`, `issn`, `email`, `alpha`, `alphanum`, `numeric`
- `boolean`, `datetime`, `timezone`, `dir`, `dirpath`, `file`, `filepath`

See [go-playground/validator documentation](https://pkg.go.dev/github.com/go-playground/validator/v10) for the complete list.

### Custom Error Messages

```go
import "github.com/quantumcycle/protego/playground"

validation.Validate(input.ID,
    playground.FromTagWithMessage[string]("uuid4", "must be a valid UUID v4"),
)
```

### Example: Wrapping a custom library

```go
import (
    "your-company/internal-validators"
    "github.com/quantumcycle/protego/validation"
)

// Wrap your existing validators
func CompanyEmailDomain() validation.Validator[string] {
    return func(v string) error {
        if !internalValidators.IsCompanyEmail(v) {
            return fmt.Errorf("must be a @company.com email")
        }
        return nil
    }
}

func (input Input) Validate() error {
    return errors.Join(
        validation.Validate(input.Email, validation.Required[string](), CompanyEmailDomain()),
        validation.Validate(input.Username, validation.Length(3, 20)),
    )
}
```

### Creating a validator factory

```go
// Create a helper to wrap any validation library
func FromLibraryX[T any](libraryValidator interface{}) validation.Validator[T] {
    return func(v T) error {
        // Call your library's validation function
        return libraryX.Validate(v, libraryValidator)
    }
}

// Now use it everywhere
var (
    IsUUID = FromLibraryX[string](libraryX.UUID)
    IsURL  = FromLibraryX[string](libraryX.URL)
)
```

## Examples

### Basic Validation

```go
func (input MyInput) Validate() error {
    return errors.Join(
        validation.Validate(input.AssetID,
            validation.Required[string](),
            validation.IsInt(),
        ),
        validation.Validate(input.Currency,
            validation.Required[string](),
            validation.InSlice(true, []string{"USD", "EUR", "GBP"}),
        ),
    )
}
```

### Nested Struct Validation

```go
type Address struct {
    Street string
    City   string
}

func (a Address) Validate() error {
    return errors.Join(
        validation.Validate(a.Street, validation.Required[string]()),
        validation.Validate(a.City, validation.Required[string]()),
    )
}

type User struct {
    Name    string
    Address Address
}

func (u User) Validate() error {
    return errors.Join(
        validation.Validate(u.Name, validation.Required[string]()),
        validation.ValidateNested(u.Address), // Calls Address.Validate()
    )
}
```

### Optional Fields

```go
import (
    "github.com/quantumcycle/protego/validation"
    "github.com/quantumcycle/protego/playground"
)

type User struct {
    Email    *string
    Age      *int
    Priority *int
}

func (input User) Validate() error {
    return errors.Join(
        validation.Validate(input.Email,
            validation.NilOr(playground.IsEmail),
        ),
        validation.Validate(input.Age,
            validation.NilOr(validation.Range(0, 120)),
        ),
        validation.Validate(input.Priority,
            validation.NilOr(validation.In(false, 1, 2, 3)),
        ),
    )
}
```

### Slice Validation

```go
import (
    "github.com/quantumcycle/protego/validation"
    "github.com/quantumcycle/protego/playground"
)

func (input CreateBatch) Validate() error {
    return errors.Join(
        validation.Validate(input.Emails,
            validation.NotEmpty[string](),
            validation.MaxItems[string](100),
            validation.Each(playground.IsEmail),
        ),
        validation.Validate(input.Tags,
            validation.UniqueItems[string](),
        ),
    )
}
```

### Map Validation

```go
import (
    "github.com/quantumcycle/protego/validation"
    "github.com/quantumcycle/protego/playground"
)

func (input CreateVersion) Validate() error {
    // For map[string]string
    return validation.ValidateStringMap(
        input.Metadata,
        true, // allow extra keys
        validation.MapKey("version", true, validation.Required[string](), playground.IsSemver),
        validation.MapKey("author", true, validation.Required[string]()),
    )
}

func (input SubmitVersion) Validate() error {
    // For map[string]any (JSON)
    return validation.ValidateAnyMap(
        input.Manifest,
        true, // allow extra keys
        validation.MapKey("name", true, validation.StringValidator(validation.Required[string]())),
        validation.MapKey("version", true, validation.StringValidator(playground.IsSemver)),
    )
}
```

### Custom Error Messages

```go
func (input SearchInput) Validate() error {
    return errors.Join(
        validation.Validate(input.Count,
            validation.WithMessage(
                validation.Required[int](),
                "count cannot be zero",
            ),
            validation.Range(1, 1000),
        ),
    )
}
```

### Conditional Validation

```go
func (input Order) Validate() error {
    return errors.Join(
        validation.Validate(input.ShippingAddress,
            validation.When(
                input.RequiresShipping,
                validation.Required[string](),
            ),
        ),
        validation.Validate(input.RejectionReason,
            validation.Unless(
                input.IsApproved,
                validation.Required[string](),
            ),
        ),
    )
}
```

### Custom Validators

```go
// Define a reusable custom validator
func IsValidDomain() validation.Validator[string] {
    return func(v string) error {
        segments := strings.Split(v, ".")
        if len(segments) < 2 {
            return fmt.Errorf("namespace must have at least 2 parts")
        }
        for _, segment := range segments {
            if !isValidDomainSegment(segment) {
                return fmt.Errorf("invalid part: %q", segment)
            }
        }
        return nil
    }
}

// Use it
func (input Website) Validate() error {
    return validation.Validate(input.Domain,
        validation.Required[string](),
        IsValidDomain(),
    )
}

// Or define inline
func (input Post) Validate() error {
    return validation.Validate(input.Content,
        validation.Custom(func(v string) error {
            if containsProfanity(v) {
                return fmt.Errorf("contains inappropriate content")
            }
            return nil
        }),
    )
}
```

### Combining Validators

```go
import (
    "github.com/quantumcycle/protego/validation"
    "github.com/quantumcycle/protego/playground"
)

// Create a complex password validator
func StrongPassword() validation.Validator[string] {
    return validation.And(
        validation.MinLength(8),
        validation.MatchesPattern(`[A-Z]`), // uppercase letter
        validation.MatchesPattern(`[a-z]`), // lowercase letter
        validation.MatchesPattern(`[0-9]`), // digit
    )
}

// Accept either email or phone
func ContactValidator() validation.Validator[string] {
    return validation.Or(
        playground.IsEmail,
        validation.MatchesPattern(`^\+\d{10,}$`), // phone pattern
    )
}
```

## Benefits Over Struct Tags

| Feature | Struct Tags | This Package        |
|---------|-----------|---------------------|
| Type safety |  Runtime only |  Compile-time       |
| Multiple use cases |  Single use | Reusable validators |
| Complex logic |  Limited | Full Go code        |
| Composable |  No | Yes                 |
| IDE support |  String tags | Full autocomplete   |
| Refactoring |  Broken tags | Compiler catches    |


// Package playground provides pre-defined validators using go-playground/validator tags.
// These are common validators that users would typically need.
//
// This package serves as both:
// 1. Convenient pre-built validators for common use cases
// 2. Examples showing how to create your own validators using FromTag()
//
// Usage:
//
//	import "github.com/quantumcycle/protego/playground"
//
//	func (input Input) Validate() error {
//	    return errors.Join(
//	        validation.Validate(input.Email, playground.IsEmail),
//	        validation.Validate(input.URL, playground.IsURL),
//	        validation.Validate(input.ID, playground.IsUUID4),
//	    )
//	}
package playground

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/quantumcycle/protego/validation"
)

var sharedValidator = validator.New()

// FromTag creates a type-safe validation.Validator[T] from a go-playground/validator tag string.
// This allows you to leverage ANY of go-playground's 100+ built-in validators.
//
// Example:
//
//	// Create reusable validators
//	var IsUUID = playground.FromTag[string]("uuid")
//	var IsCustomPort = playground.FromTag[int]("min=1,max=65535")
//
//	// Or use inline
//	validation.Validate(id, playground.FromTag[string]("uuid4"))
//
// See https://pkg.go.dev/github.com/go-playground/validator/v10 for all available tags.
func FromTag[T any](tag string) validation.Validator[T] {
	return func(v T) error {
		return sharedValidator.Var(v, tag)
	}
}

// FromTagWithMessage creates a validator from a go-playground tag with a custom error message.
//
// Example:
//
//	validation.Validate(id,
//	    playground.FromTagWithMessage[string]("uuid4", "must be a valid UUID v4"),
//	)
func FromTagWithMessage[T any](tag, message string) validation.Validator[T] {
	return func(v T) error {
		if err := sharedValidator.Var(v, tag); err != nil {
			return fmt.Errorf("%s", message)
		}
		return nil
	}
}

// String Format Validators
var (
	IsEmail        = FromTag[string]("email")
	IsURL          = FromTag[string]("url")
	IsURI          = FromTag[string]("uri")
	IsAlpha        = FromTag[string]("alpha")
	IsAlphanumeric = FromTag[string]("alphanum")
	IsNumeric      = FromTag[string]("numeric")
	IsLowercase    = FromTag[string]("lowercase")
	IsUppercase    = FromTag[string]("uppercase")
	IsASCII        = FromTag[string]("ascii")
)

// Network Validators
var (
	IsIPv4     = FromTag[string]("ipv4")
	IsIPv6     = FromTag[string]("ipv6")
	IsIP       = FromTag[string]("ip")
	IsCIDR     = FromTag[string]("cidr")
	IsMAC      = FromTag[string]("mac")
	IsHostname = FromTag[string]("hostname")
	IsFQDN     = FromTag[string]("fqdn")
)

// Identifier Validators
var (
	IsUUID  = FromTag[string]("uuid")
	IsUUID3 = FromTag[string]("uuid3")
	IsUUID4 = FromTag[string]("uuid4")
	IsUUID5 = FromTag[string]("uuid5")
	IsULID  = FromTag[string]("ulid")
)

// Encoding Validators
var (
	IsBase64    = FromTag[string]("base64")
	IsBase64URL = FromTag[string]("base64url")
	IsHex       = FromTag[string]("hexadecimal")
	IsHexColor  = FromTag[string]("hexcolor")
	IsJSON      = FromTag[string]("json")
	IsJWT       = FromTag[string]("jwt")
)

// Payment & Finance Validators
var (
	IsCreditCard = FromTag[string]("credit_card")
	IsBTC        = FromTag[string]("btc_addr")
	IsETH        = FromTag[string]("eth_addr")
	IsIBAN       = FromTag[string]("iban")
)

// Version & Format Validators
var (
	IsSemver = FromTag[string]("semver")
	IsISBN   = FromTag[string]("isbn")
	IsISBN10 = FromTag[string]("isbn10")
	IsISBN13 = FromTag[string]("isbn13")
)

// Phone & Communication Validators
var (
	IsE164 = FromTag[string]("e164")
)

// Geographic Validators
var (
	IsLatitude  = FromTag[string]("latitude")
	IsLongitude = FromTag[string]("longitude")
)

// Boolean Validators
var (
	IsBoolean = FromTag[string]("boolean")
)

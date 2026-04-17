package validation

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// FieldDef holds a deferred validation for a single struct field.
type FieldDef[S any] struct {
	validate func() error
}

// Field creates a FieldDef that binds a pointer-to-field to its validators.
// The field name is resolved automatically from the struct via unsafe pointer arithmetic.
func Field[S any, T any](s *S, fieldPtr *T, validators ...Validator[T]) FieldDef[S] {
	return FieldDef[S]{
		validate: func() error {
			name := resolveFieldName(s, fieldPtr)
			if err := Validate(*fieldPtr, validators...); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			return nil
		},
	}
}

// ValidateStruct runs all FieldDefs and joins their errors.
func ValidateStruct[S any](fields ...FieldDef[S]) error {
	errs := make([]error, 0, len(fields))
	for _, f := range fields {
		errs = append(errs, f.validate())
	}
	return errors.Join(errs...)
}

// resolveFieldName finds the struct field name by matching the field pointer offset.
// Falls back to "unknown" if not found. Respects the "json" tag if present.
func resolveFieldName[S any, T any](s *S, fieldPtr *T) string {
	structAddr := uintptr(unsafe.Pointer(s))       //nolint:gosec // intentional: computing field offset by pointer arithmetic for field name resolution
	fieldAddr := uintptr(unsafe.Pointer(fieldPtr)) //nolint:gosec // intentional: computing field offset by pointer arithmetic for field name resolution
	offset := fieldAddr - structAddr

	t := reflect.TypeOf(*s)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Offset == offset {
			if tag, ok := f.Tag.Lookup("json"); ok && tag != "" && tag != "-" {
				// strip options like omitempty
				for j, c := range tag {
					if c == ',' {
						return tag[:j]
					}
				}
				return tag
			}
			return f.Name
		}
	}
	return "unknown"
}

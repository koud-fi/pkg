package errx

import (
	"errors"
	"fmt"
)

// ValidationError indicates that some input field failed validation.
type ValidationError struct {
	// Field is the name of the invalid field, e.g. "email", "age", "password".
	Field string
	// Reason describes why it failed, e.g. "required", "too short", "invalid format".
	Reason string
	// Value is the actual value that didn’t pass validation.
	Value any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %q (%v): %s", e.Field, e.Value, e.Reason)
}

// NewValidation constructs a ValidationError for the given field, reason, and value.
func NewValidation(field, reason string, value any) error {
	return wrap(&ValidationError{
		Field:  field,
		Reason: reason,
		Value:  value,
	})
}

// IsValidation reports whether err is (or wraps) a ValidationError.
func IsValidation(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

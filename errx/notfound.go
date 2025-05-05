package errx

import (
	"errors"
	"fmt"
)

// NotFoundError represents a “not found” for some kind of resource.
type NotFoundError struct {
	// Kind is the type of thing you were looking for, e.g. "User", "Order", "Config".
	Kind string
	// Key is the lookup key, e.g. an ID, a name, etc.
	Key any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %v", e.Kind, e.Key)
}

// IsNotFound reports whether err is (or wraps) a NotFoundError.
func IsNotFound(err error) bool {
	var nf *NotFoundError
	return errors.As(err, &nf)
}

// NewNotFound constructs a NotFoundError for the given kind and key.
func NewNotFound(kind string, key any) error {
	return Wrap(&NotFoundError{Kind: kind, Key: key})
}

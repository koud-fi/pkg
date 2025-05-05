package errx

import (
	"errors"
	"runtime"
)

// Wrap is the common error wrapping logic used by all E functions.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	// If it’s already our Error type with a stack trace, return it.
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	// Otherwise, wrap it with a stack trace.
	var (
		pcs [maxStackDepth]uintptr
		// Skip 3 frames: runtime.Callers, errors.E, and wrapError.
		n     = runtime.Callers(3, pcs[:])
		trace = make([]uintptr, n) // TODO: Should we pool these buffers?
	)
	copy(trace, pcs[:n])
	return &Error{cause: err, pcs: trace}
}

package errx

import (
	"errors"
	"fmt"
	"io"
	"runtime"
)

const maxStackDepth = 32

type Error struct {
	cause error
	pcs   []uintptr // program counter
}

// E wraps err (if non-nil) capturing the current stack.
func E(err error) error {
	return wrap(err)
}

// E1 wraps an error while preserving a single return value.
func E1[T any](v T, err error) (T, error) {
	return v, wrap(err)
}

// E2 wraps an error while preserving two return values.
func E2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2, error) {
	return v1, v2, wrap(err)
}

// E3 wraps an error while preserving three return values.
func E3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3, error) {
	return v1, v2, v3, wrap(err)
}

// E4 wraps an error while preserving four return values.
func E4[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4, error) {
	return v1, v2, v3, v4, wrap(err)
}

// New is alias for errors.New, that wraps the error with a stack trace.
func New(text string) error { return wrap(errors.New(text)) }

// Fmt is an alias for fmt.Errorf, that wraps the error with a stack trace.
func Fmt(format string, args ...any) error {
	return wrap(fmt.Errorf(format, args...))
}

func (e *Error) Error() string { return e.cause.Error() }

// Unwrap allows errors.Is / errors.As to reach the cause.
func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// First print the error message
			fmt.Fprintf(s, "%s\n", e.cause.Error())

			// Then unwind the pcs into human-readable frames
			frames := runtime.CallersFrames(e.pcs)
			for {
				frame, more := frames.Next()
				fmt.Fprintf(s, "  %s\n    %s:%d\n", frame.Function, frame.File, frame.Line)
				if !more {
					break
				}
			}
			return
		}
		fallthrough // for %v without +, fall through to %s
	case 's':
		io.WriteString(s, e.cause.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.cause.Error())
	}
}

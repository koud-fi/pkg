package rx

type Maybe[T any] struct {
	v  T
	ok bool
}

func Just[T any](v T) Maybe[T]    { return Maybe[T]{v: v, ok: true} }
func Partial[T any](v T) Maybe[T] { return Maybe[T]{v: v} }
func None[T any]() Maybe[T]       { return Maybe[T]{} }

func (m Maybe[T]) Value() T { return m.v }
func (m Maybe[_]) Ok() bool { return m.ok }

func IsOk[T any](v Maybe[T]) bool { return v.Ok() }

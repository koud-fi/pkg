package auth

import (
	"context"
	"fmt"
	"sync"
)

type (
	authContextKey             struct{}
	authContextValue[User any] struct {
		authFn  func(context.Context) (User, error)
		once    sync.Once
		authErr error

		hasUser bool
		user    User
	}
)

func ContextWithAuthFunc[User any](
	parent context.Context, fn func(context.Context) (User, error),
) context.Context {
	return context.WithValue(parent, authContextKey{}, &authContextValue[User]{
		authFn: fn,
	})
}

func ContextWithUser[User any](
	parent context.Context, user User,
) context.Context {
	return context.WithValue(parent, authContextKey{}, &authContextValue[User]{
		hasUser: true,
		user:    user,
	})
}

func ContextUser[User any](ctx context.Context) (User, error) {
	var zero User
	v := ctx.Value(authContextKey{})
	if v == nil {
		return zero, ErrUnauthorized
	}
	val, ok := v.(*authContextValue[User])
	if !ok {
		return zero, fmt.Errorf("user ID type mismatch, got %T expected %T", v, zero)
	}
	if val.hasUser {
		return val.user, nil
	}
	if val.authFn == nil {
		return zero, ErrUnauthorized
	}
	val.once.Do(func() {
		user, err := val.authFn(ctx)
		if err != nil {
			val.authErr = fmt.Errorf("context auth: %w", err)
			return
		}
		val.user = user
		val.hasUser = true
	})
	if val.authErr != nil {
		return zero, val.authErr
	}
	return val.user, nil
}

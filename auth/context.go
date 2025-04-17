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
		mu      sync.RWMutex
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
	user, ok := val.cachedUser()
	if ok {
		return user, nil
	}
	return val.resolveUser(ctx)
}

func (acv *authContextValue[User]) cachedUser() (User, bool) {
	acv.mu.RLock()
	defer acv.mu.RUnlock()

	if acv.hasUser {
		return acv.user, true
	}
	var zero User
	return zero, false
}

func (acv *authContextValue[User]) resolveUser(ctx context.Context) (User, error) {
	var zero User
	if acv.authFn == nil { // authFn is only set in constructor so it's safe to check without locking
		return zero, ErrUnauthorized
	}
	acv.mu.Lock()
	defer acv.mu.Unlock()

	if acv.authErr != nil {
		return zero, acv.authErr
	}
	user, err := acv.authFn(ctx)
	if err != nil {
		acv.authErr = fmt.Errorf("context auth: %w", err)
		return zero, acv.authErr
	}
	acv.user = user
	acv.hasUser = true
	return user, nil
}

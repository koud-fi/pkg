package auth

import (
	"context"
	"fmt"
)

type userIDKey struct{}

func ContextWithUser[User any](
	parent context.Context, user User,
) context.Context {
	return context.WithValue(parent, userIDKey{}, user)
}

func ContextUser[User any](ctx context.Context) (User, error) {
	v := ctx.Value(userIDKey{})
	if v == nil {
		var zero User
		return zero, ErrUnauthorized
	}
	user, ok := v.(User)
	if !ok {
		var zero User
		return zero, fmt.Errorf("user ID type mismatch, got %T expected %T", v, zero)
	}
	return user, nil
}

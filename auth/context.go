package auth

import (
	"context"
	"fmt"
)

type userIDKey struct{}

func ContextWithUserID[UserID comparable](
	parent context.Context, userID UserID,
) context.Context {
	return context.WithValue(parent, userIDKey{}, userID)
}

func ContextUserID[UserID comparable](ctx context.Context) (UserID, error) {
	v := ctx.Value(userIDKey{})
	if v == nil {
		var zero UserID
		return zero, ErrUnauthorized
	}
	userID, ok := v.(UserID)
	if !ok {
		var zero UserID
		return zero, fmt.Errorf("user ID type mismatch, got %T expected %T", v, zero)
	}
	return userID, nil
}

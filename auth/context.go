package auth

import "context"

type identityKey struct{}

func ContextWithIdentity(parent context.Context, identity string) context.Context {
	if identity == "" {
		return parent
	}
	return context.WithValue(parent, identityKey{}, identity)
}

func Identity(ctx context.Context) (string, error) {
	if identity, ok := ctx.Value(identityKey{}).(string); ok {
		return identity, nil
	}
	return "", ErrUnauthorized
}

package password

import (
	"context"
	"fmt"

	"github.com/koud-fi/pkg/auth"
)

type (
	Authenticator[User any] struct {
		userLookup UserLookupFunc[User]
	}
	UserLookupFunc[User any] func(
		ctx context.Context, it auth.IdentityType, identity string,
	) (User, []Hash, error)
)

func NewAuthenticator[User any](
	userLookup UserLookupFunc[User],
) *Authenticator[User] {
	return &Authenticator[User]{userLookup: userLookup}
}

func (a *Authenticator[User]) Authenticate(ctx context.Context, payload auth.Payload) (User, error) {
	user, passwords, err := a.userLookup(ctx, payload.IdentityType, payload.Identity)
	if err != nil {
		return user, fmt.Errorf("lookup user: %w", err)
	}
	for _, proof := range payload.Proofs {
		if proof.Type != auth.Password {
			continue
		}
		for _, password := range passwords {
			if err := Compare(proof.Value, password); err != nil {
				continue
			}
			return user, nil
		}
	}
	return user, auth.ErrBadCredentials
}

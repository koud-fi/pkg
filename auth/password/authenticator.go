package password

import (
	"fmt"

	"github.com/koud-fi/pkg/auth"
)

type (
	Authenticator[UserID comparable] struct {
		userLookup UserLookupFunc[UserID]
	}
	UserLookupFunc[UserID comparable] func(
		it auth.IdentityType, identity string,
	) (UserID, []Password, error)
)

func NewAuthenticator[UserID comparable](
	userLookup UserLookupFunc[UserID],
) *Authenticator[UserID] {
	return &Authenticator[UserID]{userLookup: userLookup}
}

func (a *Authenticator[UserID]) Authenticate(payload auth.Payload) (UserID, error) {
	userID, passwords, err := a.userLookup(payload.IdentityType, payload.Identity)
	if err != nil {
		return userID, fmt.Errorf("lookup user: %w", err)
	}
	for _, proof := range payload.Proofs {
		if proof.Type != auth.Password {
			continue
		}
		for _, password := range passwords {
			if err := Compare(proof.Value, password); err != nil {
				continue
			}
			return userID, nil
		}
	}
	return userID, auth.ErrBadCredentials
}

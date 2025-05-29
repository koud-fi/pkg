package password

import (
	"github.com/koud-fi/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12

type Hash []byte

func NewHash(plain string) (Hash, error) {
	return bcrypt.GenerateFromPassword([]byte(plain), DefaultCost)
}

func NewHashWithCost(plain string, cost int) (Hash, error) {
	return bcrypt.GenerateFromPassword([]byte(plain), cost)
}

func Compare(plain string, to Hash) error {
	if err := bcrypt.CompareHashAndPassword(to, []byte(plain)); err != nil {
		return auth.ErrBadCredentials
	}
	return nil
}

func (Hash) String() string {
	return "************" // Avoid leaking passwords
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.String() + `"`), nil
}

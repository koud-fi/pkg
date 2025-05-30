package password

import (
	"github.com/koud-fi/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

type Hash []byte

func NewHash(plain string) (Hash, error) {
	return NewHashWithConfig(plain, DefaultConfig)
}

func NewHashWithConfig(plain string, conf Config) (Hash, error) {
	if err := conf.Validate(plain); err != nil {
		return nil, err
	}
	return bcrypt.GenerateFromPassword([]byte(plain), conf.Cost)
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

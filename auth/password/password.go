package password

import (
	"github.com/koud-fi/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 12

type Password []byte

func Hash(plain string) (Password, error) {
	return bcrypt.GenerateFromPassword([]byte(plain), defaultCost)
}

func Compare(plain string, to Password) error {
	if err := bcrypt.CompareHashAndPassword(to, []byte(plain)); err != nil {
		return auth.ErrBadCredentials
	}
	return nil
}

func (pw Password) String() string {
	return "************" // Avoid leaking passwords
}

func (pw Password) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pw.String() + `"`), nil
}

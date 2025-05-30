package password

import (
	"github.com/koud-fi/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

var DefaultConfig = Config{
	Cost: 12,
}

type Hash []byte

type Config struct {
	Cost int `json:"-"`
}

func (conf Config) NewHash(plain string) (Hash, error) {
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

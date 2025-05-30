package password

import (
	"errors"
	"fmt"

	"github.com/koud-fi/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost = bcrypt.MinCost
	MaxCost = bcrypt.MaxCost
)

var (
	DefaultConfig = Config{
		Cost:      12,
		MinLength: 8,
	}
	ErrPasswordTooShort = errors.New("password is too short")
)

type Hash []byte

type Config struct {
	Cost      int `json:"-"`
	MinLength int

	// TODO: more config options
}

func (conf Config) NewHash(plain string) (Hash, error) {
	if err := conf.Validate(plain); err != nil {
		return nil, err
	}
	return bcrypt.GenerateFromPassword([]byte(plain), conf.Cost)
}

func (conf Config) Validate(plain string) error {
	if len(plain) < conf.MinLength {
		return fmt.Errorf("%w: must be at least %d characters",
			ErrPasswordTooShort, conf.MinLength)
	}
	return nil
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

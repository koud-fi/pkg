package password

import (
	"github.com/koud-fi/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 12

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), defaultCost)
}

func Compare(password string, hash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return auth.ErrBadCredentials
	}
	return nil
}

// TODO

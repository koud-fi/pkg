package password

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost          = bcrypt.MinCost
	MaxCost          = bcrypt.MaxCost
	DefaultCost      = 12
	DefaultMaxLength = 72
)

var (
	DefaultConfig = Config{
		Cost:      DefaultCost,
		MinLength: 8, // Based on NIST SP 800-63B
		MaxLength: DefaultMaxLength,
	}
	ErrPasswordTooShort         = errors.New("password is too short")
	ErrPasswordTooLong          = errors.New("password is too long")
	ErrInsufficientUppercase    = errors.New("password requires more uppercase letters")
	ErrInsufficientLowercase    = errors.New("password requires more lowercase letters")
	ErrInsufficientDigits       = errors.New("password requires more digits")
	ErrInsufficientSpecialChars = errors.New("password requires more special characters")
	ErrTooManyRepeatingChars    = errors.New("password has too many repeating characters")
)

// Config for password hashing and validation, zero values are ignored.
type Config struct {
	Cost      int `json:"-"`          // bcrypt cost, between MinCost and MaxCost, default is DefaultCost
	MinLength int `json:",omitempty"` // Minimum number of characters required
	MaxLength int `json:",omitempty"` // Maximum number of bytes allowed

	RequireUppercase    int `json:",omitempty"` // Minimum number of uppercase letters required
	RequireLowercase    int `json:",omitempty"` // Minimum number of lowercase letters required
	RequireDigits       int `json:",omitempty"` // Minimum number of digits required
	RequireSpecialChars int `json:",omitempty"` // Minimum number of special characters required

	MaxRepeatingChars int `json:",omitempty"` // Maximum allowed consecutive repeating characters
}

func (conf Config) Validate(plain string) error {
	if len(plain) < conf.MinLength {
		return fmt.Errorf("%w: must be at least %d characters",
			ErrPasswordTooShort, conf.MinLength)
	}
	if conf.MaxLength > 0 && len([]byte(plain)) > conf.MaxLength {
		return fmt.Errorf("%w: must be at most %d bytes",
			ErrPasswordTooLong, conf.MaxLength)
	}

	// Count character types
	var uppercase, lowercase, digits, specialChars int
	for _, r := range plain {
		switch {
		case unicode.IsUpper(r):
			uppercase++
		case unicode.IsLower(r):
			lowercase++
		case unicode.IsDigit(r):
			digits++
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			specialChars++
		}
	}

	// Check character requirements
	if conf.RequireUppercase > 0 && uppercase < conf.RequireUppercase {
		return fmt.Errorf("%w: need at least %d, found %d",
			ErrInsufficientUppercase, conf.RequireUppercase, uppercase)
	}
	if conf.RequireLowercase > 0 && lowercase < conf.RequireLowercase {
		return fmt.Errorf("%w: need at least %d, found %d",
			ErrInsufficientLowercase, conf.RequireLowercase, lowercase)
	}
	if conf.RequireDigits > 0 && digits < conf.RequireDigits {
		return fmt.Errorf("%w: need at least %d, found %d",
			ErrInsufficientDigits, conf.RequireDigits, digits)
	}
	if conf.RequireSpecialChars > 0 && specialChars < conf.RequireSpecialChars {
		return fmt.Errorf("%w: need at least %d, found %d",
			ErrInsufficientSpecialChars, conf.RequireSpecialChars, specialChars)
	}

	// Check repeating characters
	if conf.MaxRepeatingChars > 0 {
		var (
			maxRepeating     = 0
			currentRepeating = 1
		)
		for i := 1; i < len(plain); i++ {
			if plain[i] == plain[i-1] {
				currentRepeating++
			} else {
				if currentRepeating > maxRepeating {
					maxRepeating = currentRepeating
				}
				currentRepeating = 1
			}
		}
		if currentRepeating > maxRepeating {
			maxRepeating = currentRepeating
		}
		if maxRepeating > conf.MaxRepeatingChars {
			return fmt.Errorf("%w: found %d consecutive repeating characters, maximum allowed is %d",
				ErrTooManyRepeatingChars, maxRepeating, conf.MaxRepeatingChars)
		}
	}
	return nil
}

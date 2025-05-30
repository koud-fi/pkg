package password

import (
	"errors"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		password string
		wantErr  error
	}{
		// Length validation tests
		{
			name:     "password too short",
			config:   Config{MinLength: 8, MaxLength: 20},
			password: "short",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "password too long",
			config:   Config{MinLength: 8, MaxLength: 10},
			password: "this_password_is_way_too_long",
			wantErr:  ErrPasswordTooLong,
		},
		{
			name:     "password length valid",
			config:   Config{MinLength: 8, MaxLength: 20},
			password: "validlength",
			wantErr:  nil,
		},

		// Character requirement tests
		{
			name:     "insufficient uppercase",
			config:   Config{RequireUppercase: 2},
			password: "passwordwithoneA",
			wantErr:  ErrInsufficientUppercase,
		},
		{
			name:     "sufficient uppercase",
			config:   Config{RequireUppercase: 2},
			password: "passwordWithTwoAB",
			wantErr:  nil,
		},
		{
			name:     "insufficient lowercase",
			config:   Config{RequireLowercase: 3},
			password: "PASSWORDwi",
			wantErr:  ErrInsufficientLowercase,
		},
		{
			name:     "sufficient lowercase",
			config:   Config{RequireLowercase: 3},
			password: "PASSWORDwith",
			wantErr:  nil,
		},
		{
			name:     "insufficient digits",
			config:   Config{RequireDigits: 2},
			password: "password1",
			wantErr:  ErrInsufficientDigits,
		},
		{
			name:     "sufficient digits",
			config:   Config{RequireDigits: 2},
			password: "password12",
			wantErr:  nil,
		},
		{
			name:     "insufficient special chars",
			config:   Config{RequireSpecialChars: 2},
			password: "password!",
			wantErr:  ErrInsufficientSpecialChars,
		},
		{
			name:     "sufficient special chars",
			config:   Config{RequireSpecialChars: 2},
			password: "password!@",
			wantErr:  nil,
		},

		// Repeating character tests
		{
			name:     "too many repeating chars",
			config:   Config{MaxRepeatingChars: 2},
			password: "passsword",
			wantErr:  ErrTooManyRepeatingChars,
		},
		{
			name:     "acceptable repeating chars",
			config:   Config{MaxRepeatingChars: 2},
			password: "password",
			wantErr:  nil,
		},
		{
			name:     "repeating chars at end",
			config:   Config{MaxRepeatingChars: 2},
			password: "passworddd",
			wantErr:  ErrTooManyRepeatingChars,
		},

		// Complex validation tests
		{
			name: "all requirements met",
			config: Config{
				MinLength:           12,
				MaxLength:           20,
				RequireUppercase:    1,
				RequireLowercase:    1,
				RequireDigits:       1,
				RequireSpecialChars: 1,
				MaxRepeatingChars:   2,
			},
			password: "Password123!",
			wantErr:  nil,
		},
		{
			name: "multiple requirements failed",
			config: Config{
				MinLength:           12,
				RequireUppercase:    2,
				RequireLowercase:    2,
				RequireDigits:       2,
				RequireSpecialChars: 2,
			},
			password: "Pass1!",
			wantErr:  ErrPasswordTooShort, // First error encountered
		},

		// Edge cases
		{
			name:     "empty password",
			config:   Config{MinLength: 1},
			password: "",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "zero requirements",
			config:   Config{},
			password: "anything",
			wantErr:  nil,
		},
		{
			name:     "unicode characters",
			config:   Config{RequireSpecialChars: 1},
			password: "password™",
			wantErr:  nil,
		},
		{
			name:     "zero max length means unrestricted",
			config:   Config{MaxLength: 0},
			password: "this_is_a_very_long_password_that_would_normally_be_rejected_if_there_was_a_length_limit",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate(tt.password)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err == nil {
				t.Errorf("Config.Validate() error = nil, wantErr %v", tt.wantErr)
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	// Test that DefaultConfig has sensible defaults
	if DefaultConfig.Cost != DefaultCost {
		t.Errorf("DefaultConfig.Cost = %d, want %d", DefaultConfig.Cost, DefaultCost)
	}
	if DefaultConfig.MinLength != 8 {
		t.Errorf("DefaultConfig.MinLength = %d, want 8", DefaultConfig.MinLength)
	}
	if DefaultConfig.MaxLength != DefaultMaxLength {
		t.Errorf("DefaultConfig.MaxLength = %d, want %d", DefaultConfig.MaxLength, DefaultMaxLength)
	}

	// Test that a reasonable password passes default validation
	err := DefaultConfig.Validate("password123")
	if err != nil {
		t.Errorf("DefaultConfig.Validate() failed for reasonable password: %v", err)
	}
}

func TestConfig_Validate_CharacterCounting(t *testing.T) {
	// Test specific character counting edge cases
	tests := []struct {
		name     string
		config   Config
		password string
		wantErr  error
	}{
		{
			name:     "mixed case and numbers",
			config:   Config{RequireUppercase: 1, RequireLowercase: 1, RequireDigits: 1},
			password: "Aa1",
			wantErr:  nil,
		},
		{
			name:     "punctuation vs symbols",
			config:   Config{RequireSpecialChars: 3},
			password: "password!@#",
			wantErr:  nil,
		},
		{
			name:     "unicode special characters",
			config:   Config{RequireSpecialChars: 1},
			password: "password©",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate(tt.password)
			if tt.wantErr == nil && err != nil {
				t.Errorf("Config.Validate() error = %v, wantErr nil", err)
			} else if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package oauth

import (
	"errors"
	"fmt"
)

const (
	BearerToken TokenType = "Bearer"
	MACToken    TokenType = "MAC"
	BasicToken  TokenType = "Basic"
)

type (
	Token struct {
		AccessToken  string    `json:"access_token"`
		ExpiresIn    int       `json:"expires_in"`
		ExtExpiresIn int       `json:"ext_expires_in,omitempty"` // microsoft specific
		RefreshToken string    `json:"refresh_token,omitempty"`
		Scope        string    `json:"scope"`
		TokenType    TokenType `json:"token_type"`
		IDToken      string    `json:"id_token,omitempty"`
	}
	TokenType string
)

func (t Token) HTTPHeader() (string, error) {
	if t.AccessToken == "" {
		return "", errors.New("empty access token")
	}
	switch t.TokenType {
	case BearerToken:
		return string(t.TokenType) + " " + t.AccessToken, nil
	default:
		return "", fmt.Errorf("unsupported token type: %s", t.TokenType)
	}
}

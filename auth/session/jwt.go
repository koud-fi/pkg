package session

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/koud-fi/pkg/auth"
)

const defaultJWTLifetime = time.Hour

type (
	JWTAuthenticator[UserID comparable] struct {
		issuer   string
		secret   string
		userIDFn JWTUserIDFunc[UserID]
	}
	JWTUserIDFunc[UserID comparable] func(jwt.MapClaims) (UserID, error)
)

var _ auth.Authenticator[any] = &JWTAuthenticator[any]{}

func NewJWTAuthenticator[UserID comparable](
	issuer string, secret string, userIDFn JWTUserIDFunc[UserID],
) *JWTAuthenticator[UserID] {
	return &JWTAuthenticator[UserID]{
		issuer:   issuer,
		secret:   secret,
		userIDFn: userIDFn,
	}
}

func (a *JWTAuthenticator[UserID]) NewToken(userID UserID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": a.issuer,
			"sub": userID,
			"exp": time.Now().Add(defaultJWTLifetime).Unix(),

			// TODO: add more claims to enhance security

		})
	singedToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return singedToken, nil
}

func (a *JWTAuthenticator[UserID]) Authenticate(
	_ context.Context, payload auth.Payload,
) (UserID, error) {
	var zero UserID
	for _, proof := range payload.Proofs {
		if proof.Type != auth.JWT {
			continue
		}
		token, err := jwt.Parse(proof.Value, func(token *jwt.Token) (any, error) {
			issuer, err := token.Claims.GetIssuer()
			if err != nil {
				return nil, fmt.Errorf("get issuer: %w", err)
			}
			if issuer != a.issuer {
				return nil, fmt.Errorf("invalid issuer: %s, expected %s", issuer, a.issuer)
			}
			return []byte(a.secret), nil
		})
		if err != nil {
			return zero, fmt.Errorf("parse token: %w", err)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return zero, fmt.Errorf("invalid token claims")
		}
		userID, err := a.userIDFn(claims)
		if err != nil {
			return zero, fmt.Errorf("get user ID: %w", err)
		}
		return userID, nil
	}
	return zero, auth.ErrBadCredentials
}

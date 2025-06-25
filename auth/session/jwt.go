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
	JWTConfig struct {
		expireTime time.Duration
	}

	JWTAuthenticator[User any] struct {
		issuer string
		secret string
		userFn JWTUserFunc[User]
		JWTConfig
	}
	JWTUserFunc[User any] func(context.Context, jwt.Claims) (User, error)
	JWTOption             func(*JWTConfig)
)

var _ auth.Authenticator[any] = &JWTAuthenticator[any]{}

func WithExpireTime(expireTime time.Duration) JWTOption {
	return func(c *JWTConfig) {
		c.expireTime = expireTime
	}
}

func NewJWTAuthenticator[User any](
	issuer string, secret string, userFn JWTUserFunc[User], opts ...JWTOption,
) *JWTAuthenticator[User] {
	config := JWTConfig{
		expireTime: defaultJWTLifetime,
	}
	for _, opt := range opts {
		opt(&config)
	}

	return &JWTAuthenticator[User]{
		issuer:    issuer,
		secret:    secret,
		userFn:    userFn,
		JWTConfig: config,
	}
}

func (a *JWTAuthenticator[_]) NewToken(subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": a.issuer,
			"sub": subject,
			"exp": time.Now().Add(a.expireTime).Unix(),

			// TODO: Add more claims to enhance security

		})
	singedToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return singedToken, nil
}

func (a *JWTAuthenticator[User]) Authenticate(
	ctx context.Context, payload auth.Payload,
) (User, error) {
	var zero User
	for _, proof := range payload.Proofs {
		if proof.Type != auth.Token {
			continue
		}
		token, err := jwt.Parse(proof.Value, func(token *jwt.Token) (any, error) {

			// TOOD: Does this actually work?

			return []byte(a.secret), nil
		})
		if err != nil {
			return zero, fmt.Errorf("parse token: %w", err)
		}

		// TODO: Improve security and validation of the token

		user, err := a.userFn(ctx, token.Claims)
		if err != nil {
			return zero, fmt.Errorf("get user: %w", err)
		}
		return user, nil
	}
	return zero, auth.ErrBadCredentials
}

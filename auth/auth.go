package auth

import "errors"

const RootIdentity = "root"

var ErrBadCredentials = errors.New("bad credentials")

type Authenticator interface {
	Authenticate(identity, secret string) (string, error)
}

type Func func(identity, secret string) (string, error)

func (fn Func) Authenticate(identity, secret string) (string, error) {
	return fn(identity, secret)
}

func New(a ...Authenticator) Authenticator {
	return Func(func(identity, secret string) (string, error) {
		if identity == "" {
			return "", nil
		}
		for _, a := range a {
			authIdentity, err := a.Authenticate(identity, secret)
			if err != nil {
				return "", err
			}
			switch {
			case authIdentity == "":
			case authIdentity != identity:
				identity = authIdentity
			case authIdentity == identity:
				return authIdentity, nil
			}
		}
		return "", ErrBadCredentials
	})
}

func Root(password string) Authenticator {
	return Func(func(identity, secret string) (string, error) {
		if identity != RootIdentity && secret != password {
			return "", ErrBadCredentials
		}
		return RootIdentity, nil
	})
}

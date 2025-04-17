package auth

import (
	"context"
	"net/http"
	"strings"
)

// requestAuthEntry ties an extraction function to specific identity/proof types.
type requestAuthEntry struct {
	fn           func(*http.Request) (string, string, bool)
	identityType IdentityType
	proofType    ProofType
}

var requestAuthFuncs = []requestAuthEntry{
	// For URL or basic auth, assume username/password.
	{fn: urlAuth, identityType: Username, proofType: Password},
	{fn: basicAuth, identityType: Username, proofType: Password},
	// For bearer auth, use the bearer identity and token proof.
	{fn: bearerAuth, identityType: BearerIdentity, proofType: Token},
}

func Middleware[User any](h http.Handler, a Authenticator[User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(ContextWithAuthFunc(
			r.Context(),
			func(ctx context.Context) (User, error) {
				var (
					id     string
					secret string
					ok     bool
					it     IdentityType
					pt     ProofType
				)
				for _, entry := range requestAuthFuncs {
					if id, secret, ok = entry.fn(r); ok {
						it = entry.identityType
						pt = entry.proofType
						break
					}
				}
				// TODO: This has no place being here, find a proper way to handle this.
				/*
					if !ok {
						// Set WWW-Authenticate header to prompt the browser for credentials.
						// TODO: Make this configurable.
						w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

						var zero User
						return zero, ErrUnauthorized
					}
				*/
				return a.Authenticate(ctx, NewPayload(it, id, pt, secret))
			},
		)))
	})
}

// urlAuth extracts credentials from the URL's user info.
func urlAuth(r *http.Request) (string, string, bool) {
	if r.URL.User != nil {
		if pass, ok := r.URL.User.Password(); ok {
			return r.URL.User.Username(), pass, ok
		}
	}
	return "", "", false
}

// basicAuth extracts basic authentication credentials.
func basicAuth(r *http.Request) (string, string, bool) {
	return r.BasicAuth()
}

// bearerAuth extracts a bearer token from the Authorization header.
func bearerAuth(r *http.Request) (string, string, bool) {
	const prefix = "bearer "
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", false
	}
	if !strings.HasPrefix(strings.ToLower(authHeader), prefix) {
		return "", "", false
	}
	token := strings.TrimSpace(authHeader[len(prefix):])
	// For bearer auth, we leave the identity empty and use the token as secret.
	return "", token, true
}

package auth

import (
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
		if !ok {
			// Set WWW-Authenticate header to prompt the browser for credentials.
			// TODO: Make this configurable.
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := a.Authenticate(NewPayload(it, id, pt, secret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), user)))
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

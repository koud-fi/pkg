package auth

import (
	"net/http"
	"strings"
)

var requestAuthFuncs = []func(*http.Request) (string, string, bool){
	urlAuth, basicAuth, bearerAuth,
}

func Middleware(h http.Handler, a Authenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var identity, secret string
		for _, fn := range requestAuthFuncs {
			var ok bool
			if identity, secret, ok = fn(r); ok {
				break
			}
		}
		identity, err := a.Authenticate(identity, secret)
		if err != nil {

			// TODO: custom error handlers

			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.ServeHTTP(w, r.WithContext(ContextWithIdentity(r.Context(), identity)))
	})
}

func urlAuth(r *http.Request) (string, string, bool) {
	if r.URL.User != nil {
		if pass, ok := r.URL.User.Password(); ok {
			return r.URL.User.Username(), pass, ok
		}
	}
	return "", "", false
}

func basicAuth(r *http.Request) (string, string, bool) { return r.BasicAuth() }

func bearerAuth(r *http.Request) (string, string, bool) {
	const prefix = "bearer "
	auth := strings.ToLower(r.Header.Get("Authorization"))
	if !strings.HasPrefix(auth, prefix) {
		return "", "", false
	}
	return BearerIdentity, strings.TrimPrefix(auth, prefix), true
}

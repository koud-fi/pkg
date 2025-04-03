package oauth

import (
	"net/http"
)

const HTTPAuthPath = "auth"
const HTTPCallbackPath = "callback"

func HTTPHandler(
	conf Config,
	handlerFn func(http.ResponseWriter, *http.Request, Token),
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/"+HTTPAuthPath, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, RedirectURL(conf), http.StatusFound)
	})
	mux.HandleFunc("/"+HTTPCallbackPath, func(w http.ResponseWriter, r *http.Request) {
		token, err := ParseCallback(conf, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		handlerFn(w, r, token)
	})
	return mux
}

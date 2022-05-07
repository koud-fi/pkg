package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/koud-fi/pkg/proc"
	"github.com/koud-fi/pkg/proc/router"
)

type Server struct {
	r router.Router
}

func New(r router.Router) Server {
	return Server{r: r}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	switch r.Method {
	case http.MethodHead, http.MethodGet:
	default:
		endpoint = r.Method + endpoint
	}
	var params proc.Params
	switch contentTypeBase(r) {
	case "text/json", "application/json":
		params = proc.ParamFunc(json.NewDecoder(r.Body).Decode)

	// TODO: multi-part form

	default:
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		params = proc.ParamMap(r.Form)
	}
	out, err := s.r.Invoke(r.Context(), endpoint, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	proc.WriteOutput(w, out)
}

func contentTypeBase(r *http.Request) string {
	return strings.SplitN(r.Header.Get("Content-Type"), " ", 2)[0]
}

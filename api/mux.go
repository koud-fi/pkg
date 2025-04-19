package api

import "net/http"

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{ServeMux: http.NewServeMux()}
}

// TODO: methods

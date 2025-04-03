package endpoint_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koud-fi/pkg/assert"
	"github.com/koud-fi/pkg/endpoint"
	"github.com/koud-fi/pkg/fetch"
)

type HelloInput struct {
	Name string
}

func Hello(_ context.Context, in HelloInput) (string, error) {
	if in.Name == "" {
		in.Name = "Teppo"
	}
	return fmt.Sprintf("Hello, %s!", in.Name), nil
}

func TestEndpoint(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/hello", endpoint.New(Hello))

	rrec := httptest.NewRecorder()

	mux.ServeHTTP(rrec, assert.Must(fetch.Get("/hello").HttpRequest()))
	if rrec.Body.String() != "Hello, Teppo!" {
		t.Fatalf("expected body to be 'Hello, Teppo!', got '%s'", rrec.Body.String())
	}
}

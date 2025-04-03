package endpoint_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

	// no arguments
	rrec := httptest.NewRecorder()
	req, err := fetch.Get("/hello").HttpRequest()
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	mux.ServeHTTP(rrec, req)
	if rrec.Body.String() != "Hello, Teppo!" {
		t.Fatalf("expected body to be 'Hello, Teppo!', got '%s'", rrec.Body.String())
	}

	// url query arguments
	rrec = httptest.NewRecorder()
	req, err = fetch.Get("/hello?name=Seppo").HttpRequest()
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	mux.ServeHTTP(rrec, req)
	if rrec.Body.String() != "Hello, Seppo!" {
		t.Fatalf("expected body to be 'Hello, Seppo!', got '%s'", rrec.Body.String())
	}
}

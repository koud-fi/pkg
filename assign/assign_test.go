package assign_test

import (
	"reflect"
	"testing"

	"github.com/koud-fi/pkg/assign"
)

type Address struct {
	Street string
	City   string
	Zip    int
}

type Friend struct {
	Name  string
	Email string
}

type ComplexDemo struct {
	Title      string
	Subtitle   *string
	Count      int
	Price      float64
	Active     bool
	Tags       []string
	Attributes map[string]string
	Address    Address
	Friends    []Friend
	Scores     []float64
	Metadata   map[string]any
}

func TestValue(t *testing.T) {
	// Input demonstrates various conversion cases:
	// - Some fields are provided as JSON strings.
	// - Some as native maps.
	// - Some require basic conversion from strings.
	input := map[string]any{
		"Title":    "Complex Example",
		"Subtitle": "Subtitle Example",
		"Count":    "7",
		"Price":    "19.99",
		"Active":   "true",
		"Tags":     `["go", "reflection", "conversion"]`,
		"Attributes": map[string]any{
			"version": "1.0",
			"env":     "production",
		},
		"Address": map[string]any{
			"Street": "123 Main St",
			"City":   "Gopher City",
			"Zip":    "12345",
		},
		"Friends": `[{"Name": "Alice", "Email": "alice@example.com"},
		             {"Name": "Bob", "Email": "bob@example.com"}]`,
		"Scores":   "[95.5, 82.3, 76]",
		"Metadata": `{"key1": "value1", "key2": 2}`,
	}
	var demo ComplexDemo
	if err := assign.Value(&demo, input); err != nil {
		t.Fatalf("Unmarshal failed: %+v", err)
	}
	if demo.Title != "Complex Example" {
		t.Errorf("Expected Title %q, got %q", "Complex Example", demo.Title)
	}
	if demo.Count != 7 {
		t.Errorf("Expected Count 7, got %d", demo.Count)
	}
	if demo.Price != 19.99 {
		t.Errorf("Expected Price 19.99, got %f", demo.Price)
	}
	if !demo.Active {
		t.Error("Expected Active true")
	}
	if !reflect.DeepEqual(demo.Tags, []string{"go", "reflection", "conversion"}) {
		t.Errorf("Tags mismatch: %v", demo.Tags)
	}
	if demo.Attributes["env"] != "production" {
		t.Errorf("Expected Attributes.env %q, got %q", "production", demo.Attributes["env"])
	}
	if demo.Address.City != "Gopher City" {
		t.Errorf("Expected Address.City %q, got %q", "Gopher City", demo.Address.City)
	}
	if len(demo.Friends) != 2 {
		t.Errorf("Expected 2 Friends, got %d", len(demo.Friends))
	}
	if len(demo.Scores) != 3 {
		t.Errorf("Expected 3 Scores, got %d", len(demo.Scores))
	}
	if demo.Metadata["key1"] != "value1" {
		t.Errorf("Expected Metadata key1 %q, got %v", "value1", demo.Metadata["key1"])
	}
}

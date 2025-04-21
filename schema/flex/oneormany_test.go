package flex

import (
	"encoding/json"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

// Test with string values using JSON.
func TestOneOrManyStrings_JSON(t *testing.T) {
	var o OneOrMany[string]

	// Unmarshal a single string.
	if err := json.Unmarshal([]byte(`"hello"`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected := OneOrMany[string]{"hello"}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	// Marshal: with a single element it should be a bare string.
	b, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `"hello"` {
		t.Errorf("Expected \"hello\", got %s", b)
	}

	// Unmarshal multiple strings.
	if err := json.Unmarshal([]byte(`["hello", "world"]`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected = OneOrMany[string]{"hello", "world"}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	// Marshal multiple elements should be an array.
	b, err = json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `["hello","world"]` {
		t.Errorf("Expected [\"hello\",\"world\"], got %s", b)
	}
}

// Test with numeric values using JSON.
func TestOneOrManyInts_JSON(t *testing.T) {
	var o OneOrMany[int]

	// Unmarshal a single number.
	if err := json.Unmarshal([]byte(`1`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected := OneOrMany[int]{1}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	// Marshal single element.
	b, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `1` {
		t.Errorf("Expected 1, got %s", b)
	}

	// Unmarshal multiple numbers.
	if err := json.Unmarshal([]byte(`[1,2,3]`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected = OneOrMany[int]{1, 2, 3}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	// Marshal multiple numbers.
	b, err = json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `[1,2,3]` {
		t.Errorf("Expected [1,2,3], got %s", b)
	}
}

// A simple struct to test object unmarshaling.
type Person struct {
	Name string `json:"name"`
}

// Test with objects using JSON.
func TestOneOrManyStructs_JSON(t *testing.T) {
	var o OneOrMany[Person]

	// Unmarshal a single object.
	if err := json.Unmarshal([]byte(`{"name": "Alice"}`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected := OneOrMany[Person]{{Name: "Alice"}}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	// Marshal a single object.
	b, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `{"name":"Alice"}` {
		t.Errorf("Expected {\"name\":\"Alice\"}, got %s", b)
	}

	// Unmarshal multiple objects.
	if err := json.Unmarshal([]byte(`[{"name": "Alice"}, {"name": "Bob"}]`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected = OneOrMany[Person]{{Name: "Alice"}, {Name: "Bob"}}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	// Marshal multiple objects.
	b, err = json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `[{"name":"Alice"},{"name":"Bob"}]` {
		t.Errorf("Expected array JSON, got %s", b)
	}
}

// Test with nested OneOrMany types using JSON.
func TestOneOrManyNested_JSON(t *testing.T) {
	// Outer type is OneOrMany of OneOrMany[string].
	var o OneOrMany[OneOrMany[string]]

	// When the JSON is an array of strings, it should be interpreted as a single inner OneOrMany.
	if err := json.Unmarshal([]byte(`["a", "b"]`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected := OneOrMany[OneOrMany[string]]{OneOrMany[string]{"a", "b"}}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	// Marshal: since there is one element, it will marshal as a single value.
	b, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	// For a single inner element, it collapses to the inner value.
	if string(b) != `["a","b"]` {
		t.Errorf("Expected [\"a\",\"b\"], got %s", b)
	}

	// Now test with explicit nested arrays.
	if err := json.Unmarshal([]byte(`[["a"], ["b", "c"]]`), &o); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	expected = OneOrMany[OneOrMany[string]]{
		OneOrMany[string]{"a"},
		OneOrMany[string]{"b", "c"},
	}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	b, err = json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	// Note: the first inner value collapses, so we expect "a" not ["a"].
	if string(b) != `["a",["b","c"]]` {
		t.Errorf("Expected nested array JSON, got %s", b)
	}
}

// Test with strings using YAML.
func TestOneOrManyStrings_YAML(t *testing.T) {
	var o OneOrMany[string]

	// YAML single string.
	yamlData := `hello`
	if err := yaml.Unmarshal([]byte(yamlData), &o); err != nil {
		t.Fatalf("YAML Unmarshal error: %v", err)
	}
	expected := OneOrMany[string]{"hello"}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	b, err := yaml.Marshal(o)
	if err != nil {
		t.Fatalf("YAML Marshal error: %v", err)
	}
	// YAML output for a single scalar usually ends with a newline.
	if string(b) != "hello\n" {
		t.Errorf("Expected 'hello', got %s", b)
	}

	// YAML multiple strings.
	yamlData = `
- hello
- world
`
	if err := yaml.Unmarshal([]byte(yamlData), &o); err != nil {
		t.Fatalf("YAML Unmarshal error: %v", err)
	}
	expected = OneOrMany[string]{"hello", "world"}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected %v, got %v", expected, o)
	}

	b, err = yaml.Marshal(o)
	if err != nil {
		t.Fatalf("YAML Marshal error: %v", err)
	}
	if string(b) != "- hello\n- world\n" {
		t.Errorf("Expected sequence YAML, got %s", b)
	}
}

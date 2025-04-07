package kvtest

import (
	"context"
	"fmt"

	"github.com/koud-fi/pkg/kv"
)

// TestKV tests a kv.Storage[string, string] implementation.
// It returns an error if any test assertion fails.
func TestKV(s kv.Storage[string, string]) error {
	ctx := context.Background()

	// Setup keys.
	testData := map[string]string{
		"a":         "av",
		"bb":        "bv",
		"dddd":      "dv",
		"ccc":       "cv",
		"eeeee":     "ev",
		"ffffff":    "fv",
		"gggg/ggg":  "gv",
		"hhhh/hhhh": "hv",
	}

	// Insert key/value pairs.
	for k, v := range testData {
		if err := s.Set(ctx, k, v); err != nil {
			return fmt.Errorf("Set failed for key %q: %w", k, err)
		}
	}

	// Delete one key.
	if err := s.Del(ctx, "bb"); err != nil {
		return fmt.Errorf("Del failed for key 'bb': %w", err)
	}

	// Get test: existing key.
	if v, err := s.Get(ctx, "a"); err != nil {
		return fmt.Errorf("Get failed for key 'a': %w", err)
	} else if v != "av" {
		return fmt.Errorf("unexpected value for key 'a': got %q, want %q", v, "av")
	}

	// Get test: deleted key should return ErrNotFound.
	if v, err := s.Get(ctx, "bb"); err != kv.ErrNotFound {
		return fmt.Errorf("expected ErrNotFound for key 'bb', got value %q and error %v", v, err)
	}

	// Get test: another existing key.
	if v, err := s.Get(ctx, "gggg/ggg"); err != nil {
		return fmt.Errorf("Get failed for key 'gggg/ggg': %w", err)
	} else if v != "gv" {
		return fmt.Errorf("unexpected value for key 'gggg/ggg': got %q, want %q", v, "gv")
	}

	// Test Scan if the storage implements the ScanReader interface.
	scanner, ok := s.(kv.ScanReader[string, string])
	if ok {
		seq, errFn := scanner.Scan(ctx)

		// Expected scan result: the testData without the deleted key "bb".
		expected := map[string]string{
			"a":         "av",
			"dddd":      "dv",
			"ccc":       "cv",
			"eeeee":     "ev",
			"ffffff":    "fv",
			"gggg/ggg":  "gv",
			"hhhh/hhhh": "hv",
		}
		got := make(map[string]string)
		for pair := range seq {
			got[pair.Key()] = pair.Value()
		}
		if len(got) != len(expected) {
			return fmt.Errorf("scan: expected %d elements, got %d", len(expected), len(got))
		}
		for k, want := range expected {
			if got[k] != want {
				return fmt.Errorf("scan: for key %q, expected %q, got %q", k, want, got[k])
			}
		}
		if err := errFn(); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
	}
	return nil
}

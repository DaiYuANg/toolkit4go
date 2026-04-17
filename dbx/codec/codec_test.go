package codec

import (
	"errors"
	"fmt"
	"testing"
)

func TestUnknownError(t *testing.T) {
	err := &UnknownError{Name: "csv"}
	if !errors.Is(err, ErrUnknown) {
		t.Fatal("errors.Is(err, ErrUnknown) should be true")
	}

	wrapped := fmt.Errorf("mapper init: %w", err)
	var target *UnknownError
	if !errors.As(wrapped, &target) {
		t.Fatal("errors.As should succeed on wrapped error")
	}
	if target.Name != "csv" {
		t.Fatalf("expected Name=%q, got %q", "csv", target.Name)
	}
}

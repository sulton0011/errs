package errs

import (
	"testing"
)

func TestNew(t *testing.T) {
	message := "An error occurred"
	err := New(message)

	if err == nil {
		t.Fatal("Expected a non-nil error")
	}

	if err.Error() != message {
		t.Fatalf("Expected error message to be '%s', got '%s'", message, err.Error())
	}

	if Unwrap(err) != message {
		t.Fatalf("Expected unwrapped message to be '%s', got '%s'", message, Unwrap(err))
	}
}

func TestIs(t *testing.T) {
	originalErr := New("Test error")

	if !Is(originalErr, originalErr) {
		t.Fatal("Expected Is to return true for the same error")
	}

	if Is(originalErr, nil) {
		t.Fatal("Expected Is to return false when target is nil and error is not")
	}

	if !Is(nil, nil) {
		t.Fatal("Expected Is to return true for two nil errors")
	}
}

func TestIsNil(t *testing.T) {
	var nilErr error // This is a nil error

	if !IsNil(nilErr) {
		t.Fatal("Expected IsNil to return true for a nil error")
	}

	notNilErr := New("Test error")
	if IsNil(notNilErr) {
		t.Fatal("Expected IsNil to return false for a non-nil error")
	}
}
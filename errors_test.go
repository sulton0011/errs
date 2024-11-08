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

func TestNew_NilMessage(t *testing.T) {
	err := New("")
	if err == nil {
		t.Fatal("Expected a non-nil error")
	}

	if err.Error() != "" {
		t.Fatalf("Expected error message to be empty, got '%s'", err.Error())
	}

	if Unwrap(err) != "" {
		t.Fatalf("Expected unwrapped message to be empty, got '%s'", Unwrap(err))
	}
}

func TestNewF_NilFormat(t *testing.T) {
	err := NewF("")
	if err == nil {
		t.Fatal("Expected a non-nil error")
	}

	if err.Error() != "" {
		t.Fatalf("Expected formatted error message to be empty, got '%s'", err.Error())
	}
}

func TestError_NilErrorString(t *testing.T) {
	var e *errorString
	if e.Error() != "" {
		t.Fatalf("Expected Error() on nil *errorString to return empty string, got '%s'", e.Error())
	}
}

func TestIs_NilTarget(t *testing.T) {
	// Test case where both errors are nil
	if !Is(nil, nil) {
		t.Fatal("Expected Is(nil, nil) to return true")
	}

	// Test case where only the target is nil
	err := New("Test error")
	if Is(err, nil) {
		t.Fatal("Expected Is(err, nil) to return false when err is not nil")
	}
}

func TestIs_ErrorNil(t *testing.T) {
	// Test case where error is nil but target is not
	target := New("Test error")
	if Is(nil, target) {
		t.Fatal("Expected Is(nil, target) to return false when target is not nil")
	}
}

func TestIsNil_NilError(t *testing.T) {
	if !IsNil(nil) {
		t.Fatal("Expected IsNil(nil) to return true for a nil error")
	}
}

func TestIsNil_NonNilError(t *testing.T) {
	err := New("Test error")
	if IsNil(err) {
		t.Fatal("Expected IsNil(err) to return false for a non-nil error")
	}
}

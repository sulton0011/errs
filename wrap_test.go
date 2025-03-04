package errs

import (
	"testing"
)

func TestWrap(t *testing.T) {
	// Test wrapping a nil error
	var nilErr error
	result := Wrap(nilErr, "Context message")
	if result != nil {
		t.Fatal("Expected nil result when wrapping a nil error")
	}

	// Test wrapping a non-nil error
	originalErr := New("Original error")
	wrappedErr := Wrap(originalErr, "Additional context")

	if wrappedErr == nil {
		t.Fatal("Expected a non-nil wrapped error")
	}

	expectedMessage := "Original error"
	if wrappedErr.Error() != expectedMessage {
		t.Fatalf("Expected error message to be '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}

	expectedMessage = "Additional context ---> Original error"
	if Unwrap(wrappedErr) != expectedMessage {
		t.Fatalf("Expected unwrapped message to be '%s', got '%s'", expectedMessage, Unwrap(wrappedErr))
	}
}

func TestWrapF(t *testing.T) {
	// Test wrapping a nil error with formatting
	var nilErr error
	result := WrapF(nilErr, "Formatted error: %s", "extra info")
	if result != nil {
		t.Fatal("Expected nil result when wrapping a nil error")
	}

	// Test wrapping a non-nil error with formatting
	originalErr := New("Original error")
	wrappedErr := WrapF(originalErr, "Formatted error: %s", "additional info")

	if wrappedErr == nil {
		t.Fatal("Expected a non-nil wrapped error")
	}

	expectedMessage := "Original error"
	if wrappedErr.Error() != expectedMessage {
		t.Fatalf("Expected error message to be '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}

	expectedMessage = "Formatted error: additional info ---> Original error"
	if Unwrap(wrappedErr) != expectedMessage {
		t.Fatalf("Expected unwrapped message to be '%s', got '%s'", expectedMessage, Unwrap(wrappedErr))
	}
}

func TestLog(t *testing.T) {
	// Test logging a nil error
	req := struct{}{}
	var nilErr error
	Log(nilErr, req, "This should not log anything")

	// Test logging a non-nil error
	err := New("Test error")
	Log(err, req, "Logging test error with request context")
	// Here we are not verifying the output because logging is asynchronous.
	// You can enhance this by implementing a way to capture logs in tests if necessary.
}

func TestUnwrap(t *testing.T) {
	// Test unwrapping a non-wrapped error
	originalErr := New("Original error")
	unwrappedMessage := Unwrap(originalErr)

	if unwrappedMessage != originalErr.Error() {
		t.Fatalf("Expected unwrapped message to be '%s', got '%s'", originalErr.Error(), unwrappedMessage)
	}

	// Test unwrapping a wrapped error
	wrappedErr := Wrap(originalErr, "Context")
	unwrappedMessage = Unwrap(wrappedErr)

	if unwrappedMessage == originalErr.Error() {
		t.Fatalf("Expected unwrapped message to be '%s', got '%s'", originalErr.Error(), unwrappedMessage)
	}

	expectedMessage := "Context ---> Original error"
	if unwrappedMessage != expectedMessage {
		t.Fatalf("Expected unwrapped message to be '%s', got '%s'", expectedMessage, unwrappedMessage)
	}
}

func TestUnwrapE(t *testing.T) {
	// Test unwrapping a non-wrapped error
	originalErr := New("Original error")
	unwrappedErr := UnwrapE(originalErr)

	if unwrappedErr.Error() != originalErr.Error() {
		t.Fatalf("Expected unwrapped error to be '%v', got '%v'", originalErr, unwrappedErr)
	}

	// Test unwrapping a wrapped error
	wrappedErr := Wrap(originalErr, "Context")
	unwrappedErr = UnwrapE(wrappedErr)

	if unwrappedErr == nil {
		t.Fatal("Expected unwrapped error to be non-nil when wrapped error is provided")
	}

	expectedMessage := "Context ---> Original error"
	if unwrappedErr.Error() != expectedMessage {
		t.Fatalf("Expected unwrapped error message to be '%s', got '%s'", expectedMessage, unwrappedErr.Error())
	}

	// Test unwrapping a wrapped error with an empty message
	emptyWrappedErr := Wrap(originalErr, "")
	unwrappedErr = UnwrapE(emptyWrappedErr)

	if unwrappedErr == nil {
		t.Fatalf("Expected unwrapped error to be nil when wrapped error has an empty message")
	}
}

func TestWrap_NilError(t *testing.T) {
	err := Wrap(nil, "additional context")
	if err != nil {
		t.Errorf("Wrap(nil) = %v, want nil", err)
	}
}

func TestWrapF_NilError(t *testing.T) {
	err := WrapF(nil, "failed to %s", "execute")
	if err != nil {
		t.Errorf("WrapF(nil) = %v, want nil", err)
	}
}

func TestUnwrap_NilError(t *testing.T) {
	result := Unwrap(nil)
	if result != "" {
		t.Errorf("Unwrap(nil) = %v, want empty string", result)
	}
}

func TestUnwrapE_NilError(t *testing.T) {
	err := UnwrapE(nil)
	if err != nil {
		t.Errorf("UnwrapE(nil) = %v, want nil", err)
	}
}

func TestLog_NilError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Log(nil) panicked, want no panic")
		}
	}()

	Log(nil, "request", "some message")
}

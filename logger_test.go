package errs

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf)

	logger.Error("Test error", slog.String("key", "value"))

	if buf.Len() == 0 {
		t.Fatal("Expected logs to be written")
	}
}

func TestNewTextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := newTextLogger(&buf)

	logger.Error("Test text error", slog.String("key", "value"))

	if buf.Len() == 0 {
		t.Fatal("Expected logs to be written")
	}
}

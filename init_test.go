package errs

import (
	"os"
	"testing"
)

func TestSetLogTypes(t *testing.T) {
	// Redirect output to a temp file
	tempFile, err := os.CreateTemp("", "logfile.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Test setting log types
	SetLogTypes(LogTypeJSON, LogTypeText, LogTypeFile)

	if len(slogLoggers) != 3 {
		t.Fatalf("Expected 3 loggers, got %d", len(slogLoggers))
	}
}

func TestSetLogFile(t *testing.T) {
	err := SetLogFile("test.log")
	if err != nil {
		t.Fatal("Expected no error, got", err)
	}

	// Check if the log file was created
	if _, err := os.Stat("test.log"); os.IsNotExist(err) {
		t.Fatal("Expected log file to exist")
	}

	// Clean up
	os.Remove("test.log")
}

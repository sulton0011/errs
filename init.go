package errs

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// Define custom logging levels.
type logType string

const (
	LogTypeJSON logType = "JSON" // JSON logging format.
	LogTypeText logType = "TEXT" // Text logging format.
	LogTypeFile logType = "FILE" // File logging format.

	DefaultLogFile = "log/logger.json" // Default log file path.
)

// Variables to manage the loggers and logging levels.
var (
	fileLogger  *os.File
	slogLoggers []*slog.Logger // List of loggers.
)

// Initializes loggers when the package is first loaded.
func init() {
	SetLogTypes(LogTypeJSON)
}

// SetLogTypes configures the logging types to be used (e.g., JSON, text, file).
// Accepts a variadic list of log types and sets up the corresponding loggers.
func SetLogTypes(types ...logType) {
	slogLoggers = []*slog.Logger{}
	for _, t := range types {
		switch t {
		case LogTypeJSON:
			slogLoggers = append(slogLoggers, newJSONLogger(os.Stderr))
		case LogTypeText:
			slogLoggers = append(slogLoggers, newTextLogger(os.Stderr))
		case LogTypeFile:
			if fileLogger == nil {
				_ = SetLogFile(DefaultLogFile)
			}
			slogLoggers = append(slogLoggers, newFileLogger(fileLogger))
		default:
			fmt.Printf("Unknown log type: %s\n", t)
		}
	}
}

// SetLogFile sets the log file path and configures the file logger.
// This function should be called only for the Master level to log errors to a file.
func SetLogFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("invalid file path")
	}

	// Ensure the directory for the log file exists.
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Attempt to create or open the specified file with appropriate flags and permissions.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to set log file: %w", err)
	}

	// Create the file-based logger and update loggers.
	fileLogger = file

	return nil
}

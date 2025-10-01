package errs

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// Define custom logging levels.
type LogType string

const (
	LogTypeJSON    LogType = "JSON"            // JSON logging format.
	LogTypeText    LogType = "TEXT"            // Text logging format.
	LogTypeFile    LogType = "FILE"            // File logging format.
	DefaultLogFile         = "log/logger.json" // Default log file path.
)

// Global variable to hold the bot instance
var bot *BroadcastBot

// Variables to manage the loggers and logging levels.
var (
	separator  string
	fileLogger *os.File
	jsonLogger *slog.Logger
	jsonBuf    *bytes.Buffer

	slogLoggers []*slog.Logger // List of loggers.
)

// init initializes loggers when the package is first loaded.
// It sets up the loggers using the LogTypeJSON constant by default.
func init() {
	SetLogTypes(LogTypeJSON)
	SetSupervisorErr(" ---> ")
}

// SetLogTypes configures the logging types to be used (e.g., JSON, text, file).
// Accepts a variadic list of log types and sets up the corresponding loggers.
//
// The function iterates through the provided log types and creates a new logger for each type.
// It supports three types: LogTypeJSON, LogTypeText, and LogTypeFile.
//
// For LogTypeJSON, it creates a JSON logger using the newJSONLogger function and appends it to the slogLoggers slice.
// For LogTypeText, it creates a text logger using the newTextLogger function and appends it to the slogLoggers slice.
// For LogTypeFile, it checks if the fileLogger is nil. If it is, it calls the SetLogFile function with the DefaultLogFile path.
// Then, it creates a file logger using the newFileLogger function and appends it to the slogLoggers slice.
//
// If an unknown log type is encountered, it prints an error message to the console.
//
// Note: The function does not return any value.
func SetLogTypes(types ...LogType) {
	slogLoggers = []*slog.Logger{}
	for _, t := range types {
		switch t {
		case LogTypeJSON:
			jsonLogger = newJSONLogger(os.Stderr)
			slogLoggers = append(slogLoggers, jsonLogger)
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
//
// Parameters:
// filePath (string): The path to the log file. If the directory does not exist, it will be created.
//
// Returns:
// error: An error if the file path is invalid, if the log directory cannot be created, or if the log file cannot be opened.
// If no error occurs, it returns nil.
func SetLogFile(filePath string) error {
	if filePath == "" {
		return New("invalid file path")
	}

	// Ensure the directory for the log file exists.
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return Wrap(err, "failed to create log directory")
	}

	// Attempt to create or open the specified file with appropriate flags and permissions.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return Wrap(err, "failed to set log file")
	}

	// Create the file-based logger and update loggers.
	fileLogger = file

	return nil
}


// SetSupervisorErr sets the separator used for logging errors from the supervisor.
// This separator is appended to the error message before logging it.
//
// Parameters:
// sep (string): The separator to be used for logging errors.
//
// Return:
// None.
func SetSupervisorErr(sep string) {
	separator = sep
}

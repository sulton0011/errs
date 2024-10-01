package errs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Define custom logging levels.
type level string

const (
	LevelLocal     level = "LOCAL"   // Local development logging level.
	LevelStaging   level = "STAGING" // Staging environment logging level.
	LevelMaster    level = "MASTER"  // Production (master) environment logging level.
	DefaultLogFile       = "log/logger.json"
)

// Variables to manage the loggers and logging levels.
var (
	currentLevel   = LevelMaster // Default logging level.
	slogJSONLogger *slog.Logger  // JSON logger for structured logging.
	slogTextLogger *slog.Logger  // Text logger for human-readable logs.
	fileLogger     *slog.Logger  // File logger for writing logs to a file.
)

// Initializes loggers when the package is first loaded.
func init() {
	updateLoggers()
}

// UpdateLevel sets a new logging level and reinitializes the loggers.
func UpdateLevel(newLevel level) {
	currentLevel = newLevel
	updateLoggers()
} // SetLogFile sets the log file path and configures the file logger.
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
	fileLogger = newFileLogger(file)
	updateLoggers()
	return nil
}

// updateLoggers creates new JSON and Text loggers based on the current settings.
func updateLoggers() {
	// Create new JSON and text loggers for stdout output.
	stdoutLogger := os.Stderr // Default to standard error output.
	slogJSONLogger = newJSONLogger(stdoutLogger)
	slogTextLogger = newTextLogger(stdoutLogger)
}

// newJSONLogger creates a new JSON logger for structured logging with no source path.
func newJSONLogger(output io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
		AddSource: false, // Disable source file and line number information.
		Level:     slog.LevelError,
	}))
}

// newTextLogger creates a new text logger for human-readable logging with no source path.
func newTextLogger(output io.Writer) *slog.Logger {
	return slog.New(newPrettyHandler(output, slog.LevelError))
}

func newPrettyHandler(output io.Writer, lvl slog.Level) slog.Handler {
	return &prettyHandler{
		Handler: slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     lvl,
		}),
		l: log.New(output, "", log.LstdFlags),
	}
}

// newFileLogger creates a logger that writes logs to a specified file without source paths.
func newFileLogger(file *os.File) *slog.Logger {
	return slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: false, // Disable source file and line number information for file logging.
		Level:     slog.LevelError,
	}))
}

// prettyHandler - custom handler for pretty text-based log output.
type prettyHandler struct {
	slog.Handler
	l *log.Logger
}

// Handle implements slog.Handler and formats logs in a custom pretty text style.
func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelStr := r.Level.String() + ":"
	levelStr = color.RedString(levelStr)

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	// Format JSON strings.
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	// Remove escaping of ">" symbol.
	formattedString := strings.ReplaceAll(string(b), `\u003e`, `>`)
	message := color.CyanString(r.Message)

	h.l.Println(levelStr, message, color.WhiteString(formattedString))
	return nil
}

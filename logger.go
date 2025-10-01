package errs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
)

// newJSONLogger creates a new JSON logger for structured logging with no source path.
func newJSONLogger(output io.Writer) *slog.Logger {
	jsonBuf = &bytes.Buffer{}
	mw := io.MultiWriter(output, jsonBuf)

	return slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
		AddSource: false, // Disable source file and line number information.
		Level:     slog.LevelError,
	}))
}

// newTextLogger creates a new text logger for human-readable logging with no source path.
func newTextLogger(output io.Writer) *slog.Logger {
	return slog.New(newPrettyHandler(output, slog.LevelError))
}

// newPrettyHandler creates a custom pretty handler for formatted logging.
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
	formattedJSON, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	formattedMessage := strings.ReplaceAll(string(formattedJSON), `\u003e`, `>`)
	h.l.Println(levelStr, color.CyanString(r.Message), color.WhiteString(formattedMessage))
	return nil
}

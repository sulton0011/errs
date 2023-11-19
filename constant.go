package errs

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type Level string

const (

	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

var (
	levelLogger Level = LevelError

	stdoutLogger io.Writer
)

func UpdateLevel(level Level) {
	levelLogger = level
}

func UpdateStdout(stdout io.Writer) {
	stdoutLogger = stdout
}

var (
	slogs = slog.New(slog.NewJSONHandler(getStdoutLogger(stdoutLogger), &slog.HandlerOptions{
		AddSource: true,
		Level: slogLevel(levelLogger),
	}))

	logs = NewSlog(getStdoutLogger(stdoutLogger), PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel(levelLogger),
		},
	})
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func slogLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	}

	return slog.LevelInfo
}

func getStdoutLogger(stdout io.Writer) io.Writer {
	if stdout != nil {
		return stdout
	}

	return os.Stdout
}

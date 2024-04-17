package errs

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type level string

const (
	LevelDebug level = "DEBUG"
	LevelInfo  level = "INFO"
	LevelWarn  level = "WARN"
	LevelError level = "ERROR"
)

var (
	levelLogger = LevelError

	stdoutLogger io.Writer

	slogs *slog.Logger
	logs  *slog.Logger
)

func UpdateLevel(level level) {
	levelLogger = level
}

func UpdateStdout(stdout *io.Writer) {
	stdoutLogger = *stdout

	slogs = slog.New(slog.NewJSONHandler(*stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slogLevel(levelLogger),
	}))

	logs = newSlog(*stdout, prettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel(levelLogger),
		},
	})
}

func sLoggerDefault() *slog.Logger {
	if slogs != nil {
		return slogs
	}
	slogs = slog.New(slog.NewJSONHandler(getStdoutLogger(stdoutLogger), &slog.HandlerOptions{
		// AddSource: true,
		Level: slogLevel(levelLogger),
	}))
	return slogs
}

func loggerDefault() *slog.Logger {
	if logs != nil {
		return logs
	}
	logs = newSlog(getStdoutLogger(stdoutLogger), prettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			// AddSource: true,
			Level: slogLevel(levelLogger),
		},
	})
	return logs

}

type prettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type prettyHandler struct {
	slog.Handler
	l *log.Logger
}

func slogLevel(level level) slog.Level {
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

	return os.Stderr
}

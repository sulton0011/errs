package errs

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"log/slog"

	"github.com/fatih/color"
)

// Уровни логирования
type level string

const (
	LevelLocal   level = "LOCAL"   // Для локальной разработки (все уровни логирования).
	LevelStaging level = "STAGING" // Для тестирования (ошибки и предупреждения).
	LevelMaster  level = "MASTER"  // Для продакшена (только ошибки).
)

// Глобальные переменные для управления логгерами и потоками вывода.
var (
	currentLevel   = LevelMaster           // Текущий уровень логирования.
	stdoutLogger   io.Writer = os.Stderr   // Логгер по умолчанию.
	slogJSONLogger *slog.Logger            // JSON логгер.
	slogTextLogger *slog.Logger            // Текстовый логгер.
)

// Инициализация логгеров при запуске.
func init() {
	updateLoggers()
}

// UpdateLevel обновляет текущий уровень логирования и пересоздает логгеры.
func UpdateLevel(newLevel level) {
	currentLevel = newLevel
	updateLoggers()
}

// UpdateStdout обновляет поток вывода логов и пересоздает логгеры.
func UpdateStdout(stdout io.Writer) {
	if stdout != nil {
		stdoutLogger = stdout
	}
	updateLoggers()
}

// updateLoggers создает новые JSON и текстовые логгеры на основе текущих настроек.
func updateLoggers() {
	slogJSONLogger = createLogger(stdoutLogger, slog.LevelError, "json")
	slogTextLogger = createLogger(stdoutLogger, slog.LevelInfo, "text")
}

// createLogger создает новый логгер на основе переданного типа вывода (json или text).
func createLogger(output io.Writer, baseLevel slog.Level, format string) *slog.Logger {
	switch format {
	case "json":
		return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     determineLogLevel(),
		}))
	case "text":
		return slog.New(newPrettyHandler(output, determineLogLevel()))
	default:
		return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     baseLevel,
		}))
	}
}

// determineLogLevel определяет уровень логирования на основе текущего уровня.
func determineLogLevel() slog.Level {
	switch currentLevel {
	case LevelLocal:
		return slog.LevelDebug // Все уровни логирования (DEBUG и выше).
	case LevelStaging:
		return slog.LevelWarn // Только предупреждения и ошибки.
	case LevelMaster:
		return slog.LevelError // Только ошибки.
	default:
		return slog.LevelError
	}
}

// prettyHandler - структура для красивого текстового вывода логов.
type prettyHandler struct {
	slog.Handler
	l *log.Logger
}

// newPrettyHandler создает новый обработчик для красивого текстового вывода логов.
func newPrettyHandler(output io.Writer, lvl slog.Level) slog.Handler {
	return &prettyHandler{
		Handler: slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     lvl,
		}),
		l: log.New(output, "", log.LstdFlags),
	}
}

// Handle реализует интерфейс slog.Handler и форматирует логи в текстовом виде.
func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	// Определение цвета и форматирование строки уровня логирования.
	levelStr := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		levelStr = color.MagentaString(levelStr)
	case slog.LevelInfo:
		levelStr = color.BlueString(levelStr)
	case slog.LevelWarn:
		levelStr = color.YellowString(levelStr)
	case slog.LevelError:
		levelStr = color.RedString(levelStr)
	}

	// Сбор атрибутов из логов в словарь.
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	// Форматирование JSON строк.
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	// Удаление экранирования символов.
	formattedString := strings.ReplaceAll(string(b), `\u003e`, `>`)
	message := color.CyanString(r.Message)

	// Вывод логов с цветным форматированием.
	h.l.Println(levelStr, message, color.WhiteString(formattedString))
	return nil
}

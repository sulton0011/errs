package errs

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

func NewSlog(out io.Writer, opts PrettyHandlerOptions) *slog.Logger {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return slog.New(h)
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, ``, `  `)
	if err != nil {
		return err
	}

	s := strings.ReplaceAll(string(b), `\u003e`, `>`)

	timeStr := r.Time.Format("[02-01-2006 15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(s))

	return nil
}

package logger

import (
	"log/slog"
	"os"
)

func Init() {
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Key = "t"
			return a
		}
		if a.Key == slog.LevelKey {
			a.Key = "l"
			return a
		}
		if a.Key == slog.MessageKey {
			a.Key = "msg"
			return a
		}
		return a
	}
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource:   true,
			Level:       slog.LevelDebug,
			ReplaceAttr: replaceAttr,
		},
	)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

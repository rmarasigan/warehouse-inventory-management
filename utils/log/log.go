package log

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func Init() {
	var level = new(slog.LevelVar)

	config := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, attribute slog.Attr) slog.Attr {
			if attribute.Key == slog.TimeKey {
				return slog.String("timestamp", attribute.Value.Time().Format("2006-01-02 3:04:05 PM"))

			} else if attribute.Key == slog.SourceKey {
				var (
					source   = attribute.Value.Any().(*slog.Source)
					id       = strings.LastIndex(source.File, "/")
					line     = strconv.Itoa(source.Line)
					filename = source.File[id+1:]
				)

				return slog.String(slog.SourceKey, slog.StringValue(filename+":"+line).String())
			}

			return attribute
		},
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, config)))
	level.Set(slog.LevelInfo)
}

func attributes(level slog.Level, msg string, args ...any) {
	var logger = slog.With(
		slog.Group("keys", args...),
	)

	switch level {
	case slog.LevelInfo:
		logger.Info(msg)

	case slog.LevelWarn:
		logger.Warn(msg)

	case slog.LevelError:
		logger.Error(msg)
	}
}

func Info(msg string, args ...any) {
	attributes(slog.LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	attributes(slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	attributes(slog.LevelError, msg, args...)
}

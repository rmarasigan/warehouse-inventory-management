package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	LevelOK    = slog.Level(2)
	LevelPanic = slog.Level(16)
)

var LevelNames = map[slog.Leveler]string{
	LevelOK:    "OK",
	LevelPanic: "PANIC",
}

func Init() {
	var level = new(slog.LevelVar)

	config := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, attribute slog.Attr) slog.Attr {
			switch attribute.Key {
			case slog.TimeKey:
				return slog.String("timestamp", attribute.Value.Time().Format("2006-01-02 3:04:05 PM"))

			case slog.LevelKey:
				level := attribute.Value.Any().(slog.Level)

				label, exists := LevelNames[level]
				if !exists {
					label = level.String()
				}
				attribute.Value = slog.StringValue(label)

				return attribute

			case slog.MessageKey:
				attribute.Key = "message"
				return attribute
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
	case LevelOK:
		logger.Log(context.TODO(), LevelOK, msg)

	case slog.LevelInfo:
		logger.Info(msg)

	case slog.LevelWarn:
		logger.Warn(msg)

	case slog.LevelError:
		_, filename, line, ok := runtime.Caller(2)
		if !ok {
			filename = "unknown"
			line = 0
		}

		logger.With("source", fmt.Sprintf("%s:%d", filename, line)).Error(msg)
	}
}

func OK(msg string, args ...any) {
	attributes(LevelOK, msg, args...)
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

func Panic(args ...any) {
	if msg := recover(); msg != nil {
		_, filename, line, ok := runtime.Caller(2)
		if !ok {
			filename = "unknown"
			line = 0
		}

		var logger = slog.With(
			slog.Group("keys", args...),
			slog.String("source", fmt.Sprintf("%s:%d", filename, line)),
			slog.Any("trace", strings.Split(string(debug.Stack()), "\n")),
		)

		logger.Log(context.TODO(), LevelPanic, slog.AnyValue(msg).String())
	}
}

package trail

import (
	"fmt"
	"strconv"
	"time"
)

const (
	FGRed    = 31
	FGCyan   = 36
	FGBlue   = 34
	FGYellow = 33
	Reset    = "\033[0m"

	LevelOK    = "OK"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

// colorize adds color to the timestamp and level.
func colorize(code int, level string) string {
	var timestamp = time.Now().Format("2006-01-02 3:04:05 PM")
	return fmt.Sprintf("\033[\x1b[%s;1m%s | %s%s", strconv.Itoa(code), timestamp, level, Reset)
}

// print logs the message with the appropriate level and color.
func print(level, input string, args ...any) {
	// Format input if there are additional arguments.
	if len(args) > 0 {
		input = fmt.Sprintf(input, args...)
	}

	// Print the message with appropriate color and level.
	switch level {
	case LevelOK:
		fmt.Println(colorize(FGCyan, LevelOK) + "\t" + input)

	case LevelInfo:
		fmt.Println(colorize(FGBlue, LevelInfo) + "\t" + input)

	case LevelWarn:
		fmt.Println(colorize(FGYellow, LevelWarn) + "\t" + input)

	case LevelError:
		fmt.Println(colorize(FGRed, LevelError) + "\t" + input)
	}
}

func OK(msg string, args ...any) {
	print(LevelOK, msg, args...)
}

func Info(msg string, args ...any) {
	print(LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	print(LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	print(LevelError, msg, args...)
}

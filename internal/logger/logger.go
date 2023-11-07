package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	sl *slog.Logger
}

func New(level string) *Logger {
	var lvl slog.Level
	switch level {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "INFO":
		lvl = slog.LevelInfo
	case "WARN":
		lvl = slog.LevelWarn
	case "ERROR":
		lvl = slog.LevelError
	}

	sl := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	return &Logger{sl}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.sl.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.sl.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.sl.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.sl.Error(msg, args...)
}

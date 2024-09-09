package logger

import (
	"log/slog"
	"os"
)

var defaultLogLevel *slog.LevelVar
var defaultLogger *slog.Logger

func init() {
	SetLogLevel(slog.LevelDebug)
	Init(defaultLogLevel, true)
}

func Init(level slog.Leveler, IsAddSource bool) {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: IsAddSource,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

func SetLogLevel(level slog.Level) {
	if defaultLogLevel == nil {
		defaultLogLevel = &slog.LevelVar{}
	}
	defaultLogLevel.Set(level)
}

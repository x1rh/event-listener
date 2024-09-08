package logger

import (
	"log/slog"
	"os"
)

var defaultLogLevel *slog.LevelVar
var defaultLogger *slog.Logger

func init() {
	level := slog.LevelDebug
	SetLogLevel(&level)
	opts := &slog.HandlerOptions{
		Level:     defaultLogLevel,
		AddSource: true,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

func SetLogLevel(level *slog.Level) {
	if defaultLogLevel == nil {
		defaultLogLevel = &slog.LevelVar{}
	}
	defaultLogLevel.Set(*level)
}


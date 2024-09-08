package logger

import (
	"testing"
	"log/slog"
)

func TestLog(t *testing.T) {
	slog.Debug("debug message")
	slog.Info("info message")
	slog.Warn("warning message")
	slog.Error("error message")
}

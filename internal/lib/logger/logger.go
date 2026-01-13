package logger

import (
	"log/slog"
	"os"
)

// NewLogger консструктор создания логгера
func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}

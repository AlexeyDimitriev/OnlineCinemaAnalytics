package logger

import (
	"os"
	"log/slog"
)

func NewLogger() *slog.Logger {
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(handler)
}

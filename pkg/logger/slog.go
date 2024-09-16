package logger

import (
	"log"
	"os"

	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

type loggerWrapper struct {
	slog *slog.Logger
}

func newSlogLogger(level string) *slog.Logger {
	var sLogger *slog.Logger

	switch level {
	case envLocal:
		sLogger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		sLogger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log.Fatalf("invalid environment: %s", level)
	}

	return sLogger
}

package logger

import (
	"log"
	"log/slog"
	"os"
)

const envLocal = "local"

func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log.Fatalf("Undefinded env: %s", env)
	}

	return logger
}

package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func New(mode string) *slog.Logger {
	var logger *slog.Logger

	switch mode {
	case "local":
		logger = slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, TimeFormat: "15:04:05"}),
		)

	case "dev":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case "prod":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		logger = slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, TimeFormat: "15:04:05"}),
		)
	}
	slog.SetDefault(logger)
	return logger
}

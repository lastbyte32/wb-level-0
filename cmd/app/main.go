package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lmittmann/tint"

	"github.com/lastbyte32/wb-level-0/internal/application"
	"github.com/lastbyte32/wb-level-0/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer cancel()
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to read config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger := Logger(cfg.Mode)

	app := application.New(cfg, logger)
	if err := app.Run(ctx); err != nil {
		slog.Error("failed to run app", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func Logger(mode string) *slog.Logger {
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

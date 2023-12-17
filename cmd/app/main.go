package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbyte32/wb-level-0/internal/application"
	"github.com/lastbyte32/wb-level-0/internal/config"
	"github.com/lastbyte32/wb-level-0/pkg/logger"
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
	sLogger := logger.New(cfg.Mode)
	app := application.New(cfg, sLogger)
	if err := app.Run(ctx); err != nil {
		slog.Error("failed to run app", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

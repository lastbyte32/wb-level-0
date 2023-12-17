package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/lastbyte32/wb-level-0/internal/config"
	"github.com/lastbyte32/wb-level-0/internal/handler/amqp"
	"github.com/lastbyte32/wb-level-0/internal/handler/rest"
	"github.com/lastbyte32/wb-level-0/internal/storage"
	"github.com/lastbyte32/wb-level-0/internal/usecase"
	"github.com/lastbyte32/wb-level-0/pkg/client/pgsql"
	"github.com/lastbyte32/wb-level-0/pkg/client/stan"
)

type app struct {
	cfg    *config.App
	logger *slog.Logger
}

func New(c *config.App, l *slog.Logger) *app {
	return &app{
		cfg:    c,
		logger: l,
	}
}

func (a *app) Run(ctx context.Context) error {
	stanMgr := stan.New(a.cfg.Nats.URL, a.cfg.Nats.ClusterID)
	sc, err := stanMgr.NewConnection(a.cfg.Nats.ClientID)
	if err != nil {
		panic(err)
	}
	a.logger.Info("connected to NATS streaming")

	pg, err := pgsql.NewPostgres(a.cfg.DataBase.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.logger.Info("connected to database")
	cache := storage.NewMemory()
	db := storage.NewPgsql(pg)

	orderUseCase := usecase.New(db, cache)
	countOrderInCache, err := orderUseCase.WarmUpCache(ctx)
	if err != nil {
		return fmt.Errorf("failed to warm up cache: %w", err)

	}

	a.logger.Info("warming up cache", slog.Int("orders", countOrderInCache))

	amqpHandler := amqp.New(sc, orderUseCase, a.logger)
	if err := amqpHandler.Subscribe(a.cfg.Nats.Subject); err != nil {
		return fmt.Errorf("failed to subscribe to subject %s", err)
	}

	restHandler := rest.New(orderUseCase)

	a.logger.Info("subscribed to new orders")

	http.HandleFunc("/api/v1/order/", restHandler.OrderByID)

	httpServer := http.Server{
		Addr: ":8081",
	}
	go func() {
		a.logger.Info("shutting down watcher start")
		<-ctx.Done()
		a.logger.Warn("signal received, shutting down")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		if err := sc.Close(); err != nil {
			a.logger.Error("failed to close NATS streaming")
		}
		if err := pg.Close(); err != nil {
			a.logger.Error("failed to close database")
		}
	}()

	a.logger.Info("http server started")
	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}
	return nil
}

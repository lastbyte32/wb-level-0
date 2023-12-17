package amqp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/stan.go"

	"github.com/lastbyte32/wb-level-0/internal/model"
)

type creator interface {
	CreateOrder(ctx context.Context, order model.Order) error
}

type handler struct {
	service creator
	stan    stan.Conn
	logger  *slog.Logger
}

func New(stan stan.Conn, service creator, logger *slog.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger,
		stan:    stan,
	}
}

func (h *handler) Subscribe(subject string) error {
	_, err := h.stan.Subscribe(subject, func(msg *stan.Msg) {
		h.logger.Info("received message")
		var order model.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			h.logger.Error("failed to unmarshal json", slog.String("error", err.Error()))
			return
		}
		if err := h.service.CreateOrder(context.Background(), order); err != nil {
			h.logger.Error("failed to create order", slog.String("error", err.Error()))
			return
		}
		h.logger.Info("created order")
	})
	return err
}

package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/lastbyte32/wb-level-0/internal/model"
)

type storage interface {
	GetByID(ctx context.Context, id string) (model.Order, error)
	Store(ctx context.Context, id string, order model.Order) error
	All(ctx context.Context) (map[string]model.Order, error)
}

type orderUC struct {
	db  storage
	mem storage
}

func New(db, mem storage) *orderUC {
	return &orderUC{
		db:  db,
		mem: mem,
	}
}

func (u *orderUC) CreateOrder(ctx context.Context, order model.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("failed to generate id: %w", err)
	}
	if err := u.mem.Store(ctx, id.String(), order); err != nil {
		return fmt.Errorf("failed to store order on memory: %w", err)
	}
	if err := u.db.Store(ctx, id.String(), order); err != nil {
		return fmt.Errorf("failed to store order on db: %w", err)
	}
	return nil
}

func (u *orderUC) GetOrderByID(ctx context.Context, id string) (model.Order, error) {
	order, err := u.mem.GetByID(ctx, id)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}
	if order.OrderUID == "" {
		return model.Order{}, model.ErrNotFound
	}
	return order, nil
}

func (u *orderUC) GetOrderFromDB(ctx context.Context, id string) (model.Order, error) {
	order, err := u.db.GetByID(ctx, id)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}
	if order.OrderUID == "" {
		return model.Order{}, model.ErrNotFound
	}
	return order, nil
}

func (u *orderUC) GetAll(ctx context.Context) (map[string]model.Order, error) {
	orders, err := u.db.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (u *orderUC) WarmUpCache(ctx context.Context) (int, error) {
	orders, err := u.db.All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get orders: %w", err)
	}

	for id, order := range orders {
		if err := u.mem.Store(ctx, id, order); err != nil {
			fmt.Printf("failed to store order on memory: %v\n", err)
		}
	}

	return len(orders), nil
}

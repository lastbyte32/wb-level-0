package storage

import (
	"context"
	"sync"

	"github.com/lastbyte32/wb-level-0/internal/model"
)

type memory struct {
	cache map[string]model.Order
	mu    sync.RWMutex
}

func (c *memory) All(ctx context.Context) (map[string]model.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	orders := make(map[string]model.Order)
	for id, order := range c.cache {
		orders[id] = order
	}
	return orders, nil

}

func NewMemory() *memory {
	return &memory{
		cache: make(map[string]model.Order),
	}
}

func (c *memory) Store(ctx context.Context, id string, order model.Order) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.cache[id] = order
	return nil
}

func (c *memory) GetByID(ctx context.Context, id string) (model.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.cache[id]
	if !ok {
		return model.Order{}, model.ErrNotFound
	}
	return order, nil
}

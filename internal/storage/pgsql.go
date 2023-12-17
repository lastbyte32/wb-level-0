package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lastbyte32/wb-level-0/internal/model"
)

type pgsql struct {
	db *sql.DB
}

func NewPgsql(conn *sql.DB) *pgsql {
	return &pgsql{db: conn}
}

func (c *pgsql) Store(ctx context.Context, id string, order model.Order) error {
	query := "INSERT INTO orders (id, data) VALUES ($1, $2)"
	jsonData, err := order.JSON()
	if err != nil {
		return fmt.Errorf("failed to marshal order to JSON: %w", err)
	}
	_, err = c.db.ExecContext(ctx, query, id, jsonData)
	if err != nil {
		return fmt.Errorf("failed to store order: %w", err)
	}
	return nil
}

func (c *pgsql) GetByID(ctx context.Context, id string) (model.Order, error) {
	query := "SELECT data FROM orders WHERE id = $1"
	var jsonData string
	if err := c.db.QueryRowContext(ctx, query, id).Scan(&jsonData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, model.ErrNotFound
		}
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}
	var order model.Order
	if err := json.Unmarshal([]byte(jsonData), &order); err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return order, nil
}

func (c *pgsql) All(ctx context.Context) (map[string]model.Order, error) {
	query := "SELECT id, data FROM orders"
	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	defer rows.Close()

	orders := make(map[string]model.Order)
	for rows.Next() {
		var id, jsonData string
		if err := rows.Scan(&id, &jsonData); err != nil {
			return nil, fmt.Errorf("failed to scan order data: %w", err)
		}

		var order model.Order
		if err := json.Unmarshal([]byte(jsonData), &order); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %w", err)
		}

		orders[id] = order
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}

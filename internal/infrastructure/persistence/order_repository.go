// internal/infrastructure/persistence/order_repository.go
package persistence

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
	"fmt"
	"time"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, user_id, symbol, order_type, side, quantity, price, status, tachibana_order_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		order.ID, "user_id", order.Symbol, order.OrderType, order.Side, // user_id は仮の値
		order.Quantity, order.Price, order.Status, order.TachibanaOrderID, time.Now(), time.Now(),
	)
	return err
}

func (r *orderRepository) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	query := `
		SELECT id, user_id, symbol, order_type, side, quantity, price, status, tachibana_order_id, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	order := &domain.Order{}
	var userID string // user_id は使わないので仮の変数
	err := row.Scan(
		&order.ID, &userID, &order.Symbol, &order.OrderType, &order.Side,
		&order.Quantity, &order.Price, &order.Status, &order.TachibanaOrderID, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Order not found
		}
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders
		SET symbol = $2, order_type = $3, side = $4, quantity = $5, price = $6, status = $7, tachibana_order_id = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.Symbol, order.OrderType, order.Side,
		order.Quantity, order.Price, order.Status, order.TachibanaOrderID, time.Now(),
	)
	return err
}

// UpdateOrderStatus は注文のステータスを更新
func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	query := `
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `
	res, err := r.db.ExecContext(ctx, query, orderID, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order not found: %s", orderID)
	}
	return nil
}

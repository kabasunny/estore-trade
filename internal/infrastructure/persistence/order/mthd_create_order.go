// internal/infrastructure/persistence/order/mthd_create_order.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
        INSERT INTO orders (id, symbol, order_type, side, quantity, price, trigger_price, filled_quantity, average_price, status, tachibana_order_id, commission, expire_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `
	_, err := r.db.ExecContext(ctx, query,
		order.UUID,
		order.Symbol,
		order.OrderType,
		order.Side,
		order.Quantity,
		order.Price,
		order.TriggerPrice,
		order.FilledQuantity,
		order.AveragePrice,
		order.Status,
		order.TachibanaOrderID,
		order.Commission,
		order.ExpireAt,
		time.Now(),
		time.Now(),
	)
	return err
}

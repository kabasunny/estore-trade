// internal/infrastructure/persistence/order/mthd_update_order.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

func (r *orderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	query := `
        UPDATE orders
        SET symbol = $2, order_type = $3, side = $4, quantity = $5, price = $6, trigger_price = $7, filled_quantity = $8, average_price = $9, status = $10, tachibana_order_id = $11, commission = $12, expire_at = $13, updated_at = $14
        WHERE id = $1
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
	)
	return err
}

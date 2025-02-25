// internal/infrastructure/persistence/order/mthd_get_order.go
package order

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
)

func (r *orderRepository) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	query := `
        SELECT id, symbol, order_type, side, quantity, price, trigger_price, filled_quantity, average_price, status, tachibana_order_id, commission, expire_at, created_at, updated_at
        FROM orders
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)

	order := &domain.Order{}
	err := row.Scan(
		&order.ID,
		&order.Symbol,
		&order.OrderType,
		&order.Side,
		&order.Quantity,
		&order.Price,
		&order.TriggerPrice,
		&order.FilledQuantity,
		&order.AveragePrice,
		&order.Status,
		&order.TachibanaOrderID,
		&order.Commission,
		&order.ExpireAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return order, nil
}

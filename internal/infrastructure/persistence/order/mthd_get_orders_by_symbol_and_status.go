// internal/infrastructure/persistence/order/mthd_get_orders_by_symbol_and_status.go
package order

import (
	"context"
	"estore-trade/internal/domain"
)

func (r *orderRepository) GetOrdersBySymbolAndStatus(ctx context.Context, symbol, status string) ([]*domain.Order, error) {
	query := `
        SELECT id, symbol, order_type, side, quantity, price, trigger_price, filled_quantity, average_price, status, tachibana_order_id, commission, expire_at, created_at, updated_at
        FROM orders
        WHERE symbol = $1 AND status = $2
    `
	rows, err := r.db.QueryContext(ctx, query, symbol, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order := &domain.Order{}
		err := rows.Scan(
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
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

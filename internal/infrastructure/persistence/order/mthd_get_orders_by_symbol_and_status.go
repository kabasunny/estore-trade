// internal/infrastructure/persistence/order/mthd_get_orders_by_symbol_and_status.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

func (r *orderRepository) GetOrdersBySymbolAndStatus(ctx context.Context, symbol, status string) ([]*domain.Order, error) {
	var orders []*domain.Order
	result := r.db.WithContext(ctx).Where("symbol = ? AND status = ?", symbol, status).Find(&orders) // Find を使用
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get orders by symbol and status: %w", result.Error)
	}
	return orders, nil
}

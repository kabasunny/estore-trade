// internal/infrastructure/persistence/order/mthd_update_order.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

func (r *orderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	result := r.db.WithContext(ctx).Save(order) // Save を使用 (既存レコードの更新)
	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}
	return nil
}

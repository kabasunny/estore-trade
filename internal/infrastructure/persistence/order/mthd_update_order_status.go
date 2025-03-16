// internal/infrastructure/persistence/order/mthd_update_order_status.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	result := r.db.WithContext(ctx).Model(&domain.Order{}).Where("uuid = ?", orderID).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found: %s", orderID)
	}
	return nil
}

// internal/infrastructure/persistence/order/mthd_create_order.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	result := r.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		return fmt.Errorf("failed to create order: %w", result.Error)
	}
	return nil
}

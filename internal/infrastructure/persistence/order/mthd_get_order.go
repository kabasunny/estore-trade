// internal/infrastructure/persistence/order/mthd_get_order.go
package order

import (
	"context"
	"errors"
	"estore-trade/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

func (r *orderRepository) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	order := &domain.Order{}
	result := r.db.WithContext(ctx).First(order, "uuid = ?", id) // First を使用
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 見つからない場合は nil, nil を返す
		}
		return nil, fmt.Errorf("failed to get order: %w", result.Error)
	}
	return order, nil
}

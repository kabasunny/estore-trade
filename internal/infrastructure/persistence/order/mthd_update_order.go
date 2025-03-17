// mthd_update_order.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

// mthd_update_order.go
func (r *orderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	result := r.db.WithContext(ctx).Model(&domain.Order{}).Where("uuid = ?", order.UUID).Updates(map[string]interface{}{
		"average_price":      order.AveragePrice,
		"commission":         order.Commission,
		"filled_quantity":    order.FilledQuantity,
		"side":               order.Side,
		"status":             order.Status,
		"symbol":             order.Symbol,
		"tachibana_order_id": order.TachibanaOrderID,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}
	if result.RowsAffected == 0 { // 追加: 更新された行がない場合
		return gorm.ErrRecordNotFound // 追加: gorm.ErrRecordNotFound を返す
	}
	return nil
}

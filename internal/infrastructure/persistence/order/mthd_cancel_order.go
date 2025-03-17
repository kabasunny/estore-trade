// mthd_cancel_order.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

func (r *orderRepository) CancelOrder(ctx context.Context, orderID string) error {
	// GORM を使って、注文のステータスを "canceled" に更新
	tx := r.db.WithContext(ctx).Model(&domain.Order{}).Where("uuid = ?", orderID).Update("status", "canceled")
	if tx.Error != nil {
		return fmt.Errorf("failed to cancel order: %w", tx.Error)
	}
	// 更新された行がない場合は、エラー
	if tx.RowsAffected == 0 {
		tx.Rollback() //追加
		return fmt.Errorf("order not found: %s", orderID)
	}

	return nil
}

// mthd_update_order_status.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	//Model(&domain.Order{})により、ordersテーブルが対象となる
	//Where("uuid = ?", orderID)により、uuidがorderIDと一致するレコードが対象となる
	//Update("status", status)により、statusカラムが指定されたstatusの値に更新される。
	result := r.db.WithContext(ctx).Model(&domain.Order{}).Where("uuid = ?", orderID).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 更新対象がない場合はエラー
	}
	return nil
}

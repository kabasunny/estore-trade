// internal/usecase/mthd_cancel_order.go
package usecase

import (
	"context"
)

// 指定された注文IDの注文をキャンセルする
func (uc *tradingUsecase) CancelOrder(ctx context.Context, orderID string) error {
	// APIクライアントを呼び出して注文をキャンセル
	err := uc.tachibanaClient.CancelOrder(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}

// internal/usecase/mthd_cancel_order.go
package usecase

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// CancelOrder は指定された注文IDの注文をキャンセルします。
func (uc *tradingUsecase) CancelOrder(ctx context.Context, orderID string) error {
	// APIクライアントを呼び出して注文をキャンセル
	err := uc.tachibanaClient.CancelOrder(ctx, orderID)
	if err != nil {
		uc.logger.Error("Failed to cancel order via Tachibana API", zap.String("orderID", orderID), zap.Error(err))
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// TachibanaClient 側でキャンセルのログが出力されるため、ここではログ出力しない

	return nil
}

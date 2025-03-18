// internal/usecase/mthd_get_order_status.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

// GetOrderStatus は指定された注文IDの注文状況をAPIから取得します。
func (uc *tradingUsecase) GetOrderStatus(ctx context.Context, orderID string, orderDate string) (*domain.Order, error) {
	// APIクライアントを呼び出して注文状況を取得
	order, err := uc.tachibanaClient.GetOrderStatus(ctx, orderID, orderDate)
	if err != nil {
		uc.logger.Error("Failed to get order status from Tachibana API", zap.String("orderID", orderID), zap.Error(err))
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}

	// ここでは、TachibanaClientから取得したOrder情報をそのまま返す
	return order, nil
}

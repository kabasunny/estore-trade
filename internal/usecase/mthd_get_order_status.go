// internal/usecase/mthd_get_order_status.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
)

// 指定された注文IDの注文状況を取得する
func (uc *tradingUsecase) GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error) {
	// APIクライアントを呼び出して注文状況を取得
	orderStatus, err := uc.tachibanaClient.GetOrderStatus(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return orderStatus, nil
}

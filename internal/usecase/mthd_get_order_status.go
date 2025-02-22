// internal/usecase/mthd_get_order_status.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
)

func (uc *tradingUsecase) GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error) {
	orderStatus, err := uc.tachibanaClient.GetOrderStatus(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return orderStatus, nil
}

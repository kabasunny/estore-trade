// internal/usecase/mthd_cancel_order.go
package usecase

import (
	"context"
)

func (uc *tradingUsecase) CancelOrder(ctx context.Context, orderID string) error {
	err := uc.tachibanaClient.CancelOrder(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}

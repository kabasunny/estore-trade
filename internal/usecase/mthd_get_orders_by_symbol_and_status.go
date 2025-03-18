// internal/usecase/mthd_get_orders_by_symbol_and_status.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

// GetOrdersBySymbolAndStatus は指定された銘柄コードとステータスの注文リストをDBから取得します。
func (uc *tradingUsecase) GetOrdersBySymbolAndStatus(ctx context.Context, symbol string, status string) ([]*domain.Order, error) {
	// OrderRepository を使用して、DBから注文情報を取得
	orders, err := uc.orderRepo.GetOrdersBySymbolAndStatus(ctx, symbol, status)
	if err != nil {
		uc.logger.Error("Failed to get orders from DB", zap.String("symbol", symbol), zap.String("status", status), zap.Error(err))
		return nil, fmt.Errorf("failed to get orders from DB: %w", err)
	}

	return orders, nil
}

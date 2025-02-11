package usecase

import (
	"context"
	"estore-trade/internal/domain"
)

// TradingUsecase defines the interface for trading operations.
type TradingUsecase interface {
	PlaceOrder(ctx context.Context, userID, password string, order *domain.Order) (*domain.Order, error)
	GetOrderStatus(ctx context.Context, userID, password string, orderID string) (*domain.Order, error)
	CancelOrder(ctx context.Context, userID, password string, orderID string) error
	// 他のユースケース
	// (例: GetAccountInfo, RunStrategy, GetMarketData,  etc.)
}

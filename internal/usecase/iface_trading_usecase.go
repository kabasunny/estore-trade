// internal/usecase/iface_trading_usecase.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
)

// 取引操作のためのインターフェース
type TradingUsecase interface {
	// スイングトレードの注文は基本的に逆指値注文となる
	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)

	GetOrderStatus(ctx context.Context, orderID string, orderDate string) (*domain.Order, error)

	CancelOrder(ctx context.Context, orderID string) error

	GetOrdersBySymbolAndStatus(ctx context.Context, symbol string, status string) ([]*domain.Order, error)
	// 他のユースケース (今のところ不要)
}

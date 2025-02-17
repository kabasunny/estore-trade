// internal/usecase/trading.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
)

// TradingUsecase defines the interface for trading operations.
type TradingUsecase interface {
	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) // 引数から userID, password を削除
	GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error)  // 引数から userID, password を削除
	CancelOrder(ctx context.Context, orderID string) error                      // 引数から userID, password を削除
	// EventStreamからイベントを受け取るためのチャネル (読み取り専用)
	GetEventChannelReader() <-chan domain.OrderEvent //修正
	// EventStreamにイベントを送信するためのチャネル(書き込み専用)
	GetEventChannelWriter() chan<- domain.OrderEvent //修正
	HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error
	// 他のユースケース
	// (例: GetAccountInfo, RunStrategy, GetMarketData,  etc.)
}

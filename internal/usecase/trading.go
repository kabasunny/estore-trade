// internal/usecase/trading.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
)

// 取引操作のためのインターフェース
type TradingUsecase interface {
	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) // 引数から userID, password を削除
	GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error)  // 引数から userID, password を削除
	CancelOrder(ctx context.Context, orderID string) error                      // 引数から userID, password を削除
	// EventStreamからイベントを受け取るためのチャネル (読み取り専用)
	GetEventChannelReader() <-chan domain.OrderEvent
	// EventStreamにイベントを送信するためのチャネル(書き込み専用)
	GetEventChannelWriter() chan<- domain.OrderEvent
	HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error
	// 他のユースケース
	// (例: GetAccountInfo, RunStrategy, GetMarketData,  etc.)
}

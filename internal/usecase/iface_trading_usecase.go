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
	// 上流層より注文依頼、クライアントに注文依依頼後、クライアントに注文状態確認を依頼し、注文状態格納をリポジトリに依頼する
	// 成行きの場合は、クライアント確認で約定している場合もある　この場合、上流からの依頼が重複情報となる
	// 逆指値の場合は、約定しない場合もある

	//OrderEventを受け取り、DBの注文情報を更新
	UpdateOrderByEvent(ctx context.Context, event *domain.OrderEvent) error
	// 上流層がイベントストリームで約定確認、上流層からの状態更新依頼を受ける

	// GetOrderStatus は注文状況を取得
	GetOrderStatus(ctx context.Context, orderID string, orderDate string) (*domain.Order, error)
	// クライアント確認後、リポジトリの齟齬の修正も行う

	// CancelOrder は指定された注文IDの注文をキャンセル
	CancelOrder(ctx context.Context, orderID string) error
	// 上流層より注文依頼、クライアントに注文依依頼後、クライアントに注文状態確認を依頼し、注文状態格納をリポジトリに依頼する
	// 現状ほとんど使用しないと想定

	// GetOrdersBySymbolAndStatus は指定された銘柄コードとステータスの注文リストを取得
	GetOrdersBySymbolAndStatus(ctx context.Context, symbol string, status string) ([]*domain.Order, error)
	// 同一銘柄に資金を集中すリスクを回避するため、特定の銘柄に対する現在のポジションを把握する
}

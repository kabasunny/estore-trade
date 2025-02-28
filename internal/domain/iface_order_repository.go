// internal/domain/iface_order_repository.go
package domain

import (
	"context"
)

// OrderRepository は注文（Order）データの永続化操作を抽象化
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error                                     // 新しい注文を作成
	GetOrder(ctx context.Context, id string) (*Order, error)                                 // 指定されたIDの注文を取得
	GetOrdersBySymbolAndStatus(ctx context.Context, symbol, status string) ([]*Order, error) // 銘柄とステータスで検索
	UpdateOrder(ctx context.Context, order *Order) error                                     // 注文のデータ更新 (ステータスを含む全てのフィールド)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error              // 注文ステータス更新 (ステータスフィールドのみ更新)
	CancelOrder(ctx context.Context, orderID string) error                                   // 注文のキャンセル(DBから削除ではなく、StatusをCancelledにする)
}

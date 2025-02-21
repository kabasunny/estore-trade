package domain

import (
	"context"
)

// 注文（Order）データの永続化操作を抽象化
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error                        // 新しい注文を作成する
	GetOrder(ctx context.Context, id string) (*Order, error)                    // 指定されたIDの注文を取得する
	UpdateOrder(ctx context.Context, order *Order) error                        // 指定された注文のデータを更新する
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error // 指定された注文IDの注文のステータスを更新する
}

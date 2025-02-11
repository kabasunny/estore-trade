package domain

import "context"

// データの永続化層における操作を抽象化し、ビジネスロジックとデータベース操作の結合度を低減するために使用

// 注文（Order）データの永続化操作を抽象化
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, id string) (*Order, error)
	UpdateOrder(ctx context.Context, order *Order) error
	// 他の必要なメソッド (例: CancelOrder, GetOrdersByStatus, etc.)
}

// 取引アカウント（Account）データの永続化操作を抽象化
type AccountRepository interface {
	GetAccount(ctx context.Context, id string) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
}

// 他のリポジトリインターフェースを必要に応じて定義

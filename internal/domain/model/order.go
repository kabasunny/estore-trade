// internal/domain/model/order.go
package model

import (
	"context"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model                   // ID, CreatedAt, UpdatedAt, DeletedAt を追加
	TachibanaOrderNumber string  `gorm:"type:varchar(255)"`
	IssueCode            string  `gorm:"type:varchar(10);not null"`
	AccountID            uint    `gorm:"not null"`
	Side                 string  `gorm:"type:varchar(10);not null"` // e.g., 'long', 'short'
	OrderType            string  `gorm:"type:varchar(50);not null"` // e.g., 'market', 'limit'
	Quantity             int     `gorm:"not null"`
	Price                float64 `gorm:"null"`                      // NULL for market orders
	Status               string  `gorm:"type:varchar(50);not null"` // e.g., 'pending', 'filled', 'cancelled'

	TachibanaOrderID string // 立花証券側の注文ID
	Symbol           string // 銘柄コード
	// 他の必要なフィールドを追加 (数量、価格など)
}

// TableName overrides the table name used by User to `profiles`
func (Order) TableName() string {
	return "orders"
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, id int) (*Order, error)
	UpdateOrder(ctx context.Context, order *Order) error
	UpdateOrderStatus(ctx context.Context, id int, status string) error
	CancelOrder(ctx context.Context, id int) error
	GetOrdersBySymbolAndStatus(ctx context.Context, symbol string, status string) ([]Order, error)
}

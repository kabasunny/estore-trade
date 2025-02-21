// internal/domain/model.go
package domain

import "time"

// 自動売買システムの中核となるデータ構造

// 株式の注文
type Order struct {
	ID               string
	Symbol           string
	Side             string // "buy" or "sell"
	OrderType        string // "market", "limit", etc.
	Price            float64
	Quantity         int
	Status           string // "pending", "filled", "canceled", etc.
	TachibanaOrderID string // 立花証券側の注文ID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// internal/domain/strct_order.go
package domain

import (
	"time"
)

// Order 株式の注文
type Order struct {
	ID               string    // 注文ID (UUID)
	Symbol           string    // 銘柄コード
	Side             string    // 売買区分 ("buy" or "sell")
	OrderType        string    // 注文種別 ("market", "limit", "stop", "stop_limit")
	Price            float64   // 注文価格 (指値、逆指値の場合)
	TriggerPrice     float64   // トリガー価格 (逆指値の場合)
	Quantity         int       // 注文数量
	FilledQuantity   int       // 約定数量
	AveragePrice     float64   // 平均約定価格
	Status           string    // 注文ステータス ("pending", "filled", "partially_filled", "canceled", "rejected")
	TachibanaOrderID string    // 立花証券側の注文ID
	Commission       float64   // 手数料
	ExpireAt         time.Time // 注文有効期限
	CreatedAt        time.Time // 注文作成日時
	UpdatedAt        time.Time // 注文最終更新日時
	// ParentOrderID string // 親注文ID (今回は使用しない)
}

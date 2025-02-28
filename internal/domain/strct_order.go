// internal/domain/order.go
package domain

import (
	"time"
)

// Order 株式の注文
type Order struct {
	ID               string    `json:"id"`                 // 注文ID (UUID)
	Symbol           string    `json:"symbol"`             // 銘柄コード
	Side             string    `json:"side"`               // 売買区分 ("buy" or "sell")
	OrderType        string    `json:"order_type"`         // 注文種別 ("market", "limit", "stop", "stop_limit")
	Price            float64   `json:"price"`              // 注文価格 (指値、逆指値の場合)
	TriggerPrice     float64   `json:"trigger_price"`      // トリガー価格 (逆指値の場合)
	Quantity         int       `json:"quantity"`           // 注文数量
	FilledQuantity   int       `json:"filled_quantity"`    // 約定数量
	AveragePrice     float64   `json:"average_price"`      // 平均約定価格
	Status           string    `json:"status"`             // 注文ステータス
	TachibanaOrderID string    `json:"tachibana_order_id"` // 立花証券側の注文ID
	Commission       float64   `json:"commission"`         // 手数料
	ExpireAt         time.Time `json:"expire_at"`          // 注文有効期限
	CreatedAt        time.Time `json:"created_at"`         // 注文作成日時
	UpdatedAt        time.Time `json:"updated_at"`         // 注文最終更新日時

	Condition string `json:"condition"` // 執行条件
}

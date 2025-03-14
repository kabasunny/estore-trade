package domain

import (
	"time"
)

// Order 株式の注文（既存の構造体を拡張）
type Order struct {
	UUID             string    `json:"id"`                 // 注文ID (UUID)
	Symbol           string    `json:"symbol"`             // 銘柄コード
	Side             string    `json:"side"`               // 売買区分 ("long" or "short")
	OrderType        string    `json:"order_type"`         // 注文種別
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

	Condition  string     `json:"condition"`   // 執行条件（"credit_open"（信用新規）など）
	MarketCode string     `json:"market_code"` // 市場コード
	Positions  []Position `json:"positions"`   // 信用返済時の建玉情報（追加）

	AfterTriggerOrderType string  // "market" or "limit"
	AfterTriggerPrice     float64 // トリガー後の指値 (トリガー後指値の場合のみ)

}

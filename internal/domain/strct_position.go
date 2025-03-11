// internal/domain/strct_position.go
package domain

import (
	"time"
)

// Position ポジション
type Position struct {
	ID              string    `json:"id"` // 建玉ID (UUID) <- これを追加
	Symbol          string    `json:"symbol"`
	Side            string    `json:"side"` //"long" or "short"
	Quantity        int       `json:"quantity"`
	Price           float64   `json:"price"`
	OpenDate        time.Time `json:"open_date"`
	DueDate         string    `json:"due_date"`          // 返済期限 (YYYYMMDD)
	MarginTradeType string    `json:"margin_trade_type"` // 信用取引区分
}

// internal/domain/strct_position.go
package domain

import (
	"time"
)

// Position ポジション
type Position struct {
	Symbol   string    `json:"symbol"`    // 銘柄コード
	Side     string    `json:"side"`      //"long" or "short"
	Quantity int       `json:"quantity"`  // 保有数量
	Price    float64   `json:"price"`     // 平均取得単価
	OpenDate time.Time `json:"open_date"` // ポジションを建てた日付
	//LastPrice          float64   // 現在の価格 (時価, 今回は使用しない)
	//UnrealizedProfitLoss float64 // 評価損益 (今回は使用しない)
}

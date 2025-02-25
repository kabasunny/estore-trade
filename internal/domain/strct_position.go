// internal/domain/strct_position.go
package domain

import (
	"time"
)

// Position ポジション
type Position struct {
	Symbol               string    // 銘柄コード
	Side                 string    // "long" (買い) or "short" (売り)
	Quantity             int       // 保有数量
	Price                float64   // 平均取得単価
	OpenDate             time.Time // ポジションを建てた日付
	LastPrice            float64   // 現在の価格 (時価, 今回は使用しない)
	UnrealizedProfitLoss float64   // 評価損益 (今回は使用しない)
}

// internal/domain/strct_position.go
package domain

import (
	"time"
)

// Position ポジション
type Position struct {
	ID              string
	Symbol          string
	Side            string // string のまま
	Quantity        int
	Price           float64
	OpenDate        time.Time
	DueDate         string // YYYYMMDD
	MarginTradeType string // 制度信用、一般信用など
	OrderID         string // 注文番号  //★ここ
}

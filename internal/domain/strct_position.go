// internal/domain/strct_position.go
package domain

import (
	"time"
)

// Position ポジション
type Position struct {
	ID              string // 立花証券APIにおける sTategyokuNumber (新規建玉番号) に対応
	Symbol          string
	Side            string // string のまま
	Quantity        int
	Price           float64
	OpenDate        time.Time
	DueDate         string // YYYYMMDD
	MarginTradeType string // 制度信用、一般信用など
	OrderID         string // 注文番号  //★ここ
}

// internal/domain/order_event.go
package domain

import "time"

type OrderEvent struct {
	EventType string    // "EC", "NS", "SS", "US" など
	EventNo   int       // p_ENO
	Order     *Order    // 更新された注文情報 (ECの場合)
	Timestamp time.Time // p_date
	// ... 他の必要なフィールド (エラー情報など)
	// 例:
	ErrNo  int    // p_errno (エラーの場合)
	ErrMsg string // p_err (エラーの場合)
}

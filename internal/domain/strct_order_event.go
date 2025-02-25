// internal/domain/strct_order_event.go
package domain

import "time"

// OrderEvent 注文イベント
type OrderEvent struct {
	EventType string    // イベントタイプ ("EC", "NS", "SS", "US" など)
	EventNo   int       // イベント番号 (p_ENO)
	Order     *Order    // 更新された注文情報 (ECの場合)
	Timestamp time.Time // イベント発生時刻 (p_date)
	ErrNo     int       // エラー番号 (p_errno, 0以外の場合)
	ErrMsg    string    // エラーメッセージ (p_err)
	// EventData map[string]interface{} // イベントに関する追加情報 (今回は使用しない)
}

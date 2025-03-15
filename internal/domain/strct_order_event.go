package domain

import "time"

// OrderEvent は、立花証券からのイベント通知を表すドメインモデル
type OrderEvent struct {
	EventType        string        // イベントタイプ ("EC", "NS", "SS", "US", "ST", "KP", "FD" など)
	EventNo          string        // イベント番号 (p_ENO)
	Order            *Order        // 更新された注文情報 (ECの場合)
	News             *News         // ニュース情報 (NSの場合)
	System           *SystemStatus // システムステータス (SS, US, ST, KPの場合) // 既存フィールドを優先
	Market           *Market       // 市場情報 (FDの場合)
	Timestamp        time.Time     // イベント発生時刻 (p_date)
	IsFirstEvent     bool          // アラートフラグ (p_ALT) 1:初回通知, 0:再送通知  // 追加
	Provider         string        // プロバイダ (p_PV)  // 追加
	NotificationType string        // 通知種別(p_NT) //Orderにあるが、OrderEventにも追加
}

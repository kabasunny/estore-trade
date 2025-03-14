package domain

import (
	"time"
)

// Order 株式の注文
type Order struct {
	UUID                  string     `json:"id"`                       // 注文ID (UUID) ... アプリケーション内でユニークなID (APIには直接関係ない)
	Symbol                string     `json:"symbol"`                   // 銘柄コード (例: "7974") ... sIssueCode
	Side                  string     `json:"side"`                     // 売買区分 ("long" or "short") ... sBaibaiKubun (1: 売, 3: 買)
	OrderType             string     `json:"order_type"`               // 注文種別 ("market", "limit", "stop", "stop_limit", "credit_close_market", "credit_close_limit", "credit_close_stop", "credit_close_stop_limit") ... sCondition, sGyakusasiOrderType
	Price                 float64    `json:"price"`                    // 注文価格 (指値、逆指値の場合) ... sOrderPrice (成行の場合は "0", 逆指値の場合は "*" だが、ここではトリガー後の価格)
	TriggerPrice          float64    `json:"trigger_price"`            // トリガー価格 (逆指値の場合) ... sGyakusasiZyouken
	Quantity              int        `json:"quantity"`                 // 注文数量 ... sOrderSuryou
	FilledQuantity        int        `json:"filled_quantity"`          // 約定数量 (APIからのレスポンス、イベントで更新) ... p_CREXSR, p_EXSR
	AveragePrice          float64    `json:"average_price"`            // 平均約定価格 (APIからのレスポンスで更新) ... p_EXPR
	Status                string     `json:"status"`                   // 注文ステータス (APIからのレスポンス、イベントで更新) ... p_ODST
	TachibanaOrderID      string     `json:"tachibana_order_id"`       // 立花証券側の注文ID (APIからのレスポンス) ... p_ON
	Commission            float64    `json:"commission"`               // 手数料 (APIからのレスポンスには含まれない)
	ExpireAt              time.Time  `json:"expire_at"`                // 注文有効期限 ... sOrderExpireDay (APIには YYYYMMDD 形式で渡す)
	CreatedAt             time.Time  `json:"created_at"`               // 注文作成日時 (APIからのレスポンス) ... p_ED (営業日 YYYYMMDD), p_EXDT (通知日時 YYYYMMDDHHMMSS)
	UpdatedAt             time.Time  `json:"updated_at"`               // 注文最終更新日時
	ExecutionType         string     `json:"execution_type"`           // 執行条件 (寄付: "opening", 引け: "closing", 成行: "market", 指値: "limit" など) ... sCondition
	TradeType             string     `json:"trade_type"`               // 取引区分 (現物: "spot", 信用新規: "credit_open", 信用返済: "credit_close") ... sGenkinShinyouKubun, sBaibaiKubun (現引:7, 現渡:5)
	MarketCode            string     `json:"market_code"`              // 市場コード ... sSizyouC
	Positions             []Position `json:"positions"`                // 信用返済時の建玉情報 ... aCLMKabuHensaiData (sTategyokuNumber, sTatebiZyuni, sOrderSuryou)
	AfterTriggerOrderType string     `json:"after_trigger_order_type"` // "market" or "limit"  (トリガー後注文種別) ... sGyakusasiOrderType, sGyakusasiPrice
	AfterTriggerPrice     float64    `json:"after_trigger_price"`      //  (トリガー後指値) ... sGyakusasiPrice
	NotificationType      string     `json:"notification_type"`        // 追加: p_NT (通知種別)
}

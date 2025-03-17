// internal/domain/strct_order.go
package domain

import (
	"time"

	"gorm.io/gorm" // GORM をインポート
)

// Order 株式の注文
type Order struct {
	gorm.Model                 // ID, CreatedAt, UpdatedAt, DeletedAt を追加
	UUID             string    `gorm:"type:varchar(255);unique;not null"` // UUID
	Symbol           string    `gorm:"type:varchar(10);not null"`         // p_IC 銘柄コード
	Side             string    `gorm:"type:varchar(5);not null"`          // "long" or "short"  p_BBKB 売買区分
	OrderType        string    `gorm:"type:varchar(50);not null"`         // 注文タイプ ("market", "limit", "stop" など)
	Quantity         int       `gorm:"not null"`                          // p_CRSR 注文数量
	Price            float64   // p_CRPR 注文価格 (指値の場合)
	TriggerPrice     float64   // 逆指値トリガー価格
	FilledQuantity   int       `gorm:"default:0"` // p_CREXSR 約定済数量
	AveragePrice     float64   // p_EXPR 約定値段
	Status           string    `gorm:"type:varchar(50);not null"` // p_ODST 注文ステータス
	TachibanaOrderID string    `gorm:"type:varchar(255)"`         // p_ON 注文番号
	Commission       float64   //手数料
	ExpireAt         time.Time // p_LMIT 有効期限 (YYYYMMDD)

	//DBに保存不要なフィールドには、gorm:"-"
	Positions                      []Position `gorm:"-"`                                          // ポジション情報 (API からは取得しない)
	ExecutionType                  string     `json:"execution_type"           gorm:"-"`          // p_CRSJ 執行条件
	TradeType                      string     `json:"trade_type"               gorm:"-"`          // p_THKB 取引区分
	MarketCode                     string     `json:"market_code"              gorm:"-"`          // p_MC 市場コード
	AfterTriggerOrderType          string     `json:"after_trigger_order_type" gorm:"-"`          // 逆指値トリガー後の注文種別 ("market" or "limit")
	AfterTriggerPrice              float64    `json:"after_trigger_price"      gorm:"-"`          // 逆指値トリガー後の価格
	NotificationType               string     `json:"notification_type"        gorm:"-"`          // p_NT 通知種別
	BusinessDate                   string     `json:"business_date"            gorm:"-"`          // p_ED 営業日 (YYYYMMDD)
	ParentOrderNumber              string     `json:"parent_order_number"      gorm:"-"`          // p_OON 親注文番号
	ProductType                    int        `json:"product_type"             gorm:"-"`          // p_ST 商品種別
	PriceType                      string     `json:"price_type"               gorm:"-"`          // p_CRPRKB 注文値段区分
	CanceledQuantity               int        `json:"canceled_quantity"        gorm:"-"`          // p_CRTKSR 取消数量
	ExpiredQuantity                int        `json:"expired_quantity"         gorm:"-"`          // p_CREPSR 失効数量
	CarryOverFlag                  string     `json:"carry_over_flag"           gorm:"-"`         // p_KOFG 繰越フラグ
	ModifyCancelStatus             string     `json:"modify_cancel_status"     gorm:"-"`          // p_TTST 訂正取消ステータス
	ExecutionStatus                string     `json:"execution_status"         gorm:"-"`          // p_EXST 約定ステータス
	TaxCategory                    string     `json:"tax_category"              gorm:"-"`         // p_JKK 譲渡益課税区分
	Channel                        string     `json:"channel"                  gorm:"-"`          // p_CHNL チャネル
	ExchangeInvalidationReasonCode string     `json:"exchange_invalidation_reason_code" gorm:"-"` // p_EPRC 失効理由コード(取引所からの値)
	ExchangeExecutionQuantity      int        `json:"exchange_execution_quantity"      gorm:"-"`  // p_EXSR 約定数量(取引所からの値)
	ExchangeErrorCode              string     `json:"exchange_error_code"              gorm:"-"`  // p_EXRC 取引所エラーコード(取引所からの値)
	SymbolName                     string     `json:"symbol_name"                      gorm:"-"`  // p_IN 銘柄名称
}

// TableName overrides the table name
func (Order) TableName() string {
	return "orders"
}

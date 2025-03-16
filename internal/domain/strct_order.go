// internal/domain/strct_order.go
package domain

import (
	"time"

	"gorm.io/gorm" // GORM をインポート
)

// Order 株式の注文
type Order struct {
	gorm.Model              // ID, CreatedAt, UpdatedAt, DeletedAt を追加
	UUID             string `gorm:"type:varchar(255);unique;not null"`
	Symbol           string `gorm:"type:varchar(10);not null"`
	Side             string `gorm:"type:varchar(5);not null"` // "long" or "short"
	OrderType        string `gorm:"type:varchar(50);not null"`
	Quantity         int    `gorm:"not null"`
	Price            float64
	TriggerPrice     float64
	FilledQuantity   int `gorm:"default:0"`
	AveragePrice     float64
	Status           string `gorm:"type:varchar(50);not null"`
	TachibanaOrderID string `gorm:"type:varchar(255)"`
	Commission       float64
	ExpireAt         time.Time
	//DBに保存不要なフィールドには、gorm:"-"
	ExecutionType                  string  `json:"execution_type"           gorm:"-"` // 執行条件
	TradeType                      string  `json:"trade_type"               gorm:"-"` // 取引区分
	MarketCode                     string  `json:"market_code"              gorm:"-"` // 市場コード
	AfterTriggerOrderType          string  `json:"after_trigger_order_type" gorm:"-"` // "market" or "limit"
	AfterTriggerPrice              float64 `json:"after_trigger_price"      gorm:"-"`
	NotificationType               string  `json:"notification_type"        gorm:"-"`          // 追加: p_NT (通知種別)
	BusinessDate                   string  `json:"business_date"                    gorm:"-"`  // 営業日 (p_ED) // 追加
	ParentOrderNumber              string  `json:"parent_order_number"              gorm:"-"`  // 親注文番号 (p_OON) // 追加
	ProductType                    int     `json:"product_type"                     gorm:"-"`  // 商品種別 (p_ST) // 追加
	PriceType                      string  `json:"price_type"                     gorm:"-"`    // 注文値段区分 (p_CRPRKB) // 追加
	CanceledQuantity               int     `json:"canceled_quantity"                gorm:"-"`  // 取消数量 (p_CRTKSR) // 追加
	ExpiredQuantity                int     `json:"expired_quantity"                 gorm:"-"`  // 失効数量 (p_CREPSR) // 追加
	CarryOverFlag                  string  `json:"carry_over_flag"                   gorm:"-"` // 繰越フラグ (p_KOFG) // 追加
	ModifyCancelStatus             string  `json:"modify_cancel_status"             gorm:"-"`  // 訂正取消ステータス (p_TTST) // 追加
	ExecutionStatus                string  `json:"execution_status"                gorm:"-"`   // 約定ステータス (p_EXST) // 追加
	TaxCategory                    string  `json:"tax_category"                    gorm:"-"`   // 譲渡益課税区分 (p_JKK) // 追加
	Channel                        string  `json:"channel"                          gorm:"-"`  //チャネル(p_CHNL) //追加
	ExchangeInvalidationReasonCode string  `json:"exchange_invalidation_reason_code" gorm:"-"` // 失効理由コード(取引所からの値) (p_EPRC) // 追加
	ExchangeExecutionQuantity      int     `json:"exchange_execution_quantity"      gorm:"-"`  // 約定数量(取引所からの値) (p_EXSR) // 追加
	ExchangeErrorCode              string  `json:"exchange_error_code"              gorm:"-"`  // 取引所エラーコード(取引所からの値) (p_EXRC) // 追加
	SymbolName                     string  `json:"symbol_name"                      gorm:"-"`  // 銘柄名称(p_IN) //追加
}

// TableName overrides the table name
func (Order) TableName() string {
	return "orders"
}

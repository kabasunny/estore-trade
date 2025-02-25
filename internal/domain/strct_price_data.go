// internal/domain/strct_price_data.go
package domain

// PriceData は株価データ (日付、始値、高値、安値、終値、出来高) を表す構造体
type PriceData struct {
	Date      string  `json:"sZyoukaiDay"`
	Open      float64 `json:"sHajimene,string"`
	High      float64 `json:"sTakane,string"`
	Low       float64 `json:"sYasune,string"`
	Close     float64 `json:"sOwarine,string"`
	Volume    int     `json:"sDekidaka,string"`
	IssueCode string  `json:"sIssueCode"`
}

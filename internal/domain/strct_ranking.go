// internal/domain/ranking.go
package domain

import "time"

// Ranking はランキング情報を表す構造体
type Ranking struct {
	ID           int       `json:"id"` // DBのID
	Rank         int       `json:"rank"`
	IssueCode    string    `json:"issue_code"`
	TradingValue float64   `json:"trading_value"`
	CreatedAt    time.Time `json:"created_at"`
}

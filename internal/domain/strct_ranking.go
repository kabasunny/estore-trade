// internal/domain/ranking.go
package domain

import "time"

// 売買代金ランキングの情報を表す構造体
type Ranking struct {
	Rank         int       `json:"rank"`
	IssueCode    string    `json:"issue_code"`
	TradingValue float64   `json:"trading_value"`
	CreatedAt    time.Time `json:"created_at"` // ランキング生成日時
}

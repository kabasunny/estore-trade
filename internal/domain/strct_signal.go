// internal/domain/signal.go
package domain

import "time"

// 売買シグナルの情報を表す構造体
type Signal struct {
	ID        int       `json:"id"` // シグナルID
	IssueCode string    `json:"issue_code"`
	Side      string    `json:"side"`       // "buy" or "sell"
	Priority  int       `json:"priority"`   // 優先度 (例: 1, 2, 3, ...)
	CreatedAt time.Time `json:"created_at"` // シグナル生成日時
}

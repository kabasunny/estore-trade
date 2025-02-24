// internal/domain/ranking.go
package domain

// 取引対象とする銘柄の情報を表す構造体
type TargetIssue struct {
	IssueCode string `json:"issue_code"`
}

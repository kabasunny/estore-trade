// internal/batch/ranking/util_create_target_issue_list.go
package ranking

import (
	"estore-trade/internal/domain"
)

// CreateTargetIssueList ランキングから取引対象銘柄のリストを作成する(大文字に)
func CreateTargetIssueList(ranking []domain.Ranking, limit int) []domain.TargetIssue {
	var targetIssues []domain.TargetIssue
	// TODO: 実際にはランキング上位N件を抽出
	for i := 0; i < len(ranking) && i < limit; i++ {
		targetIssues = append(targetIssues, domain.TargetIssue{IssueCode: ranking[i].IssueCode})
	}

	return targetIssues
}

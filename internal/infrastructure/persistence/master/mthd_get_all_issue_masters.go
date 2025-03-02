// internal/infrastructure/persistence/master/mthd_get_all_issue_masters.go
package master

import (
	"context"
	"estore-trade/internal/domain"
)

// 全てのIssueMasterをDBから取得
func (r *masterDataRepository) getAllIssueMasters(ctx context.Context) ([]domain.IssueMaster, error) {
	query := `SELECT issue_code, issue_name, trading_unit, tokutei_f FROM issue_masters` //仮
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []domain.IssueMaster
	for rows.Next() {
		var issue domain.IssueMaster
		if err := rows.Scan(&issue.IssueCode, &issue.IssueName, &issue.TradingUnit, &issue.TokuteiF); err != nil { //仮
			return nil, err
		}
		issues = append(issues, issue)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return issues, nil
}

// internal/infrastructure/persistence/master/mthd_get_all_issue_market_regulations.go
package master

import (
	"context"
	"estore-trade/internal/domain"
)

// IssueMarketRegulationをDBから取得
func (r *masterDataRepository) getAllIssueMarketRegulations(ctx context.Context) ([]domain.IssueMarketRegulation, error) {
	query := `
        SELECT issue_code, listed_market, stop_kubun, genbutu_urituke,
               seido_sinyou_sinki_uritate, ippan_sinyou_sinki_uritate, sinyou_syutyu_kubun
        FROM issue_market_regulations` // 仮のテーブル名とカラム名
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issueRegulations []domain.IssueMarketRegulation
	for rows.Next() {
		var ir domain.IssueMarketRegulation
		err := rows.Scan(
			&ir.IssueCode,
			&ir.ListedMarket,
			&ir.StopKubun,
			&ir.GenbutuUrituke,
			&ir.SeidoSinyouSinkiUritate,
			&ir.IppanSinyouSinkiUritate,
			&ir.SinyouSyutyuKubun,
		)
		if err != nil {
			return nil, err
		}
		issueRegulations = append(issueRegulations, ir)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return issueRegulations, nil
}

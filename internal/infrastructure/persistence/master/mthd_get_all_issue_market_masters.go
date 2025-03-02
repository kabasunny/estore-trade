// internal/infrastructure/persistence/master/mthd_get_all_issue_market_masters.go
package master

import (
	"context"
	"estore-trade/internal/domain"
)

// IssueMarketMasterをDBから取得
func (r *masterDataRepository) getAllIssueMarketMasters(ctx context.Context) ([]domain.IssueMarketMaster, error) {
	query := `
        SELECT issue_code, market_code, price_range_min, price_range_max, sinyou_c, previous_close,
               issue_kubun_c, zyouzyou_kubun, call_price_unit_number, call_price_unit_number_yoku
        FROM issue_market_masters` // 仮のテーブル名とカラム名
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issueMarkets []domain.IssueMarketMaster
	for rows.Next() {
		var im domain.IssueMarketMaster
		err := rows.Scan(
			&im.IssueCode,
			&im.MarketCode,
			&im.PriceRangeMin,
			&im.PriceRangeMax,
			&im.SinyouC,
			&im.PreviousClose,
			&im.IssueKubunC,
			&im.ZyouzyouKubun,
			&im.CallPriceUnitNumber,
			&im.CallPriceUnitNumberYoku,
		)
		if err != nil {
			return nil, err
		}
		issueMarkets = append(issueMarkets, im)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return issueMarkets, nil
}

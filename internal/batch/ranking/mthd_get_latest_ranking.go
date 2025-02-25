// internal/batch/ranking/mthd_get_latest_ranking.go
package ranking

import (
	"context"
	"estore-trade/internal/domain"
)

func (r *rankingRepository) GetLatestRanking(ctx context.Context, limit int) ([]domain.Ranking, error) {
	query := `
        SELECT rank, issue_code, trading_value, created_at
        FROM rankings
        ORDER BY created_at DESC, rank ASC
        LIMIT $1
    `
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranking []domain.Ranking
	for rows.Next() {
		var item domain.Ranking
		if err := rows.Scan(&item.Rank, &item.IssueCode, &item.TradingValue, &item.CreatedAt); err != nil {
			return nil, err
		}
		ranking = append(ranking, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ranking, nil
}

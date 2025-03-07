package ranking

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

func (r *rankingRepository) GetLatestRanking(ctx context.Context, limit int) ([]domain.Ranking, error) {
	query := `
    SELECT rank, issue_code, trading_value, created_at
    FROM rankings
    ORDER BY created_at DESC, id DESC, rank ASC
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
		var createdAtStr string // created_at を文字列として受け取る
		if err := rows.Scan(&item.Rank, &item.IssueCode, &item.TradingValue, &createdAtStr); err != nil {
			return nil, err
		}

		// 文字列を time.Time に変換 (実際のDBのフォーマットに合わせる)
		// 例: "2006-01-02 15:04:05.999999999-07:00"
		createdAt, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", createdAtStr)
		if err != nil {
			return nil, err // パースに失敗したらエラー
		}
		item.CreatedAt = createdAt

		ranking = append(ranking, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ranking, nil
}

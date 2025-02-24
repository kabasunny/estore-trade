// internal/batch/ranking/repository.go
package ranking

import (
	"context"
	"estore-trade/internal/domain"
)

// ランキングデータの永続化を担当するインターフェース
type RankingRepository interface {
	SaveRanking(ctx context.Context, ranking []domain.Ranking) error
	// 必要に応じて、GetRanking, GetLatestRanking などのメソッドを追加
}

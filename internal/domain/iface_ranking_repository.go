// internal/domain/iface_ranking_repository.go
package domain // パッケージ名を domain に変更

import (
	"context"
)

// RankingRepository は、ランキングデータの永続化を担当するインターフェース
type RankingRepository interface {
	SaveRanking(ctx context.Context, ranking []Ranking) error
	GetLatestRanking(ctx context.Context, limit int) ([]Ranking, error)
}

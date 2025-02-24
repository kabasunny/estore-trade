// internal/batch/ranking/util_calculate_ranking.go
package ranking

import (
	"context"
	"estore-trade/internal/domain"

	//"estore-trade/internal/infrastructure/persistence/tachibana" //tachibanaClientを使うため　不要
	"sort" //追加
	"time"
)

// CalculateRanking は売買代金ランキングを計算し、上位N件の銘柄リストを返す
func CalculateRanking(ctx context.Context, marketData []marketDataItem) ([]domain.Ranking, error) { //tachibanaClientを削除
	// 3. 売買代金の計算とランキング作成
	var ranking []domain.Ranking
	for _, data := range marketData {
		// TODO: 実際には株価と出来高を掛け合わせて売買代金を計算
		tradingValue := data.Price * float64(data.Volume) // 仮の計算
		ranking = append(ranking, domain.Ranking{
			IssueCode:    data.IssueCode,
			TradingValue: tradingValue,
			CreatedAt:    time.Now(),
		})
	}

	// 売買代金で降順にソート
	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].TradingValue > ranking[j].TradingValue // 大小を逆にすることで降順
	})

	// ランキングに順位を付与
	for i := range ranking {
		ranking[i].Rank = i + 1
	}

	return ranking, nil
}

// internal/batch/ranking/util_calculate_ranking.go
package ranking

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"sort"
	"time"
)

// CalculateRanking は売買代金ランキングを計算します。
func CalculateRanking(ctx context.Context, client tachibana.TachibanaClient, issueCodes []string) ([]domain.Ranking, error) {
	const days = 5 // 過去5日間のデータを使用

	var ranking []domain.Ranking
	for _, issueCode := range issueCodes {
		// 過去5日間の株価データを取得 (TachibanaClient を使用)
		// 実際には、GetPriceHistory のようなメソッドが必要になる
		// prices, err := client.GetPriceHistory(ctx, issueCode, days)
		// if err != nil {
		//     return nil, fmt.Errorf("failed to get price history for %s: %w", issueCode, err)
		// }

		// 仮のデータ (TachibanaClient に GetPriceHistory が実装されるまで)
		prices := make([]domain.PriceData, days)
		for i := range prices {
			prices[i] = domain.PriceData{
				Date:   time.Now().AddDate(0, 0, -i).Format("20060102"),
				Close:  1000.0 + float64(i), // 仮の株価
				Volume: 1000 + i,            // 仮の出来高
			}
		}

		// 売買代金を計算
		var totalTradingValue float64
		for _, price := range prices {
			tradingValue := price.Close * float64(price.Volume)
			totalTradingValue += tradingValue
		}
		averageTradingValue := totalTradingValue / float64(days)

		ranking = append(ranking, domain.Ranking{
			IssueCode:    issueCode,
			TradingValue: averageTradingValue,
			CreatedAt:    time.Now(),
		})
	}

	// 売買代金で降順にソート
	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].TradingValue > ranking[j].TradingValue
	})

	// ランキングに順位を付与
	for i := range ranking {
		ranking[i].Rank = i + 1
	}

	return ranking, nil
}

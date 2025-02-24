// internal/batch/ranking/calculate.go
package ranking

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana" //tachibanaClientを使うため
	"time"
)

// CalculateRanking は売買代金ランキングを計算し、上位N件の銘柄リストを返す
func CalculateRanking(ctx context.Context, client tachibana.TachibanaClient) ([]domain.Ranking, error) {
	// 1. 全銘柄コードの取得 (TODO: 実際にはマスタデータから取得)
	//    allIssueCodes := getAllIssueCodes(client) // 全銘柄コードを取得する関数 (仮)
	allIssueCodes := []string{"7203", "8306", "9432"} // トヨタ、UFJ、NTTのコードを仮で使う

	// 2. 株価・出来高の取得 (TODO: 実際には CLMMfdsGetMarketPrice を使う)
	marketData, err := getMarketData(ctx, client, allIssueCodes)
	if err != nil {
		return nil, err
	}

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

	// TODO: 売買代金で降順にソート (ここでは省略)

	return ranking, nil
}

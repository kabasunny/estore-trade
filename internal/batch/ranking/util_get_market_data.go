// internal/batch/ranking/calculate.go
package ranking

import (
	"context"
	"estore-trade/internal/infrastructure/persistence/tachibana" //tachibanaClientを使うため
	"fmt"
)

// 仮実装.　tachibanaClientからデータを取得する
func getMarketData(ctx context.Context, client tachibana.TachibanaClient, issueCodes []string) ([]marketDataItem, error) {
	// ここに実際には、tachibanaClient.CLMMfdsGetMarketPrice()を呼び出してデータを取得するコードを記述
	// (今はモックデータ)
	requestURL, err := client.GetRequestURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get request URL: %w", err)
	}
	fmt.Println("requestURL:", requestURL)

	var marketData []marketDataItem //仮
	return marketData, nil          // 仮のデータを返す
}

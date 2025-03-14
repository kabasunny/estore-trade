package tachibana_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder_CloseCreditMarketSell(t *testing.T) {
	t.Run("正常系: 信用成行売り返済（信用成行買いポジションに対応）", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		ctx := context.Background()

		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// 1. エントリー注文 (信用成行買い)
		entryOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long",
			OrderType:  "market",
			Condition:  "credit_open", // 信用新規
			Quantity:   100,
			MarketCode: "00", // 東証
		}
		placedEntryOrder, err := client.PlaceOrder(ctx, entryOrder)
		assert.NoError(t, err)
		assert.NotEmpty(t, placedEntryOrder.TachibanaOrderID)

		time.Sleep(3 * time.Second) // 約定を待つ

		// 2. GetPositions で建玉情報を取得 (リトライ処理付き)
		var positions []domain.Position
		maxRetries := 5
		retryInterval := 2 * time.Second
		for i := 0; i < maxRetries; i++ {
			fmt.Printf("GetPositions retry: %d\n", i+1) // リトライ回数をログ出力
			positions, err = client.GetPositions(ctx)
			if err != nil {
				fmt.Printf("GetPositions error: %v\n", err) // エラー内容をログ出力
				time.Sleep(retryInterval)
				continue
			}
			if len(positions) > 0 { // 建玉が取得できたら(空でなければ)
				break // リトライループを抜ける
			}
			time.Sleep(retryInterval) // 少し待ってからリトライ
		}
		assert.NoError(t, err)                                               // リトライ後もエラーならテスト失敗
		assert.NotEmpty(t, positions, "No positions returned after retries") // リトライ後も空なら失敗

		// 3. 信用買いの建玉をより厳密に探して、成行売りで返済
		found := false
		today := time.Now().Format("20060102") // 今日の日付を YYYYMMDD 形式で取得
		for _, p := range positions {
			// 検索条件: 銘柄、Side、数量、建日
			if p.Symbol == entryOrder.Symbol &&
				p.Side == "long" &&
				p.Quantity >= entryOrder.Quantity && // 建玉の数量 >= 注文数量
				p.OpenDate.Format("20060102") == today { // 建日が今日

				found = true
				exitOrder := &domain.Order{
					Symbol:     p.Symbol,
					Side:       "short",
					OrderType:  "credit_close_market",
					Quantity:   entryOrder.Quantity, // エントリー時の数量で返済
					MarketCode: "00",
					Positions: []domain.Position{
						{
							ID:       p.ID,
							Quantity: entryOrder.Quantity, // エントリー時の数量
						},
					},
				}

				placedExitOrder, err := client.PlaceOrder(ctx, exitOrder)
				assert.NoError(t, err)
				assert.NotNil(t, placedExitOrder)
				assert.NotEmpty(t, placedExitOrder.TachibanaOrderID)
				assert.Equal(t, "pending", placedExitOrder.Status)

				time.Sleep(1 * time.Second)
				break //  見つけたらループを抜ける
			}
		}
		assert.True(t, found, "Matching position not found")
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order/place_order_close_credit_market_sell_test.go

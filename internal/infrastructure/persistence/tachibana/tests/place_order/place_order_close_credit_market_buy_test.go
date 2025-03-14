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

func TestPlaceOrder_CloseCreditMarketBuy(t *testing.T) {
	t.Run("正常系: 信用成行買い返済（信用成行売りポジションに対応）", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		ctx := context.Background()

		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// 1. エントリー注文 (信用成行売り)
		entryOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "short",
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
			fmt.Printf("GetPositions retry: %d\n", i+1)
			positions, err = client.GetPositions(ctx)
			if err != nil {
				fmt.Printf("GetPositions error: %v\n", err)
				time.Sleep(retryInterval)
				continue
			}
			if len(positions) > 0 {
				break
			}
			time.Sleep(retryInterval)
		}
		assert.NoError(t, err)
		assert.NotEmpty(t, positions, "No positions returned after retries")

		// 3. 信用売りの建玉を探して、成行買いで返済
		found := false
		today := time.Now().Format("20060102") // 今日の日付を取得
		for _, p := range positions {
			// 検索条件: 銘柄、Side、数量、建日
			if p.Symbol == entryOrder.Symbol &&
				p.Side == "short" && // 売り建玉
				p.Quantity >= entryOrder.Quantity && // 建玉の数量 >= 注文数量
				p.OpenDate.Format("20060102") == today { // 建日が今日

				found = true
				exitOrder := &domain.Order{
					Symbol:     p.Symbol,
					Side:       "long",                // 買い
					OrderType:  "credit_close_market", // 信用返済成行
					Quantity:   entryOrder.Quantity,   // エントリー時の数量で返済
					MarketCode: "00",
					Positions: []domain.Position{
						{
							ID:       p.ID,                // 建玉番号
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
				break
			}
		}
		assert.True(t, found, "Matching position not found")
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order/place_order_close_credit_market_buy_test.go

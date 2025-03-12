package tachibana_test

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder_CloseCreditMarketBuy(t *testing.T) {
	t.Run("正常系: 信用成行買い返済（信用成行売りポジションに対応）", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		// 1. エントリー注文 (信用成行売り)
		entryOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "short",
			OrderType:  "market",
			Condition:  "credit_open", // 信用新規
			Quantity:   100,
			MarketCode: "00", // 東証
		}
		placedEntryOrder, err := client.PlaceOrder(context.Background(), entryOrder)
		assert.NoError(t, err)
		assert.NotEmpty(t, placedEntryOrder.TachibanaOrderID)

		time.Sleep(3 * time.Second) // 約定を待つ

		// 2. GetPositions で建玉情報を取得
		positions, err := client.GetPositions(context.Background())
		assert.NoError(t, err)

		// 3. 信用売りの建玉を探して、成行買いで返済
		for _, p := range positions {
			if p.Symbol == entryOrder.Symbol && p.Side == "short" && p.Quantity == entryOrder.Quantity {
				exitOrder := &domain.Order{
					Symbol:     p.Symbol,
					Side:       "long",                // 買い
					OrderType:  "credit_close_market", // 信用返済成行
					Condition:  "",
					Quantity:   p.Quantity, // 全量返済
					MarketCode: "00",
					Positions: []domain.Position{
						{
							ID:       p.ID,       // 建玉番号
							Quantity: p.Quantity, // 建玉数量 (全量)
						},
					},
				}

				placedExitOrder, err := client.PlaceOrder(context.Background(), exitOrder)
				if err != nil {
					t.Fatalf("PlaceOrder for exitOrder failed: %v", err)
				}
				assert.NotNil(t, placedExitOrder)
				assert.NotEmpty(t, placedExitOrder.ID)
				assert.Equal(t, "pending", placedExitOrder.Status)

				time.Sleep(1 * time.Second) // 少し待つ
			}
		}
	})
}

// internal/infrastructure/persistence/tachibana/tests/place_order_close_credit_stop_sell_test.go
// close: 返済注文
// credit: 信用取引
// stop: 逆指値
// sell: 売り
package tachibana_test

import (
	"context"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder_CloseCreditStopSell(t *testing.T) {
	t.Run("正常系: 信用逆指値売り（信用成行買いポジションに対応）", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		// 1. エントリー注文 (信用成行買い)
		//    テスト用の建玉を事前に作成
		entryOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long",
			OrderType:  "market",
			Condition:  "credit_open", // 信用新規
			Quantity:   100,
			MarketCode: "00", // 東証
		}
		placedEntryOrder, err := client.PlaceOrder(context.Background(), entryOrder)
		assert.NoError(t, err)
		assert.NotEmpty(t, placedEntryOrder.TachibanaOrderID)

		time.Sleep(3 * time.Second)

		// 2. GetPositions で建玉情報を取得 (既存の建玉 + 今回の建玉)
		positions, err := client.GetPositions(context.Background())
		assert.NoError(t, err)

		// 3. 信用買いの建玉を探して、逆指値売りで返済
		//    条件を p.Quantity == entryOrder.Quantity に変更
		for _, p := range positions {
			// 今回のエントリー注文、および、過去の買い建玉が対象
			if p.Symbol == entryOrder.Symbol && p.Side == "long" && p.Quantity == entryOrder.Quantity { //ここを修正
				exitOrder := &domain.Order{
					Symbol:       p.Symbol,
					Side:         "short", // 売り
					OrderType:    "credit_close_stop",
					Condition:    "",
					Quantity:     p.Quantity, // 全量返済
					TriggerPrice: 9500.0,     // 逆指値トリガー価格
					MarketCode:   "00",
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
				assert.NotEmpty(t, placedExitOrder.UUID)           // 注文ID
				assert.Equal(t, "pending", placedExitOrder.Status) // ステータス

				time.Sleep(1 * time.Second) // 少し待つ (逆指値注文はすぐには約定しない)
			}
		}
	})

	// 異常系テストケース (API へのリクエストは送信されない)

	// t.Run("異常系: 無効な Side", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil) // Login を追加
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())
	// 	order := &domain.Order{
	// 		OrderType:    "credit_close_stop",
	// 		Side:         "invalid_side", // 無効な値
	// 		Symbol:       "7974",
	// 		Quantity:     100,
	// 		TriggerPrice: 9500.0,
	// 		Positions: []domain.Position{
	// 			{ID: "test-position-id", Quantity: 100},
	// 		},
	// 		MarketCode: "00",
	// 	}
	// 	_, err = client.PlaceOrder(context.Background(), order)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "invalid order side")
	// })

	// t.Run("異常系: 無効な Quantity", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil) // Login を追加
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())
	// 	order := &domain.Order{
	// 		OrderType:    "credit_close_stop",
	// 		Side:         "sell",
	// 		Symbol:       "7974",
	// 		Quantity:     0, // 無効な値
	// 		TriggerPrice: 9500.0,
	// 		Positions: []domain.Position{
	// 			{ID: "test-position-id", Quantity: 100},
	// 		},
	// 		MarketCode: "00",
	// 	}
	// 	_, err = client.PlaceOrder(context.Background(), order)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "invalid order quantity")
	// })

	// t.Run("異常系: 無効な TriggerPrice", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil) // Login を追加
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())
	// 	order := &domain.Order{
	// 		OrderType:    "credit_close_stop",
	// 		Side:         "sell",
	// 		Symbol:       "7974",
	// 		Quantity:     100,
	// 		TriggerPrice: 0, // 無効な値
	// 		Positions: []domain.Position{
	// 			{ID: "test-position-id", Quantity: 100},
	// 		},
	// 		MarketCode: "00",
	// 	}
	// 	_, err = client.PlaceOrder(context.Background(), order)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "invalid trigger price")
	// })

	// t.Run("異常系: Positions が空", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil) // Login を追加
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())

	// 	order := &domain.Order{
	// 		OrderType:    "credit_close_stop",
	// 		Side:         "sell",
	// 		Symbol:       "7974",
	// 		Quantity:     100, // ここは有効な値にしておく
	// 		TriggerPrice: 9500.0,
	// 		Positions:    []domain.Position{}, // 空の Positions スライス
	// 		MarketCode:   "00",
	// 	}
	// 	_, err = client.PlaceOrder(context.Background(), order)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "no positions specified for credit close order") // 期待されるエラー
	// })

	// t.Run("異常系: Positions 内の Quantity が無効", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil) // Login を追加
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())

	// 	order := &domain.Order{
	// 		OrderType:    "credit_close_stop",
	// 		Side:         "sell",
	// 		Symbol:       "7974",
	// 		Quantity:     1, // 全体の数量は有効な値
	// 		TriggerPrice: 9500.0,
	// 		Positions: []domain.Position{
	// 			{ID: "test-position-id", Quantity: 0}, // 無効な値
	// 		},
	// 		MarketCode: "00",
	// 	}
	// 	_, err = client.PlaceOrder(context.Background(), order)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "invalid position quantity") // 期待されるエラー
	// })
	// t.Run("異常系: Positions 内の ID が空", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil) // Login を追加
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())
	// 	order := &domain.Order{
	// 		OrderType:    "credit_close_stop",
	// 		Side:         "sell",
	// 		Symbol:       "7974",
	// 		Quantity:     100,
	// 		TriggerPrice: 9500.0,
	// 		Positions: []domain.Position{
	// 			{ID: "", Quantity: 100}, // 無効な値（空文字列）
	// 		},
	// 		MarketCode: "00",
	// 	}
	// 	_, err = client.PlaceOrder(context.Background(), order)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "invalid position ID") // 期待されるエラー
	// })
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order_close_credit_stop_sell_test.go

// internal/infrastructure/persistence/tachibana/tests/place_order_close_credit_stop_buy_test.go
// close: 返済注文
// credit: 信用取引
// stop: 逆指値
// buy: 買い
package tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder_CreditCloseStopBuy(t *testing.T) {
	// t.Run("正常系: 信用逆指値買い（信用成行売りポジションに対応）", func(t *testing.T) {
	// 	client := tachibana.CreateTestClient(t, nil)
	// 	err := client.Login(context.Background(), nil)
	// 	assert.NoError(t, err)
	// 	defer client.Logout(context.Background())

	// 	// 1. エントリー注文 (信用成行売り)
	// 	entryOrder := &domain.Order{
	// 		Symbol:     "7974", // 例: 任天堂
	// 		Side:       "sell",
	// 		OrderType:  "market",
	// 		Condition:  "credit_open", // 信用新規
	// 		Quantity:   100,
	// 		MarketCode: "00", // 東証
	// 	}
	// 	_, err = client.PlaceOrder(context.Background(), entryOrder)
	// 	assert.NoError(t, err)

	// 	// 2. 約定のシミュレーション (3秒待機) - 本来は約定イベントを監視
	// 	time.Sleep(3 * time.Second)

	// 	// 3. GetPositions で建玉情報を取得
	// 	positions, err := client.GetPositions(context.Background())
	// 	assert.NoError(t, err)
	// 	// 取得した建玉の中に、先ほど建てた売り建玉があることを確認
	// 	found := false
	// 	var targetPosition domain.Position
	// 	for _, p := range positions {
	// 		if p.Symbol == entryOrder.Symbol && p.Side == "short" && p.Quantity == entryOrder.Quantity {
	// 			found = true
	// 			targetPosition = p
	// 			break
	// 		}
	// 	}
	// 	assert.True(t, found, "建てたはずの売り建玉が見つかりません")

	// 	// 4. エグジット注文 (信用逆指値買い)
	// 	exitOrder := &domain.Order{
	// 		Symbol:       targetPosition.Symbol,             // 建玉の銘柄コード
	// 		Side:         "buy",                             // 買い
	// 		OrderType:    "credit_close_stop",               // 信用返済逆指値
	// 		Condition:    "",                                // 信用返済
	// 		Quantity:     targetPosition.Quantity,           // 建玉の数量
	// 		TriggerPrice: 10500.0,                           // 逆指値のトリガー価格 (現在価格より高い価格)
	// 		MarketCode:   "00",                              // 東証
	// 		Positions:    []domain.Position{targetPosition}, // 返済対象の建玉情報
	// 	}
	// 	placedOrder, err := client.PlaceOrder(context.Background(), exitOrder)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, placedOrder)
	// 	assert.NotEmpty(t, placedOrder.ID)
	// 	assert.Equal(t, "pending", placedOrder.Status)

	// 	// (オプション) 約定のシミュレーション (数秒待機)
	// 	time.Sleep(3 * time.Second)
	// 	// ここで、GetOrderStatus などを使って注文状況を確認することも可能
	// })

	// 異常系のテストケースは、この後必要に応じて追加
	// 異常系テストケース (API へのリクエストは送信されない)

	t.Run("異常系: 無効な Side", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil) // Login を追加
		assert.NoError(t, err)
		defer client.Logout(context.Background())
		order := &domain.Order{
			OrderType:    "credit_close_stop",
			Side:         "invalid_side", // 無効な値
			Symbol:       "7974",
			Quantity:     100,
			TriggerPrice: 10500.0,
			Positions: []domain.Position{
				{ID: "test-position-id", Quantity: 100},
			},
			MarketCode: "00",
		}
		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order side")
	})

	t.Run("異常系: 無効な Quantity", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil) // Login を追加
		assert.NoError(t, err)
		defer client.Logout(context.Background())
		order := &domain.Order{
			OrderType:    "credit_close_stop",
			Side:         "buy",
			Symbol:       "7974",
			Quantity:     0, // 無効な値
			TriggerPrice: 10500.0,
			Positions: []domain.Position{
				{ID: "test-position-id", Quantity: 100},
			},
			MarketCode: "00",
		}
		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order quantity")
	})

	t.Run("異常系: 無効な TriggerPrice", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil) // Login を追加
		assert.NoError(t, err)
		defer client.Logout(context.Background())
		order := &domain.Order{
			OrderType:    "credit_close_stop",
			Side:         "buy",
			Symbol:       "7974",
			Quantity:     100,
			TriggerPrice: 0, // 無効な値
			Positions: []domain.Position{
				{ID: "test-position-id", Quantity: 100},
			},
			MarketCode: "00",
		}
		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid trigger price")
	})

	t.Run("異常系: Positions が空", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			OrderType:    "credit_close_stop",
			Side:         "buy",
			Symbol:       "7974",
			Quantity:     100,
			TriggerPrice: 10500.0,
			Positions:    []domain.Position{}, // 空の Positions スライス
			MarketCode:   "00",
		}
		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no positions specified for credit close order") // ここを修正
	})

	t.Run("異常系: Positions 内の Quantity が無効", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			OrderType:    "credit_close_stop",
			Side:         "buy",
			Symbol:       "7974",
			Quantity:     1, // 全体の数量は有効な値
			TriggerPrice: 10500.0,
			Positions: []domain.Position{
				{ID: "test-position-id", Quantity: 0}, // 無効な値
			},
			MarketCode: "00",
		}
		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid position quantity") // ここを修正
	})

	t.Run("異常系: Positions 内の ID が空", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			OrderType:    "credit_close_stop",
			Side:         "buy",
			Symbol:       "7974",
			Quantity:     100,
			TriggerPrice: 10500.0,
			Positions: []domain.Position{
				{ID: "", Quantity: 100}, // 無効な値（空文字列）
			},
			MarketCode: "00",
		}
		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid position ID") // ここを修正
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order_close_credit_stop_buy_test.go

// internal/infrastructure/persistence/tachibana/tests/place_order_close_spot_stop_sell_test.go
// close: 返済注文
// spot: 現物取引
// stop: 逆指値
// sell: 売り
package tachibana_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder_SpotStopSell(t *testing.T) {
	t.Run("正常系: 現物逆指値売り（現物成行買いポジションに対応）", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		// 1. エントリー注文 (現物成行買い)
		entryOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "buy",
			OrderType:  "market",
			Condition:  "", // 現物
			Quantity:   100,
			MarketCode: "00", // 東証
		}
		_, err = client.PlaceOrder(context.Background(), entryOrder)
		assert.NoError(t, err)

		// 2. 約定のシミュレーション (3秒待機)
		time.Sleep(3 * time.Second)

		// 3. エグジット注文 (現物逆指値売り)
		exitOrder := &domain.Order{
			Symbol:       "7974", // 例: 任天堂
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "", // 現物
			Quantity:     100,
			TriggerPrice: 9000.0, // 逆指値条件価格 (例: 現在価格が10000円として、9000円以下になったら)
			Price:        0,      // 逆指値トリガー後、指値で売りたい場合は、ここに価格を設定
			MarketCode:   "00",   // 東証
		}
		placedOrder, err := client.PlaceOrder(context.Background(), exitOrder)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.NotEmpty(t, placedOrder.UUID)
		assert.Equal(t, "pending", placedOrder.Status)

	})

	// 以降、異常系のテストケースを追加

	// API にリクエストが送信 *されない* ケース (エントリー注文不要)
	t.Run("異常系: 無効なSide", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil) // Login 不要

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "invalid_side", // 無効なSide
			OrderType:    "stop",
			Condition:    "",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err := client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order side")
	})

	t.Run("異常系: 無効な数量", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil) // Login 不要

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "",
			Quantity:     0, // 無効な数量
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err := client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order quantity")
	})

	t.Run("異常系: 無効なトリガー価格", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil) // Login 不要

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "",
			Quantity:     100,
			TriggerPrice: 0, // 無効なトリガー価格
			MarketCode:   "00",
		}
		_, err := client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid trigger price")
	})

	// API にリクエストが送信 *される* ケース (エントリー注文が必要)

	t.Run("異常系: 無効な銘柄コード", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		// エントリー注文 (現物成行買い)
		entryOrder := &domain.Order{
			Symbol:     "7974",
			Side:       "buy",
			OrderType:  "market",
			Condition:  "",
			Quantity:   100,
			MarketCode: "00",
		}
		_, err = client.PlaceOrder(context.Background(), entryOrder)
		assert.NoError(t, err)
		time.Sleep(3 * time.Second)

		// エグジット注文 (無効な銘柄コード)
		order := &domain.Order{
			Symbol:       "invalid_code", // 無効な銘柄コード
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API returned an error")
	})

	t.Run("異常系: APIエラー (ログインせずに注文)", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		// ログインしない

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err := client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request URL not found, need to Login")
	})

	t.Run("異常系: コンテキストキャンセル", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		// エントリー注文 (現物成行買い)
		entryOrder := &domain.Order{
			Symbol:     "7974",
			Side:       "buy",
			OrderType:  "market",
			Condition:  "",
			Quantity:   100,
			MarketCode: "00",
		}
		_, err = client.PlaceOrder(context.Background(), entryOrder)
		assert.NoError(t, err)
		time.Sleep(3 * time.Second)

		// エグジット注文 (コンテキストキャンセル)
		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // キャンセル

		_, err = client.PlaceOrder(ctx, order)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled)) // キャンセルエラー
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order_close_spot_stop_sell_test.go

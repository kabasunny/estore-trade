// internal/infrastructure/persistence/tachibana/tests/place_order_credit_market_sell_with_stop_test.go
// credit: 信用取引
// market: 成行注文
// sell: 売り
// with_stop: 逆指値付き
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

func TestPlaceOrder_CreditMarketSellWithStop(t *testing.T) {
	t.Run("正常系: 信用成行き売り（逆指値付き）注文が成功すること", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "7974",        // 例: 任天堂
			Side:         "sell",        // 売り
			OrderType:    "stop",        // 逆指値
			Condition:    "credit_open", // 信用新規
			Quantity:     100,
			TriggerPrice: 9000.0, // 逆指値条件価格 (例: 現在価格が10000円として、9000円以下になったら)
			Price:        0,      // TriggerPriceに達したら、成行き
			MarketCode:   "00",   // 東証
		}

		placedOrder, err := client.PlaceOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.NotEmpty(t, placedOrder.UUID)
		assert.Equal(t, "pending", placedOrder.Status)

		// 1秒待機 (逆指値注文はすぐに約定しない場合があるため)
		time.Sleep(1 * time.Second)
	})

	/// 異常系テストケース
	t.Run("異常系: 無効な銘柄コード", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "invalid_code", // 無効な銘柄コード
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "credit_open",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API returned an error") // エラーメッセージの確認
	})

	t.Run("異常系: 無効なSide", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "invalid_side", // 無効なSide
			OrderType:    "stop",
			Condition:    "credit_open",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order side") // エラーメッセージの確認
	})

	t.Run("異常系: 無効な数量", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "credit_open",
			Quantity:     0, // 無効な数量
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order quantity") // 例
	})

	t.Run("異常系: 無効なトリガー価格", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "credit_open",
			Quantity:     100,
			TriggerPrice: 0, // 無効なトリガー価格
			MarketCode:   "00",
		}
		_, err = client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid trigger price")
	})

	t.Run("異常系: APIエラー (ログインせずに注文)", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		// ログインしない

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "credit_open",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		_, err := client.PlaceOrder(context.Background(), order)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request URL not found, need to Login") // エラーメッセージの確認
	})

	t.Run("異常系: コンテキストキャンセル", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil) //ログイン
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "sell",
			OrderType:    "stop",
			Condition:    "credit_open",
			Quantity:     100,
			TriggerPrice: 9000.0,
			MarketCode:   "00",
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // キャンセル

		_, err = client.PlaceOrder(ctx, order) // contextを渡す
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled)) // キャンセルエラー
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order_credit_market_sell_with_stop_test.go

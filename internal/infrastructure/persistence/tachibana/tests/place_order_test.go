// internal/infrastructure/persistence/tachibana/tests/place_order_test.go
package tachibana_test

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder(t *testing.T) {
	// 正常系のテスト(成行注文)
	t.Run("正常系: 成行注文が成功すること", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:    "7974", // 任天堂
			Side:      "buy",
			OrderType: "market",
			Quantity:  100,
		}

		placedOrder, err := client.PlaceOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.NotEmpty(t, placedOrder.ID)             // 注文IDが設定されている
		assert.Equal(t, "pending", placedOrder.Status) // ステータスが"pending"

		// 1秒待機
		time.Sleep(1 * time.Second)
	})

	// 正常系のテスト (指値注文)
	t.Run("正常系: 指値注文が成功すること", func(t *testing.T) {

		// 1秒待機
		// time.Sleep(2 * time.Second)
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:    "7974",
			Side:      "buy",
			OrderType: "limit",
			Quantity:  100,
			Price:     10000.0, // 指値価格
		}

		placedOrder, err := client.PlaceOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.NotEmpty(t, placedOrder.ID)
		assert.Equal(t, "pending", placedOrder.Status)
		// 1秒待機
		time.Sleep(1 * time.Second)
	})

	// 正常系のテスト（逆指値）
	t.Run("正常系: 逆指値注文が成功すること", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, nil)
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)
		defer client.Logout(context.Background())

		order := &domain.Order{
			Symbol:       "7974",
			Side:         "buy",
			OrderType:    "stop",
			Quantity:     100,
			TriggerPrice: 4000.0,  // 逆指値条件価格
			Price:        10500.0, // 逆指値がトリガーされた後の注文価格（指値）
		}

		placedOrder, err := client.PlaceOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.NotNil(t, placedOrder)
		assert.NotEmpty(t, placedOrder.ID)
		assert.Equal(t, "pending", placedOrder.Status)
		// 1秒待機
		time.Sleep(1 * time.Second)
	})

}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/place_order_test.go

package order_test

import (
	"context"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/order"
	"estore-trade/test/docker"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestOrderRepository_UpdateOrder(t *testing.T) {
	db, cleanup, err := docker.SetupTestDatabase() // SetupTestDatabase を呼び出す
	require.NoError(t, err)
	defer cleanup()

	repo := order.NewOrderRepository(db)

	t.Run("正常系: 注文情報を更新できること", func(t *testing.T) {
		// テストデータの準備 (CreateOrder を使う)
		order := &domain.Order{
			UUID:             uuid.NewString(),
			Symbol:           "7203",
			Side:             "long", // long に変更
			OrderType:        "market",
			Quantity:         100,
			Status:           "pending", // 初期ステータス
			TachibanaOrderID: "tachibana-order-id",
		}
		err = repo.CreateOrder(context.Background(), order)
		assert.NoError(t, err)

		// 更新する値
		updatedOrder := &domain.Order{
			UUID:             order.UUID, // 更新対象を特定するために UUID は必須
			Symbol:           "7203",     // Symbolは変更しない
			Side:             "short",    // long/short を使用
			Status:           "filled",   // filled (約定済み) に更新
			FilledQuantity:   100,
			AveragePrice:     750.0,
			TachibanaOrderID: "new-tachibana-id",
			Commission:       10.0,
			ExpireAt:         time.Now().Add(24 * time.Hour),
			//UpdatedAt:        time.Now(), //gormで管理
		}

		// UpdateOrder の呼び出し
		err = repo.UpdateOrder(context.Background(), updatedOrder)
		assert.NoError(t, err)

		// 更新後の注文情報を取得
		retrievedOrder, err := repo.GetOrder(context.Background(), order.UUID)
		assert.NoError(t, err)

		// 更新された値を確認
		if assert.NotNil(t, retrievedOrder) {
			assert.Equal(t, "short", retrievedOrder.Side)    // long/short を使用
			assert.Equal(t, "filled", retrievedOrder.Status) // filled に更新されている
			assert.Equal(t, 100, retrievedOrder.FilledQuantity)
			assert.Equal(t, 750.0, retrievedOrder.AveragePrice)
			assert.Equal(t, "new-tachibana-id", retrievedOrder.TachibanaOrderID)
			assert.Equal(t, 10.0, retrievedOrder.Commission)
		}
	})
	// update_order_test.go

	t.Run("異常系：存在しない注文情報を更新", func(t *testing.T) {
		// 存在しないUUID
		nonExistentUUID := uuid.NewString()

		updatedOrder := &domain.Order{
			UUID:   nonExistentUUID, // 存在しない UUID
			Symbol: "7203",
			Side:   "short", // long/short を使用
			Status: "filled",
		}
		err = repo.UpdateOrder(context.Background(), updatedOrder)
		assert.Error(t, err)                           // 更新に失敗することを期待
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound) // 追加: gorm.ErrRecordNotFound であることを確認
	})
}

// go test -v ./internal/infrastructure/persistence/order/tests/update_order_test.go

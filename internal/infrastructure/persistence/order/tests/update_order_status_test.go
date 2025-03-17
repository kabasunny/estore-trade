package order_test

import (
	"context"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/order"
	"estore-trade/test/docker" // docker パッケージをインポート

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestOrderRepository_UpdateOrderStatus(t *testing.T) {
	db, cleanup, err := docker.SetupTestDatabase() // SetupTestDatabase を呼び出す
	require.NoError(t, err)
	defer cleanup()

	repo := order.NewOrderRepository(db)

	orderID := uuid.NewString()
	newStatus := "filled"

	t.Run("正常系: 注文ステータスを更新できること", func(t *testing.T) {
		// テストデータの準備 (CreateOrder を使う)
		order := &domain.Order{
			UUID:             orderID,
			Symbol:           "7203",
			Side:             "buy",
			OrderType:        "market",
			Quantity:         100,
			Status:           "pending", // 初期ステータス
			TachibanaOrderID: "tachibana-order-id",
		}
		err = repo.CreateOrder(context.Background(), order)
		assert.NoError(t, err)

		// UpdateOrderStatus の呼び出し
		err = repo.UpdateOrderStatus(context.Background(), orderID, newStatus)
		assert.NoError(t, err)

		// 更新後の注文情報を取得
		updatedOrder, err := repo.GetOrder(context.Background(), orderID)
		assert.NoError(t, err)
		if assert.NotNil(t, updatedOrder) {
			assert.Equal(t, newStatus, updatedOrder.Status)
		}
	})

	t.Run("異常系: 存在しない注文のステータスを更新しようとするとエラーになること", func(t *testing.T) {
		// UpdateOrderStatus の呼び出し (存在しない UUID を指定)
		nonExistentUUID := uuid.NewString()
		err = repo.UpdateOrderStatus(context.Background(), nonExistentUUID, newStatus)
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound) //追加
	})
}

// go test -v ./internal/infrastructure/persistence/order/tests/update_order_status_test.go

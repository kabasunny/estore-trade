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
)

func TestOrderRepository_CancelOrder(t *testing.T) {
	db, cleanup, err := docker.SetupTestDatabase() // SetupTestDatabase を呼び出す
	require.NoError(t, err)
	defer cleanup()

	repo := order.NewOrderRepository(db)

	t.Run("正常系: 注文をキャンセルできること", func(t *testing.T) {
		// テストデータの準備 (CreateOrder を使う)
		order := &domain.Order{
			UUID:             uuid.NewString(),
			Symbol:           "7203",
			Side:             "buy",
			OrderType:        "market",
			Quantity:         100,
			Status:           "pending", // 初期ステータス
			TachibanaOrderID: "tachibana-order-id",
		}
		err = repo.CreateOrder(context.Background(), order)
		assert.NoError(t, err)

		// CancelOrder の呼び出し
		err = repo.CancelOrder(context.Background(), order.UUID)
		assert.NoError(t, err)

		// 注文のステータスが "canceled" に更新されていることを確認
		canceledOrder, err := repo.GetOrder(context.Background(), order.UUID)
		assert.NoError(t, err)
		if assert.NotNil(t, canceledOrder) {
			assert.Equal(t, "canceled", canceledOrder.Status)
		}
	})

	t.Run("異常系: 存在しない注文をキャンセルしようとするとエラーになること", func(t *testing.T) {
		// CancelOrder の呼び出し (存在しない UUID を指定)
		nonExistentUUID := uuid.NewString()
		err := repo.CancelOrder(context.Background(), nonExistentUUID)
		assert.Error(t, err) // エラーが発生することを期待
		assert.Contains(t, err.Error(), "order not found")
	})

	// 他の異常系テストケース (必要に応じて追加)
}

// go test -v ./internal/infrastructure/persistence/order/tests/cancel_order_test.go

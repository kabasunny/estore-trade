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

func TestOrderRepository_GetOrder(t *testing.T) {
	db, cleanup, err := docker.SetupTestDatabase() // SetupTestDatabase を呼び出す
	require.NoError(t, err)
	defer cleanup()

	repo := order.NewOrderRepository(db)

	t.Run("正常系: 存在する注文を取得できること", func(t *testing.T) {
		// テストデータの準備 (CreateOrder を使う)
		order := &domain.Order{
			UUID:             uuid.NewString(),
			Symbol:           "7203",
			Side:             "buy",
			OrderType:        "market",
			Quantity:         100,
			Status:           "pending",
			TachibanaOrderID: "tachibana-order-id",
		}
		err = repo.CreateOrder(context.Background(), order)
		assert.NoError(t, err)

		// GetOrder の呼び出し
		retrievedOrder, err := repo.GetOrder(context.Background(), order.UUID)

		// 結果の検証
		assert.NoError(t, err)
		if assert.NotNil(t, retrievedOrder) { // nil チェックを追加
			assert.Equal(t, order.UUID, retrievedOrder.UUID)
			assert.Equal(t, order.Symbol, retrievedOrder.Symbol)
			assert.Equal(t, order.Side, retrievedOrder.Side)
			assert.Equal(t, order.OrderType, retrievedOrder.OrderType)
			assert.Equal(t, order.Quantity, retrievedOrder.Quantity)
			assert.Equal(t, order.Status, retrievedOrder.Status)
			assert.Equal(t, order.TachibanaOrderID, retrievedOrder.TachibanaOrderID)
			// 他のフィールドも必要に応じて比較
		}
	})

	t.Run("異常系: 存在しない注文を取得しようとすると nil, nil が返ること", func(t *testing.T) {
		// GetOrder の呼び出し (存在しない UUID を指定)
		nonExistentUUID := uuid.NewString()
		order, err := repo.GetOrder(context.Background(), nonExistentUUID)

		// 結果の検証
		assert.NoError(t, err) // GORM は ErrRecordNotFound を返さない
		assert.Nil(t, order)
	})
}

// go test -v ./internal/infrastructure/persistence/order/tests/get_order_test.go

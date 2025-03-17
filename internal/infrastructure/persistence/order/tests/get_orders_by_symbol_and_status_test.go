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

func TestOrderRepository_GetOrdersBySymbolAndStatus(t *testing.T) {
	db, cleanup, err := docker.SetupTestDatabase() // SetupTestDatabase を呼び出す
	require.NoError(t, err)
	defer cleanup()

	repo := order.NewOrderRepository(db)

	symbol := "7203"
	status := "pending"

	t.Run("正常系: 指定されたシンボルとステータスの注文が複数取得できること", func(t *testing.T) {
		// テストデータの準備 (CreateOrder を使う)
		expectedOrders := []*domain.Order{
			{UUID: uuid.NewString(), Symbol: symbol, Side: "long", OrderType: "market", Quantity: 100, Status: status},
			{UUID: uuid.NewString(), Symbol: symbol, Side: "short", OrderType: "limit", Price: 1500, Quantity: 50, Status: status},
		}
		for _, order := range expectedOrders {
			err = repo.CreateOrder(context.Background(), order)
			assert.NoError(t, err)
		}

		orders, err := repo.GetOrdersBySymbolAndStatus(context.Background(), symbol, status)
		assert.NoError(t, err)
		assert.Len(t, orders, 2)
		if len(orders) > 1 { //追加
			assert.Equal(t, expectedOrders[0].UUID, orders[0].UUID)
			assert.Equal(t, expectedOrders[1].UUID, orders[1].UUID)
		}

	})

	t.Run("正常系: 指定されたシンボルとステータスの注文が存在しない場合、空のスライスが返ること", func(t *testing.T) {
		orders, err := repo.GetOrdersBySymbolAndStatus(context.Background(), "non-existent-symbol", status)
		assert.NoError(t, err)
		assert.Empty(t, orders) // 空のスライス
	})
}

// go test -v ./internal/infrastructure/persistence/order/tests/get_orders_by_symbol_and_status_test.go

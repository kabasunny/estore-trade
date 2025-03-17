// internal/infrastructure/persistence/order/tests/create_order_test.go
package order_test

import (
	"context"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/order"
	"estore-trade/test/docker"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderRepository_CreateOrder(t *testing.T) {
	db, cleanup, err := docker.SetupTestDatabase()
	require.NoError(t, err)
	defer cleanup()

	repo := order.NewOrderRepository(db)

	order := &domain.Order{
		UUID:             uuid.NewString(),
		Symbol:           "7974",
		Side:             "long",
		OrderType:        "market",
		Quantity:         100,
		Status:           "pending",
		TachibanaOrderID: "tachibana-order-id-123",
	}

	err = repo.CreateOrder(context.Background(), order)
	assert.NoError(t, err)

	// GORM を使って DB からデータを取得
	var got domain.Order
	err = db.First(&got, "uuid = ?", order.UUID).Error // First を使用
	require.NoError(t, err)
	assert.Equal(t, order.UUID, got.UUID)
	assert.Equal(t, order.Symbol, got.Symbol)
	// ... 他のフィールドも検証 ...
}

// go test -v ./internal/infrastructure/persistence/order/tests/create_order_test.go

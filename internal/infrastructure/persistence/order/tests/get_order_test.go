// internal/infrastructure/persistence/order/tests/get_order_test.go
package order

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/order"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_GetOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) // order. を追加
	orderID := "order-1"

	expectedOrder := &domain.Order{
		UUID:             orderID,
		Symbol:           "7203",
		Side:             "buy",
		OrderType:        "market",
		Quantity:         100,
		Status:           "pending",
		TachibanaOrderID: "tachibana-order-id",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "symbol", "order_type", "side", "quantity", "price", "trigger_price", "filled_quantity", "average_price", "status", "tachibana_order_id", "commission", "expire_at", "created_at", "updated_at"}).
		AddRow(expectedOrder.UUID, expectedOrder.Symbol, expectedOrder.OrderType, expectedOrder.Side,
			expectedOrder.Quantity, expectedOrder.Price, expectedOrder.TriggerPrice,
			expectedOrder.FilledQuantity, expectedOrder.AveragePrice, expectedOrder.Status,
			expectedOrder.TachibanaOrderID, expectedOrder.Commission, expectedOrder.ExpireAt,
			expectedOrder.CreatedAt, expectedOrder.UpdatedAt)

	mock.ExpectQuery("^SELECT (.+) FROM orders WHERE id =").WithArgs(orderID).WillReturnRows(rows)

	order, err := repo.GetOrder(context.Background(), orderID)
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, expectedOrder.UUID, order.UUID)
	assert.Equal(t, expectedOrder.Symbol, order.Symbol)

	// 他のフィールドも必要に応じて比較

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOrderRepository_GetOrder_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) // order. を追加

	orderID := "non-existent-order"

	mock.ExpectQuery("^SELECT (.+) FROM orders WHERE id =").WithArgs(orderID).WillReturnError(sql.ErrNoRows)

	order, err := repo.GetOrder(context.Background(), orderID)

	assert.NoError(t, err) // NotFound はエラーではない
	assert.Nil(t, order)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// internal/infrastructure/persistence/order/tests/get_orders_by_symbol_and_status_test.go
package order

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/order"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_GetOrdersBySymbolAndStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) // order.を削除

	symbol := "7203"
	status := "pending"
	expectedOrders := []*domain.Order{
		{UUID: "order-1", Symbol: symbol, Side: "buy", OrderType: "market", Quantity: 100, Status: status},
		{UUID: "order-2", Symbol: symbol, Side: "sell", OrderType: "limit", Price: 1500, Quantity: 50, Status: status},
	}

	rows := sqlmock.NewRows([]string{"id", "symbol", "order_type", "side", "quantity", "price", "trigger_price", "filled_quantity", "average_price", "status", "tachibana_order_id", "commission", "expire_at", "created_at", "updated_at"})
	for _, order := range expectedOrders {
		rows.AddRow(order.UUID, order.Symbol, order.OrderType, order.Side, order.Quantity, order.Price,
			order.TriggerPrice, order.FilledQuantity, order.AveragePrice, order.Status,
			order.TachibanaOrderID, order.Commission, order.ExpireAt, order.CreatedAt, order.UpdatedAt)
	}
	mock.ExpectQuery("^SELECT (.+) FROM orders WHERE symbol = (.+) AND status = (.+)").WithArgs(symbol, status).WillReturnRows(rows)

	orders, err := repo.GetOrdersBySymbolAndStatus(context.Background(), symbol, status)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, expectedOrders[0].UUID, orders[0].UUID)
	assert.Equal(t, expectedOrders[1].UUID, orders[1].UUID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

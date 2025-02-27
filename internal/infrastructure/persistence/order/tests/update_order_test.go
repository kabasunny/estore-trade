// internal/infrastructure/persistence/order/tests/update_order_test.go
package order

import (
	"context"
	"errors"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/order"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_UpdateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) //order.を削除

	order := &domain.Order{
		ID:               "order-1",
		Symbol:           "7203",
		Side:             "sell", // 更新
		Status:           "filled",
		FilledQuantity:   100,
		AveragePrice:     750.0,
		TachibanaOrderID: "new-tachibana-id",
		Commission:       10.0,
		ExpireAt:         time.Now().Add(24 * time.Hour),
		UpdatedAt:        time.Now(),
	}
	// ExecContextが呼ばれ、エラーがnilであることをモックで定義
	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE orders
        SET symbol = $2, order_type = $3, side = $4, quantity = $5, price = $6, trigger_price = $7, filled_quantity = $8, average_price = $9, status = $10, tachibana_order_id = $11, commission = $12, expire_at = $13, updated_at = $14
        WHERE id = $1
    `)).WithArgs(
		order.ID,
		order.Symbol,
		order.OrderType,
		order.Side,
		order.Quantity,
		order.Price,
		order.TriggerPrice,
		order.FilledQuantity,
		order.AveragePrice,
		order.Status,
		order.TachibanaOrderID,
		order.Commission,
		order.ExpireAt,
		sqlmock.AnyArg(), // UpdatedAt
	).WillReturnResult(sqlmock.NewResult(0, 1)) // 0, 1 -> LastInsertId, RowsAffected

	err = repo.UpdateOrder(context.Background(), order)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestOrderRepository_UpdateOrder_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) //order.を削除

	order := &domain.Order{
		ID:               "order-1",
		Symbol:           "7203",
		Side:             "sell", // 更新
		Status:           "filled",
		FilledQuantity:   100,
		AveragePrice:     750.0,
		TachibanaOrderID: "new-tachibana-id",
		Commission:       10.0,
		ExpireAt:         time.Now().Add(24 * time.Hour),
		UpdatedAt:        time.Now(),
	}

	// 期待されるエラー
	expectedError := errors.New("database error")

	// ExecContext がエラーを返すようにモック
	mock.ExpectExec(regexp.QuoteMeta(`
	UPDATE orders
	SET symbol = $2, order_type = $3, side = $4, quantity = $5, price = $6, trigger_price = $7, filled_quantity = $8, average_price = $9, status = $10, tachibana_order_id = $11, commission = $12, expire_at = $13, updated_at = $14
	WHERE id = $1
	`)).WithArgs(
		order.ID,
		order.Symbol,
		order.OrderType,
		order.Side,
		order.Quantity,
		order.Price,
		order.TriggerPrice,
		order.FilledQuantity,
		order.AveragePrice,
		order.Status,
		order.TachibanaOrderID,
		order.Commission,
		order.ExpireAt,
		sqlmock.AnyArg(), // UpdatedAt
	).WillReturnError(expectedError)

	err = repo.UpdateOrder(context.Background(), order)
	assert.Error(t, err) // エラーが発生することを期待
	//assert.Equal(t, expectedError, err) // エラーが期待通りか確認

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

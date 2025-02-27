// internal/infrastructure/persistence/order/tests/cancel_order_test.go

package order

import (
	"context"
	"errors"
	"estore-trade/internal/infrastructure/persistence/order"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_CancelOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	repo := order.NewOrderRepository(db) //order.を追加
	orderID := "order-1"

	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `)).WithArgs(orderID, "canceled", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.CancelOrder(context.Background(), orderID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestOrderRepository_CancelOrder_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) //order.を追加

	orderID := "non-existent-order"

	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `)).WithArgs(orderID, "canceled", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.CancelOrder(context.Background(), orderID)
	assert.Error(t, err) // NotFoundはエラー

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestOrderRepository_CancelOrder_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) //order.を追加

	orderID := "order-1"
	expectedError := errors.New("database error")
	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `)).WithArgs(orderID, "canceled", sqlmock.AnyArg()).WillReturnError(expectedError)

	err = repo.CancelOrder(context.Background(), orderID)
	assert.Error(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

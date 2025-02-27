// internal/infrastructure/persistence/order/tests/update_order_status_test.go
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

func TestOrderRepository_UpdateOrderStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) //order.を削除

	orderID := "order-1"
	newStatus := "filled"

	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `)).WithArgs(orderID, newStatus, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateOrderStatus(context.Background(), orderID, newStatus)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOrderRepository_UpdateOrderStatus_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) //order.を削除

	orderID := "non-existent-order"
	newStatus := "filled"

	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `)).WithArgs(orderID, newStatus, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateOrderStatus(context.Background(), orderID, newStatus)
	assert.Error(t, err) // NotFoundはエラー

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOrderRepository_UpdateOrderStatus_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db) //order.を削除

	orderID := "order-1"
	newStatus := "filled"
	expectedError := errors.New("database error")

	mock.ExpectExec(regexp.QuoteMeta(`
    UPDATE orders
    SET status = $2, updated_at = $3
    WHERE id = $1
    `)).WithArgs(orderID, newStatus, sqlmock.AnyArg()).WillReturnError(expectedError)

	err = repo.UpdateOrderStatus(context.Background(), orderID, newStatus)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// internal/infrastructure/persistence/order/tests/create_order_test.go
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

// MockDB は sql.DB のモック(testify/mock は不要)
type MockDB struct {
	sqlmock.Sqlmock
}

func TestOrderRepository_CreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New() // モックのDBとモックオブジェクトを取得
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db)

	order := &domain.Order{
		UUID:             "order-1",
		Symbol:           "7203",
		Side:             "buy",
		OrderType:        "market",
		Quantity:         100,
		Status:           "pending",
		TachibanaOrderID: "tachibana-order-id",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// クエリと期待される引数をモックに設定
	// SQL クエリを order パッケージのものに修正
	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO orders (id, symbol, order_type, side, quantity, price, trigger_price, filled_quantity, average_price, status, tachibana_order_id, commission, expire_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `)).WithArgs(
		order.UUID,
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
		sqlmock.AnyArg(), // CreatedAt
		sqlmock.AnyArg(), // UpdatedAt
	).WillReturnResult(sqlmock.NewResult(1, 1)) // 挿入されたID, 影響行数

	err = repo.CreateOrder(context.Background(), order)
	assert.NoError(t, err)

	// モックの設定がすべて満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestOrderRepository_CreateOrder_Error(t *testing.T) {
	db, mock, err := sqlmock.New() // sqlmock のモックを使用
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := order.NewOrderRepository(db)

	order := &domain.Order{
		UUID:             "order-1",
		Symbol:           "7203",
		Side:             "buy",
		OrderType:        "market",
		Quantity:         100,
		Status:           "pending",
		TachibanaOrderID: "tachibana-order-id",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	expectedError := errors.New("database error")

	// SQL クエリを order パッケージのものに修正
	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO orders (id, symbol, order_type, side, quantity, price, trigger_price, filled_quantity, average_price, status, tachibana_order_id, commission, expire_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `)).WithArgs(
		order.UUID,
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
		sqlmock.AnyArg(), // CreatedAt
		sqlmock.AnyArg(), // UpdatedAt
	).WillReturnError(expectedError)

	err = repo.CreateOrder(context.Background(), order)
	assert.Error(t, err)
	// エラーメッセージの比較ではなく、エラーが発生したことだけを確認
	// (sqlmock のエラーメッセージは詳細で、変更される可能性があるため)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

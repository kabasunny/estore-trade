// get_by_user_id_test.go
package account

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/account"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_GetAccountByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := account.NewAccountRepository(db)

	userID := "user1"
	expectedPositions := []domain.Position{
		{Symbol: "7203", Quantity: 100, Price: 1500, Side: "buy"},
		{Symbol: "8306", Quantity: 200, Price: 750, Side: "sell"},
	}
	expectedAccount := &domain.Account{
		ID:               "test-id",
		UserID:           userID,
		AccountType:      "special",
		Balance:          10000,
		AvailableBalance: 10000,
		Margin:           0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Positions:        expectedPositions,
	}

	// GetAccountByUserID のクエリに対する期待値を設定
	rows := sqlmock.NewRows([]string{"id", "user_id", "account_type", "balance", "available_balance", "margin", "created_at", "updated_at"}).
		AddRow(expectedAccount.ID, expectedAccount.UserID, expectedAccount.AccountType, expectedAccount.Balance,
			expectedAccount.AvailableBalance, expectedAccount.Margin, expectedAccount.CreatedAt, expectedAccount.UpdatedAt)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE user_id =").WithArgs(userID).WillReturnRows(rows)

	// getPositions のクエリに対する期待値を設定
	positionRows := sqlmock.NewRows([]string{"symbol", "quantity", "price", "side"})
	for _, p := range expectedPositions {
		positionRows.AddRow(p.Symbol, p.Quantity, p.Price, p.Side)
	}
	mock.ExpectQuery("^SELECT (.+) FROM positions WHERE account_id =").WithArgs(expectedAccount.ID).WillReturnRows(positionRows)

	// テスト実行
	account, err := repo.GetAccountByUserID(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, expectedAccount, account)

	// モックの設定がすべて満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_GetAccountByUserID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := account.NewAccountRepository(db)

	userID := "non-existent-user"

	// sql.ErrNoRows を返すようにモック
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE user_id =").WithArgs(userID).WillReturnError(sql.ErrNoRows)

	account, err := repo.GetAccountByUserID(context.Background(), userID)

	assert.NoError(t, err) // ErrNoRows はエラーとして扱わない
	assert.Nil(t, account)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

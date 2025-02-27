// internal/infrastructure/persistence/account/tests/get_account_test.go
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

func TestAccountRepository_GetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := account.NewAccountRepository(db)

	accountID := "test-id"
	expectedPositions := []domain.Position{
		{Symbol: "7203", Quantity: 100, Price: 1500, Side: "buy"},
		{Symbol: "8306", Quantity: 200, Price: 750, Side: "sell"},
	}
	expectedAccount := &domain.Account{
		ID:               accountID,
		UserID:           "user1",
		AccountType:      "special",
		Balance:          10000,
		AvailableBalance: 10000,
		Margin:           0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Positions:        expectedPositions,
	}

	// GetAccount のクエリに対する期待値を設定
	rows := sqlmock.NewRows([]string{"id", "user_id", "account_type", "balance", "available_balance", "margin", "created_at", "updated_at"}).
		AddRow(expectedAccount.ID, expectedAccount.UserID, expectedAccount.AccountType, expectedAccount.Balance,
			expectedAccount.AvailableBalance, expectedAccount.Margin, expectedAccount.CreatedAt, expectedAccount.UpdatedAt)
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE id =").WithArgs(accountID).WillReturnRows(rows)

	// getPositions のクエリに対する期待値を設定
	positionRows := sqlmock.NewRows([]string{"symbol", "quantity", "price", "side"})
	for _, p := range expectedPositions {
		positionRows.AddRow(p.Symbol, p.Quantity, p.Price, p.Side)
	}
	mock.ExpectQuery("^SELECT (.+) FROM positions WHERE account_id =").WithArgs(accountID).WillReturnRows(positionRows)

	// テスト実行
	account, err := repo.GetAccount(context.Background(), accountID)

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, expectedAccount, account)

	// モックの設定がすべて満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAccountRepository_GetAccount_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := account.NewAccountRepository(db)

	accountID := "non-existent-id"

	// sql.ErrNoRows を返すようにモック
	mock.ExpectQuery("^SELECT (.+) FROM accounts WHERE id =").WithArgs(accountID).WillReturnError(sql.ErrNoRows)

	account, err := repo.GetAccount(context.Background(), accountID)

	assert.NoError(t, err) // ErrNoRows はエラーとして扱わない
	assert.Nil(t, account)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

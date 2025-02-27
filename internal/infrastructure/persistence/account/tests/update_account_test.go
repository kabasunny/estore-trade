// internal/infrastructure/persistence/account/tests/update_account_test.go
package account

import (
	"context"
	"errors"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/account"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

func TestAccountRepository_UpdateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := account.NewAccountRepository(db)

	account := &domain.Account{
		ID:               "test-id",
		UserID:           "user1",
		AccountType:      "special",
		Balance:          12000, // 更新後の値
		AvailableBalance: 11000,
		Margin:           0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	// ExecContextが呼ばれ、エラーがnilであることをモックで定義
	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE accounts
        SET user_id = $2, account_type = $3, balance = $4, available_balance = $5, margin = $6, updated_at = $7
        WHERE id = $1
    `)).WithArgs(
		account.ID,
		account.UserID,
		account.AccountType,
		account.Balance,
		account.AvailableBalance,
		account.Margin,
		sqlmock.AnyArg(), // AnyTime{} を sqlmock.AnyArg() に変更
	).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateAccount(context.Background(), account)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestAccountRepository_UpdateAccount_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := account.NewAccountRepository(db)

	account := &domain.Account{
		ID:               "test-id",
		UserID:           "user1",
		AccountType:      "special",
		Balance:          12000,
		AvailableBalance: 10000,
		Margin:           0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 期待されるエラー
	expectedError := errors.New("database error")

	// ExecContext がエラーを返すようにモック
	mock.ExpectExec(regexp.QuoteMeta(`
	UPDATE accounts
	SET user_id = $2, account_type = $3, balance = $4, available_balance = $5, margin = $6, updated_at = $7
	WHERE id = $1
	`)).WithArgs(
		account.ID,
		account.UserID,
		account.AccountType,
		account.Balance,
		account.AvailableBalance,
		account.Margin,
		sqlmock.AnyArg(), // AnyTime{} を sqlmock.AnyArg() に変更
	).WillReturnError(expectedError)

	err = repo.UpdateAccount(context.Background(), account)
	assert.Error(t, err) // エラーが発生することを期待
	//assert.Equal(t, expectedError, err) // エラーが期待通りか確認

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
